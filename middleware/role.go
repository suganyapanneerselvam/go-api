package middleware

import "net/http"

func AdminOnly(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value("role")

		if role != "admin" {
			http.Error(w, "Admin only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
