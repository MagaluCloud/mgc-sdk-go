package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func newResourceTypesTestServer() (*httptest.Server, ResourceTypeService) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v0/resource-types" && r.Method == http.MethodGet {
			handleListResourceTypes(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))

	cfg := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithTimeout(20*time.Second),
	)
	c := New(cfg, WithBasePath(client.MgcUrl(ts.URL)))
	return ts, c.ResourceTypes()
}

func newResourceTypesErrorService() (ResourceTypeService, func()) {
	ts := errorServer()
	cfg := client.NewMgcClient(client.WithAPIKey("test-api-key"))
	svc := New(cfg, WithBasePath(client.MgcUrl(ts.URL))).ResourceTypes()
	return svc, ts.Close
}

func handleListResourceTypes(w http.ResponseWriter, r *http.Request) {
	response := ResourceTypeListResponse{
		Results: []ResourceType{
			{
				Name:      []ResourceEnum{ResourceInstancesVM},
				Product:   []ProductEnum{ProductVirtualMachine},
				CreatedAt: "2024-01-01T00:00:00Z",
			},
			{
				Name:      []ResourceEnum{ResourceVolumes},
				Product:   []ProductEnum{ProductBlockStorage},
				CreatedAt: "2024-01-02T00:00:00Z",
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func TestResourceTypeService_List(t *testing.T) {
	ts, svc := newResourceTypesTestServer()
	defer ts.Close()

	t.Run("lists all resource types successfully", func(t *testing.T) {
		types, err := svc.List(context.Background(), ListResourceTypesOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(types) != 2 {
			t.Fatalf("expected 2 resource types, got %d", len(types))
		}
		if len(types[0].Name) == 0 || types[0].Name[0] != ResourceInstancesVM {
			t.Errorf("expected first type 'instances-vm', got %v", types[0].Name)
		}
		if len(types[1].Name) == 0 || types[1].Name[0] != ResourceVolumes {
			t.Errorf("expected second type 'volumes', got %v", types[1].Name)
		}
	})

	t.Run("lists resource types with name filter", func(t *testing.T) {
		name := ResourceInstancesVM
		types, err := svc.List(context.Background(), ListResourceTypesOptions{
			Name: &name,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(types) == 0 {
			t.Error("expected at least one resource type")
		}
	})

	t.Run("lists resource types with product filter", func(t *testing.T) {
		product := ProductVirtualMachine
		types, err := svc.List(context.Background(), ListResourceTypesOptions{
			Product: &product,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(types) == 0 {
			t.Error("expected at least one resource type")
		}
	})

	t.Run("lists resource types with pagination", func(t *testing.T) {
		types, err := svc.List(context.Background(), ListResourceTypesOptions{
			Limit:  helpers.IntPtr(10),
			Offset: helpers.IntPtr(0),
			Sort:   helpers.StrPtr("name"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(types) == 0 {
			t.Error("expected at least one resource type")
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newResourceTypesErrorService()
		defer closeServer()
		_, err := svc.List(context.Background(), ListResourceTypesOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
