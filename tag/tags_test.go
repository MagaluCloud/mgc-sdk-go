package tag

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

const tagPayload = `{
	"name": "kubernetes-expenses",
	"color": "f54927",
	"kinds": ["finops"],
	"description": "tag to monitor expenses with environments",
	"created_at": "2025-06-20T18:36:18.919454",
	"updated_at": "2025-06-21T18:36:18.919454",
	"values": [
		{
			"name": "test-labs",
			"description": "tag value to monitor expenses with test-labs",
			"created_at": "2025-06-20T18:36:18.919454",
			"updated_at": null
		}
	]
}`

func formatTime(t *testing.T, value utils.LocalDateTimeWithoutZone) string {
	t.Helper()
	return time.Time(value).Format(utils.LocalDateTimeWithoutZoneLayout)
}

func TestTagService_List(t *testing.T) {
	t.Parallel()

	t.Run("tags of the tenant with their values", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags", r.URL.Path)
			assertEqual(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": [` + tagPayload + `]}`))
		}))
		defer server.Close()

		tags, err := testTagClient(server.URL).Tags().List(context.Background(), ListTagsOptions{})

		assertNoError(t, err)
		assertEqual(t, 1, len(tags))
		assertEqual(t, "kubernetes-expenses", tags[0].Name)
		assertEqual(t, "f54927", *tags[0].Color)
		assertEqual(t, "tag to monitor expenses with environments", *tags[0].Description)
		assertEqual(t, 1, len(tags[0].Kinds))
		assertEqual(t, TagKindFinops, tags[0].Kinds[0])
		assertEqual(t, "2025-06-20T18:36:18.919454", formatTime(t, tags[0].CreatedAt))
		assertEqual(t, "2025-06-21T18:36:18.919454", formatTime(t, *tags[0].UpdatedAt))
		assertEqual(t, 1, len(tags[0].Values))
		assertEqual(t, "test-labs", tags[0].Values[0].Name)
		if tags[0].Values[0].UpdatedAt != nil {
			t.Error("Expected a value never updated to have a nil UpdatedAt")
		}
		if tags[0].Values[0].Tag != nil {
			t.Error("Expected a value embedded in a tag to carry no Tag reference")
		}
	})

	t.Run("tag without description or color", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": [{
				"name": "no-metadata",
				"color": null,
				"description": null,
				"kinds": [],
				"values": [],
				"created_at": "2025-06-20T18:36:18.919454",
				"updated_at": null
			}]}`))
		}))
		defer server.Close()

		tags, err := testTagClient(server.URL).Tags().List(context.Background(), ListTagsOptions{})

		assertNoError(t, err)
		assertEqual(t, 1, len(tags))
		if tags[0].Description != nil || tags[0].Color != nil {
			t.Error("Expected null description and color to decode as nil, not empty string")
		}
	})

	t.Run("kind the SDK does not know yet", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": [{
				"name": "future",
				"kinds": ["secops"],
				"values": [],
				"description": null,
				"created_at": "2025-06-20T18:36:18.919454",
				"updated_at": null
			}]}`))
		}))
		defer server.Close()

		tags, err := testTagClient(server.URL).Tags().List(context.Background(), ListTagsOptions{})

		assertNoError(t, err)
		assertEqual(t, TagKind("secops"), tags[0].Kinds[0])
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": []}`))
		}))
		defer server.Close()

		tags, err := testTagClient(server.URL).Tags().List(context.Background(), ListTagsOptions{})

		assertNoError(t, err)
		assertEqual(t, 0, len(tags))
	})

	t.Run("invalid response body", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"results": `))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Tags().List(context.Background(), ListTagsOptions{})

		assertError(t, err)
	})

	t.Run("api rejects the request", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"detail": [{"msg": "invalid color"}]}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Tags().List(context.Background(), ListTagsOptions{})

		assertError(t, err)
	})

	t.Run("api failure is retried and then reported", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Tags().List(context.Background(), ListTagsOptions{})

		assertError(t, err)
	})
}

func TestTagService_List_Filters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		opts      ListTagsOptions
		wantQuery map[string][]string
	}{
		{
			name:      "no filters sends no query",
			opts:      ListTagsOptions{},
			wantQuery: map[string][]string{},
		},
		{
			name: "pagination",
			opts: ListTagsOptions{
				Limit:  helpers.IntPtr(10),
				Offset: helpers.IntPtr(20),
				Sort:   helpers.StrPtr("name:asc"),
			},
			wantQuery: map[string][]string{
				"_limit":  {"10"},
				"_offset": {"20"},
				"_sort":   {"name:asc"},
			},
		},
		{
			name: "name and color",
			opts: ListTagsOptions{
				Name:  helpers.StrPtr("kubernetes-expenses"),
				Color: helpers.StrPtr("F54927"),
			},
			wantQuery: map[string][]string{
				"name":  {"kubernetes-expenses"},
				"color": {"f54927"},
			},
		},
		{
			name:      "each kind is sent as its own parameter",
			opts:      ListTagsOptions{Kinds: []TagKind{TagKindFinops, "secops"}},
			wantQuery: map[string][]string{"kinds": {"finops", "secops"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				assertEqual(t, len(tt.wantQuery), len(query))
				for param, want := range tt.wantQuery {
					got := query[param]
					assertEqual(t, len(want), len(got), param)
					for i := range want {
						if i < len(got) {
							assertEqual(t, want[i], got[i], param)
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"results": []}`))
			}))
			defer server.Close()

			_, err := testTagClient(server.URL).Tags().List(context.Background(), tt.opts)

			assertNoError(t, err)
		})
	}

	t.Run("invalid color is rejected before the request", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Expected no request for an invalid color")
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Tags().
			List(context.Background(), ListTagsOptions{Color: helpers.StrPtr("#f54927")})

		assertValidationError(t, err, "color")
	})
}

