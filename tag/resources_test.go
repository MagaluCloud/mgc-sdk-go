package tag

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

const resourcePayload = `{
	"id": "31201f93-f5f4-4cf1-ba9c-bfed0717f4ac",
	"external_id": "d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b",
	"resource_type": {
		"name": "vm.instance",
		"product": "virtual-machine"
	},
	"region": "br-se1",
	"created_at": "2025-06-20T18:36:18.919454",
	"updated_at": null,
	"last_tag_associated_at": "2025-06-23T18:36:18.919454",
	"tags": [
		{"name": "kubernetes-expenses", "value": "test-labs"},
		{"name": "environment", "value": "production"}
	]
}`

func TestResourceService_List(t *testing.T) {
	t.Parallel()

	t.Run("tagged resources with their tags", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resources", r.URL.Path)
			assertEqual(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": [` + resourcePayload + `]}`))
		}))
		defer server.Close()

		resources, err := testTagClient(server.URL).Resources().
			List(context.Background(), ListResourcesOptions{})

		assertNoError(t, err)
		assertEqual(t, 1, len(resources))

		resource := resources[0]
		assertEqual(t, "31201f93-f5f4-4cf1-ba9c-bfed0717f4ac", resource.ID)
		assertEqual(t, "d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b", resource.ExternalID)
		assertEqual(t, ResourceTypeName("vm.instance"), resource.ResourceType.Name)
		assertEqual(t, Product("virtual-machine"), resource.ResourceType.Product)
		assertEqual(t, "br-se1", resource.Region)
		assertEqual(t, "2025-06-20T18:36:18.919454", formatTime(t, resource.CreatedAt))
		if resource.UpdatedAt != nil {
			t.Error("Expected a resource never updated to have a nil UpdatedAt")
		}
		if resource.LastTagAssociatedAt == nil {
			t.Fatal("Expected LastTagAssociatedAt to be set")
		}
		assertEqual(t, "2025-06-23T18:36:18.919454", formatTime(t, *resource.LastTagAssociatedAt))
		assertEqual(t, 2, len(resource.Tags))
		assertEqual(t, "kubernetes-expenses", resource.Tags[0].Name)
		assertEqual(t, "test-labs", resource.Tags[0].Value)
		assertEqual(t, "environment", resource.Tags[1].Name)
		assertEqual(t, "production", resource.Tags[1].Value)
	})

	t.Run("resource that never had a tag attached", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": [{
				"id": "31201f93-f5f4-4cf1-ba9c-bfed0717f4ac",
				"external_id": "d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b",
				"resource_type": {"name": "os.bucket", "product": "object-storage"},
				"region": "global",
				"created_at": "2025-06-20T18:36:18.919454",
				"updated_at": null,
				"last_tag_associated_at": null,
				"tags": []
			}]}`))
		}))
		defer server.Close()

		resources, err := testTagClient(server.URL).Resources().
			List(context.Background(), ListResourcesOptions{})

		assertNoError(t, err)
		if resources[0].LastTagAssociatedAt != nil {
			t.Error("Expected LastTagAssociatedAt to be nil")
		}
		assertEqual(t, 0, len(resources[0].Tags))
	})

	t.Run("filters and pagination", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			assertEqual(t, 6, len(query))
			assertEqual(t, "d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b", query.Get("external_id"))
			assertEqual(t, "k8s.cluster", query.Get("resource_type_name"))
			assertEqual(t, "br-ne1", query.Get("region"))
			assertEqual(t, "50", query.Get("_limit"))
			assertEqual(t, "10", query.Get("_offset"))
			assertEqual(t, "external_id:asc", query.Get("_sort"))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": []}`))
		}))
		defer server.Close()

		resourceTypeName := ResourceTypeName("k8s.cluster")
		_, err := testTagClient(server.URL).Resources().
			List(context.Background(), ListResourcesOptions{
				ExternalID:       helpers.StrPtr("d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b"),
				ResourceTypeName: &resourceTypeName,
				Region:           helpers.StrPtr("br-ne1"),
				Limit:            helpers.IntPtr(50),
				Offset:           helpers.IntPtr(10),
				Sort:             helpers.StrPtr("external_id:asc"),
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

		_, err := testTagClient(server.URL).Resources().
			List(context.Background(), ListResourcesOptions{})

		assertNoError(t, err)
	})

	t.Run("invalid response body", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": {}}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Resources().
			List(context.Background(), ListResourcesOptions{})

		assertError(t, err)
	})

	t.Run("api failure is reported", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"detail": [{"loc": ["query", "region"], "msg": "invalid region"}]}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Resources().
			List(context.Background(), ListResourcesOptions{})

		assertError(t, err)
	})
}

func TestResourceService_Get(t *testing.T) {
	t.Parallel()

	t.Run("resource with its tags", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resources/d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b", r.URL.Path)
			assertEqual(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(resourcePayload))
		}))
		defer server.Close()

		resource, err := testTagClient(server.URL).Resources().
			Get(context.Background(), "d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b")

		assertNoError(t, err)
		assertEqual(t, ResourceTypeName("vm.instance"), resource.ResourceType.Name)
		assertEqual(t, 2, len(resource.Tags))
	})

	t.Run("external ID with characters that need escaping", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resources/my-bucket%2Freports%2Fjune.csv", r.URL.EscapedPath())
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(resourcePayload))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Resources().
			Get(context.Background(), "my-bucket/reports/june.csv")

		assertNoError(t, err)
	})

	t.Run("empty external ID is rejected before the request", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Resources().
			Get(context.Background(), "")

		assertValidationError(t, err, "external_id")
	})

	t.Run("api failure is reported", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"detail": "not a valid resource"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Resources().
			Get(context.Background(), "d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b")

		assertError(t, err)
	})
}

