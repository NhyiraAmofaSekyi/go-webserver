package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("hqwebfoeuyrh38y24-821")

// Function to generate a new JWT for a given name
func GenerateJWT(name string) (string, error) {

	if name == "" {
		return "", fmt.Errorf("no string provided")
	}
	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name, // Include the name in the token
		"nbf":  time.Now().Unix(),
		"exp":  expirationTime,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
		// log.Fatalf("Error parsing token: %v", err)
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Type assertion to extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		fmt.Println("Invalid token or failed claims assertion")
		return nil, fmt.Errorf("invalid token or failed claims assertion")
	}
}
