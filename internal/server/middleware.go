package server

import (
	"context"
	"net/http"
	"strings"
)

type userIDKeyType struct{}

var userIDKey = userIDKeyType{}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := h.authService.ValidateToken(parts[1])
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
