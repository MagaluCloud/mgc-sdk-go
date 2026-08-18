package tag

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestResourceTypeService_List(t *testing.T) {
	t.Parallel()

	t.Run("resource types with their products", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resource-types", r.URL.Path)
			assertEqual(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": [
				{
					"name": "vm.instance",
					"product": "virtual-machine",
					"created_at": "2025-06-20T18:36:18.919454",
					"updated_at": null
				},
				{
					"name": "os.bucket",
					"product": "object-storage",
					"created_at": "2025-06-20T18:36:18.919454",
					"updated_at": "2025-06-21T18:36:18.919454"
				}
			]}`))
		}))
		defer server.Close()

		resourceTypes, err := testTagClient(server.URL).ResourceTypes().
			List(context.Background(), ListResourceTypesOptions{})

		assertNoError(t, err)
		assertEqual(t, 2, len(resourceTypes))
		assertEqual(t, ResourceTypeName("vm.instance"), resourceTypes[0].Name)
		assertEqual(t, Product("virtual-machine"), resourceTypes[0].Product)
		assertEqual(t, "2025-06-20T18:36:18.919454", formatTime(t, resourceTypes[0].CreatedAt))
		if resourceTypes[0].UpdatedAt != nil {
			t.Error("Expected a resource type never updated to have a nil UpdatedAt")
		}
		assertEqual(t, ResourceTypeName("os.bucket"), resourceTypes[1].Name)
		assertEqual(t, Product("object-storage"), resourceTypes[1].Product)
	})

	t.Run("filters and pagination", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			assertEqual(t, 5, len(query))
			assertEqual(t, "k8s.cluster", query.Get("name"))
			assertEqual(t, "kubernetes", query.Get("product"))
			assertEqual(t, "50", query.Get("_limit"))
			assertEqual(t, "0", query.Get("_offset"))
			assertEqual(t, "name:asc", query.Get("_sort"))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": []}`))
		}))
		defer server.Close()

		name := ResourceTypeName("k8s.cluster")
		product := Product("kubernetes")
		_, err := testTagClient(server.URL).ResourceTypes().
			List(context.Background(), ListResourceTypesOptions{
				Name:    &name,
				Product: &product,
				Limit:   helpers.IntPtr(50),
				Offset:  helpers.IntPtr(0),
				Sort:    helpers.StrPtr("name:asc"),
			})

		assertNoError(t, err)
	})

	t.Run("no filters sends no query", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, 0, len(r.URL.Query()))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": []}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).ResourceTypes().
			List(context.Background(), ListResourceTypesOptions{})

		assertNoError(t, err)
	})

	t.Run("invalid response body", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": {}}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).ResourceTypes().
			List(context.Background(), ListResourceTypesOptions{})

		assertError(t, err)
	})

	t.Run("api failure is retried and then reported", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).ResourceTypes().
			List(context.Background(), ListResourceTypesOptions{})

		assertError(t, err)
	})
}
