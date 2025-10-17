package blockstorage

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

func TestSnapshotService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       SnaphotListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "basic list",
			opts: SnaphotListOptions{},
			response: `{
				"snapshots": [
					{"id": "snap1", "name": "backup1"},
					{"id": "snap2", "name": "backup2"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "with pagination",
			opts: SnaphotListOptions{
				Limit:  helpers.IntPtr(1),
				Offset: helpers.IntPtr(1),
				Sort:   helpers.StrPtr("name:desc"),
			},
			response:   `{"snapshots": [{"id": "snap2"}]}`,
			statusCode: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "with expansion",
			opts: SnaphotListOptions{
				Expand: []SnapshotExpand{SnapshotVolumeExpand},
			},
			response:   `{"snapshots": [{"id": "snap1", "volume": {"id": "vol1"}}]}`,
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
				if r.URL.Path != "/volume/v1/snapshots" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				q := r.URL.Query()
				if tt.opts.Limit != nil && q.Get("_limit") != strconv.Itoa(*tt.opts.Limit) {
					t.Errorf("limit mismatch: got %s", q.Get("_limit"))
				}
				if tt.opts.Offset != nil && q.Get("_offset") != strconv.Itoa(*tt.opts.Offset) {
					t.Errorf("offset mismatch: got %s", q.Get("_offset"))
				}
				if tt.opts.Sort != nil && q.Get("_sort") != *tt.opts.Sort {
					t.Errorf("sort mismatch: got %s", q.Get("_sort"))
				}
				if len(tt.opts.Expand) > 0 && q.Get("expand") != strings.Join(tt.opts.Expand, ",") {
					t.Errorf("expand mismatch: got %s", q.Get("expand"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientSnaphots(server.URL)
			resp, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Snapshots) != tt.wantCount {
				t.Errorf("got %d snapshots, want %d", len(resp.Snapshots), tt.wantCount)
			}
		})
	}
}

func TestSnapshotService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateSnapshotRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateSnapshotRequest{
				Name:        "backup",
				Volume:      &IDOrName{ID: helpers.StrPtr("vol1")},
				Description: helpers.StrPtr("test backup"),
				Type:        helpers.StrPtr("daily"),
			},
			response:   `{"id": "snap1"}`,
			statusCode: http.StatusOK,
			wantID:     "snap1",
		},
		{
			name: "missing volume",
			request: CreateSnapshotRequest{
				Name: "invalid",
			},
			response:   `{"error": "volume required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/volume/v1/snapshots" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				var req CreateSnapshotRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("error decoding request: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientSnaphots(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id != tt.wantID {
				t.Errorf("got ID %q, want %q", id, tt.wantID)
			}
		})
	}
}

func TestSnapshotService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expand     []string
		response   string
		statusCode int
		want       *Snapshot
		wantErr    bool
	}{
		{
			name: "existing snapshot",
			id:   "snap1",
			response: `{
				"id": "snap1",
				"name": "backup",
				"status": "completed"
			}`,
			statusCode: http.StatusOK,
			want: &Snapshot{
				ID:     "snap1",
				Name:   "backup",
				Status: SnapshotStatusCompleted,
			},
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
				expectedPath := fmt.Sprintf("/volume/v1/snapshots/%s", tt.id)
				if r.URL.Path != expectedPath {
					t.Errorf("unexpected path: got %s, want %s", r.URL.Path, expectedPath)
				}

				if len(tt.expand) > 0 {
					gotExpand := r.URL.Query().Get("expand")
					if gotExpand != strings.Join(tt.expand, ",") {
						t.Errorf("expand mismatch: got %s", gotExpand)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientSnaphots(server.URL)
			snap, err := client.Get(context.Background(), tt.id, tt.expand)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if snap.ID != tt.want.ID || snap.Name != tt.want.Name || snap.Status != tt.want.Status {
				t.Errorf("got %+v, want %+v", snap, tt.want)
			}
		})
	}
}

