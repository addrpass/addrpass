package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type ShareHandler struct {
	shares  *service.ShareService
	baseURL string
}

func NewShareHandler(shares *service.ShareService, baseURL string) *ShareHandler {
	return &ShareHandler{shares: shares, baseURL: baseURL}
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

	addr, err := h.shares.Resolve(r.Context(), token, pin, ip, userAgent)
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

	writeJSON(w, http.StatusOK, model.ResolveResponse{Address: *addr})
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
