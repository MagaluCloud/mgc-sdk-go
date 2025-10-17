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

func testInstanceClient(baseURL string) InstanceService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Instances()
}

func TestInstanceService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListInstanceOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "basic list",
			response: `{
				"meta": {"total": 2},
				"results": [
					{"id": "inst1", "name": "instance1"},
					{"id": "inst2", "name": "instance2"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "with filters and pagination",
			opts: ListInstanceOptions{
				Limit:          helpers.IntPtr(10),
				Offset:         helpers.IntPtr(5),
				Status:         instanceStatusPtr(InstanceStatusActive),
				EngineID:       helpers.StrPtr("postgres"),
				VolumeSize:     helpers.IntPtr(100),
				ExpandedFields: []string{"replicas", "parameters"},
			},
			response: `{
				"meta": {"total": 1},
				"results": [{"id": "inst1", "status": "ACTIVE"}]
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
				assertEqual(t, "/database/v2/instances", r.URL.Path)
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
				if len(tt.opts.ExpandedFields) > 0 {
					assertEqual(t, strings.Join(tt.opts.ExpandedFields, ","), query.Get("_expand"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			result, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(result.Results))
		})
	}
}

func instanceStatusPtr(InstanceStatusActive InstanceStatus) *InstanceStatus {
	return &InstanceStatusActive
}

func TestInstanceService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		opts       GetInstanceOptions
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "existing instance",
			id:   "inst1",
			opts: GetInstanceOptions{
				ExpandedFields: []string{"replicas"},
			},
			response: `{
				"id": "inst1",
				"name": "test-instance",
				"status": "ACTIVE"
			}`,
			statusCode: http.StatusOK,
			wantID:     "inst1",
			wantErr:    false,
		},
		{
			name:       "not found",
			id:         "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v2/instances/%s", tt.id), r.URL.Path)
				query := r.URL.Query()

				if len(tt.opts.ExpandedFields) > 0 {
					assertEqual(t, strings.Join(tt.opts.ExpandedFields, ","), query.Get("_expand"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			instance, err := client.Get(context.Background(), tt.id, tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, instance.ID)
		})
	}
}

func TestInstanceService_List_PaginationMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/instances", r.URL.Path)
		query := r.URL.Query()

		assertEqual(t, "10", query.Get("_limit"))
		assertEqual(t, "20", query.Get("_offset"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"meta": {
				"page": {
					"offset": 20,
					"limit": 10,
					"count": 10,
					"total": 50,
					"max_limit": 100
				},
				"filters": [
					{"field": "status", "value": "ACTIVE"}
				]
			},
			"results": [
				{"id": "inst1", "name": "instance1", "status": "ACTIVE"},
				{"id": "inst2", "name": "instance2", "status": "ACTIVE"}
			]
		}`))
	}))
	defer server.Close()

	client := testInstanceClient(server.URL)
	offset := 20
	limit := 10
	status := InstanceStatusActive
	result, err := client.List(context.Background(), ListInstanceOptions{
		Offset: &offset,
		Limit:  &limit,
		Status: &status,
	})

	assertNoError(t, err)

	// Validate results
	assertEqual(t, 2, len(result.Results))
	assertEqual(t, "inst1", result.Results[0].ID)
	assertEqual(t, "inst2", result.Results[1].ID)

	// Validate pagination metadata
	assertEqual(t, 20, result.Meta.Page.Offset)
	assertEqual(t, 10, result.Meta.Page.Limit)
	assertEqual(t, 10, result.Meta.Page.Count)
	assertEqual(t, 50, result.Meta.Page.Total)
	assertEqual(t, 100, result.Meta.Page.MaxLimit)

	// Validate filters metadata
	assertEqual(t, 1, len(result.Meta.Filters))
	assertEqual(t, "status", result.Meta.Filters[0].Field)
	assertEqual(t, "ACTIVE", result.Meta.Filters[0].Value)
}

func TestInstanceService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    InstanceCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: InstanceCreateRequest{
				Name:     "new-instance",
				User:     "admin",
				Password: "secret",
				Volume: InstanceVolumeRequest{
					Size: 100,
					Type: "nvme",
				},
			},
			response:   `{"id": "inst-new"}`,
			statusCode: http.StatusOK,
			wantID:     "inst-new",
			wantErr:    false,
		},
		{
			name: "invalid request",
			request: InstanceCreateRequest{
				Name: "missing-password",
			},
			response:   `{"error": "password required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/database/v2/instances", r.URL.Path)

				var req InstanceCreateRequest
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, tt.request.Name, req.Name)
				assertEqual(t, tt.request.User, req.User)
				assertEqual(t, tt.request.Password, req.Password)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
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

