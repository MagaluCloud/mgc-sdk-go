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

func testListenerClient(baseURL string) NetworkListenerService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkListeners()
}

func TestNetworkListenerService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		backendID  string
		request    CreateNetworkListenerRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name:      "successful creation",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"id": "listener-123"}`,
			statusCode: http.StatusOK,
			want:       "listener-123",
			wantErr:    false,
		},
		{
			name:      "server error",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:      "bad request - invalid protocol",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "INVALID",
				Port:     80,
			},
			response:   `{"error": "invalid protocol"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:      "bad request - invalid port",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     -1,
			},
			response:   `{"error": "invalid port number"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:      "unauthorized access",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:      "forbidden access",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name:      "load balancer not found",
			lbID:      "invalid-lb",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "load balancer not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:      "backend not found",
			lbID:      "lb-123",
			backendID: "invalid-backend",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "backend not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:      "conflict - port already in use",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkListenerRequest{
				Name:     "test-listener",
				Protocol: "HTTP",
				Port:     80,
			},
			response:   `{"error": "port 80 is already in use"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/listeners", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				backendID := r.URL.Query().Get("backend_id")
				assertEqual(t, tt.backendID, backendID)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testListenerClient(server.URL)
			listener, err := client.Create(context.Background(), tt.lbID, tt.backendID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, listener.ID)
		})
	}
}

func TestNetworkListenerService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		listenerID string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "existing listener",
			lbID:       "lb-123",
			listenerID: "listener-123",
			response: `{
				"id": "listener-123",
				"name": "test-listener",
				"description": "Test listener",
				"protocol": "HTTP",
				"port": 80,
				"backend_id": "backend-123",
				"tls_certificate_id": "cert-123",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent listener",
			lbID:       "lb-123",
			listenerID: "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/listeners/%s", tt.lbID, tt.listenerID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testListenerClient(server.URL)
			listener, err := client.Get(context.Background(), tt.lbID, tt.listenerID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, "listener-123", listener.ID)
			assertEqual(t, "test-listener", listener.Name)
		})
	}
}

func TestNetworkListenerService_List(t *testing.T) {
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
			name: "successful list with multiple listeners",
			lbID: "lb-123",
			response: `{
				"meta": {
					"current_page": 1,
					"total_count": 2,
					"total_pages": 1,
					"total_results": 2
				},
				"results": [
					{"id": "listener-1", "name": "test1", "protocol": "HTTP", "port": 80, "backend_id": "backend-1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
					{"id": "listener-2", "name": "test2", "protocol": "HTTPS", "port": 443, "backend_id": "backend-2", "tls_certificate_id": "cert-1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/listeners", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testListenerClient(server.URL)
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

func TestNetworkListenerService_ListAll(t *testing.T) {
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
							"self": "/listeners?_limit=50&_offset=0"
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
							"id": "listener-1",
							"name": "test1",
							"protocol": "HTTP",
							"port": 80,
							"backend_id": "backend-1",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						},
						{
							"id": "listener-2",
							"name": "test2",
							"protocol": "HTTPS",
							"port": 443,
							"backend_id": "backend-2",
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
							"self": "/listeners?_limit=50&_offset=0",
							"next": "/listeners?_limit=50&_offset=50"
						},
						"page": {
							"count": 50,
							"limit": 50,
							"offset": 0,
							"total": 75
						}
					},
					"results": [` + generateListenerResults(1, 50) + `]
				}`,
				`{
					"meta": {
						"links": {
							"self": "/listeners?_limit=50&_offset=50",
							"previous": "/listeners?_limit=50&_offset=0"
						},
						"page": {
							"count": 25,
							"limit": 50,
							"offset": 50,
							"total": 75
						}
					},
					"results": [` + generateListenerResults(51, 25) + `]
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
							"self": "/listeners?_limit=50&_offset=0"
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/listeners", tt.lbID), r.URL.Path)
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

			client := testListenerClient(server.URL)
			listeners, err := client.ListAll(context.Background(), tt.lbID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(listeners))
		})
	}
}

func generateListenerResults(start, count int) string {
	results := make([]string, count)
	for i := 0; i < count; i++ {
		id := start + i
		results[i] = fmt.Sprintf(`{
			"id": "listener-%d",
			"name": "test%d",
			"protocol": "HTTP",
			"port": 80,
			"backend_id": "backend-%d",
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z"
		}`, id, id, id)
	}
	return strings.Join(results, ",")
}

func TestNetworkListenerService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		listenerID string
		request    UpdateNetworkListenerRequest
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful update",
			lbID:       "lb-123",
			listenerID: "listener-123",
			request: UpdateNetworkListenerRequest{
				TLSCertificateID: stringPtr("updated-listener"),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent listener",
			lbID:       "lb-123",
			listenerID: "invalid",
			request: UpdateNetworkListenerRequest{
				TLSCertificateID: stringPtr("updated-listener"),
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/listeners/%s", tt.lbID, tt.listenerID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testListenerClient(server.URL)
			err := client.Update(context.Background(), tt.lbID, tt.listenerID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkListenerService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		listenerID string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			lbID:       "lb-123",
			listenerID: "listener-123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent listener",
			lbID:       "lb-123",
			listenerID: "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/listeners/%s", tt.lbID, tt.listenerID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testListenerClient(server.URL)
			err := client.Delete(context.Background(), tt.lbID, tt.listenerID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkListenerService_Create_NewRequestError(t *testing.T) {
	t.Parallel()

	// Usar um contexto cancelado para forçar erro no newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela imediatamente

	client := testListenerClient("http://dummy-url")

	req := CreateNetworkListenerRequest{
		Name:     "test-listener",
		Protocol: "HTTP",
		Port:     80,
	}

	_, err := client.Create(ctx, "lb-123", "backend-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkListenerService_Get_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testListenerClient("http://dummy-url")

	_, err := client.Get(ctx, "lb-123", "listener-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkListenerService_List_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testListenerClient("http://dummy-url")

	_, err := client.List(ctx, "lb-123", ListNetworkLoadBalancerRequest{})

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkListenerService_ListAll_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testListenerClient("http://dummy-url")

	_, err := client.ListAll(ctx, "lb-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkListenerService_Update_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testListenerClient("http://dummy-url")

	req := UpdateNetworkListenerRequest{
		TLSCertificateID: stringPtr("updated-listener"),
	}

	err := client.Update(ctx, "lb-123", "listener-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkListenerService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testListenerClient("http://dummy-url")

	err := client.Delete(ctx, "lb-123", "listener-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
