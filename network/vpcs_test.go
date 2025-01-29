package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestVPCService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateVPCRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful create",
			request: CreateVPCRequest{
				Name:        "prod-vpc",
				Description: "Production VPC",
			},
			response:   `{"id": "vpc1", "status": "creating"}`,
			statusCode: http.StatusCreated,
			wantID:     "vpc1",
			wantErr:    false,
		},
		{
			name: "missing name",
			request: CreateVPCRequest{
				Description: "Invalid VPC",
			},
			response:   `{"error": "name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/vpcs", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req CreateVPCRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.request.Name, req.Name)
				assertEqual(t, tt.request.Description, req.Description)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestVPCService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "vpc1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent vpc",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "vpc not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v1/vpcs/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestVPCService_Rename(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		newName    string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful rename",
			id:         "vpc1",
			newName:    "new-prod-vpc",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "invalid name",
			id:         "vpc1",
			newName:    "",
			statusCode: http.StatusBadRequest,
			response:   `{"error": "name cannot be empty"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/vpcs/%s/rename", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPatch, r.Method)

				var req RenameVPCRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.newName, req.Name)

				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			err := client.Rename(context.Background(), tt.id, tt.newName)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestVPCService_ListPorts(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		detailed   bool
		opts       ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:     "detailed ports list",
			vpcID:    "vpc1",
			detailed: true,
			opts: ListOptions{
				Limit:  helpers.IntPtr(10),
				Offset: helpers.IntPtr(20),
			},
			response: `{
				"ports": [
					{"id": "port1", "name": "web-port"},
					{"id": "port2", "name": "db-port"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:     "simplified ports list",
			vpcID:    "vpc1",
			detailed: false,
			response: `{
				"ports_simplified": [
					{"id": "port1", "ip_address": [{"ip_address": "10.0.0.1"}]},
					{"id": "port2", "ip_address": [{"ip_address": "10.0.0.2"}]}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/vpcs/%s/ports", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				query := r.URL.Query()
				if tt.opts.Limit != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Limit), query.Get("_limit"))
				}
				assertEqual(t, strconv.FormatBool(tt.detailed), query.Get("detailed"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			result, err := client.ListPorts(context.Background(), tt.vpcID, tt.detailed, tt.opts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)

			switch res := result.(type) {
			case []PortResponse:
				assertEqual(t, tt.wantCount, len(res))
			case []PortSimpleResponse:
				assertEqual(t, tt.wantCount, len(res))
			default:
				t.Errorf("unexpected response type")
			}
		})
	}
}

func TestVPCService_CreatePort(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		request    PortCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:  "successful port creation",
			vpcID: "vpc1",
			request: PortCreateRequest{
				Name:           "web-port",
				HasPIP:         true,
				Subnets:        []string{"subnet1"},
				SecurityGroups: []string{"sg1"},
			},
			response:   `{"id": "port-new"}`,
			statusCode: http.StatusCreated,
			wantID:     "port-new",
			wantErr:    false,
		},
		{
			name:  "missing subnets",
			vpcID: "vpc1",
			request: PortCreateRequest{
				Name: "invalid-port",
			},
			response:   `{"error": "subnets required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/vpcs/%s/ports", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req PortCreateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.request.Name, req.Name)
				assertEqualSlice(t, tt.request.Subnets, req.Subnets)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			id, err := client.CreatePort(context.Background(), tt.vpcID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestVPCService_ListPublicIPs(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:  "successful list",
			vpcID: "vpc1",
			response: `{
				"public_ips": [
					{"id": "ip1", "public_ip": "203.0.113.1", "created_at": "2024-01-01T00:00:00"},
					{"id": "ip2", "public_ip": "203.0.113.2"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			vpcID:      "vpc1",
			response:   `{"public_ips": []}`,
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/vpcs/%s/public_ips", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			ips, err := client.ListPublicIPs(context.Background(), tt.vpcID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(ips))
		})
	}
}

func TestVPCService_CreatePublicIP(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		request    PublicIPCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:  "successful create",
			vpcID: "vpc1",
			request: PublicIPCreateRequest{
				Description: "Web server IP",
			},
			response:   `{"id": "ip-new"}`,
			statusCode: http.StatusCreated,
			wantID:     "ip-new",
			wantErr:    false,
		},
		{
			name:    "empty request",
			vpcID:   "vpc1",
			request: PublicIPCreateRequest{},
			response: `{
				"id": "ip-auto"
			}`,
			statusCode: http.StatusCreated,
			wantID:     "ip-auto",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/vpcs/%s/public_ips", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req PublicIPCreateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.request.Description, req.Description)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			id, err := client.CreatePublicIP(context.Background(), tt.vpcID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestVPCService_ListSubnets(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:  "successful list",
			vpcID: "vpc1",
			response: `{
				"subnets": [
					{"id": "subnet1", "name": "web-subnet"},
					{"id": "subnet2", "name": "db-subnet"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			vpcID:      "vpc1",
			response:   `{"subnets": []}`,
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/vpcs/%s/subnets", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			subnets, err := client.ListSubnets(context.Background(), tt.vpcID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(subnets))
		})
	}
}

func TestVPCService_CreateSubnet(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		request    SubnetCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:  "successful create",
			vpcID: "vpc1",
			request: SubnetCreateRequest{
				Name:      "web-subnet",
				CIDRBlock: "10.0.0.0/24",
				IPVersion: 4,
			},
			response:   `{"id": "subnet-new"}`,
			statusCode: http.StatusCreated,
			wantID:     "subnet-new",
			wantErr:    false,
		},
		{
			name:  "invalid CIDR",
			vpcID: "vpc1",
			request: SubnetCreateRequest{
				Name:      "invalid",
				CIDRBlock: "invalid",
			},
			response:   `{"error": "invalid CIDR"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/vpcs/%s/subnets", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req SubnetCreateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.request.Name, req.Name)
				assertEqual(t, tt.request.CIDRBlock, req.CIDRBlock)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testVPCClient(server.URL)
			id, err := client.CreateSubnet(context.Background(), tt.vpcID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func testVPCClient(baseURL string) VPCService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).VPCs()
}
