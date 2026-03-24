package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type DelegationHandler struct {
	delegations *service.DelegationService
}

func NewDelegationHandler(delegations *service.DelegationService) *DelegationHandler {
	return &DelegationHandler{delegations: delegations}
}

// CreateByUser — share owner delegates to a business
func (h *DelegationHandler) CreateByUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.CreateDelegationRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ShareID == "" || req.ToBusinessID == "" {
		writeError(w, http.StatusBadRequest, "share_id and to_business_id are required")
		return
	}

	d, err := h.delegations.CreateByUser(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShareNotFound):
			writeError(w, http.StatusNotFound, "share not found")
		case errors.Is(err, service.ErrCannotDelegateUp):
			writeError(w, http.StatusBadRequest, "cannot delegate higher scope than the share")
		default:
			writeError(w, http.StatusInternalServerError, "failed to create delegation")
		}
		return
	}

	writeJSON(w, http.StatusCreated, d)
}

// ListForShare — share owner views all delegations
func (h *DelegationHandler) ListForShare(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	shareID := chi.URLParam(r, "shareId")

	delegations, err := h.delegations.ListForShare(r.Context(), userID, shareID)
	if err != nil {
		if errors.Is(err, service.ErrShareNotFound) {
			writeError(w, http.StatusNotFound, "share not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to list delegations")
		return
	}

	writeJSON(w, http.StatusOK, delegations)
}

// Revoke — share owner revokes a delegation
func (h *DelegationHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	delegationID := chi.URLParam(r, "id")

	err := h.delegations.Revoke(r.Context(), userID, delegationID)
	if err != nil {
		if errors.Is(err, service.ErrDelegationNotFound) {
			writeError(w, http.StatusNotFound, "delegation not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to revoke delegation")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}
