package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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
			client := New(tt.apiKey, tt.opts...)
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

func TestCoreClient_NewRequest(t *testing.T) {
	client := New("test-api-key")

	tests := []struct {
		name     string
		method   string
		path     string
		body     interface{}
		ctxFunc  func() context.Context
		wantErr  bool
		checkReq func(*testing.T, *http.Request)
	}{
		{
			name:    "valid GET request",
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			ctxFunc: context.Background,
			wantErr: false,
			checkReq: func(t *testing.T, req *http.Request) {
				if req.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", req.Method)
				}
			},
		},
		{
			name:   "valid POST request with body",
			method: http.MethodPost,
			path:   "/test",
			body:   mockRequest{Data: "test"},
			ctxFunc: func() context.Context {
				return context.WithValue(context.Background(), RequestIDKey, "test-request-id")
			},
			wantErr: false,
			checkReq: func(t *testing.T, req *http.Request) {
				if req.Header.Get("X-Request-ID") != "test-request-id" {
					t.Error("expected X-Request-ID header")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := client.NewRequest(tt.ctxFunc(), tt.method, tt.path, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.checkReq != nil {
				tt.checkReq(t, req)
			}
		})
	}
}

func TestCoreClient_Do(t *testing.T) {
	tests := []struct {
		name           string
		setupServer    func() *httptest.Server
		setupContext   func() context.Context
		expectedResult interface{}
		wantErr        bool
	}{
		{
			name: "successful request",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(mockResponse{Message: "success"})
				}))
			},
			setupContext: context.Background,
			wantErr:      false,
		},
		{
			name: "retry on 500",
			setupServer: func() *httptest.Server {
				attempts := 0
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if attempts < 2 {
						attempts++
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					json.NewEncoder(w).Encode(mockResponse{Message: "success"})
				}))
			},
			setupContext: context.Background,
			wantErr:      false,
		},
		{
			name: "timeout",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2 * time.Second)
					json.NewEncoder(w).Encode(mockResponse{Message: "success"})
				}))
			},
			setupContext: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
				return ctx
			},
			wantErr: true,
		},
		{
			name: "retry on 429 too many requests",
			setupServer: func() *httptest.Server {
				attempts := 0
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if attempts < 2 {
						attempts++
						w.WriteHeader(http.StatusTooManyRequests)
						return
					}
					json.NewEncoder(w).Encode(mockResponse{Message: "success"})
				}))
			},
			setupContext: context.Background,
			wantErr:      false,
		},
		{
			name: "bad json response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"message": "invalid json`))
				}))
			},
			setupContext: context.Background,
			wantErr:      true,
		},
		{
			name: "empty response body with 204",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}))
			},
			setupContext: context.Background,
			wantErr:      false,
		},
		{
			name: "context cancellation",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(100 * time.Millisecond)
					json.NewEncoder(w).Encode(mockResponse{Message: "success"})
				}))
			},
			setupContext: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					time.Sleep(50 * time.Millisecond)
					cancel()
				}()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			client := New("test-api-key",
				WithBaseURL(MgcUrl(server.URL)),
				WithTimeout(10*time.Second),
				WithRetryConfig(3,
					100*time.Millisecond,
					500*time.Millisecond,
					1.5,
				))

			ctx := tt.setupContext()
			req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			var response mockResponse
			_, err = client.Do(ctx, req, &response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRetryLogic(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		json.NewEncoder(w).Encode(mockResponse{Message: "success"})
	}))
	defer server.Close()

	client := New("test-api-key",
		WithBaseURL(MgcUrl(server.URL)),
		WithRetryConfig(3, 100*time.Millisecond, 1*time.Second, 2.0))

	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	var response mockResponse
	_, err := client.Do(context.Background(), req, &response)
	if err != nil {
		t.Errorf("Expected successful retry, got error: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRequestHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "test-api-key" {
			t.Error("Missing or incorrect X-API-Key header")
		}
		if r.Header.Get("User-Agent") != DefaultUserAgent {
			t.Error("Missing or incorrect User-Agent header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Missing or incorrect Content-Type header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New("test-api-key", WithBaseURL(MgcUrl(server.URL)))
	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	_, err := client.Do(context.Background(), req, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestResponseStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantRetry  bool
		wantErr    bool
	}{
		{"400 Bad Request", http.StatusBadRequest, false, true},
		{"401 Unauthorized", http.StatusUnauthorized, false, true},
		{"403 Forbidden", http.StatusForbidden, false, true},
		{"404 Not Found", http.StatusNotFound, false, true},
		{"429 Too Many Requests", http.StatusTooManyRequests, true, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true, true},
		{"502 Bad Gateway", http.StatusBadGateway, true, true},
		{"503 Service Unavailable", http.StatusServiceUnavailable, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := New("test-api-key",
				WithBaseURL(MgcUrl(server.URL)),
				WithRetryConfig(
					2,
					50*time.Millisecond,
					100*time.Millisecond,
					1.5,
				))

			req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
			_, err := client.Do(context.Background(), req, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRequest_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    interface{}
		wantErr string
	}{
		{
			name:    "invalid method",
			method:  string([]byte{0x7f}), // Use an invalid character for method
			path:    "/test",
			wantErr: "error creating request",
		},
		{
			name:    "invalid URL characters",
			method:  http.MethodGet,
			path:    string([]byte{0x7f}),
			wantErr: "error creating request",
		},
		{
			name:   "invalid body type",
			method: http.MethodPost,
			path:   "/test",
			body: struct {
				Ch chan int
			}{
				Ch: make(chan int),
			},
			wantErr: "error marshalling body",
		},
	}

	client := New("test-api-key")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.NewRequest(context.Background(), tt.method, tt.path, tt.body)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if err != nil && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestClient_InvalidConfiguration(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() *CoreClient
		testFunc  func(*testing.T, *CoreClient)
	}{
		{
			name: "nil http client",
			setupFunc: func() *CoreClient {
				return New("test-api-key", func(c *Config) {
					c.HTTPClient = nil
				})
			},
			testFunc: func(t *testing.T, client *CoreClient) {
				req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
				_, err := client.Do(context.Background(), req, nil)
				if err == nil || !strings.Contains(err.Error(), "HTTP client is nil") {
					t.Errorf("expected 'HTTP client is nil' error, got %v", err)
				}
			},
		},
		{
			name: "invalid base URL",
			setupFunc: func() *CoreClient {
				return New("test-api-key", WithBaseURL(MgcUrl("http://invalid\x7f.com")))
			},
			testFunc: func(t *testing.T, client *CoreClient) {
				_, err := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
				if err == nil {
					t.Error("expected error with invalid base URL")
				}
			},
		},
		{
			name: "zero timeout",
			setupFunc: func() *CoreClient {
				return New("test-api-key", WithTimeout(0))
			},
			testFunc: func(t *testing.T, client *CoreClient) {
				if client.config.Timeout != 0 {
					t.Error("expected zero timeout")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setupFunc()
			tt.testFunc(t, client)
		})
	}
}

func TestResponseError_Handling(t *testing.T) {
	tests := []struct {
		name         string
		setupServer  func() *httptest.Server
		expectedBody string
		wantErr      bool
	}{
		{
			name: "malformed JSON error response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error": malformed`))
				}))
			},
			wantErr: true,
		},
		{
			name: "empty error response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				}))
			},
			wantErr: true,
		},
		{
			name: "response with invalid content type",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("not json"))
				}))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			client := New("test-api-key", WithBaseURL(MgcUrl(server.URL)))
			req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
			var response interface{}
			_, err := client.Do(context.Background(), req, &response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaxRetryAttemptsReached(t *testing.T) {
	maxAttempts := 3
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := New("test-api-key",
		WithBaseURL(MgcUrl(server.URL)),
		WithRetryConfig(maxAttempts, 10*time.Millisecond, 50*time.Millisecond, 1.5))

	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("expected error after max retry attempts, got nil")
	}

	expectedMsgs := []string{
		"max retry attempts reached",
		"HTTP error: 503",
	}

	for _, msg := range expectedMsgs {
		if !strings.Contains(err.Error(), msg) {
			t.Errorf("expected error containing %q, got: %v", msg, err)
		}
	}

	if attemptCount != maxAttempts {
		t.Errorf("expected %d attempts, got %d", maxAttempts, attemptCount)
	}
}