func TestInstanceService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "inst1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "not found",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v2/instances/%s", tt.id), r.URL.Path)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestInstanceService_Update(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    DatabaseInstanceUpdateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "update backup settings",
			id:   "inst1",
			request: DatabaseInstanceUpdateRequest{
				BackupRetentionDays: helpers.IntPtr(7),
				BackupStartAt:       helpers.StrPtr("02:00"),
				ParameterGroupID:    helpers.StrPtr("pg-id"),
			},
			response: `{
				"id": "inst1",
				"backup_retention_days": 7,
				"backup_start_at": "02:00",
				"parameter_group_id": "pg-id"
			}`,
			statusCode: http.StatusOK,
			wantID:     "inst1",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v2/instances/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPatch, r.Method)

				var req DatabaseInstanceUpdateRequest
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, *tt.request.BackupRetentionDays, *req.BackupRetentionDays)
				assertEqual(t, *tt.request.BackupStartAt, *req.BackupStartAt)
				assertEqual(t, *tt.request.ParameterGroupID, *req.ParameterGroupID)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			result, err := client.Update(context.Background(), tt.id, tt.request)

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

func TestInstanceService_Resize(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    InstanceResizeRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "resize instance type",
			id:   "inst1",
			request: InstanceResizeRequest{
				InstanceTypeID: helpers.StrPtr("type-large"),
				Volume: &InstanceVolumeResizeRequest{
					Size: 200,
					Type: "nvme",
				},
			},
			response: `{
				"id": "inst1",
				"instance_type_id": "type-large",
				"volume": {"size": 200}
			}`,
			statusCode: http.StatusOK,
			wantID:     "inst1",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v2/instances/%s/resize", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req InstanceResizeRequest
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, *tt.request.InstanceTypeID, *req.InstanceTypeID)
				assertEqual(t, tt.request.Volume.Size, req.Volume.Size)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
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

func TestInstanceService_StartStop(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         string
		response   string
		statusCode int
		wantStatus InstanceStatus
		wantErr    bool
	}{
		{
			name:       "start instance",
			method:     "Start",
			id:         "inst1",
			response:   `{"id": "inst1", "status": "STARTING"}`,
			statusCode: http.StatusOK,
			wantStatus: InstanceStatusStarting,
			wantErr:    false,
		},
		{
			name:       "stop instance",
			method:     "Stop",
			id:         "inst1",
			response:   `{"id": "inst1", "status": "STOPPING"}`,
			statusCode: http.StatusOK,
			wantStatus: InstanceStatusStopping,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var expectedPath string
				switch tt.method {
				case "Start":
					expectedPath = fmt.Sprintf("/database/v2/instances/%s/start", tt.id)
				case "Stop":
					expectedPath = fmt.Sprintf("/database/v2/instances/%s/stop", tt.id)
				}
				assertEqual(t, expectedPath, r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			var result *InstanceDetail
			var err error

			switch tt.method {
			case "Start":
				result, err = client.Start(context.Background(), tt.id)
			case "Stop":
				result, err = client.Stop(context.Background(), tt.id)
			}

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantStatus, result.Status)
		})
	}
}

