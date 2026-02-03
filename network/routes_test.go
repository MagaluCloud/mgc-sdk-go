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
		opts       *ListRouteOptions
		response   string
		statusCode int
		want       *ListResponse
		wantErr    bool
	}{
		{
			name:  "successful list",
			vpcID: "vpc-1",
			opts:  &ListRouteOptions{},
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
						Next:     helpers.StrPtr("?_offset=4&_limit=2"),
						Previous: helpers.StrPtr("?_offset=0&_limit=2"),
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
				Result: []RouteDetail{
					{
						ID:              "route-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
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
			opts: &ListRouteOptions{
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
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
						Next:     helpers.StrPtr("?_offset=4&_limit=2"),
						Previous: helpers.StrPtr("?_offset=0&_limit=2"),
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
				Result: []RouteDetail{
					{
						ID:              "route-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
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
			opts: &ListRouteOptions{
				Zone: "a",
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
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
						Next:     helpers.StrPtr("?_offset=4&_limit=2"),
						Previous: helpers.StrPtr("?_offset=0&_limit=2"),
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
				Result: []RouteDetail{
					{
						ID:              "route-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
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
			opts: &ListRouteOptions{
				Sort: "description:desc",
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
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
						Next:     helpers.StrPtr("?_offset=4&_limit=2"),
						Previous: helpers.StrPtr("?_offset=0&_limit=2"),
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
				Result: []RouteDetail{
					{
						ID:              "route-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
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
			name:  "successful list without meta links next and previous",
			vpcID: "vpc-1",
			opts:  &ListRouteOptions{},
			response: `{
				"meta": {
					"links": {
						"next": null,
						"previous": null,
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
						Next:     nil,
						Previous: nil,
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
				Result: []RouteDetail{
					{
						ID:              "route-1",
						PortID:          "port-1",
						CIDRDestination: "192.168.1.1",
						Description:     "Description",
						NextHop:         "192.168.1.1",
						Type:            "default",
						Status:          "processing",
					},
					{
						ID:              "route-2",
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
			opts:       &ListRouteOptions{},
			response:   `{"error": "vpc id not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			vpcID:      "vpc-1",
			opts:       &ListRouteOptions{},
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
				if tt.opts.Zone != "" {
					assertEqual(t, tt.opts.Zone, query.Get("zone"))
				}
				if tt.opts.Sort != "" {
					assertEqual(t, tt.opts.Sort, query.Get("sort"))
				}
				if tt.opts.Page != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Page), query.Get("page"))
				}
				if tt.opts.ItemsPerPage != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.ItemsPerPage), query.Get("items_per_page"))
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

			if tt.want.Meta.Links.Next == nil {
				assertEqual(t, tt.want.Meta.Links.Next, routes.Meta.Links.Next)
			} else {
				assertEqual(t, *tt.want.Meta.Links.Next, *routes.Meta.Links.Next)
			}
			if tt.want.Meta.Links.Previous == nil {
				assertEqual(t, tt.want.Meta.Links.Previous, routes.Meta.Links.Previous)
			} else {
				assertEqual(t, *tt.want.Meta.Links.Previous, *routes.Meta.Links.Previous)
			}
			assertEqual(t, tt.want.Meta.Links.Self, routes.Meta.Links.Self)

			assertEqual(t, tt.want.Meta.Page.Count, routes.Meta.Page.Count)
			assertEqual(t, tt.want.Meta.Page.Limit, routes.Meta.Page.Limit)
			assertEqual(t, tt.want.Meta.Page.MaxItemsPerPage, routes.Meta.Page.MaxItemsPerPage)
			assertEqual(t, tt.want.Meta.Page.Offset, routes.Meta.Page.Offset)
			assertEqual(t, tt.want.Meta.Page.Total, routes.Meta.Page.Total)

			for i, route := range routes.Result {
				assertEqual(t, tt.want.Result[i].ID, route.ID)
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

func TestRouteService_List_InvalidSortValue(t *testing.T) {
	tests := []struct {
		name string
		sort string
		err  string
	}{
		{
			name: "invalid format",
			sort: "test",
			err:  "invalid sort format, expected field:asc|desc",
		},
		{
			name: "invalid field",
			sort: "test:asc",
			err:  `invalid sort field: "test", allowed fields are: id, port_id, vpc_id, description, cidr_destination, type, status`,
		},
		{
			name: "invalid sort direction",
			sort: "description:test",
			err:  "invalid sort direction, expected asc or desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testRouteClient("test")
			_, err := client.List(context.Background(), "123", &ListRouteOptions{
				Sort: tt.sort,
			})

			assertEqual(t, tt.err, err.Error())
		})
	}
}

func TestRouteService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		vpcID      string
		opts       *ListAllRoutesOptions
		response   string
		statusCode int
		want       []RouteDetail
		wantErr    bool
	}{
		{
			name:  "successful list all",
			vpcID: "vpc-1",
			opts:  &ListAllRoutesOptions{},
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
						"total": 200
					}
				},
				"result": [
					{
						"id": "route-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
			want: []RouteDetail{
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
			},
			wantErr: false,
		},
		{
			name:  "successful list all with 2 pages",
			vpcID: "vpc-1",
			opts:  &ListAllRoutesOptions{},
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
						"total": 199
					}
				},
				"result": [
					{
						"id": "route-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
			want: []RouteDetail{
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
			},
			wantErr: false,
		},
		{
			name:  "successful list all with 3 pages",
			vpcID: "vpc-1",
			opts:  &ListAllRoutesOptions{},
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
						"total": 201
					}
				},
				"result": [
					{
						"id": "route-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
			want: []RouteDetail{
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
			},
			wantErr: false,
		},
		{
			name:  "successful list with zone",
			vpcID: "vpc-1",
			opts: &ListAllRoutesOptions{
				Zone: "a",
			},
			response: `{
				"meta": {
					"links": {
						"next": "?_offset=4&_limit=2",
						"previous": "?_offset=0&_limit=2",
						"self": "?_offset=0&_limit=2"
					},
					"page": {
						"count": 2,
						"limit": 2,
						"max_items_per_page": 100,
						"offset": 0,
						"total": 2
					}
				},
				"result": [
					{
						"id": "route-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
			want: []RouteDetail{
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
			},
			wantErr: false,
		},
		{
			name:  "successful list with sort",
			vpcID: "vpc-1",
			opts: &ListAllRoutesOptions{
				Sort: "description:asc",
			},
			response: `{
				"meta": {
					"links": {
						"next": "?_offset=4&_limit=2",
						"previous": "?_offset=0&_limit=2",
						"self": "?_offset=0&_limit=2"
					},
					"page": {
						"count": 2,
						"limit": 2,
						"max_items_per_page": 100,
						"offset": 0,
						"total": 2
					}
				},
				"result": [
					{
						"id": "route-1",
						"port_id": "port-1",
						"cidr_destination": "192.168.1.1",
						"description": "Description",
						"next_hop": "192.168.1.1",
						"type": "default",
						"status": "processing"
					},
					{
						"id": "route-2",
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
			want: []RouteDetail{
				{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				{
					ID:              "route-2",
					PortID:          "port-2",
					CIDRDestination: "192.168.2.1",
					Description:     "Description",
					NextHop:         "192.168.2.1",
					Type:            "default",
					Status:          "processing",
				},
			},
			wantErr: false,
		},
		{
			name:       "non-existent vpc id",
			vpcID:      "invalid",
			opts:       &ListAllRoutesOptions{},
			response:   `{"error": "vpc id not found"}`,
			statusCode: http.StatusNotFound,
			want:       []RouteDetail{},
			wantErr:    true,
		},
		{
			name:       "server error",
			vpcID:      "vpc-1",
			opts:       &ListAllRoutesOptions{},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			want:       []RouteDetail{},
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
				if tt.opts.Zone != "" {
					assertEqual(t, tt.opts.Zone, query.Get("zone"))
				}
				if tt.opts.Sort != "" {
					assertEqual(t, tt.opts.Sort, query.Get("sort"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testRouteClient(server.URL)
			routes, err := client.ListAll(context.Background(), tt.vpcID, tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))

				return
			}

			assertNoError(t, err)

			assertEqual(t, len(tt.want), len(routes))

			for i, route := range routes {
				assertEqual(t, tt.want[i].ID, route.ID)
				assertEqual(t, tt.want[i].PortID, route.PortID)
				assertEqual(t, tt.want[i].CIDRDestination, route.CIDRDestination)
				assertEqual(t, tt.want[i].Description, route.Description)
				assertEqual(t, tt.want[i].NextHop, route.NextHop)
				assertEqual(t, tt.want[i].Type, route.Type)
				assertEqual(t, tt.want[i].Status, route.Status)
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
				RouteDetail: RouteDetail{
					ID:              "route-1",
					PortID:          "port-1",
					CIDRDestination: "192.168.1.1",
					Description:     "Description",
					NextHop:         "192.168.1.1",
					Type:            "default",
					Status:          "processing",
				},
				VpcID: "vpc-1",
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

func TestRouteService_Create_InvalidBody(t *testing.T) {
	tests := []struct {
		name string
		body CreateRequest
		err  string
	}{
		{
			name: "empty port_id",
			body: CreateRequest{
				CIDRDestination: "192.168.1.1",
			},
			err: "port_id cannot be empty",
		},
		{
			name: "empty cidr_destination",
			body: CreateRequest{
				PortID: "port-1",
			},
			err: "cidr_destination cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testRouteClient("test")
			_, err := client.Create(context.Background(), "123", tt.body)

			assertEqual(t, tt.err, err.Error())
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

func TestValidateSortValue_Valid(t *testing.T) {
	validSorts := []string{
		"id:asc",
		"port_id:desc",
		"description:asc",
		"CIDR_DESTINATION:DESC",
	}

	for _, sort := range validSorts {
		t.Run(sort, func(t *testing.T) {
			t.Parallel()

			if err := validateSortValue(sort); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
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
