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
			apiKey:  "test-api-key",
			opts:    nil,
			wantErr: false,
		},
		{
			name:   "client with custom options",
			apiKey: "test-api-key",
			opts: []Option{
				WithBaseURL(BrNe1),
				WithTimeout(5 * time.Second),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMgcClient(tt.apiKey, tt.opts...)
			if client == nil {
				t.Error("expected non-nil client")
				return
			}
			if client.config == nil {
				t.Error("expected non-nil config")
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

	client := NewMgcClient(expectedAPIKey,
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
