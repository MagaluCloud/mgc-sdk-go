package network

import (
	"context"
	"encoding/json"
	"io"
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

func TestVpcsPeeringsService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    VpcsPeeringsCreateRequest
		response   string
		statusCode int
		want       *VpcsPeeringsCreateResponse
		wantBody   string
		wantErr    bool
	}{
		{
			name: "peering between two vpcs is accepted",
			request: VpcsPeeringsCreateRequest{
				Name:        "peering-prod-to-db",
				Description: helpers.StrPtr("connection between production and database vpcs"),
				VPCs: VpcsPeeringsCreateVpcs{
					RequesterVpcID: "vpc-requester",
					AccepterVpcID:  "vpc-accepter",
				},
			},
			response: `{
				"id": "peering-1",
				"status": "pending"
			}`,
			statusCode: http.StatusAccepted,
			want: &VpcsPeeringsCreateResponse{
				ID:     "peering-1",
				Status: VpcsPeeringStatusPending,
			},
			wantBody: `{
				"name": "peering-prod-to-db",
				"description": "connection between production and database vpcs",
				"vpcs": {
					"requester_vpc_id": "vpc-requester",
					"accepter_vpc_id": "vpc-accepter"
				}
			}`,
			wantErr: false,
		},
		{
			name: "description is optional",
			request: VpcsPeeringsCreateRequest{
				Name: "peering-minimal",
				VPCs: VpcsPeeringsCreateVpcs{
					RequesterVpcID: "vpc-requester",
					AccepterVpcID:  "vpc-accepter",
				},
			},
			response: `{
				"id": "peering-2",
				"status": "pending"
			}`,
			statusCode: http.StatusAccepted,
			want: &VpcsPeeringsCreateResponse{
				ID:     "peering-2",
				Status: VpcsPeeringStatusPending,
			},
			wantBody: `{
				"name": "peering-minimal",
				"vpcs": {
					"requester_vpc_id": "vpc-requester",
					"accepter_vpc_id": "vpc-accepter"
				}
			}`,
			wantErr: false,
		},
		{
			name: "peering already exists",
			request: VpcsPeeringsCreateRequest{
				Name: "peering-duplicated",
				VPCs: VpcsPeeringsCreateVpcs{
					RequesterVpcID: "vpc-requester",
					AccepterVpcID:  "vpc-accepter",
				},
			},
			response:   `{"detail": "peering already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "server error",
			request: VpcsPeeringsCreateRequest{
				Name: "peering-prod-to-db",
				VPCs: VpcsPeeringsCreateVpcs{
					RequesterVpcID: "vpc-requester",
					AccepterVpcID:  "vpc-accepter",
				},
			},
			response:   `{"detail": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "invalid response body",
			request: VpcsPeeringsCreateRequest{
				Name: "peering-prod-to-db",
				VPCs: VpcsPeeringsCreateVpcs{
					RequesterVpcID: "vpc-requester",
					AccepterVpcID:  "vpc-accepter",
				},
			},
			response:   `{invalid json}`,
			statusCode: http.StatusAccepted,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/vpcs_peerings", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				if tt.wantBody != "" {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("failed to read request body: %v", err)
					}
					assertEqual(t, canonicalJSON(t, tt.wantBody), canonicalJSON(t, string(body)))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testPeeringClient(server.URL)
			got, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want.ID, got.ID)
			assertEqual(t, tt.want.Status, got.Status)
		})
	}
}

func TestVpcsPeeringsService_CreateValidation(t *testing.T) {
	tests := []struct {
		name    string
		request VpcsPeeringsCreateRequest
		err     string
	}{
		{
			name: "peering without name",
			request: VpcsPeeringsCreateRequest{
				VPCs: VpcsPeeringsCreateVpcs{
					RequesterVpcID: "vpc-requester",
					AccepterVpcID:  "vpc-accepter",
				},
			},
			err: "name cannot be empty",
		},
		{
			name: "peering without requester vpc",
			request: VpcsPeeringsCreateRequest{
				Name: "peering-prod-to-db",
				VPCs: VpcsPeeringsCreateVpcs{
					AccepterVpcID: "vpc-accepter",
				},
			},
			err: "requester_vpc_id cannot be empty",
		},
		{
			name: "peering without accepter vpc",
			request: VpcsPeeringsCreateRequest{
				Name: "peering-prod-to-db",
				VPCs: VpcsPeeringsCreateVpcs{
					RequesterVpcID: "vpc-requester",
				},
			},
			err: "accepter_vpc_id cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testPeeringClient("test")
			_, err := client.Create(context.Background(), tt.request)

			assertError(t, err)
			assertEqual(t, tt.err, err.Error())
		})
	}
}

