package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type BusinessHandler struct {
	businesses *service.BusinessService
}

func NewBusinessHandler(businesses *service.BusinessService) *BusinessHandler {
	return &BusinessHandler{businesses: businesses}
}

func (h *BusinessHandler) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.CreateBusinessRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	biz, err := h.businesses.CreateBusiness(r.Context(), userID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create business")
		return
	}

	writeJSON(w, http.StatusCreated, biz)
}

func (h *BusinessHandler) ListBusinesses(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	businesses, err := h.businesses.ListBusinesses(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list businesses")
		return
	}

	writeJSON(w, http.StatusOK, businesses)
}

func (h *BusinessHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	businessID := chi.URLParam(r, "businessId")

	var req model.CreateAPIKeyRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	resp, err := h.businesses.CreateAPIKey(r.Context(), userID, businessID, req)
	if err != nil {
		if errors.Is(err, service.ErrBusinessNotFound) {
			writeError(w, http.StatusNotFound, "business not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create API key")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *BusinessHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	businessID := chi.URLParam(r, "businessId")

	keys, err := h.businesses.ListAPIKeys(r.Context(), userID, businessID)
	if err != nil {
		if errors.Is(err, service.ErrBusinessNotFound) {
			writeError(w, http.StatusNotFound, "business not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to list API keys")
		return
	}

	writeJSON(w, http.StatusOK, keys)
}

func (h *BusinessHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	keyID := chi.URLParam(r, "keyId")

	err := h.businesses.RevokeAPIKey(r.Context(), userID, keyID)
	if err != nil {
		if errors.Is(err, service.ErrAPIKeyNotFound) {
			writeError(w, http.StatusNotFound, "API key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to revoke API key")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

func (h *BusinessHandler) OAuthToken(w http.ResponseWriter, r *http.Request) {
	var req model.OAuthTokenRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.businesses.OAuthToken(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidGrant):
			writeError(w, http.StatusBadRequest, "unsupported grant_type, use client_credentials")
		case errors.Is(err, service.ErrInvalidClient):
			writeError(w, http.StatusUnauthorized, "invalid client credentials")
		default:
			writeError(w, http.StatusInternalServerError, "failed to generate token")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
