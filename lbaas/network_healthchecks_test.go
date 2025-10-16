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

func testHealthCheckClient(baseURL string) NetworkHealthCheckService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkHealthChecks()
}

func TestNetworkHealthCheckService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		request    CreateNetworkHealthCheckRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			lbID: "lb-123",
			request: CreateNetworkHealthCheckRequest{
				Name:     "test-hc",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"id": "hc-123"}`,
			statusCode: http.StatusOK,
			want:       "hc-123",
			wantErr:    false,
		},
		{
			name: "server error",
			lbID: "lb-123",
			request: CreateNetworkHealthCheckRequest{
				Name:     "test-hc",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "bad request - invalid protocol",
			lbID: "lb-123",
			request: CreateNetworkHealthCheckRequest{
				Name:     "test-hc",
				Protocol: "INVALID",
				Port:     80,
			},
			response:   `{"error": "invalid protocol"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized access",
			lbID: "lb-123",
			request: CreateNetworkHealthCheckRequest{
				Name:     "test-hc",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden access",
			lbID: "lb-123",
			request: CreateNetworkHealthCheckRequest{
				Name:     "test-hc",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "load balancer not found",
			lbID: "invalid-lb",
			request: CreateNetworkHealthCheckRequest{
				Name:     "test-hc",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "load balancer not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "conflict - healthcheck already exists",
			lbID: "lb-123",
			request: CreateNetworkHealthCheckRequest{
				Name:     "existing-hc",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "healthcheck with this name already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/health-checks", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testHealthCheckClient(server.URL)
			hc, err := client.Create(context.Background(), tt.lbID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, hc.ID)
		})
	}
}

func TestNetworkHealthCheckService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		hcID       string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "existing health check",
			lbID: "lb-123",
			hcID: "hc-123",
			response: `{
				"id": "hc-123",
				"name": "test-hc",
				"description": "Test health check",
				"protocol": "HTTP",
				"path": "/health",
				"port": 80,
				"healthy_status_code": 200,
				"interval_seconds": 30,
				"timeout_seconds": 5,
				"initial_delay_seconds": 10,
				"healthy_threshold_count": 3,
				"unhealthy_threshold_count": 3,
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent health check",
			lbID:       "lb-123",
			hcID:       "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "unauthorized access",
			lbID:       "lb-123",
			hcID:       "hc-123",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "forbidden access",
			lbID:       "lb-123",
			hcID:       "hc-123",
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name:       "load balancer not found",
			lbID:       "invalid-lb",
			hcID:       "hc-123",
			response:   `{"error": "load balancer not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			lbID:       "lb-123",
			hcID:       "hc-123",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/health-checks/%s", tt.lbID, tt.hcID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testHealthCheckClient(server.URL)
			hc, err := client.Get(context.Background(), tt.lbID, tt.hcID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, "hc-123", hc.ID)
			assertEqual(t, "test-hc", hc.Name)
		})
	}
}

