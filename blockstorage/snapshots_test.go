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
		opts       ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "basic list",
			opts: ListOptions{},
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
			opts: ListOptions{
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
			opts: ListOptions{
				Expand: []string{SnapshotVolumeExpand},
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
			snapshots, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(snapshots) != tt.wantCount {
				t.Errorf("got %d snapshots, want %d", len(snapshots), tt.wantCount)
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
				Description: "test backup",
				Type:        "daily",
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
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func testClientSnaphots(baseURL string) SnapshotService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Snapshots()
}
