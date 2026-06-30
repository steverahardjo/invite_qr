package public

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"invite_qr/internal/server"

	db_gen "invite_qr/db/db_gen"

	"github.com/danielgtaylor/huma/v2"
	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

type publicService interface {
	GetParticipantByExternalID(ctx context.Context, externalID string) (*db_gen.Participant, error)
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

func (h *Handler) RegisterHumaRoutes(api huma.API) {
	huma.Get(api, "/api/invite/{token}", func(ctx context.Context, input *struct {
		Token string `path:"token"`
	}) (*struct {
		Body *db_gen.Participant
	}, error) {
		if input.Token == "" {
			return nil, huma.Error400BadRequest("missing token")
		}
		participant, err := h.service.GetParticipantByExternalID(ctx, input.Token)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get participant", err)
		}
		return &struct{ Body *db_gen.Participant }{Body: participant}, nil
	})

	huma.Get(api, "/api/user", func(ctx context.Context, input *struct {
		ID string `query:"id"`
	}) (*struct {
		Body *db_gen.Participant
	}, error) {
		logger := server.LoggerFromContext(ctx)
		logger.Info("user details requested", zap.String("id", input.ID))
		if input.ID == "" {
			return nil, huma.Error400BadRequest("missing id")
		}
		participant, err := h.service.GetParticipantByExternalID(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get participant", err)
		}
		return &struct{ Body *db_gen.Participant }{Body: participant}, nil
	})

	huma.Get(api, "/api/qr", func(ctx context.Context, input *struct {
		ParticipantID string `query:"participant_id"`
	}) (*struct {
		ContentType string `header:"Content-Type"`
		Body        []byte
	}, error) {
		logger := server.LoggerFromContext(ctx)
		if input.ParticipantID == "" {
			return nil, huma.Error400BadRequest("missing participant_id")
		}

		logger.Info("generating qr code", zap.String("participant_id", input.ParticipantID))

		qr, err := qrcode.New(input.ParticipantID, qrcode.Medium)
		if err != nil {
			logger.Error("failed to generate qr code", zap.Error(err))
			return nil, huma.Error500InternalServerError("failed to generate qr code", err)
		}

		png, err := qr.PNG(256)
		if err != nil {
			logger.Error("failed to render qr code as png", zap.Error(err))
			return nil, huma.Error500InternalServerError("failed to render qr code", err)
		}

		return &struct {
			ContentType string `header:"Content-Type"`
			Body        []byte
		}{
			ContentType: "image/png",
			Body:        png,
		}, nil
	})
}
