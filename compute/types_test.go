package compute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMachineTypeService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       InstanceTypeListOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
		checkQuery func(*testing.T, *http.Request)
	}{
		{
			name: "basic list",
			opts: InstanceTypeListOptions{},
			response: `{
				"instance_types": [
					{"id": "mt1", "name": "small", "vcpus": 2, "ram": 4096, "disk": 50},
					{"id": "mt2", "name": "medium", "vcpus": 4, "ram": 8192, "disk": 100}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: InstanceTypeListOptions{
				Limit:  intPtr(1),
				Offset: intPtr(1),
			},
			response: `{
				"instance_types": [
					{"id": "mt2", "name": "medium", "vcpus": 4, "ram": 8192, "disk": 100}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_limit") != "1" {
					t.Errorf("expected limit=1, got %s", r.URL.Query().Get("_limit"))
				}
				if r.URL.Query().Get("_offset") != "1" {
					t.Errorf("expected offset=1, got %s", r.URL.Query().Get("_offset"))
				}
			},
		},
		{
			name: "with sorting",
			opts: InstanceTypeListOptions{
				Sort: strPtr("vcpus:asc"),
			},
			response: `{
				"instance_types": [
					{"id": "mt1", "name": "small", "vcpus": 2, "ram": 4096},
					{"id": "mt2", "name": "medium", "vcpus": 4, "ram": 8192}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_sort") != "vcpus:asc" {
					t.Errorf("expected sort=vcpus:asc, got %s", r.URL.Query().Get("_sort"))
				}
			},
		},
		{
			name: "with availability zone",
			opts: InstanceTypeListOptions{
				AvailabilityZone: "zone1",
			},
			response: `{
				"instance_types": [
					{"id": "mt1", "name": "small", "vcpus": 2, "ram": 4096, "availability_zones": ["zone1"]}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("availability-zone") != "zone1" {
					t.Errorf("expected availability-zone=zone1, got %s", r.URL.Query().Get("availability-zone"))
				}
			},
		},
		{
			name:       "server error",
			opts:       InstanceTypeListOptions{},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "empty response",
			opts:       InstanceTypeListOptions{},
			response:   "",
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			opts:       InstanceTypeListOptions{},
			response:   `{"instance_types": [{"id": "broken"}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid pagination values",
			opts: InstanceTypeListOptions{
				Limit:  intPtr(-1),
				Offset: intPtr(-1),
			},
			response:   `{"error": "invalid pagination parameters"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_limit") != "-1" {
					t.Errorf("expected limit=-1, got %s", r.URL.Query().Get("_limit"))
				}
				if r.URL.Query().Get("_offset") != "-1" {
					t.Errorf("expected offset=-1, got %s", r.URL.Query().Get("_offset"))
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
			got, err := client.InstanceTypes().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("List() got %v instance types, want %v", len(got), tt.want)
			}
		})
	}
}
