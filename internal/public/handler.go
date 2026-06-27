package public

import (
	"context"
	"encoding/json"
	"invite_qr/internal/server"
	"net/http"

	db_gen "invite_qr/db/db_gen"

	"go.uber.org/zap"
)

type publicService interface {
	GetParticipantByExternalID(ctx context.Context, externalID string) (*db_gen.Participant, error)
	UpdateParticipantAccessed(ctx context.Context, externalID string, email string, waNumber string) error
}

type Handler struct {
	service publicService
}

func NewHandler(service publicService) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) MarkAttendance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := server.LoggerFromContext(r.Context())
		logger.Info("mark attendance request received")
		var req struct {
			ParticipantID string `json:"participant_id"`
			Email         string `json:"email"`
			WaNumber      string `json:"wa_number"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := h.service.UpdateParticipantAccessed(r.Context(), req.ParticipantID, req.Email, req.WaNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("failed to update participant accessed", zap.Error(err))
			return
		}
		logger.Info("participant accessed updated", zap.String("participant_id", req.ParticipantID))
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) SendQRCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := server.LoggerFromContext(r.Context())
		logger.Info("send qr code request received")

		w.WriteHeader(http.StatusOK)
	}
}
