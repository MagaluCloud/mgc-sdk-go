package lbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func testBackendClient(baseURL string) NetworkBackendService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkBackends()
}

func floatPtr(f float64) *float64 {
	return &f
}

func TestNetworkBackendService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		request    CreateBackendRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			lbID: "lb-123",
			request: CreateBackendRequest{
				Name:             "test-backend",
				BalanceAlgorithm: "round_robin",
				TargetsType:      "instance",
			},
			response:   `{"id": "backend-123"}`,
			statusCode: http.StatusOK,
			want:       "backend-123",
			wantErr:    false,
		},
		{
			name: "successful creation with all fields",
			lbID: "lb-123",
			request: CreateBackendRequest{
				Name:                                "test-backend",
				Description:                         stringPtr("Test backend description"),
				BalanceAlgorithm:                    "least_connections",
				TargetsType:                         "raw",
				PanicThreshold:                      floatPtr(75.0),
				HealthCheckName:                     stringPtr("hc-test"),
				HealthCheckID:                       stringPtr("hc-123"),
				CloseConnectionsOnHostHealthFailure: boolPtr(true),
				Targets: &[]NetworkBackendInstanceTargetRequest{
					{
						NicID:     stringPtr("nic-123"),
						IPAddress: stringPtr("192.168.1.10"),
						Port:      8080,
					},
				},
			},
			response:   `{"id": "backend-456"}`,
			statusCode: http.StatusOK,
			want:       "backend-456",
			wantErr:    false,
		},
		{
			name: "server error",
			lbID: "lb-123",
			request: CreateBackendRequest{
				Name:             "test-backend",
				BalanceAlgorithm: "round_robin",
				TargetsType:      "instance",
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "bad request - invalid balance algorithm",
			lbID: "lb-123",
			request: CreateBackendRequest{
				Name:             "test-backend",
				BalanceAlgorithm: "invalid_algorithm",
				TargetsType:      "instance",
			},
			response:   `{"error": "invalid balance algorithm"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized access",
			lbID: "lb-123",
			request: CreateBackendRequest{
				Name:             "test-backend",
				BalanceAlgorithm: "round_robin",
				TargetsType:      "instance",
			},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden access",
			lbID: "lb-123",
			request: CreateBackendRequest{
				Name:             "test-backend",
				BalanceAlgorithm: "round_robin",
				TargetsType:      "instance",
			},
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "load balancer not found",
			lbID: "invalid-lb",
			request: CreateBackendRequest{
				Name:             "test-backend",
				BalanceAlgorithm: "round_robin",
				TargetsType:      "instance",
			},
			response:   `{"error": "load balancer not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "conflict - backend already exists",
			lbID: "lb-123",
			request: CreateBackendRequest{
				Name:             "existing-backend",
				BalanceAlgorithm: "round_robin",
				TargetsType:      "instance",
			},
			response:   `{"error": "backend with this name already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			id, err := client.Create(context.Background(), tt.lbID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, id)
		})
	}
}

func TestNetworkBackendService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		backendID  string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:      "existing backend",
			lbID:      "lb-123",
			backendID: "backend-123",
			response: `{
				"id": "backend-123",
				"name": "test-backend",
				"balance_algorithm": "round_robin",
				"targets_type": "instance",
				"health_check_id": "hc-123",
				"close_connections_on_host_health_failure": false,
				"targets": [],
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent backend",
			lbID:       "lb-123",
			backendID:  "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "unauthorized access",
			lbID:       "lb-123",
			backendID:  "backend-123",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "forbidden access",
			lbID:       "lb-123",
			backendID:  "backend-123",
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s", tt.lbID, tt.backendID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			backend, err := client.Get(context.Background(), tt.lbID, tt.backendID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, "backend-123", backend.ID)
			assertEqual(t, "test-backend", backend.Name)
			assertEqual(t, BackendBalanceAlgorithm("round_robin"), backend.BalanceAlgorithm)
			assertEqual(t, BackendType("instance"), backend.TargetsType)
			assertEqual(t, false, *backend.CloseConnectionsOnHostHealthFailure)
		})
	}
}

