package tag

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func newValuesTestServer() (*httptest.Server, TagValueService) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// /v0/tags/{tag_name}/values
		// /v0/tags/{tag_name}/values/{value_name}
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v0/tags/"), "/")
		// parts[0] = tag_name, parts[1] = "values", parts[2] = value_name (optional)
		switch {
		case len(parts) == 2 && parts[1] == "values" && r.Method == http.MethodGet:
			handleListTagValues(w, r)
		case len(parts) == 2 && parts[1] == "values" && r.Method == http.MethodPost:
			handleCreateTagValue(w, r)
		case len(parts) == 3 && parts[1] == "values" && r.Method == http.MethodGet:
			handleGetTagValue(w, r, parts[2])
		case len(parts) == 3 && parts[1] == "values" && r.Method == http.MethodPatch:
			handleUpdateTagValue(w, r, parts[2])
		case len(parts) == 3 && parts[1] == "values" && r.Method == http.MethodDelete:
			handleDeleteTagValue(w, r, parts[2])
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))

	cfg := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithTimeout(20*time.Second),
	)
	c := New(cfg, WithBasePath(client.MgcUrl(ts.URL)))
	return ts, c.Values()
}

func newValuesErrorService() (TagValueService, func()) {
	ts := errorServer()
	cfg := client.NewMgcClient(client.WithAPIKey("test-api-key"))
	svc := New(cfg, WithBasePath(client.MgcUrl(ts.URL))).Values()
	return svc, ts.Close
}

// --- handlers ---

