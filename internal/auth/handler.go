// Package auth provides authentication and authorization handlers,
// including JWT-based login and middleware for protected routes.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
)

// authenticator defines the contract for admin login operations.
// The production implementation is *JwtService; tests can provide a mock.
type authenticator interface {
	// LoginAdmin validates admin credentials and returns a signed JWT token.
	LoginAdmin(ctx context.Context, username, password string) (string, error)
}

// JwtHandler handles admin authentication HTTP requests.
type JwtHandler struct {
	service authenticator
}

// NewJwtHandler creates a JwtHandler backed by the given authenticator.
func NewJwtHandler(service authenticator) *JwtHandler {
	return &JwtHandler{service: service}
}

// LoginAdmin returns an HTTP handler that accepts a JSON body with "username"
// and "password", validates credentials via the authenticator, and returns
// a signed JWT token on success or 401 on failure.
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
