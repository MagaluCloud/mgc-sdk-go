package client

import (
	"log/slog"
	"net/http"
	"time"
)

const (
	RequestIDKey XRequestID = "x-request-id"

	DefaultUserAgent = "mgc-sdk-go"

	DefaultMaxAttempts     = 3
	DefaultInitialInterval = 1 * time.Second
	DefaultMaxInterval     = 30 * time.Second
	DefaultBackoffFactor   = 2.0
	DefaultTimeout         = 15 * time.Minute
)

type (
	XRequestID string
)

type CoreClient struct {
	config *Config
}

func NewMgcClient(apiKey string, opts ...Option) *CoreClient {
	cfg := &Config{
		HTTPClient: http.DefaultClient,
		Logger:     slog.Default(),
		APIKey:     apiKey,
		UserAgent:  DefaultUserAgent,
		BaseURL:    BrSe1,
		Timeout:    DefaultTimeout,
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
	return &CoreClient{config: cfg}
}

func (c *CoreClient) GetConfig() *Config {
	return c.config
}
