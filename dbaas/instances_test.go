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
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(result))
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
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, instance.ID)
		})
	}
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
			},
			response: `{
				"id": "inst1",
				"backup_retention_days": 7,
				"backup_start_at": "02:00"
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

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceClient(server.URL)
			result, err := client.Update(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
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
		snapshots, err := client.ListSnapshots(context.Background(), "inst1", ListSnapshotOptions{
			Limit: helpers.IntPtr(10),
			Type:  snapshotTypePtr(SnapshotTypeAutomated),
		})

		assertNoError(t, err)
		assertEqual(t, 1, len(snapshots))
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
