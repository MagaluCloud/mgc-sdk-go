package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		name       string
		clusterID  string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:      "successful get cluster",
			clusterID: "cluster-123",
			response: `{
				"id": "cluster-123",
				"name": "production"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "invalid cluster ID",
			clusterID: "",
			wantErr:   true,
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
			_, err := client.Clusters().Get(context.Background(), tt.clusterID, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClusterService_Update(t *testing.T) {
	tests := []struct {
		name        string
		clusterID   string
		request     AllowedCIDRsUpdateRequest
		response    string
		statusCode  int
		wantUpdated bool
		wantErr     bool
	}{
		{
			name:      "successful update",
			clusterID: "cluster-123",
			request: AllowedCIDRsUpdateRequest{
				AllowedCIDRs: []string{"192.168.1.0/24"},
			},
			response:    `{"allowed_cidrs": ["192.168.1.0/24"]}`,
			statusCode:  http.StatusOK,
			wantUpdated: true,
			wantErr:     false,
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

			if !tt.wantErr && len(result.AllowedCIDRs) != len(tt.request.AllowedCIDRs) {
				t.Errorf("Update() cidrs atualizados = %d, esperado %d",
					len(result.AllowedCIDRs), len(tt.request.AllowedCIDRs))
			}
		})
	}
}

func TestClusterService_GetKubeConfig(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:      "valid kubeconfig",
			clusterID: "cluster-123",
			response: `{
				"apiVersion": "v1",
				"clusters": [{"name": "test"}]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
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
			_, err := client.Clusters().GetKubeConfig(context.Background(), tt.clusterID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetKubeConfig() error = %v, wantErr %v", err, tt.wantErr)
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
