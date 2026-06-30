package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
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

func (h *JwtHandler) RegisterHumaRoutes(api huma.API) {
	type LoginInput struct {
		Body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
	}
	type LoginOutput struct {
		Body struct {
			Token string `json:"token"`
		}
	}

	huma.Post(api, "/api/admin/login", func(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
		token, err := h.service.LoginAdmin(ctx, input.Body.Username, input.Body.Password)
		if err != nil {
			return nil, huma.Error401Unauthorized("invalid credentials")
		}
		return &LoginOutput{Body: struct{ Token string `json:"token"` }{Token: token}}, nil
	})
}
