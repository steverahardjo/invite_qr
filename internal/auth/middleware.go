package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": "token_expired", "message": "Your session has expired. Please log in again."}`))
					return
				}

				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