func TestRequestIDHandling(t *testing.T) {
	requestIDValue := "test-request-id-123"
	requestIDReceived := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestIDReceived = r.Header.Get("X-Request-ID")
		w.Header().Set("X-Request-ID", requestIDReceived)
		json.NewEncoder(w).Encode(mockResponse{Message: "success"})
	}))
	defer server.Close()

	client := New("test-api-key", WithBaseURL(MgcUrl(server.URL)))

	ctx := context.WithValue(context.Background(), RequestIDKey, requestIDValue)
	req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var response mockResponse
	_, err = client.Do(ctx, req, &response)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if requestIDReceived != requestIDValue {
		t.Errorf("Expected X-Request-ID %q, got %q", requestIDValue, requestIDReceived)
	}

	// Test with invalid RequestIDKey type
	ctx = context.WithValue(context.Background(), RequestIDKey, 123) // Invalid type
	req, err = client.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if req.Header.Get("X-Request-ID") != "" {
		t.Error("Expected no X-Request-ID header when invalid type is provided")
	}
}

func TestRequestIDHandling_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		requestIDValue interface{}
		wantHeader    string
		wantLogMsg    string
	}{
		{
			name:           "valid string request ID",
			requestIDValue: "test-id-123",
			wantHeader:    "test-id-123",
			wantLogMsg:    "X-Request-ID found in context",
		},
		{
			name:           "invalid type request ID",
			requestIDValue: 123,
			wantHeader:    "",
			wantLogMsg:    "X-Request-ID in context is not a string",
		},
		{
			name:           "empty string request ID",
			requestIDValue: "",
			wantHeader:    "",
			wantLogMsg:    "X-Request-ID found in context",
		},
		{
			name:           "nil request ID",
			requestIDValue: nil,
			wantHeader:    "",
			wantLogMsg:    "X-Request-ID not found in context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New("test-api-key")
			ctx := context.Background()
			if tt.requestIDValue != nil {
				ctx = context.WithValue(ctx, RequestIDKey, tt.requestIDValue)
			}

			req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			if got := req.Header.Get("X-Request-ID"); got != tt.wantHeader {
				t.Errorf("RequestID header = %q, want %q", got, tt.wantHeader)
			}
		})
	}
}

