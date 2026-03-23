package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type AddressHandler struct {
	addresses *service.AddressService
}

func NewAddressHandler(addresses *service.AddressService) *AddressHandler {
	return &AddressHandler{addresses: addresses}
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.CreateAddressRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Line1 == "" || req.City == "" || req.PostCode == "" || req.Country == "" {
		writeError(w, http.StatusBadRequest, "line1, city, post_code, and country are required")
		return
	}

	addr, err := h.addresses.Create(r.Context(), userID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create address")
		return
	}

	writeJSON(w, http.StatusCreated, addr)
}

func (h *AddressHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	addresses, err := h.addresses.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list addresses")
		return
	}

	writeJSON(w, http.StatusOK, addresses)
}

func (h *AddressHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	addressID := chi.URLParam(r, "id")

	addr, err := h.addresses.Get(r.Context(), userID, addressID)
	if err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			writeError(w, http.StatusNotFound, "address not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get address")
		return
	}

	writeJSON(w, http.StatusOK, addr)
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	addressID := chi.URLParam(r, "id")

	var req model.UpdateAddressRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	addr, err := h.addresses.Update(r.Context(), userID, addressID, req)
	if err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			writeError(w, http.StatusNotFound, "address not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update address")
		return
	}

	writeJSON(w, http.StatusOK, addr)
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	addressID := chi.URLParam(r, "id")

	err := h.addresses.Delete(r.Context(), userID, addressID)
	if err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			writeError(w, http.StatusNotFound, "address not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete address")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
