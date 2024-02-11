// utils/cors.go
package utils

import (
	"net/http"
)

// CorsWrapper wraps an http.Handler with CORS headers.
func CorsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set your CORS headers here
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Check if the request is for the OPTIONS method, return immediately
		if r.Method == "OPTIONS" {
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		h.ServeHTTP(w, r)
	})
}
