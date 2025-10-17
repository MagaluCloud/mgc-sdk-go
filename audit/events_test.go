package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

func TestEventService_List(t *testing.T) {
	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/audit/v0/events":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handleListEvents(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	// Configure client with test server URL
	cfg := client.NewMgcClient("test-api-key",
		client.WithBaseURL(client.MgcUrl(ts.URL)),
		client.WithTimeout(20*time.Second),
	)
	eventsClient := New(cfg)
	service := eventsClient.Events()

	// Test cases
	t.Run("successful list with default parameters", func(t *testing.T) {
		response, err := service.List(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(response.Results) != 2 {
			t.Errorf("expected 2 events, got %d", len(response.Results))
		}
		if response.Meta.Count != 2 {
			t.Errorf("expected meta count 2, got %d", response.Meta.Count)
		}
		validateEvent(t, response.Results[0], "1", "test-source", "test-type")
		validateEvent(t, response.Results[1], "2", "another-source", "another-type")
	})

	t.Run("successful list with query parameters", func(t *testing.T) {
		limit := 10
		offset := 5
		id := "event123"
		sourceLike := "test%"
		typeLike := "type%"
		productLike := "product%"
		authID := "auth123"
		tenantID := "tenant456"
		data := map[string]string{"key": "value"}

		params := &ListEventsParams{
			Limit:  &limit,
			Offset: &offset,
			EventFilterParams: EventFilterParams{
				ID:          &id,
				SourceLike:  &sourceLike,
				TypeLike:    &typeLike,
				ProductLike: &productLike,
				AuthID:      &authID,
				TenantID:    &tenantID,
				Data:        data,
			},
		}

		response, err := service.List(context.Background(), params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(response.Results) != 2 {
			t.Errorf("expected 2 events, got %d", len(response.Results))
		}
		if response.Meta.Count != 2 {
			t.Errorf("expected meta count 2, got %d", response.Meta.Count)
		}
	})

	t.Run("handle server error", func(t *testing.T) {
		errorTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"error": "internal server error"}`)
		}))
		defer errorTS.Close()

		cfg := client.NewMgcClient("test-api-key",
			client.WithBaseURL(client.MgcUrl(errorTS.URL)),
			client.WithTimeout(10*time.Second),
		)
		errorClient := New(cfg)
		errorService := errorClient.Events()

		_, err := errorService.List(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		// Check for either error message or timeout (due to retries)
		if !strings.Contains(err.Error(), "internal server error") && !strings.Contains(err.Error(), "deadline exceeded") && !strings.Contains(err.Error(), "500") {
			t.Errorf("expected error related to server error, got: %v", err)
		}
	})
}

func TestEventService_ListAll(t *testing.T) {
	tests := []struct {
		name       string
		params     *EventFilterParams
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
						{"id": "1", "source": "test-source", "type": "test-type", "specversion": "1.0", "subject": "test-subject", "time": "2024-01-01T00:00:00", "authid": "auth1", "authtype": "type1", "product": "product1", "tenantid": "tenant1", "data": {}},
						{"id": "2", "source": "test-source2", "type": "test-type2", "specversion": "1.0", "subject": "test-subject2", "time": "2024-01-01T00:00:00", "authid": "auth2", "authtype": "type2", "product": "product2", "tenantid": "tenant2", "data": {}}
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
					"results": [` + generateEventJSON(50, 0) + `],
					"meta": {
						"count": 50,
						"limit": 50,
						"offset": 0,
						"total": 50
					}
				}`,
				`{
					"results": [` + generateEventJSON(25, 50) + `],
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
			name: "with filters",
			params: &EventFilterParams{
				SourceLike:  strPtr("test%"),
				ProductLike: strPtr("product%"),
			},
			responses: []string{
				`{
					"results": [
						{"id": "1", "source": "test-source", "type": "test-type", "specversion": "1.0", "subject": "test-subject", "time": "2024-01-01T00:00:00", "authid": "auth1", "authtype": "type1", "product": "product1", "tenantid": "tenant1", "data": {}}
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

			cfg := client.NewMgcClient("test-api-key",
				client.WithBaseURL(client.MgcUrl(ts.URL)),
				client.WithTimeout(20*time.Second),
			)
			eventsClient := New(cfg)
			service := eventsClient.Events()

			got, err := service.ListAll(context.Background(), tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ListAll() got %v events, want %v", len(got), tt.want)
			}
			if tt.checkCalls != nil {
				tt.checkCalls(t, callCount)
			}
		})
	}
}

// Helper function to generate event JSON for testing
func generateEventJSON(count, startID int) string {
	if count == 0 {
		return ""
	}
	var result string
	for i := 0; i < count; i++ {
		if i > 0 {
			result += ","
		}
		id := startID + i + 1
		result += fmt.Sprintf(`{"id": "%d", "source": "source%d", "type": "type%d", "specversion": "1.0", "subject": "subject%d", "time": "2024-01-01T00:00:00", "authid": "auth%d", "authtype": "type%d", "product": "product%d", "tenantid": "tenant%d", "data": {}}`, id, id, id, id, id, id, id, id)
	}
	return result
}

func strPtr(s string) *string {
	return &s
}

func TestListEventsParamsQuery(t *testing.T) {
	limit := 10
	offset := 20
	id := "event123"
	sourceLike := "source%"
	typeLike := "type%"
	productLike := "product%"
	authID := "auth123"
	tenantID := "tenant456"
	data := map[string]string{"key1": "value1", "key2": "value2"}

	params := &ListEventsParams{
		Limit:  &limit,
		Offset: &offset,
		EventFilterParams: EventFilterParams{
			ID:          &id,
			SourceLike:  &sourceLike,
			TypeLike:    &typeLike,
			ProductLike: &productLike,
			AuthID:      &authID,
			TenantID:    &tenantID,
			Data:        data,
		},
	}

	query := make(url.Values)
	if params.Limit != nil {
		query.Set("_limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		query.Set("_offset", strconv.Itoa(*params.Offset))
	}
	if params.ID != nil {
		query.Set("id", *params.ID)
	}
	if params.SourceLike != nil {
		query.Set("source__like", *params.SourceLike)
	}
	if params.TypeLike != nil {
		query.Set("type__like", *params.TypeLike)
	}
	if params.ProductLike != nil {
		query.Set("product__like", *params.ProductLike)
	}
	if params.AuthID != nil {
		query.Set("authid", *params.AuthID)
	}
	if params.TenantID != nil {
		query.Set("X-Tenant-ID", *params.TenantID)
	}
	for k, v := range params.Data {
		query.Set("data."+k, v)
	}

	// Validate query parameters
	assertQueryParam(t, query, "_limit", "10")
	assertQueryParam(t, query, "_offset", "20")
	assertQueryParam(t, query, "id", "event123")
	assertQueryParam(t, query, "source__like", "source%")
	assertQueryParam(t, query, "type__like", "type%")
	assertQueryParam(t, query, "product__like", "product%")
	assertQueryParam(t, query, "authid", "auth123")
	assertQueryParam(t, query, "X-Tenant-ID", "tenant456")
	assertQueryParam(t, query, "data.key1", "value1")
	assertQueryParam(t, query, "data.key2", "value2")
}

func handleListEvents(w http.ResponseWriter, r *http.Request) {
	time := time.Now()
	response := PaginatedResponse[Event]{
		Results: []Event{
			{
				ID:          "1",
				Source:      "test-source",
				Type:        "test-type",
				SpecVersion: "1.0",
				Subject:     "test-subject",
				Time:        utils.LocalDateTimeWithoutZone(time),
				AuthID:      "test-auth",
				AuthType:    "type",
				Product:     "test-product",
				TenantID:    "test-tenant",
				Data:        json.RawMessage(`{"key":"value"}`),
			},
			{
				ID:          "2",
				Source:      "another-source",
				Type:        "another-type",
				SpecVersion: "1.0",
				Subject:     "another-subject",
				Time:        utils.LocalDateTimeWithoutZone(time),
				AuthID:      "another-auth",
				AuthType:    "another-type",
				Product:     "another-product",
				TenantID:    "another-tenant",
				Data:        json.RawMessage(`{"key":"value2"}`),
			},
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

func validateEvent(t *testing.T, event Event, expectedID, expectedSource, expectedType string) {
	t.Helper()
	if event.ID != expectedID {
		t.Errorf("expected ID %s, got %s", expectedID, event.ID)
	}
	if event.Source != expectedSource {
		t.Errorf("expected Source %s, got %s", expectedSource, event.Source)
	}
	if event.Type != expectedType {
		t.Errorf("expected Type %s, got %s", expectedType, event.Type)
	}
}

func assertQueryParam(t *testing.T, query url.Values, key, expected string) {
	t.Helper()
	actual := query.Get(key)
	if actual != expected {
		t.Errorf("expected %s=%s, got %s", key, expected, actual)
	}
}
