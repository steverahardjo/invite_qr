package invitation

import (
	"context"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
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
