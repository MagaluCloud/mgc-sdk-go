package containerregistry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
func TestRegistriesService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    *RegistryRequest
		response   string
		statusCode int
		want       *RegistryResponse
		wantErr    bool
	}{
		{
			name: "successful create",
			request: &RegistryRequest{
				Name: "test-registry",
			},
			response: `{
				"id": "reg-123",
				"name": "test-registry",
				"storage_usage_bytes": 1024,
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			want: &RegistryResponse{
				ID:        "reg-123",
				Name:      "test-registry",
				Storage:   1024,
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			request: &RegistryRequest{
				Name: "",
			},
			response:   `{"error": "name cannot be empty"}`,
			statusCode: http.StatusBadRequest,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "malformed json response",
			request:    &RegistryRequest{Name: "test"},
			response:   `{"id": "reg-123"`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "server error",
			request:    &RegistryRequest{Name: "test"},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
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
			got, err := client.Registries().Create(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("Create() got ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Name != tt.want.Name {
					t.Errorf("Create() got Name = %v, want %v", got.Name, tt.want.Name)
				}
			}
		})
	}
}

func TestRegistriesService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListOptions
		response   string
		statusCode int
		want       *ListRegistriesResponse
		wantErr    bool
	}{
		{
			name: "successful list",
			opts: ListOptions{},
			response: `{
				"results": [
					{
						"id": "reg-123",
						"name": "test-registry",
						"storage_usage_bytes": 1024,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListRegistriesResponse{
				Registries: []RegistryResponse{
					{
						ID:        "reg-123",
						Name:      "test-registry",
						Storage:   1024,
						CreatedAt: "2024-01-01T00:00:00Z",
						UpdatedAt: "2024-01-01T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "empty response",
			opts:       ListOptions{},
			response:   "",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			opts:       ListOptions{},
			response:   `{"results": [{"id": "reg-123"`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			opts:       ListOptions{},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			want:       nil,
			wantErr:    true,
		},
		{
			name: "list with pagination",
			opts: ListOptions{
				Limit:  intPtr(10),
				Offset: intPtr(20),
			},
			response: `{
				"results": [
					{
						"id": "reg-123",
						"name": "test-registry",
						"storage_usage_bytes": 1024,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListRegistriesResponse{
				Registries: []RegistryResponse{
					{
						ID:        "reg-123",
						Name:      "test-registry",
						Storage:   1024,
						CreatedAt: "2024-01-01T00:00:00Z",
						UpdatedAt: "2024-01-01T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "list with sorting",
			opts: ListOptions{
				Sort: strPtr("name:asc"),
			},
			response: `{
				"results": [
					{
						"id": "reg-123",
						"name": "test-registry",
						"storage_usage_bytes": 1024,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListRegistriesResponse{
				Registries: []RegistryResponse{
					{
						ID:        "reg-123",
						Name:      "test-registry",
						Storage:   1024,
						CreatedAt: "2024-01-01T00:00:00Z",
						UpdatedAt: "2024-01-01T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "list with expand",
			opts: ListOptions{
				Expand: []string{"storage"},
			},
			response: `{
				"results": [
					{
						"id": "reg-123",
						"name": "test-registry",
						"storage_usage_bytes": 1024,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListRegistriesResponse{
				Registries: []RegistryResponse{
					{
						ID:        "reg-123",
						Name:      "test-registry",
						Storage:   1024,
						CreatedAt: "2024-01-01T00:00:00Z",
						UpdatedAt: "2024-01-01T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			got, err := client.Registries().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if len(got.Registries) != len(tt.want.Registries) {
					t.Errorf("List() got %v registries, want %v", len(got.Registries), len(tt.want.Registries))
				}
			}
		})
	}
}

func TestRegistriesService_Get(t *testing.T) {
	tests := []struct {
		name       string
		registryID string
		response   string
		statusCode int
		want       *RegistryResponse
		wantErr    bool
	}{
		{
			name:       "successful get",
			registryID: "reg-123",
			response: `{
				"id": "reg-123",
				"name": "test-registry",
				"storage_usage_bytes": 1024,
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			want: &RegistryResponse{
				ID:        "reg-123",
				Name:      "test-registry",
				Storage:   1024,
				CreatedAt: "2024-01-01T00:00:00Z",
				UpdatedAt: "2024-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name:       "not found",
			registryID: "nonexistent",
			response:   `{"error": "registry not found"}`,
			statusCode: http.StatusNotFound,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "malformed response",
			registryID: "reg-123",
			response:   `{"id": "reg-123"`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "server error",
			registryID: "reg-123",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			got, err := client.Registries().Get(context.Background(), tt.registryID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("Get() got ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Name != tt.want.Name {
					t.Errorf("Get() got Name = %v, want %v", got.Name, tt.want.Name)
				}
			}
		})
	}
}

func TestRegistriesService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		registryID string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			registryID: "reg-123",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "not found",
			registryID: "nonexistent",
			response:   `{"error": "registry not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			registryID: "reg-123",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "server error",
			registryID: "reg-123",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE method, got %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Registries().Delete(context.Background(), tt.registryID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegistriesService_Concurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results": []}`))
	}))
	defer server.Close()

	client := testClient(server.URL)
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.Registries().List(ctx, ListOptions{})
			if err != nil {
				t.Errorf("concurrent List() error = %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
