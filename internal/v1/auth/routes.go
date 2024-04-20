package auth

import (
	"net/http"
)

// NewRouter returns a new http.ServeMux with v1 routes configured
func NewRouter() *http.ServeMux {
	authRouter := http.NewServeMux()

	authRouter.HandleFunc("POST /signIn", SignIn) // Note the path is just "/healthz" now
	authRouter.HandleFunc("GET /SignOut", SignOut)

	return authRouter
}
