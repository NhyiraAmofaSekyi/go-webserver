// middleware/integration_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
)

func TestMiddlewareStackIntegration(t *testing.T) {

	corsConfig = nil
	corsOnce = sync.Once{}
	// Setup temporary directory for logs
	tmpDir, err := os.MkdirTemp("", "middleware_integration_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create logs directory
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.Mkdir(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	// Change to temp directory and initialize logger
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(currentDir)

	logger.Init(true)

	tests := []struct {
		name            string
		method          string
		path            string
		headers         map[string]string
		expectedStatus  int
		expectedHeaders map[string]string
		expectedInLogs  []string
	}{
		{
			name:   "Successful Request with CORS",
			method: "GET",
			path:   "/test",
			headers: map[string]string{
				"Origin": "http://example.com",
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			expectedInLogs: []string{
				"INFO",
				"GET",
				"/test",
				"200",
			},
		},
		{
			name:   "Options Request",
			method: "OPTIONS",
			path:   "/test",
			headers: map[string]string{
				"Origin": "http://example.com",
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			},
			expectedInLogs: []string{
				"INFO",
				"OPTIONS",
				"/test",
			},
		},
		{
			name:   "Server Error with CORS",
			method: "GET",
			path:   "/error",
			headers: map[string]string{
				"Origin": "http://example.com",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			expectedInLogs: []string{
				"ERROR",
				"GET",
				"/error",
				"500",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout for logging
			old := os.Stdout
			_, w, _ := os.Pipe()
			os.Stdout = w

			// Create final handler that returns the expected status
			finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "error") {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(tt.expectedStatus)
			})

			// Create middleware stack
			stack := CreateStack(
				Logging,
				CorsWrapper,
			)

			// Create test server with middleware stack
			handler := stack(finalHandler)

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(rr, req)

			// Restore stdout and get output
			w.Close()
			os.Stdout = old

			// Verify status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Verify CORS headers
			for key, expected := range tt.expectedHeaders {
				if got := rr.Header().Get(key); got != expected {
					t.Errorf("handler returned wrong header %s: got %v want %v",
						key, got, expected)
				}
			}

			// Read log file
			files, err := os.ReadDir("logs")
			if err != nil {
				t.Fatalf("Failed to read logs directory: %v", err)
			}
			if len(files) == 0 {
				t.Fatal("No log file created")
			}

			logContent, err := os.ReadFile(filepath.Join("logs", files[0].Name()))
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			logOutput := string(logContent)

			// Verify log contents
			for _, expected := range tt.expectedInLogs {
				if !strings.Contains(logOutput, expected) {
					t.Errorf("Log does not contain expected content %q\nLog output: %s",
						expected, logOutput)
				}
			}
		})
	}
}