func TestSnapshotService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "snap1",
			statusCode: http.StatusNoContent,
		},
		{
			name:       "not found",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/volume/v1/snapshots/%s", tt.id)
				if r.URL.Path != expectedPath {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testClientSnaphots(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSnapshotService_Rename(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		newName    string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful rename",
			id:         "snap1",
			newName:    "new-name",
			statusCode: http.StatusOK,
		},
		{
			name:       "name conflict",
			id:         "snap1",
			newName:    "existing",
			statusCode: http.StatusConflict,
			response:   `{"error": "name exists"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/volume/v1/snapshots/%s/rename", tt.id)
				if r.URL.Path != expectedPath {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				var req RenameSnapshotRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("error decoding request: %v", err)
				}

				if req.Name != tt.newName {
					t.Errorf("got name %q, want %q", req.Name, tt.newName)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientSnaphots(server.URL)
			err := client.Rename(context.Background(), tt.id, tt.newName)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSnapshotService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		pages      []string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page with all results",
			pages: []string{
				`{
					"snapshots": [
						{
							"id": "snap1",
							"name": "backup1",
							"size": 10,
							"state": "available",
							"status": "completed",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z",
							"availability_zones": ["az1"],
							"type": "standard"
						},
						{
							"id": "snap2",
							"name": "backup2",
							"size": 20,
							"state": "available",
							"status": "completed",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z",
							"availability_zones": ["az1"],
							"type": "standard"
						}
					]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "multiple pages",
			pages: []string{
				`{
					"snapshots": [` + generateSnapshotResults(1, 50) + `]
				}`,
				`{
					"snapshots": [` + generateSnapshotResults(51, 25) + `]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  75,
			wantErr:    false,
		},
		{
			name: "empty results",
			pages: []string{
				`{
					"snapshots": []
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageIndex := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/volume/v1/snapshots" {
					t.Errorf("Expected path /volume/v1/snapshots, got %s", r.URL.Path)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if pageIndex < len(tt.pages) {
					w.Write([]byte(tt.pages[pageIndex]))
					pageIndex++
				} else {
					w.Write([]byte(`{"snapshots":[]}`))
				}
			}))
			defer server.Close()

			client := testClientSnaphots(server.URL)
			snapshots, err := client.ListAll(context.Background(), SnapshotFilterOptions{})

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(snapshots) != tt.wantCount {
				t.Errorf("Expected %d snapshots, got %d", tt.wantCount, len(snapshots))
			}
		})
	}
}

func generateSnapshotResults(start, count int) string {
	results := make([]string, count)
	for i := 0; i < count; i++ {
		id := start + i
		results[i] = `{
			"id": "snap` + strconv.Itoa(id) + `",
			"name": "backup` + strconv.Itoa(id) + `",
			"size": 10,
			"state": "available",
			"status": "completed",
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z",
			"availability_zones": ["az1"],
			"type": "standard"
		}`
	}
	return strings.Join(results, ",")
}

func TestSnapshotService_ListAll_WithExpand(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/volume/v1/snapshots" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		query := r.URL.Query()

		// Verify expand parameter is present
		expandValue := query.Get("expand")
		if expandValue != "volume" {
			t.Errorf("expected expand=volume, got %s", expandValue)
		}

		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Return 50 items on first page, 25 on second
		offset := query.Get("_offset")
		switch offset {
		case "0":
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 50, "count": 50, "total": 75, "max_limit": 100}},
				"snapshots": [%s]
			}`, generateSnapshotResults(0, 50))
			w.Write([]byte(response))
		case "50":
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 50, "count": 25, "total": 75, "max_limit": 100}},
				"snapshots": [%s]
			}`, generateSnapshotResults(50, 25))
			w.Write([]byte(response))
		default:
			t.Errorf("unexpected offset: %s", offset)
		}
	}))
	defer server.Close()

	client := testClientSnaphots(server.URL)
	snapshots, err := client.ListAll(context.Background(), SnapshotFilterOptions{
		Expand: []string{SnapshotVolumeExpand},
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Should have fetched all 75 snapshots
	if len(snapshots) != 75 {
		t.Errorf("expected 75 snapshots, got %d", len(snapshots))
	}

	// Should have made exactly 2 requests
	if requestCount != 2 {
		t.Errorf("made %d requests, want 2", requestCount)
	}
}

func TestSnapshotService_ListAll_NewRequestError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testClientSnaphots("http://dummy-url")

	_, err := client.ListAll(ctx, SnapshotFilterOptions{})

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func testClientSnaphots(baseURL string) SnapshotService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Snapshots()
}
