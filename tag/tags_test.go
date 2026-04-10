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

// --- helpers ---

func newTagTestServer() (*httptest.Server, *TagClient) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v0/tags" && r.Method == http.MethodGet:
			handleListTags(w, r)
		case r.URL.Path == "/v0/tags" && r.Method == http.MethodPost:
			handleCreateTag(w, r)
		case strings.HasPrefix(r.URL.Path, "/v0/tags/") && r.Method == http.MethodGet:
			handleGetTag(w, r)
		case strings.HasPrefix(r.URL.Path, "/v0/tags/") && r.Method == http.MethodPatch:
			handleUpdateTag(w, r)
		case strings.HasPrefix(r.URL.Path, "/v0/tags/") && r.Method == http.MethodDelete:
			handleDeleteTag(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))

	cfg := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithTimeout(20*time.Second),
	)
	c := New(cfg, WithBasePath(client.MgcUrl(ts.URL)))
	return ts, c
}

func errorServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error": "server error"}`)
	}))
}

func newServiceWithErrorServer() (TagService, func()) {
	ts := errorServer()
	cfg := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
	)
	svc := New(cfg, WithBasePath(client.MgcUrl(ts.URL))).Tags()
	return svc, ts.Close
}

// --- request handlers ---

func handleListTags(w http.ResponseWriter, r *http.Request) {
	response := ListTagsResponse{
		Results: []Tag{
			{
				Name:        "environment",
				Description: "Environment tag",
				Color:       "ff0000ff",
				Kinds:       []string{"virtual-machine"},
				Values: []TagValue{
					{Name: "production", Description: "Prod env", CreatedAt: "2024-01-01T00:00:00Z"},
				},
				CreatedAt: "2024-01-01T00:00:00Z",
			},
			{
				Name:        "team",
				Description: "Team tag",
				Color:       "00ff00ff",
				Kinds:       []string{"virtual-machine", "volume"},
				CreatedAt:   "2024-01-02T00:00:00Z",
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleCreateTag(w http.ResponseWriter, r *http.Request) {
	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response := Tag{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Kinds:       req.Kinds,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}
	for _, v := range req.Values {
		response.Values = append(response.Values, TagValue{
			Name:        v.Name,
			Description: v.Description,
			CreatedAt:   "2024-01-01T00:00:00Z",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleUpdateTag(w http.ResponseWriter, r *http.Request) {
	tagName := strings.TrimPrefix(r.URL.Path, "/v0/tags/")

	var req UpdateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	desc := "updated description"
	if req.Description != nil {
		desc = *req.Description
	}
	color := "aabbccdd"
	if req.Color != nil {
		color = *req.Color
	}

	response := Tag{
		Name:        tagName,
		Description: desc,
		Color:       color,
		Kinds:       req.Kinds,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetTag(w http.ResponseWriter, r *http.Request) {
	tagName := strings.TrimPrefix(r.URL.Path, "/v0/tags/")

	response := Tag{
		Name:        tagName,
		Description: "A tag",
		Color:       "ff0000ff",
		Kinds:       []string{"virtual-machine"},
		CreatedAt:   "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDeleteTag(w http.ResponseWriter, r *http.Request) {
	tagName := strings.TrimPrefix(r.URL.Path, "/v0/tags/")

	response := Tag{
		Name:      tagName,
		Color:     "ff0000ff",
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// --- tests ---

func TestTagService_List(t *testing.T) {
	ts, c := newTagTestServer()
	defer ts.Close()
	svc := c.Tags()

	t.Run("lists all tags successfully", func(t *testing.T) {
		tags, err := svc.List(context.Background(), ListOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tags) != 2 {
			t.Fatalf("expected 2 tags, got %d", len(tags))
		}
		if tags[0].Name != "environment" {
			t.Errorf("expected first tag name 'environment', got %s", tags[0].Name)
		}
		if tags[1].Name != "team" {
			t.Errorf("expected second tag name 'team', got %s", tags[1].Name)
		}
	})

	t.Run("lists tags with name filter", func(t *testing.T) {
		tags, err := svc.List(context.Background(), ListOptions{
			Name: helpers.StrPtr("env"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tags) == 0 {
			t.Error("expected at least one tag in response")
		}
	})

	t.Run("returns values nested in tags", func(t *testing.T) {
		tags, err := svc.List(context.Background(), ListOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tags[0].Values) != 1 {
			t.Fatalf("expected 1 value in first tag, got %d", len(tags[0].Values))
		}
		if tags[0].Values[0].Name != "production" {
			t.Errorf("expected value name 'production', got %s", tags[0].Values[0].Name)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newServiceWithErrorServer()
		defer closeServer()
		_, err := svc.List(context.Background(), ListOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "server error") && !strings.Contains(err.Error(), "500") {
			t.Errorf("expected server error in message, got: %v", err)
		}
	})
}

func TestTagService_Create(t *testing.T) {
	ts, c := newTagTestServer()
	defer ts.Close()
	svc := c.Tags()

	t.Run("creates tag successfully", func(t *testing.T) {
		tag, err := svc.Create(context.Background(), CreateTagRequest{
			Name:        "my-tag",
			Description: "A test tag",
			Color:       "ff0000ff",
			Kinds:       []string{"virtual-machine"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tag.Name != "my-tag" {
			t.Errorf("expected name 'my-tag', got %s", tag.Name)
		}
		if tag.Color != "ff0000ff" {
			t.Errorf("expected color 'ff0000ff', got %s", tag.Color)
		}
	})

	t.Run("creates tag with values", func(t *testing.T) {
		tag, err := svc.Create(context.Background(), CreateTagRequest{
			Name:  "env-tag",
			Color: "aabbccdd",
			Values: []CreateTagValueRequest{
				{Name: "production", Description: "Production environment"},
				{Name: "staging", Description: "Staging environment"},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tag.Values) != 2 {
			t.Fatalf("expected 2 values, got %d", len(tag.Values))
		}
		if tag.Values[0].Name != "production" {
			t.Errorf("expected value 'production', got %s", tag.Values[0].Name)
		}
	})

	t.Run("converts uppercase color to lowercase", func(t *testing.T) {
		tag, err := svc.Create(context.Background(), CreateTagRequest{
			Name:  "upper-color-tag",
			Color: "FF0000FF",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tag.Color != "ff0000ff" {
			t.Errorf("expected color normalized to 'ff0000ff', got %s", tag.Color)
		}
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), CreateTagRequest{
			Name:  "",
			Color: "ff0000ff",
		})
		if err == nil {
			t.Fatal("expected error for empty name, got nil")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "name" {
			t.Errorf("expected field 'name', got %s", valErr.Field)
		}
	})

	t.Run("returns error when value name is empty", func(t *testing.T) {
		_, err := svc.Create(context.Background(), CreateTagRequest{
			Name:  "valid-name",
			Color: "ff0000ff",
			Values: []CreateTagValueRequest{
				{Name: ""},
			},
		})
		if err == nil {
			t.Fatal("expected error for empty value name, got nil")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "values[0].name" {
			t.Errorf("expected field 'values[0].name', got %s", valErr.Field)
		}
	})

	t.Run("returns error when color format is invalid", func(t *testing.T) {
		_, err := svc.Create(context.Background(), CreateTagRequest{
			Name:  "valid-name",
			Color: "gg112233",
		})
		if err == nil {
			t.Fatal("expected error for invalid color, got nil")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "color" {
			t.Errorf("expected field 'color', got %s", valErr.Field)
		}
	})

	t.Run("returns error when color has wrong length", func(t *testing.T) {
		_, err := svc.Create(context.Background(), CreateTagRequest{
			Name:  "valid-name",
			Color: "ff0000",
		})
		if err == nil {
			t.Fatal("expected error for short color, got nil")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "color" {
			t.Errorf("expected field 'color', got %s", valErr.Field)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newServiceWithErrorServer()
		defer closeServer()
		_, err := svc.Create(context.Background(), CreateTagRequest{
			Name:  "tag",
			Color: "ff0000ff",
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagService_Update(t *testing.T) {
	ts, c := newTagTestServer()
	defer ts.Close()
	svc := c.Tags()

	t.Run("updates tag successfully", func(t *testing.T) {
		desc := "new description"
		tag, err := svc.Update(context.Background(), "environment", UpdateTagRequest{
			Description: &desc,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tag.Name != "environment" {
			t.Errorf("expected name 'environment', got %s", tag.Name)
		}
		if tag.Description != "new description" {
			t.Errorf("expected description 'new description', got %s", tag.Description)
		}
	})

	t.Run("updates color and normalizes uppercase to lowercase", func(t *testing.T) {
		color := "AABBCCDD"
		tag, err := svc.Update(context.Background(), "my-tag", UpdateTagRequest{
			Color: &color,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tag.Color != "aabbccdd" {
			t.Errorf("expected color 'aabbccdd', got %s", tag.Color)
		}
	})

	t.Run("returns error when tag name is empty", func(t *testing.T) {
		_, err := svc.Update(context.Background(), "", UpdateTagRequest{})
		if err == nil {
			t.Fatal("expected error for empty tag name, got nil")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "tag_name" {
			t.Errorf("expected field 'tag_name', got %s", valErr.Field)
		}
	})

	t.Run("returns error when color format is invalid", func(t *testing.T) {
		color := "zzzzzzzz"
		_, err := svc.Update(context.Background(), "my-tag", UpdateTagRequest{
			Color: &color,
		})
		if err == nil {
			t.Fatal("expected error for invalid color, got nil")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "color" {
			t.Errorf("expected field 'color', got %s", valErr.Field)
		}
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		svc, closeServer := newServiceWithErrorServer()
		defer closeServer()
		_, err := svc.Update(context.Background(), "tag", UpdateTagRequest{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagService_Delete(t *testing.T) {
	ts, c := newTagTestServer()
	defer ts.Close()
	svc := c.Tags()

	t.Run("deletes tag successfully", func(t *testing.T) {
		tag, err := svc.Delete(context.Background(), "environment")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tag.Name != "environment" {
			t.Errorf("expected name 'environment', got %s", tag.Name)
		}
	})

	t.Run("returns error when tag name is empty", func(t *testing.T) {
		_, err := svc.Delete(context.Background(), "")
		if err == nil {
			t.Fatal("expected error for empty tag name, got nil")
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
		svc, closeServer := newServiceWithErrorServer()
		defer closeServer()
		_, err := svc.Delete(context.Background(), "tag")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTagService_List_Filters(t *testing.T) {
	ts, c := newTagTestServer()
	defer ts.Close()
	svc := c.Tags()

	t.Run("lists tags with color filter (normalizes uppercase)", func(t *testing.T) {
		tags, err := svc.List(context.Background(), ListOptions{
			Color: helpers.StrPtr("FF0000FF"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tags) == 0 {
			t.Error("expected at least one tag")
		}
	})

	t.Run("returns error when color filter has invalid format", func(t *testing.T) {
		_, err := svc.List(context.Background(), ListOptions{
			Color: helpers.StrPtr("zzzzzzzz"),
		})
		if err == nil {
			t.Fatal("expected error for invalid color filter, got nil")
		}
		var valErr *client.ValidationError
		if !isValidationError(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
		if valErr.Field != "color" {
			t.Errorf("expected field 'color', got %s", valErr.Field)
		}
	})

	t.Run("lists tags with pagination options", func(t *testing.T) {
		tags, err := svc.List(context.Background(), ListOptions{
			Limit:  helpers.IntPtr(10),
			Offset: helpers.IntPtr(0),
			Sort:   helpers.StrPtr("name"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tags) == 0 {
			t.Error("expected at least one tag")
		}
	})
}

func TestTagService_Get(t *testing.T) {
	ts, c := newTagTestServer()
	defer ts.Close()
	svc := c.Tags()

	t.Run("gets a tag successfully", func(t *testing.T) {
		tag, err := svc.Get(context.Background(), "environment")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tag.Name != "environment" {
			t.Errorf("expected name 'environment', got %s", tag.Name)
		}
		if tag.Color != "ff0000ff" {
			t.Errorf("expected color 'ff0000ff', got %s", tag.Color)
		}
	})

	t.Run("returns error when tag name is empty", func(t *testing.T) {
		_, err := svc.Get(context.Background(), "")
		if err == nil {
			t.Fatal("expected error for empty tag name, got nil")
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
		svc, closeServer := newServiceWithErrorServer()
		defer closeServer()
		_, err := svc.Get(context.Background(), "tag")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestNormalizeColor(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		wantError bool
	}{
		{
			name:     "valid lowercase color",
			input:    "ff0000ff",
			expected: "ff0000ff",
		},
		{
			name:     "valid uppercase color converted to lowercase",
			input:    "FF0000FF",
			expected: "ff0000ff",
		},
		{
			name:     "valid mixed case color converted to lowercase",
			input:    "aAbBcCdD",
			expected: "aabbccdd",
		},
		{
			name:      "invalid color - too short",
			input:     "ff0000f",
			wantError: true,
		},
		{
			name:      "invalid color - too long",
			input:     "ff0000fff",
			wantError: true,
		},
		{
			name:      "invalid color - non-hex characters",
			input:     "gg0000ff",
			wantError: true,
		},
		{
			name:      "invalid color - empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "invalid color - with hash prefix",
			input:     "#ff0000f",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeColor(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// isValidationError checks if err is a *client.ValidationError and assigns it to target.
func isValidationError(err error, target **client.ValidationError) bool {
	if v, ok := err.(*client.ValidationError); ok {
		*target = v
		return true
	}
	return false
}
