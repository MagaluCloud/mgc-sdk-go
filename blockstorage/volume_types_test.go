package blockstorage

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestVolumeTypeService_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           ListVolumeTypesOptions
		response       string
		statusCode     int
		wantCount      int
		wantErr        bool
		checkQueries   map[string]string
		wantMetaOffset int
		wantMetaLimit  int
		wantMetaTotal  int
	}{
		{
			name:           "basic list",
			opts:           ListVolumeTypesOptions{},
			response:       `{"meta": {"page": {"offset": 0, "limit": 50, "count": 2, "total": 2, "max_limit": 100}}, "types": [{"id": "type1", "name": "SSD"}, {"id": "type2", "name": "HDD"}]}`,
			statusCode:     http.StatusOK,
			wantCount:      2,
			wantMetaOffset: 0,
			wantMetaLimit:  50,
			wantMetaTotal:  2,
			checkQueries: map[string]string{
				"availability-zone": "",
				"name":              "",
				"allows-encryption": "",
			},
		},
		{
			name: "filter by availability zone",
			opts: ListVolumeTypesOptions{
				AvailabilityZone: "zone-a",
			},
			response:       `{"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}}, "types": [{"id": "type1", "name": "SSD"}]}`,
			statusCode:     http.StatusOK,
			wantCount:      1,
			wantMetaOffset: 0,
			wantMetaLimit:  50,
			wantMetaTotal:  1,
			checkQueries: map[string]string{
				"availability-zone": "zone-a",
				"name":              "",
				"allows-encryption": "",
			},
		},
		{
			name: "filter by name",
			opts: ListVolumeTypesOptions{
				Name: "SSD",
			},
			response:       `{"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}}, "types": [{"id": "type1", "name": "SSD"}]}`,
			statusCode:     http.StatusOK,
			wantCount:      1,
			wantMetaOffset: 0,
			wantMetaLimit:  50,
			wantMetaTotal:  1,
			checkQueries: map[string]string{
				"name": "SSD",
			},
		},
		{
			name: "filter by encryption support",
			opts: ListVolumeTypesOptions{
				AllowsEncryption: helpers.BoolPtr(true),
			},
			response:       `{"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}}, "types": [{"id": "type3", "name": "Encrypted"}]}`,
			statusCode:     http.StatusOK,
			wantCount:      1,
			wantMetaOffset: 0,
			wantMetaLimit:  50,
			wantMetaTotal:  1,
			checkQueries: map[string]string{
				"allows-encryption": "true",
			},
		},
		{
			name: "combined filters",
			opts: ListVolumeTypesOptions{
				AvailabilityZone: "zone-b",
				Name:             "NVMe",
				AllowsEncryption: helpers.BoolPtr(false),
			},
			response:       `{"meta": {"page": {"offset": 0, "limit": 50, "count": 1, "total": 1, "max_limit": 100}}, "types": [{"id": "type4", "name": "NVMe"}]}`,
			statusCode:     http.StatusOK,
			wantCount:      1,
			wantMetaOffset: 0,
			wantMetaLimit:  50,
			wantMetaTotal:  1,
			checkQueries: map[string]string{
				"availability-zone": "zone-b",
				"name":              "NVMe",
				"allows-encryption": "false",
			},
		},
		{
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "empty response",
			response:   "",
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "invalid json",
			response:   `{"types": [{"id": "type1", "name": ]}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/volume/v1/volume-types" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				query := r.URL.Query()
				for param, expected := range tt.checkQueries {
					actual := query.Get(param)
					if actual != expected {
						t.Errorf("query param %s: got %s, want %s", param, actual, expected)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientTypes(server.URL)
			resp, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if resp == nil {
				t.Error("expected response, got nil")
				return
			}

			if len(resp.Types) != tt.wantCount {
				t.Errorf("got %d types, want %d", len(resp.Types), tt.wantCount)
			}

			if tt.wantMetaOffset != 0 || tt.wantMetaLimit != 0 || tt.wantMetaTotal != 0 {
				if resp.Meta.Page.Offset != tt.wantMetaOffset {
					t.Errorf("got meta offset %d, want %d", resp.Meta.Page.Offset, tt.wantMetaOffset)
				}
				if resp.Meta.Page.Limit != tt.wantMetaLimit {
					t.Errorf("got meta limit %d, want %d", resp.Meta.Page.Limit, tt.wantMetaLimit)
				}
				if resp.Meta.Page.Total != tt.wantMetaTotal {
					t.Errorf("got meta total %d, want %d", resp.Meta.Page.Total, tt.wantMetaTotal)
				}
			}
		})
	}
}

func TestVolumeTypeService_List_QueryParams(t *testing.T) {
	tests := []struct {
		name         string
		opts         ListVolumeTypesOptions
		expectParams map[string]string
	}{
		{
			name: "allows encryption true",
			opts: ListVolumeTypesOptions{
				AllowsEncryption: helpers.BoolPtr(true),
			},
			expectParams: map[string]string{
				"allows-encryption": "true",
			},
		},
		{
			name: "allows encryption false",
			opts: ListVolumeTypesOptions{
				AllowsEncryption: helpers.BoolPtr(false),
			},
			expectParams: map[string]string{
				"allows-encryption": "false",
			},
		},
		{
			name: "all filters combined",
			opts: ListVolumeTypesOptions{
				AvailabilityZone: "zone-c",
				Name:             "Fast",
				AllowsEncryption: helpers.BoolPtr(true),
			},
			expectParams: map[string]string{
				"availability-zone": "zone-c",
				"name":              "Fast",
				"allows-encryption": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				for param, expected := range tt.expectParams {
					actual := query.Get(param)
					if actual != expected {
						t.Errorf("query param %s: got %s, want %s", param, actual, expected)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"meta": {"page": {"offset": 0, "limit": 50, "count": 0, "total": 0, "max_limit": 100}}, "types": []}`))
			}))
			defer server.Close()

			client := testClientTypes(server.URL)
			_, err := client.List(context.Background(), tt.opts)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestVolumeTypeService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "single page",
			response:   `{"meta": {"page": {"offset": 0, "limit": 50, "count": 2, "total": 2, "max_limit": 100}}, "types": [{"id": "type1", "name": "SSD"}, {"id": "type2", "name": "HDD"}]}`,
			statusCode: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "empty result",
			response:   `{"meta": {"page": {"offset": 0, "limit": 50, "count": 0, "total": 0, "max_limit": 100}}, "types": []}`,
			statusCode: http.StatusOK,
			wantCount:  0,
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
				if r.URL.Path != "/volume/v1/volume-types" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClientTypes(server.URL)
			types, err := client.ListAll(context.Background(), VolumeTypeFilterOptions{})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(types) != tt.wantCount {
				t.Errorf("got %d types, want %d", len(types), tt.wantCount)
			}
		})
	}
}

