package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func testClusterClient(baseURL string) ClusterService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Clusters()
}

func TestClusterService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListClustersOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "basic list",
			response: `{
				"results": [
					{
						"id": "cluster-1",
						"name": "test-cluster-1",
						"engine_id": "postgres-13",
						"instance_type_id": "db.t3.micro",
						"parameter_group_id": "pg-1",
						"volume": {"size": 20, "type": "gp2"},
						"status": "ACTIVE",
						"addresses": [{"access": "PUBLIC", "type": "READ_WRITE", "address": "test1.db.example.com", "port": "5432"}],
						"apply_parameters_pending": false,
						"backup_retention_days": 7,
						"backup_start_at": "01:00",
						"created_at": "2023-01-01T00:00:00Z"
					},
					{
						"id": "cluster-2",
						"name": "test-cluster-2",
						"engine_id": "mysql-8",
						"instance_type_id": "db.t3.small",
						"parameter_group_id": "pg-2",
						"volume": {"size": 50, "type": "gp2"},
						"status": "PENDING",
						"addresses": [{"access": "PUBLIC", "type": "READ_WRITE", "address": "test2.db.example.com", "port": "3306"}],
						"apply_parameters_pending": false,
						"backup_retention_days": 7,
						"backup_start_at": "02:00",
						"created_at": "2023-01-02T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "with filters",
			opts: ListClustersOptions{
				Limit:            helpers.IntPtr(10),
				Offset:           helpers.IntPtr(5),
				Status:           Ptr(ClusterStatusActive),
				EngineID:         helpers.StrPtr("postgres-13"),
				VolumeSize:       helpers.IntPtr(20),
				VolumeSizeGt:     helpers.IntPtr(10),
				VolumeSizeGte:    helpers.IntPtr(15),
				VolumeSizeLt:     helpers.IntPtr(100),
				VolumeSizeLte:    helpers.IntPtr(50),
				ParameterGroupID: helpers.StrPtr("pg-1"),
			},
			response: `{
				"results": [
					{
						"id": "cluster-1",
						"name": "test-cluster-1",
						"engine_id": "postgres-13",
						"status": "ACTIVE"
					}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
			wantErr:    false,
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
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/database/v2/clusters", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				query := r.URL.Query()
				if tt.opts.Limit != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Limit), query.Get("_limit"))
				}
				if tt.opts.Offset != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Offset), query.Get("_offset"))
				}
				if tt.opts.Status != nil {
					assertEqual(t, string(*tt.opts.Status), query.Get("status"))
				}
				if tt.opts.EngineID != nil {
					assertEqual(t, *tt.opts.EngineID, query.Get("engine_id"))
				}
				if tt.opts.VolumeSize != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.VolumeSize), query.Get("volume.size"))
				}
				if tt.opts.VolumeSizeGt != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.VolumeSizeGt), query.Get("volume.size__gt"))
				}
				if tt.opts.VolumeSizeGte != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.VolumeSizeGte), query.Get("volume.size__gte"))
				}
				if tt.opts.VolumeSizeLt != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.VolumeSizeLt), query.Get("volume.size__lt"))
				}
				if tt.opts.VolumeSizeLte != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.VolumeSizeLte), query.Get("volume.size__lte"))
				}
				if tt.opts.ParameterGroupID != nil {
					assertEqual(t, *tt.opts.ParameterGroupID, query.Get("parameter_group_id"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			result, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(result))
		})
	}
}

