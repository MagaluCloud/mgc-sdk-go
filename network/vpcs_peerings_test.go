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
	createdAt := time.Date(2026, 7, 23, 20, 21, 8, 861443000, time.UTC)
	updated := time.Date(2026, 7, 24, 17, 47, 13, 459860000, time.UTC)

	tests := []struct {
		name       string
		opts       *ListVpcsPeeringsOptions
		wantQuery  map[string]string
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
						"next": null,
						"previous": null,
						"self": "?_offset=0&_limit=1"
					},
					"page": {
						"count": 1,
						"limit": 1,
						"max_items_per_page": 100,
						"offset": 0,
						"total": 1
					}
				},
				"result": [
					{
						"created_at": "2026-07-23T20:21:08.861443",
						"description": "desc",
						"id": "866b0e4d-59d7-42e6-82d7-aa5fcb72360d",
						"members": [
							{
								"direct_role": "requester",
								"id": "ec32dc9f-679d-4c25-9f3f-b0b37aa101a0",
								"vpc_id": "2ea7d752-c185-422b-9780-2d07109d8afe"
							},
							{
								"direct_role": "accepter",
								"id": "759afe65-f9d5-408e-a4ac-c430dbbd6d35",
								"vpc_id": "732d0b32-b9c0-44da-8766-72044e59ef12"
							}
						],
						"name": "name",
						"status": "pending",
						"updated": "2026-07-24T17:47:13.459860"
					}
				]
			}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *ListVpcsPeeringsResponse) {
				assertEqual(t, 1, got.Meta.Page.Total)
				assertEqual(t, 100, got.Meta.Page.MaxItemsPerPage)
				assertEqual(t, "?_offset=0&_limit=1", got.Meta.Links.Self)
				assertEqual(t, true, got.Meta.Links.Next == nil)
				assertEqual(t, true, got.Meta.Links.Previous == nil)

				assertEqual(t, 1, len(got.Result))
				peering := got.Result[0]
				assertEqual(t, "866b0e4d-59d7-42e6-82d7-aa5fcb72360d", peering.ID)
				assertEqual(t, "name", peering.Name)
				assertEqual(t, "desc", *peering.Description)
				assertEqual(t, VpcsPeeringStatusPending, peering.Status)
				assertEqual(t, true, time.Time(*peering.CreatedAt).Equal(createdAt))
				assertEqual(t, true, time.Time(*peering.Updated).Equal(updated))

				assertEqual(t, 2, len(peering.Members))
				assertEqual(t, "ec32dc9f-679d-4c25-9f3f-b0b37aa101a0", peering.Members[0].ID)
				assertEqual(t, "2ea7d752-c185-422b-9780-2d07109d8afe", peering.Members[0].VpcID)
				assertEqual(t, VpcsPeeringDirectRoleRequester, peering.Members[0].DirectRole)
				assertEqual(t, VpcsPeeringDirectRoleAccepter, peering.Members[1].DirectRole)
			},
		},
		{
			name: "second page of peerings",
			opts: &ListVpcsPeeringsOptions{
				Limit:  helpers.IntPtr(1),
				Offset: helpers.IntPtr(1),
			},
			wantQuery: map[string]string{"_limit": "1", "_offset": "1"},
			response: `{
				"meta": {
					"links": {
						"next": null,
						"previous": "?_offset=0&_limit=1",
						"self": "?_offset=1&_limit=1"
					},
					"page": {"count": 1, "limit": 1, "max_items_per_page": 100, "offset": 1, "total": 2}
				},
				"result": [
					{"id": "peering-2", "name": "second", "status": "created", "members": []}
				]
			}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *ListVpcsPeeringsResponse) {
				assertEqual(t, 1, got.Meta.Page.Offset)
				assertEqual(t, 2, got.Meta.Page.Total)
				assertEqual(t, "?_offset=0&_limit=1", *got.Meta.Links.Previous)
				assertEqual(t, "peering-2", got.Result[0].ID)
			},
		},
		{
			name:      "peerings of a specific vpc",
			opts:      &ListVpcsPeeringsOptions{VpcID: "vpc-requester"},
			wantQuery: map[string]string{"vpc_id": "vpc-requester"},
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
			response:   `{"meta": {"links": {"self": "?"}, "page": {"count": 1, "limit": 10, "max_items_per_page": 100, "offset": 0, "total": 1}}, "result": [{"id": "peering-2", "name": "peering-minimal", "status": "pending", "members": []}]}`,
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

				query := r.URL.Query()
				assertEqual(t, len(tt.wantQuery), len(query))
				for param, want := range tt.wantQuery {
					assertEqual(t, want, query.Get(param))
				}

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
			peeringID: "866b0e4d-59d7-42e6-82d7-aa5fcb72360d",
			response: `{
				"members": [
					{
						"direct_role": "requester",
						"id": "ec32dc9f-679d-4c25-9f3f-b0b37aa101a0",
						"vpc_id": "2ea7d752-c185-422b-9780-2d07109d8afe"
					},
					{
						"direct_role": "accepter",
						"id": "759afe65-f9d5-408e-a4ac-c430dbbd6d35",
						"vpc_id": "732d0b32-b9c0-44da-8766-72044e59ef12"
					}
				],
				"status": "pending_route_table",
				"vpc_peering_id": "866b0e4d-59d7-42e6-82d7-aa5fcb72360d"
			}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *VpcsPeeringMembers) {
				assertEqual(t, "866b0e4d-59d7-42e6-82d7-aa5fcb72360d", got.ID)
				assertEqual(t, VpcsPeeringStatusPendingRouteTable, got.Status)
				assertEqual(t, 2, len(got.Members))
				assertEqual(t, "2ea7d752-c185-422b-9780-2d07109d8afe", got.Members[0].VpcID)
				assertEqual(t, VpcsPeeringDirectRoleRequester, got.Members[0].DirectRole)
				assertEqual(t, "732d0b32-b9c0-44da-8766-72044e59ef12", got.Members[1].VpcID)
				assertEqual(t, VpcsPeeringDirectRoleAccepter, got.Members[1].DirectRole)
			},
		},
		{
			name:       "status the SDK does not know yet",
			peeringID:  "peering-3",
			response:   `{"vpc_peering_id": "peering-3", "status": "some_future_status", "members": []}`,
			statusCode: http.StatusOK,
			check: func(t *testing.T, got *VpcsPeeringMembers) {
				assertEqual(t, VpcsPeeringStatus("some_future_status"), got.Status)
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
			name:       "non-existent peering",
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

// ListAll pages with _offset until Meta.Page.Total is covered; the handler
// answers each offset with a distinct result so a page fetched twice, skipped
// or out of order fails the assertions.
func TestVpcsPeeringsService_ListAll(t *testing.T) {
	t.Parallel()

	peeringJSON := func(id string) string {
		return `{"id": "` + id + `", "name": "peering-` + id + `", "status": "created", "members": []}`
	}
	metaJSON := func(offset, count, total int) string {
		return `{"links": {"self": "?"}, "page": {"count": ` + strconv.Itoa(count) +
			`, "limit": 100, "max_items_per_page": 100, "offset": ` + strconv.Itoa(offset) +
			`, "total": ` + strconv.Itoa(total) + `}}`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/network/v1/vpcs_peerings", r.URL.Path)
		assertEqual(t, "vpc-1", r.URL.Query().Get("vpc_id"))
		assertEqual(t, "100", r.URL.Query().Get("_limit"))

		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Query().Get("_offset") {
		case "0":
			io.WriteString(w, `{"meta": `+metaJSON(0, 2, 102)+`, "result": [`+peeringJSON("peering-1")+`, `+peeringJSON("peering-2")+`]}`)
		case "100":
			io.WriteString(w, `{"meta": `+metaJSON(100, 1, 102)+`, "result": [`+peeringJSON("peering-3")+`]}`)
		default:
			t.Errorf("unexpected offset %q", r.URL.Query().Get("_offset"))
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	client := testPeeringClient(server.URL)
	got, err := client.ListAll(context.Background(), &ListAllVpcsPeeringsOptions{VpcID: "vpc-1"})

	assertNoError(t, err)
	assertEqual(t, 3, len(got))
	assertEqual(t, "peering-1", got[0].ID)
	assertEqual(t, "peering-2", got[1].ID)
	assertEqual(t, "peering-3", got[2].ID)
}

func TestVpcsPeeringsService_ListAllEmpty(t *testing.T) {
	t.Parallel()

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		assertEqual(t, "", r.URL.Query().Get("vpc_id"))

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"meta": {"links": {"self": "?"}, "page": {"count": 0, "limit": 100, "max_items_per_page": 100, "offset": 0, "total": 0}}, "result": []}`)
	}))
	defer server.Close()

	client := testPeeringClient(server.URL)
	got, err := client.ListAll(context.Background(), nil)

	assertNoError(t, err)
	assertEqual(t, 0, len(got))
	assertEqual(t, 1, requests)
}

func TestVpcsPeeringsService_ListAllError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := testPeeringClient(server.URL)
	_, err := client.ListAll(context.Background(), nil)

	assertError(t, err)
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
