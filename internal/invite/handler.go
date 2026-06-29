// Package invite provides HTTP handlers and services for sending event
// invitations via WhatsApp and email, including single and bulk delivery.
package invite

import (
	"context"
	"net/http"
	"strconv"
)

// inviteService defines the contract for sending invitations.
// The production implementation is *Service; tests can provide a mock.
type inviteService interface {
	// BulkSendInvite sends invitations to all unsent participants.
	BulkSendInvite(ctx context.Context, eventTitle string) error
	// SendInviteOnetime sends a single invitation to the specified guest.
	SendInviteOnetime(ctx context.Context, guestID int32, eventTitle string, email string, waNumber string, name string) error
}

// Handler holds the dependencies needed by the invite HTTP handlers.
type Handler struct {
	service inviteService
}

// NewHandler creates a Handler backed by the given invite service.
func NewHandler(service inviteService) *Handler {
	return &Handler{
		service: service,
	}
}

// HandleBulkInvite returns an HTTP handler that triggers a bulk invitation
// send for all unsent participants using the provided context for lifecycle
// management (e.g. cancellation on server shutdown).
func (h *Handler) HandleBulkInvite(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.service.BulkSendInvite(ctx, "My Event")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// HandleSendInviteOnetime returns an HTTP handler that sends a single invitation
// to a specific guest identified by guest_id query parameter, with optional
// email, wa_number, and name overrides.
func (h *Handler) HandleSendInviteOnetime(eventTitle string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestID := r.URL.Query().Get("guest_id")
		email := r.URL.Query().Get("email")
		waNumber := r.URL.Query().Get("wa_number")
		name := r.URL.Query().Get("name")

		if guestID == "" {
			http.Error(w, "missing guest_id", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(guestID, 10, 32)
		if err != nil {
			http.Error(w, "invalid guest_id", http.StatusBadRequest)
			return
		}

		if email == "" && waNumber == "" {
			http.Error(w, "either email or wa_number is required", http.StatusBadRequest)
			return
		}

		if err := h.service.SendInviteOnetime(
			r.Context(),
			int32(id),
			eventTitle,
			email,
			waNumber,
			name,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
