package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestFlavorService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list flavors with pagination",
			opts: ListOptions{
				Limit:  helpers.IntPtr(2),
				Offset: helpers.IntPtr(1),
			},
			response: `{
				"results": [
					{
						"nodepool": [{"name": "gp1.small"}],
						"controlplane": [{"name": "cp1.medium"}]
					}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "invalid response format",
			response:   `{"invalid": "data"}`,
			statusCode: http.StatusOK,
			wantErr:    true,
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
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Flavors().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				total := len(result.Results[0].NodePool) + len(result.Results[0].ControlPlane)
				if total != tt.want {
					t.Errorf("List() got = %d, want %d", total, tt.want)
				}
			}
		})
	}
}

func TestFlavorService_List_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client := testClient(server.URL)
	_, err := client.Flavors().List(ctx, ListOptions{})

	if err == nil {
		t.Error("Expected context timeout error, got nil")
	}
}
