package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func testInstanceTypeClient(baseURL string) InstanceTypeService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
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
				"meta": {"total": 6},
				"results": [
					{"id": "type1_mysql8", "name": "small", "vcpu": "1", "ram": "2GB", "status": "ACTIVE", "engine_id": "mysql8_id", "compatible_product": "SINGLE_INSTANCE"},
					{"id": "type2_mysql8", "name": "medium", "vcpu": "2", "ram": "4GB", "status": "ACTIVE", "engine_id": "mysql8_id", "compatible_product": "SINGLE_INSTANCE"},
					{"id": "type3_mysql8", "name": "medium", "vcpu": "2", "ram": "4GB", "status": "DEPRECATED", "engine_id": "mysql8_id", "compatible_product": "SINGLE_INSTANCE"},
					{"id": "type1_mysql84", "name": "small", "vcpu": "1", "ram": "2GB", "status": "ACTIVE", "engine_id": "mysql84_id", "compatible_product": "SINGLE_INSTANCE_REPLICA"},
					{"id": "type2_mysql84", "name": "medium", "vcpu": "2", "ram": "4GB", "status": "ACTIVE", "engine_id": "mysql84_id", "compatible_product": "CLUSTER"},
					{"id": "type1_postgres16", "name": "small", "vcpu": "1", "ram": "2GB", "status": "ACTIVE", "engine_id": "postgresql16_id", "compatible_product": "SINGLE_INSTANCE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  6,
			wantErr:    false,
		},
		{
			name: "with filters and pagination",
			opts: ListInstanceTypeOptions{
				Limit:             helpers.IntPtr(10),
				Offset:            helpers.IntPtr(5),
				Status:            helpers.StrPtr("ACTIVE"),
				EngineID:          helpers.StrPtr("mysql8_id"),
				CompatibleProduct: helpers.StrPtr("SINGLE_INSTANCE"),
			},
			response: `{
				"meta": {"total": 1},
				"results": [{"id": "type1_mysql8", "status": "ACTIVE", "engine_id": "mysql8_id", "compatible_product": "SINGLE_INSTANCE"}]
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
				assertEqual(t, "/database/v2/instance-types", r.URL.Path)
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
				if tt.opts.EngineID != nil {
					assertEqual(t, string(*tt.opts.EngineID), query.Get("engine_id"))
				}
				if tt.opts.EngineID != nil {
					assertEqual(t, string(*tt.opts.CompatibleProduct), query.Get("compatible_product"))
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
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(result.Results))
		})
	}
}

func TestInstanceTypeService_List_PaginationMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/instance-types", r.URL.Path)
		query := r.URL.Query()

		assertEqual(t, "10", query.Get("_limit"))
		assertEqual(t, "20", query.Get("_offset"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"meta": {
				"page": {
					"offset": 20,
					"limit": 10,
					"count": 10,
					"total": 50,
					"max_limit": 100
				},
				"filters": [
					{"field": "status", "value": "ACTIVE"}
				]
			},
			"results": [
				{"id": "type1", "name": "small", "vcpu": "1", "ram": "2GB", "status": "ACTIVE"},
				{"id": "type2", "name": "medium", "vcpu": "2", "ram": "4GB", "status": "ACTIVE"}
			]
		}`))
	}))
	defer server.Close()

	client := testInstanceTypeClient(server.URL)
	offset := 20
	limit := 10
	status := "ACTIVE"
	result, err := client.List(context.Background(), ListInstanceTypeOptions{
		Offset: &offset,
		Limit:  &limit,
		Status: &status,
	})

	assertNoError(t, err)

	// Validate results
	assertEqual(t, 2, len(result.Results))
	assertEqual(t, "type1", result.Results[0].ID)
	assertEqual(t, "type2", result.Results[1].ID)

	// Validate pagination metadata
	assertEqual(t, 20, result.Meta.Page.Offset)
	assertEqual(t, 10, result.Meta.Page.Limit)
	assertEqual(t, 10, result.Meta.Page.Count)
	assertEqual(t, 50, result.Meta.Page.Total)
	assertEqual(t, 100, result.Meta.Page.MaxLimit)

	// Validate filters metadata
	assertEqual(t, 1, len(result.Meta.Filters))
	assertEqual(t, "status", result.Meta.Filters[0].Field)
	assertEqual(t, "ACTIVE", result.Meta.Filters[0].Value)
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
				"status": "ACTIVE",
				"engine_id": "mysql-8_id"
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
				assertEqual(t, fmt.Sprintf("/database/v2/instance-types/%s", tt.id), r.URL.Path)
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
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, result.ID)
		})
	}
}

