package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

func TestSubnetPoolService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "successful list with pagination",
			opts: ListOptions{
				Limit:  helpers.IntPtr(10),
				Offset: helpers.IntPtr(20),
				Sort:   helpers.StrPtr("name"),
			},
			response: `{
				"meta": {
					"page": {"limit": 10, "offset": 20, "count": 5, "total": 100},
					"links": {"self": "/network/v0/subnetpools"}
				},
				"results": [
					{"id": "pool1", "name": "test-pool1"},
					{"id": "pool2", "name": "test-pool2"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			response:   `{"meta": {}, "results": []}`,
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "invalid parameters",
			opts:       ListOptions{Limit: helpers.IntPtr(1000)},
			response:   `{"error": "invalid limit"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v0/subnetpools", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				query := r.URL.Query()
				if tt.opts.Limit != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Limit), query.Get("_limit"))
				}
				if tt.opts.Offset != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Offset), query.Get("_offset"))
				}
				if tt.opts.Sort != nil {
					assertEqual(t, *tt.opts.Sort, query.Get("_sort"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSubnetPoolClient(server.URL)
			pools, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(pools))
		})
	}
}

func TestSubnetPoolService_Get(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")

	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		want       *SubnetPoolDetailsResponse
		wantErr    bool
	}{
		{
			name: "existing pool",
			id:   "pool1",
			response: `{
				"id": "pool1",
				"name": "test-pool",
				"cidr": "10.0.0.0/16",
				"ip_version": 4,
				"created_at": "2024-01-01T00:00:00"
			}`,
			statusCode: http.StatusOK,
			want: &SubnetPoolDetailsResponse{
				ID:        "pool1",
				Name:      "test-pool",
				CIDR:      helpers.StrPtr("10.0.0.0/16"),
				IPVersion: 4,
				CreatedAt: utils.LocalDateTimeWithoutZone(createdAt),
			},
			wantErr: false,
		},
		{
			name:       "non-existent pool",
			id:         "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/subnetpools/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSubnetPoolClient(server.URL)
			pool, err := client.Get(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want.ID, pool.ID)
			assertEqual(t, tt.want.Name, pool.Name)
			assertEqual(t, *tt.want.CIDR, *pool.CIDR)
		})
	}
}

func TestSubnetPoolService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateSubnetPoolRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful create",
			request: CreateSubnetPoolRequest{
				Name:        "test-pool",
				Description: "test description",
				CIDR:        helpers.StrPtr("10.0.0.0/16"),
			},
			response:   `{"id": "pool-new"}`,
			statusCode: http.StatusOK,
			wantID:     "pool-new",
			wantErr:    false,
		},
		{
			name: "missing required field",
			request: CreateSubnetPoolRequest{
				Description: "invalid",
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
				assertEqual(t, "/network/v0/subnetpools", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				if !tt.wantErr {
					var req CreateSubnetPoolRequest
					err := json.NewDecoder(r.Body).Decode(&req)
					assertNoError(t, err)
					assertEqual(t, tt.request.Name, req.Name)
					assertEqual(t, tt.request.Description, req.Description)
					assertEqual(t, *tt.request.CIDR, *req.CIDR)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSubnetPoolClient(server.URL)
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

func TestSubnetPoolService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "pool1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent pool",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/subnetpools/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testSubnetPoolClient(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestSubnetPoolService_BookCIDR(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    BookCIDRRequest
		response   string
		statusCode int
		wantCIDR   string
		wantErr    bool
	}{
		{
			name: "book by cidr",
			id:   "pool1",
			request: BookCIDRRequest{
				CIDR: helpers.StrPtr("10.0.1.0/24"),
			},
			response:   `{"cidr": "10.0.1.0/24"}`,
			statusCode: http.StatusOK,
			wantCIDR:   "10.0.1.0/24",
			wantErr:    false,
		},
		{
			name:       "invalid request",
			id:         "pool1",
			request:    BookCIDRRequest{},
			response:   `{"error": "cidr or mask required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/subnetpools/%s/book_cidr", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req BookCIDRRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				if tt.request.CIDR != nil {
					assertEqual(t, *tt.request.CIDR, *req.CIDR)
				} else {
					assertEqual(t, tt.request.CIDR, req.CIDR)
				}
				assertEqual(t, tt.request.Mask, req.Mask)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSubnetPoolClient(server.URL)
			resp, err := client.BookCIDR(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCIDR, resp.CIDR)
		})
	}
}

func TestSubnetPoolService_UnbookCIDR(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    UnbookCIDRRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful unbook",
			id:         "pool1",
			request:    UnbookCIDRRequest{CIDR: "10.0.1.0/24"},
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent cidr",
			id:         "pool1",
			request:    UnbookCIDRRequest{CIDR: "10.0.9.0/24"},
			statusCode: http.StatusNotFound,
			response:   `{"error": "cidr not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/subnetpools/%s/unbook_cidr", tt.id), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req UnbookCIDRRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.request.CIDR, req.CIDR)

				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testSubnetPoolClient(server.URL)
			err := client.UnbookCIDR(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func testSubnetPoolClient(baseURL string) SubnetPoolService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).SubnetPools()
}
