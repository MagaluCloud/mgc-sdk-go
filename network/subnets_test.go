package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

func TestSubnetService_Get(t *testing.T) {
	timeparsed, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	localDatew := utils.LocalDateTimeWithoutZone(timeparsed)

	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		want       *SubnetResponseDetail
		wantErr    bool
	}{
		{
			name: "successful get",
			id:   "subnet1",
			response: `{
				"id": "subnet1",
				"name": "prod-subnet",
				"vpc_id": "vpc1",
				"cidr_block": "10.0.0.0/24",
				"gateway_ip": "10.0.0.1",
				"dns_nameservers": ["8.8.8.8"],
				"dhcp_pools": [{"start": "10.0.0.100", "end": "10.0.0.200"}],
				"created_at": "2024-01-01T00:00:00.000000",
				"updated": "2024-01-01T00:00:00.000000"
			}`,
			statusCode: http.StatusOK,
			want: &SubnetResponseDetail{
				SubnetResponse: SubnetResponse{
					ID:        "subnet1",
					Name:      helpers.StrPtr("prod-subnet"),
					VPCID:     "vpc1",
					CIDRBlock: "10.0.0.0/24",
					CreatedAt: &localDatew,
					Updated:   &localDatew,
				},
				GatewayIP:      "10.0.0.1",
				DNSNameservers: []string{"8.8.8.8"},
				DHCPPools: []DHCPPoolStr{
					{Start: "10.0.0.100", End: "10.0.0.200"},
				},
			},
			wantErr: false,
		},
		{
			name:       "non-existent subnet",
			id:         "invalid",
			response:   `{"error": "subnet not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			id:         "subnet1",
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
				assertEqual(t, fmt.Sprintf("/network/v0/subnets/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSubnetClient(server.URL)
			subnet, err := client.Get(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want.ID, subnet.ID)
			assertEqual(t, *tt.want.Name, *subnet.Name)
			assertEqual(t, tt.want.GatewayIP, subnet.GatewayIP)
			assertEqual(t, len(tt.want.DNSNameservers), len(subnet.DNSNameservers))
			assertEqual(t, len(tt.want.DHCPPools), len(subnet.DHCPPools))
		})
	}
}

func TestSubnetService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "subnet1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent subnet",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "subnet not found"}`,
			wantErr:    true,
		},
		{
			name:       "conflict error",
			id:         "subnet-in-use",
			statusCode: http.StatusConflict,
			response:   `{"error": "subnet in use"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/subnets/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testSubnetClient(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestSubnetService_Update(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    SubnetPatchRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful update dns",
			id:   "subnet1",
			request: SubnetPatchRequest{
				DNSNameservers: &[]string{"8.8.8.8", "1.1.1.1"},
			},
			response:   `{"id": "subnet1"}`,
			statusCode: http.StatusOK,
			wantID:     "subnet1",
			wantErr:    false,
		},
		{
			name: "empty dns servers",
			id:   "subnet1",
			request: SubnetPatchRequest{
				DNSNameservers: &[]string{},
			},
			response:   `{"id": "subnet1"}`,
			statusCode: http.StatusOK,
			wantID:     "subnet1",
			wantErr:    false,
		},
		{
			name:    "invalid request",
			id:      "subnet1",
			request: SubnetPatchRequest{},
			response: `{
				"error": "at least one field must be provided"
			}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "non-existent subnet",
			id:   "invalid",
			request: SubnetPatchRequest{
				DNSNameservers: &[]string{"8.8.8.8"},
			},
			response:   `{"error": "subnet not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/subnets/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPatch, r.Method)

				if !tt.wantErr {
					var req SubnetPatchRequest
					err := json.NewDecoder(r.Body).Decode(&req)
					assertNoError(t, err)

					assertEqualSlice(t, *tt.request.DNSNameservers, *req.DNSNameservers)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSubnetClient(server.URL)
			resp, err := client.Update(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, resp.ID)
		})
	}
}

func testSubnetClient(baseURL string) SubnetService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Subnets()
}

func assertEqualSlice(t *testing.T, want, got []string) {
	t.Helper()
	assertEqual(t, len(want), len(got))
	for i := range want {
		assertEqual(t, want[i], got[i])
	}
}
