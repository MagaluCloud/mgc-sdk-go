package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestParameterService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		opts       ListParametersOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "basic list",
			opts: ListParametersOptions{
				ParameterGroupID: "g1",
			},
			response: `{
                "meta": {
                  "offset": 0,
                  "limit": 10,
                  "count": 2,
                  "total": 2,
                  "max_limit": 100
                },
                "results": [
                  {"id": "p1", "name": "param1", "value": 1},
                  {"id": "p2", "name": "param2", "value": "v2"}
                ]
            }`,
			statusCode: http.StatusOK,
			want:       2,
		},
		{
			name: "with pagination",
			opts: ListParametersOptions{
				ParameterGroupID: "g1",
				Offset:           helpers.IntPtr(1),
				Limit:            helpers.IntPtr(1),
			},
			response: `{
                "meta": {
                  "offset": 1,
                  "limit": 1,
                  "count": 1,
                  "total": 2,
                  "max_limit": 100
                },
                "results": [
                  {"id": "p2", "name": "param2", "value": "v2"}
                ]
            }`,
			statusCode: http.StatusOK,
			want:       1,
		},
		{
			name:       "server error",
			opts:       ListParametersOptions{ParameterGroupID: "g1"},
			response:   `{"error":"internal"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "nil response body",
			opts:       ListParametersOptions{ParameterGroupID: "g1"},
			response:   ``,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			opts:       ListParametersOptions{ParameterGroupID: "g1"},
			response:   `{"meta":,`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientParameters(server.URL)
			got, err := client.List(context.Background(), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got.Results) != tt.want {
				t.Errorf("List() got %d, want %d", len(got.Results), tt.want)
			}
		})
	}
}

func TestParameterService_List_PaginationMetadata(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"meta": {
				"page": {
					"offset": 10,
					"limit": 5,
					"count": 5,
					"total": 20,
					"max_limit": 100
				},
				"filters": [
					{"field": "parameter_group_id", "value": "group-123"}
				]
			},
			"results": [
				{"id": "p1", "name": "param1", "value": 1},
				{"id": "p2", "name": "param2", "value": 2}
			]
		}`))
	}))
	defer server.Close()

	client := testClientParameters(server.URL)
	offset := 10
	limit := 5
	result, err := client.List(context.Background(), ListParametersOptions{
		ParameterGroupID: "group-123",
		Offset:           &offset,
		Limit:            &limit,
	})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// Validate results
	if len(result.Results) != 2 {
		t.Errorf("List() got %d results, want 2", len(result.Results))
	}

	// Validate pagination metadata
	if result.Meta.Page.Offset != 10 {
		t.Errorf("Meta.Page.Offset = %d, want 10", result.Meta.Page.Offset)
	}
	if result.Meta.Page.Limit != 5 {
		t.Errorf("Meta.Page.Limit = %d, want 5", result.Meta.Page.Limit)
	}
	if result.Meta.Page.Count != 5 {
		t.Errorf("Meta.Page.Count = %d, want 5", result.Meta.Page.Count)
	}
	if result.Meta.Page.Total != 20 {
		t.Errorf("Meta.Page.Total = %d, want 20", result.Meta.Page.Total)
	}
	if result.Meta.Page.MaxLimit != 100 {
		t.Errorf("Meta.Page.MaxLimit = %d, want 100", result.Meta.Page.MaxLimit)
	}

	// Validate filters metadata
	if len(result.Meta.Filters) != 1 {
		t.Errorf("Meta.Filters length = %d, want 1", len(result.Meta.Filters))
	}
}

