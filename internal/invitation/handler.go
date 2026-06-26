package invitation

import (
	"context"
	"database/sql"
	db "invite_qr/db/db_gen"
	"net/http"
	"net/url"
	"time"
)

type Handler struct {
	db      *sql.DB
	queries *db.Queries
	service *Service
}

func NewHandler(dbConn *sql.DB, whatsappSender *WhatsappSender, emailSender *EmailSender, baseWebURL *url.URL, hourLimit time.Time, dateLimit time.Time) *Handler {
	servobj := NewService(dbConn,
		whatsappSender,
		emailSender,
		dateLimit,
		hourLimit,
		baseWebURL,
	)
	return &Handler{
		service: servobj,
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
