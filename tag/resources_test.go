package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func newResourcesTestServer() (*httptest.Server, TagValueResourceService) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// /v0/tags/{tag_name}/values/{value_name}/resources
		// /v0/tags/{tag_name}/values/{value_name}/resources/{resource_id}
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v0/tags/"), "/")
		// parts[0]=tag_name, parts[1]="values", parts[2]=value_name, parts[3]="resources", parts[4]=resource_id (opt)
		switch {
		case len(parts) == 4 && parts[1] == "values" && parts[3] == "resources" && r.Method == http.MethodGet:
			handleListResources(w, r)
		case len(parts) == 4 && parts[1] == "values" && parts[3] == "resources" && r.Method == http.MethodPost:
			handleCreateResource(w, r)
		case len(parts) == 5 && parts[1] == "values" && parts[3] == "resources" && r.Method == http.MethodGet:
			handleGetResource(w, r, parts[4])
		case len(parts) == 5 && parts[1] == "values" && parts[3] == "resources" && r.Method == http.MethodDelete:
			handleDeleteResource(w, r, parts[4])
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))

	cfg := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithTimeout(20*time.Second),
	)
	c := New(cfg, WithBasePath(client.MgcUrl(ts.URL)))
	return ts, c.Resources()
}

func newResourcesErrorService() (TagValueResourceService, func()) {
	ts := errorServer()
	cfg := client.NewMgcClient(client.WithAPIKey("test-api-key"))
	svc := New(cfg, WithBasePath(client.MgcUrl(ts.URL))).Resources()
	return svc, ts.Close
}

// --- handlers ---

