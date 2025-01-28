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
		events, err := service.List(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(events) != 2 {
			t.Errorf("expected 2 events, got %d", len(events))
		}
		validateEvent(t, events[0], "1", "test-source", "test-type")
		validateEvent(t, events[1], "2", "another-source", "another-type")
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
			Limit:       &limit,
			Offset:      &offset,
			ID:          &id,
			SourceLike:  &sourceLike,
			TypeLike:    &typeLike,
			ProductLike: &productLike,
			AuthID:      &authID,
			TenantID:    &tenantID,
			Data:        data,
		}

		events, err := service.List(context.Background(), params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(events) != 2 {
			t.Errorf("expected 2 events, got %d", len(events))
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
			client.WithTimeout(1*time.Second),
		)
		errorClient := New(cfg)
		errorService := errorClient.Events()

		_, err := errorService.List(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !contains(err.Error(), "internal server error") {
			t.Errorf("expected 'internal server error' in error, got: %v", err)
		}
	})
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
		Limit:       &limit,
		Offset:      &offset,
		ID:          &id,
		SourceLike:  &sourceLike,
		TypeLike:    &typeLike,
		ProductLike: &productLike,
		AuthID:      &authID,
		TenantID:    &tenantID,
		Data:        data,
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
			Count: 2,
			Total: 2,
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
