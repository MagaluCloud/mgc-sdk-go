package client

import (
	"log/slog"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	BackoffFactor   float64
}

type Config struct {
	APIKey        string
	BaseURL       MgcUrl
	UserAgent     string
	Logger        *slog.Logger
	HTTPClient    *http.Client
	Timeout       time.Duration
	RetryConfig   RetryConfig
	ContentType   string
	CustomHeaders map[string]string
}

type Option func(*Config)

func WithAPIKey(key string) Option {
	return func(c *Config) {
		c.APIKey = key
	}
}

func WithBaseURL(url MgcUrl) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

func WithUserAgent(ua string) Option {
	return func(c *Config) {
		c.UserAgent = ua
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithRetryConfig(maxAttempts int, initialInterval, maxInterval time.Duration, backoffFactor float64) Option {
	return func(c *Config) {
		c.RetryConfig = RetryConfig{
			MaxAttempts:     maxAttempts,
			InitialInterval: initialInterval,
			MaxInterval:     maxInterval,
			BackoffFactor:   backoffFactor,
		}
	}
}

func WithCustomHeader(key, value string) Option {
	return func(c *Config) {
		if c.CustomHeaders == nil {
			c.CustomHeaders = make(map[string]string)
		}
		c.CustomHeaders[key] = value
	}
}
