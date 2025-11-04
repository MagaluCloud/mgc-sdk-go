package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestEventTypeService(t *testing.T) {
	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/audit/v0/event-types":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handleListEventTypes(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	// Configure client with test server URL
	cfg := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(ts.URL)),
		client.WithTimeout(20*time.Second),
	)
	eventsClient := New(cfg)
	service := eventsClient.EventTypes()

	// Subtests for List method
	t.Run("List", func(t *testing.T) {
		t.Run("successful list with default options", func(t *testing.T) {
			response, err := service.List(context.Background(), nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(response.Results) != 2 {
				t.Errorf("expected 2 event types, got %d", len(response.Results))
			}
			if response.Results[0].Type != "type1" {
				t.Errorf("expected type1, got %s", response.Results[0].Type)
			}
			if response.Results[1].Type != "type2" {
				t.Errorf("expected type2, got %s", response.Results[1].Type)
			}
		})

		t.Run("successful list with query parameters", func(t *testing.T) {
			limit := 10
			offset := 5
			tenantID := "tenant123"
			params := &ListEventTypesParams{
				Limit:  &limit,
				Offset: &offset,
				EventTypeFilterParams: EventTypeFilterParams{
					TenantID: &tenantID,
				},
			}

			response, err := service.List(context.Background(), params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(response.Results) != 2 {
				t.Errorf("expected 2 event types, got %d", len(response.Results))
			}
			if response.Meta.Count != 2 {
				t.Errorf("expected meta count 2, got %d", response.Meta.Count)
			}
		})

		t.Run("error from server", func(t *testing.T) {
			errorTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"error": "server error"}`)
			}))
			defer errorTs.Close()

			cfg := client.NewMgcClient(client.WithAPIKey("test-api-key"),
				client.WithBaseURL(client.MgcUrl(errorTs.URL)),
				client.WithTimeout(10*time.Second),
			)
			errorClient := New(cfg)
			errorService := errorClient.EventTypes()

			_, err := errorService.List(context.Background(), nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})
}

func TestEventTypeService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		params     *EventTypeFilterParams
		responses  []string
		want       int
		wantErr    bool
		checkCalls func(*testing.T, int)
	}{
		{
			name:   "single page",
			params: nil,
			responses: []string{
				`{
					"results": [
						{"type": "type1"},
						{"type": "type2"}
					],
					"meta": {
						"count": 2,
						"limit": 50,
						"offset": 0,
						"total": 2
					}
				}`,
			},
			want:    2,
			wantErr: false,
		},
		{
			name:   "multiple pages",
			params: nil,
			responses: []string{
				`{
					"results": [` + generateEventTypeJSON(50, 0) + `],
					"meta": {
						"count": 50,
						"limit": 50,
						"offset": 0,
						"total": 50
					}
				}`,
				`{
					"results": [` + generateEventTypeJSON(25, 50) + `],
					"meta": {
						"count": 75,
						"limit": 50,
						"offset": 50,
						"total": 25
					}
				}`,
			},
			want:    75,
			wantErr: false,
			checkCalls: func(t *testing.T, calls int) {
				if calls != 2 {
					t.Errorf("expected 2 API calls, got %d", calls)
				}
			},
		},
		{
			name: "with tenant filter",
			params: &EventTypeFilterParams{
				TenantID: strPtr("tenant123"),
			},
			responses: []string{
				`{
					"results": [
						{"type": "type1"}
					],
					"meta": {
						"count": 1,
						"limit": 50,
						"offset": 0,
						"total": 1
					}
				}`,
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "empty results",
			params: nil,
			responses: []string{
				`{
					"results": [],
					"meta": {
						"count": 0,
						"limit": 50,
						"offset": 0,
						"total": 0
					}
				}`,
			},
			want:    0,
			wantErr: false,
		},
		{
			name:   "server error",
			params: nil,
			responses: []string{
				`{"error": "internal server error"}`,
				`{"error": "internal server error"}`,
				`{"error": "internal server error"}`,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if callCount >= len(tt.responses) {
					t.Errorf("unexpected API call #%d", callCount+1)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				w.Write([]byte(tt.responses[callCount]))
				callCount++
			}))
			defer ts.Close()

			cfg := client.NewMgcClient(client.WithAPIKey("test-api-key"),
				client.WithBaseURL(client.MgcUrl(ts.URL)),
				client.WithTimeout(20*time.Second),
			)
			eventsClient := New(cfg)
			service := eventsClient.EventTypes()

			got, err := service.ListAll(context.Background(), tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ListAll() got %v event types, want %v", len(got), tt.want)
			}
			if tt.checkCalls != nil {
				tt.checkCalls(t, callCount)
			}
		})
	}
}

// Helper function to generate event type JSON for testing
func generateEventTypeJSON(count, startID int) string {
	if count == 0 {
		return ""
	}
	var result string
	for i := 0; i < count; i++ {
		if i > 0 {
			result += ","
		}
		id := startID + i + 1
		result += fmt.Sprintf(`{"type": "type%d"}`, id)
	}
	return result
}

// handleListEventTypes handles GET requests to /v0/event-types and returns mock data
func handleListEventTypes(w http.ResponseWriter, r *http.Request) {
	response := PaginatedResponse[EventType]{
		Results: []EventType{
			{Type: "type1"},
			{Type: "type2"},
		},
		Meta: PaginatedMeta{
			Count:  2,
			Limit:  50,
			Offset: 0,
			Total:  2,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TestListEventTypesParamsQuery verifies query parameter construction from ListEventTypesParams
func TestListEventTypesParamsQuery(t *testing.T) {
	params := ListEventTypesParams{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(20),
		EventTypeFilterParams: EventTypeFilterParams{
			TenantID: helpers.StrPtr("tenant123"),
		},
	}

	query := make(url.Values)
	if params.Limit != nil {
		query.Set("_limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		query.Set("_offset", strconv.Itoa(*params.Offset))
	}
	if params.TenantID != nil {
		query.Set("X-Tenant-ID", *params.TenantID)
	}

	if query.Get("_limit") != "10" {
		t.Errorf("expected _limit=10, got %s", query.Get("_limit"))
	}
	if query.Get("_offset") != "20" {
		t.Errorf("expected _offset=20, got %s", query.Get("_offset"))
	}
	if query.Get("X-Tenant-ID") != "tenant123" {
		t.Errorf("expected X-Tenant-ID=tenant123, got %s", query.Get("X-Tenant-ID"))
	}
}
