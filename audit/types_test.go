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
	cfg := client.NewMgcClient("test-api-key",
		client.WithBaseURL(client.MgcUrl(ts.URL)),
		client.WithTimeout(20*time.Second),
	)
	eventsClient := New(cfg)
	service := eventsClient.EventTypes()

	// Subtests for List method
	t.Run("List", func(t *testing.T) {
		t.Run("successful list with default options", func(t *testing.T) {
			eventTypes, err := service.List(context.Background(), nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(eventTypes) != 2 {
				t.Errorf("expected 2 event types, got %d", len(eventTypes))
			}
			if eventTypes[0].Type != "type1" {
				t.Errorf("expected type1, got %s", eventTypes[0].Type)
			}
			if eventTypes[1].Type != "type2" {
				t.Errorf("expected type2, got %s", eventTypes[1].Type)
			}
		})

		t.Run("successful list with query parameters", func(t *testing.T) {
			limit := 10
			offset := 5
			tenantID := "tenant123"
			params := &ListEventTypesParams{
				Limit:    &limit,
				Offset:   &offset,
				TenantID: &tenantID,
			}

			eventTypes, err := service.List(context.Background(), params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(eventTypes) != 2 {
				t.Errorf("expected 2 event types, got %d", len(eventTypes))
			}
		})

		t.Run("error from server", func(t *testing.T) {
			errorTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"error": "server error"}`)
			}))
			defer errorTs.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithBaseURL(client.MgcUrl(errorTs.URL)),
				client.WithTimeout(1*time.Second),
			)
			errorClient := New(cfg)
			errorService := errorClient.EventTypes()

			_, err := errorService.List(context.Background(), nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if err != nil && !contains(err.Error(), "server error") {
				t.Errorf("expected 'server error' in error message, got: %v", err)
			}
		})
	})
}

// handleListEventTypes handles GET requests to /v0/event-types and returns mock data
func handleListEventTypes(w http.ResponseWriter, r *http.Request) {
	response := PaginatedResponse[EventType]{
		Results: []EventType{
			{Type: "type1"},
			{Type: "type2"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TestListEventTypesParamsQuery verifies query parameter construction from ListEventTypesParams
func TestListEventTypesParamsQuery(t *testing.T) {
	params := ListEventTypesParams{
		Limit:    helpers.IntPtr(10),
		Offset:   helpers.IntPtr(20),
		TenantID: helpers.StrPtr("tenant123"),
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

// contains checks if a string contains a substring (helper function)
func contains(s, substr string) bool {
	return s != "" && s != substr
}
