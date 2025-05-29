package lbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestNetworkBackendService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		request    CreateNetworkBackendRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateNetworkBackendRequest{
				LoadBalancerID:   "lb-123",
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
			name: "server error",
			request: CreateNetworkBackendRequest{
				LoadBalancerID:   "lb-123",
				Name:             "test-backend",
				BalanceAlgorithm: "round_robin",
				TargetsType:      "instance",
			},
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends", tt.request.LoadBalancerID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
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
				"targets": []
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
			backend, err := client.Get(context.Background(), GetNetworkBackendRequest{
				LoadBalancerID: tt.lbID,
				BackendID:      tt.backendID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, "backend-123", backend.ID)
			assertEqual(t, "test-backend", backend.Name)
		})
	}
}

func TestNetworkBackendService_List(t *testing.T) {
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
			name: "successful list with multiple backends",
			lbID: "lb-123",
			response: `{
				"meta": {
					"current_page": 1,
					"total_count": 2,
					"total_pages": 1,
					"total_results": 2
				},
				"results": [
					{"id": "backend-1", "name": "test1", "balance_algorithm": "round_robin", "targets_type": "instance", "targets": []},
					{"id": "backend-2", "name": "test2", "balance_algorithm": "least_connections", "targets_type": "raw", "targets": []}
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			backends, err := client.List(context.Background(), ListNetworkBackendRequest{
				LoadBalancerID: tt.lbID,
			})

			if tt.wantErr {
				assertError(t, err)
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
		statusCode int
		wantErr    bool
	}{
		{
			name:      "successful update",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: UpdateNetworkBackendRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "backend-123",
				Name:           stringPtr("updated-backend"),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "non-existent backend",
			lbID:      "lb-123",
			backendID: "invalid",
			request: UpdateNetworkBackendRequest{
				LoadBalancerID: "lb-123",
				BackendID:      "invalid",
				Name:           stringPtr("updated-backend"),
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s", tt.lbID, tt.backendID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testBackendClient(server.URL)
			err := client.Update(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
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
			name:       "non-existent backend",
			lbID:       "lb-123",
			backendID:  "invalid",
			statusCode: http.StatusNotFound,
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
			err := client.Delete(context.Background(), DeleteNetworkBackendRequest{
				LoadBalancerID: tt.lbID,
				BackendID:      tt.backendID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}
