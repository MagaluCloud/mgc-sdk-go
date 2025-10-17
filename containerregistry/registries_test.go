package containerregistry

import (
	"context"
	"fmt"
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
		name          string
		opts          RegistryListOptions
		response      string
		statusCode    int
		expectedQuery map[string]string
		want          *ListRegistriesResponse
		wantErr       bool
	}{
		{
			name: "successful list",
			opts: RegistryListOptions{},
			response: `{
				"meta": {
					"page": {
						"count": 1,
						"limit": 20,
						"offset": 0,
						"total": 1
					}
				},
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
			statusCode:    http.StatusOK,
			expectedQuery: map[string]string{},
			want: &ListRegistriesResponse{
				Results: []RegistryResponse{
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
			name:          "empty response",
			opts:          RegistryListOptions{},
			expectedQuery: map[string]string{},
			response:      "",
			statusCode:    http.StatusOK,
			want:          nil,
			wantErr:       true,
		},
		{
			name:          "malformed json",
			opts:          RegistryListOptions{},
			expectedQuery: map[string]string{},
			response:      `{"results": [{"id": "reg-123"`,
			statusCode:    http.StatusOK,
			want:          nil,
			wantErr:       true,
		},
		{
			name:          "unauthorized",
			opts:          RegistryListOptions{},
			expectedQuery: map[string]string{},
			response:      `{"error": "unauthorized"}`,
			statusCode:    http.StatusUnauthorized,
			want:          nil,
			wantErr:       true,
		},
		{
			name: "list with pagination",
			opts: RegistryListOptions{
				Limit:  intPtr(10),
				Offset: intPtr(20),
			},
			response: `{
				"meta": {
					"page": {
						"count": 1,
						"limit": 10,
						"offset": 20,
						"total": 100
					}
				},
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
			expectedQuery: map[string]string{
				"_limit":  "10",
				"_offset": "20",
			},
			want: &ListRegistriesResponse{
				Results: []RegistryResponse{
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
			opts: RegistryListOptions{
				RegistryFilterOptions: RegistryFilterOptions{
					Sort: strPtr("name:asc"),
				},
			},
			response: `{
				"meta": {
					"page": {
						"count": 1,
						"limit": 20,
						"offset": 0,
						"total": 1
					}
				},
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
			expectedQuery: map[string]string{
				"_sort": "name:asc",
			},
			want: &ListRegistriesResponse{
				Results: []RegistryResponse{
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
			name: "list with multiple options",
			opts: RegistryListOptions{
				Limit:  intPtr(20),
				Offset: intPtr(10),
				RegistryFilterOptions: RegistryFilterOptions{
					Sort: strPtr("created_at"),
				},
			},
			response: `{
				"meta": {
					"page": {
						"count": 1,
						"limit": 20,
						"offset": 10,
						"total": 50
					}
				},
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
			expectedQuery: map[string]string{
				"_limit":  "20",
				"_offset": "10",
				"_sort":   "created_at",
			},
			want: &ListRegistriesResponse{
				Results: []RegistryResponse{
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
				query := r.URL.Query()
				for key, expectedValue := range tt.expectedQuery {
					if actualValue := query.Get(key); actualValue != expectedValue {
						t.Errorf("expected query param %s=%s, got %s", key, expectedValue, actualValue)
					}
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
				if len(got.Results) != len(tt.want.Results) {
					t.Errorf("List() got %v registries, want %v", len(got.Results), len(tt.want.Results))
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
			_, err := client.Registries().List(ctx, RegistryListOptions{})
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

func TestRegistriesService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		filterOpts RegistryFilterOptions
		responses  []string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "successful list all - single page",
			filterOpts: RegistryFilterOptions{},
			responses: []string{
				`{
					"meta": {
						"page": {
							"count": 2,
							"limit": 50,
							"offset": 0,
							"total": 2
						}
					},
					"results": [
						{
							"id": "reg-123",
							"name": "test-registry-1",
							"storage_usage_bytes": 1024,
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						},
						{
							"id": "reg-456",
							"name": "test-registry-2",
							"storage_usage_bytes": 2048,
							"created_at": "2024-01-02T00:00:00Z",
							"updated_at": "2024-01-02T00:00:00Z"
						}
					]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "successful list all - multiple pages",
			filterOpts: RegistryFilterOptions{},
			responses: []string{
				`{
					"meta": {
						"page": {
							"count": 50,
							"limit": 50,
							"offset": 0,
							"total": 75
						}
					},
					"results": [` + generateRegistryJSONArray(50) + `]
				}`,
				`{
					"meta": {
						"page": {
							"count": 25,
							"limit": 50,
							"offset": 50,
							"total": 75
						}
					},
					"results": [` + generateRegistryJSONArray(25) + `]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  75,
			wantErr:    false,
		},
		{
			name:       "empty results",
			filterOpts: RegistryFilterOptions{},
			responses: []string{
				`{
					"meta": {
						"page": {
							"count": 0,
							"limit": 50,
							"offset": 0,
							"total": 0
						}
					},
					"results": []
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "list all with sort",
			filterOpts: RegistryFilterOptions{Sort: strPtr("name")},
			responses: []string{
				`{
					"meta": {
						"page": {
							"count": 1,
							"limit": 50,
							"offset": 0,
							"total": 1
						}
					},
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
			},
			statusCode: http.StatusOK,
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "error on first page",
			filterOpts: RegistryFilterOptions{},
			responses:  []string{`{"error": "internal server error"}`},
			statusCode: http.StatusInternalServerError,
			wantCount:  0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if requestCount < len(tt.responses) {
					w.Write([]byte(tt.responses[requestCount]))
					requestCount++
				}
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Registries().ListAll(context.Background(), tt.filterOpts)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("ListAll() got %v registries, want %v", len(got), tt.wantCount)
			}
		})
	}
}

// Helper function to generate registry JSON array for testing pagination
func generateRegistryJSONArray(count int) string {
	var registries []string
	for i := 0; i < count; i++ {
		registries = append(registries, fmt.Sprintf(`{
			"id": "reg-%d",
			"name": "test-registry-%d",
			"storage_usage_bytes": %d,
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z"
		}`, i, i, 1024*(i+1)))
	}
	result := ""
	for i, reg := range registries {
		if i > 0 {
			result += ","
		}
		result += reg
	}
	return result
}
