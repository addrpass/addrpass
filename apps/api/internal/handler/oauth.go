package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type OAuthHandler struct {
	oauth     *service.OAuthService
	addresses *service.AddressService
}

func NewOAuthHandler(oauth *service.OAuthService, addresses *service.AddressService) *OAuthHandler {
	return &OAuthHandler{oauth: oauth, addresses: addresses}
}

// CreateApp — registers an OAuth app for a business (protected)
func (h *OAuthHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	businessID := chi.URLParam(r, "businessId")

	var req model.CreateOAuthAppRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || len(req.RedirectURIs) == 0 {
		writeError(w, http.StatusBadRequest, "name and redirect_uris are required")
		return
	}

	resp, err := h.oauth.CreateApp(r.Context(), userID, businessID, req)
	if err != nil {
		if errors.Is(err, service.ErrBusinessNotFound) {
			writeError(w, http.StatusNotFound, "business not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create OAuth app")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// Authorize — returns consent page data (app info + user's addresses).
// The frontend renders the consent screen using this data.
// GET /api/v1/oauth/authorize?client_id=X&redirect_uri=Y&scope=Z&state=W
func (h *OAuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	if clientID == "" || redirectURI == "" {
		writeError(w, http.StatusBadRequest, "client_id and redirect_uri are required")
		return
	}
	if scope == "" {
		scope = "full"
	}

	app, err := h.oauth.GetAppByClientID(r.Context(), clientID)
	if err != nil {
		if errors.Is(err, service.ErrOAuthAppNotFound) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get application")
		return
	}

	if err := h.oauth.ValidateRedirectURI(app, redirectURI); err != nil {
		writeError(w, http.StatusBadRequest, "invalid redirect_uri")
		return
	}

	addresses, err := h.addresses.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list addresses")
		return
	}

	writeJSON(w, http.StatusOK, model.ConsentPageData{
		App:       *app,
		Addresses: addresses,
		Scope:     scope,
		State:     state,
	})
}

// Consent — user approves sharing an address. Creates an auth code and returns redirect URL.
// POST /api/v1/oauth/consent
func (h *OAuthHandler) Consent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.ConsentRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ClientID == "" || req.RedirectURI == "" || req.AddressID == "" {
		writeError(w, http.StatusBadRequest, "client_id, redirect_uri, and address_id are required")
		return
	}

	code, err := h.oauth.CreateAuthorizationCode(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOAuthAppNotFound):
			writeError(w, http.StatusNotFound, "application not found")
		case errors.Is(err, service.ErrInvalidRedirectURI):
			writeError(w, http.StatusBadRequest, "invalid redirect_uri")
		default:
			writeError(w, http.StatusInternalServerError, "failed to create authorization code")
		}
		return
	}

	// Build redirect URL with code
	redirectURL := req.RedirectURI + "?code=" + code
	if req.State != "" {
		redirectURL += "&state=" + req.State
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"redirect_url": redirectURL,
		"code":         code,
	})
}

// Exchange — exchanges authorization code for access token + share token.
// POST /api/v1/oauth/token (with grant_type=authorization_code)
func (h *OAuthHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	var req model.TokenExchangeRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.GrantType != "authorization_code" {
		writeError(w, http.StatusBadRequest, "grant_type must be authorization_code")
		return
	}

	resp, err := h.oauth.ExchangeCode(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAuthCodeNotFound):
			writeError(w, http.StatusBadRequest, "invalid authorization code")
		case errors.Is(err, service.ErrAuthCodeExpired):
			writeError(w, http.StatusBadRequest, "authorization code expired")
		case errors.Is(err, service.ErrAuthCodeUsed):
			writeError(w, http.StatusBadRequest, "authorization code already used")
		case errors.Is(err, service.ErrInvalidClient):
			writeError(w, http.StatusUnauthorized, "invalid client credentials")
		case errors.Is(err, service.ErrClientMismatch):
			writeError(w, http.StatusBadRequest, "client_id mismatch")
		case errors.Is(err, service.ErrRedirectURIMismatch):
			writeError(w, http.StatusBadRequest, "redirect_uri mismatch")
		default:
			writeError(w, http.StatusInternalServerError, "failed to exchange code")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
