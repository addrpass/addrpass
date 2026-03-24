package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type ShareHandler struct {
	shares    *service.ShareService
	webhooks  *service.WebhookService
	baseURL   string
	jwtSecret string
}

func NewShareHandler(shares *service.ShareService, webhooks *service.WebhookService, baseURL, jwtSecret string) *ShareHandler {
	return &ShareHandler{shares: shares, webhooks: webhooks, baseURL: baseURL, jwtSecret: jwtSecret}
}

func (h *ShareHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.CreateShareRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AddressID == "" {
		writeError(w, http.StatusBadRequest, "address_id is required")
		return
	}

	share, err := h.shares.Create(r.Context(), userID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create share")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"share": share,
		"url":   "https://addrpass.com/resolve?t=" + share.Token,
	})
}

func (h *ShareHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	shares, err := h.shares.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list shares")
		return
	}

	writeJSON(w, http.StatusOK, shares)
}

func (h *ShareHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	shareID := chi.URLParam(r, "id")

	err := h.shares.Revoke(r.Context(), userID, shareID)
	if err != nil {
		if errors.Is(err, service.ErrShareNotFound) {
			writeError(w, http.StatusNotFound, "share not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to revoke share")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

func (h *ShareHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	shareID := chi.URLParam(r, "id")

	err := h.shares.Delete(r.Context(), userID, shareID)
	if err != nil {
		if errors.Is(err, service.ErrShareNotFound) {
			writeError(w, http.StatusNotFound, "share not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete share")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShareHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	pin := r.URL.Query().Get("pin")
	ip := r.RemoteAddr
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ip = fwd
	}
	userAgent := r.Header.Get("User-Agent")

	// Check if caller is an authenticated business
	businessID, businessName := h.extractBusinessIdentity(r)

	result, err := h.shares.Resolve(r.Context(), token, pin, ip, userAgent, businessID, businessName)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShareNotFound):
			writeError(w, http.StatusNotFound, "share not found")
		case errors.Is(err, service.ErrShareExpired):
			writeError(w, http.StatusGone, "share has expired")
		case errors.Is(err, service.ErrShareRevoked):
			writeError(w, http.StatusGone, "share has been revoked")
		case errors.Is(err, service.ErrMaxAccesses):
			writeError(w, http.StatusGone, "maximum accesses reached")
		case errors.Is(err, service.ErrInvalidPin):
			writeError(w, http.StatusForbidden, "invalid pin")
		default:
			writeError(w, http.StatusInternalServerError, "failed to resolve share")
		}
		return
	}

	// Dispatch webhook asynchronously
	if h.webhooks != nil {
		go h.webhooks.DispatchAccessEvent(result.OwnerID, model.AccessEvent{
			ShareID:      result.ShareID,
			Token:        result.Token,
			IP:           ip,
			UserAgent:    userAgent,
			BusinessName: businessName,
			Scope:        string(result.Scope),
		})
	}

	writeJSON(w, http.StatusOK, model.ResolveResponse{
		Address: result.Address,
		Scope:   string(result.Scope),
	})
}

// extractBusinessIdentity checks if the request has a business Bearer token
// and returns the business ID and name. Returns empty strings for non-business callers.
func (h *ShareHandler) extractBusinessIdentity(r *http.Request) (string, string) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ""
	}

	token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ""
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "business" {
		return "", ""
	}

	bizID, _ := claims["sub"].(string)
	bizName, _ := claims["business_name"].(string)
	return bizID, bizName
}

func (h *ShareHandler) QRCode(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	url := "https://addrpass.com/resolve?t=" + token

	png, err := qrcode.Encode(url, qrcode.Medium, 512)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate QR code")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(png)
}

func (h *ShareHandler) AccessLogs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	shareID := chi.URLParam(r, "id")

	logs, err := h.shares.GetAccessLogs(r.Context(), userID, shareID)
	if err != nil {
		if errors.Is(err, service.ErrShareNotFound) {
			writeError(w, http.StatusNotFound, "share not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get access logs")
		return
	}

	writeJSON(w, http.StatusOK, logs)
}
