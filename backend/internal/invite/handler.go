package invite

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
)

type inviteService interface {
	BulkSendInvite(ctx context.Context, eventTitle string) error
	SendInviteOnetime(ctx context.Context, guestID int32, eventTitle string, email string, waNumber string, name string) error
}

type Handler struct {
	service inviteService
}

func NewHandler(service inviteService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleBulkInvite(eventTitle string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.service.BulkSendInvite(r.Context(), eventTitle)
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

func (h *Handler) RegisterHumaRoutes(api huma.API, eventTitle string) {
	huma.Post(api, "/api/bulk-invite", func(ctx context.Context, input *struct{}) (*struct {
		Body struct{}
	}, error) {
		if err := h.service.BulkSendInvite(ctx, eventTitle); err != nil {
			return nil, huma.Error500InternalServerError("failed to send bulk invite", err)
		}
		return &struct{ Body struct{} }{}, nil
	})

	huma.Get(api, "/api/send-invite", func(ctx context.Context, input *struct {
		GuestID  int32  `query:"guest_id"`
		Email    string `query:"email"`
		WaNumber string `query:"wa_number"`
		Name     string `query:"name"`
	}) (*struct {
		Body struct{}
	}, error) {
		if input.GuestID == 0 {
			return nil, huma.Error400BadRequest("missing guest_id")
		}
		if input.Email == "" && input.WaNumber == "" {
			return nil, huma.Error400BadRequest("either email or wa_number is required")
		}
		if err := h.service.SendInviteOnetime(ctx, input.GuestID, eventTitle, input.Email, input.WaNumber, input.Name); err != nil {
			return nil, huma.Error500InternalServerError("failed to send invite", err)
		}
		return &struct{ Body struct{} }{}, nil
	})
}