func TestTagService_Get(t *testing.T) {
	t.Parallel()

	t.Run("tag with its values", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags/kubernetes-expenses", r.URL.Path)
			assertEqual(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tagPayload))
		}))
		defer server.Close()

		tag, err := testTagClient(server.URL).Tags().Get(context.Background(), "kubernetes-expenses")

		assertNoError(t, err)
		assertEqual(t, "kubernetes-expenses", tag.Name)
		assertEqual(t, 1, len(tag.Values))
	})

	t.Run("name with characters that need escaping", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags/my tag [prod]", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tagPayload))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Tags().Get(context.Background(), "my tag [prod]")

		assertNoError(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Tags().Get(context.Background(), "")

		assertValidationError(t, err, "name")
	})

	t.Run("non-existent tag", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"detail": "tag not found"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Tags().Get(context.Background(), "missing")

		assertError(t, err)
	})
}

func TestTagService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      CreateTagRequest
		wantBody string
	}{
		{
			name: "tag with every field",
			req: CreateTagRequest{
				Name:        "kubernetes-expenses",
				Description: helpers.StrPtr("tag to monitor expenses with environments"),
				Color:       helpers.StrPtr("F54927"),
				Kinds:       []TagKind{TagKindFinops},
				Values: []CreateTagValueRequest{
					{Name: "test-labs", Description: helpers.StrPtr("labs")},
				},
			},
			wantBody: `{
				"name": "kubernetes-expenses",
				"description": "tag to monitor expenses with environments",
				"color": "f54927",
				"kinds": ["finops"],
				"values": [{"name": "test-labs", "description": "labs"}]
			}`,
		},
		{
			name:     "only the name is required",
			req:      CreateTagRequest{Name: "minimal"},
			wantBody: `{"name": "minimal"}`,
		},
		{
			name: "value without description",
			req: CreateTagRequest{
				Name:   "with-values",
				Values: []CreateTagValueRequest{{Name: "test-labs"}},
			},
			wantBody: `{"name": "with-values", "values": [{"name": "test-labs"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/tags/v0/tags", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				body, err := io.ReadAll(r.Body)
				assertNoError(t, err)
				assertEqual(t, canonicalJSON(t, tt.wantBody), canonicalJSON(t, string(body)))

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tagPayload))
			}))
			defer server.Close()

			tag, err := testTagClient(server.URL).Tags().Create(context.Background(), tt.req)

			assertNoError(t, err)
			assertEqual(t, "kubernetes-expenses", tag.Name)
		})
	}

	validations := []struct {
		name      string
		req       CreateTagRequest
		wantField string
	}{
		{
			name:      "empty name",
			req:       CreateTagRequest{},
			wantField: "name",
		},
		{
			name:      "invalid color",
			req:       CreateTagRequest{Name: "tag", Color: helpers.StrPtr("red")},
			wantField: "color",
		},
		{
			name:      "value without name",
			req:       CreateTagRequest{Name: "tag", Values: []CreateTagValueRequest{{Name: ""}}},
			wantField: "values[0].name",
		},
	}

	for _, tt := range validations {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Expected no request for invalid input")
			}))
			defer server.Close()

			_, err := testTagClient(server.URL).Tags().Create(context.Background(), tt.req)

			assertValidationError(t, err, tt.wantField)
		})
	}

	t.Run("name already taken", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"detail": "tag already exists"}`))
		}))
		defer server.Close()

		_, err := testTagClient(server.URL).Tags().
			Create(context.Background(), CreateTagRequest{Name: "duplicated"})

		assertError(t, err)
	})
}

