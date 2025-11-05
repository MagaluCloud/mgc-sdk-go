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

func testEngineClient(baseURL string) EngineService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Engines()
}

func TestEngineService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListEngineOptions
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
					{"id": "postgres-16", "name": "PostgreSQL", "version": "16", "status": "PREVIEW"},
					{"id": "mysql-8", "name": "MySQL", "version": "8.0", "status": "ACTIVE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "with pagination and status filter",
			opts: ListEngineOptions{
				Limit:  helpers.IntPtr(10),
				Offset: helpers.IntPtr(5),
				Status: helpers.StrPtr("PREVIEW"),
			},
			response: `{
				"meta": {"total": 1},
				"results": [{"id": "postgres-16", "status": "PREVIEW"}]
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
				assertEqual(t, "/database/v2/engines", r.URL.Path)
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

			client := testEngineClient(server.URL)
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

func TestEngineService_List_PaginationMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/engines", r.URL.Path)
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
				{"id": "engine-1", "name": "PostgreSQL", "version": "16", "status": "ACTIVE"},
				{"id": "engine-2", "name": "MySQL", "version": "8.0", "status": "ACTIVE"}
			]
		}`))
	}))
	defer server.Close()

	client := testEngineClient(server.URL)
	offset := 20
	limit := 10
	status := "ACTIVE"
	result, err := client.List(context.Background(), ListEngineOptions{
		Offset: &offset,
		Limit:  &limit,
		Status: &status,
	})

	assertNoError(t, err)

	// Validate results
	assertEqual(t, 2, len(result.Results))
	assertEqual(t, "engine-1", result.Results[0].ID)
	assertEqual(t, "engine-2", result.Results[1].ID)

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

func TestEngineService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "preview engine",
			id:   "postgres-16",
			response: `{
				"id": "postgres-16",
				"name": "PostgreSQL",
				"version": "16",
				"status": "PREVIEW"
			}`,
			statusCode: http.StatusOK,
			wantID:     "postgres-16",
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
				assertEqual(t, fmt.Sprintf("/database/v2/engines/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testEngineClient(server.URL)
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

func TestEngineService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		filterOpts EngineFilterOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page",
			response: `{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 2, "total": 2, "max_limit": 100}},
				"results": [
					{"id": "postgres-16", "name": "PostgreSQL", "version": "16", "status": "PREVIEW"},
					{"id": "mysql-8", "name": "MySQL", "version": "8.0", "status": "ACTIVE"}
				]
			}`,
			statusCode: http.StatusOK,
			wantCount:  2,
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
			filterOpts: EngineFilterOptions{
				Status: helpers.StrPtr("PREVIEW"),
			},
			response: `{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 1, "total": 1, "max_limit": 100}},
				"results": [
					{"id": "postgres-16", "status": "PREVIEW"}
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
				assertEqual(t, "/database/v2/engines", r.URL.Path)

				query := r.URL.Query()

				// Verify filter parameters
				if tt.filterOpts.Status != nil {
					assertEqual(t, *tt.filterOpts.Status, query.Get("status"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testEngineClient(server.URL)
			engines, err := client.ListAll(context.Background(), tt.filterOpts)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(engines))
		})
	}
}

func TestEngineService_ListAll_MultiplePagesWithPagination(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/engines", r.URL.Path)

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
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "engine-%d", "status": "ACTIVE"}`, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 25, "total": 75, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case "25":
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "engine-%d", "status": "ACTIVE"}`, i+26)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 25, "limit": 25, "count": 25, "total": 75, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case "50":
			results := `[`
			for i := 0; i < 25; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "engine-%d", "status": "ACTIVE"}`, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 25, "count": 25, "total": 75, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case "75":
			response := `{
				"meta": {"page": {"offset": 75, "limit": 25, "count": 0, "total": 75, "max_limit": 100}},
				"results": []
			}`
			w.Write([]byte(response))
		default:
			t.Errorf("unexpected offset: %s", offset)
		}

		requestCount++
	}))
	defer server.Close()

	client := testEngineClient(server.URL)
	engines, err := client.ListAll(context.Background(), EngineFilterOptions{})

	assertNoError(t, err)
	assertEqual(t, 75, len(engines))
	assertEqual(t, 4, requestCount)
}

func TestEngineService_ListAll_WithFilters(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/database/v2/engines", r.URL.Path)

		query := r.URL.Query()

		// Verify filter parameters are present
		if query.Get("status") != "ACTIVE" {
			t.Errorf("expected status=ACTIVE, got %s", query.Get("status"))
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
				results += fmt.Sprintf(`{"id": "engine-%d", "status": "ACTIVE"}`, i+1)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 0, "limit": 25, "count": 25, "total": 60, "max_limit": 100}},
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
				results += fmt.Sprintf(`{"id": "engine-%d", "status": "ACTIVE"}`, i+26)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 25, "limit": 25, "count": 25, "total": 60, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		case 2:
			// Third page with 10 results (< limit triggers stop)
			results := `[`
			for i := 0; i < 10; i++ {
				if i > 0 {
					results += ","
				}
				results += fmt.Sprintf(`{"id": "engine-%d", "status": "ACTIVE"}`, i+51)
			}
			results += `]`
			response := fmt.Sprintf(`{
				"meta": {"page": {"offset": 50, "limit": 25, "count": 10, "total": 60, "max_limit": 100}},
				"results": %s
			}`, results)
			w.Write([]byte(response))
		default:
			t.Errorf("unexpected extra request: %d", requestCount)
		}

		requestCount++
	}))
	defer server.Close()

	client := testEngineClient(server.URL)
	engines, err := client.ListAll(context.Background(), EngineFilterOptions{
		Status: helpers.StrPtr("ACTIVE"),
	})

	assertNoError(t, err)
	assertEqual(t, 60, len(engines))
	assertEqual(t, 3, requestCount)

	// Verify all engines have ACTIVE status
	for _, engine := range engines {
		assertEqual(t, "ACTIVE", engine.Status)
	}
}
