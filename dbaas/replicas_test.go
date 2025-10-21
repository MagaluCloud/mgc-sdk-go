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

// Helper functions
func assertEqual(t *testing.T, expected, actual any, msgAndArgs ...any) {
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

func testClient(baseURL string) ReplicaService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return NewReplicaService(New(core))
}

func TestReplicaService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListReplicaOptions
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
					{"id": "rep1", "name": "replica1"},
					{"id": "rep2", "name": "replica2"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: ListReplicaOptions{
				Limit:  helpers.IntPtr(1),
				Offset: helpers.IntPtr(1),
			},
			response: `{
				"meta": {"total": 1},
				"results": [{"id": "rep2", "name": "replica2"}]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
			wantErr:    false,
		},
		{
			name: "filter by source",
			opts: ListReplicaOptions{
				SourceID: helpers.StrPtr("src1"),
			},
			response: `{
				"meta": {"total": 1},
				"results": [{"id": "rep1", "source_id": "src1"}]
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
				assertEqual(t, "/database/v2/replicas", r.URL.Path)

				query := r.URL.Query()
				if tt.opts.Limit != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Limit), query.Get("_limit"))
				}
				if tt.opts.Offset != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Offset), query.Get("_offset"))
				}
				if tt.opts.SourceID != nil {
					assertEqual(t, *tt.opts.SourceID, query.Get("source_id"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(result.Results))
			// Verify metadata is returned
			if result != nil && tt.wantCount > 0 {
				assertEqual(t, true, result.Meta.Page.Total >= 0)
			}
		})
	}
}

func TestReplicaService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "existing replica",
			id:   "rep1",
			response: `{
				"id": "rep1",
				"name": "test-replica",
				"status": "ACTIVE"
			}`,
			statusCode: http.StatusOK,
			wantID:     "rep1",
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
				assertEqual(t, fmt.Sprintf("/database/v2/replicas/%s", tt.id), r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			replica, err := client.Get(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, replica.ID)
		})
	}
}

