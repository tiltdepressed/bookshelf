// Package middleware
package middleware

import (
	"bookshelf/pkg/utils"
	"context"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(tokenHeader, "Bearer ")
		if tokenString == tokenHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
