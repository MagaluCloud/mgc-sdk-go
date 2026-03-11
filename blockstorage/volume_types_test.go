package blockstorage

import (
	"context"
	"net/http"
	"net/http/httptest"
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
				AvailabilityZone: helpers.StrPtr("zone-a"),
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
				Name: helpers.StrPtr("SSD"),
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
				AvailabilityZone: helpers.StrPtr("zone-b"),
				Name:             helpers.StrPtr("NVMe"),
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
				AvailabilityZone: helpers.StrPtr("zone-c"),
				Name:             helpers.StrPtr("Fast"),
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

func testClientTypes(baseURL string) VolumeTypeService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).VolumeTypes()
}
