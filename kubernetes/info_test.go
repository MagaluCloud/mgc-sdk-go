package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInfoService_ListFlavors(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list flavors",
			response: `{
				"nodepool": [
					{"name": "gp1.small", "vcpu": 2, "ram": 4096}
				],
				"controlplane": [
					{"name": "cp1.medium", "vcpu": 4, "ram": 8192}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty response",
			response:   `{}`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Info().ListFlavors(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListFlavors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				total := len(result.NodePool) + len(result.ControlPlane)
				if total != tt.want {
					t.Errorf("ListFlavors() got = %d, want %d", total, tt.want)
				}
			}
		})
	}
}

func TestInfoService_ListVersions(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list versions",
			response: `{
				"results": [
					{"version": "v1.30.2", "deprecated": false},
					{"version": "v1.29.5", "deprecated": true}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty response",
			response:   `{"results": []}`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name:       "invalid response",
			response:   `{"invalid": "data"}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Info().ListVersions(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("ListVersions() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestInfoService_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client := testClient(server.URL)
	_, err := client.Info().ListFlavors(ctx)

	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
}
