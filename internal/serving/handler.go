// handle to serve the frontend dynamic routing
package serving

import (
	"context"
	"encoding/json"
	"invite_qr/cmd"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type Handler struct {
	service *Service
}

func (h *Handler) GetUserDetails(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := cmd.LoggerFromContext(ctx)
		logger.Info("User details are requestd: ", zap.String("id", r.URL.Query().Get("id")))
		id := strings.Split(r.URL.Query().Get("id"), "")[0]
		user, err := h.service.GetParticipantName(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = user
	}
}

func (h *Handler) MarkAttendance(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := cmd.LoggerFromContext(ctx)
		logger.Info("Mark attendance request received")
		var req struct {
			ParticipantID string `json:"participant_id"`
			Email         string `json:"email"`
			WaNumber      string `json:"wa_number"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := h.service.UpdateParticipantAccessed(ctx, req.ParticipantID, req.Email, req.WaNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Failed to update participant accessed", zap.Error(err))
			return
		}
		logger.Info("Participant accessed updated", zap.String("participant_id", req.ParticipantID))
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) SendQRCode(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := cmd.LoggerFromContext(ctx)
		logger.Info("Send QR code request received")

		w.WriteHeader(http.StatusOK)
	}
}
