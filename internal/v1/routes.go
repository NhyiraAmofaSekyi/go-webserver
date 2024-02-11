// v1/routes.go
package v1

import (
	"net/http"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
)

// NewRouter returns a new http.ServeMux with v1 routes configured
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthzHandler) // Note the path is just "/healthz" now
	mux.HandleFunc("/secure", middleware.AuthMiddleware(secureHandler))
	return mux
}
