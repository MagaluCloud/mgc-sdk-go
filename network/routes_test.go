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

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestRouteService_List(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		opts       ListRouteOptions
		response   string
		statusCode int
		want       *ListResponse
		wantErr    bool
	}{
		{
			name:  "successful list",
			vpcID: "vpc-1",
			opts:  ListRouteOptions{},
			response: `{
				"meta": {
					"links": {
						"next": "?_offset=4&_limit=2",
						"previous": "?_offset=0&_limit=2",
						"self": "?_offset=0&_limit=2"
					},
					"page": {
						"count": 6,
						"limit": 2,
						"max_items_per_page": 100,
						"offset": 0,
						"total": 6
					}
				},
				"result": [
					{
						"id": "route-1",
						"vpc_id": "vpc-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
						"vpc_id": "vpc-2",
						"port_id": "port-2",
						"cidr_destination": "192.168.2.1",
						"description": "Description",
						"next_hop": "192.168.2.1",
						"type": "default",
						"status": "processing"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListResponse{
				Meta: ListMeta{
					Links: ListLinks{
						Next:     "?_offset=4&_limit=2",
						Previous: "?_offset=0&_limit=2",
						Self:     "?_offset=0&_limit=2",
					},
					Page: ListPage{
						Count:           6,
						Limit:           2,
						MaxItemsPerPage: 100,
						Offset:          0,
						Total:           6,
					},
				},
				Result: []Route{
					{
						ID:              "route-1",
						VpcID:           "vpc-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
						VpcID:           "vpc-2",
						PortID:          "port-2",
						CIDRDestination: "192.168.2.1",
						Description:     "Description",
						NextHop:         "192.168.2.1",
						Type:            "default",
						Status:          "processing",
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "successful list with pagination",
			vpcID: "vpc-1",
			opts: ListRouteOptions{
				Page:         helpers.IntPtr(2),
				ItemsPerPage: helpers.IntPtr(2),
			},
			response: `{
				"meta": {
					"links": {
						"next": "?_offset=4&_limit=2",
						"previous": "?_offset=0&_limit=2",
						"self": "?_offset=2&_limit=2"
					},
					"page": {
						"count": 6,
						"limit": 2,
						"max_items_per_page": 100,
						"offset": 2,
						"total": 6
					}
				},
				"result": [
					{
						"id": "route-1",
						"vpc_id": "vpc-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
						"vpc_id": "vpc-2",
						"port_id": "port-2",
						"cidr_destination": "192.168.2.1",
						"description": "Description",
						"next_hop": "192.168.2.1",
						"type": "default",
						"status": "processing"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListResponse{
				Meta: ListMeta{
					Links: ListLinks{
						Next:     "?_offset=4&_limit=2",
						Previous: "?_offset=0&_limit=2",
						Self:     "?_offset=2&_limit=2",
					},
					Page: ListPage{
						Count:           6,
						Limit:           2,
						MaxItemsPerPage: 100,
						Offset:          2,
						Total:           6,
					},
				},
				Result: []Route{
					{
						ID:              "route-1",
						VpcID:           "vpc-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
						VpcID:           "vpc-2",
						PortID:          "port-2",
						CIDRDestination: "192.168.2.1",
						Description:     "Description",
						NextHop:         "192.168.2.1",
						Type:            "default",
						Status:          "processing",
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "successful list with zone",
			vpcID: "vpc-1",
			opts: ListRouteOptions{
				Zone: helpers.StrPtr("br-se1"),
			},
			response: `{
				"meta": {
					"links": {
						"next": "?_offset=4&_limit=2",
						"previous": "?_offset=0&_limit=2",
						"self": "?_offset=0&_limit=2"
					},
					"page": {
						"count": 6,
						"limit": 2,
						"max_items_per_page": 100,
						"offset": 0,
						"total": 6
					}
				},
				"result": [
					{
						"id": "route-1",
						"vpc_id": "vpc-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
						"vpc_id": "vpc-2",
						"port_id": "port-2",
						"cidr_destination": "192.168.2.1",
						"description": "Description",
						"next_hop": "192.168.2.1",
						"type": "default",
						"status": "processing"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListResponse{
				Meta: ListMeta{
					Links: ListLinks{
						Next:     "?_offset=4&_limit=2",
						Previous: "?_offset=0&_limit=2",
						Self:     "?_offset=0&_limit=2",
					},
					Page: ListPage{
						Count:           6,
						Limit:           2,
						MaxItemsPerPage: 100,
						Offset:          0,
						Total:           6,
					},
				},
				Result: []Route{
					{
						ID:              "route-1",
						VpcID:           "vpc-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
						VpcID:           "vpc-2",
						PortID:          "port-2",
						CIDRDestination: "192.168.2.1",
						Description:     "Description",
						NextHop:         "192.168.2.1",
						Type:            "default",
						Status:          "processing",
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "successful list with sort",
			vpcID: "vpc-1",
			opts: ListRouteOptions{
				Sort: helpers.StrPtr("desc"),
			},
			response: `{
				"meta": {
					"links": {
						"next": "?_offset=4&_limit=2",
						"previous": "?_offset=0&_limit=2",
						"self": "?_offset=0&_limit=2"
					},
					"page": {
						"count": 6,
						"limit": 2,
						"max_items_per_page": 100,
						"offset": 0,
						"total": 6
					}
				},
				"result": [
					{
						"id": "route-1",
						"vpc_id": "vpc-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
						"vpc_id": "vpc-2",
						"port_id": "port-2",
						"cidr_destination": "192.168.2.1",
						"description": "Description",
						"next_hop": "192.168.2.1",
						"type": "default",
						"status": "processing"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want: &ListResponse{
				Meta: ListMeta{
					Links: ListLinks{
						Next:     "?_offset=4&_limit=2",
						Previous: "?_offset=0&_limit=2",
						Self:     "?_offset=0&_limit=2",
					},
					Page: ListPage{
						Count:           6,
						Limit:           2,
						MaxItemsPerPage: 100,
						Offset:          0,
						Total:           6,
					},
				},
				Result: []Route{
					{
						ID:              "route-1",
						VpcID:           "vpc-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
						VpcID:           "vpc-2",
						PortID:          "port-2",
						CIDRDestination: "192.168.2.1",
						Description:     "Description",
						NextHop:         "192.168.2.1",
						Type:            "default",
						Status:          "processing",
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "non-existent vpc id",
			vpcID:      "invalid",
			opts:       ListRouteOptions{},
			response:   `{"error": "vpc id not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			vpcID:      "vpc-1",
			opts:       ListRouteOptions{},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v1/vpcs/%s/route_table/routes", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				query := r.URL.Query()
				if tt.opts.Zone != nil {
					assertEqual(t, *tt.opts.Zone, query.Get("_zone"))
				}
				if tt.opts.Sort != nil {
					assertEqual(t, *tt.opts.Sort, query.Get("_sort"))
				}
				if tt.opts.Page != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Page), query.Get("_page"))
				}
				if tt.opts.ItemsPerPage != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.ItemsPerPage), query.Get("_items_per_page"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testRouteClient(server.URL)
			routes, err := client.List(context.Background(), tt.vpcID, tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))

				return
			}

			assertNoError(t, err)

			assertEqual(t, tt.want.Meta.Links.Next, routes.Meta.Links.Next)
			assertEqual(t, tt.want.Meta.Links.Previous, routes.Meta.Links.Previous)
			assertEqual(t, tt.want.Meta.Links.Self, routes.Meta.Links.Self)

			assertEqual(t, tt.want.Meta.Page.Count, routes.Meta.Page.Count)
			assertEqual(t, tt.want.Meta.Page.Limit, routes.Meta.Page.Limit)
			assertEqual(t, tt.want.Meta.Page.MaxItemsPerPage, routes.Meta.Page.MaxItemsPerPage)
			assertEqual(t, tt.want.Meta.Page.Offset, routes.Meta.Page.Offset)
			assertEqual(t, tt.want.Meta.Page.Total, routes.Meta.Page.Total)

			for i, route := range routes.Result {
				assertEqual(t, tt.want.Result[i].ID, route.ID)
				assertEqual(t, tt.want.Result[i].VpcID, route.VpcID)
				assertEqual(t, tt.want.Result[i].PortID, route.PortID)
				assertEqual(t, tt.want.Result[i].CIDRDestination, route.CIDRDestination)
				assertEqual(t, tt.want.Result[i].Description, route.Description)
				assertEqual(t, tt.want.Result[i].NextHop, route.NextHop)
				assertEqual(t, tt.want.Result[i].Type, route.Type)
				assertEqual(t, tt.want.Result[i].Status, route.Status)
			}
		})
	}
}

func TestRouteService_Get(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		routeID    string
		response   string
		statusCode int
		want       *Route
		wantErr    bool
	}{
		{
			name:    "successful get",
			vpcID:   "vpc-1",
			routeID: "route-1",
			response: `{
				"id": "route-1",
				"vpc_id": "vpc-1",
				"port_id": "port-1",
				"cidr_destination": "192.168.1.1",
				"description": "Description",
				"next_hop": "192.168.1.1",
				"type": "default",
				"status": "processing"
			}`,
			statusCode: http.StatusOK,
			want: &Route{
				ID:              "route-1",
				VpcID:           "vpc-1",
				PortID:          "port-1",
				CIDRDestination: "192.168.1.1",
				Description:     "Description",
				NextHop:         "192.168.1.1",
				Type:            "default",
				Status:          "processing",
			},
			wantErr: false,
		},
		{
			name:       "non-existent vpc id",
			vpcID:      "invalid",
			routeID:    "route-1",
			response:   `{"error": "vpc id not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "non-existent route id",
			vpcID:      "vpc-1",
			routeID:    "invalid",
			response:   `{"error": "route id not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			vpcID:      "vpc-1",
			routeID:    "route-1",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v1/vpcs/%s/route_table/%s", tt.vpcID, tt.routeID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testRouteClient(server.URL)
			route, err := client.Get(context.Background(), tt.vpcID, tt.routeID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))

				return
			}

			assertNoError(t, err)

			assertEqual(t, tt.want.ID, route.ID)
			assertEqual(t, tt.want.VpcID, route.VpcID)
			assertEqual(t, tt.want.PortID, route.PortID)
			assertEqual(t, tt.want.CIDRDestination, route.CIDRDestination)
			assertEqual(t, tt.want.Description, route.Description)
			assertEqual(t, tt.want.NextHop, route.NextHop)
			assertEqual(t, tt.want.Type, route.Type)
			assertEqual(t, tt.want.Status, route.Status)
		})
	}
}

func TestRouteService_Create(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		body       CreateRequest
		response   string
		statusCode int
		want       *CreateResponse
		wantErr    bool
	}{
		{
			name:  "successful create",
			vpcID: "vpc-1",
			body: CreateRequest{
				PortID:          "port-1",
				CIDRDestination: "192.168.1.1",
			},
			response: `{
				"id": "route-1",
				"status": "processing"
			}`,
			statusCode: http.StatusOK,
			want: &CreateResponse{
				ID:     "route-1",
				Status: "processing",
			},
			wantErr: false,
		},
		{
			name:  "successful create with description",
			vpcID: "vpc-1",
			body: CreateRequest{
				PortID:          "port-1",
				CIDRDestination: "192.168.1.1",
				Description:     helpers.StrPtr("Description"),
			},
			response: `{
				"id": "route-1",
				"status": "processing"
			}`,
			statusCode: http.StatusOK,
			want: &CreateResponse{
				ID:     "route-1",
				Status: "processing",
			},
			wantErr: false,
		},
		{
			name:  "non-existent vpc id",
			vpcID: "invalid",
			body: CreateRequest{
				PortID:          "port-1",
				CIDRDestination: "192.168.1.1",
			},
			response:   `{"error": "vpc id not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:  "server error",
			vpcID: "vpc-1",
			body: CreateRequest{
				PortID:          "port-1",
				CIDRDestination: "192.168.1.1",
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v1/vpcs/%s/route_table/routes", tt.vpcID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req CreateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)

				assertEqual(t, tt.body.CIDRDestination, req.CIDRDestination)
				assertEqual(t, tt.body.PortID, req.PortID)

				if tt.body.Description != nil {
					assertEqual(t, *tt.body.Description, *req.Description)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testRouteClient(server.URL)
			route, err := client.Create(context.Background(), tt.vpcID, tt.body)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))

				return
			}

			assertNoError(t, err)

			assertEqual(t, tt.want.ID, route.ID)
			assertEqual(t, tt.want.Status, route.Status)
		})
	}
}

func TestRouteService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		routeID    string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			vpcID:      "vpc-1",
			routeID:    "route-1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent vpc id",
			vpcID:      "invalid",
			routeID:    "route-1",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "non-existent route id",
			vpcID:      "vpc-1",
			routeID:    "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			vpcID:      "vpc-1",
			routeID:    "route-1",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v1/vpcs/%s/route_table/%s", tt.vpcID, tt.routeID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
			}))

			defer server.Close()

			client := testRouteClient(server.URL)
			err := client.Delete(context.Background(), tt.vpcID, tt.routeID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))

				return
			}

			assertNoError(t, err)
		})
	}
}

func testRouteClient(baseURL string) RouteService {
	httpClient := &http.Client{}

	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))

	return New(core).Routes()
}
