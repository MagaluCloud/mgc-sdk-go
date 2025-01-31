package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNodePoolService_List(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		opts       ListOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name:      "successful list node pools",
			clusterID: "cluster-123",
			response: `{
				"results": [
					{"id": "pool1", "name": "default-pool"},
					{"id": "pool2", "name": "system-pool"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:      "invalid cluster ID",
			clusterID: "",
			wantErr:   true,
		},
		{
			name:       "server error",
			clusterID:  "cluster-123",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Nodepools().List(context.Background(), tt.clusterID, tt.opts)

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

func TestNodePoolService_Create(t *testing.T) {
	tests := []struct {
		name         string
		clusterID    string
		request      CreateNodePoolRequest
		response     string
		statusCode   int
		wantID       string
		wantReplicas int
		wantErr      bool
	}{
		{
			name:      "successful create node pool",
			clusterID: "cluster-123",
			request: CreateNodePoolRequest{
				Name:     "new-pool",
				Flavor:   "gp1.small",
				Replicas: 3,
			},
			response: `{
				"id": "pool-new",
				"name": "new-pool",
				"replicas": 3
			}`,
			statusCode:   http.StatusCreated,
			wantID:       "pool-new",
			wantReplicas: 3,
			wantErr:      false,
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
			result, err := client.Nodepools().Create(context.Background(), tt.clusterID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.ID != tt.wantID || result.Replicas != tt.wantReplicas {
					t.Errorf("Create() got = %v, want ID %s with %d replicas", result, tt.wantID, tt.wantReplicas)
				}
			}
		})
	}
}

func TestNodePoolService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		nodePoolID string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			clusterID:  "cluster-123",
			nodePoolID: "pool-456",
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
			err := client.Nodepools().Delete(context.Background(), tt.clusterID, tt.nodePoolID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