func TestNetworkBackendService_List(t *testing.T) {
	t.Parallel()
	sorted := "created_at"
	tests := []struct {
		name       string
		lbID       string
		options    ListNetworkLoadBalancerRequest
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with multiple backends",
			lbID: "lb-123",
			options: ListNetworkLoadBalancerRequest{
				Limit:  intPtr(10),
				Offset: intPtr(0),
				Sort:   stringPtr(sorted),
			},
			response: `{
				"meta": {
					"current_page": 1,
					"total_count": 2,
					"total_pages": 1,
					"total_results": 2
				},
				"results": [
					{
						"id": "backend-1",
						"name": "test1",
						"balance_algorithm": "round_robin",
						"targets_type": "instance",
						"close_connections_on_host_health_failure": false,
						"targets": [],
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					},
					{
						"id": "backend-2",
						"name": "test2",
						"balance_algorithm": "least_connections",
						"targets_type": "raw",
						"close_connections_on_host_health_failure": true,
						"targets": [],
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "empty list",
			lbID: "lb-123",
			options: ListNetworkLoadBalancerRequest{
				Limit:  intPtr(10),
				Offset: intPtr(0),
			},
			response: `{
				"meta": {
					"current_page": 1,
					"total_count": 0,
					"total_pages": 0,
					"total_results": 0
				},
				"results": []
			}`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name: "list with pagination",
			lbID: "lb-123",
			options: ListNetworkLoadBalancerRequest{
				Limit:  intPtr(1),
				Offset: intPtr(1),
				Sort:   &sorted,
			},
			response: `{
				"meta": {
					"current_page": 2,
					"total_count": 2,
					"total_pages": 2,
					"total_results": 1
				},
				"results": [
					{
						"id": "backend-2",
						"name": "test2",
						"balance_algorithm": "least_connections",
						"targets_type": "raw",
						"close_connections_on_host_health_failure": false,
						"targets": [],
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name:       "server error",
			lbID:       "lb-123",
			options:    ListNetworkLoadBalancerRequest{},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "unauthorized access",
			lbID:       "lb-123",
			options:    ListNetworkLoadBalancerRequest{},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				// Verify query parameters if specified
				if tt.options.Limit != nil {
					assertEqual(t, strconv.Itoa(*tt.options.Limit), r.URL.Query().Get("_limit"))
				}
				if tt.options.Offset != nil {
					assertEqual(t, strconv.Itoa(*tt.options.Offset), r.URL.Query().Get("_offset"))
				}
				if tt.options.Sort != nil {
					assertEqual(t, *tt.options.Sort, r.URL.Query().Get("_sort"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			backends, err := client.List(context.Background(), tt.lbID, tt.options)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(backends))
		})
	}
}

func TestNetworkBackendService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		backendID  string
		request    UpdateNetworkBackendRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name:      "successful update with panic threshold",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: UpdateNetworkBackendRequest{
				PanicThreshold: intPtr(50),
			},
			response:   `{"id": "backend-123"}`,
			statusCode: http.StatusOK,
			want:       "backend-123",
			wantErr:    false,
		},
		{
			name:      "successful update with health check",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: UpdateNetworkBackendRequest{
				HealthCheckID: stringPtr("hc-456"),
			},
			response:   `{"id": "backend-123"}`,
			statusCode: http.StatusOK,
			want:       "backend-123",
			wantErr:    false,
		},
		{
			name:      "successful update with all fields",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: UpdateNetworkBackendRequest{
				HealthCheckID:                       stringPtr("hc-789"),
				PanicThreshold:                      intPtr(75),
				CloseConnectionsOnHostHealthFailure: boolPtr(true),
			},
			response:   `{"id": "backend-123"}`,
			statusCode: http.StatusOK,
			want:       "backend-123",
			wantErr:    false,
		},
		{
			name:      "non-existent backend",
			lbID:      "lb-123",
			backendID: "invalid",
			request: UpdateNetworkBackendRequest{
				PanicThreshold: intPtr(50),
			},
			response:   `{"error": "backend not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:      "bad request - invalid panic threshold",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: UpdateNetworkBackendRequest{
				PanicThreshold: intPtr(-10),
			},
			response:   `{"error": "invalid panic threshold"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:      "unauthorized access",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: UpdateNetworkBackendRequest{
				PanicThreshold: intPtr(50),
			},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s", tt.lbID, tt.backendID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			id, err := client.Update(context.Background(), tt.lbID, tt.backendID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				assertEqual(t, "", id)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, id)
		})
	}
}

func TestNetworkBackendService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		backendID  string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			lbID:       "lb-123",
			backendID:  "backend-123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "successful deletion with no content",
			lbID:       "lb-123",
			backendID:  "backend-456",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent backend",
			lbID:       "lb-123",
			backendID:  "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "unauthorized access",
			lbID:       "lb-123",
			backendID:  "backend-123",
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "forbidden access",
			lbID:       "lb-123",
			backendID:  "backend-123",
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name:       "server error",
			lbID:       "lb-123",
			backendID:  "backend-123",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s", tt.lbID, tt.backendID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			err := client.Delete(context.Background(), tt.lbID, tt.backendID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkBackendService_Create_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendClient("http://dummy-url")

	req := CreateBackendRequest{
		Name:             "test-backend",
		BalanceAlgorithm: "round_robin",
		TargetsType:      "instance",
	}

	_, err := client.Create(ctx, "lb-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkBackendService_Get_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendClient("http://dummy-url")

	_, err := client.Get(ctx, "lb-123", "backend-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkBackendService_List_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendClient("http://dummy-url")

	_, err := client.List(ctx, "lb-123", ListNetworkLoadBalancerRequest{})

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkBackendService_Update_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendClient("http://dummy-url")

	req := UpdateNetworkBackendRequest{
		PanicThreshold: intPtr(50),
	}

	_, err := client.Update(ctx, "lb-123", "backend-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkBackendService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendClient("http://dummy-url")

	err := client.Delete(ctx, "lb-123", "backend-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

// Helper function for int pointers
func intPtr(i int) *int {
	return &i
}
