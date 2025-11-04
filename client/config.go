package client

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// RetryConfig contains configuration for retry behavior.
type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	BackoffFactor   float64
}

// Config contains all configuration options for the client.
type Config struct {
	APIKey        string
	JWToken       string
	BaseURL       MgcUrl
	UserAgent     string
	Logger        *slog.Logger
	HTTPClient    *http.Client
	Timeout       time.Duration
	RetryConfig   RetryConfig
	ContentType   string
	CustomHeaders map[string]string
}

// Option is a function type that modifies the client configuration.
// Options are used to customize the client behavior during initialization.
type Option func(*Config)

// WithAPIKey sets the API key for authentication.
// This option is required for all API operations.
func WithAPIKey(key string) Option {
	return func(c *Config) {
		c.APIKey = key
	}
}

// WithBaseURL sets the base URL for API requests.
// This option allows specifying a custom endpoint for the API.
func WithBaseURL(url MgcUrl) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithUserAgent sets the user agent string for HTTP requests.
// This option allows customizing the user agent header.
func WithUserAgent(ua string) Option {
	return func(c *Config) {
		c.UserAgent = ua
	}
}

// WithJWToken sets the JWToken for authentication.
// This option allows specifying a custom JWToken for authentication.
func WithJWToken(token string) Option {
	return func(c *Config) {
		c.JWToken, _ = strings.CutPrefix(token, "Bearer ")
	}
}

// WithLogger sets the logger instance for client operations.
// This option allows customizing logging behavior.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithHTTPClient sets the HTTP client for making requests.
// This option allows using a custom HTTP client with specific settings.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithTimeout sets the timeout for HTTP requests.
// This option controls how long to wait for responses.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithRetryConfig sets the retry configuration for failed requests.
// This option allows customizing retry behavior with exponential backoff.
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

// WithCustomHeader adds a custom HTTP header to all requests.
// This option allows adding additional headers for specific requirements.
func WithCustomHeader(key, value string) Option {
	return func(c *Config) {
		if c.CustomHeaders == nil {
			c.CustomHeaders = make(map[string]string)
		}
		c.CustomHeaders[key] = value
	}
}
