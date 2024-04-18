package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("sample")

type AuthUserIDKey string

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
		log.Fatalf("Error parsing token: %v", err)
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

// func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader != "auth" {
// 			// Use the RespondWithError utility to send an unauthorized response
// 			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
// 			return
// 		}
// 		// ParseJWT()
// 		ctx := context.WithValue(r.Context(), AuthUserID, authHeader)
// 		req := r.WithContext(ctx)
// 		next.ServeHTTP(w, req)
// 	}
// }

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
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		// Retrieve the name claim from the token
		name, ok := claims["name"].(string)
		if !ok {
			http.Error(w, "Unauthorized - Name claim missing", http.StatusUnauthorized)
			return
		}

		// Add the name to the request context
		ctx := context.WithValue(r.Context(), AuthUserID, name)
		req := r.WithContext(ctx)

		// Continue with the pipeline
		next.ServeHTTP(w, req)
	}
}
