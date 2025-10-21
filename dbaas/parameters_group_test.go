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

func TestParameterGroupService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		opts       ListParameterGroupsOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "basic list",
			opts: ListParameterGroupsOptions{},
			response: `{
				"meta": {"total": 2},
				"results": [
					{"id": "pg1", "name": "param-group-1", "description": "test group 1", "type": "USER", "engine_id": "eng1"},
					{"id": "pg2", "name": "param-group-2", "description": "test group 2", "type": "SYSTEM", "engine_id": "eng1"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: ListParameterGroupsOptions{
				Limit:  helpers.IntPtr(1),
				Offset: helpers.IntPtr(1),
			},
			response: `{
				"meta": {"total": 2},
				"results": [
					{"id": "pg2", "name": "param-group-2", "description": "test group 2", "type": "SYSTEM", "engine_id": "eng1"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name: "filter by type",
			opts: ListParameterGroupsOptions{
				Type: paramGroupTypePtr(ParameterGroupTypeUser),
			},
			response: `{
				"meta": {"total": 1},
				"results": [
					{"id": "pg1", "name": "param-group-1", "description": "test group 1", "type": "USER", "engine_id": "eng1"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name: "filter by engine id",
			opts: ListParameterGroupsOptions{
				EngineID: helpers.StrPtr("eng2"),
			},
			response: `{
				"meta": {"total": 1},
				"results": [
					{"id": "pg3", "name": "param-group-3", "description": "test group 3", "type": "SYSTEM", "engine_id": "eng2"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name:       "server error",
			opts:       ListParameterGroupsOptions{},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "invalid pagination",
			opts: ListParameterGroupsOptions{
				Limit: helpers.IntPtr(-1),
			},
			response:   `{"error": "invalid limit"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "nil response body",
			opts:       ListParameterGroupsOptions{},
			response:   "",
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			opts:       ListParameterGroupsOptions{},
			response:   `{"meta": {"total": 1}, "results": [{"id": "broken"}`,
			statusCode: http.StatusOK,
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

			client := testClientParamerts(server.URL)
			got, err := client.List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got.Results) != tt.want {
				t.Errorf("List() got %v parameter groups, want %v", len(got.Results), tt.want)
			}
		})
	}
}

func TestParameterGroupService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		req        ParameterGroupCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			req: ParameterGroupCreateRequest{
				Name:        "test-param-group",
				EngineID:    "eng1",
				Description: helpers.StrPtr("Test parameter group"),
			},
			response:   `{"id": "pg1"}`,
			statusCode: http.StatusOK,
			wantID:     "pg1",
			wantErr:    false,
		},
		{
			name: "empty name",
			req: ParameterGroupCreateRequest{
				Name:     "",
				EngineID: "eng1",
			},
			response:   `{"error": "name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid engine id",
			req: ParameterGroupCreateRequest{
				Name:     "test-param-group",
				EngineID: "invalid-engine",
			},
			response:   `{"error": "invalid engine id"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "server error",
			req: ParameterGroupCreateRequest{
				Name:     "test-param-group",
				EngineID: "eng1",
			},
			response:   `{"error": "internal error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "duplicate name",
			req: ParameterGroupCreateRequest{
				Name:     "existing-group",
				EngineID: "eng1",
			},
			response:   `{"error": "parameter group name already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "malformed json response",
			req: ParameterGroupCreateRequest{
				Name:     "test-param-group",
				EngineID: "eng1",
			},
			response:   `{"id": "pg1"`,
			statusCode: http.StatusOK,
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

			client := testClientParamerts(server.URL)
			got, err := client.Create(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.ID != tt.wantID {
				t.Errorf("Create() got = %v, want %v", got.ID, tt.wantID)
			}
		})
	}
}

func TestParameterGroupService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "existing parameter group",
			id:   "pg1",
			response: `{
				"id": "pg1",
				"name": "param-group-1",
				"description": "test group 1",
				"type": "USER",
				"engine_id": "eng1"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existing parameter group",
			id:         "invalid",
			response:   `{"error": "Parameter group not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "empty id",
			id:         "",
			response:   `{"error": "ID cannot be empty"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "server error",
			id:         "error",
			response:   `{"error": "Internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "malformed response",
			id:         "pg1",
			response:   `{"id": "pg1", "name":}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.id == "" {
					// Test error handling in the client before request is sent
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientParamerts(server.URL)
			got, err := client.Get(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("Get() expected error, got nil")
					return
				}
			} else {
				if err != nil {
					t.Errorf("Get() unexpected error: %v", err)
					return
				}
				if got == nil {
					t.Error("Get() got nil, want parameter group")
					return
				}
				if got.ID != tt.id {
					t.Errorf("Get() got ID = %v, want %v", got.ID, tt.id)
				}
			}
		})
	}
}

func TestParameterGroupService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		id         string
		req        ParameterGroupUpdateRequest
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful update",
			id:   "pg1",
			req: ParameterGroupUpdateRequest{
				Name:        helpers.StrPtr("updated-name"),
				Description: helpers.StrPtr("Updated description"),
			},
			response: `{
				"id": "pg1",
				"name": "updated-name",
				"description": "Updated description",
				"type": "USER",
				"engine_id": "eng1"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "update name only",
			id:   "pg1",
			req: ParameterGroupUpdateRequest{
				Name: helpers.StrPtr("updated-name"),
			},
			response: `{
				"id": "pg1",
				"name": "updated-name",
				"description": "Original description",
				"type": "USER",
				"engine_id": "eng1"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "update description only",
			id:   "pg1",
			req: ParameterGroupUpdateRequest{
				Description: helpers.StrPtr("Updated description"),
			},
			response: `{
				"id": "pg1",
				"name": "original-name",
				"description": "Updated description",
				"type": "USER",
				"engine_id": "eng1"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty id",
			id:         "",
			req:        ParameterGroupUpdateRequest{Name: helpers.StrPtr("updated-name")},
			response:   `{"error": "ID cannot be empty"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "system parameter group",
			id:   "sys1",
			req: ParameterGroupUpdateRequest{
				Name: helpers.StrPtr("updated-name"),
			},
			response:   `{"error": "cannot update system parameter group"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "non-existing parameter group",
			id:   "invalid",
			req: ParameterGroupUpdateRequest{
				Name: helpers.StrPtr("updated-name"),
			},
			response:   `{"error": "parameter group not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "duplicate name",
			id:   "pg1",
			req: ParameterGroupUpdateRequest{
				Name: helpers.StrPtr("existing-name"),
			},
			response:   `{"error": "parameter group name already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "server error",
			id:   "pg1",
			req: ParameterGroupUpdateRequest{
				Name: helpers.StrPtr("updated-name"),
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.id == "" {
					// Test error handling in the client before request is sent
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientParamerts(server.URL)
			got, err := client.Update(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got == nil {
					t.Error("Update() got nil, want parameter group")
					return
				}
				if got.ID != tt.id {
					t.Errorf("Update() got ID = %v, want %v", got.ID, tt.id)
				}
			}
		})
	}
}

func TestParameterGroupService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "pg1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "empty id",
			id:         "",
			response:   `{"error": "ID cannot be empty"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "non-existing parameter group",
			id:         "invalid",
			response:   `{"error": "parameter group not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "system parameter group",
			id:         "sys1",
			response:   `{"error": "cannot delete system parameter group"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "parameter group in use",
			id:         "in-use",
			response:   `{"error": "parameter group is in use"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name:       "server error",
			id:         "error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.id == "" {
					// Test error handling in the client before request is sent
					return
				}

				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testClientParamerts(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParameterGroupService_Concurrent(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"meta": {"total": 0}, "results": []}`))
	}))
	defer server.Close()

	client := testClientParamerts(server.URL)
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool)
	for range 10 {
		go func() {
			_, err := client.List(ctx, ListParameterGroupsOptions{})
			if err != nil {
				t.Errorf("concurrent List() error = %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for range 10 {
		<-done
	}
}

func paramGroupTypePtr(t ParameterGroupType) *ParameterGroupType {
	return &t
}

func testClientParamerts(baseURL string) ParameterGroupService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return &parameterGroupService{New(core)}
}

func TestParameterGroupService_List_PaginationMetadata(t *testing.T) {
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
					{"field": "type", "value": "USER"}
				]
			},
			"results": [
				{"id": "pg1", "name": "group1", "type": "USER", "engine_id": "eng1"},
				{"id": "pg2", "name": "group2", "type": "USER", "engine_id": "eng1"}
			]
		}`))
	}))
	defer server.Close()

	client := testClientParamerts(server.URL)
	offset := 10
	limit := 5
	pgType := ParameterGroupTypeUser
	result, err := client.List(context.Background(), ListParameterGroupsOptions{
		Offset: &offset,
		Limit:  &limit,
		Type:   &pgType,
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

func TestParameterGroupService_ListAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		filterOpts ParameterGroupFilterOptions
		pages      []string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "single page",
			filterOpts: ParameterGroupFilterOptions{},
			pages: []string{
				`{
					"meta": {"page": {"offset": 0, "limit": 25, "count": 3, "total": 3}},
					"results": [
						{"id": "pg1", "name": "group1", "type": "USER", "engine_id": "eng1"},
						{"id": "pg2", "name": "group2", "type": "SYSTEM", "engine_id": "eng1"},
						{"id": "pg3", "name": "group3", "type": "USER", "engine_id": "eng2"}
					]
				}`,
			},
			wantCount: 3,
		},
		{
			name: "multiple pages",
			filterOpts: ParameterGroupFilterOptions{
				Type: paramGroupTypePtr(ParameterGroupTypeUser),
			},
			pages: []string{
				func() string {
					results := `[`
					for i := 0; i < 50; i++ {
						if i > 0 {
							results += ","
						}
						results += `{"id": "pg` + fmt.Sprintf("%d", i) + `", "name": "group` + fmt.Sprintf("%d", i) + `", "type": "USER", "engine_id": "eng1"}`
					}
					results += `]`
					return `{
						"meta": {"page": {"offset": 0, "limit": 25, "count": 50, "total": 60}},
						"results": ` + results + `
					}`
				}(),
				`{
					"meta": {"page": {"offset": 50, "limit": 25, "count": 10, "total": 60}},
					"results": [
						{"id": "pg50", "name": "group50", "type": "USER", "engine_id": "eng1"},
						{"id": "pg51", "name": "group51", "type": "USER", "engine_id": "eng1"},
						{"id": "pg52", "name": "group52", "type": "USER", "engine_id": "eng1"},
						{"id": "pg53", "name": "group53", "type": "USER", "engine_id": "eng1"},
						{"id": "pg54", "name": "group54", "type": "USER", "engine_id": "eng1"},
						{"id": "pg55", "name": "group55", "type": "USER", "engine_id": "eng1"},
						{"id": "pg56", "name": "group56", "type": "USER", "engine_id": "eng1"},
						{"id": "pg57", "name": "group57", "type": "USER", "engine_id": "eng1"},
						{"id": "pg58", "name": "group58", "type": "USER", "engine_id": "eng1"},
						{"id": "pg59", "name": "group59", "type": "USER", "engine_id": "eng1"}
					]
				}`,
			},
			wantCount: 60,
		},
		{
			name: "with engine filter",
			filterOpts: ParameterGroupFilterOptions{
				EngineID: helpers.StrPtr("eng2"),
			},
			pages: []string{
				`{
					"meta": {"page": {"offset": 0, "limit": 25, "count": 1, "total": 1}},
					"results": [
						{"id": "pg1", "name": "group1", "type": "USER", "engine_id": "eng2"}
					]
				}`,
			},
			wantCount: 1,
		},
		{
			name:       "empty result",
			filterOpts: ParameterGroupFilterOptions{},
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
		t.Run(tt.name, func(t *testing.T) {
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

			client := testClientParamerts(server.URL)
			got, err := client.ListAll(context.Background(), tt.filterOpts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("ListAll() got %d parameter groups, want %d", len(got), tt.wantCount)
			}
		})
	}
}
