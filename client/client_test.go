package client

import (
	"testing"
	"time"
)

type mockResponse struct {
	Message string `json:"message"`
}

type mockRequest struct {
	Data string `json:"data"`
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "valid client creation",
			opts:    []Option{WithAPIKey("test-api-key")},
			wantErr: false,
		},
		{
			name: "client with custom options",
			opts: []Option{
				WithBaseURL(BrNe1),
				WithTimeout(5 * time.Second),
				WithAPIKey("test-api-key"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMgcClient(tt.opts...)
			if client == nil {
				t.Error("expected non-nil client")
				return
			}
			if client.config.APIKey != tt.apiKey {
				t.Errorf("expected API key %s, got %s", tt.apiKey, client.config.APIKey)
			}
		})
	}
}

func TestCoreClient_GetConfig(t *testing.T) {
	// Arrange
	expectedAPIKey := "test-api-key"
	expectedTimeout := 5 * time.Second

	client := NewMgcClient(WithAPIKey(expectedAPIKey),
		WithTimeout(expectedTimeout))

	// Act
	config := client.GetConfig()

	// Assert
	if config == nil {
		t.Fatal("expected non-nil config")
	}
	if config.APIKey != expectedAPIKey {
		t.Errorf("expected APIKey %s, got %s", expectedAPIKey, config.APIKey)
	}
	if config.Timeout != expectedTimeout {
		t.Errorf("expected Timeout %v, got %v", expectedTimeout, config.Timeout)
	}
}

func TestCoreClient_GetConfig_WithJWToken(t *testing.T) {
	expectedJWToken := "test-jwt-token"
	expectedTimeout := 5 * time.Second

	client := NewMgcClient(WithJWToken(expectedJWToken),
		WithTimeout(expectedTimeout))

	config := client.GetConfig()

	if config == nil {
		t.Fatal("expected non-nil config")
	}
	if config.JWToken != expectedJWToken {
		t.Errorf("expected JWToken %s, got %s", expectedJWToken, config.JWToken)
	}
	if config.Timeout != expectedTimeout {
		t.Errorf("expected Timeout %v, got %v", expectedTimeout, config.Timeout)
	}
}
