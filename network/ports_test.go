package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
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

func testClient(baseURL string) PortService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Ports()
}

func TestPortService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with multiple ports",
			response: `[
				{"id": "port1", "name": "test1"},
				{"id": "port2", "name": "test2"}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			response:   `[]`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v0/ports", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			ports, err := client.List(context.Background())

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(ports))
		})
	}
}

func TestPortService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		portID     string
		response   string
		statusCode int
		want       *PortResponse
		wantErr    bool
	}{
		{
			name:   "existing port",
			portID: "port1",
			response: `{
				"id": "port1",
				"name": "test-port",
				"vpc_id": "vpc1",
				"security_groups": ["sg1", "sg2"],
				"public_ip": [{"public_ip_id": "ip1", "public_ip": "203.0.113.5"}],
				"ip_address": [{"ip_address": "10.0.0.2", "subnet_id": "subnet1"}]
			}`,
			statusCode: http.StatusOK,
			want: &PortResponse{
				ID:             helpers.StrPtr("port1"),
				Name:           helpers.StrPtr("test-port"),
				VPCID:          helpers.StrPtr("vpc1"),
				SecurityGroups: &[]string{"sg1", "sg2"},
				PublicIP: &[]PublicIpResponsePort{
					{
						PublicIPID: helpers.StrPtr("ip1"),
						PublicIP:   helpers.StrPtr("203.0.113.5"),
					},
				},
				IPAddress: &[]IpAddress{
					{
						IPAddress: "10.0.0.2",
						SubnetID:  "subnet1",
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "non-existent port",
			portID:     "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			portID:     "port1",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/ports/%s", tt.portID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			port, err := client.Get(context.Background(), tt.portID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, *tt.want.ID, *port.ID)
			assertEqual(t, *tt.want.Name, *port.Name)
			assertEqual(t, *tt.want.VPCID, *port.VPCID)
			assertEqual(t, len(*tt.want.SecurityGroups), len(*port.SecurityGroups))
			assertEqual(t, len(*tt.want.PublicIP), len(*port.PublicIP))
			assertEqual(t, len(*tt.want.IPAddress), len(*port.IPAddress))
		})
	}
}

func TestPortService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		portID     string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			portID:     "port1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent port",
			portID:     "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			portID:     "port1",
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/ports/%s", tt.portID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Delete(context.Background(), tt.portID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestPortService_AttachSecurityGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		portID          string
		securityGroupID string
		statusCode      int
		response        string
		wantErr         bool
	}{
		{
			name:            "successful attach",
			portID:          "port1",
			securityGroupID: "sg1",
			statusCode:      http.StatusOK,
			wantErr:         false,
		},
		{
			name:            "port not found",
			portID:          "invalid",
			securityGroupID: "sg1",
			statusCode:      http.StatusNotFound,
			response:        `{"error": "port not found"}`,
			wantErr:         true,
		},
		{
			name:            "security group not found",
			portID:          "port1",
			securityGroupID: "invalid",
			statusCode:      http.StatusBadRequest,
			response:        `{"error": "security group not found"}`,
			wantErr:         true,
		},
		{
			name:            "server error",
			portID:          "port1",
			securityGroupID: "sg1",
			statusCode:      http.StatusInternalServerError,
			response:        `{"error": "internal server error"}`,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/network/v0/ports/%s/attach/%s", tt.portID, tt.securityGroupID)
				assertEqual(t, expectedPath, r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.AttachSecurityGroup(context.Background(), tt.portID, tt.securityGroupID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestPortService_DetachSecurityGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		portID          string
		securityGroupID string
		statusCode      int
		response        string
		wantErr         bool
	}{
		{
			name:            "successful detach",
			portID:          "port1",
			securityGroupID: "sg1",
			statusCode:      http.StatusOK,
			wantErr:         false,
		},
		{
			name:            "port not found",
			portID:          "invalid",
			securityGroupID: "sg1",
			statusCode:      http.StatusNotFound,
			response:        `{"error": "port not found"}`,
			wantErr:         true,
		},
		{
			name:            "security group not attached",
			portID:          "port1",
			securityGroupID: "sg2",
			statusCode:      http.StatusConflict,
			response:        `{"error": "security group not attached"}`,
			wantErr:         true,
		},
		{
			name:            "server error",
			portID:          "port1",
			securityGroupID: "sg1",
			statusCode:      http.StatusInternalServerError,
			response:        `{"error": "internal server error"}`,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/network/v0/ports/%s/detach/%s", tt.portID, tt.securityGroupID)
				assertEqual(t, expectedPath, r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.DetachSecurityGroup(context.Background(), tt.portID, tt.securityGroupID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestPortService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		portID     string
		request    PortUpdateRequest
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful update",
			portID:     "port1",
			request:    PortUpdateRequest{IPSpoofingGuard: helpers.BoolPtr(false)},
			statusCode: http.StatusNoContent,
		},
		{
			name:       "update failed - port not found",
			portID:     "port2",
			request:    PortUpdateRequest{IPSpoofingGuard: helpers.BoolPtr(true)},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/ports/%s", tt.portID), r.URL.Path)
				assertEqual(t, http.MethodPatch, r.Method)
				w.WriteHeader(tt.statusCode)

				if !tt.wantErr {
					var req PortUpdateRequest
					err := json.NewDecoder(r.Body).Decode(&req)
					assertNoError(t, err)

					assertEqual(t, *tt.request.IPSpoofingGuard, *req.IPSpoofingGuard)

				}
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.Update(context.Background(), tt.portID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}
