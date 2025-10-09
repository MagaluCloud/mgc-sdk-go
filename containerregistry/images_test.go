package containerregistry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImagesService_List(t *testing.T) {
	tests := []struct {
		name           string
		registryID     string
		repositoryName string
		opts           ListOptions
		response       string
		statusCode     int
		expectedQuery  map[string]string
		want           *ImagesResponse
		wantErr        bool
	}{
		{
			name:           "successful list images",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			response: `{
				"results": [
					{
						"digest": "sha256:123",
						"size_bytes": 1024,
						"pushed_at": "2024-01-01T00:00:00Z",
						"pulled_at": "2024-01-02T00:00:00Z",
						"manifest_media_type": "application/vnd.docker.distribution.manifest.v2+json",
						"media_type": "application/vnd.docker.container.image.v1+json",
						"tags": ["latest", "v1.0"],
						"tags_details": [
							{
								"name": "latest",
								"pushed_at": "2024-01-01T00:00:00Z",
								"pulled_at": "2024-01-02T00:00:00Z",
								"signed": true
							}
						]
					}
				]
			}`,
			statusCode:    http.StatusOK,
			expectedQuery: map[string]string{},
			want: &ImagesResponse{
				Results: []ImageResponse{
					{
						Digest:            "sha256:123",
						SizeBytes:         1024,
						PushedAt:          "2024-01-01T00:00:00Z",
						PulledAt:          "2024-01-02T00:00:00Z",
						ManifestMediaType: "application/vnd.docker.distribution.manifest.v2+json",
						MediaType:         "application/vnd.docker.container.image.v1+json",
						Tags:              []string{"latest", "v1.0"},
						TagsDetails: []ImageTagResponse{
							{
								Name:     "latest",
								PushedAt: "2024-01-01T00:00:00Z",
								PulledAt: "2024-01-02T00:00:00Z",
								Signed:   true,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:           "list images with limit",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			opts: ListOptions{
				Limit: intPtr(10),
			},
			response: `{
				"results": []
			}`,
			statusCode:    http.StatusOK,
			expectedQuery: map[string]string{"_limit": "10"},
			want: &ImagesResponse{
				Results: []ImageResponse{},
			},
			wantErr: false,
		},
		{
			name:           "list images with offset",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			opts: ListOptions{
				Offset: intPtr(5),
			},
			response: `{
				"results": []
			}`,
			statusCode:    http.StatusOK,
			expectedQuery: map[string]string{"_offset": "5"},
			want: &ImagesResponse{
				Results: []ImageResponse{},
			},
			wantErr: false,
		},
		{
			name:           "list images with sort",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			opts: ListOptions{
				Sort: strPtr("pushed_at"),
			},
			response: `{
				"results": []
			}`,
			statusCode:    http.StatusOK,
			expectedQuery: map[string]string{"_sort": "pushed_at"},
			want: &ImagesResponse{
				Results: []ImageResponse{},
			},
			wantErr: false,
		},
		{
			name:           "list images with expand",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			opts: ListOptions{
				Expand: []string{"tags", "manifest"},
			},
			response: `{
				"results": []
			}`,
			statusCode:    http.StatusOK,
			expectedQuery: map[string]string{"_expand": "tags,manifest"},
			want: &ImagesResponse{
				Results: []ImageResponse{},
			},
			wantErr: false,
		},
		{
			name:           "list images with multiple options",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			opts: ListOptions{
				Limit:  intPtr(20),
				Offset: intPtr(10),
				Sort:   strPtr("created_at"),
				Expand: []string{"tags"},
			},
			response: `{
				"results": []
			}`,
			statusCode: http.StatusOK,
			expectedQuery: map[string]string{
				"_limit":  "20",
				"_offset": "10",
				"_sort":   "created_at",
				"_expand": "tags",
			},
			want: &ImagesResponse{
				Results: []ImageResponse{},
			},
			wantErr: false,
		},
		{
			name:           "empty response",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			response:       "",
			statusCode:     http.StatusOK,
			expectedQuery:  map[string]string{},
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "malformed json",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			response:       `{"results": [{"digest": "sha256:123"`,
			statusCode:     http.StatusOK,
			expectedQuery:  map[string]string{},
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "server error",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			response:       `{"error": "internal server error"}`,
			statusCode:     http.StatusInternalServerError,
			expectedQuery:  map[string]string{},
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "unauthorized",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			response:       `{"error": "unauthorized"}`,
			statusCode:     http.StatusUnauthorized,
			expectedQuery:  map[string]string{},
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "not found",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			response:       `{"error": "repository not found"}`,
			statusCode:     http.StatusNotFound,
			expectedQuery:  map[string]string{},
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
			got, err := client.Images().List(context.Background(), tt.registryID, tt.repositoryName, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if len(got.Results) != len(tt.want.Results) {
					t.Errorf("List() got %v results, want %v", len(got.Results), len(tt.want.Results))
				}
			}
		})
	}
}

func TestImagesService_Delete(t *testing.T) {
	tests := []struct {
		name           string
		registryID     string
		repositoryName string
		digestOrTag    string
		statusCode     int
		response       string
		wantErr        bool
	}{
		{
			name:           "successful delete",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "latest",
			statusCode:     http.StatusNoContent,
			wantErr:        false,
		},
		{
			name:           "not found",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "nonexistent",
			statusCode:     http.StatusNotFound,
			response:       `{"error": "image not found"}`,
			wantErr:        true,
		},
		{
			name:           "unauthorized",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "latest",
			statusCode:     http.StatusUnauthorized,
			response:       `{"error": "unauthorized"}`,
			wantErr:        true,
		},
		{
			name:           "server error",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "latest",
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
			err := client.Images().Delete(context.Background(), tt.registryID, tt.repositoryName, tt.digestOrTag)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImagesService_Get(t *testing.T) {
	tests := []struct {
		name           string
		registryID     string
		repositoryName string
		digestOrTag    string
		response       string
		statusCode     int
		want           *ImageResponse
		wantErr        bool
	}{
		{
			name:           "successful get image",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "latest",
			response: `{
				"digest": "sha256:123",
				"size_bytes": 1024,
				"pushed_at": "2024-01-01T00:00:00Z",
				"pulled_at": "2024-01-02T00:00:00Z",
				"manifest_media_type": "application/vnd.docker.distribution.manifest.v2+json",
				"media_type": "application/vnd.docker.container.image.v1+json",
				"tags": ["latest"],
				"tags_details": [
					{
						"name": "latest",
						"pushed_at": "2024-01-01T00:00:00Z",
						"pulled_at": "2024-01-02T00:00:00Z",
						"signed": true
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ImageResponse{
				Digest:            "sha256:123",
				SizeBytes:         1024,
				PushedAt:          "2024-01-01T00:00:00Z",
				PulledAt:          "2024-01-02T00:00:00Z",
				ManifestMediaType: "application/vnd.docker.distribution.manifest.v2+json",
				MediaType:         "application/vnd.docker.container.image.v1+json",
				Tags:              []string{"latest"},
				TagsDetails: []ImageTagResponse{
					{
						Name:     "latest",
						PushedAt: "2024-01-01T00:00:00Z",
						PulledAt: "2024-01-02T00:00:00Z",
						Signed:   true,
					},
				},
			},
			wantErr: false,
		},
		{
			name:           "not found",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "nonexistent",
			response:       `{"error": "image not found"}`,
			statusCode:     http.StatusNotFound,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "malformed response",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "latest",
			response:       `{"digest": "sha256:123"`,
			statusCode:     http.StatusOK,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "empty response",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "latest",
			response:       "",
			statusCode:     http.StatusOK,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "server error",
			registryID:     "reg-123",
			repositoryName: "repo-test",
			digestOrTag:    "latest",
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
			got, err := client.Images().Get(context.Background(), tt.registryID, tt.repositoryName, tt.digestOrTag)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.Digest != tt.want.Digest {
					t.Errorf("Get() got digest %v, want %v", got.Digest, tt.want.Digest)
				}
				if len(got.Tags) != len(tt.want.Tags) {
					t.Errorf("Get() got %v tags, want %v", len(got.Tags), len(tt.want.Tags))
				}
			}
		})
	}
}

func TestImagesService_Concurrent(t *testing.T) {
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
			_, err := client.Images().List(ctx, "reg-123", "repo-test", ListOptions{})
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
