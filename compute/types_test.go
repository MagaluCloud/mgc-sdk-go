package compute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestMachineTypeService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       InstanceTypeListOptions
		response   string
		statusCode int
		wantCount  int
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
			wantCount:  2,
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
			wantCount:  1,
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
			wantCount:  2,
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
			wantCount:  1,
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
			if !tt.wantErr {
				if len(got.InstanceTypes) != tt.wantCount {
					t.Errorf("List() got %v instance types, want %v", len(got.InstanceTypes), tt.wantCount)
				}
				if got.Meta.Page.Count != tt.wantCount {
					t.Errorf("List() meta count = %v, want %v", got.Meta.Page.Count, tt.wantCount)
				}
			}
		})
	}
}

func TestInstanceTypeService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		opts       InstanceTypeFilterOptions
		responses  []string
		want       int
		wantErr    bool
		checkCalls func(*testing.T, int)
	}{
		{
			name: "single page",
			opts: InstanceTypeFilterOptions{},
			responses: []string{
				`{
					"instance_types": [
						{"id": "mt1", "name": "small", "vcpus": 2, "ram": 4096, "disk": 50},
						{"id": "mt2", "name": "medium", "vcpus": 4, "ram": 8192, "disk": 100}
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
			opts: InstanceTypeFilterOptions{},
			responses: []string{
				`{
					"instance_types": [` +
					generateInstanceTypeJSON(50, 0) + `
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
					"instance_types": [` +
					generateInstanceTypeJSON(25, 50) + `
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
			name: "with availability zone filter",
			opts: InstanceTypeFilterOptions{
				AvailabilityZone: "zone1",
			},
			responses: []string{
				`{
					"instance_types": [
						{"id": "mt1", "name": "small", "vcpus": 2, "ram": 4096, "availability_zones": ["zone1"]}
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
		{
			name: "with sorting",
			opts: InstanceTypeFilterOptions{
				Sort: strPtr("vcpus:asc"),
			},
			responses: []string{
				`{
					"instance_types": [
						{"id": "mt1", "name": "small", "vcpus": 2, "ram": 4096},
						{"id": "mt2", "name": "medium", "vcpus": 4, "ram": 8192}
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
			name: "empty results",
			opts: InstanceTypeFilterOptions{},
			responses: []string{
				`{
					"instance_types": [],
					"meta": {
						"page": {
							"offset": 0,
							"limit": 50,
							"count": 0,
							"total": 0
						}
					}
				}`,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "server error",
			opts: InstanceTypeFilterOptions{},
			responses: []string{
				`{"error": "internal server error"}`,
				`{"error": "internal server error"}`,
				`{"error": "internal server error"}`,
			},
			wantErr: true,
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
				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				w.Write([]byte(tt.responses[callCount]))
				callCount++
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.InstanceTypes().ListAll(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ListAll() got %v instance types, want %v", len(got), tt.want)
			}
			if tt.checkCalls != nil {
				tt.checkCalls(t, callCount)
			}
		})
	}
}

// Helper function to generate instance type JSON for testing
func generateInstanceTypeJSON(count, startID int) string {
	if count == 0 {
		return ""
	}
	var result string
	for i := 0; i < count; i++ {
		if i > 0 {
			result += ","
		}
		id := startID + i + 1
		result += `{"id": "mt` + strconv.Itoa(id) + `", "name": "type` + strconv.Itoa(id) + `", "vcpus": 2, "ram": 4096, "disk": 50}`
	}
	return result
}