func TestReplicaService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    ReplicaCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: ReplicaCreateRequest{
				SourceID: "src1",
				Name:     "test-replica",
			},
			response:   `{"id": "rep1"}`,
			statusCode: http.StatusOK,
			wantID:     "rep1",
			wantErr:    false,
		},
		{
			name: "invalid request",
			request: ReplicaCreateRequest{
				Name: "missing-source",
			},
			response:   `{"error": "source_id required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/database/v2/replicas", r.URL.Path)

				var req ReplicaCreateRequest
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, tt.request.SourceID, req.SourceID)
				assertEqual(t, tt.request.Name, req.Name)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
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

func TestReplicaService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "rep1",
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
				assertEqual(t, fmt.Sprintf("/database/v2/replicas/%s", tt.id), r.URL.Path)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
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

func TestReplicaService_Resize(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    ReplicaResizeRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "resize instance type",
			id:   "rep1",
			request: ReplicaResizeRequest{
				InstanceTypeID: helpers.StrPtr("type-large"),
				Volume: &InstanceVolumeResizeRequest{
					Size: 200,
					Type: "nvme",
				},
			},
			response: `{
				"id": "rep1",
				"instance_type_id": "type-large",
				"volume": {"size": 200}
			}`,
			statusCode: http.StatusOK,
			wantID:     "rep1",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v2/replicas/%s/resize", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req ReplicaResizeRequest
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, *tt.request.InstanceTypeID, *req.InstanceTypeID)
				assertEqual(t, tt.request.Volume.Size, req.Volume.Size)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
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

func TestReplicaService_StartStop(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "successful start",
			method:     "Start",
			id:         "rep1",
			response:   `{"id": "rep1", "status": "STARTING"}`,
			statusCode: http.StatusOK,
			wantID:     "rep1",
			wantErr:    false,
		},
		{
			name:       "already running",
			method:     "Start",
			id:         "rep-running",
			response:   `{"error": "already running"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name:       "successful stop",
			method:     "Stop",
			id:         "rep1",
			response:   `{"id": "rep1", "status": "STOPPING"}`,
			statusCode: http.StatusOK,
			wantID:     "rep1",
			wantErr:    false,
		},
		{
			name:       "already stopped",
			method:     "Stop",
			id:         "rep-stopped",
			response:   `{"error": "already stopped"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var expectedPath string
				if tt.method == "Start" {
					expectedPath = fmt.Sprintf("/database/v2/replicas/%s/start", tt.id)
				} else {
					expectedPath = fmt.Sprintf("/database/v2/replicas/%s/stop", tt.id)
				}
				assertEqual(t, expectedPath, r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			var result *ReplicaDetailResponse
			var err error

			if tt.method == "Start" {
				result, err = client.Start(context.Background(), tt.id)
			} else {
				result, err = client.Stop(context.Background(), tt.id)
			}

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

func TestReplicaService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		opts       ReplicaFilterOptions
		pages      []string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page",
			pages: []string{
				`{
					"meta": {
						"page": {"count": 2, "limit": 25, "offset": 0, "total": 2}
					},
					"results": [
						{"id": "rep1", "name": "replica1"},
						{"id": "rep2", "name": "replica2"}
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
				func() string {
					results := `[`
					for i := 0; i < 25; i++ {
						if i > 0 {
							results += ","
						}
						results += fmt.Sprintf(`{"id": "rep%d", "name": "replica%d"}`, i+1, i+1)
					}
					results += `]`
					return fmt.Sprintf(`{
						"meta": {"page": {"offset": 0, "limit": 25, "count": 25, "total": 60}},
						"results": %s
					}`, results)
				}(),
				func() string {
					results := `[`
					for i := 0; i < 25; i++ {
						if i > 0 {
							results += ","
						}
						results += fmt.Sprintf(`{"id": "rep%d", "name": "replica%d"}`, i+26, i+26)
					}
					results += `]`
					return fmt.Sprintf(`{
						"meta": {"page": {"offset": 25, "limit": 25, "count": 25, "total": 60}},
						"results": %s
					}`, results)
				}(),
				func() string {
					results := `[`
					for i := 0; i < 10; i++ {
						if i > 0 {
							results += ","
						}
						results += fmt.Sprintf(`{"id": "rep%d", "name": "replica%d"}`, i+51, i+51)
					}
					results += `]`
					return fmt.Sprintf(`{
						"meta": {"page": {"offset": 50, "limit": 25, "count": 10, "total": 60}},
						"results": %s
					}`, results)
				}(),
			},
			statusCode: http.StatusOK,
			wantCount:  60,
			wantErr:    false,
		},
		{
			name: "empty results",
			pages: []string{
				`{
					"meta": {
						"page": {"count": 0, "limit": 25, "offset": 0, "total": 0}
					},
					"results": []
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
		{
			name: "with source filter",
			opts: ReplicaFilterOptions{
				SourceID: helpers.StrPtr("src1"),
			},
			pages: []string{
				`{
					"meta": {
						"filters": [{"field": "source_id", "value": "src1"}],
						"page": {"count": 1, "limit": 25, "offset": 0, "total": 1}
					},
					"results": [
						{"id": "rep1", "source_id": "src1", "name": "replica1"}
					]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "server error",
			pages:      []string{`{"error": "internal server error"}`},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/database/v2/replicas", r.URL.Path)

				query := r.URL.Query()
				if tt.opts.SourceID != nil {
					assertEqual(t, *tt.opts.SourceID, query.Get("source_id"))
				}

				// Determine which page to return based on offset
				offset := query.Get("_offset")
				currentPage := 0
				if offset != "" {
					offsetInt, _ := strconv.Atoi(offset)
					currentPage = offsetInt / 25
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if currentPage < len(tt.pages) {
					w.Write([]byte(tt.pages[currentPage]))
				} else {
					// Return empty results if we've run out of pages
					w.Write([]byte(`{"meta": {"page": {"count": 0, "limit": 25, "offset": 0, "total": 0}}, "results": []}`))
				}
			}))
			defer server.Close()

			client := testClient(server.URL)
			results, err := client.ListAll(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(results))

		})
	}
}
