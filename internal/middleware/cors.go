package middleware

import (
	"net/http"
	"strings"
	"sync"
)

type CorsConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         string // Optional
	Credentials    bool   // Optional, for Access-Control-Allow-Credentials
}

var (
	corsConfig *CorsConfig
	corsOnce   sync.Once

	DefaultCorsConfig = &CorsConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
)

// validateAndSetDefaults ensures all required fields have values
func validateAndSetDefaults(config *CorsConfig) *CorsConfig {
	// Create new config to avoid modifying the input
	validConfig := &CorsConfig{
		Credentials: config.Credentials, // Keep original credentials setting
		MaxAge:      config.MaxAge,      // Keep original max age setting
	}

	// Validate and set Origins
	if len(config.AllowedOrigins) == 0 {
		validConfig.AllowedOrigins = DefaultCorsConfig.AllowedOrigins
	} else {
		validConfig.AllowedOrigins = config.AllowedOrigins
	}

	// Validate and set Methods
	if len(config.AllowedMethods) == 0 {
		validConfig.AllowedMethods = DefaultCorsConfig.AllowedMethods
	} else {
		// Ensure OPTIONS is included if not present
		hasOptions := false
		for _, method := range config.AllowedMethods {
			if method == "OPTIONS" {
				hasOptions = true
				break
			}
		}
		if !hasOptions {
			validConfig.AllowedMethods = append(config.AllowedMethods, "OPTIONS")
		} else {
			validConfig.AllowedMethods = config.AllowedMethods
		}
	}

	// Validate and set Headers
	if len(config.AllowedHeaders) == 0 {
		validConfig.AllowedHeaders = DefaultCorsConfig.AllowedHeaders
	} else {
		validConfig.AllowedHeaders = config.AllowedHeaders
	}

	return validConfig
}

// SetCorsConfig sets the global CORS configuration with validation
func SetCorsConfig(config *CorsConfig) error {
	if config == nil {
		corsConfig = DefaultCorsConfig
		return nil
	}

	corsOnce.Do(func() {
		corsConfig = validateAndSetDefaults(config)
	})
	return nil
}

// GetCorsConfig returns the current CORS configuration
func GetCorsConfig() *CorsConfig {
	if corsConfig == nil {
		corsConfig = DefaultCorsConfig
	}
	return corsConfig
}

// CorsWrapper wraps an http.Handler with CORS headers.
func CorsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := GetCorsConfig()
		origin := r.Header.Get("Origin")

		if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			for _, allowedOrigin := range config.AllowedOrigins {
				if allowedOrigin == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// Handle methods
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))

		// Handle headers
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

		// Handle credentials if set
		if config.Credentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle max age if set
		if config.MaxAge != "" {
			w.Header().Set("Access-Control-Max-Age", config.MaxAge)
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		h.ServeHTTP(w, r)
	})
}
