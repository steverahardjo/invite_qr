package participant

import (
	"encoding/json"
	"net/http"

	db "invite_qr/db/db_gen"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		db:      pool,
		queries: db.New(pool),
	}
}

func (h *Handler) ListParticipants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		participants, err := h.queries.ListParticipants(
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
		if _, err := h.queries.InsertParticipant(
			r.Context(),
			db.InsertParticipantParams{
				Name:     req.Name,
				Email:    req.Email,
				WaNumber: req.WaNumber,
			},
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) UpdateParticipantAccessed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ParticipantID int64 `json:"participant_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := h.queries.UpdateParticipantAccessed(
			r.Context(),
			db.UpdateParticipantAccessedParams{
				ParticipantID: req.ParticipantID,
			},
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
