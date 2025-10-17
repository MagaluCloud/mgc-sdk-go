package compute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestSnapshotService_List(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		opts       SnapshotListOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
		checkQuery func(*testing.T, *http.Request)
	}{
		{
			name: "basic list",
			opts: SnapshotListOptions{},
			response: `{
				"snapshots": [
					{"id": "snap1", "name": "test1", "created_at": "` + now.Format(time.RFC3339) + `"},
					{"id": "snap2", "name": "test2", "created_at": "` + now.Format(time.RFC3339) + `"}
				],
				"meta": {
					"page": {
						"offset": 0,
						"limit": 50,
						"count": 2,
						"total": 2
					}
				}
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: SnapshotListOptions{
				Limit:  intPtr(1),
				Offset: intPtr(1),
			},
			response: `{
				"snapshots": [
					{"id": "snap2", "name": "test2", "created_at": "` + now.Format(time.RFC3339) + `"}
				],
				"meta": {
					"page": {
						"offset": 1,
						"limit": 1,
						"count": 1,
						"total": 2
					}
				}
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_limit") != "1" || r.URL.Query().Get("_offset") != "1" {
					t.Error("pagination parameters not set correctly")
				}
			},
		},
		{
			name: "with expand",
			opts: SnapshotListOptions{
				Expand: []SnapshotExpand{SnapshotImageExpand},
			},
			response: `{
				"snapshots": [
					{
						"id": "snap1",
						"name": "test1",
						"created_at": "` + now.Format(time.RFC3339) + `",
						"instance": {
							"id": "inst1",
							"image": {"id": "img1"}
						}
					}
				],
				"meta": {
					"page": {
						"offset": 0,
						"limit": 50,
						"count": 1,
						"total": 1
					}
				}
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("expand") != string(SnapshotImageExpand) {
					t.Error("expand parameter not set correctly")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkQuery != nil {
					tt.checkQuery(t, r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Snapshots().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Snapshots) != tt.want {
				t.Errorf("List() got %v snapshots, want %v", len(got.Snapshots), tt.want)
			}
			if !tt.wantErr && got.Meta.Page.Count != tt.want {
				t.Errorf("List() meta count = %v, want %v", got.Meta.Page.Count, tt.want)
			}
		})
	}
}

func TestSnapshotService_ListAll(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		opts       SnapshotFilterOptions
		responses  []string
		want       int
		wantErr    bool
		checkCalls func(*testing.T, int)
	}{
		{
			name: "single page",
			opts: SnapshotFilterOptions{},
			responses: []string{
				`{
					"snapshots": [
						{"id": "snap1", "name": "test1", "created_at": "` + now.Format(time.RFC3339) + `"},
						{"id": "snap2", "name": "test2", "created_at": "` + now.Format(time.RFC3339) + `"}
					],
					"meta": {
						"page": {
							"offset": 0,
							"limit": 50,
							"count": 2,
							"total": 2
						}
					}
				}`,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "multiple pages",
			opts: SnapshotFilterOptions{},
			responses: []string{
				`{
					"snapshots": [` +
					generateSnapshotJSON(50, 0, now) + `
					],
					"meta": {
						"page": {
							"offset": 0,
							"limit": 50,
							"count": 50,
							"total": 75
						}
					}
				}`,
				`{
					"snapshots": [` +
					generateSnapshotJSON(25, 50, now) + `
					],
					"meta": {
						"page": {
							"offset": 50,
							"limit": 50,
							"count": 25,
							"total": 75
						}
					}
				}`,
			},
			want:    75,
			wantErr: false,
			checkCalls: func(t *testing.T, calls int) {
				if calls != 2 {
					t.Errorf("expected 2 API calls, got %d", calls)
				}
			},
		},
		{
			name: "with expand",
			opts: SnapshotFilterOptions{
				Expand: []SnapshotExpand{SnapshotImageExpand},
			},
			responses: []string{
				`{
					"snapshots": [
						{
							"id": "snap1",
							"name": "test1",
							"created_at": "` + now.Format(time.RFC3339) + `",
							"instance": {
								"id": "inst1",
								"image": {"id": "img1"}
							}
						}
					],
					"meta": {
						"page": {
							"offset": 0,
							"limit": 50,
							"count": 1,
							"total": 1
						}
					}
				}`,
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if callCount >= len(tt.responses) {
					t.Errorf("unexpected API call #%d", callCount+1)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responses[callCount]))
				callCount++
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Snapshots().ListAll(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ListAll() got %v snapshots, want %v", len(got), tt.want)
			}
			if tt.checkCalls != nil {
				tt.checkCalls(t, callCount)
			}
		})
	}
}

func TestSnapshotService_Create(t *testing.T) {
	tests := []struct {
		name       string
		req        CreateSnapshotRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			req: CreateSnapshotRequest{
				Name:     "test-snapshot",
				Instance: IDOrName{ID: strPtr("inst1")},
			},
			response:   `{"id": "snap1"}`,
			statusCode: http.StatusOK,
			wantID:     "snap1",
			wantErr:    false,
		},
		{
			name: "instance not found",
			req: CreateSnapshotRequest{
				Name:     "test-snapshot",
				Instance: IDOrName{ID: strPtr("invalid")},
			},
			response:   `{"error": "instance not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "invalid request",
			req:        CreateSnapshotRequest{},
			response:   `{"error": "invalid request"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			gotID, err := client.Snapshots().Create(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("Create() got = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestSnapshotService_Get(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		id         string
		expand     []SnapshotExpand
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "existing snapshot",
			id:   "snap1",
			response: `{
				"id": "snap1",
				"name": "test-snapshot",
				"created_at": "` + now.Format(time.RFC3339) + `"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "not found",
			id:         "invalid",
			response:   `{"error": "snapshot not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:   "with expand",
			id:     "snap1",
			expand: []SnapshotExpand{SnapshotImageExpand, SnapshotMachineTypeExpand},
			response: `{
				"id": "snap1",
				"name": "test-snapshot",
				"created_at": "` + now.Format(time.RFC3339) + `",
				"instance": {
					"id": "inst1",
					"image": {"id": "img1"},
					"machine_type": {"id": "mt1"}
				}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Snapshots().Get(context.Background(), tt.id, tt.expand)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("Get() got = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

// Helper function to generate snapshot JSON for testing
func generateSnapshotJSON(count, startID int, baseTime time.Time) string {
	if count == 0 {
		return ""
	}
	var result string
	for i := 0; i < count; i++ {
		if i > 0 {
			result += ","
		}
		result += `{"id": "snap` + strconv.Itoa(startID+i+1) + `", "name": "test` + strconv.Itoa(startID+i+1) + `", "created_at": "` + baseTime.Format(time.RFC3339) + `"}`
	}
	return result
}

func TestSnapshotService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "snap1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "not found",
			id:         "invalid",
			response:   `{"error": "snapshot not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "in use",
			id:         "in-use",
			response:   `{"error": "snapshot in use"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Snapshots().Delete(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSnapshotService_Rename(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		newName    string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful rename",
			id:         "snap1",
			newName:    "new-name",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "name in use",
			id:         "snap1",
			newName:    "existing-name",
			response:   `{"error": "name already in use"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Snapshots().Rename(context.Background(), tt.id, tt.newName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSnapshotService_Restore(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		req        RestoreSnapshotRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful restore",
			id:   "snap1",
			req: RestoreSnapshotRequest{
				Name:        "restored-instance",
				MachineType: IDOrName{ID: strPtr("mt1")},
			},
			response:   `{"id": "inst1"}`,
			statusCode: http.StatusOK,
			wantID:     "inst1",
			wantErr:    false,
		},
		{
			name: "invalid machine type",
			id:   "snap1",
			req: RestoreSnapshotRequest{
				Name:        "restored-instance",
				MachineType: IDOrName{ID: strPtr("invalid")},
			},
			response:   `{"error": "invalid machine type"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			gotID, err := client.Snapshots().Restore(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("Restore() got = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestSnapshotService_Copy(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		req        CopySnapshotRequest
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful copy",
			id:   "snap1",
			req: CopySnapshotRequest{
				DestinationRegion: "region2",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid region",
			id:   "snap1",
			req: CopySnapshotRequest{
				DestinationRegion: "invalid",
			},
			response:   `{"error": "invalid region"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Snapshots().Copy(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
