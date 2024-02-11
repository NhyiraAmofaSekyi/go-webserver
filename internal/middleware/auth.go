package middleware

import (
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "auth" {
			// Use the RespondWithError utility to send an unauthorized response
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	}
}
