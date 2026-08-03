package tag

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

const tagValuePayload = `{
	"name": "test-labs",
	"description": "tag value to monitor expenses with test-labs",
	"created_at": "2025-06-20T18:36:18.919454",
	"updated_at": "2025-06-21T18:36:18.919454",
	"tag": {
		"name": "kubernetes-expenses",
		"description": "tag to monitor expenses with environments"
	}
}`

func TestTagValueService_List(t *testing.T) {
	t.Parallel()

	t.Run("values of a tag", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags/kubernetes-expenses/values", r.URL.Path)
			assertEqual(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": [` + tagValuePayload + `]}`))
		}))
		defer server.Close()

		values, err := testTagClient(server.URL).Values().
			List(context.Background(), "kubernetes-expenses", ListTagValuesOptions{})

		assertNoError(t, err)
		assertEqual(t, 1, len(values))
		assertEqual(t, "test-labs", values[0].Name)
		assertEqual(t, "tag value to monitor expenses with test-labs", *values[0].Description)
		assertEqual(t, "2025-06-20T18:36:18.919454", formatTime(t, values[0].CreatedAt))
		assertEqual(t, "kubernetes-expenses", values[0].Tag.Name)
	})

	t.Run("filters and pagination", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			assertEqual(t, 4, len(query))
			assertEqual(t, "test-labs", query.Get("name"))
			assertEqual(t, "5", query.Get("_limit"))
			assertEqual(t, "10", query.Get("_offset"))
			assertEqual(t, "name:asc", query.Get("_sort"))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": []}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Values().
			List(context.Background(), "kubernetes-expenses", ListTagValuesOptions{
				Name:   helpers.StrPtr("test-labs"),
				Limit:  helpers.IntPtr(5),
				Offset: helpers.IntPtr(10),
				Sort:   helpers.StrPtr("name:asc"),
			})

		assertNoError(t, err)
	})

	t.Run("empty tag name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Values().
			List(context.Background(), "", ListTagValuesOptions{})

		assertValidationError(t, err, "tag_name")
	})

	t.Run("invalid response body", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`not json`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Values().
			List(context.Background(), "kubernetes-expenses", ListTagValuesOptions{})

		assertError(t, err)
	})

	t.Run("api failure is retried and then reported", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Values().
			List(context.Background(), "kubernetes-expenses", ListTagValuesOptions{})

		assertError(t, err)
	})
}

func TestTagValueService_Get(t *testing.T) {
	t.Parallel()

	t.Run("value with the tag that owns it", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags/kubernetes-expenses/values/test-labs", r.URL.Path)
			assertEqual(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tagValuePayload))
		}))
		defer server.Close()

		value, err := testTagClient(server.URL).Values().
			Get(context.Background(), "kubernetes-expenses", "test-labs")

		assertNoError(t, err)
		assertEqual(t, "test-labs", value.Name)
		assertEqual(t, "kubernetes-expenses", value.Tag.Name)
		assertEqual(t, "tag to monitor expenses with environments", *value.Tag.Description)
	})

	t.Run("names with characters that need escaping", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags/my tag/values/prod (main)", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tagValuePayload))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Values().
			Get(context.Background(), "my tag", "prod (main)")

		assertNoError(t, err)
	})

	t.Run("empty tag name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Values().
			Get(context.Background(), "", "test-labs")

		assertValidationError(t, err, "tag_name")
	})

	t.Run("empty value name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Values().
			Get(context.Background(), "kubernetes-expenses", "")

		assertValidationError(t, err, "value_name")
	})

	t.Run("non-existent value", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"detail": "value not found"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Values().
			Get(context.Background(), "kubernetes-expenses", "missing")

		assertError(t, err)
	})
}

