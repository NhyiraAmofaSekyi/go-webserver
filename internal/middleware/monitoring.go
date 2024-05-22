package middleware

import (
	"net/http"

	monitoring "github.com/NhyiraAmofaSekyi/go-webserver/internal/monitoring"
)

func Monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		monitoring.HttpRequestsTotal.WithLabelValues("api", r.Method, r.URL.Path).Inc()

		next.ServeHTTP(wrapped, r)

	})
}
