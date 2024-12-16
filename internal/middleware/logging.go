package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
	monitoring "github.com/NhyiraAmofaSekyi/go-webserver/internal/monitoring"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

type ReqTime string

const ReqStartTime ReqTime = "reqStartTime"

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := context.WithValue(r.Context(), ReqStartTime, start)
		req := r.WithContext(ctx)
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, req)

		duration := time.Since(start)
		if wrapped.statusCode > 499 {
			logger.Error("%d %s %s %v", wrapped.statusCode, req.Method, req.URL.Path, duration)
			monitoring.HttpRequestErrorsTotal.WithLabelValues("api", req.Method, req.URL.Path, http.StatusText(wrapped.statusCode)).Inc()
		} else {
			logger.Info("%d %s %s %v", wrapped.statusCode, req.Method, req.URL.Path, duration)
		}

		monitoring.HttpRequestsTotal.WithLabelValues("api", req.Method, req.URL.Path).Inc()
		monitoring.HttpRequestDuration.WithLabelValues("api", req.Method, req.URL.Path).Observe(time.Since(start).Seconds())
	})
}
