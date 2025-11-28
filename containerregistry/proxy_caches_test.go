package containerregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestProxyCachesService_List(t *testing.T) {
	tests := []struct {
		name          string
		opts          ProxyCacheListOptions
		response      string
		statusCode    int
		expectedQuery map[string]string
		want          *ListProxyCachesResponse
		wantErr       bool
	}{
		{
			name: "successful list proxy-caches",
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
						"id": "id-1",
						"name": "proxy-1",
						"provider": "docker-hub",
						"url": "https://hub.docker.com/repositories",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					},
					{
						"id": "id-2",
						"name": "proxy-2",
						"provider": "docker-hub",
						"url": "https://hub.docker.com/repositories",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListProxyCachesResponse{
				Meta: Meta{Page: Page{Count: 2, Limit: 0, Offset: 0, Total: 2}},
				Results: []ProxyCache{
					{
						ID:        "id-1",
						Name:      "proxy-1",
						Provider:  "docker-hub",
						URL:       "https://hub.docker.com/repositories",
						CreatedAt: "2024-01-02T00:00:00Z",
						UpdatedAt: "2024-01-02T00:00:00Z",
					},
					{
						ID:        "id-2",
						Name:      "proxy-2",
						Provider:  "docker-hub",
						URL:       "https://hub.docker.com/repositories",
						CreatedAt: "2024-01-02T00:00:00Z",
						UpdatedAt: "2024-01-02T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "empty response",
			response:   "",
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			response:   `{"results": [{"name": "broken"`,
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
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			want:       nil,
			wantErr:    true,
		},
		{
			name: "list with pagination",
			opts: ProxyCacheListOptions{
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
						"id": "id-1",
						"name": "proxy-1",
						"provider": "docker-hub",
						"url": "https://hub.docker.com/repositories",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListProxyCachesResponse{
				Meta: Meta{Page: Page{Count: 1, Limit: 10, Offset: 20, Total: 30}},
				Results: []ProxyCache{
					{
						ID:        "id-1",
						Name:      "proxy-1",
						Provider:  "docker-hub",
						URL:       "https://hub.docker.com/repositories",
						CreatedAt: "2024-01-02T00:00:00Z",
						UpdatedAt: "2024-01-02T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "list with sorting",
			opts: ProxyCacheListOptions{
				Sort: strPtr("name:asc"),
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
						"id": "id-1",
						"name": "proxy-1",
						"provider": "docker-hub",
						"url": "https://hub.docker.com/repositories",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListProxyCachesResponse{
				Meta: Meta{Page: Page{Count: 1, Limit: 0, Offset: 0, Total: 1}},
				Results: []ProxyCache{
					{
						ID:        "id-1",
						Name:      "proxy-1",
						Provider:  "docker-hub",
						URL:       "https://hub.docker.com/repositories",
						CreatedAt: "2024-01-02T00:00:00Z",
						UpdatedAt: "2024-01-02T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/container-registry/v0/proxy-caches" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

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
			got, err := client.ProxyCaches().List(context.Background(), tt.opts)

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

func TestProxyCachesService_ListAll(t *testing.T) {
	tests := []struct {
		name          string
		opts          ProxyCacheListAllOptions
		responses     []string
		statusCode    int
		expectedQuery map[string]string
		wantCount     int
		wantErr       bool
	}{
		{
			name: "successful list all - single page",
			opts: ProxyCacheListAllOptions{},
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
							"id": "id-1",
							"name": "proxy-1",
							"provider": "docker-hub",
							"url": "https://hub.docker.com/repositories",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-02T00:00:00Z"
						},
						{
							"id": "id-2",
							"name": "proxy-2",
							"provider": "docker-hub",
							"url": "https://hub.docker.com/repositories",
							"created_at": "2024-01-01T00:00:00Z",
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
			name: "successful list all - multiple pages",
			opts: ProxyCacheListAllOptions{},
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
					"results": [` + generateProxyCachesJSONArray(50) + `]
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
					"results": [` + generateProxyCachesJSONArray(25) + `]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  75,
			wantErr:    false,
		},
		{
			name: "empty results",
			opts: ProxyCacheListAllOptions{},
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
			opts:       ProxyCacheListAllOptions{},
			responses:  []string{`{"error": "internal server error"}`},
			statusCode: http.StatusInternalServerError,
			wantCount:  0,
			wantErr:    true,
		},
		{
			name: "successful list all - with sorting",
			opts: ProxyCacheListAllOptions{
				Sort: strPtr("name:asc"),
			},
			expectedQuery: map[string]string{
				"_sort": "name:asc",
			},
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
							"id": "id-1",
							"name": "proxy-1",
							"provider": "docker-hub",
							"url": "https://hub.docker.com/repositories",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-02T00:00:00Z"
						},
						{
							"id": "id-2",
							"name": "proxy-2",
							"provider": "docker-hub",
							"url": "https://hub.docker.com/repositories",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-02T00:00:00Z"
						}
					]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/container-registry/v0/proxy-caches" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

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

				if requestCount < len(tt.responses) {
					w.Write([]byte(tt.responses[requestCount]))
					requestCount++
				}
			}))

			defer server.Close()

			client := testClient(server.URL)
			got, err := client.ProxyCaches().ListAll(context.Background(), tt.opts)

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

func TestProxyCachesService_Get(t *testing.T) {
	tests := []struct {
		name         string
		proxyCacheID string
		response     string
		statusCode   int
		want         *GetProxyCacheResponse
		wantErr      bool
	}{
		{
			name:         "successful get proxy-cache",
			proxyCacheID: "id-1",
			response: `{
				"id": "id-1",
				"name": "proxy-1",
				"description": "Description proxy-cache",
				"provider": "docker-hub",
				"url": "https://hub.docker.com/repositories",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-02T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			want: &GetProxyCacheResponse{
				ProxyCache: ProxyCache{ID: "id-1",
					Name:      "proxy-1",
					Provider:  "docker-hub",
					URL:       "https://hub.docker.com/repositories",
					CreatedAt: "2024-01-02T00:00:00Z",
					UpdatedAt: "2024-01-02T00:00:00Z",
				},
				Description: "Description proxy-cache",
			},
			wantErr: false,
		},
		{
			name:         "proxy-cache not found",
			proxyCacheID: "id-invalid",
			response:     `{"error": "proxy-cache not found"}`,
			statusCode:   http.StatusNotFound,
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "malformed response",
			proxyCacheID: "id-1",
			response:     `{"name": "broken"`,
			statusCode:   http.StatusOK,
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "server error",
			proxyCacheID: "id-1",
			response:     `{"error": "internal server error"}`,
			statusCode:   http.StatusInternalServerError,
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/container-registry/v0/proxy-caches/%s", tt.proxyCacheID) {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testClient(server.URL)
			got, err := client.ProxyCaches().Get(context.Background(), tt.proxyCacheID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil {
				if got.Name != tt.want.Name {
					t.Errorf("Get() got name %v, want %v", got.Name, tt.want.Name)
				}

				if got.Description != tt.want.Description {
					t.Errorf("Get() got description %v, want %v", got.Description, tt.want.Description)
				}
			}
		})
	}
}

func TestProxyCachesService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateProxyCacheRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateProxyCacheRequest{
				Name:     "proxy-cache",
				Provider: "docker-hub",
				URL:      "https://hub.docker.com/repositories",
			},
			response:   `{"id": "1"}`,
			statusCode: http.StatusOK,
			wantID:     "1",
		},
		{
			name: "missing name",
			request: CreateProxyCacheRequest{
				Provider: "docker-hub",
				URL:      "https://hub.docker.com/repositories",
			},
			response:   `{"error": "name required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing provider",
			request: CreateProxyCacheRequest{
				Name: "proxy-cache",
				URL:  "https://hub.docker.com/repositories",
			},
			response:   `{"error": "provider required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing url",
			request: CreateProxyCacheRequest{
				Name:     "proxy-cache",
				Provider: "docker-hub",
			},
			response:   `{"error": "url required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/container-registry/v0/proxy-caches" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}

				var req CreateProxyCacheRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("error decoding request: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testClient(server.URL)
			got, err := client.ProxyCaches().Create(context.Background(), tt.request)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.ID != tt.wantID {
				t.Errorf("got ID %q, want %q", got.ID, tt.wantID)
			}
		})
	}
}

func TestProxyCachesService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		proxyCacheID string
		statusCode   int
		response     string
		wantErr      bool
	}{
		{
			name:         "successful delete",
			proxyCacheID: "id-123",
			statusCode:   http.StatusNoContent,
			wantErr:      false,
		},
		{
			name:         "proxy-cache not found",
			proxyCacheID: "id-123",
			statusCode:   http.StatusNotFound,
			response:     `{"error": "proxy-cache not found"}`,
			wantErr:      true,
		},
		{
			name:         "unauthorized",
			proxyCacheID: "id-123",
			statusCode:   http.StatusUnauthorized,
			response:     `{"error": "unauthorized"}`,
			wantErr:      true,
		},
		{
			name:         "server error",
			proxyCacheID: "id-123",
			statusCode:   http.StatusInternalServerError,
			response:     `{"error": "internal server error"}`,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/container-registry/v0/proxy-caches/%s", tt.proxyCacheID) {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

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
			err := client.ProxyCaches().Delete(context.Background(), tt.proxyCacheID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyCachesService_Update(t *testing.T) {
	tests := []struct {
		name         string
		proxyCacheID string
		request      UpdateProxyCacheRequest
		response     string
		statusCode   int
		wantID       string
		wantErr      bool
	}{
		{
			name:         "successful update",
			proxyCacheID: "id-1",
			request: UpdateProxyCacheRequest{
				Name: helpers.StrPtr("new-proxy-cache"),
			},
			response: `{
				"id": "id-1",
				"name": "new-proxy-cache",
				"provider": "docker-hub",
				"url": "https://hub.docker.com/repositories",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-02T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			wantID:     "id-1",
		},
		{
			name:         "proxy-cache not found",
			proxyCacheID: "invalid-id",
			request: UpdateProxyCacheRequest{
				URL: helpers.StrPtr("https://hub.docker.com/repositories"),
			},
			response:   `{"error": "invalid id"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:         "unauthorized",
			proxyCacheID: "id-1",
			statusCode:   http.StatusUnauthorized,
			response:     `{"error": "unauthorized"}`,
			wantErr:      true,
		},
		{
			name:         "server error",
			proxyCacheID: "id-1",
			statusCode:   http.StatusInternalServerError,
			response:     `{"error": "internal server error"}`,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/container-registry/v0/proxy-caches/%s", tt.proxyCacheID) {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				if r.Method != http.MethodPatch {
					t.Errorf("expected PATCH method, got %s", r.Method)
				}

				var req CreateProxyCacheRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("error decoding request: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testClient(server.URL)
			got, err := client.ProxyCaches().Update(context.Background(), tt.proxyCacheID, tt.request)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.ID != tt.wantID {
				t.Errorf("got ID %q, want %q", got.ID, tt.wantID)
			}
		})
	}
}

func TestProxyCachesService_ListStatus(t *testing.T) {
	tests := []struct {
		name         string
		proxyCacheID string
		response     string
		statusCode   int
		want         *ListProxyCacheStatusResponse
		wantErr      bool
	}{
		{
			name:         "successful list proxy-cache status",
			proxyCacheID: "id-1",
			response: `{
				"message": "Registry is healthy and reachable",
				"status": "healthy"
			}`,
			statusCode: http.StatusOK,
			want: &ListProxyCacheStatusResponse{
				Message: "Registry is healthy and reachable",
				Status:  "healthy",
			},
			wantErr: false,
		},
		{
			name:         "proxy-cache not found",
			proxyCacheID: "id-invalid",
			response:     `{"error": "proxy-cache not found"}`,
			statusCode:   http.StatusNotFound,
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "malformed response",
			proxyCacheID: "id-1",
			response:     `{"message": "broken"`,
			statusCode:   http.StatusOK,
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "server error",
			proxyCacheID: "id-1",
			response:     `{"error": "internal server error"}`,
			statusCode:   http.StatusInternalServerError,
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/container-registry/v0/proxy-caches/%s/status", tt.proxyCacheID) {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testClient(server.URL)
			got, err := client.ProxyCaches().ListStatus(context.Background(), tt.proxyCacheID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil {
				if got.Message != tt.want.Message {
					t.Errorf("Get() got message %v, want %v", got.Message, tt.want.Message)
				}

				if got.Status != tt.want.Status {
					t.Errorf("Get() got status %v, want %v", got.Status, tt.want.Status)
				}
			}
		})
	}
}

func TestProxyCachesService_CreateStatus(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateProxyCacheStatusRequest
		response   string
		statusCode int
		want       *CreateProxyCacheStatusResponse
		wantErr    bool
	}{
		{
			name: "successful creation of proxy-cache status",
			request: CreateProxyCacheStatusRequest{
				Provider:     "docker-hub",
				URL:          "https://hub.docker.com/repositories",
				AccessKey:    "test@gmail.com",
				AccessSecret: "teste.123",
			},
			response: `{
				"message": "Proxy cache credentials are valid.",
				"status": "valid"
			}`,
			statusCode: http.StatusOK,
			want: &CreateProxyCacheStatusResponse{
				Message: "Proxy cache credentials are valid.",
				Status:  "valid",
			},
			wantErr: false,
		},
		{
			name: "missing provider",
			request: CreateProxyCacheStatusRequest{
				URL:          "https://hub.docker.com/repositories",
				AccessKey:    "test@gmail.com",
				AccessSecret: "teste.123",
			},
			response:   `{"error": "provider required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing url",
			request: CreateProxyCacheStatusRequest{
				Provider:     "docker-hub",
				AccessKey:    "test@gmail.com",
				AccessSecret: "teste.123",
			},
			response:   `{"error": "url required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing access key",
			request: CreateProxyCacheStatusRequest{
				Provider:     "docker-hub",
				URL:          "https://hub.docker.com/repositories",
				AccessSecret: "teste.123",
			},
			response:   `{"error": "access key required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing access secret",
			request: CreateProxyCacheStatusRequest{
				Provider:  "docker-hub",
				URL:       "https://hub.docker.com/repositories",
				AccessKey: "test@gmail.com",
			},
			response:   `{"error": "access secret required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "malformed response",
			request: CreateProxyCacheStatusRequest{
				Provider:     "docker-hub",
				URL:          "https://hub.docker.com/repositories",
				AccessKey:    "test@gmail.com",
				AccessSecret: "teste.123",
			},
			response:   `{"message": "broken"`,
			statusCode: http.StatusOK,
			want:       nil,
			wantErr:    true,
		},
		{
			name: "server error",
			request: CreateProxyCacheStatusRequest{
				Provider:     "docker-hub",
				URL:          "https://hub.docker.com/repositories",
				AccessKey:    "test@gmail.com",
				AccessSecret: "teste.123",
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/container-registry/v0/proxy-caches/status" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testClient(server.URL)
			got, err := client.ProxyCaches().CreateStatus(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil {
				if got.Message != tt.want.Message {
					t.Errorf("Get() got message %v, want %v", got.Message, tt.want.Message)
				}

				if got.Status != tt.want.Status {
					t.Errorf("Get() got status %v, want %v", got.Status, tt.want.Status)
				}
			}
		})
	}
}

// Helper function to generate proxy-caches JSON array for testing pagination
func generateProxyCachesJSONArray(count int) string {
	var repositories []string
	for i := 0; i < count; i++ {
		repositories = append(repositories, fmt.Sprintf(`{
			"id": "id-%d",
			"name": "test-proxy-cache",
			"provider": "docker-hub",
			"url": "https://hub.docker.com/repositories",
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z"
		}`, i+1))
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
