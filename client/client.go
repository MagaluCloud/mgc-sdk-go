package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/internal/retry"
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
		BaseURL:    BrNe1,
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

func (c *CoreClient) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	c.config.Logger.Debug("creating new request",
		"method", method,
		"path", path,
		"hasBody", body != nil)

	url := c.config.BaseURL.String() + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			c.config.Logger.Error("failed to marshal request body",
				"error", err,
				"method", method,
				"path", path)
			return nil, fmt.Errorf("error marshalling body: %w", err)
		}
		bodyReader = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		c.config.Logger.Error("failed to create request",
			"error", err,
			"url", url,
			"method", method)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if requestIDVal := ctx.Value(RequestIDKey); requestIDVal != nil {
		if requestID, ok := requestIDVal.(string); ok {
			c.config.Logger.Info("X-Request-ID found in context", "requestID", requestID)
			req.Header.Set("X-Request-ID", requestID)
		} else {
			c.config.Logger.Warn("X-Request-ID in context is not a string")
		}
	}

	c.config.Logger.Debug("setting request headers",
		"apiKey", "redacted",
		"userAgent", c.config.UserAgent)

	req.Header.Set("X-API-Key", c.config.APIKey)
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Do executes an HTTP request and processes the response.
// If v is provided, the response body will be JSON decoded into it.
// Returns the parsed response and an error if the request fails,
// the response status is not 2xx, or if there are JSON decoding issues.
func (c *CoreClient) Do(ctx context.Context, req *http.Request, v any) (any, error) {
	c.config.Logger.Debug("starting request execution",
		"method", req.Method,
		"url", req.URL.String(),
		"expectResponse", v != nil)

	if c.config.HTTPClient == nil {
		return nil, fmt.Errorf("HTTP client is nil")
	}

	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	var lastError error
	for attempt := 0; attempt < c.config.RetryConfig.MaxAttempts; attempt++ {
		if attempt > 0 {
			backoff := retry.GetNextBackoff(attempt-1, c.config.RetryConfig.BackoffFactor, c.config.RetryConfig.InitialInterval, c.config.RetryConfig.MaxInterval)
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}

		c.config.Logger.Info("making request",
			"method", req.Method,
			"url", req.URL.String(),
			"attempt", attempt+1)

		resp, err := c.config.HTTPClient.Do(req.Clone(ctx))
		if err != nil {
			lastError = err
			continue
		}

		defer resp.Body.Close()

		if xRequestID := resp.Header.Get("X-Request-ID"); xRequestID != "" {
			c.config.Logger.Info("X-Request-ID received in response", "requestID", xRequestID)
		} else {
			c.config.Logger.Info("X-Request-ID not found in response")
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			lastError = NewHTTPError(resp)
			if !retry.ShouldRetry(resp.StatusCode) {
				return nil, lastError
			}
			continue
		}

		if v != nil && resp.StatusCode != http.StatusNoContent {
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
				return nil, fmt.Errorf("error decoding response: %w", err)
			}
			return v, nil
		}

		return nil, nil
	}

	return nil, fmt.Errorf("max retry attempts reached: %w", lastError)
}
