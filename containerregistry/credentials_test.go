package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

// Helper functions
func testClient(baseURL string) *ContainerRegistryClient {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core)
}

func TestCredentialsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       *CredentialsResponse
		wantErr    bool
	}{
		{
			name: "successful credentials fetch",
			response: `{
				"username": "test-user",
				"password": "test-pass",
				"email": "test@example.com"
			}`,
			statusCode: http.StatusOK,
			want: &CredentialsResponse{
				Username: "test-user",
				Password: "test-pass",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			response:   `{"username": "broken"`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "empty response",
			response:   "",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "empty response body",
			response:   "",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "nil response",
			response:   "null",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "bad request",
			response:   `{"error": "bad request"}`,
			statusCode: http.StatusBadRequest,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "service unavailable",
			response:   `{"error": "service unavailable"}`,
			statusCode: http.StatusServiceUnavailable,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("=== " + tt.name + " ===")
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Credentials().Get(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.Username != tt.want.Username {
					t.Errorf("Username = %v, want %v", got.Username, tt.want.Username)
				}
				if got.Password != tt.want.Password {
					t.Errorf("Password = %v, want %v", got.Password, tt.want.Password)
				}
				if got.Email != tt.want.Email {
					t.Errorf("Email = %v, want %v", got.Email, tt.want.Email)
				}
			}
		})
	}
}

func TestCredentialsService_ResetPassword(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       *CredentialsResponse
		wantErr    bool
	}{
		{
			name: "successful password reset",
			response: `{
				"username": "test-user",
				"password": "new-password",
				"email": "test@example.com"
			}`,
			statusCode: http.StatusOK,
			want: &CredentialsResponse{
				Username: "test-user",
				Password: "new-password",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			response:   `{"username": "broken"`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "empty response body",
			response:   "",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "nil response",
			response:   "null",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "too many requests",
			response:   `{"error": "rate limit exceeded"}`,
			statusCode: http.StatusTooManyRequests,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "gateway timeout",
			response:   `{"error": "gateway timeout"}`,
			statusCode: http.StatusGatewayTimeout,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "invalid json structure",
			response:   `{"username": 123, "password": true}`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Credentials().ResetPassword(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.Username != tt.want.Username {
					t.Errorf("Username = %v, want %v", got.Username, tt.want.Username)
				}
				if got.Password != tt.want.Password {
					t.Errorf("Password = %v, want %v", got.Password, tt.want.Password)
				}
				if got.Email != tt.want.Email {
					t.Errorf("Email = %v, want %v", got.Email, tt.want.Email)
				}
			}
		})
	}
}

func TestCredentialsService_Concurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"username": "test-user", "password": "test-pass", "email": "test@example.com"}`))
	}))
	defer server.Close()

	client := testClient(server.URL)
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.Credentials().Get(ctx)
			if err != nil {
				t.Errorf("concurrent Get() error = %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestCredentialsService_InvalidServerResponses(t *testing.T) {
	tests := []struct {
		name           string
		serverBehavior func(w http.ResponseWriter, r *http.Request)
		wantErr        bool
	}{
		{
			name: "connection closed",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				hj, ok := w.(http.Hijacker)
				if !ok {
					t.Fatal("couldn't hijack connection")
				}
				conn, _, err := hj.Hijack()
				if err != nil {
					t.Fatal(err)
				}
				conn.Close()
			},
			wantErr: true,
		},
		{
			name: "partial response",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"username": "test"`)) // Resposta JSON incompleta
			},
			wantErr: true,
		},
		{
			name: "wrong content type",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"username": "test"}`))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverBehavior))
			defer server.Close()

			client := testClient(server.URL)
			_, err := client.Credentials().Get(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
