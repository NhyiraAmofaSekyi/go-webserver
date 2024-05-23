package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
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
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
