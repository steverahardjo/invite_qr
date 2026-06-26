// handle to serve the frontend dynamic routing
package serving

import (
	"context"
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