func TestResourceService_AttachTags(t *testing.T) {
	t.Parallel()

	t.Run("attaches tags and their values", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resources/d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b/tags", r.URL.Path)
			assertEqual(t, http.MethodPost, r.Method)

			body, err := io.ReadAll(r.Body)
			assertNoError(t, err)
			expected := `{"tags": [
				{"name": "kubernetes-expenses", "value": "test-labs"},
				{"name": "environment", "value": "production"}
			]}`
			assertEqual(t, canonicalJSON(t, expected), canonicalJSON(t, string(body)))

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(resourcePayload))
		}))
		defer server.Close()

		resource, err := testTagClient(server.URL).Resources().
			AttachTags(context.Background(), "d9f3a1b2-6c5d-4e8f-9a0b-1c2d3e4f5a6b", AttachTagsRequest{
				Tags: []AttachTag{
					{Name: "kubernetes-expenses", Value: "test-labs"},
					{Name: "environment", Value: "production"},
				},
			})

		assertNoError(t, err)
		assertEqual(t, 2, len(resource.Tags))
		assertEqual(t, "kubernetes-expenses", resource.Tags[0].Name)
		assertEqual(t, "test-labs", resource.Tags[0].Value)
	})

	t.Run("external ID with characters that need escaping", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resources/my-bucket%2Freports%2Fjune.csv/tags", r.URL.EscapedPath())
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(resourcePayload))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Resources().
			AttachTags(context.Background(), "my-bucket/reports/june.csv", AttachTagsRequest{
				Tags: []AttachTag{{Name: "environment", Value: "production"}},
			})

		assertNoError(t, err)
	})

	t.Run("invalid requests are rejected before the request", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Expected no request to be sent")
		}))
		defer server.Close()

		tests := []struct {
			name       string
			externalID string
			req        AttachTagsRequest
			field      string
		}{
			{
				name:       "empty external ID",
				externalID: "",
				req:        AttachTagsRequest{Tags: []AttachTag{{Name: "environment", Value: "production"}}},
				field:      "external_id",
			},
			{
				name:       "no tags",
				externalID: "d9f3a1b2",
				req:        AttachTagsRequest{},
				field:      "tags",
			},
			{
				name:       "tag without a name",
				externalID: "d9f3a1b2",
				req:        AttachTagsRequest{Tags: []AttachTag{{Value: "production"}}},
				field:      "tags[0].name",
			},
			{
				name:       "tag without a value",
				externalID: "d9f3a1b2",
				req: AttachTagsRequest{Tags: []AttachTag{
					{Name: "environment", Value: "production"},
					{Name: "kubernetes-expenses"},
				}},
				field: "tags[1].value",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := testTagClient(server.URL).Resources().
					AttachTags(context.Background(), tt.externalID, tt.req)

				assertValidationError(t, err, tt.field)
			})
		}
	})

	t.Run("api failure is reported", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"detail": "tag not found"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Resources().
			AttachTags(context.Background(), "d9f3a1b2", AttachTagsRequest{
				Tags: []AttachTag{{Name: "environment", Value: "production"}},
			})

		assertError(t, err)
	})
}

func TestResourceService_DetachTag(t *testing.T) {
	t.Parallel()

	t.Run("detaches a tag", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resources/d9f3a1b2/tags/environment", r.URL.Path)
			assertEqual(t, http.MethodDelete, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(resourcePayload))
		}))
		defer server.Close()

		err := testTagClient(server.URL).Resources().
			DetachTag(context.Background(), "d9f3a1b2", "environment")

		assertNoError(t, err)
	})

	t.Run("names with characters that need escaping", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/resources/my-bucket%2Freports/tags/my%20tag%20%5Bprod%5D", r.URL.EscapedPath())
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(resourcePayload))
		}))
		defer server.Close()

		err := testTagClient(server.URL).Resources().
			DetachTag(context.Background(), "my-bucket/reports", "my tag [prod]")

		assertNoError(t, err)
	})

	t.Run("invalid requests are rejected before the request", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Expected no request to be sent")
		}))
		defer server.Close()

		tests := []struct {
			name       string
			externalID string
			tagName    string
			field      string
		}{
			{name: "empty external ID", externalID: "", tagName: "environment", field: "external_id"},
			{name: "empty tag name", externalID: "d9f3a1b2", tagName: "", field: "tag_name"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := testTagClient(server.URL).Resources().
					DetachTag(context.Background(), tt.externalID, tt.tagName)

				assertValidationError(t, err, tt.field)
			})
		}
	})

	t.Run("api failure is reported", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"detail": "tag not attached"}`))
		}))
		defer server.Close()

		err := testTagClient(server.URL).Resources().
			DetachTag(context.Background(), "d9f3a1b2", "environment")

		assertError(t, err)
	})
}
