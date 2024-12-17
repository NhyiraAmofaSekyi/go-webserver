package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		code          int
		expectedLevel string
		exepectedCode int
	}{
		{
			name:  "Error Log",
			code:  http.StatusInternalServerError,
			level: "ERROR",
		},
		{
			name:  "Info Log",
			code:  http.StatusOK,
			level: "INFO",
		},
	}
	tmpDir, err := os.MkdirTemp("", "logger_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Change to temp directory for test
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger.Init(true)

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test logs

			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.code)
			})

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)

			// Create response recorder
			rr := httptest.NewRecorder()

			logging := Logging(handler)

			logging.ServeHTTP(rr, req)

		})
	}

	// Close writer and restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check log file
	files, err := os.ReadDir("logs")
	if err != nil {
		t.Fatalf("Failed to read logs directory: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("No log file created")
	}

	logContent, err := os.ReadFile("logs/" + files[0].Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	fileOutput := string(logContent)

	// Verify each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check stdout
			if !strings.Contains(output, strconv.Itoa(tt.code)) {
				t.Errorf("Stdout does not contain code: %d", tt.code)
			}
			if !strings.Contains(output, tt.level) {
				t.Errorf("Stdout does not contain level: %s", tt.level)
			}

			// Check file output
			if !strings.Contains(fileOutput, strconv.Itoa(tt.code)) {
				t.Errorf("Stdout does not contain code: %d", tt.code)
			}
			if !strings.Contains(fileOutput, tt.level) {
				t.Errorf("Stdout does not contain level: %s", tt.level)
			}
		})
	}

}