func TestVolumeTypeService_ListAll_MultiplePagesWithPagination(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/volume/v1/volume-types" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		query := r.URL.Query()
		offset := query.Get("_offset")
		limit := query.Get("_limit")

		if limit != "50" {
			t.Errorf("expected limit 50, got %s", limit)
		}

		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Simulate pagination: first page has 50 items, second page has 25
		switch offset {
		case "0":
			// First page: 50 items
			types := make([]string, 50)
			for i := 0; i < 50; i++ {
				types[i] = fmt.Sprintf(`{"id": "type%d", "name": "Type%d"}`, i+1, i+1)
			}
			response := fmt.Sprintf(`{"meta": {"page": {"offset": 0, "limit": 50, "count": 50, "total": 75, "max_limit": 100}}, "types": [%s]}`,
				strings.Join(types, ","))
			w.Write([]byte(response))
		case "50":
			// Second page: 25 items
			types := make([]string, 25)
			for i := 0; i < 25; i++ {
				types[i] = fmt.Sprintf(`{"id": "type%d", "name": "Type%d"}`, i+51, i+51)
			}
			response := fmt.Sprintf(`{"meta": {"page": {"offset": 50, "limit": 50, "count": 25, "total": 75, "max_limit": 100}}, "types": [%s]}`,
				strings.Join(types, ","))
			w.Write([]byte(response))
		default:
			t.Errorf("unexpected offset: %s", offset)
		}
	}))
	defer server.Close()

	client := testClientTypes(server.URL)
	types, err := client.ListAll(context.Background(), VolumeTypeFilterOptions{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Should have fetched all 75 types across 2 pages
	if len(types) != 75 {
		t.Errorf("got %d types, want 75", len(types))
	}

	// Should have made exactly 2 requests
	if requestCount != 2 {
		t.Errorf("made %d requests, want 2", requestCount)
	}

	// Verify first and last items
	if types[0].ID != "type1" {
		t.Errorf("first type ID: got %s, want type1", types[0].ID)
	}
	if types[74].ID != "type75" {
		t.Errorf("last type ID: got %s, want type75", types[74].ID)
	}
}

func TestVolumeTypeService_ListAll_WithFilters(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/volume/v1/volume-types" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		query := r.URL.Query()

		// Verify filter parameters are present
		if query.Get("availability-zone") != "zone-a" {
			t.Errorf("expected availability-zone=zone-a, got %s", query.Get("availability-zone"))
		}
		if query.Get("allows-encryption") != "true" {
			t.Errorf("expected allows-encryption=true, got %s", query.Get("allows-encryption"))
		}
		if query.Get("_sort") != "name:asc" {
			t.Errorf("expected _sort=name:asc, got %s", query.Get("_sort"))
		}

		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Return 50 items on first page, 25 on second
		offset := query.Get("_offset")
		switch offset {
		case "0":
			types := make([]string, 50)
			for i := 0; i < 50; i++ {
				types[i] = fmt.Sprintf(`{"id": "type%d", "name": "Type%d", "availability_zones": ["zone-a"], "allows_encryption": true}`, i+1, i+1)
			}
			response := fmt.Sprintf(`{
				"meta": {
					"filters": [
						{"field": "availability-zone", "value": "zone-a"},
						{"field": "allows-encryption", "value": "true"}
					],
					"page": {"offset": 0, "limit": 50, "count": 50, "total": 75, "max_limit": 100}
				},
				"types": [%s]
			}`, strings.Join(types, ","))
			w.Write([]byte(response))
		case "50":
			types := make([]string, 25)
			for i := 0; i < 25; i++ {
				types[i] = fmt.Sprintf(`{"id": "type%d", "name": "Type%d", "availability_zones": ["zone-a"], "allows_encryption": true}`, i+51, i+51)
			}
			response := fmt.Sprintf(`{
				"meta": {
					"filters": [
						{"field": "availability-zone", "value": "zone-a"},
						{"field": "allows-encryption", "value": "true"}
					],
					"page": {"offset": 50, "limit": 50, "count": 25, "total": 75, "max_limit": 100}
				},
				"types": [%s]
			}`, strings.Join(types, ","))
			w.Write([]byte(response))
		default:
			t.Errorf("unexpected offset: %s", offset)
		}
	}))
	defer server.Close()

	client := testClientTypes(server.URL)
	types, err := client.ListAll(context.Background(), VolumeTypeFilterOptions{
		AvailabilityZone: "zone-a",
		AllowsEncryption: helpers.BoolPtr(true),
		Sort:             helpers.StrPtr("name:asc"),
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Should have fetched all 75 types
	if len(types) != 75 {
		t.Errorf("expected 75 types, got %d", len(types))
	}

	// Should have made exactly 2 requests
	if requestCount != 2 {
		t.Errorf("made %d requests, want 2", requestCount)
	}

	// Verify all types have the filtered properties
	for _, vt := range types {
		if vt.AllowsEncryption != true {
			t.Errorf("expected allows_encryption true, got %v", vt.AllowsEncryption)
		}
	}
}

func testClientTypes(baseURL string) VolumeTypeService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).VolumeTypes()
}
