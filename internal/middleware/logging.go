package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
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
		if wrapped.statusCode >= 500 {
			logger.Error("%d %s %s %v", wrapped.statusCode, r.Method, r.URL.Path, duration)
		} else if wrapped.statusCode >= 400 {
			logger.Info("%d %s %s %v", wrapped.statusCode, r.Method, r.URL.Path, duration)
		} else {
			logger.Info("%d %s %s %v", wrapped.statusCode, r.Method, r.URL.Path, duration)
		}
	})
}
