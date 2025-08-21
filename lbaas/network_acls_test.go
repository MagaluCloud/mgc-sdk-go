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
		lbID       string
		request    CreateNetworkACLRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			lbID: "lb-123",
			request: CreateNetworkACLRequest{
				Name:           stringPtr("test-acl"),
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
			name: "successful creation without name",
			lbID: "lb-456",
			request: CreateNetworkACLRequest{
				Ethertype:      "IPv6",
				Protocol:       "UDP",
				RemoteIPPrefix: "2001:db8::/32",
				Action:         "deny",
			},
			response:   `{"id": "acl-456"}`,
			statusCode: http.StatusOK,
			want:       "acl-456",
			wantErr:    false,
		},
		{
			name: "server error",
			lbID: "lb-789",
			request: CreateNetworkACLRequest{
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/acls", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testACLClient(server.URL)
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
		{
			name:       "server error",
			lbID:       "lb-456",
			aclID:      "acl-456",
			statusCode: http.StatusInternalServerError,
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
			err := client.Delete(context.Background(), tt.lbID, tt.aclID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkACLService_Replace(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		request    UpdateNetworkACLRequest
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful replace",
			lbID: "lb-123",
			request: UpdateNetworkACLRequest{
				Acls: []CreateNetworkACLRequest{
					{
						Name:           stringPtr("acl-1"),
						Ethertype:      "IPv4",
						Protocol:       "TCP",
						RemoteIPPrefix: "192.168.1.0/24",
						Action:         "allow",
					},
					{
						Ethertype:      "IPv4",
						Protocol:       "UDP",
						RemoteIPPrefix: "10.0.0.0/8",
						Action:         "deny",
					},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "server error",
			lbID: "lb-456",
			request: UpdateNetworkACLRequest{
				Acls: []CreateNetworkACLRequest{
					{
						Ethertype:      "IPv4",
						Protocol:       "TCP",
						RemoteIPPrefix: "192.168.1.0/24",
						Action:         "allow",
					},
				},
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/acls", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testACLClient(server.URL)
			err := client.Replace(context.Background(), tt.lbID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
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
		Ethertype:      "IPv4",
		Protocol:       "TCP",
		RemoteIPPrefix: "192.168.1.0/24",
		Action:         "allow",
	}

	_, err := client.Create(ctx, "lb-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkACLService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testACLClient("http://dummy-url")

	err := client.Delete(ctx, "lb-123", "acl-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkACLService_Replace_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testACLClient("http://dummy-url")

	req := UpdateNetworkACLRequest{
		Acls: []CreateNetworkACLRequest{
			{
				Ethertype:      "IPv4",
				Protocol:       "TCP",
				RemoteIPPrefix: "192.168.1.0/24",
				Action:         "allow",
			},
		},
	}

	err := client.Replace(ctx, "lb-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
