package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

// Helper functions
func assertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v but got %v. %v", expected, actual, msgAndArgs)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

// assertValidationError checks that err is a *client.ValidationError for the given API field.
func assertValidationError(t *testing.T, err error, field string) {
	t.Helper()
	validationErr, ok := err.(*client.ValidationError)
	if !ok {
		t.Fatalf("Expected *client.ValidationError but got %T: %v", err, err)
	}
	assertEqual(t, field, validationErr.Field)
}

// canonicalJSON re-encodes a JSON document so two encodings of the same object compare equal.
func canonicalJSON(t *testing.T, raw string) string {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatalf("invalid JSON %q: %v", raw, err)
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("could not re-encode JSON: %v", err)
	}
	return string(out)
}

func testTagClient(baseURL string) *TagClient {
	core := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithHTTPClient(&http.Client{}),
	)
	return New(core, WithGlobalBasePath(client.MgcUrl(baseURL)))
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("nil core returns nil", func(t *testing.T) {
		t.Parallel()

		if New(nil) != nil {
			t.Error("Expected nil client for nil core")
		}
	})

	t.Run("defaults to the global endpoint", func(t *testing.T) {
		t.Parallel()

		core := client.NewMgcClient(client.WithAPIKey("test-api-key"))
		tagClient := New(core)

		req, err := tagClient.newRequest(context.Background(), http.MethodGet, "/v0/tags", nil)

		assertNoError(t, err)
		assertEqual(t, client.Global.String()+"/tags/v0/tags", req.URL.String())
		assertEqual(t, core, tagClient.CoreClient)
	})

	t.Run("global base path can be overridden", func(t *testing.T) {
		t.Parallel()

		core := client.NewMgcClient(client.WithAPIKey("test-api-key"))
		tagClient := New(core, WithGlobalBasePath("https://tags.example.com"))

		req, err := tagClient.newRequest(context.Background(), http.MethodGet, "/v0/tags", nil)

		assertNoError(t, err)
		assertEqual(t, "https://tags.example.com/tags/v0/tags", req.URL.String())
	})

	t.Run("keeps the endpoint of the core client, which is shared with other services", func(t *testing.T) {
		t.Parallel()

		core := client.NewMgcClient(client.WithBaseURL(client.BrNe1))

		tagClient := New(core)
		if _, err := tagClient.newRequest(context.Background(), http.MethodGet, "/v0/tags", nil); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}

		assertEqual(t, client.BrNe1, core.GetConfig().BaseURL)
	})

	t.Run("keeps the credentials and transport configured by the caller", func(t *testing.T) {
		t.Parallel()

		core := client.NewMgcClient(
			client.WithJWToken("Bearer caller-token"),
			client.WithCustomHeader("x-tenant-id", "caller-tenant"),
		)
		callerTransport := core.GetConfig().HTTPClient.Transport

		New(core)

		assertEqual(t, "Bearer caller-token", core.GetConfig().JWToken)
		assertEqual(t, "caller-tenant", core.GetConfig().CustomHeaders["x-tenant-id"])
		assertEqual(t, callerTransport, core.GetConfig().HTTPClient.Transport)
	})
}

func TestTagClient_newRequest(t *testing.T) {
	t.Parallel()

	t.Run("prefixes the path with the tags base path", func(t *testing.T) {
		t.Parallel()

		req, err := testTagClient("https://tags.example.com").
			newRequest(context.Background(), http.MethodGet, "/v0/tags", nil)

		assertNoError(t, err)
		assertEqual(t, "https://tags.example.com/tags/v0/tags", req.URL.String())
		assertEqual(t, "test-api-key", req.Header.Get("X-API-Key"))
		assertEqual(t, "application/json", req.Header.Get("Content-Type"))
	})

	t.Run("sends the tenant configured in the client", func(t *testing.T) {
		t.Parallel()

		core := client.NewMgcClient(client.WithCustomHeader("x-tenant-id", "core-tenant"))
		req, err := New(core, WithTenantID("tag-tenant")).
			newRequest(context.Background(), http.MethodGet, "/v0/tags", nil)

		assertNoError(t, err)
		assertEqual(t, "tag-tenant", req.Header.Get("x-tenant-id"))
	})

	t.Run("without a tenant the header of the core client is kept", func(t *testing.T) {
		t.Parallel()

		core := client.NewMgcClient(client.WithCustomHeader("x-tenant-id", "core-tenant"))
		req, err := New(core).newRequest(context.Background(), http.MethodGet, "/v0/tags", nil)

		assertNoError(t, err)
		assertEqual(t, "core-tenant", req.Header.Get("x-tenant-id"))
	})

	t.Run("body that cannot be encoded returns error", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").
			newRequest(context.Background(), http.MethodPost, "/v0/tags", make(chan int))

		assertError(t, err)
	})
}

func TestTagClient_Services(t *testing.T) {
	t.Parallel()

	tagClient := testTagClient("https://tags.example.com")

	if _, ok := tagClient.Tags().(*tagService); !ok {
		t.Error("Tags() did not return a *tagService")
	}
	if _, ok := tagClient.Values().(*tagValueService); !ok {
		t.Error("Values() did not return a *tagValueService")
	}
	if _, ok := tagClient.ResourceTypes().(*resourceTypeService); !ok {
		t.Error("ResourceTypes() did not return a *resourceTypeService")
	}
	if _, ok := tagClient.Resources().(*resourceService); !ok {
		t.Error("Resources() did not return a *resourceService")
	}
}
