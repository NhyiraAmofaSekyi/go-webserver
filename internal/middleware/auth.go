package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/auth"
	"github.com/NhyiraAmofaSekyi/go-webserver/utils"
)

var hmacSampleSecret = []byte("sample")

type AuthUserIDKey string
type ServiceKey string

const Skey ServiceKey = "service"

const AuthUserID AuthUserIDKey = "middleware.auth.userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Split the authorization header to separate the bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized - Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header
		tokenString := parts[1]

		// Parse the JWT and validate it
		claims, err := auth.ParseJWT(tokenString)
		if err != nil {
			utils.RespondWithJSON(w, 403, map[string]string{"message": "unauthorised"})
			return
		}

		name, ok := claims["name"].(string)
		if !ok {
			utils.RespondWithJSON(w, 400, map[string]string{"message": "bad request"})
			return
		}

		// Retrieve the name claim from the token
		if exp, ok := claims["exp"].(float64); ok {
			currentTime := time.Now().Unix()
			if int64(exp) < currentTime {
				utils.RespondWithJSON(w, 403, map[string]string{"message": "forbidden"})
				return
			}
		}

		// Add the name to the request context
		ctx := context.WithValue(r.Context(), AuthUserID, name)
		req := r.WithContext(ctx)

		// Continue with the pipeline
		next.ServeHTTP(w, req)
	}
}

func ClearSessionCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   false, // Set to true in production when using HTTPS
		HttpOnly: true,  // Prevent JavaScript access to cookie
	})
}
