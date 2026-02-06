package middleware

import (
	//"context"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value("role")

		if role != "admin" {
			http.Error(w, "Admin access only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