func handleListTagValues(w http.ResponseWriter, r *http.Request) {
	response := ListTagValuesResponse{
		Results: []TagValue{
			{Name: "production", Description: "Prod", CreatedAt: "2024-01-01T00:00:00Z"},
			{Name: "staging", Description: "Stage", CreatedAt: "2024-01-02T00:00:00Z"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleCreateTagValue(w http.ResponseWriter, r *http.Request) {
	var req CreateTagValueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	response := TagValue{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetTagValue(w http.ResponseWriter, r *http.Request, valueName string) {
	response := TagValue{
		Name:        valueName,
		Description: "A tag value",
		CreatedAt:   "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleUpdateTagValue(w http.ResponseWriter, r *http.Request, valueName string) {
	var req UpdateTagValueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	desc := "updated"
	if req.Description != nil {
		desc = *req.Description
	}
	response := TagValue{
		Name:        valueName,
		Description: desc,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteTagValue(w http.ResponseWriter, r *http.Request, valueName string) {
	response := TagValue{
		Name:      valueName,
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// --- tests ---

func TestTagValueService_List(t *testing.T) {
	ts, svc := newValuesTestServer()
	defer ts.Close()

	t.Run("lists all values for a tag", func(t *testing.T) {
		values, err := svc.List(context.Background(), "environment", ListTagValuesOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(values) != 2 {
			t.Fatalf("expected 2 values, got %d", len(values))
		}
		if values[0].Name != "production" {
			t.Errorf("expected 'production', got %s", values[0].Name)
		}
		if values[1].Name != "staging" {
			t.Errorf("expected 'staging', got %s", values[1].Name)
		}
	})

	t.Run("lists values with filters", func(t *testing.T) {
		values, err := svc.List(context.Background(), "environment", ListTagValuesOptions{
			Name:   helpers.StrPtr("prod"),
			Limit:  helpers.IntPtr(10),
			Offset: helpers.IntPtr(0),
			Sort:   helpers.StrPtr("name"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(values) == 0 {
			t.Error("expected at least one value")
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.List(context.Background(), "", ListTagValuesOptions{})
		if err == nil {
			t.Fatal("expected error for empty tag_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "tag_name" {
			t.Errorf("expected field 'tag_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newValuesErrorService()
		defer closeServer()
		_, err := svc.List(context.Background(), "my-tag", ListTagValuesOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagValueService_Get(t *testing.T) {
	ts, svc := newValuesTestServer()
	defer ts.Close()

	t.Run("gets a tag value successfully", func(t *testing.T) {
		v, err := svc.Get(context.Background(), "environment", "production")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Name != "production" {
			t.Errorf("expected name 'production', got %s", v.Name)
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.Get(context.Background(), "", "production")
		if err == nil {
			t.Fatal("expected error for empty tag_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "tag_name" {
			t.Errorf("expected field 'tag_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error when value_name is empty", func(t *testing.T) {
		_, err := svc.Get(context.Background(), "environment", "")
		if err == nil {
			t.Fatal("expected error for empty value_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "value_name" {
			t.Errorf("expected field 'value_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newValuesErrorService()
		defer closeServer()
		_, err := svc.Get(context.Background(), "my-tag", "my-value")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagValueService_Create(t *testing.T) {
	ts, svc := newValuesTestServer()
	defer ts.Close()

	t.Run("creates a tag value successfully", func(t *testing.T) {
		v, err := svc.Create(context.Background(), "environment", CreateTagValueRequest{
			Name:        "production",
			Description: "Production environment",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Name != "production" {
			t.Errorf("expected name 'production', got %s", v.Name)
		}
		if v.Description != "Production environment" {
			t.Errorf("expected description 'Production environment', got %s", v.Description)
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), "", CreateTagValueRequest{Name: "production"})
		if err == nil {
			t.Fatal("expected error for empty tag_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "tag_name" {
			t.Errorf("expected field 'tag_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error when value name is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), "environment", CreateTagValueRequest{Name: ""})
		if err == nil {
			t.Fatal("expected error for empty value name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "name" {
			t.Errorf("expected field 'name', got %s", valErr.Field)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newValuesErrorService()
		defer closeServer()
		_, err := svc.Create(context.Background(), "my-tag", CreateTagValueRequest{Name: "v"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagValueService_Update(t *testing.T) {
	ts, svc := newValuesTestServer()
	defer ts.Close()

	t.Run("updates a tag value successfully", func(t *testing.T) {
		desc := "New description"
		v, err := svc.Update(context.Background(), "environment", "production", UpdateTagValueRequest{
			Description: &desc,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Name != "production" {
			t.Errorf("expected name 'production', got %s", v.Name)
		}
		if v.Description != "New description" {
			t.Errorf("expected description 'New description', got %s", v.Description)
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.Update(context.Background(), "", "production", UpdateTagValueRequest{})
		if err == nil {
			t.Fatal("expected error for empty tag_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "tag_name" {
			t.Errorf("expected field 'tag_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error when value_name is empty", func(t *testing.T) {
		_, err := svc.Update(context.Background(), "environment", "", UpdateTagValueRequest{})
		if err == nil {
			t.Fatal("expected error for empty value_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "value_name" {
			t.Errorf("expected field 'value_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newValuesErrorService()
		defer closeServer()
		_, err := svc.Update(context.Background(), "my-tag", "my-value", UpdateTagValueRequest{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagValueService_Delete(t *testing.T) {
	ts, svc := newValuesTestServer()
	defer ts.Close()

	t.Run("deletes a tag value successfully", func(t *testing.T) {
		v, err := svc.Delete(context.Background(), "environment", "production")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Name != "production" {
			t.Errorf("expected name 'production', got %s", v.Name)
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.Delete(context.Background(), "", "production")
		if err == nil {
			t.Fatal("expected error for empty tag_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "tag_name" {
			t.Errorf("expected field 'tag_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error when value_name is empty", func(t *testing.T) {
		_, err := svc.Delete(context.Background(), "environment", "")
		if err == nil {
			t.Fatal("expected error for empty value_name")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "value_name" {
			t.Errorf("expected field 'value_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newValuesErrorService()
		defer closeServer()
		_, err := svc.Delete(context.Background(), "my-tag", "my-value")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// prevent "imported and not used" for fmt
var _ = fmt.Sprintf
