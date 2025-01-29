package containerregistry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRepositoriesService_List(t *testing.T) {
	tests := []struct {
		name       string
		registryID string
		opts       ListOptions
		response   string
		statusCode int
		want       *RepositoriesResponse
		wantErr    bool
	}{
		{
			name:       "successful list repositories",
			registryID: "reg-123",
			response: `{
				"goal": {
					"total": 2
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
				Goal: AmountRepositoryResponse{Total: 2},
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
			opts: ListOptions{
				Limit:  intPtr(10),
				Offset: intPtr(20),
			},
			response: `{
				"goal": {"total": 30},
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
				Goal: AmountRepositoryResponse{Total: 30},
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
			opts: ListOptions{
				Sort: strPtr("name:asc"),
			},
			response: `{
				"goal": {"total": 1},
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
				Goal: AmountRepositoryResponse{Total: 1},
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
				if got.Goal.Total != tt.want.Goal.Total {
					t.Errorf("List() got total %v, want %v", got.Goal.Total, tt.want.Goal.Total)
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
		w.Write([]byte(`{"goal": {"total": 0}, "results": []}`))
	}))
	defer server.Close()

	client := testClient(server.URL)
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.Repositories().List(ctx, "reg-123", ListOptions{})
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
