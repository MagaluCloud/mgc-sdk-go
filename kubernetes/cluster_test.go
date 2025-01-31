package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClusterService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list clusters",
			response: `{
				"results": [
					{"id": "cluster1", "name": "prod-cluster"},
					{"id": "cluster2", "name": "staging-cluster"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "successful list clusters",
			response: `{
				"results": [
					{"id": "cluster1", "name": "prod-cluster"},
					{"id": "cluster2", "name": "staging-cluster"}
				]
			}`,
			opts: ListOptions{
				Limit:  intPtr(10),
				Offset: intPtr(0),
				Expand: []string{"network"},
				Sort:   strPtr("name"),
			},
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty response",
			response:   `[]`,
			statusCode: http.StatusOK,
			want:       0,
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
			result, err := client.Clusters().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("List() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestClusterService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    ClusterRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful create cluster",
			request: ClusterRequest{
				Name:    "new-cluster",
				Version: "v1.30.2",
			},
			response:   `{"id": "cluster-new"}`,
			statusCode: http.StatusCreated,
			wantID:     "cluster-new",
			wantErr:    false,
		},
		{
			name: "invalid request",
			request: ClusterRequest{
				Name:    "",
				Version: "v1.30.2",
			},
			wantErr: true,
		},
		{
			name: "server error",
			request: ClusterRequest{
				Name:    "new-cluster",
				Version: "v1.30.2",
			},
			wantErr: true,
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
			result, err := client.Clusters().Create(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.ID != tt.wantID {
				t.Errorf("Create() got ID = %s, want %s", result.ID, tt.wantID)
			}
		})
	}
}

func TestClusterService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			clusterID:  "cluster-123",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "invalid cluster ID",
			clusterID:  "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "already deleted",
			clusterID:  "cluster-404",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			clusterID:  "cluster-123",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Clusters().Delete(context.Background(), tt.clusterID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClusterService_Get(t *testing.T) {
	tests := []struct {
		name        string
		clusterID   string
		response    string
		statusCode  int
		wantID      string
		wantVersion string
		wantErr     bool
	}{
		{
			name:        "successful get",
			clusterID:   "cluster-123",
			response:    `{"id": "cluster-123", "version": "v1.30.2"}`,
			statusCode:  http.StatusOK,
			wantID:      "cluster-123",
			wantVersion: "v1.30.2",
			wantErr:     false,
		},
		{
			name:       "invalid cluster ID",
			clusterID:  "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			clusterID:  "cluster-404",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Clusters().Get(context.Background(), tt.clusterID, []string{"network"})

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.ID != tt.wantID {
					t.Errorf("Get() ID = %s, want %s", result.ID, tt.wantID)
				}
				if result.Version != tt.wantVersion {
					t.Errorf("Get() Version = %s, want %s", result.Version, tt.wantVersion)
				}
			}
		})
	}
}

func TestClusterService_Update(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		request    AllowedCIDRsUpdateRequest
		response   string
		statusCode int
		wantCIDRs  int
		wantErr    bool
	}{
		{
			name:      "successful update",
			clusterID: "cluster-123",
			request: AllowedCIDRsUpdateRequest{
				AllowedCIDRs: []string{"192.168.1.0/24", "10.0.0.0/8"},
			},
			response:   `{"allowed_cidrs": ["192.168.1.0/24", "10.0.0.0/8"]}`,
			statusCode: http.StatusOK,
			wantCIDRs:  2,
			wantErr:    false,
		},
		{
			name:       "empty CIDRs list",
			clusterID:  "cluster-123",
			request:    AllowedCIDRsUpdateRequest{},
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "server error",
			clusterID:  "cluster-123",
			request:    AllowedCIDRsUpdateRequest{},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "empty cluster ID",
			clusterID:  "",
			request:    AllowedCIDRsUpdateRequest{},
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Clusters().Update(context.Background(), tt.clusterID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result.AllowedCIDRs) != tt.wantCIDRs {
				t.Errorf("Update() CIDRs = %d, want %d", len(result.AllowedCIDRs), tt.wantCIDRs)
			}
		})
	}
}

func TestClusterService_GetKubeConfig(t *testing.T) {
	tests := []struct {
		name        string
		clusterID   string
		response    string
		statusCode  int
		wantContent string
		wantErr     bool
	}{
		{
			name:        "valid kubeconfig",
			clusterID:   "cluster-123",
			response:    "apiVersion: v1\nclusters:\n- cluster: {}\n",
			statusCode:  http.StatusOK,
			wantContent: "v1",
			wantErr:     false,
		},
		{
			name:       "empty response",
			clusterID:  "cluster-123",
			response:   "",
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "empty cluster ID",
			clusterID:  "",
			response:   "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/yaml")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Clusters().GetKubeConfig(context.Background(), tt.clusterID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetKubeConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !strings.Contains(result.APIVersion, tt.wantContent) {
				t.Errorf("GetKubeConfig() content = %s, want containing %s", result.APIVersion, tt.wantContent)
			}
		})
	}
}

func TestClusterService_NodePoolOperations(t *testing.T) {
	client := testClient("http://dummy")

	t.Run("NodePools service exists", func(t *testing.T) {
		service := client.Nodepools()
		if service == nil {
			t.Error("NodePools service não disponível")
		}
	})
}

func TestClusterService_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		clusterID string
		wantErr   bool
	}{
		{
			name:      "empty cluster ID",
			clusterID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient("http://dummy")
			_, err := client.Clusters().Get(context.Background(), tt.clusterID, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("Validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClusterService_EdgeCases(t *testing.T) {
	t.Run("context timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := testClient(server.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := client.Clusters().Get(ctx, "cluster-123", []string{})
		if err == nil {
			t.Error("Esperado erro de timeout")
		}
	})

	t.Run("invalid response format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"invalid": "data`))
		}))
		defer server.Close()

		client := testClient(server.URL)
		_, err := client.Clusters().Get(context.Background(), "cluster-123", []string{})
		if err == nil {
			t.Error("Esperado erro de parsing")
		}
	})
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
