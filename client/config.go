// Package client provides the core client functionality for the MagaluCloud SDK.
// This package contains the main client implementation, configuration options, and error handling.
package client

import (
	"log/slog"
	"net/http"
	"time"
)

// RetryConfig contains configuration for retry behavior.
// This structure defines how the client should handle retries when requests fail.
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int
	// InitialInterval is the initial delay between retry attempts
	InitialInterval time.Duration
	// MaxInterval is the maximum delay between retry attempts
	MaxInterval time.Duration
	// BackoffFactor is the exponential backoff multiplier
	BackoffFactor float64
}

// Config contains all configuration options for the client.
// This structure holds all the settings that control the client's behavior.
type Config struct {
	// APIKey is the authentication key for MagaluCloud services
	APIKey string
	// BaseURL is the base URL for API requests
	BaseURL MgcUrl
	// UserAgent is the user agent string for HTTP requests
	UserAgent string
	// Logger is the logger instance for client operations
	Logger *slog.Logger
	// HTTPClient is the HTTP client for making requests
	HTTPClient *http.Client
	// Timeout is the timeout for HTTP requests
	Timeout time.Duration
	// RetryConfig contains retry behavior configuration
	RetryConfig RetryConfig
	// ContentType is the content type for HTTP requests
	ContentType string
	// CustomHeaders contains additional HTTP headers to include in requests
	CustomHeaders map[string]string
}

// Option is a function type that modifies the client configuration.
// Options are used to customize the client behavior during initialization.
type Option func(*Config)

// WithAPIKey sets the API key for authentication.
// This option is required for all API operations.
//
// Parameters:
//   - key: The API key for MagaluCloud services
//
// Returns:
//   - Option: A configuration option function
func WithAPIKey(key string) Option {
	return func(c *Config) {
		c.APIKey = key
	}
}

// WithBaseURL sets the base URL for API requests.
// This option allows specifying a custom endpoint for the API.
//
// Parameters:
//   - url: The base URL for API requests
//
// Returns:
//   - Option: A configuration option function
func WithBaseURL(url MgcUrl) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithUserAgent sets the user agent string for HTTP requests.
// This option allows customizing the user agent header.
//
// Parameters:
//   - ua: The user agent string
//
// Returns:
//   - Option: A configuration option function
func WithUserAgent(ua string) Option {
	return func(c *Config) {
		c.UserAgent = ua
	}
}

// WithLogger sets the logger instance for client operations.
// This option allows customizing logging behavior.
//
// Parameters:
//   - logger: The logger instance to use
//
// Returns:
//   - Option: A configuration option function
func WithLogger(logger *slog.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithHTTPClient sets the HTTP client for making requests.
// This option allows using a custom HTTP client with specific settings.
//
// Parameters:
//   - client: The HTTP client to use
//
// Returns:
//   - Option: A configuration option function
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithTimeout sets the timeout for HTTP requests.
// This option controls how long to wait for responses.
//
// Parameters:
//   - timeout: The timeout duration for HTTP requests
//
// Returns:
//   - Option: A configuration option function
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithRetryConfig sets the retry configuration for failed requests.
// This option allows customizing retry behavior with exponential backoff.
//
// Parameters:
//   - maxAttempts: Maximum number of retry attempts
//   - initialInterval: Initial delay between retries
//   - maxInterval: Maximum delay between retries
//   - backoffFactor: Exponential backoff multiplier
//
// Returns:
//   - Option: A configuration option function
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
//
// Parameters:
//   - key: The header key
//   - value: The header value
//
// Returns:
//   - Option: A configuration option function
func WithCustomHeader(key, value string) Option {
	return func(c *Config) {
		if c.CustomHeaders == nil {
			c.CustomHeaders = make(map[string]string)
		}
		c.CustomHeaders[key] = value
	}
}