func TestTagValueService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      CreateTagValueRequest
		wantBody string
	}{
		{
			name:     "value with description",
			req:      CreateTagValueRequest{Name: "test-labs", Description: helpers.StrPtr("labs")},
			wantBody: `{"name": "test-labs", "description": "labs"}`,
		},
		{
			name:     "description is optional",
			req:      CreateTagValueRequest{Name: "test-labs"},
			wantBody: `{"name": "test-labs"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/tags/v0/tags/kubernetes-expenses/values", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				body, err := io.ReadAll(r.Body)
				assertNoError(t, err)
				assertEqual(t, canonicalJSON(t, tt.wantBody), canonicalJSON(t, string(body)))

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tagValuePayload))
			}))
			defer server.Close()

			value, err := testTagClient(server.URL).Values().
				Create(context.Background(), "kubernetes-expenses", tt.req)

			assertNoError(t, err)
			assertEqual(t, "test-labs", value.Name)
		})
	}

	t.Run("empty tag name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Values().
			Create(context.Background(), "", CreateTagValueRequest{Name: "test-labs"})

		assertValidationError(t, err, "tag_name")
	})

	t.Run("empty value name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Values().
			Create(context.Background(), "kubernetes-expenses", CreateTagValueRequest{})

		assertValidationError(t, err, "name")
	})

	t.Run("name already taken within the tag", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"detail": "value already exists"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Values().
			Create(context.Background(), "kubernetes-expenses", CreateTagValueRequest{Name: "test-labs"})

		assertError(t, err)
	})
}

func TestTagValueService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      UpdateTagValueRequest
		wantBody string
	}{
		{
			name:     "new description",
			req:      UpdateTagValueRequest{Description: helpers.StrPtr("new description")},
			wantBody: `{"description": "new description"}`,
		},
		{
			name:     "nil description clears it",
			req:      UpdateTagValueRequest{},
			wantBody: `{"description": null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/tags/v0/tags/kubernetes-expenses/values/test-labs", r.URL.Path)
				assertEqual(t, http.MethodPatch, r.Method)

				body, err := io.ReadAll(r.Body)
				assertNoError(t, err)
				assertEqual(t, canonicalJSON(t, tt.wantBody), canonicalJSON(t, string(body)))

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tagValuePayload))
			}))
			defer server.Close()

			value, err := testTagClient(server.URL).Values().
				Update(context.Background(), "kubernetes-expenses", "test-labs", tt.req)

			assertNoError(t, err)
			assertEqual(t, "test-labs", value.Name)
		})
	}

	t.Run("empty tag name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Values().
			Update(context.Background(), "", "test-labs", UpdateTagValueRequest{})

		assertValidationError(t, err, "tag_name")
	})

	t.Run("empty value name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Values().
			Update(context.Background(), "kubernetes-expenses", "", UpdateTagValueRequest{})

		assertValidationError(t, err, "value_name")
	})
}

func TestTagValueService_Delete(t *testing.T) {
	t.Parallel()

	t.Run("value is removed", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags/kubernetes-expenses/values/test-labs", r.URL.Path)
			assertEqual(t, http.MethodDelete, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tagValuePayload))
		}))
		defer server.Close()

		err := testTagClient(server.URL).Values().
			Delete(context.Background(), "kubernetes-expenses", "test-labs")

		assertNoError(t, err)
	})

	t.Run("empty tag name", func(t *testing.T) {
		t.Parallel()

		err := testTagClient("https://tags.example.com").Values().
			Delete(context.Background(), "", "test-labs")

		assertValidationError(t, err, "tag_name")
	})

	t.Run("empty value name", func(t *testing.T) {
		t.Parallel()

		err := testTagClient("https://tags.example.com").Values().
			Delete(context.Background(), "kubernetes-expenses", "")

		assertValidationError(t, err, "value_name")
	})

	t.Run("non-existent value", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"detail": "value not found"}`))
		}))
		defer server.Close()

		err := testTagClient(server.URL).Values().
			Delete(context.Background(), "kubernetes-expenses", "missing")

		assertError(t, err)
	})
}
