// Package public provides HTTP handlers for guest-facing endpoints including
// fetching invite details, viewing user details, and requesting QR code delivery.
package public

import (
	"context"
	"encoding/json"
	"invite_qr/internal/server"
	"net/http"
	"strconv"

	db_gen "invite_qr/db/db_gen"

	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

// publicService defines the contract for guest-facing participant operations.
// The production implementation is *Service; tests can provide a mock.
type publicService interface {
	// GetParticipantByExternalID looks up a participant by their external UUID.
	GetParticipantByExternalID(ctx context.Context, externalID string) (*db_gen.Participant, error)
}

// Handler holds the dependencies needed by the public HTTP handlers.
type Handler struct {
	service publicService
}

// NewHandler creates a Handler backed by the given service.
func NewHandler(service publicService) *Handler {
	return &Handler{service: service}
}

// HandleGetInvite returns an HTTP handler that reads the participant's external
// UUID from the URL path, looks up the participant, and returns their details as JSON.
func (h *Handler) HandleGetInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		if token == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}
		participant, err := h.service.GetParticipantByExternalID(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(participant)
	}
}

// GetUserDetails returns an HTTP handler that reads an "id" query parameter,
// looks up the participant, and returns their details as JSON.
func (h *Handler) GetUserDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := server.LoggerFromContext(r.Context())
		id := r.URL.Query().Get("id")
		logger.Info("user details requested", zap.String("id", id))
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		participant, err := h.service.GetParticipantByExternalID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(participant)
	}
}

// SendQRCode returns an HTTP handler that accepts a participant_id query parameter,
// generates a QR code PNG containing the external UUID, and serves it as an image.
func (h *Handler) SendQRCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := server.LoggerFromContext(r.Context())
		id := r.URL.Query().Get("participant_id")
		if id == "" {
			http.Error(w, "missing participant_id", http.StatusBadRequest)
			return
		}

		logger.Info("generating qr code", zap.String("participant_id", id))

		qr, err := qrcode.New(id, qrcode.Medium)
		if err != nil {
			logger.Error("failed to generate qr code", zap.Error(err))
			http.Error(w, "failed to generate qr code", http.StatusInternalServerError)
			return
		}

		png, err := qr.PNG(256)
		if err != nil {
			logger.Error("failed to render qr code as png", zap.Error(err))
			http.Error(w, "failed to render qr code", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(png)))
		w.Write(png)
	}
}
