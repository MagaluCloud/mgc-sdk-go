package lbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

// Helper functions
func assertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v but got %v. %v", expected, actual, msgAndArgs)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func testLoadBalancerClient(baseURL string) NetworkLoadBalancerService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkLoadBalancers()
}

func TestNetworkLoadBalancerService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		request    CreateNetworkLoadBalancerRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []NetworkBackendRequest{},
			},
			response:   `{"id": "lb-123"}`,
			statusCode: http.StatusOK,
			want:       "lb-123",
			wantErr:    false,
		},
		{
			name: "bad request - invalid request data",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "",
				Visibility: "invalid",
				VPCID:      "",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []NetworkBackendRequest{},
			},
			response:   `{"error": "invalid request: name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized - invalid credentials",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []NetworkBackendRequest{},
			},
			response:   `{"error": "unauthorized access"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden - insufficient permissions",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []NetworkBackendRequest{},
			},
			response:   `{"error": "forbidden: insufficient permissions"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "conflict - resource already exists",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "existing-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []NetworkBackendRequest{},
			},
			response:   `{"error": "load balancer with name 'existing-lb' already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "server error",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []NetworkBackendRequest{},
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
				assertEqual(t, "/load-balancer/v0beta1/network-load-balancers", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
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

func TestNetworkLoadBalancerService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "existing load balancer",
			lbID: "lb-123",
			response: `{
				"id": "lb-123",
				"name": "test-lb",
				"type": "proxy",
				"visibility": "external",
				"status": "running",
				"listeners": [],
				"backends": [],
				"health_checks": [],
				"public_ips": [],
				"tls_certificates": [],
				"acls": [],
				"vpc_id": "vpc-123",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent load balancer",
			lbID:       "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "unauthorized access",
			lbID:       "lb-123",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "forbidden access",
			lbID:       "lb-123",
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			lb, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{LoadBalancerID: tt.lbID})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, "lb-123", lb.ID)
			assertEqual(t, "test-lb", lb.Name)
		})
	}
}

func TestNetworkLoadBalancerService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with multiple load balancers",
			response: `{
				"results": [
					{"id": "lb-1", "name": "test1", "type": "proxy", "visibility": "external", "status": "running", "listeners": [], "backends": [], "health_checks": [], "public_ips": [], "tls_certificates": [], "acls": [], "vpc_id": "vpc-1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
					{"id": "lb-2", "name": "test2", "type": "proxy", "visibility": "internal", "status": "running", "listeners": [], "backends": [], "health_checks": [], "public_ips": [], "tls_certificates": [], "acls": [], "vpc_id": "vpc-2", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			response:   `{"results": []}`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name:       "unauthorized access",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "forbidden access",
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name:       "server error",
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
				assertEqual(t, "/load-balancer/v0beta1/network-load-balancers", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			lbs, err := client.List(context.Background(), ListNetworkLoadBalancerRequest{})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(lbs))
		})
	}
}

func TestNetworkLoadBalancerService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		request    UpdateNetworkLoadBalancerRequest
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful update",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
				Name:           stringPtr("updated-lb"),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non-existent load balancer",
			lbID: "invalid",
			request: UpdateNetworkLoadBalancerRequest{
				LoadBalancerID: "invalid",
				Name:           stringPtr("updated-lb"),
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "bad request - invalid data",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
				Name:           stringPtr(""),
			},
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized access",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
				Name:           stringPtr("updated-lb"),
			},
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden access",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
				Name:           stringPtr("updated-lb"),
			},
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "conflict - name already exists",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
				Name:           stringPtr("existing-name"),
			},
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "server error",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
				Name:           stringPtr("updated-lb"),
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			err := client.Update(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkLoadBalancerService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		request    DeleteNetworkLoadBalancerRequest
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful deletion",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "deletion with delete public IP",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
				DeletePublicIP: boolPtr(true),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non-existent load balancer",
			lbID: "invalid",
			request: DeleteNetworkLoadBalancerRequest{
				LoadBalancerID: "invalid",
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "unauthorized access",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
			},
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden access",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
			},
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "conflict - resource in use",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
			},
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "server error",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				LoadBalancerID: "lb-123",
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)

				// Check query parameter if DeletePublicIP is set
				if tt.request.DeletePublicIP != nil {
					assertEqual(t, "true", r.URL.Query().Get("delete_public_ip"))
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			err := client.Delete(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

// Helper functions for pointer values
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
