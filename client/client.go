// Package client provides the core client functionality for the MagaluCloud SDK.
// This package contains the main client implementation, configuration options, and error handling.
package client

import (
	"log/slog"
	"net/http"
	"time"
)

// Default configuration constants for the client.
const (
	RequestIDKey           XRequestID = "x-request-id"
	DefaultUserAgent                  = "mgc-sdk-go"
	DefaultMaxAttempts                = 3
	DefaultInitialInterval            = 1 * time.Second
	DefaultMaxInterval                = 30 * time.Second
	DefaultBackoffFactor              = 2.0
	DefaultTimeout                    = 15 * time.Minute
)

// XRequestID represents a request ID type for tracking requests.
type XRequestID string

// CoreClient represents the main client for interacting with MagaluCloud APIs.
// It encapsulates the configuration and provides methods for making HTTP requests.
type CoreClient struct {
	config Config
}

// NewMgcClient creates a new instance of CoreClient with the specified API key and options.
// The client is configured with sensible defaults and can be customized using the provided options.
func NewMgcClient(opts ...Option) *CoreClient {
	cfg := &Config{
		HTTPClient:  http.DefaultClient,
		Logger:      slog.Default(),
		APIKey:      "",
		JWToken:     "",
		UserAgent:   DefaultUserAgent,
		BaseURL:     BrSe1,
		Timeout:     DefaultTimeout,
		ContentType: "application/json",
		RetryConfig: RetryConfig{
			MaxAttempts:     DefaultMaxAttempts,
			InitialInterval: DefaultInitialInterval,
			MaxInterval:     DefaultMaxInterval,
			BackoffFactor:   DefaultBackoffFactor,
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	cfg.Logger.Debug("creating new core client",
		"baseURL", cfg.BaseURL.String(),
		"userAgent", cfg.UserAgent)
	return &CoreClient{config: *cfg}
}

// GetConfig returns a pointer to the client's configuration.
// This method allows access to the current configuration for inspection or modification.
func (c *CoreClient) GetConfig() *Config {
	return &c.config
}
