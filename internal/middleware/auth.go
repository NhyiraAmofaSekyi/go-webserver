package middleware

import (
	"context"
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/utils"
)

type AuthUserIDKey string

const AuthUserID AuthUserIDKey = "middleware.auth.userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "auth" {
			// Use the RespondWithError utility to send an unauthorized response
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), AuthUserID, authHeader)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	}
}