func TestInstanceService_Snapshots(t *testing.T) {
	t.Run("ListSnapshots", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/database/v2/instances/inst1/snapshots", r.URL.Path)
			query := r.URL.Query()
			assertEqual(t, "10", query.Get("_limit"))
			assertEqual(t, "AUTOMATED", query.Get("type"))

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"meta": {"total": 1}, "results": [{"id": "snap1"}]}`))
		}))
		defer server.Close()

		client := testInstanceClient(server.URL)
		result, err := client.ListSnapshots(context.Background(), "inst1", ListSnapshotOptions{
			Limit: helpers.IntPtr(10),
			Type:  snapshotTypePtr(SnapshotTypeAutomated),
		})

		assertNoError(t, err)
		assertEqual(t, 1, len(result.Results))
	})

	t.Run("CreateSnapshot", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/database/v2/instances/inst1/snapshots", r.URL.Path)

			var req SnapshotCreateRequest
			json.NewDecoder(r.Body).Decode(&req)
			assertEqual(t, "daily-backup", req.Name)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id": "snap-new"}`))
		}))
		defer server.Close()

		client := testInstanceClient(server.URL)
		result, err := client.CreateSnapshot(context.Background(), "inst1", SnapshotCreateRequest{
			Name: "daily-backup",
		})

		assertNoError(t, err)
		assertEqual(t, "snap-new", result.ID)
	})

	t.Run("GetSnapshot", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/database/v2/instances/inst1/snapshots/snap1", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id": "snap1", "status": "AVAILABLE"}`))
		}))
		defer server.Close()

		client := testInstanceClient(server.URL)
		snapshot, err := client.GetSnapshot(context.Background(), "inst1", "snap1")

		assertNoError(t, err)
		assertEqual(t, "snap1", snapshot.ID)
	})

	t.Run("DeleteSnapshot", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/database/v2/instances/inst1/snapshots/snap1", r.URL.Path)
			assertEqual(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client := testInstanceClient(server.URL)
		err := client.DeleteSnapshot(context.Background(), "inst1", "snap1")

		assertNoError(t, err)
	})

	t.Run("UpdateSnapshot", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/database/v2/instances/inst1/snapshots/snap1", r.URL.Path)
			assertEqual(t, http.MethodPatch, r.Method)

			var req SnapshotUpdateRequest
			json.NewDecoder(r.Body).Decode(&req)
			assertEqual(t, "updated-name", req.Name)
			assertEqual(t, "updated description", *req.Description)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "snap1",
				"name": "updated-name",
				"description": "updated description"
			}`))
		}))
		defer server.Close()

		client := testInstanceClient(server.URL)
		result, err := client.UpdateSnapshot(context.Background(), "inst1", "snap1", SnapshotUpdateRequest{
			Name:        "updated-name",
			Description: helpers.StrPtr("updated description"),
		})

		assertNoError(t, err)
		assertEqual(t, "snap1", result.ID)
		assertEqual(t, "updated-name", result.Name)
		assertEqual(t, "updated description", result.Description)
	})

	t.Run("RestoreSnapshot", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/database/v2/instances/inst1/snapshots/snap1/restore", r.URL.Path)
			assertEqual(t, http.MethodPost, r.Method)

			var req RestoreSnapshotRequest
			json.NewDecoder(r.Body).Decode(&req)
			assertEqual(t, "restored-instance", req.Name)
			assertEqual(t, "type-large", req.InstanceTypeID)
			assertEqual(t, 100, req.Volume.Size)
			assertEqual(t, "nvme", req.Volume.Type)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "new-inst"}`))
		}))
		defer server.Close()

		client := testInstanceClient(server.URL)
		result, err := client.RestoreSnapshot(context.Background(), "inst1", "snap1", RestoreSnapshotRequest{
			Name:           "restored-instance",
			InstanceTypeID: "type-large",
			Volume: &InstanceVolumeRequest{
				Size: 100,
				Type: "nvme",
			},
		})

		assertNoError(t, err)
		assertEqual(t, "new-inst", result.ID)
	})
}

func snapshotTypePtr(SnapshotTypeAutomated SnapshotType) *SnapshotType {
	return &SnapshotTypeAutomated
}

func TestInstanceService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		filterOpts InstanceFilterOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page",
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 2, "total": 2, "max_limit": 100}},
				"results": [
					{"id": "inst1", "name": "instance1", "status": "ACTIVE"},
					{"id": "inst2", "name": "instance2", "status": "ACTIVE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "empty result",
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 0, "total": 0, "max_limit": 100}},
				"results": []
			}`,
			statusCode: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "with status filter",
			filterOpts: InstanceFilterOptions{
				Status: instanceStatusPtr(InstanceStatusActive),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}},
				"results": [
					{"id": "inst1", "status": "ACTIVE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "with engine_id filter",
			filterOpts: InstanceFilterOptions{
				EngineID: helpers.StrPtr("postgres-16"),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}},
				"results": [
					{"id": "inst1", "engine_id": "postgres-16"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "with volume size filters",
			filterOpts: InstanceFilterOptions{
				VolumeSizeGte: helpers.IntPtr(100),
				VolumeSizeLte: helpers.IntPtr(500),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 2, "total": 2, "max_limit": 100}},
				"results": [
					{"id": "inst1", "volume": {"size": 100, "type": "nvme"}},
					{"id": "inst2", "volume": {"size": 200, "type": "nvme"}}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "with expanded fields",
			filterOpts: InstanceFilterOptions{
				ExpandedFields: []string{"replicas", "parameters"},
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}},
				"results": [
					{"id": "inst1", "replicas": []}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
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
				assertEqual(t, "/database/v2/instances", r.URL.Path)

				query := r.URL.Query()

				// Verify filter parameters
				if tt.filterOpts.Status != nil {
					assertEqual(t, string(*tt.filterOpts.Status), query.Get("status"))
				}
				if tt.filterOpts.EngineID != nil {
					assertEqual(t, *tt.filterOpts.EngineID, query.Get("engine_id"))
				}
				if tt.filterOpts.VolumeSize != nil {
					assertEqual(t, strconv.Itoa(*tt.filterOpts.VolumeSize), query.Get("volume.size"))
				}
				if tt.filterOpts.VolumeSizeGt != nil {
					assertEqual(t, strconv.Itoa(*tt.filterOpts.VolumeSizeGt), query.Get("volume.size__gt"))
				}
				if tt.filterOpts.VolumeSizeGte != nil {
					assertEqual(t, strconv.Itoa(*tt.filterOpts.VolumeSizeGte), query.Get("volume.size__gte"))
				}
				if tt.filterOpts.VolumeSizeLt != nil {
					assertEqual(t, strconv.Itoa(*tt.filterOpts.VolumeSizeLt), query.Get("volume.size__lt"))
				}
				if tt.filterOpts.VolumeSizeLte != nil {
					assertEqual(t, strconv.Itoa(*tt.filterOpts.VolumeSizeLte), query.Get("volume.size__lte"))
				}
				if len(tt.filterOpts.ExpandedFields) > 0 {
					assertEqual(t, strings.Join(tt.filterOpts.ExpandedFields, ","), query.Get("_expand"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			instances, err := client.ListAll(context.Background(), tt.filterOpts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(instances))
		})
	}
}

func TestInstanceService_ListAll_MultiplePagesWithPagination(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/instances", r.URL.Path)

		query := r.URL.Query()
		offset := query.Get("_offset")
		limit := query.Get("_limit")

		if limit != "50" {
			t.Errorf("expected limit 50, got %s", limit)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if requestCount == 0 {
			// First page
			if offset != "0" {
				t.Errorf("expected offset 0, got %s", offset)
			}
			results := `[`
			for i := 0; i < 50; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "inst-%d", "name": "instance-%d", "status": "ACTIVE"}`, i+1, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 50, "total": 70, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		} else if requestCount == 1 {
			// Second page
			if offset != "50" {
				t.Errorf("expected offset 50, got %s", offset)
			}
			results := `[`
			for i := 0; i < 20; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "inst-%d", "name": "instance-%d", "status": "ACTIVE"}`, i+51, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 50, "count": 20, "total": 70, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		}

		requestCount++
	}))
	defer server.Close()

	client := testInstanceClient(server.URL)
	instances, err := client.ListAll(context.Background(), InstanceFilterOptions{})

	assertNoError(t, err)
	assertEqual(t, 70, len(instances))
	assertEqual(t, 2, requestCount)
}

