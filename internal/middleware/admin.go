package middleware

import (
	"bookshelf/pkg/utils"
	"net/http"
)

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user").(*utils.Claims)
		if !ok {
			http.Error(w, "User information not found", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			http.Error(w, "Admin privileges required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
