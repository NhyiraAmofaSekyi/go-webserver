package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/v1/auth"
)

func TestHealthzHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthzHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := map[string]string{"status": "ok", "route": "v1"}
	var actual map[string]string
	body, _ := io.ReadAll(rr.Body)
	err = json.Unmarshal(body, &actual)
	if err != nil {
		t.Fatal("Could not unmarshal response:", err)
	}

	if actual["status"] != expected["status"] || actual["route"] != expected["route"] {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestSecureHandler(t *testing.T) {

	handler := middleware.AuthMiddleware(SecureHandler)

	// Test with a valid token
	t.Run("valid token", func(t *testing.T) {
		token := auth.GenerateTestToken("testuser", true)
		request, _ := http.NewRequest("GET", "/v1/secure", nil)
		request.Header.Set("Authorization", "Bearer "+token)

		responseRecorder := httptest.NewRecorder()
		handler(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusOK {
			t.Errorf("Expected status OK for valid token, got %v", status)
		}

	})

	// Test with an invalid token
	t.Run("invalid token", func(t *testing.T) {
		token := auth.GenerateTestToken("testuser", false)
		request, _ := http.NewRequest("GET", "/v1/secure", nil)
		request.Header.Set("Authorization", "Bearer "+token)

		responseRecorder := httptest.NewRecorder()
		handler(responseRecorder, request)

		if status := responseRecorder.Code; status != http.StatusForbidden {
			t.Errorf("Expected status Forbidden for invalid token, got %v", status)
		}
	})
}
