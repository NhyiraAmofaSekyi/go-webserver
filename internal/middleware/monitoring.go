package middleware

import (
	"net/http"
	"time"

	monitoring "github.com/NhyiraAmofaSekyi/go-webserver/internal/monitoring"
)

func Monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		if wrapped.statusCode >= 500 {
			monitoring.HttpRequestErrorsTotal.WithLabelValues("api", r.Method, r.URL.Path, http.StatusText(wrapped.statusCode)).Inc()
		}
		monitoring.HttpRequestsTotal.WithLabelValues("api", r.Method, r.URL.Path).Inc()
		monitoring.HttpRequestDuration.WithLabelValues("api", r.Method, r.URL.Path).Observe(float64(time.Since(start)))

	})
}
