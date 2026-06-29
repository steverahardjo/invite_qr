// Package admin provides HTTP handlers for admin-level participant management
// including listing, creating, and marking attendance.
package admin

import (
	"context"
	"encoding/json"
	"net/http"

	db "invite_qr/db/db_gen"

	"github.com/google/uuid"
)

// adminStore defines the data access methods required by the admin handlers.
// It is satisfied by *db.Queries for production use and can be mocked in tests.
type adminStore interface {
	// ListParticipants returns a paginated list of participants.
	ListParticipants(ctx context.Context, arg db.ListParticipantsParams) ([]db.Participant, error)
	// InsertParticipant creates a new participant and returns it.
	InsertParticipant(ctx context.Context, arg db.InsertParticipantParams) (db.Participant, error)
	// UpdateParticipantAccessedByExternalID marks a participant as accessed by their external UUID.
	UpdateParticipantAccessedByExternalID(ctx context.Context, externalID uuid.UUID) (db.Participant, error)
}

// Handler holds the dependencies needed by the admin HTTP handlers.
type Handler struct {
	store adminStore
}

// NewHandler creates a new Handler with the given store implementation.
func NewHandler(store adminStore) *Handler {
	return &Handler{store: store}
}

// ListParticipants returns an HTTP handler that fetches up to 100 participants
// and responds with a JSON array.
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

// AddParticipant returns an HTTP handler that decodes a JSON body with
// name, email, and wa_number, creates a new participant, and responds 201 Created.
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

// MarkAttendance returns an HTTP handler that accepts a JSON body with
// participant_id (external UUID from QR code scan), and marks the
// participant as having attended the event.
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