func TestInstanceService_ListAll_WithFilters(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/instances", r.URL.Path)

		query := r.URL.Query()

		// Verify filter parameters are present
		if query.Get("status") != "ACTIVE" {
			t.Errorf("expected status=ACTIVE, got %s", query.Get("status"))
		}
		if query.Get("engine_id") != "postgres-16" {
			t.Errorf("expected engine_id=postgres-16, got %s", query.Get("engine_id"))
		}
		if query.Get("_expand") != "replicas" {
			t.Errorf("expected _expand=replicas, got %s", query.Get("_expand"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if requestCount == 0 {
			// First page with 50 results
			results := `[`
			for i := 0; i < 50; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "inst-%d", "status": "ACTIVE", "engine_id": "postgres-16", "replicas": []}`, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 50, "total": 55, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		} else if requestCount == 1 {
			// Second page with 5 results
			results := `[`
			for i := 0; i < 5; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "inst-%d", "status": "ACTIVE", "engine_id": "postgres-16", "replicas": []}`, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 50, "count": 5, "total": 55, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		}

		requestCount++
	}))
	defer server.Close()

	client := testInstanceClient(server.URL)
	instances, err := client.ListAll(context.Background(), InstanceFilterOptions{
		Status:         instanceStatusPtr(InstanceStatusActive),
		EngineID:       helpers.StrPtr("postgres-16"),
		ExpandedFields: []string{"replicas"},
	})

	assertNoError(t, err)
	assertEqual(t, 55, len(instances))
	assertEqual(t, 2, requestCount)

	// Verify all instances have the correct status
	for _, inst := range instances {
		assertEqual(t, InstanceStatusActive, inst.Status)
	}
}

