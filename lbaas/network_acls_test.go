package lbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func testACLClient(baseURL string) NetworkACLService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkACLs()
}

func TestNetworkACLService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		request    CreateNetworkACLRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateNetworkACLRequest{
				LoadBalancerID: "lb-123",
				Ethertype:      "IPv4",
				Protocol:       "TCP",
				RemoteIPPrefix: "192.168.1.0/24",
				Action:         "allow",
			},
			response:   `{"id": "acl-123"}`,
			statusCode: http.StatusOK,
			want:       "acl-123",
			wantErr:    false,
		},
		{
			name: "server error",
			request: CreateNetworkACLRequest{
				LoadBalancerID: "lb-123",
				Ethertype:      "IPv4",
				Protocol:       "TCP",
				RemoteIPPrefix: "192.168.1.0/24",
				Action:         "allow",
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/acls", tt.request.LoadBalancerID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testACLClient(server.URL)
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

func TestNetworkACLService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		aclID      string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:  "existing ACL",
			lbID:  "lb-123",
			aclID: "acl-123",
			response: `{
				"id": "acl-123",
				"name": "test-acl",
				"ethertype": "IPv4",
				"protocol": "TCP",
				"remote_ip_prefix": "192.168.1.0/24",
				"action": "allow"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent ACL",
			lbID:       "lb-123",
			aclID:      "invalid",
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/acls/%s", tt.lbID, tt.aclID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testACLClient(server.URL)
			acl, err := client.Get(context.Background(), GetNetworkACLRequest{
				LoadBalancerID: tt.lbID,
				NetworkACLID:   tt.aclID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, "acl-123", acl.ID)
			assertEqual(t, "IPv4", acl.Ethertype)
		})
	}
}

func TestNetworkACLService_List(t *testing.T) {
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
			name: "successful list with multiple ACLs",
			lbID: "lb-123",
			response: `[
				{"id": "acl-1", "name": "test1", "ethertype": "IPv4", "protocol": "TCP", "remote_ip_prefix": "192.168.1.0/24", "action": "allow"},
				{"id": "acl-2", "name": "test2", "ethertype": "IPv4", "protocol": "UDP", "remote_ip_prefix": "10.0.0.0/8", "action": "deny"}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			lbID:       "lb-123",
			response:   `[]`,
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/acls", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testACLClient(server.URL)
			acls, err := client.List(context.Background(), ListNetworkACLRequest{
				LoadBalancerID: tt.lbID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(acls))
		})
	}
}

func TestNetworkACLService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		aclID      string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			lbID:       "lb-123",
			aclID:      "acl-123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent ACL",
			lbID:       "lb-123",
			aclID:      "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/acls/%s", tt.lbID, tt.aclID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testACLClient(server.URL)
			err := client.Delete(context.Background(), DeleteNetworkACLRequest{
				LoadBalancerID: tt.lbID,
				ID:             tt.aclID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}