func TestInstanceTypeService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		filterOpts InstanceTypeFilterOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page",
			response: `{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 3, "total": 3, "max_limit": 100}},
				"results": [
					{"id": "type1", "name": "small", "vcpu": "1", "ram": "2GB", "status": "ACTIVE"},
					{"id": "type2", "name": "medium", "vcpu": "2", "ram": "4GB", "status": "ACTIVE"},
					{"id": "type3", "name": "large", "vcpu": "4", "ram": "8GB", "status": "ACTIVE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  3,
		},
		{
			name: "empty result",
			response: `{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 0, "total": 0, "max_limit": 100}},
				"results": []
			}`,
			statusCode: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "with status filter",
			filterOpts: InstanceTypeFilterOptions{
				Status: helpers.StrPtr("ACTIVE"),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 2, "total": 2, "max_limit": 100}},
				"results": [
					{"id": "type1", "status": "ACTIVE"},
					{"id": "type2", "status": "ACTIVE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "with engine_id filter",
			filterOpts: InstanceTypeFilterOptions{
				EngineID: helpers.StrPtr("mysql8_id"),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 1, "total": 1, "max_limit": 100}},
				"results": [
					{"id": "type1_mysql8", "engine_id": "mysql8_id"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "with compatible_product filter",
			filterOpts: InstanceTypeFilterOptions{
				CompatibleProduct: helpers.StrPtr("SINGLE_INSTANCE"),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 1, "total": 1, "max_limit": 100}},
				"results": [
					{"id": "type1", "compatible_product": "SINGLE_INSTANCE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  1,
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
				assertEqual(t, "/database/v2/instance-types", r.URL.Path)

				query := r.URL.Query()

				// Verify filter parameters
				if tt.filterOpts.Status != nil {
					assertEqual(t, *tt.filterOpts.Status, query.Get("status"))
				}
				if tt.filterOpts.EngineID != nil {
					assertEqual(t, *tt.filterOpts.EngineID, query.Get("engine_id"))
				}
				if tt.filterOpts.CompatibleProduct != nil {
					assertEqual(t, *tt.filterOpts.CompatibleProduct, query.Get("compatible_product"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testInstanceTypeClient(server.URL)
			instanceTypes, err := client.ListAll(context.Background(), tt.filterOpts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(instanceTypes))
		})
	}
}

func TestInstanceTypeService_ListAll_MultiplePagesWithPagination(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/instance-types", r.URL.Path)

		query := r.URL.Query()
		offset := query.Get("_offset")
		limit := query.Get("_limit")

		if limit != "25" {
			t.Errorf("expected limit 25, got %s", limit)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch offset {
		case "0":
			// First page: 25 items
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "type-%d", "status": "ACTIVE"}`, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 25, "total": 80, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case "25":
			// Second page: 25 items
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "type-%d", "status": "ACTIVE"}`, i+26)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 25, "limit": 25, "count": 25, "total": 80, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case "50":
			// Third page: 25 items
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "type-%d", "status": "ACTIVE"}`, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 25, "count": 25, "total": 80, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case "75":
			// Fourth page: remaining 5 items (break condition)
			results := `[`
			for i := 0; i < 5; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "type-%d", "status": "ACTIVE"}`, i+76)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 75, "limit": 25, "count": 5, "total": 80, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		default:
			t.Errorf("unexpected offset: %s", offset)
		}

		requestCount++
	}))
	defer server.Close()

	client := testInstanceTypeClient(server.URL)
	instanceTypes, err := client.ListAll(context.Background(), InstanceTypeFilterOptions{})

	assertNoError(t, err)
	assertEqual(t, 80, len(instanceTypes))
	assertEqual(t, 4, requestCount)
}

func TestInstanceTypeService_ListAll_WithFilters(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/instance-types", r.URL.Path)

		query := r.URL.Query()

		// Verify filter parameters are present
		if query.Get("status") != "ACTIVE" {
			t.Errorf("expected status=ACTIVE, got %s", query.Get("status"))
		}
		if query.Get("engine_id") != "mysql8_id" {
			t.Errorf("expected engine_id=mysql8_id, got %s", query.Get("engine_id"))
		}
		if query.Get("compatible_product") != "SINGLE_INSTANCE" {
			t.Errorf("expected compatible_product=SINGLE_INSTANCE, got %s", query.Get("compatible_product"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch requestCount {
		case 0:
			// First page with 25 results
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "type-%d", "status": "ACTIVE", "engine_id": "mysql8_id", "compatible_product": "SINGLE_INSTANCE"}`, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 25, "total": 65, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case 1:
			// Second page with 25 results
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "type-%d", "status": "ACTIVE", "engine_id": "mysql8_id", "compatible_product": "SINGLE_INSTANCE"}`, i+26)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 25, "limit": 25, "count": 25, "total": 65, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case 2:
			// Third page with 15 results
			results := `[`
			for i := 0; i < 15; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "type-%d", "status": "ACTIVE", "engine_id": "mysql8_id", "compatible_product": "SINGLE_INSTANCE"}`, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 25, "count": 15, "total": 65, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		default:
			t.Errorf("unexpected extra request: %d", requestCount)
		}

		requestCount++
	}))
	defer server.Close()

	client := testInstanceTypeClient(server.URL)
	instanceTypes, err := client.ListAll(context.Background(), InstanceTypeFilterOptions{
		Status:            helpers.StrPtr("ACTIVE"),
		EngineID:          helpers.StrPtr("mysql8_id"),
		CompatibleProduct: helpers.StrPtr("SINGLE_INSTANCE"),
	})

	assertNoError(t, err)
	assertEqual(t, 65, len(instanceTypes))
	assertEqual(t, 3, requestCount)

	// Verify all instance types have the correct compatible_product filter applied
	for _, it := range instanceTypes {
		assertEqual(t, "SINGLE_INSTANCE", it.CompatibleProduct)
	}
}