func TestTagService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      UpdateTagRequest
		wantBody string
	}{
		{
			name: "every field",
			req: UpdateTagRequest{
				Description: helpers.StrPtr("new description"),
				Color:       helpers.StrPtr("FFFFFF"),
				Kinds:       &[]TagKind{TagKindFinops},
			},
			wantBody: `{"description": "new description", "color": "ffffff", "kinds": ["finops"]}`,
		},
		{
			name:     "only the fields that were set",
			req:      UpdateTagRequest{Description: helpers.StrPtr("new description")},
			wantBody: `{"description": "new description"}`,
		},
		{
			name:     "empty description clears it",
			req:      UpdateTagRequest{Description: helpers.StrPtr("")},
			wantBody: `{"description": null}`,
		},
		{
			name:     "empty color clears it",
			req:      UpdateTagRequest{Color: helpers.StrPtr("")},
			wantBody: `{"color": null}`,
		},
		{
			name:     "empty kinds clears the list",
			req:      UpdateTagRequest{Kinds: &[]TagKind{}},
			wantBody: `{"kinds": []}`,
		},
		{
			name:     "nothing to update sends an empty object",
			req:      UpdateTagRequest{},
			wantBody: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/tags/v0/tags/kubernetes-expenses", r.URL.Path)
				assertEqual(t, http.MethodPatch, r.Method)

				body, err := io.ReadAll(r.Body)
				assertNoError(t, err)
				assertEqual(t, canonicalJSON(t, tt.wantBody), canonicalJSON(t, string(body)))

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tagPayload))
			}))
			defer server.Close()

			tag, err := testTagClient(server.URL).Tags().
				Update(context.Background(), "kubernetes-expenses", tt.req)

			assertNoError(t, err)
			assertEqual(t, "kubernetes-expenses", tag.Name)
		})
	}

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Tags().
			Update(context.Background(), "", UpdateTagRequest{})

		assertValidationError(t, err, "name")
	})

	t.Run("invalid color", func(t *testing.T) {
		t.Parallel()

		_, err := testTagClient("https://tags.example.com").Tags().
			Update(context.Background(), "tag", UpdateTagRequest{Color: helpers.StrPtr("f5492")})

		assertValidationError(t, err, "color")
	})
}

func TestTagService_Delete(t *testing.T) {
	t.Parallel()

	t.Run("tag is removed", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertEqual(t, "/tags/v0/tags/kubernetes-expenses", r.URL.Path)
			assertEqual(t, http.MethodDelete, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tagPayload))
		}))
		defer server.Close()

		err := testTagClient(server.URL).Tags().Delete(context.Background(), "kubernetes-expenses")

		assertNoError(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()

		err := testTagClient("https://tags.example.com").Tags().Delete(context.Background(), "")

		assertValidationError(t, err, "name")
	})

	t.Run("non-existent tag", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"detail": "tag not found"}`))
		}))
		defer server.Close()

		err := testTagClient(server.URL).Tags().Delete(context.Background(), "missing")

		assertError(t, err)
	})
}

func TestNormalizeColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		color   string
		want    string
		wantErr bool
	}{
		{name: "lowercase hex", color: "f54927", want: "f54927"},
		{name: "uppercase hex is normalized", color: "F54927", want: "f54927"},
		{name: "mixed case is normalized", color: "F54aB7", want: "f54ab7"},
		{name: "with hash prefix", color: "#f54927", wantErr: true},
		{name: "too short", color: "f5492", wantErr: true},
		{name: "too long", color: "f549277", wantErr: true},
		{name: "not hexadecimal", color: "zzzzzz", wantErr: true},
		{name: "empty", color: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeColor(tt.color)

			if tt.wantErr {
				assertValidationError(t, err, "color")
				return
			}
			assertNoError(t, err)
			assertEqual(t, tt.want, got)
		})
	}
}
