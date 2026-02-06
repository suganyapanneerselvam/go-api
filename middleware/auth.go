package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key") // Later move to env

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Format: Bearer <token>
		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// ✅ Get role from JWT
		role := claims["role"]

		// ✅ Store role in context
		ctx := context.WithValue(r.Context(), "role", role)

		// ✅ Continue with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
