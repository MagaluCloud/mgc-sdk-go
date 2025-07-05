package network

import (
	"context"
	"encoding/json"
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

func TestNatGatewayService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateNatGatewayRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful create",
			request: CreateNatGatewayRequest{
				Name:        "prod-nat",
				Description: helpers.StrPtr("Production NAT Gateway"),
				Zone:        "zone1",
				VPCID:       "vpc1",
			},
			response:   `{"id": "nat1", "status": "creating"}`,
			statusCode: http.StatusCreated,
			wantID:     "nat1",
			wantErr:    false,
		},
		{
			name: "missing required fields",
			request: CreateNatGatewayRequest{
				Description: helpers.StrPtr("Invalid NAT Gateway"),
			},
			response:   `{"error": "name, zone and vpc_id are required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/nat_gateways", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req CreateNatGatewayRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.request.Name, req.Name)
				assertEqual(t, tt.request.Zone, req.Zone)
				assertEqual(t, tt.request.VPCID, req.VPCID)
				if tt.request.Description != nil {
					assertEqual(t, *tt.request.Description, *req.Description)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testNatGatewayClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestNatGatewayService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "nat1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent nat gateway",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "nat gateway not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/nat_gateways/"+tt.id, r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)

				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testNatGatewayClient(server.URL)
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

func TestNatGatewayService_Get(t *testing.T) {
	basetime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	basetimeWithoutZone := utils.LocalDateTimeWithoutZone(basetime)
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		want       *NatGatewayDetailsResponse
		wantErr    bool
	}{
		{
			name: "successful get",
			id:   "nat1",
			response: `{
				"id": "nat1",
				"name": "prod-nat",
				"description": "Production NAT Gateway",
				"vpc_id": "vpc1",
				"zone": "zone1",
				"nat_gateway_ip": "10.0.0.1",
				"status": "active",
				"created_at": "` + basetime.Format(utils.LocalDateTimeWithoutZoneLayout) + `",
				"updated": "` + basetime.Format(utils.LocalDateTimeWithoutZoneLayout) + `"
			}`,
			statusCode: http.StatusOK,
			want: &NatGatewayDetailsResponse{
				ID:           helpers.StrPtr("nat1"),
				Name:         helpers.StrPtr("prod-nat"),
				Description:  helpers.StrPtr("Production NAT Gateway"),
				VPCID:        helpers.StrPtr("vpc1"),
				Zone:         helpers.StrPtr("zone1"),
				NatGatewayIP: helpers.StrPtr("10.0.0.1"),
				Status:       "active",
				CreatedAt:    &basetimeWithoutZone,
				Updated:      &basetimeWithoutZone,
			},
			wantErr: false,
		},
		{
			name:       "non-existent nat gateway",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "nat gateway not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/nat_gateways/"+tt.id, r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testNatGatewayClient(server.URL)
			got, err := client.Get(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, *tt.want.ID, *got.ID)
			assertEqual(t, *tt.want.Name, *got.Name)
			assertEqual(t, *tt.want.Description, *got.Description)
			assertEqual(t, *tt.want.VPCID, *got.VPCID)
			assertEqual(t, *tt.want.Zone, *got.Zone)
			assertEqual(t, *tt.want.NatGatewayIP, *got.NatGatewayIP)
			assertEqual(t, tt.want.Status, got.Status)
		})
	}
}

func TestNatGatewayService_List(t *testing.T) {
	basetime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	basetimeWithoutZone := utils.LocalDateTimeWithoutZone(basetime)
	tests := []struct {
		name       string
		vpcID      string
		opts       ListOptions
		response   string
		statusCode int
		want       []NatGatewayResponse
		wantErr    bool
	}{
		{
			name:  "successful list",
			vpcID: "vpc1",
			opts: ListOptions{
				Limit:  helpers.IntPtr(10),
				Offset: helpers.IntPtr(0),
			},
			response: `{
				"meta": {
					"page": {
						"limit": 10,
						"offset": 0,
						"count": 1,
						"total": 1,
						"max_items_per_page": 100
					},
					"links": {
						"self": "/v1/nat_gateways?page=1&items_per_page=10"
					}
				},
				"result": [{
					"id": "nat1",
					"name": "prod-nat",
					"description": "Production NAT Gateway",
					"vpc_id": "vpc1",
					"zone": "zone1",
					"nat_gateway_ip": "10.0.0.1",
					"status": "active",
					"created_at": "` + basetime.Format(utils.LocalDateTimeWithoutZoneLayout) + `",
					"updated": "` + basetime.Format(utils.LocalDateTimeWithoutZoneLayout) + `"
				}]
			}`,
			statusCode: http.StatusOK,
			want: []NatGatewayResponse{
				{
					ID:           helpers.StrPtr("nat1"),
					Name:         helpers.StrPtr("prod-nat"),
					Description:  helpers.StrPtr("Production NAT Gateway"),
					VPCID:        helpers.StrPtr("vpc1"),
					Zone:         helpers.StrPtr("zone1"),
					NatGatewayIP: helpers.StrPtr("10.0.0.1"),
					Status:       "active",
					CreatedAt:    &basetimeWithoutZone,
					Updated:      &basetimeWithoutZone,
				},
			},
			wantErr: false,
		},
		{
			name:  "empty list",
			vpcID: "vpc2",
			opts: ListOptions{
				Limit:  helpers.IntPtr(10),
				Offset: helpers.IntPtr(0),
			},
			response: `{
				"meta": {
					"page": {
						"limit": 10,
						"offset": 0,
						"count": 0,
						"total": 0,
						"max_items_per_page": 100
					},
					"links": {
						"self": "/v1/nat_gateways?page=1&items_per_page=10"
					}
				},
				"result": []
			}`,
			statusCode: http.StatusOK,
			want:       []NatGatewayResponse{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/nat_gateways", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				assertEqual(t, tt.vpcID, r.URL.Query().Get("vpc_id"))
				assertEqual(t, "1", r.URL.Query().Get("page"))
				assertEqual(t, "10", r.URL.Query().Get("items_per_page"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testNatGatewayClient(server.URL)
			got, err := client.List(context.Background(), tt.vpcID, tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, len(tt.want), len(got))
			if len(got) > 0 {
				assertEqual(t, *tt.want[0].ID, *got[0].ID)
				assertEqual(t, *tt.want[0].Name, *got[0].Name)
				assertEqual(t, *tt.want[0].Description, *got[0].Description)
				assertEqual(t, *tt.want[0].VPCID, *got[0].VPCID)
				assertEqual(t, *tt.want[0].Zone, *got[0].Zone)
				assertEqual(t, *tt.want[0].NatGatewayIP, *got[0].NatGatewayIP)
				assertEqual(t, tt.want[0].Status, got[0].Status)
			}
		})
	}
}

func testNatGatewayClient(baseURL string) NatGatewayService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NatGateways()
}
