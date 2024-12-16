package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create temporary directory for test logs
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

	// Initialize logger
	Init(true)

	// Test messages
	tests := []struct {
		name    string
		logFunc func(string, ...interface{})
		message string
		level   string
		color   string
	}{
		{
			name:    "Debug Message",
			logFunc: Debug,
			message: "test debug message",
			level:   "DEBUG",
			color:   colorYellow,
		},
		{
			name:    "Info Message",
			logFunc: Info,
			message: "test info message",
			level:   "INFO",
			color:   colorGreen,
		},
		{
			name:    "Error Message",
			logFunc: Error,
			message: "test error message",
			level:   "ERROR",
			color:   colorRed,
		},
	}

	// Log test messages
	for _, tt := range tests {
		tt.logFunc(tt.message)
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
			if !strings.Contains(output, tt.message) {
				t.Errorf("Stdout does not contain message: %s", tt.message)
			}
			if !strings.Contains(output, tt.level) {
				t.Errorf("Stdout does not contain level: %s", tt.level)
			}
			if !strings.Contains(output, tt.color) {
				t.Errorf("Stdout does not contain color code: %s", tt.color)
			}

			// Check file output
			if !strings.Contains(fileOutput, tt.message) {
				t.Errorf("Log file does not contain message: %s", tt.message)
			}
			if !strings.Contains(fileOutput, tt.level) {
				t.Errorf("Log file does not contain level: %s", tt.level)
			}
		})
	}
}

func TestDebugMode(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "logger_debug_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(currentDir)

	tests := []struct {
		name      string
		debugMode bool
		shouldLog bool
	}{
		{
			name:      "Debug Mode On",
			debugMode: true,
			shouldLog: true,
		},
		{
			name:      "Debug Mode Off",
			debugMode: false,
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Initialize logger with test debug mode
			Init(tt.debugMode)

			// Log debug message
			Debug("debug test message")

			// Restore stdout
			w.Close()
			os.Stdout = old

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.shouldLog && !strings.Contains(output, "debug test message") {
				t.Error("Debug message not logged when it should be")
			}
			if !tt.shouldLog && strings.Contains(output, "debug test message") {
				t.Error("Debug message logged when it shouldn't be")
			}
		})
	}
}
