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
			listeners, err := client.List(context.Background(), tt.lbID, ListNetworkLoadBalancerRequest{})

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
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
