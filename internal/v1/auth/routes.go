package auth

import (
	"net/http"

	databaseCfg "github.com/NhyiraAmofaSekyi/go-webserver/internal/db"
)

// NewRouter returns a new http.ServeMux with v1 routes configured
func NewRouter(dbCFG *databaseCfg.DBConfig) *http.ServeMux {
	authRouter := http.NewServeMux()

	authRouter.HandleFunc("POST /signIn", SignIn) // Note the path is just "/healthz" now
	authRouter.HandleFunc("GET /SignOut", SignOut)

	return authRouter
}
