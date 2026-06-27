package admin

import (
	"context"
	"encoding/json"
	"net/http"

	db "invite_qr/db/db_gen"
)

type adminStore interface {
	ListParticipants(ctx context.Context, arg db.ListParticipantsParams) ([]db.Participant, error)
	InsertParticipant(ctx context.Context, arg db.InsertParticipantParams) (db.Participant, error)
	UpdateParticipantAccessed(ctx context.Context, arg db.UpdateParticipantAccessedParams) (db.Participant, error)
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

func (h *Handler) UpdateParticipantAccessed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ParticipantID int64  `json:"participant_id"`
			Email         string `json:"email"`
			WaNumber      string `json:"wa_number"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := h.store.UpdateParticipantAccessed(
			r.Context(),
			db.UpdateParticipantAccessedParams{
				ID:       int32(req.ParticipantID),
				Email:    req.Email,
				WaNumber: req.WaNumber,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