func TestParameterService_ListAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		filterOpts ParameterFilterOptions
		pages      []string
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page",
			filterOpts: ParameterFilterOptions{
				ParameterGroupID: "group-1",
			},
			pages: []string{
				`{
					"meta": {"page": {"offset": 0, "limit": 25, "count": 2, "total": 2}},
					"results": [
						{"id": "p1", "name": "param1", "value": 1},
						{"id": "p2", "name": "param2", "value": 2}
					]
				}`,
			},
			wantCount: 2,
		},
		{
			name: "multiple pages",
			filterOpts: ParameterFilterOptions{
				ParameterGroupID: "group-2",
			},
			pages: []string{
				func() string {
					results := `[`
					for i := 0; i < 50; i++ {
						if i > 0 {
							results += ","
						}
						results += fmt.Sprintf(`{"id": "p%d", "name": "param%d", "value": %d}`, i+1, i+1, i+1)
					}
					results += `]`
					return fmt.Sprintf(`{
						"meta": {"page": {"offset": 0, "limit": 25, "count": 50, "total": 75}},
						"results": %s
					}`, results)
				}(),
				func() string {
					results := `[`
					for i := 0; i < 25; i++ {
						if i > 0 {
							results += ","
						}
						results += fmt.Sprintf(`{"id": "p%d", "name": "param%d", "value": %d}`, i+51, i+51, i+51)
					}
					results += `]`
					return fmt.Sprintf(`{
						"meta": {"page": {"offset": 50, "limit": 25, "count": 25, "total": 75}},
						"results": %s
					}`, results)
				}(),
			},
			wantCount: 75,
		},
		{
			name: "empty result",
			filterOpts: ParameterFilterOptions{
				ParameterGroupID: "group-empty",
			},
			pages: []string{
				`{
					"meta": {"page": {"offset": 0, "limit": 25, "count": 0, "total": 0}},
					"results": []
				}`,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if requestCount < len(tt.pages) {
					w.Write([]byte(tt.pages[requestCount]))
					requestCount++
				} else {
					w.Write([]byte(`{"meta": {"page": {"offset": 0, "limit": 25, "count": 0, "total": 0}}, "results": []}`))
				}
			}))
			defer server.Close()

			client := testClientParameters(server.URL)
			got, err := client.ListAll(context.Background(), tt.filterOpts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("ListAll() got %d parameters, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestParameterService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		groupID    string
		req        ParameterCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "successful create",
			groupID:    "g1",
			req:        ParameterCreateRequest{Name: "p1", Value: 42},
			response:   `{"id":"p1"}`,
			statusCode: http.StatusOK,
			wantID:     "p1",
		},
		{
			name:       "server error",
			groupID:    "g1",
			req:        ParameterCreateRequest{Name: "p1", Value: 42},
			response:   `{"error":"fail"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			groupID:    "g1",
			req:        ParameterCreateRequest{Name: "p1", Value: 42},
			response:   `{"id":`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientParameters(server.URL)
			got, err := client.Create(context.Background(), tt.groupID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.ID != tt.wantID {
				t.Errorf("Create() got ID = %q, want %q", got.ID, tt.wantID)
			}
		})
	}
}

func TestParameterService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		groupID    string
		paramID    string
		req        ParameterUpdateRequest
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful update",
			groupID:    "g1",
			paramID:    "p1",
			req:        ParameterUpdateRequest{Value: "new"},
			response:   `{"id":"p1","name":"param1","value":"new"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "not found",
			groupID:    "g1",
			paramID:    "x",
			req:        ParameterUpdateRequest{Value: 0},
			response:   `{"error":"not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			groupID:    "g1",
			paramID:    "p1",
			req:        ParameterUpdateRequest{Value: 0},
			response:   `{"id":`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientParameters(server.URL)
			_, err := client.Update(context.Background(), tt.groupID, tt.paramID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParameterService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		groupID    string
		paramID    string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			groupID:    "g1",
			paramID:    "p1",
			statusCode: http.StatusNoContent,
		},
		{
			name:       "not found",
			groupID:    "g1",
			paramID:    "x",
			statusCode: http.StatusNotFound,
			response:   `{"error":"not found"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			groupID:    "g1",
			paramID:    "p1",
			statusCode: http.StatusInternalServerError,
			response:   `{"error":"fail"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClientParameters(server.URL)
			err := client.Delete(context.Background(), tt.groupID, tt.paramID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func testClientParameters(baseURL string) ParameterService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return &parameterService{New(core)}
}
