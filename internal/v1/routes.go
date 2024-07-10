// v1/routes.go
package v1

import (
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
	"github.com/NhyiraAmofaSekyi/go-webserver/internal/v1/auth"
	"github.com/NhyiraAmofaSekyi/go-webserver/internal/v1/users"
)

// NewRouter returns a new http.ServeMux with v1 routes configured
func NewRouter() *http.ServeMux {
	v1Router := http.NewServeMux()
	authRouter := auth.NewRouter()
	userRouter := users.NewRouter()

	v1Router.HandleFunc("GET /healthz", HealthzHandler) // Note the path is just "/healthz" now
	v1Router.HandleFunc("GET /secure", middleware.AuthMiddleware(SecureHandler))
	v1Router.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	v1Router.Handle("/users/", http.StripPrefix("/users", userRouter))

	return v1Router
}