func TestVpcsPeeringsService_List(t *testing.T) {
	createdAt := time.Date(2026, 7, 28, 10, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 7, 28, 11, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		opts       *ListVpcsPeeringsOptions
		wantVpcID  string
		response   string
		statusCode int
		check      func(t *testing.T, got *ListVpcsPeeringsResponse)
		wantErr    bool
	}{
		{
			name: "tenant peerings with their members",
			opts: &ListVpcsPeeringsOptions{},
			response: `{
				"meta": {
					"links": {
						"next": "?_offset=10&_limit=10",
						"previous": null,
						"self": "?_offset=0&_limit=10"
					},
					"page": {
						"count": 1,
						"limit": 10,
						"max_items_per_page": 100,
						"offset": 0,
						"total": 12
					}
				},
				"result": [
					{
						"vpc_peering_id": "peering-1",
						"name": "peering-prod-to-db",
						"description": "connection between production and database vpcs",
						"status": "created",
						"created_at": "` + createdAt.Format(utils.LocalDateTimeWithoutZoneLayout) + `",
						"updated": "` + updated.Format(utils.LocalDateTimeWithoutZoneLayout) + `",
						"members": [
							{
								"id": "member-1",
								"vpc_id": "vpc-requester",
								"direct_role": "requester"
							},
							{
								"id": "member-2",
								"vpc_id": "vpc-accepter",
								"direct_role": "accepter"
							}
						]
					}
				]
			}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *ListVpcsPeeringsResponse) {
				assertEqual(t, 12, got.Meta.Page.Total)
				assertEqual(t, 100, got.Meta.Page.MaxItemsPerPage)
				assertEqual(t, "?_offset=0&_limit=10", got.Meta.Links.Self)
				assertEqual(t, "?_offset=10&_limit=10", *got.Meta.Links.Next)

				assertEqual(t, 1, len(got.Result))
				peering := got.Result[0]
				assertEqual(t, "peering-1", peering.ID)
				assertEqual(t, "peering-prod-to-db", peering.Name)
				assertEqual(t, "connection between production and database vpcs", *peering.Description)
				assertEqual(t, VpcsPeeringStatusCreated, peering.Status)
				assertEqual(t, true, time.Time(*peering.CreatedAt).Equal(createdAt))
				assertEqual(t, true, time.Time(*peering.Updated).Equal(updated))

				assertEqual(t, 2, len(peering.Members))
				assertEqual(t, "member-1", peering.Members[0].ID)
				assertEqual(t, "vpc-requester", peering.Members[0].VpcID)
				assertEqual(t, VpcsPeeringDirectRoleRequester, peering.Members[0].DirectRole)
				assertEqual(t, VpcsPeeringDirectRoleAccepter, peering.Members[1].DirectRole)
			},
		},
		{
			name:      "peerings of a specific vpc",
			opts:      &ListVpcsPeeringsOptions{VpcID: "vpc-requester"},
			wantVpcID: "vpc-requester",
			response: `{
				"meta": {
					"links": {"self": "?_offset=0&_limit=10"},
					"page": {"count": 0, "limit": 10, "max_items_per_page": 100, "offset": 0, "total": 0}
				},
				"result": []
			}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *ListVpcsPeeringsResponse) {
				assertEqual(t, 0, len(got.Result))
				assertEqual(t, 0, got.Meta.Page.Total)
			},
		},
		{
			name: "no filter options",
			opts: nil,
			response: `{
				"meta": {
					"links": {"self": "?_offset=0&_limit=10"},
					"page": {"count": 0, "limit": 10, "max_items_per_page": 100, "offset": 0, "total": 0}
				},
				"result": []
			}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *ListVpcsPeeringsResponse) {
				assertEqual(t, 0, len(got.Result))
			},
		},
		{
			name:       "missing description becomes nil",
			opts:       &ListVpcsPeeringsOptions{},
			response:   `{"meta": {"links": {"self": "?"}, "page": {"count": 1, "limit": 10, "max_items_per_page": 100, "offset": 0, "total": 1}}, "result": [{"vpc_peering_id": "peering-2", "name": "peering-minimal", "status": "pending", "members": []}]}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *ListVpcsPeeringsResponse) {
				peering := got.Result[0]
				assertEqual(t, true, peering.Description == nil)
				assertEqual(t, true, peering.CreatedAt == nil)
				assertEqual(t, VpcsPeeringStatusPending, peering.Status)
			},
		},
		{
			name:       "server error",
			opts:       &ListVpcsPeeringsOptions{},
			response:   `{"detail": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "invalid response body",
			opts:       &ListVpcsPeeringsOptions{},
			response:   `{invalid json}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/vpcs_peerings", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				assertEqual(t, tt.wantVpcID, r.URL.Query().Get("vpc_id"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testPeeringClient(server.URL)
			got, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			tt.check(t, got)
		})
	}
}