func TestNetworkHealthCheckService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with multiple health checks",
			lbID: "lb-123",
			response: `{
				"meta": {
					"current_page": 1,
					"total_count": 2,
					"total_pages": 1,
					"total_results": 2
				},
				"results": [
					{"id": "hc-1", "name": "test1", "protocol": "HTTP", "port": 80, "healthy_status_code": 200, "interval_seconds": 30, "timeout_seconds": 5, "initial_delay_seconds": 10, "healthy_threshold_count": 3, "unhealthy_threshold_count": 3, "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
					{"id": "hc-2", "name": "test2", "protocol": "TCP", "port": 443, "interval_seconds": 30, "timeout_seconds": 5, "initial_delay_seconds": 10, "healthy_threshold_count": 3, "unhealthy_threshold_count": 3, "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "empty list",
			lbID: "lb-123",
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
			name:       "server error",
			lbID:       "lb-123",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/health-checks", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testHealthCheckClient(server.URL)
			resp, err := client.List(context.Background(), tt.lbID, ListNetworkLoadBalancerRequest{})

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(resp.Results))
		})
	}
}

func TestNetworkHealthCheckService_ListAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		lbID       string
		pages      []string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page with all results",
			lbID: "lb-123",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/health-checks?_limit=50&_offset=0"
						},
						"page": {
							"count": 2,
							"limit": 50,
							"offset": 0,
							"total": 2
						}
					},
					"results": [
						{
							"id": "hc-1",
							"name": "test1",
							"protocol": "HTTP",
							"port": 80,
							"healthy_status_code": 200,
							"interval_seconds": 30,
							"timeout_seconds": 5,
							"initial_delay_seconds": 10,
							"healthy_threshold_count": 3,
							"unhealthy_threshold_count": 3,
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						},
						{
							"id": "hc-2",
							"name": "test2",
							"protocol": "TCP",
							"port": 443,
							"healthy_status_code": 200,
							"interval_seconds": 30,
							"timeout_seconds": 5,
							"initial_delay_seconds": 10,
							"healthy_threshold_count": 3,
							"unhealthy_threshold_count": 3,
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						}
					]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "multiple pages",
			lbID: "lb-456",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/health-checks?_limit=50&_offset=0",
							"next": "/health-checks?_limit=50&_offset=50"
						},
						"page": {
							"count": 50,
							"limit": 50,
							"offset": 0,
							"total": 75
						}
					},
					"results": [` + generateHealthCheckResults(1, 50) + `]
				}`,
				`{
					"meta": {
						"links": {
							"self": "/health-checks?_limit=50&_offset=50",
							"previous": "/health-checks?_limit=50&_offset=0"
						},
						"page": {
							"count": 25,
							"limit": 50,
							"offset": 50,
							"total": 75
						}
					},
					"results": [` + generateHealthCheckResults(51, 25) + `]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  75,
			wantErr:    false,
		},
		{
			name: "empty results",
			lbID: "lb-789",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/health-checks?_limit=50&_offset=0"
						},
						"page": {
							"count": 0,
							"limit": 50,
							"offset": 0,
							"total": 0
						}
					},
					"results": []
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pageIndex := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/health-checks", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if pageIndex < len(tt.pages) {
					w.Write([]byte(tt.pages[pageIndex]))
					pageIndex++
				} else {
					w.Write([]byte(`{"meta":{"links":{"self":""},"page":{"count":0,"limit":50,"offset":0,"total":0}},"results":[]}`))
				}
			}))
			defer server.Close()

			client := testHealthCheckClient(server.URL)
			healthChecks, err := client.ListAll(context.Background(), tt.lbID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(healthChecks))
		})
	}
}

func generateHealthCheckResults(start, count int) string {
	results := make([]string, count)
	for i := 0; i < count; i++ {
		id := start + i
		results[i] = fmt.Sprintf(`{
			"id": "hc-%d",
			"name": "test%d",
			"protocol": "HTTP",
			"port": 80,
			"healthy_status_code": 200,
			"interval_seconds": 30,
			"timeout_seconds": 5,
			"initial_delay_seconds": 10,
			"healthy_threshold_count": 3,
			"unhealthy_threshold_count": 3,
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z"
		}`, id, id)
	}
	return strings.Join(results, ",")
}

func TestNetworkHealthCheckService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		hcID       string
		request    UpdateNetworkHealthCheckRequest
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful update",
			lbID: "lb-123",
			hcID: "hc-123",
			request: UpdateNetworkHealthCheckRequest{
				Path: stringPtr("updated-hc"),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non-existent health check",
			lbID: "lb-123",
			hcID: "invalid",
			request: UpdateNetworkHealthCheckRequest{
				Path: stringPtr("updated-hc"),
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/health-checks/%s", tt.lbID, tt.hcID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testHealthCheckClient(server.URL)
			err := client.Update(context.Background(), tt.lbID, tt.hcID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkHealthCheckService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		hcID       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			lbID:       "lb-123",
			hcID:       "hc-123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent health check",
			lbID:       "lb-123",
			hcID:       "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/health-checks/%s", tt.lbID, tt.hcID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testHealthCheckClient(server.URL)
			err := client.Delete(context.Background(), tt.lbID, tt.hcID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkHealthCheckService_Create_NewRequestError(t *testing.T) {
	t.Parallel()

	// Usar um contexto cancelado para forçar erro no newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela imediatamente

	client := testHealthCheckClient("http://dummy-url")

	req := CreateNetworkHealthCheckRequest{
		Name:     "test-hc",
		Protocol: "HTTP",
		Port:     80,
	}

	_, err := client.Create(ctx, "lb-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkHealthCheckService_Get_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testHealthCheckClient("http://dummy-url")

	_, err := client.Get(ctx, "lb-123", "hc-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkHealthCheckService_List_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testHealthCheckClient("http://dummy-url")

	_, err := client.List(ctx, "lb-123", ListNetworkLoadBalancerRequest{})

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkHealthCheckService_ListAll_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testHealthCheckClient("http://dummy-url")

	_, err := client.ListAll(ctx, "lb-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkHealthCheckService_Update_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testHealthCheckClient("http://dummy-url")

	req := UpdateNetworkHealthCheckRequest{
		Path: stringPtr("updated-hc"),
	}

	err := client.Update(ctx, "lb-123", "hc-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkHealthCheckService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testHealthCheckClient("http://dummy-url")

	err := client.Delete(ctx, "lb-123", "hc-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
