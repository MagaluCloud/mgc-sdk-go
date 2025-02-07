package blockstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

// Helper functions for testing
func assertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v but got %v. %v", expected, actual, msgAndArgs)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestVolumeService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "basic list",
			opts: ListOptions{},
			response: `{
				"volumes": [
					{"id": "vol1", "name": "test1"},
					{"id": "vol2", "name": "test2"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: ListOptions{
				Limit:  helpers.IntPtr(1),
				Offset: helpers.IntPtr(1),
			},
			response: `{
				"volumes": [
					{"id": "vol2", "name": "test2"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name: "with expansion",
			opts: ListOptions{
				Expand: []string{VolumeTypeExpand, VolumeAttachExpand},
			},
			response: `{
				"volumes": [
					{"id": "vol1", "type": {"id": "type1"}, "attachment": {"instance": {"id": "inst1"}}}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
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
				assertEqual(t, "/volume/v1/volumes", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			volumes, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(volumes))
		})
	}
}

func TestVolumeService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateVolumeRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateVolumeRequest{
				Name: "test-vol",
				Size: 100,
				Type: IDOrName{Name: helpers.StrPtr("ssd")},
			},
			response:   `{"id": "vol1"}`,
			statusCode: http.StatusOK,
			wantID:     "vol1",
			wantErr:    false,
		},
		{
			name: "invalid size",
			request: CreateVolumeRequest{
				Name: "test-vol",
				Size: 0,
			},
			response:   `{"error": "invalid size"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "quota exceeded",
			request: CreateVolumeRequest{
				Name: "test-vol",
				Size: 1000,
			},
			response:   `{"error": "quota exceeded"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/volumes", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestVolumeService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expand     []string
		response   string
		statusCode int
		want       *Volume
		wantErr    bool
	}{
		{
			name: "existing volume",
			id:   "vol1",
			response: `{
				"id": "vol1",
				"name": "test-vol",
				"size": 100,
				"status": "completed"
			}`,
			statusCode: http.StatusOK,
			want: &Volume{
				ID:     "vol1",
				Name:   "test-vol",
				Size:   100,
				Status: "completed",
			},
			wantErr: false,
		},
		{
			name:       "not found",
			id:         "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:   "with expansion",
			id:     "vol1",
			expand: []string{VolumeTypeExpand},
			response: `{
				"id": "vol1",
				"type": {"id": "type1"},
				"status": "completed"
			}`,
			statusCode: http.StatusOK,
			want: &Volume{
				ID:     "vol1",
				Status: "completed",
				Type:   Type{ID: "type1"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/volumes/"+tt.id, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			volume, err := client.Get(context.Background(), tt.id, tt.expand)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			if volume.ID != tt.want.ID || volume.Status != tt.want.Status {
				t.Errorf("Got volume %+v, want %+v", volume, tt.want)
			}
		})
	}
}

func TestVolumeService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "vol1",
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
		{
			name:       "attached volume",
			id:         "vol-attached",
			statusCode: http.StatusConflict,
			response:   `{"error": "volume is attached"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/volumes/"+tt.id, r.URL.Path)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestVolumeService_Rename(t *testing.T) {
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
			id:         "vol1",
			newName:    "new-name",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty name",
			id:         "vol1",
			newName:    "",
			statusCode: http.StatusBadRequest,
			response:   `{"error": "name required"}`,
			wantErr:    true,
		},
		{
			name:       "duplicate name",
			id:         "vol1",
			newName:    "existing",
			statusCode: http.StatusConflict,
			response:   `{"error": "name exists"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/volumes/"+tt.id+"/rename", r.URL.Path)

				var req RenameVolumeRequest
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, tt.newName, req.Name)

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Rename(context.Background(), tt.id, tt.newName)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestVolumeService_Extend(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    ExtendVolumeRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful extend",
			id:         "vol1",
			request:    ExtendVolumeRequest{Size: 200},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "smaller size",
			id:         "vol1",
			request:    ExtendVolumeRequest{Size: 50},
			statusCode: http.StatusBadRequest,
			response:   `{"error": "size too small"}`,
			wantErr:    true,
		},
		{
			name:       "attached volume",
			id:         "vol-attached",
			request:    ExtendVolumeRequest{Size: 200},
			statusCode: http.StatusConflict,
			response:   `{"error": "volume attached"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/volumes/"+tt.id+"/extend", r.URL.Path)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Extend(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestVolumeService_Retype(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    RetypeVolumeRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "successful retype",
			id:   "vol1",
			request: RetypeVolumeRequest{
				NewType: IDOrName{ID: helpers.StrPtr("type2")},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid type",
			id:   "vol1",
			request: RetypeVolumeRequest{
				NewType: IDOrName{},
			},
			statusCode: http.StatusBadRequest,
			response:   `{"error": "invalid type"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/volumes/"+tt.id+"/retype", r.URL.Path)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Retype(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestVolumeService_AttachDetach(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		volumeID   string
		instanceID string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful attach",
			method:     "Attach",
			volumeID:   "vol1",
			instanceID: "inst1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "already attached",
			method:     "Attach",
			volumeID:   "vol-attached",
			instanceID: "inst1",
			statusCode: http.StatusConflict,
			response:   `{"error": "already attached"}`,
			wantErr:    true,
		},
		{
			name:       "successful detach",
			method:     "Detach",
			volumeID:   "vol1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "not attached",
			method:     "Detach",
			volumeID:   "vol-unattached",
			statusCode: http.StatusConflict,
			response:   `{"error": "not attached"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.method == "Attach" {
					assertEqual(t, fmt.Sprintf("/volume/v1/volumes/%s/attach/%s", tt.volumeID, tt.instanceID), r.URL.Path)
				} else {
					assertEqual(t, fmt.Sprintf("/volume/v1/volumes/%s/detach", tt.volumeID), r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			var err error

			if tt.method == "Attach" {
				err = client.Attach(context.Background(), tt.volumeID, tt.instanceID)
			} else {
				err = client.Detach(context.Background(), tt.volumeID)
			}

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

// Helper functions
func testClient(baseURL string) VolumeService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Volumes()
}
