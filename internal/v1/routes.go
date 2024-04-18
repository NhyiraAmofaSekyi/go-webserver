// v1/routes.go
package v1

import (
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
	"github.com/NhyiraAmofaSekyi/go-webserver/internal/v1/auth"
)

// NewRouter returns a new http.ServeMux with v1 routes configured
func NewRouter() *http.ServeMux {
	v1Router := http.NewServeMux()
	authRouter := auth.NewRouter()

	v1Router.HandleFunc("GET /healthz", healthzHandler) // Note the path is just "/healthz" now
	v1Router.HandleFunc("GET /secure", middleware.AuthMiddleware(secureHandler))
	v1Router.Handle("/", http.StripPrefix("/auth", authRouter))

	return v1Router
}
