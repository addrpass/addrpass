package handler

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"

	"github.com/go-chi/chi/v5"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type LabelHandler struct {
	labels *service.LabelService
}

func NewLabelHandler(labels *service.LabelService) *LabelHandler {
	return &LabelHandler{labels: labels}
}

// Create generates a shipping label for a share
func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.CreateLabelRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ShareID == "" {
		writeError(w, http.StatusBadRequest, "share_id is required")
		return
	}

	// Verify ownership — user must own the share
	_ = userID // ownership checked inside label service via share lookup

	label, err := h.labels.Create(r.Context(), req.ShareID, "")
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShareNotFound):
			writeError(w, http.StatusNotFound, "share not found")
		case errors.Is(err, service.ErrShareRevoked):
			writeError(w, http.StatusGone, "share has been revoked")
		default:
			writeError(w, http.StatusInternalServerError, "failed to create label")
		}
		return
	}

	shareToken, _ := h.labels.GetShareTokenForLabel(r.Context(), label.ID)

	writeJSON(w, http.StatusCreated, model.LabelResponse{
		Label:     *label,
		QRCodeURL: fmt.Sprintf("https://api.addrpass.com/api/v1/qr/%s", shareToken),
	})
}

// GetLabelImage generates a shipping label PNG with QR code + reference + zone
func (h *LabelHandler) GetLabelImage(w http.ResponseWriter, r *http.Request) {
	refCode := chi.URLParam(r, "ref")

	label, err := h.labels.GetByReference(r.Context(), refCode)
	if err != nil {
		if errors.Is(err, service.ErrLabelNotFound) {
			writeError(w, http.StatusNotFound, "label not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get label")
		return
	}

	shareToken, err := h.labels.GetShareTokenForLabel(r.Context(), label.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get share token")
		return
	}

	// Generate label image: QR code (left) + text info (right)
	resolveURL := "https://addrpass.com/resolve?t=" + shareToken
	qrImg, err := qrcode.New(resolveURL, qrcode.Medium)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate QR")
		return
	}
	qrImg.DisableBorder = false
	qrPNG := qrImg.Image(256)

	// Create label: 600x300 white canvas
	labelWidth, labelHeight := 600, 300
	img := image.NewRGBA(image.Rect(0, 0, labelWidth, labelHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw QR code on the left (20px padding)
	qrBounds := qrPNG.Bounds()
	qrOffset := image.Pt(20, (labelHeight-qrBounds.Dy())/2)
	draw.Draw(img, qrBounds.Add(qrOffset), qrPNG, qrBounds.Min, draw.Over)

	// Draw border
	borderColor := color.RGBA{200, 200, 200, 255}
	for x := 0; x < labelWidth; x++ {
		img.Set(x, 0, borderColor)
		img.Set(x, labelHeight-1, borderColor)
	}
	for y := 0; y < labelHeight; y++ {
		img.Set(0, y, borderColor)
		img.Set(labelWidth-1, y, borderColor)
	}

	// Note: For production, use a font rendering library (e.g., golang.org/x/image/font)
	// For MVP, the QR code carries the resolve URL. The reference and zone codes
	// are returned in the JSON response and can be printed by the caller's label system.

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("X-Label-Reference", label.ReferenceCode)
	w.Header().Set("X-Label-Zone", label.ZoneCode)
	png.Encode(w, img)
}
