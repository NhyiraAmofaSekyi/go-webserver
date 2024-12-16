package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

// Reset CORS config before each test
func resetCorsConfig() {
	corsConfig = nil
	corsOnce = sync.Once{}
}

func TestCorsWrapper(t *testing.T) {
	tests := []struct {
		name            string
		setupConfig     *CorsConfig
		requestMethod   string
		requestOrigin   string
		requestHeaders  map[string]string
		expectedStatus  int
		expectedHeaders map[string]string
	}{
		{
			name: "Custom Global Config",
			setupConfig: &CorsConfig{
				AllowedOrigins: []string{"http://example.com"},
				AllowedMethods: []string{"GET", "POST"},
				AllowedHeaders: []string{"X-Custom-Header"},
				Credentials:    false,
				MaxAge:         "3600",
			},
			requestMethod:  "OPTIONS",
			requestOrigin:  "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://example.com",
				"Access-Control-Allow-Methods":     "GET, POST",
				"Access-Control-Allow-Headers":     "X-Custom-Header",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Max-Age":           "3600",
			},
		},
		{
			name:           "Fallback Default Config",
			setupConfig:    nil, // Don't set config, let it use default
			requestMethod:  "GET",
			requestOrigin:  "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config before each test
			resetCorsConfig()

			// Set config if provided
			if tt.setupConfig != nil {
				SetCorsConfig(tt.setupConfig)
			}

			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create request
			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}

			// Add any additional request headers
			for key, value := range tt.requestHeaders {
				req.Header.Set(key, value)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create middleware handler
			corsHandler := CorsWrapper(handler)

			// Serve request
			corsHandler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check headers
			for key, expectedValue := range tt.expectedHeaders {
				if value := rr.Header().Get(key); value != expectedValue {
					t.Errorf("handler returned wrong header %s: got %v want %v",
						key, value, expectedValue)
				}
			}
		})
	}
}

func TestValidateAndSetDefaults(t *testing.T) {
	tests := []struct {
		name           string
		inputConfig    *CorsConfig
		expectedConfig *CorsConfig
	}{
		{
			name:        "Empty Config - Should Use All Defaults",
			inputConfig: &CorsConfig{},
			expectedConfig: &CorsConfig{
				AllowedOrigins: DefaultCorsConfig.AllowedOrigins,
				AllowedMethods: DefaultCorsConfig.AllowedMethods,
				AllowedHeaders: DefaultCorsConfig.AllowedHeaders,
			},
		},
		{
			name: "Only Origins Set",
			inputConfig: &CorsConfig{
				AllowedOrigins: []string{"http://example.com"},
			},
			expectedConfig: &CorsConfig{
				AllowedOrigins: []string{"http://example.com"},
				AllowedMethods: DefaultCorsConfig.AllowedMethods,
				AllowedHeaders: DefaultCorsConfig.AllowedHeaders,
			},
		},
		{
			name: "Methods Without OPTIONS",
			inputConfig: &CorsConfig{
				AllowedMethods: []string{"GET", "POST"},
			},
			expectedConfig: &CorsConfig{
				AllowedOrigins: DefaultCorsConfig.AllowedOrigins,
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders: DefaultCorsConfig.AllowedHeaders,
			},
		},
		{
			name: "Methods With OPTIONS",
			inputConfig: &CorsConfig{
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			},
			expectedConfig: &CorsConfig{
				AllowedOrigins: DefaultCorsConfig.AllowedOrigins,
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders: DefaultCorsConfig.AllowedHeaders,
			},
		},
		{
			name: "Custom Headers Only",
			inputConfig: &CorsConfig{
				AllowedHeaders: []string{"X-Custom-Header"},
			},
			expectedConfig: &CorsConfig{
				AllowedOrigins: DefaultCorsConfig.AllowedOrigins,
				AllowedMethods: DefaultCorsConfig.AllowedMethods,
				AllowedHeaders: []string{"X-Custom-Header"},
			},
		},
		{
			name: "Credentials and MaxAge Only",
			inputConfig: &CorsConfig{
				Credentials: true,
				MaxAge:      "3600",
			},
			expectedConfig: &CorsConfig{
				AllowedOrigins: DefaultCorsConfig.AllowedOrigins,
				AllowedMethods: DefaultCorsConfig.AllowedMethods,
				AllowedHeaders: DefaultCorsConfig.AllowedHeaders,
				Credentials:    true,
				MaxAge:         "3600",
			},
		},
		{
			name: "Full Custom Config",
			inputConfig: &CorsConfig{
				AllowedOrigins: []string{"http://custom.com"},
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders: []string{"X-Custom-Header"},
				Credentials:    true,
				MaxAge:         "3600",
			},
			expectedConfig: &CorsConfig{
				AllowedOrigins: []string{"http://custom.com"},
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders: []string{"X-Custom-Header"},
				Credentials:    true,
				MaxAge:         "3600",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateAndSetDefaults(tt.inputConfig)

			// Compare AllowedOrigins
			if !reflect.DeepEqual(result.AllowedOrigins, tt.expectedConfig.AllowedOrigins) {
				t.Errorf("AllowedOrigins mismatch:\ngot: %v\nwant: %v",
					result.AllowedOrigins, tt.expectedConfig.AllowedOrigins)
			}

			// Compare AllowedMethods
			if !reflect.DeepEqual(result.AllowedMethods, tt.expectedConfig.AllowedMethods) {
				t.Errorf("AllowedMethods mismatch:\ngot: %v\nwant: %v",
					result.AllowedMethods, tt.expectedConfig.AllowedMethods)
			}

			// Compare AllowedHeaders
			if !reflect.DeepEqual(result.AllowedHeaders, tt.expectedConfig.AllowedHeaders) {
				t.Errorf("AllowedHeaders mismatch:\ngot: %v\nwant: %v",
					result.AllowedHeaders, tt.expectedConfig.AllowedHeaders)
			}

			// Compare Credentials
			if result.Credentials != tt.expectedConfig.Credentials {
				t.Errorf("Credentials mismatch:\ngot: %v\nwant: %v",
					result.Credentials, tt.expectedConfig.Credentials)
			}

			// Compare MaxAge
			if result.MaxAge != tt.expectedConfig.MaxAge {
				t.Errorf("MaxAge mismatch:\ngot: %v\nwant: %v",
					result.MaxAge, tt.expectedConfig.MaxAge)
			}
		})
	}
}

func TestSetCorsConfig(t *testing.T) {
	resetCorsConfig()

	config1 := &CorsConfig{
		AllowedOrigins: []string{"http://first.com"},
	}

	config2 := &CorsConfig{
		AllowedOrigins: []string{"http://second.com"},
	}

	// Set first config
	SetCorsConfig(config1)
	firstConfig := GetCorsConfig()

	// Try to set second config
	SetCorsConfig(config2)
	secondConfig := GetCorsConfig()

	// Verify that the second set didn't change the config (due to sync.Once)
	if firstConfig.AllowedOrigins[0] != secondConfig.AllowedOrigins[0] {
		t.Error("CORS config was modified after initial set")
	}

	if secondConfig.AllowedOrigins[0] != "http://first.com" {
		t.Error("CORS config doesn't match first set configuration")
	}
}
