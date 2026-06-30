package api

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UseAdminAuth(api huma.API, secretKey []byte) {
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		op := ctx.Operation()
		if op == nil {
			next(ctx)
			return
		}

		isAdmin := false
		for _, tag := range op.Tags {
			if tag == "admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			next(ctx)
			return
		}

		authHeader := ctx.Header("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.SetHeader("Content-Type", "application/json")
			ctx.SetStatus(http.StatusUnauthorized)
			ctx.BodyWriter().Write([]byte(`{"error":"missing or invalid Authorization header"}`))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil {
			ctx.SetHeader("Content-Type", "application/json")
			ctx.SetStatus(http.StatusUnauthorized)
			ctx.BodyWriter().Write([]byte(`{"error":"invalid or expired token"}`))
			return
		}

		next(ctx)
	})
}