func handleListResources(w http.ResponseWriter, r *http.Request) {
	response := TagValueResourceListResponse{
		Results: []TagValueResource{
			{
				ResourceID:     "vm-abc123",
				ResourceTypeID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				Region:         []RegionEnum{RegionBrSe1},
				CreatedAt:      "2024-01-01T00:00:00Z",
			},
			{
				ResourceID:     "vol-xyz789",
				ResourceTypeID: "ffffffff-0000-1111-2222-333333333333",
				Region:         []RegionEnum{RegionBrNe1},
				CreatedAt:      "2024-01-02T00:00:00Z",
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleCreateResource(w http.ResponseWriter, r *http.Request) {
	var req CreateTagValueResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	response := TagValueResource{
		ResourceID:     req.ResourceID,
		ResourceTypeID: req.ResourceTypeID,
		Region:         req.Region,
		CreatedAt:      "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetResource(w http.ResponseWriter, r *http.Request, resourceID string) {
	response := TagValueResource{
		ResourceID:     resourceID,
		ResourceTypeID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Region:         []RegionEnum{RegionBrSe1},
		CreatedAt:      "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteResource(w http.ResponseWriter, r *http.Request, resourceID string) {
	response := TagValueResource{
		ResourceID:     resourceID,
		ResourceTypeID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Region:         []RegionEnum{RegionBrSe1},
		CreatedAt:      "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// --- tests ---

func TestTagValueResourceService_List(t *testing.T) {
	ts, svc := newResourcesTestServer()
	defer ts.Close()

	t.Run("lists all resources for a tag value", func(t *testing.T) {
		resources, err := svc.List(context.Background(), "environment", "production", ListTagValueResourcesOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resources) != 2 {
			t.Fatalf("expected 2 resources, got %d", len(resources))
		}
		if resources[0].ResourceID != "vm-abc123" {
			t.Errorf("expected 'vm-abc123', got %s", resources[0].ResourceID)
		}
	})

	t.Run("lists resources with filters", func(t *testing.T) {
		region := RegionBrSe1
		resources, err := svc.List(context.Background(), "environment", "production", ListTagValueResourcesOptions{
			Region:     &region,
			ResourceID: helpers.StrPtr("vm-abc123"),
			Limit:      helpers.IntPtr(10),
			Offset:     helpers.IntPtr(0),
			Sort:       helpers.StrPtr("resource_id"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resources) == 0 {
			t.Error("expected at least one resource")
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.List(context.Background(), "", "production", ListTagValueResourcesOptions{})
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
		_, err := svc.List(context.Background(), "environment", "", ListTagValueResourcesOptions{})
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
		svc, closeServer := newResourcesErrorService()
		defer closeServer()
		_, err := svc.List(context.Background(), "tag", "value", ListTagValueResourcesOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagValueResourceService_Get(t *testing.T) {
	ts, svc := newResourcesTestServer()
	defer ts.Close()

	t.Run("gets a resource successfully", func(t *testing.T) {
		r, err := svc.Get(context.Background(), "environment", "production", "vm-abc123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.ResourceID != "vm-abc123" {
			t.Errorf("expected 'vm-abc123', got %s", r.ResourceID)
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.Get(context.Background(), "", "production", "vm-abc123")
		assertValidationError(t, err, "tag_name")
	})

	t.Run("returns error when value_name is empty", func(t *testing.T) {
		_, err := svc.Get(context.Background(), "environment", "", "vm-abc123")
		assertValidationError(t, err, "value_name")
	})

	t.Run("returns error when resource_id is empty", func(t *testing.T) {
		_, err := svc.Get(context.Background(), "environment", "production", "")
		assertValidationError(t, err, "resource_id")
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newResourcesErrorService()
		defer closeServer()
		_, err := svc.Get(context.Background(), "tag", "value", "res-id")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagValueResourceService_Create(t *testing.T) {
	ts, svc := newResourcesTestServer()
	defer ts.Close()

	t.Run("creates a resource link successfully", func(t *testing.T) {
		r, err := svc.Create(context.Background(), "environment", "production", CreateTagValueResourceRequest{
			ResourceID:     "vm-abc123",
			ResourceTypeID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			Region:         []RegionEnum{RegionBrSe1},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.ResourceID != "vm-abc123" {
			t.Errorf("expected 'vm-abc123', got %s", r.ResourceID)
		}
		if len(r.Region) != 1 || r.Region[0] != RegionBrSe1 {
			t.Errorf("expected region [%s], got %v", RegionBrSe1, r.Region)
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), "", "production", CreateTagValueResourceRequest{
			ResourceID: "vm-abc123", ResourceTypeID: "uuid", Region: []RegionEnum{RegionBrSe1},
		})
		assertValidationError(t, err, "tag_name")
	})

	t.Run("returns error when value_name is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), "environment", "", CreateTagValueResourceRequest{
			ResourceID: "vm-abc123", ResourceTypeID: "uuid", Region: []RegionEnum{RegionBrSe1},
		})
		assertValidationError(t, err, "value_name")
	})

	t.Run("returns error when resource_id is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), "environment", "production", CreateTagValueResourceRequest{
			ResourceID: "", ResourceTypeID: "uuid", Region: []RegionEnum{RegionBrSe1},
		})
		assertValidationError(t, err, "resource_id")
	})

	t.Run("returns error when resource_type_id is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), "environment", "production", CreateTagValueResourceRequest{
			ResourceID: "vm-abc123", ResourceTypeID: "", Region: []RegionEnum{RegionBrSe1},
		})
		assertValidationError(t, err, "resource_type_id")
	})

	t.Run("returns error when region is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), "environment", "production", CreateTagValueResourceRequest{
			ResourceID: "vm-abc123", ResourceTypeID: "uuid", Region: []RegionEnum{},
		})
		assertValidationError(t, err, "region")
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newResourcesErrorService()
		defer closeServer()
		_, err := svc.Create(context.Background(), "tag", "value", CreateTagValueResourceRequest{
			ResourceID: "r", ResourceTypeID: "t", Region: []RegionEnum{RegionBrSe1},
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagValueResourceService_Delete(t *testing.T) {
	ts, svc := newResourcesTestServer()
	defer ts.Close()

	t.Run("deletes a resource link successfully", func(t *testing.T) {
		r, err := svc.Delete(context.Background(), "environment", "production", "vm-abc123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.ResourceID != "vm-abc123" {
			t.Errorf("expected 'vm-abc123', got %s", r.ResourceID)
		}
	})

	t.Run("returns error when tag_name is empty", func(t *testing.T) {
		_, err := svc.Delete(context.Background(), "", "production", "vm-abc123")
		assertValidationError(t, err, "tag_name")
	})

	t.Run("returns error when value_name is empty", func(t *testing.T) {
		_, err := svc.Delete(context.Background(), "environment", "", "vm-abc123")
		assertValidationError(t, err, "value_name")
	})

	t.Run("returns error when resource_id is empty", func(t *testing.T) {
		_, err := svc.Delete(context.Background(), "environment", "production", "")
		assertValidationError(t, err, "resource_id")
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newResourcesErrorService()
		defer closeServer()
		_, err := svc.Delete(context.Background(), "tag", "value", "res-id")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// assertValidationError is a helper to check for ValidationError on a specific field
func assertValidationError(t *testing.T, err error, field string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected ValidationError for field '%s', got nil", field)
	}
	var valErr *client.ValidationError
	if !isValidationError(err, &valErr) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Field != field {
		t.Errorf("expected field '%s', got '%s'", field, valErr.Field)
	}
}