func TestInstanceService_ListAllSnapshots(t *testing.T) {
	tests := []struct {
		name       string
		instanceID string
		filterOpts SnapshotFilterOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "single page",
			instanceID: "inst-123",
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 2, "total": 2, "max_limit": 100}},
				"results": [
					{"id": "snap1", "name": "snapshot1", "status": "AVAILABLE", "type": "ON_DEMAND"},
					{"id": "snap2", "name": "snapshot2", "status": "AVAILABLE", "type": "AUTOMATED"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "empty result",
			instanceID: "inst-123",
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 0, "total": 0, "max_limit": 100}},
				"results": []
			}`,
			statusCode: http.StatusOK,
			wantCount:  0,
		},
		{
			name:       "with type filter",
			instanceID: "inst-123",
			filterOpts: SnapshotFilterOptions{
				Type: snapshotTypePtr(SnapshotTypeOnDemand),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}},
				"results": [
					{"id": "snap1", "type": "ON_DEMAND", "status": "AVAILABLE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
		},
		{
			name:       "with status filter",
			instanceID: "inst-123",
			filterOpts: SnapshotFilterOptions{
				Status: snapshotStatusPtr(SnapshotStatusAvailable),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 2, "total": 2, "max_limit": 100}},
				"results": [
					{"id": "snap1", "status": "AVAILABLE"},
					{"id": "snap2", "status": "AVAILABLE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "server error",
			instanceID: "inst-123",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v2/instances/%s/snapshots", tt.instanceID), r.URL.Path)

				query := r.URL.Query()

				// Verify filter parameters
				if tt.filterOpts.Type != nil {
					assertEqual(t, string(*tt.filterOpts.Type), query.Get("type"))
				}
				if tt.filterOpts.Status != nil {
					assertEqual(t, string(*tt.filterOpts.Status), query.Get("status"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			snapshots, err := client.ListAllSnapshots(context.Background(), tt.instanceID, tt.filterOpts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(snapshots))
		})
	}
}

func TestInstanceService_ListAllSnapshots_MultiplePagesWithPagination(t *testing.T) {
	instanceID := "inst-123"
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, fmt.Sprintf("/database/v2/instances/%s/snapshots", instanceID), r.URL.Path)

		query := r.URL.Query()
		offset := query.Get("_offset")
		limit := query.Get("_limit")

		if limit != "50" {
			t.Errorf("expected limit 50, got %s", limit)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if requestCount == 0 {
			// First page
			if offset != "0" {
				t.Errorf("expected offset 0, got %s", offset)
			}
			results := `[`
			for i := 0; i < 50; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "snap-%d", "name": "snapshot-%d", "status": "AVAILABLE", "type": "AUTOMATED"}`, i+1, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 50, "total": 75, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		} else if requestCount == 1 {
			// Second page
			if offset != "50" {
				t.Errorf("expected offset 50, got %s", offset)
			}
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "snap-%d", "name": "snapshot-%d", "status": "AVAILABLE", "type": "AUTOMATED"}`, i+51, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 50, "count": 25, "total": 75, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		}

		requestCount++
	}))
	defer server.Close()

	client := testInstanceClient(server.URL)
	snapshots, err := client.ListAllSnapshots(context.Background(), instanceID, SnapshotFilterOptions{})

	assertNoError(t, err)
	assertEqual(t, 75, len(snapshots))
	assertEqual(t, 2, requestCount)
}

func TestInstanceService_ListAllSnapshots_WithFilters(t *testing.T) {
	instanceID := "inst-123"
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, fmt.Sprintf("/database/v2/instances/%s/snapshots", instanceID), r.URL.Path)

		query := r.URL.Query()

		// Verify filter parameters are present
		if query.Get("type") != "ON_DEMAND" {
			t.Errorf("expected type=ON_DEMAND, got %s", query.Get("type"))
		}
		if query.Get("status") != "AVAILABLE" {
			t.Errorf("expected status=AVAILABLE, got %s", query.Get("status"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if requestCount == 0 {
			// First page with 50 results
			results := `[`
			for i := 0; i < 50; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "snap-%d", "status": "AVAILABLE", "type": "ON_DEMAND"}`, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 50, "total": 60, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		} else if requestCount == 1 {
			// Second page with 10 results
			results := `[`
			for i := 0; i < 10; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "snap-%d", "status": "AVAILABLE", "type": "ON_DEMAND"}`, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 50, "count": 10, "total": 60, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		}

		requestCount++
	}))
	defer server.Close()

	client := testInstanceClient(server.URL)
	snapshots, err := client.ListAllSnapshots(context.Background(), instanceID, SnapshotFilterOptions{
		Type:   snapshotTypePtr(SnapshotTypeOnDemand),
		Status: snapshotStatusPtr(SnapshotStatusAvailable),
	})

	assertNoError(t, err)
	assertEqual(t, 60, len(snapshots))
	assertEqual(t, 2, requestCount)

	// Verify all snapshots have the correct type and status
	for _, snap := range snapshots {
		assertEqual(t, SnapshotTypeOnDemand, snap.Type)
		assertEqual(t, SnapshotStatusAvailable, snap.Status)
	}
}

func snapshotStatusPtr(status SnapshotStatus) *SnapshotStatus {
	return &status
}
