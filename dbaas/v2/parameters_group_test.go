package dbaas

import (
	"context"
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
			if len(got) != tt.want {
				t.Errorf("List() got %v parameter groups, want %v", len(got), tt.want)
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
