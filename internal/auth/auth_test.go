package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		shouldError bool
	}{
		{
			name:        "Valid Username",
			userName:    "testuser",
			shouldError: false,
		},
		{
			name:        "Empty Username",
			userName:    "",
			shouldError: true, // Current implementation allows empty names
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.userName)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify token is not empty
			if token == "" {
				t.Error("Token is empty")
			}

			// Parse and verify token contents
			claims, err := ParseJWT(token)
			if err != nil {
				t.Errorf("Failed to parse generated token: %v", err)
				return
			}

			// Check username claim
			if name, ok := claims["name"].(string); !ok || name != tt.userName {
				t.Errorf("Expected name claim %s, got %v", tt.userName, claims["name"])
			}

			// Check expiration
			if exp, ok := claims["exp"].(float64); !ok {
				t.Error("Expiration claim not found or invalid type")
			} else {
				expTime := time.Unix(int64(exp), 0)
				if expTime.Before(time.Now()) {
					t.Error("Token is already expired")
				}
				if time.Until(expTime) > time.Hour+time.Second {
					t.Error("Expiration time is more than 1 hour")
				}
			}
		})
	}
}

func TestParseJWT(t *testing.T) {
	tests := []struct {
		name        string
		tokenString string
		shouldError bool
	}{
		{
			name:        "Valid Token",
			tokenString: "", // Will be populated with a valid token
			shouldError: false,
		},
		{
			name:        "Invalid Token",
			tokenString: "invalid.token.string",
			shouldError: true,
		},
		{
			name:        "Empty Token",
			tokenString: "",
			shouldError: true,
		},
		{
			name:        "Expired Token",
			tokenString: "", // Will be populated with an expired token
			shouldError: true,
		},
	}

	// Generate a valid token for the first test case
	validToken, err := GenerateJWT("testuser")
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}
	tests[0].tokenString = validToken

	// Generate an expired token for the last test case
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": "testuser",
		"nbf":  time.Now().Add(-2 * time.Hour).Unix(),
		"exp":  time.Now().Add(-1 * time.Hour).Unix(),
	})
	expiredTokenString, err := expiredToken.SignedString(hmacSampleSecret)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}
	tests[3].tokenString = expiredTokenString

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseJWT(tt.tokenString)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify claims for valid token
			if name, ok := claims["name"].(string); !ok || name == "" {
				t.Error("Name claim is missing or empty")
			}

			if exp, ok := claims["exp"].(float64); !ok {
				t.Error("Expiration claim is missing")
			} else {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					t.Error("Token is expired")
				}
			}
		})
	}
}
