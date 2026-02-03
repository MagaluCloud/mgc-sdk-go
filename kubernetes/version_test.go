package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestVersionService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       *VersionListOptions
		response   string
		statusCode int
		want       []Version
		wantErr    bool
	}{
		{
			name: "successful list versions including deprecated versions",
			opts: &VersionListOptions{
				IncludeDeprecated: true,
			},
			response: `{
				"results": [
					{"version": "v1.30.2", "deprecated": false},
					{"version": "v1.29.5", "deprecated": true}
				]
			}`,
			statusCode: http.StatusOK,
			want: []Version{
				{Version: "v1.30.2", Deprecated: false},
				{Version: "v1.29.5", Deprecated: true},
			},
			wantErr: false,
		},
		{
			name: "successful list versions excluding deprecated versions",
			opts: &VersionListOptions{
				IncludeDeprecated: false,
			},
			response: `{
				"results": [
					{"version": "v1.30.2", "deprecated": false},
					{"version": "v1.29.5", "deprecated": true}
				]
			}`,
			statusCode: http.StatusOK,
			want: []Version{
				{Version: "v1.30.2", Deprecated: false},
			},
			wantErr: false,
		},
		{
			name: "successful list versions without opts excludes deprecated versions",
			opts: nil,
			response: `{
				"results": [
					{"version": "v1.30.2", "deprecated": false},
					{"version": "v1.29.5", "deprecated": true}
				]
			}`,
			statusCode: http.StatusOK,
			want: []Version{
				{Version: "v1.30.2", Deprecated: false},
			},
			wantErr: false,
		},
		{
			name: "successful list versions without IncludeDeprecated parameter excludes deprecated versions",
			opts: &VersionListOptions{},
			response: `{
				"results": [
					{"version": "v1.30.2", "deprecated": false},
					{"version": "v1.29.5", "deprecated": true}
				]
			}`,
			statusCode: http.StatusOK,
			want: []Version{
				{Version: "v1.30.2", Deprecated: false},
			},
			wantErr: false,
		},
		{
			name: "empty response",
			opts: &VersionListOptions{
				IncludeDeprecated: false,
			},
			response:   `{"results": []}`,
			statusCode: http.StatusOK,
			want:       []Version{},
			wantErr:    false,
		},
		{
			name: "invalid response format",
			opts: &VersionListOptions{
				IncludeDeprecated: false,
			},
			response:   `{"invalid": "data"}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "server error",
			opts: &VersionListOptions{
				IncludeDeprecated: false,
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			versions, err := client.Versions().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(versions) != len(tt.want) {
					t.Errorf("List() returned %d versions, want %d", len(versions), len(tt.want))
				}

				for i, v := range versions {
					if v.Version != tt.want[i].Version || v.Deprecated != tt.want[i].Deprecated {
						t.Errorf("Version mismatch at index %d: got %+v, want %+v", i, v, tt.want[i])
					}
				}
			}
		})
	}
}

func TestVersionService_List_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client := testClient(server.URL)
	_, err := client.Versions().List(ctx, nil)

	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
}
