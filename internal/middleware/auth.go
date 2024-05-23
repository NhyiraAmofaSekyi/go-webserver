package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NhyiraAmofaSekyi/go-webserver/utils"
	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("sample")

type AuthUserIDKey string
type ServiceKey string

const Skey ServiceKey = "service"

const AuthUserID AuthUserIDKey = "middleware.auth.userID"

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token using a callback function to provide the key for verification
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used to sign the tokens
		return hmacSampleSecret, nil
	})

	if err != nil {
		// log.Fatalf("Error parsing token: %v", err)
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Type assertion to extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token claims:")
		for key, val := range claims {
			fmt.Printf("%s: %v\n", key, val)
		}
		return claims, nil
	} else {
		fmt.Println("Invalid token or failed claims assertion")
		return nil, fmt.Errorf("invalid token or failed claims assertion")
	}
}

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
		claims, err := ParseJWT(tokenString)
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
