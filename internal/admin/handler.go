package admin

import (
	"context"
	"encoding/json"
	"net/http"

	db "invite_qr/db/db_gen"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type adminStore interface {
	ListParticipants(ctx context.Context, arg db.ListParticipantsParams) ([]db.Participant, error)
	InsertParticipant(ctx context.Context, arg db.InsertParticipantParams) (db.Participant, error)
	UpdateParticipantAccessedByExternalID(ctx context.Context, externalID uuid.UUID) (db.Participant, error)
}

type Handler struct {
	store adminStore
}

func NewHandler(store adminStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ListParticipants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		participants, err := h.store.ListParticipants(
			r.Context(),
			db.ListParticipantsParams{
				Limit:  100,
				Offset: 0,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if participants == nil {
			participants = []db.Participant{}
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(participants); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) AddParticipant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			WaNumber string `json:"wa_number"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := h.store.InsertParticipant(
			r.Context(),
			db.InsertParticipantParams{
				Name:     req.Name,
				Email:    req.Email,
				WaNumber: req.WaNumber,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) MarkAttendance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ParticipantID string `json:"participant_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.ParticipantID == "" {
			http.Error(w, "missing participant_id", http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(req.ParticipantID)
		if err != nil {
			http.Error(w, "invalid participant_id", http.StatusBadRequest)
			return
		}

		_, err = h.store.UpdateParticipantAccessedByExternalID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) RegisterHumaRoutes(api huma.API) {
	huma.Get(api, "/api/admin/participants", func(ctx context.Context, input *struct{}) (*struct {
		Body []db.Participant
	}, error) {
		participants, err := h.store.ListParticipants(ctx, db.ListParticipantsParams{Limit: 100, Offset: 0})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list participants", err)
		}
		if participants == nil {
			participants = []db.Participant{}
		}
		return &struct{ Body []db.Participant }{Body: participants}, nil
	}, huma.OperationTags("admin"))

	huma.Post(api, "/api/admin/participants", func(ctx context.Context, input *struct {
		Body struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			WaNumber string `json:"wa_number"`
		}
	}) (*struct {
		Body struct{}
	}, error) {
		_, err := h.store.InsertParticipant(ctx, db.InsertParticipantParams{
			Name:     input.Body.Name,
			Email:    input.Body.Email,
			WaNumber: input.Body.WaNumber,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to add participant", err)
		}
		return &struct{ Body struct{} }{}, nil
	}, huma.OperationTags("admin"))

	huma.Post(api, "/api/admin/attendance", func(ctx context.Context, input *struct {
		Body struct {
			ParticipantID string `json:"participant_id"`
		}
	}) (*struct {
		Body struct{}
	}, error) {
		if input.Body.ParticipantID == "" {
			return nil, huma.Error400BadRequest("missing participant_id")
		}
		id, err := uuid.Parse(input.Body.ParticipantID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid participant_id")
		}
		_, err = h.store.UpdateParticipantAccessedByExternalID(ctx, id)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to mark attendance", err)
		}
		return &struct{ Body struct{} }{}, nil
	}, huma.OperationTags("admin"))
}