func TestRequestID_Retries(t *testing.T) {
	requestID := "test-retry-id"
	receivedIDs := make([]string, 0)
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedIDs = append(receivedIDs, r.Header.Get("X-Request-ID"))
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New("test-api-key",
		WithBaseURL(MgcUrl(server.URL)),
		WithRetryConfig(3, 10*time.Millisecond, 50*time.Millisecond, 1.5))

	ctx := context.WithValue(context.Background(), RequestIDKey, requestID)
	req, _ := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	_, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if len(receivedIDs) != 3 {
		t.Errorf("Expected 3 requests with IDs, got %d", len(receivedIDs))
	}
	for i, id := range receivedIDs {
		if id != requestID {
			t.Errorf("Request %d: expected ID %q, got %q", i+1, requestID, id)
		}
	}
}

func TestConcurrentRequests_DifferentRequestIDs(t *testing.T) {
	receivedIDs := make(chan string, 10)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedIDs <- r.Header.Get("X-Request-ID")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New("test-api-key", WithBaseURL(MgcUrl(server.URL)))
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		requestID := fmt.Sprintf("request-%d", i)
		go func(rid string) {
			defer wg.Done()
			ctx := context.WithValue(context.Background(), RequestIDKey, rid)
			req, _ := client.NewRequest(ctx, http.MethodGet, "/test", nil)
			_, err := client.Do(ctx, req, nil)
			if err != nil {
				t.Errorf("Request with ID %s failed: %v", rid, err)
			}
		}(requestID)
	}

	wg.Wait()
	close(receivedIDs)

	receivedMap := make(map[string]bool)
	for id := range receivedIDs {
		receivedMap[id] = true
	}

	for i := 0; i < 5; i++ {
		expectedID := fmt.Sprintf("request-%d", i)
		if !receivedMap[expectedID] {
			t.Errorf("Request ID %s not received by server", expectedID)
		}
	}
}