func TestVpcsPeeringsService_GetMembers(t *testing.T) {
	tests := []struct {
		name       string
		peeringID  string
		response   string
		statusCode int
		check      func(t *testing.T, got *VpcsPeeringMembers)
		wantErr    bool
	}{
		{
			name:      "members of both ends of the peering",
			peeringID: "peering-1",
			response: `{
				"vpc_peering_id": "peering-1",
				"status": "created",
				"members": [
					{
						"id": "member-1",
						"vpc_id": "vpc-requester",
						"direct_role": "requester"
					},
					{
						"id": "member-2",
						"vpc_id": "vpc-accepter",
						"direct_role": "accepter"
					}
				]
			}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *VpcsPeeringMembers) {
				assertEqual(t, "peering-1", got.ID)
				assertEqual(t, VpcsPeeringStatusCreated, got.Status)
				assertEqual(t, 2, len(got.Members))
				assertEqual(t, "vpc-requester", got.Members[0].VpcID)
				assertEqual(t, VpcsPeeringDirectRoleRequester, got.Members[0].DirectRole)
				assertEqual(t, "vpc-accepter", got.Members[1].VpcID)
				assertEqual(t, VpcsPeeringDirectRoleAccepter, got.Members[1].DirectRole)
			},
		},
		{
			name:       "peering being created has no members yet",
			peeringID:  "peering-2",
			response:   `{"vpc_peering_id": "peering-2", "status": "pending", "members": []}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *VpcsPeeringMembers) {
				assertEqual(t, VpcsPeeringStatusPending, got.Status)
				assertEqual(t, 0, len(got.Members))
			},
		},
		{
			name:       "peering inexistente",
			peeringID:  "invalid",
			response:   `{"detail": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			peeringID:  "peering-1",
			response:   `{"detail": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "invalid response body",
			peeringID:  "peering-1",
			response:   `{invalid json}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/vpcs_peerings/"+tt.peeringID, r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))

			defer server.Close()

			client := testPeeringClient(server.URL)
			got, err := client.GetMembers(context.Background(), tt.peeringID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			tt.check(t, got)
		})
	}
}

func TestVpcsPeeringsService_GetMembersValidation(t *testing.T) {
	t.Parallel()

	client := testPeeringClient("test")
	_, err := client.GetMembers(context.Background(), "")

	assertError(t, err)
	assertEqual(t, "vpc_peering_id cannot be empty", err.Error())
}

func TestVpcsPeeringsService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		peeringID  string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			peeringID:  "peering-1",
			statusCode: http.StatusAccepted,
			wantErr:    false,
		},
		{
			name:       "non-existent peering id",
			peeringID:  "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "peering cannot be deleted in its current state",
			peeringID:  "peering-1",
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name:       "server error",
			peeringID:  "peering-1",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v1/vpcs_peerings/"+tt.peeringID, r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
			}))

			defer server.Close()

			client := testPeeringClient(server.URL)
			err := client.Delete(context.Background(), tt.peeringID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))

				return
			}

			assertNoError(t, err)
		})
	}
}

func TestVpcsPeeringsService_DeleteValidation(t *testing.T) {
	t.Parallel()

	client := testPeeringClient("test")
	err := client.Delete(context.Background(), "")

	assertError(t, err)
	assertEqual(t, "vpc_peering_id cannot be empty", err.Error())
}

// canonicalJSON normalizes a JSON document so payloads can be compared
// regardless of key order and formatting.
func canonicalJSON(t *testing.T, raw string) string {
	t.Helper()

	var decoded any
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		t.Fatalf("invalid JSON %q: %v", raw, err)
	}

	encoded, err := json.Marshal(decoded)
	if err != nil {
		t.Fatalf("failed to encode JSON: %v", err)
	}

	return string(encoded)
}

func testPeeringClient(baseURL string) VpcsPeeringsService {
	httpClient := &http.Client{}

	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))

	return New(core).VpcsPeerings()
}
