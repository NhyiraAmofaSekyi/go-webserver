package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignIn(t *testing.T) {
	// Test for valid input
	t.Run("valid input", func(t *testing.T) {
		name := "testuser"
		body, _ := json.Marshal(map[string]string{"name": name})
		req, err := http.NewRequest("POST", "/v1/signin", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(SignIn)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the token is present and valid
		responseMap := make(map[string]string)
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Fatalf("could not parse response: %v", err)
		}

		if responseMap["status"] != "ok" || responseMap["route"] != "auth sign in" {
			t.Errorf("handler returned unexpected body: got status %v route %v, want status ok route auth sign in", responseMap["status"], responseMap["route"])
		}

		if responseMap["token"] == "" {
			t.Errorf("expected token to be not empty")
		}
	})

	// Test for invalid JSON input
	t.Run("invalid JSON", func(t *testing.T) {
		body := []byte(`{"name":}`)
		req, err := http.NewRequest("POST", "/v1/signin", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(SignIn)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