func TestClusterService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    ClusterCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: ClusterCreateRequest{
				Name:           "test-cluster",
				EngineID:       "postgres-13",
				InstanceTypeID: "db.t3.micro",
				User:           "admin",
				Password:       "password123",
				Volume: ClusterVolumeRequest{
					Size: 20,
					Type: helpers.StrPtr("gp2"),
				},
				ParameterGroupID:    helpers.StrPtr("pg-1"),
				BackupRetentionDays: helpers.IntPtr(7),
				BackupStartAt:       helpers.StrPtr("01:00"),
			},
			response: `{
				"id": "cluster-1"
			}`,
			statusCode: http.StatusAccepted,
			wantID:     "cluster-1",
			wantErr:    false,
		},
		{
			name: "validation error",
			request: ClusterCreateRequest{
				Name:     "test-cluster",
				EngineID: "invalid-engine",
			},
			response:   `{"error": "validation failed"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/database/v2/clusters", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var reqBody ClusterCreateRequest
				json.NewDecoder(r.Body).Decode(&reqBody)
				assertEqual(t, tt.request.Name, reqBody.Name)
				assertEqual(t, tt.request.EngineID, reqBody.EngineID)
				assertEqual(t, tt.request.InstanceTypeID, reqBody.InstanceTypeID)
				assertEqual(t, tt.request.User, reqBody.User)
				assertEqual(t, tt.request.Password, reqBody.Password)
				assertEqual(t, tt.request.Volume.Size, reqBody.Volume.Size)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			result, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
		})
	}
}

func TestClusterService_Get(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:      "existing cluster",
			clusterID: "cluster-1",
			response: `{
				"id": "cluster-1",
				"name": "test-cluster",
				"engine_id": "postgres-13",
				"instance_type_id": "db.t3.micro",
				"parameter_group_id": "pg-1",
				"volume": {"size": 20, "type": "gp2"},
				"status": "ACTIVE",
				"addresses": [{"access": "PUBLIC", "type": "READ_WRITE", "address": "test.db.example.com", "port": "5432"}],
				"apply_parameters_pending": false,
				"backup_retention_days": 7,
				"backup_start_at": "01:00",
				"created_at": "2023-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			wantID:     "cluster-1",
			wantErr:    false,
		},
		{
			name:       "not found",
			clusterID:  "invalid",
			response:   `{"error": "cluster not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID != "" {
					expectedPath := fmt.Sprintf("/database/v2/clusters/%s", tt.clusterID)
					assertEqual(t, expectedPath, r.URL.Path)
					assertEqual(t, http.MethodGet, r.Method)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			result, err := client.Get(context.Background(), tt.clusterID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
		})
	}
}

func TestClusterService_Update(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		request    ClusterUpdateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:      "successful update",
			clusterID: "cluster-1",
			request: ClusterUpdateRequest{
				ParameterGroupID:    helpers.StrPtr("pg-2"),
				BackupRetentionDays: helpers.IntPtr(14),
				BackupStartAt:       helpers.StrPtr("02:00"),
			},
			response: `{
				"id": "cluster-1",
				"name": "test-cluster",
				"parameter_group_id": "pg-2",
				"backup_retention_days": 14,
				"backup_start_at": "02:00",
				"status": "ACTIVE"
			}`,
			statusCode: http.StatusOK,
			wantID:     "cluster-1",
			wantErr:    false,
		},
		{
			name:       "not found",
			clusterID:  "invalid",
			request:    ClusterUpdateRequest{ParameterGroupID: helpers.StrPtr("pg-2")},
			response:   `{"error": "cluster not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID != "" {
					expectedPath := fmt.Sprintf("/database/v2/clusters/%s", tt.clusterID)
					assertEqual(t, expectedPath, r.URL.Path)
					assertEqual(t, http.MethodPatch, r.Method)

					var reqBody ClusterUpdateRequest
					json.NewDecoder(r.Body).Decode(&reqBody)
					if tt.request.ParameterGroupID != nil {
						assertEqual(t, *tt.request.ParameterGroupID, *reqBody.ParameterGroupID)
					}
					if tt.request.BackupRetentionDays != nil {
						assertEqual(t, *tt.request.BackupRetentionDays, *reqBody.BackupRetentionDays)
					}
					if tt.request.BackupStartAt != nil {
						assertEqual(t, *tt.request.BackupStartAt, *reqBody.BackupStartAt)
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			result, err := client.Update(context.Background(), tt.clusterID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
		})
	}
}

func TestClusterService_Resize(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    ClusterResizeRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "resize simultaneously instance-type and volume",
			id:   "cluster-1",
			request: ClusterResizeRequest{
				InstanceTypeID: helpers.StrPtr("type-large"),
				Volume: &ClusterVolumeResizeRequest{
					Size: 200,
					Type: "nvme",
				},
			},
			response: `{
				"id": "cluster-1",
				"instance_type_id": "type-large",
				"volume": {"size": 200}
			}`,
			statusCode: http.StatusOK,
			wantID:     "cluster-1",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v2/clusters/%s/resize", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req ClusterResizeRequest
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, *tt.request.InstanceTypeID, *req.InstanceTypeID)
				assertEqual(t, tt.request.Volume.Size, req.Volume.Size)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			result, err := client.Resize(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
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
			name:       "successful deletion",
			clusterID:  "cluster-1",
			statusCode: http.StatusAccepted,
			wantErr:    false,
		},
		{
			name:       "not found",
			clusterID:  "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID != "" {
					expectedPath := fmt.Sprintf("/database/v2/clusters/%s", tt.clusterID)
					assertEqual(t, expectedPath, r.URL.Path)
					assertEqual(t, http.MethodDelete, r.Method)
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			err := client.Delete(context.Background(), tt.clusterID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
			} else {
				assertNoError(t, err)
			}
		})
	}
}

func TestClusterService_Start(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:      "successful start",
			clusterID: "cluster-1",
			response: `{
				"id": "cluster-1",
				"name": "test-cluster",
				"status": "STARTING",
				"started_at": "2023-01-01T01:00:00Z"
			}`,
			statusCode: http.StatusAccepted,
			wantID:     "cluster-1",
			wantErr:    false,
		},
		{
			name:       "not found",
			clusterID:  "invalid",
			response:   `{"error": "cluster not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID != "" {
					expectedPath := fmt.Sprintf("/database/v2/clusters/%s/start", tt.clusterID)
					assertEqual(t, expectedPath, r.URL.Path)
					assertEqual(t, http.MethodPost, r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			result, err := client.Start(context.Background(), tt.clusterID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
		})
	}
}

func TestClusterService_Stop(t *testing.T) {
	tests := []struct {
		name       string
		clusterID  string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:      "successful stop",
			clusterID: "cluster-1",
			response: `{
				"id": "cluster-1",
				"name": "test-cluster",
				"status": "STOPPING"
			}`,
			statusCode: http.StatusAccepted,
			wantID:     "cluster-1",
			wantErr:    false,
		},
		{
			name:       "not found",
			clusterID:  "invalid",
			response:   `{"error": "cluster not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.clusterID != "" {
					expectedPath := fmt.Sprintf("/database/v2/clusters/%s/stop", tt.clusterID)
					assertEqual(t, expectedPath, r.URL.Path)
					assertEqual(t, http.MethodPost, r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClusterClient(server.URL)
			result, err := client.Stop(context.Background(), tt.clusterID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
		})
	}
}

func Ptr[T any](v T) *T {
	return &v
}
