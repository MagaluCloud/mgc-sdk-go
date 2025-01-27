package compute

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestInstanceService_List(t *testing.T) {
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
				"instances": [
					{"id": "inst1", "name": "test1"},
					{"id": "inst2", "name": "test2"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: ListOptions{
				Limit:  intPtr(1),
				Offset: intPtr(1),
			},
			response: `{
				"instances": [
					{"id": "inst2", "name": "test2"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name:       "server error",
			opts:       ListOptions{},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "invalid response format",
			opts: ListOptions{},
			response: `{
				"invalid": []
			}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid pagination",
			opts: ListOptions{
				Limit: intPtr(-1),
			},
			response:   `{"error": "invalid limit"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "nil response body",
			opts:       ListOptions{},
			response:   "",
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			opts:       ListOptions{},
			response:   `{"instances": [{"id": "broken"}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid sort parameter",
			opts: ListOptions{
				Sort: strPtr("invalid:order"),
			},
			response:   `{"error": "invalid sort parameter"}`,
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
			got, err := client.Instances().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("List() got %v instances, want %v", len(got), tt.want)
			}
		})
	}
}

func TestInstanceService_Create(t *testing.T) {
	tests := []struct {
		name       string
		req        CreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			req: CreateRequest{
				Name: "test-vm",
			},
			response:   `{"id": "inst1"}`,
			statusCode: http.StatusOK,
			wantID:     "inst1",
			wantErr:    false,
		},
		{
			name:       "empty name",
			req:        CreateRequest{},
			response:   `{"error": "name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid machine type",
			req: CreateRequest{
				Name: "test-vm",
				MachineType: IDOrName{
					Name: strPtr("invalid-type"),
				},
			},
			response:   `{"error": "invalid machine type"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "server error",
			req:        CreateRequest{Name: "test-vm"},
			response:   `{"error": "internal error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "duplicate name",
			req: CreateRequest{
				Name: "existing-vm",
			},
			response:   `{"error": "instance name already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "insufficient resources",
			req: CreateRequest{
				Name:        "test-vm",
				MachineType: IDOrName{Name: strPtr("large")},
			},
			response:   `{"error": "insufficient resources"}`,
			statusCode: http.StatusServiceUnavailable,
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
			gotID, err := client.Instances().Create(context.Background(), tt.req)

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

func TestInstanceService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expand     []string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "existing instance",
			id:   "inst1",
			response: `{
				"id": "inst1",
				"name": "test-vm"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existing instance",
			id:         "invalid",
			response:   `{"error": "Instance not found"}`,
			statusCode: http.StatusNotFound,
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
			name:   "with expansion",
			id:     "inst1",
			expand: []string{"network", "storage"},
			response: `{
				"id": "inst1",
				"name": "test-vm",
				"network": {"id": "net1"},
				"storage": {"id": "stor1"}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "malformed response",
			id:         "inst1",
			response:   `{"id": "inst1", "name":}`,
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

			client := testClient(server.URL)
			got, err := client.Instances().Get(context.Background(), tt.id, tt.expand)

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
					t.Error("Get() got nil, want instance")
					return
				}
				if got.ID != tt.id {
					t.Errorf("Get() got ID = %v, want %v", got.ID, tt.id)
				}
			}
		})
	}
}

