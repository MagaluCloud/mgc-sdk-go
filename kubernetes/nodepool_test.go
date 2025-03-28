package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
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
			opts: ListOptions{
				Limit:  intPtr(2),
				Offset: intPtr(1),
				Sort:   strPtr("name"),
			},
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
		{
			name:       "invalid cluster ID",
			clusterID:  "",
			nodePoolID: "pool-456",
			wantErr:    true,
		},
		{
			name:       "invalid cluster ID",
			clusterID:  "pool-456",
			nodePoolID: "",
			wantErr:    true,
		},
		{
			name:       "invalid node pool ID",
			clusterID:  "cluster-123",
			nodePoolID: "asdasd",
			wantErr:    true,
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

func TestNodePoolService_List_InvalidOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := testClient(server.URL)
	_, err := client.Nodepools().List(context.Background(), "cluster-123", ListOptions{
		Limit: helpers.IntPtr(-1),
	})

	if err == nil {
		t.Error("Esperado erro com opções inválidas")
	}
}

func TestNodePoolService_Create_ValidationError(t *testing.T) {
	tests := []struct {
		name      string
		clusterID string
		request   CreateNodePoolRequest
		wantErr   bool
	}{
		{
			name:      "empty cluster ID",
			clusterID: "",
			request: CreateNodePoolRequest{
				Name:     "test-pool",
				Flavor:   "gp1.small",
				Replicas: 2,
			},
			wantErr: true,
		},
		{
			name:      "invalid replicas",
			clusterID: "cluster-123",
			request: CreateNodePoolRequest{
				Name:     "test-pool",
				Flavor:   "gp1.small",
				Replicas: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := testClient("http://dummy")
			_, err := client.Nodepools().Create(context.Background(), tt.clusterID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodePoolService_Scale(t *testing.T) {
	tests := []struct {
		name         string
		clusterID    string
		nodePoolID   string
		replicas     int
		response     string
		statusCode   int
		wantReplicas int
		wantErr      bool
	}{
		{
			name:         "successful scale up",
			clusterID:    "cluster-123",
			nodePoolID:   "pool-456",
			replicas:     5,
			response:     `{"replicas": 5}`,
			statusCode:   http.StatusOK,
			wantReplicas: 5,
			wantErr:      false,
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
			result, err := client.Nodepools().Update(context.Background(), tt.clusterID, tt.nodePoolID, PatchNodePoolRequest{
				Replicas: helpers.IntPtr(tt.replicas),
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("Scale() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.Replicas != tt.wantReplicas {
				t.Errorf("Scale() replicas = %d, want %d", result.Replicas, tt.wantReplicas)
			}
		})
	}
}

func TestNodePoolService_Update(t *testing.T) {
	tests := []struct {
		name         string
		clusterID    string
		nodePoolID   string
		request      PatchNodePoolRequest
		response     string
		statusCode   int
		wantReplicas int
		wantErr      bool
	}{
		{
			name:       "successful update",
			clusterID:  "cluster-123",
			nodePoolID: "pool-456",
			request: PatchNodePoolRequest{
				Replicas: helpers.IntPtr(3),
			},
			response:     `{"replicas": 3}`,
			statusCode:   http.StatusOK,
			wantReplicas: 3,
			wantErr:      false,
		},
		{
			name:       "invalid cluster ID",
			clusterID:  "",
			nodePoolID: "pool-456",
			wantErr:    true,
		},
		{
			name:       "invalid node pool ID",
			clusterID:  "cluster-123",
			nodePoolID: "",
			wantErr:    true,
		},
		{
			name:       "invalid cluster ID",
			clusterID:  "",
			nodePoolID: "pool-456",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "invalid node pool ID",
			clusterID:  "cluster-123",
			nodePoolID: "asfsd",
			statusCode: http.StatusOK,
			response:   `{"replica`,
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
			result, err := client.Nodepools().Update(context.Background(), tt.clusterID, tt.nodePoolID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.Replicas != tt.wantReplicas {
				t.Errorf("Update() replicas = %d, want %d", result.Replicas, tt.wantReplicas)
			}
		})
	}
}

func TestNodePoolService_Get(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		nodePoolID string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "successful get",
			clusterID:  "cluster-123",
			nodePoolID: "pool-456",
			response:   `{"id": "pool-456"}`,
			statusCode: http.StatusOK,
			wantID:     "pool-456",
			wantErr:    false,
		},
		{
			name:       "invalid node pool ID",
			clusterID:  "cluster-123",
			nodePoolID: "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "invalid cluster ID",
			clusterID:  "",
			nodePoolID: "pool-456",
			wantErr:    true,
		},
		{
			name:       "invalid cluster ID",
			nodePoolID: "pool-456",
			wantErr:    true,
		},
		{
			name:    "invalid cluster ID",
			wantErr: true,
		},
		{
			name:       "invalid node pool ID",
			clusterID:  "cluster-123",
			nodePoolID: "456456",
			statusCode: http.StatusOK,
			response:   `{"replica`,
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
			result, err := client.Nodepools().Get(context.Background(), tt.clusterID, tt.nodePoolID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.ID != tt.wantID {
				t.Errorf("Get() node pool ID = %s, want %s", result.ID, tt.wantID)
			}
		})
	}
}

func TestNodeService_List(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		nodePoolID string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name:       "successful list nodes",
			clusterID:  "cluster-123",
			nodePoolID: "pool-456",
			response: `{
				"results": [
					{"id": "node1", "name": "worker-1", "role": "worker"},
					{"id": "node2", "name": "worker-2", "role": "worker"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "list with pagination",
			clusterID:  "cluster-123",
			nodePoolID: "pool-456",
			response: `{
				"results": [
					{"id": "node1", "name": "worker-1", "role": "worker"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name:       "empty cluster ID",
			clusterID:  "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
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
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Nodepools().Nodes(context.Background(), tt.clusterID, tt.nodePoolID)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("List() got = %d nodes, want %d", len(result), tt.want)
			}
		})
	}
}

func TestNodeService_Nodes(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		nodePoolID string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name:       "successful get nodes",
			clusterID:  "cluster-123",
			nodePoolID: "pool-456",
			response: `{
				"results": [
					{"id": "node1", "name": "worker-1", "role": "worker"},
					{"id": "node2", "name": "worker-2", "role": "worker"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty cluster ID",
			clusterID:  "",
			nodePoolID: "pool-456",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "empty nodepool ID",
			clusterID:  "cluster-123",
			nodePoolID: "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			clusterID:  "cluster-123",
			nodePoolID: "pool-404",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			clusterID:  "cluster-123",
			nodePoolID: "pool-456",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID == "" || tt.nodePoolID == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Nodepools().Nodes(context.Background(), tt.clusterID, tt.nodePoolID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Nodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("Nodes() got = %d nodes, want %d", len(result), tt.want)
			}
		})
	}
}
