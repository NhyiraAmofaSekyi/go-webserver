package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NhyiraAmofaSekyi/go-webserver/utils"
	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("sample")

type AuthUserIDKey string

const AuthUserID AuthUserIDKey = "middleware.auth.userID"

func NewJWT() {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Convert the secret from string to []byte

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		fmt.Println("Error in signing:", err)
		return
	}
	fmt.Println(tokenString)
}

func Parse() {
	// sample token string taken from the New example
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.KQClWEcvHJtvMXzok3gvz7kPZqCR_SNCZR0Sj37mJAs"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

}

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

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "auth" {
			// Use the RespondWithError utility to send an unauthorized response
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		// ParseJWT()
		ctx := context.WithValue(r.Context(), AuthUserID, authHeader)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	}
}
