package lbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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
		request    CreateNetworkListenerRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           80,
			},
			response:   `{"id": "listener-123"}`,
			statusCode: http.StatusOK,
			want:       "listener-123",
			wantErr:    false,
		},
		{
			name: "server error",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           80,
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "bad request - invalid protocol",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "INVALID",
				Port:           80,
			},
			response:   `{"error": "invalid protocol"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "bad request - invalid port",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           -1,
			},
			response:   `{"error": "invalid port number"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized access",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           80,
			},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden access",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           80,
			},
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "load balancer not found",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "invalid-lb",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           80,
			},
			response:   `{"error": "load balancer not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "backend not found",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "invalid-backend",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           80,
			},
			response:   `{"error": "backend not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "conflict - port already in use",
			request: CreateNetworkListenerRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           "test-listener",
				Protocol:       "HTTP",
				Port:           80,
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/listeners", tt.request.LoadBalancerID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				backendID := r.URL.Query().Get("backend_id")
				assertEqual(t, tt.request.BackendID, backendID)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testListenerClient(server.URL)
			listener, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
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
				"tls_certificate_id": "cert-123"
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
			listener, err := client.Get(context.Background(), GetNetworkListenerRequest{
				LoadBalancerID: tt.lbID,
				ListenerID:     tt.listenerID,
			})

			if tt.wantErr {
				assertError(t, err)
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
					{"id": "listener-1", "name": "test1", "protocol": "HTTP", "port": 80, "backend_id": "backend-1"},
					{"id": "listener-2", "name": "test2", "protocol": "HTTPS", "port": 443, "backend_id": "backend-2", "tls_certificate_id": "cert-1"}
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
			listeners, err := client.List(context.Background(), ListNetworkListenerRequest{
				LoadBalancerID: tt.lbID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(listeners))
		})
	}
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
				LoadBalancerID:   "lb-123",
				ListenerID:       "listener-123",
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
				LoadBalancerID:   "lb-123",
				ListenerID:       "invalid",
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
			err := client.Update(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
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
			err := client.Delete(context.Background(), DeleteNetworkListenerRequest{
				LoadBalancerID: tt.lbID,
				ListenerID:     tt.listenerID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}
