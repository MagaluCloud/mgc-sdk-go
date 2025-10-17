package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRepositoriesService_List(t *testing.T) {
	tests := []struct {
		name          string
		registryID    string
		opts          RepositoryListOptions
		response      string
		statusCode    int
		expectedQuery map[string]string
		want          *RepositoriesResponse
		wantErr       bool
	}{
		{
			name:       "successful list repositories",
			registryID: "reg-123",
			response: `{
				"meta": {
					"page": {
						"count": 2,
						"limit": 0,
						"offset": 0,
						"total": 2
					}
				},
				"results": [
					{
						"registry_name": "test-registry",
						"name": "repo1",
						"image_count": 5,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					},
					{
						"registry_name": "test-registry",
						"name": "repo2",
						"image_count": 3,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &RepositoriesResponse{
				Meta: Meta{Page: Page{Count: 2, Limit: 0, Offset: 0, Total: 2}},
				Results: []RepositoryResponse{
					{
						RegistryName: "test-registry",
						Name:         "repo1",
						ImageCount:   5,
						CreatedAt:    "2024-01-01T00:00:00Z",
						UpdatedAt:    "2024-01-02T00:00:00Z",
					},
					{
						RegistryName: "test-registry",
						Name:         "repo2",
						ImageCount:   3,
						CreatedAt:    "2024-01-01T00:00:00Z",
						UpdatedAt:    "2024-01-02T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "empty response",
			registryID: "reg-123",
			response:   "",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			registryID: "reg-123",
			response:   `{"results": [{"name": "broken"`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "registry not found",
			registryID: "nonexistent",
			response:   `{"error": "registry not found"}`,
			statusCode: http.StatusNotFound,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			registryID: "reg-123",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
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
		{
			name:       "list with pagination",
			registryID: "reg-123",
			opts: RepositoryListOptions{
				Limit:  intPtr(10),
				Offset: intPtr(20),
			},
			expectedQuery: map[string]string{
				"_limit":  "10",
				"_offset": "20",
			},
			response: `{
				"meta": {
					"page": {
						"count": 1,
						"limit": 10,
						"offset": 20,
						"total": 30
					}
				},
				"results": [
					{
						"registry_name": "test-registry",
						"name": "repo1",
						"image_count": 5,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &RepositoriesResponse{
				Meta: Meta{Page: Page{Count: 1, Limit: 10, Offset: 20, Total: 30}},
				Results: []RepositoryResponse{
					{
						RegistryName: "test-registry",
						Name:         "repo1",
						ImageCount:   5,
						CreatedAt:    "2024-01-01T00:00:00Z",
						UpdatedAt:    "2024-01-02T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "list with sorting",
			registryID: "reg-123",
			opts: RepositoryListOptions{
				RepositoryFilterOptions: RepositoryFilterOptions{
					Sort: strPtr("name:asc"),
				},
			},
			expectedQuery: map[string]string{
				"_sort": "name:asc",
			},
			response: `{
				"meta": {
					"page": {
						"count": 1,
						"limit": 0,
						"offset": 0,
						"total": 1
					}
				},
				"results": [
					{
						"registry_name": "test-registry",
						"name": "repo1",
						"image_count": 5,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &RepositoriesResponse{
				Meta: Meta{Page: Page{Count: 1, Limit: 0, Offset: 0, Total: 1}},
				Results: []RepositoryResponse{
					{
						RegistryName: "test-registry",
						Name:         "repo1",
						ImageCount:   5,
						CreatedAt:    "2024-01-01T00:00:00Z",
						UpdatedAt:    "2024-01-02T00:00:00Z",
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
			got, err := client.Repositories().List(context.Background(), tt.registryID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if len(got.Results) != len(tt.want.Results) {
					t.Errorf("List() got %v results, want %v", len(got.Results), len(tt.want.Results))
				}
				if got.Meta.Page.Total != tt.want.Meta.Page.Total {
					t.Errorf("List() got total %v, want %v", got.Meta.Page.Total, tt.want.Meta.Page.Total)
				}
			}
		})
	}
}

func TestRepositoriesService_Get(t *testing.T) {
	tests := []struct {
		name           string
		registryID     string
		repositoryName string
		response       string
		statusCode     int
		want           *RepositoryResponse
		wantErr        bool
	}{
		{
			name:           "successful get repository",
			registryID:     "reg-123",
			repositoryName: "repo1",
			response: `{
				"registry_name": "test-registry",
				"name": "repo1",
				"image_count": 5,
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-02T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			want: &RepositoryResponse{
				RegistryName: "test-registry",
				Name:         "repo1",
				ImageCount:   5,
				CreatedAt:    "2024-01-01T00:00:00Z",
				UpdatedAt:    "2024-01-02T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name:           "repository not found",
			registryID:     "reg-123",
			repositoryName: "nonexistent",
			response:       `{"error": "repository not found"}`,
			statusCode:     http.StatusNotFound,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "malformed response",
			registryID:     "reg-123",
			repositoryName: "repo1",
			response:       `{"name": "broken"`,
			statusCode:     http.StatusOK,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "server error",
			registryID:     "reg-123",
			repositoryName: "repo1",
			response:       `{"error": "internal server error"}`,
			statusCode:     http.StatusInternalServerError,
			want:           nil,
			wantErr:        true,
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
			got, err := client.Repositories().Get(context.Background(), tt.registryID, tt.repositoryName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.Name != tt.want.Name {
					t.Errorf("Get() got name %v, want %v", got.Name, tt.want.Name)
				}
				if got.ImageCount != tt.want.ImageCount {
					t.Errorf("Get() got image count %v, want %v", got.ImageCount, tt.want.ImageCount)
				}
			}
		})
	}
}

func TestRepositoriesService_Delete(t *testing.T) {
	tests := []struct {
		name           string
		registryID     string
		repositoryName string
		statusCode     int
		response       string
		wantErr        bool
	}{
		{
			name:           "successful delete",
			registryID:     "reg-123",
			repositoryName: "repo1",
			statusCode:     http.StatusNoContent,
			wantErr:        false,
		},
		{
			name:           "repository not found",
			registryID:     "reg-123",
			repositoryName: "nonexistent",
			statusCode:     http.StatusNotFound,
			response:       `{"error": "repository not found"}`,
			wantErr:        true,
		},
		{
			name:           "unauthorized",
			registryID:     "reg-123",
			repositoryName: "repo1",
			statusCode:     http.StatusUnauthorized,
			response:       `{"error": "unauthorized"}`,
			wantErr:        true,
		},
		{
			name:           "server error",
			registryID:     "reg-123",
			repositoryName: "repo1",
			statusCode:     http.StatusInternalServerError,
			response:       `{"error": "internal server error"}`,
			wantErr:        true,
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
			err := client.Repositories().Delete(context.Background(), tt.registryID, tt.repositoryName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoriesService_Concurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"meta": {"page": {"count": 0, "limit": 0, "offset": 0, "total": 0}}, "results": []}`))
	}))
	defer server.Close()

	client := testClient(server.URL)
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.Repositories().List(ctx, "reg-123", RepositoryListOptions{})
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

func TestRepositoriesService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		registryID string
		filterOpts RepositoryFilterOptions
		responses  []string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "successful list all - single page",
			registryID: "reg-123",
			filterOpts: RepositoryFilterOptions{},
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
							"registry_name": "test-registry",
							"name": "repo-1",
							"image_count": 5,
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						},
						{
							"registry_name": "test-registry",
							"name": "repo-2",
							"image_count": 10,
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
			registryID: "reg-123",
			filterOpts: RepositoryFilterOptions{},
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
					"results": [` + generateRepositoryJSONArray(50) + `]
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
					"results": [` + generateRepositoryJSONArray(25) + `]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  75,
			wantErr:    false,
		},
		{
			name:       "empty results",
			registryID: "reg-123",
			filterOpts: RepositoryFilterOptions{},
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
			name:       "error on first page",
			registryID: "reg-123",
			filterOpts: RepositoryFilterOptions{},
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
			got, err := client.Repositories().ListAll(context.Background(), tt.registryID, tt.filterOpts)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("ListAll() got %v repositories, want %v", len(got), tt.wantCount)
			}
		})
	}
}

// Helper function to generate repository JSON array for testing pagination
func generateRepositoryJSONArray(count int) string {
	var repositories []string
	for i := 0; i < count; i++ {
		repositories = append(repositories, fmt.Sprintf(`{
			"registry_name": "test-registry",
			"name": "repo-%d",
			"image_count": %d,
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z"
		}`, i, i+1))
	}
	result := ""
	for i, repo := range repositories {
		if i > 0 {
			result += ","
		}
		result += repo
	}
	return result
}
