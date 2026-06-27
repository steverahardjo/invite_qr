package auth

import (
	"context"
	"encoding/json"
	"net/http"
)

type authenticator interface {
	LoginAdmin(ctx context.Context, username, password string) (string, error)
}

type JwtHandler struct {
	service authenticator
}

func NewJwtHandler(service authenticator) *JwtHandler {
	return &JwtHandler{service: service}
}

func (h *JwtHandler) LoginAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		token, err := h.service.LoginAdmin(r.Context(), req.Username, req.Password)
		if err != nil {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
