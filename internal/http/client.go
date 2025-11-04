package mgc_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/internal/retry"
	"gopkg.in/yaml.v3"
)

// NewRequestFunc is a function that creates a new HTTP request.
type NewRequestFunc func(ctx context.Context, method, path string, body any) (*http.Request, error)

// NewRequest creates a new HTTP request with the given method, path, and body.
// It returns the request and an error if the request creation fails.
func NewRequest[T any](c *client.Config, ctx context.Context, method, path string, body *T) (*http.Request, error) {
	c.Logger.Debug("creating new request",
		"method", method,
		"path", path,
		"hasBody", body != nil)

	url := c.BaseURL.String() + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			c.Logger.Error("failed to marshal request body",
				"error", err,
				"method", method,
				"path", path)
			return nil, fmt.Errorf("error marshalling body: %w", err)
		}
		bodyReader = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		c.Logger.Error("failed to create request",
			"error", err,
			"url", url,
			"method", method)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if requestIDVal := ctx.Value(client.RequestIDKey); requestIDVal != nil {
		if requestID, ok := requestIDVal.(string); ok {
			c.Logger.Info("X-Request-ID found in context", "requestID", requestID)
			req.Header.Set("X-Request-ID", requestID)
		} else {
			c.Logger.Warn("X-Request-ID in context is not a string")
		}
	}

	c.Logger.Debug("setting request headers",
		"apiKey", "redacted",
		"userAgent", c.UserAgent)

	if c.JWToken != "" {
		req.Header.Set("Authorization", c.JWToken)
	}
	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", c.ContentType)

	if c.CustomHeaders != nil {
		for k, v := range c.CustomHeaders {
			req.Header.Set(k, v)
			c.Logger.Debug("Request with custom header", "key", k, "value", v)
		}
	}

	return req, nil
}

// Do executes an HTTP request and processes the response.
// If v is provided, the response body will be JSON decoded into it.
// Returns the parsed response and an error if the request fails,
// the response status is not 2xx, or if there are JSON decoding issues.
func Do[T any](c *client.Config, ctx context.Context, req *http.Request, v *T) (*T, error) {
	c.Logger.Debug("starting request execution",
		"method", req.Method,
		"url", req.URL.String(),
		"expectResponse", v != nil)

	if c.HTTPClient == nil {
		return nil, fmt.Errorf("HTTP client is nil")
	}

	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}

	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}

	var lastError error
	for attempt := range c.RetryConfig.MaxAttempts {
		if attempt > 0 {
			backoff := retry.GetNextBackoff(attempt-1, c.RetryConfig.BackoffFactor, c.RetryConfig.InitialInterval, c.RetryConfig.MaxInterval)
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}

		clonedReq := req.Clone(ctx)
		if len(bodyBytes) > 0 {
			clonedReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		c.Logger.Info("making request",
			"method", clonedReq.Method,
			"url", clonedReq.URL.String(),
			"attempt", attempt+1)

		resp, err := c.HTTPClient.Do(clonedReq)
		if err != nil {
			lastError = err
			continue
		}

		defer resp.Body.Close()

		if xRequestID := resp.Header.Get("X-Request-ID"); xRequestID != "" {
			c.Logger.Info("X-Request-ID received in response", "requestID", xRequestID)
		} else {
			c.Logger.Info("X-Request-ID not found in response")
		}

		if xTraceID := resp.Header.Get("X-Mgc-Trace-Id"); xTraceID != "" {
			c.Logger.Info("X-Mgc-Trace-ID received in response", "mgcTraceID", xTraceID)
		} else {
			c.Logger.Info("X-Mgc-Trace-ID not found in response")
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			lastError = client.NewHTTPError(resp)
			if !retry.ShouldRetry(resp.StatusCode) {
				return nil, lastError
			}
			continue
		}

		if v != nil && resp.StatusCode != http.StatusNoContent {
			ct := resp.Header.Get("Content-Type")
			if strings.Contains(ct, "application/x-yaml") || strings.Contains(ct, "application/yaml") {
				return decodeYamlResponse(resp, v)
			}
			// JSON is the default
			return decodeJsonResponse(resp, v)
		}

		return nil, nil
	}

	return nil, &client.RetryError{LastError: lastError, Retries: c.RetryConfig.MaxAttempts}
}

func decodeYamlResponse[T any](resp *http.Response, v *T) (*T, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var checkNull any
	if err := yaml.Unmarshal(body, &checkNull); err != nil {
		return nil, fmt.Errorf("error validating null response: %w", err)
	}
	if checkNull == nil {
		return nil, fmt.Errorf("response body is null")
	}

	if err := yaml.Unmarshal(body, v); err != nil {
		return nil, fmt.Errorf("error decoding yaml response: %w", err)
	}

	return v, nil
}

func decodeJsonResponse[T any](resp *http.Response, v *T) (*T, error) {
	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var checkNull any
	if err := json.Unmarshal(raw, &checkNull); err != nil {
		return nil, fmt.Errorf("error validating null response: %w", err)
	}
	if checkNull == nil {
		return nil, fmt.Errorf("response body is null")
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	if err := decoder.Decode(v); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return v, nil
}

// ExecuteSimpleRequestWithRespBody handles HTTP requests that require response body parsing
func ExecuteSimpleRequestWithRespBody[T any](
	ctx context.Context,
	reqf NewRequestFunc,
	configs *client.Config,
	method string,
	path string,
	body any,
	queryParams url.Values,
) (*T, error) {
	req, err := reqf(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		req.URL.RawQuery = queryParams.Encode()
	}

	var resType T
	result, err := Do(configs, ctx, req, &resType)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ExecuteSimpleRequest handles HTTP requests that do not require response body parsing
func ExecuteSimpleRequest(
	ctx context.Context,
	reqf NewRequestFunc,
	configs *client.Config,
	method string,
	path string,
	body any,
	queryParams url.Values,
) error {
	req, err := reqf(ctx, method, path, body)
	if err != nil {
		return err
	}

	if queryParams != nil {
		req.URL.RawQuery = queryParams.Encode()
	}

	_, err = Do[any](configs, ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
