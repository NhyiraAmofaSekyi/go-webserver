package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	utils "github.com/NhyiraAmofaSekyi/go-webserver/utils"
	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("sample")

// Function to generate a new JWT for a given name
func generateJWT(name string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Second).Unix()
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

func signIn(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithJSON(w, 400, fmt.Sprintf("Error passing json: %v", err))
		return
	}

	jwtToken, err := generateJWT(params.Name)
	if err != nil {
		utils.RespondWithJSON(w, 400, fmt.Sprintf("Error passing json: %v", err))
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"status": "ok", "route": "auth sign in", "token": jwtToken})
}

func signOut(w http.ResponseWriter, r *http.Request) {

	utils.RespondWithJSON(w, 200, map[string]string{"status": "ok", "route": "auth sign out"})
}
