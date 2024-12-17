package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
)

func TestMain(m *testing.M) {
	// Initialize logger before running tests
	logger.Init(true) // Set debug mode to true for testing
	m.Run()
}

// TestRespondWithError tests the error response functionality
func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name          string
		code          int
		message       string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Client Error",
			code:          http.StatusBadRequest,
			message:       "Invalid request",
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid request",
		},
		{
			name:          "Server Error",
			code:          http.StatusInternalServerError,
			message:       "Internal server error",
			expectedCode:  http.StatusInternalServerError,
			expectedError: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the function being tested
			RespondWithError(rr, tt.code, tt.message)

			// Check status code
			if rr.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedCode)
			}

			// Check Content-Type header
			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want application/json",
					contentType)
			}

			// Check response body
			var response struct {
				Error string `json:"error"`
			}
			err := json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Could not decode response body: %v", err)
			}

			if response.Error != tt.expectedError {
				t.Errorf("handler returned unexpected error message: got %v want %v",
					response.Error, tt.expectedError)
			}
		})
	}
}

// TestRespondWithJSON tests the JSON response functionality
func TestRespondWithJSON(t *testing.T) {
	tests := []struct {
		name         string
		payload      interface{}
		expectedCode int
		expectError  bool
	}{
		{
			name: "Valid JSON Response",
			payload: map[string]string{
				"message": "success",
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name: "Complex JSON Response",
			payload: struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Active bool   `json:"active"`
			}{
				ID:     1,
				Name:   "Test",
				Active: true,
			},
			expectedCode: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "Response with Channel (Invalid JSON)",
			payload: struct {
				Ch chan int
			}{
				Ch: make(chan int),
			},
			expectedCode: http.StatusOK,
			expectError:  true,
		},
	}

	defer os.RemoveAll("logs")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			RespondWithJSON(rr, tt.expectedCode, tt.payload)

			// Check status code
			if !tt.expectError && rr.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedCode)
			}

			// Check Content-Type header
			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want application/json",
					contentType)
			}

			if !tt.expectError {
				// Verify the response can be decoded back to JSON
				var response interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Could not decode response body: %v", err)
				}
			}
		})
	}
}
