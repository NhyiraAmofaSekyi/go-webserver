// v1/routes.go
package v1

import (
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
)

// NewRouter returns a new http.ServeMux with v1 routes configured
func NewRouter() *http.ServeMux {
	v1Router := http.NewServeMux()

	v1Router.HandleFunc("GET /healthz", healthzHandler) // Note the path is just "/healthz" now
	v1Router.HandleFunc("GET /secure", middleware.AuthMiddleware(secureHandler))

	return v1Router
}
