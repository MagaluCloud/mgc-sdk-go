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

func TestNetworkACLService_Create_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testACLClient("http://dummy-url")

	req := CreateNetworkACLRequest{
		LoadBalancerID: "lb-123",
		Ethertype:      "IPv4",
		Protocol:       "TCP",
		RemoteIPPrefix: "192.168.1.0/24",
		Action:         "allow",
	}

	_, err := client.Create(ctx, req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkACLService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testACLClient("http://dummy-url")

	req := DeleteNetworkACLRequest{
		LoadBalancerID: "lb-123",
		ID:             "acl-123",
	}

	err := client.Delete(ctx, req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