func TestInstanceService_Delete(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		deletePublicIP bool
		statusCode     int
		response       string
		wantErr        bool
	}{
		{
			name:           "successful delete",
			id:             "inst1",
			deletePublicIP: true,
			statusCode:     http.StatusNoContent,
			wantErr:        false,
		},
		{
			name:           "not found",
			id:             "invalid",
			deletePublicIP: false,
			statusCode:     http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:       "empty id",
			id:         "",
			response:   `{"error": "id is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "instance in use",
			id:         "in-use",
			response:   `{"error": "instance is in use"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("delete_public_ip") != strconv.FormatBool(tt.deletePublicIP) {
					t.Errorf("unexpected delete_public_ip query param: got %v", r.URL.Query().Get("delete_public_ip"))
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Instances().Delete(context.Background(), tt.id, tt.deletePublicIP)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceService_Rename(t *testing.T) {
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
			id:         "inst1",
			newName:    "new-name",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty id",
			id:         "",
			newName:    "new-name",
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "empty new name",
			id:         "inst1",
			newName:    "",
			response:   `{"error": "name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "name in use",
			id:         "inst1",
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
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Instances().Rename(context.Background(), tt.id, tt.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceService_Retype(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		req        RetypeRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "successful retype",
			id:   "inst1",
			req: RetypeRequest{
				MachineType: IDOrName{
					Name: strPtr("new-type"),
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty id",
			id:         "",
			req:        RetypeRequest{MachineType: IDOrName{Name: strPtr("new-type")}},
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "instance running",
			id:   "running",
			req: RetypeRequest{
				MachineType: IDOrName{Name: strPtr("new-type")},
			},
			response:   `{"error": "instance must be stopped"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "invalid machine type",
			id:   "inst1",
			req: RetypeRequest{
				MachineType: IDOrName{Name: strPtr("")},
			},
			response:   `{"error": "invalid machine type"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Instances().Retype(context.Background(), tt.id, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Retype() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceService_StateOperations(t *testing.T) {
	tests := []struct {
		name       string
		operation  string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "start success",
			operation:  "start",
			id:         "inst1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "stop success",
			operation:  "stop",
			id:         "inst1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "suspend success",
			operation:  "suspend",
			id:         "inst1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "invalid operation",
			operation:  "start",
			id:         "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "already started",
			operation:  "start",
			id:         "running",
			response:   `{"error": "instance is already running"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name:       "already stopped",
			operation:  "stop",
			id:         "stopped",
			response:   `{"error": "instance is already stopped"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name:       "suspend failed",
			operation:  "suspend",
			id:         "inst1",
			response:   `{"error": "suspend operation failed"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "operation timeout",
			operation:  "start",
			id:         "timeout",
			response:   `{"error": "operation timed out"}`,
			statusCode: http.StatusGatewayTimeout,
			wantErr:    true,
		},
		{
			name:       "operation in progress",
			operation:  "stop",
			id:         "busy",
			response:   `{"error": "another operation in progress"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	operations := map[string]func(*VirtualMachineClient, context.Context, string) error{
		"start": func(c *VirtualMachineClient, ctx context.Context, id string) error {
			return c.Instances().Start(ctx, id)
		},
		"stop": func(c *VirtualMachineClient, ctx context.Context, id string) error {
			return c.Instances().Stop(ctx, id)
		},
		"suspend": func(c *VirtualMachineClient, ctx context.Context, id string) error {
			return c.Instances().Suspend(ctx, id)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			operation := operations[tt.operation]
			err := operation(client, context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s() error = %v, wantErr %v", tt.operation, err, tt.wantErr)
			}
		})
	}
}

func TestInstanceService_Concurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"instances": []}`))
	}))
	defer server.Close()

	client := testClient(server.URL)
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.Instances().List(ctx, ListOptions{})
			if err != nil {
				t.Errorf("concurrent List() error = %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Helper functions
func testClient(baseURL string) *VirtualMachineClient {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core)
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

// here
func TestInstanceService_ListWithExpand(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListOptions
		response   string
		statusCode int
		wantErr    bool
		checkQuery func(*testing.T, *http.Request)
	}{
		{
			name: "single expand",
			opts: ListOptions{
				Expand: []string{"network"},
			},
			response: `{
				"instances": [{
					"id": "inst1",
					"name": "test1",
					"network": {
						"id": "net1",
						"name": "network1"
					}
				}]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if !r.URL.Query().Has("expand") {
					t.Error("expected expand parameter, got none")
					return
				}
				got := r.URL.Query().Get("expand")
				switch got {
				case "network":
					// single expand ok
				case "network,storage,machineType":
					// multiple expand ok
				case "invalid":
					// invalid expand ok
				default:
					t.Errorf("unexpected expand value: %s", got)
				}
			},
		},
		{
			name: "multiple expand",
			opts: ListOptions{
				Expand: []string{"network", "storage", "machineType"},
			},
			response: `{
				"instances": [{
					"id": "inst1",
					"name": "test1",
					"network": {"id": "net1"},
					"storage": {"id": "stor1"},
					"machineType": {"id": "mt1"}
				}]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if !r.URL.Query().Has("expand") {
					t.Error("expected expand parameter, got none")
					return
				}
				got := r.URL.Query().Get("expand")
				switch got {
				case "network":
					// single expand ok
				case "network,storage,machineType":
					// multiple expand ok
				case "invalid":
					// invalid expand ok
				default:
					t.Errorf("unexpected expand value: %s", got)
				}
			},
		},
		{
			name: "invalid expand field",
			opts: ListOptions{
				Expand: []string{"invalid"},
			},
			response:   `{"error": "invalid expand field: invalid"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
			checkQuery: func(t *testing.T, r *http.Request) {
				if !r.URL.Query().Has("expand") {
					t.Error("expected expand parameter, got none")
					return
				}
				got := r.URL.Query().Get("expand")
				switch got {
				case "network":
					// single expand ok
				case "network,storage,machineType":
					// multiple expand ok
				case "invalid":
					// invalid expand ok
				default:
					t.Errorf("unexpected expand value: %s", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if tt.checkQuery != nil {
					tt.checkQuery(t, r)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Instances().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(got) == 0 {
				t.Error("List() expected non-empty result")
			}
		})
	}
}

func TestInstanceService_GetWithExpand(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expand     []string
		response   string
		statusCode int
		wantErr    bool
		checkQuery func(*testing.T, *http.Request)
	}{
		{
			name:   "expand network",
			id:     "inst1",
			expand: []string{"network"},
			response: `{
				"id": "inst1",
				"name": "test-vm",
				"network": {
					"id": "net1",
					"name": "network1",
					"type": "private"
				}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if !r.URL.Query().Has("expand") {
					t.Error("expected expand parameter, got none")
					return
				}
				got := r.URL.Query().Get("expand")
				switch got {
				case "network":
					// single expand ok
				case "network,storage":
					// multiple expand ok
				case "invalid":
					// invalid expand ok
				default:
					t.Errorf("unexpected expand value: %s", got)
				}
			},
		},
		{
			name:   "multiple expands",
			id:     "inst1",
			expand: []string{"network", "storage"},
			response: `{
				"id": "inst1",
				"name": "test-vm",
				"network": {"id": "net1"},
				"storage": {"id": "stor1"}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if !r.URL.Query().Has("expand") {
					t.Error("expected expand parameter, got none")
					return
				}
				got := r.URL.Query().Get("expand")
				switch got {
				case "network":
					// single expand ok
				case "network,storage":
					// multiple expand ok
				case "invalid":
					// invalid expand ok
				default:
					t.Errorf("unexpected expand value: %s", got)
				}
			},
		},
		{
			name:       "empty expand",
			id:         "inst1",
			expand:     []string{},
			response:   `{"id": "inst1", "name": "test-vm"}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if query := r.URL.Query().Encode(); query != "" {
					t.Errorf("expected empty query, got %s", query)
				}
			},
		},
		{
			name:       "invalid expand field",
			id:         "inst1",
			expand:     []string{"invalid"},
			response:   `{"error": "invalid expand field"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
			checkQuery: func(t *testing.T, r *http.Request) {
				if !r.URL.Query().Has("expand") {
					t.Error("expected expand parameter, got none")
					return
				}
				got := r.URL.Query().Get("expand")
				switch got {
				case "network":
					// single expand ok
				case "network,storage":
					// multiple expand ok
				case "invalid":
					// invalid expand ok
				default:
					t.Errorf("unexpected expand value: %s", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if tt.checkQuery != nil {
					tt.checkQuery(t, r)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Instances().Get(context.Background(), tt.id, tt.expand)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Error("Get() got nil, want instance")
				} else if got.ID != tt.id {
					t.Errorf("Get() got ID = %v, want %v", got.ID, tt.id)
				}
			}
		})
	}
}

func TestInstanceService_AttachNetworkInterface(t *testing.T) {
	tests := []struct {
		name       string
		req        NICRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "successful attach",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("inst1")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("nic1")},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "instance not found",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("invalid")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("nic1")},
				},
			},
			response:   `{"error": "instance not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "interface not found",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("inst1")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("invalid")},
				},
			},
			response:   `{"error": "network interface not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "interface already attached",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("inst1")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("nic1")},
				},
			},
			response:   `{"error": "interface already attached"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/compute/v1/instances/network-interface/attach"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Instances().AttachNetworkInterface(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttachNetworkInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceService_DetachNetworkInterface(t *testing.T) {
	tests := []struct {
		name       string
		req        NICRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "successful detach",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("inst1")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("nic1")},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "instance not found",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("invalid")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("nic1")},
				},
			},
			response:   `{"error": "instance not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "interface not found",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("inst1")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("invalid")},
				},
			},
			response:   `{"error": "network interface not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "primary interface",
			req: NICRequest{
				Instance: IDOrName{ID: strPtr("inst1")},
				Network: NICRequestInterface{
					Interface: IDOrName{ID: strPtr("primary")},
				},
			},
			response:   `{"error": "cannot detach primary interface"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/compute/v1/instances/network-interface/detach"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Instances().DetachNetworkInterface(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetachNetworkInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceService_GetFirstWindowsPassword(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		want       *WindowsPasswordResponse
		wantErr    bool
	}{
		{
			name: "successful password retrieval",
			id:   "inst1",
			response: `{
				"instance": {
					"id": "inst1",
					"password": "P@ssw0rd123",
					"created_at": "2023-01-01T00:00:00Z",
					"user": "Administrator"
				}
			}`,
			statusCode: http.StatusOK,
			want: &WindowsPasswordResponse{
				Instance: WindowsPasswordInstance{
					ID:        "inst1",
					Password:  "P@ssw0rd123",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					User:      "Administrator",
				},
			},
			wantErr: false,
		},
		{
			name:       "instance not found",
			id:         "invalid",
			response:   `{"error": "instance not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "empty id",
			id:         "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "non-windows instance",
			id:         "linux-inst",
			response:   `{"error": "not a Windows instance"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "password not ready",
			id:         "new-inst",
			response:   `{"error": "password not yet available"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/compute/v1/instances/config/%s/first-windows-password", tt.id)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Instances().GetFirstWindowsPassword(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFirstWindowsPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("GetFirstWindowsPassword() got nil, want response")
					return
				}
				if got.Instance.ID != tt.want.Instance.ID {
					t.Errorf("GetFirstWindowsPassword() got ID = %v, want %v", got.Instance.ID, tt.want.Instance.ID)
				}
				if got.Instance.Password != tt.want.Instance.Password {
					t.Errorf("GetFirstWindowsPassword() got Password = %v, want %v", got.Instance.Password, tt.want.Instance.Password)
				}
				if !got.Instance.CreatedAt.Equal(tt.want.Instance.CreatedAt) {
					t.Errorf("GetFirstWindowsPassword() got CreatedAt = %v, want %v", got.Instance.CreatedAt, tt.want.Instance.CreatedAt)
				}
			}
		})
	}
}
