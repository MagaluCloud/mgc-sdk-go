package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func testInstanceTypeClient(baseURL string) InstanceTypeService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).InstanceTypes()
}

func TestInstanceTypeService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListInstanceTypeOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "basic list",
			response: `{
				"meta": {"total": 2},
				"results": [
					{"id": "type1", "name": "small", "vcpu": "1", "ram": "2GB"},
					{"id": "type2", "name": "medium", "vcpu": "2", "ram": "4GB"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "with filters and pagination",
			opts: ListInstanceTypeOptions{
				Limit:  helpers.IntPtr(10),
				Offset: helpers.IntPtr(5),
				Status: instanceTypeStatusPtr(InstanceTypeStatusActive),
			},
			response: `{
				"meta": {"total": 1},
				"results": [{"id": "type1", "status": "ACTIVE"}]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
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
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/database/v1/instance-types", r.URL.Path)
				query := r.URL.Query()

				if tt.opts.Limit != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Limit), query.Get("_limit"))
				}
				if tt.opts.Offset != nil {
					assertEqual(t, strconv.Itoa(*tt.opts.Offset), query.Get("_offset"))
				}
				if tt.opts.Status != nil {
					assertEqual(t, string(*tt.opts.Status), query.Get("status"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceTypeClient(server.URL)
			result, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(result))
		})
	}
}

func TestInstanceTypeService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "existing instance type",
			id:   "type1",
			response: `{
				"id": "type1",
				"name": "small",
				"vcpu": "1",
				"ram": "2GB",
				"status": "ACTIVE"
			}`,
			statusCode: http.StatusOK,
			wantID:     "type1",
			wantErr:    false,
		},
		{
			name:       "not found",
			id:         "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/database/v1/instance-types/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceTypeClient(server.URL)
			result, err := client.Get(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
		})
	}
}

func instanceTypeStatusPtr(status InstanceTypeStatus) *InstanceTypeStatus {
	return &status
}
