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
			opts: ListOptions{
				Limit:  helpers.IntPtr(2),
				Offset: helpers.IntPtr(1),
				Sort:   helpers.StrPtr("name"),
				Expand: []string{"controlplane", "nodepool"},
			},
			name:       "invalid response format",
			response:   `{"invalid": "`,
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
				total := len((*result)[0].NodePool) + len((*result)[0].ControlPlane)
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

func TestFlavorService_List_InvalidOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := testClient(server.URL)
	_, err := client.Flavors().List(context.Background(), ListOptions{
		Limit: helpers.IntPtr(-1),
	})

	if err == nil {
		t.Error("Esperado erro com opções inválidas")
	}
}

func TestFlavorService_List_EmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results": []}`))
	}))
	defer server.Close()

	client := testClient(server.URL)
	result, err := client.Flavors().List(context.Background(), ListOptions{})

	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	if len(*result) != 0 {
		t.Errorf("Esperado 0 resultados, obtido %d", len(*result))
	}
}

func TestFlavorService_List_AuthorizationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := testClient(server.URL)
	_, err := client.Flavors().List(context.Background(), ListOptions{})

	if err == nil {
		t.Error("Esperado erro de autorização")
	}
}

func TestFlavorService_List_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalid": "json`)) // JSON malformado
	}))
	defer server.Close()

	client := testClient(server.URL)
	_, err := client.Flavors().List(context.Background(), ListOptions{})

	if err == nil {
		t.Error("Esperado erro de parsing JSON")
	}
}

func TestFlavorService_List_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		response string
	}{
		{
			name:     "empty node pool",
			response: `{"results": [{"nodepool": [], "controlplane": []}]}`,
		},
		{
			name:     "mixed flavors",
			response: `{"results": [{"nodepool": [{"name":"n1"}], "controlplane": [{"name":"c1"}]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			_, err := client.Flavors().List(context.Background(), ListOptions{})

			if err != nil {
				t.Errorf("Erro inesperado: %v", err)
			}
		})
	}
}
