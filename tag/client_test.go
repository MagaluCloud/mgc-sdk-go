package tag

import (
	"context"
	"net/http"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func newTestCoreClient() *client.CoreClient {
	httpClient := &http.Client{}
	return client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl("http://test-api.com")),
		client.WithHTTPClient(httpClient),
	)
}

func TestNew(t *testing.T) {
	t.Run("returns nil when core is nil", func(t *testing.T) {
		c := New(nil)
		if c != nil {
			t.Error("expected nil client when core is nil")
		}
	})

	t.Run("returns valid client with core", func(t *testing.T) {
		core := newTestCoreClient()
		c := New(core)
		if c == nil {
			t.Fatal("expected non-nil client")
		}
		if c.CoreClient != core {
			t.Error("expected CoreClient to match the provided core")
		}
	})

	t.Run("sets global base URL by default", func(t *testing.T) {
		core := newTestCoreClient()
		c := New(core)
		if c.GetConfig().BaseURL != client.Global {
			t.Errorf("expected BaseURL %s, got %s", client.Global, c.GetConfig().BaseURL)
		}
	})

	t.Run("WithBasePath overrides base URL", func(t *testing.T) {
		core := newTestCoreClient()
		custom := client.MgcUrl("http://custom-endpoint")
		c := New(core, WithBasePath(custom))
		if c.GetConfig().BaseURL != custom {
			t.Errorf("expected BaseURL %s, got %s", custom, c.GetConfig().BaseURL)
		}
	})
}

func TestTagClient_newRequest(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    interface{}
		wantErr bool
	}{
		{
			name:    "valid GET request",
			method:  http.MethodGet,
			path:    "/v0/tags",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/v0/tags",
			body:    map[string]string{"name": "test-tag"},
			wantErr: false,
		},
		{
			name:    "valid PATCH request",
			method:  http.MethodPatch,
			path:    "/v0/tags/my-tag",
			body:    map[string]string{"description": "updated"},
			wantErr: false,
		},
		{
			name:    "valid DELETE request",
			method:  http.MethodDelete,
			path:    "/v0/tags/my-tag",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "invalid body returns error",
			method:  http.MethodPost,
			path:    "/v0/tags",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newTestCoreClient()
			c := New(core)

			req, err := c.newRequest(context.Background(), tt.method, tt.path, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req == nil {
				t.Fatal("expected non-nil request")
			}

			expectedPath := DefaultBasePath + tt.path
			if req.URL.Path != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.URL.Path)
			}
			if req.Method != tt.method {
				t.Errorf("expected method %s, got %s", tt.method, req.Method)
			}
		})
	}
}

func TestTagClient_newRequest_Headers(t *testing.T) {
	core := newTestCoreClient()
	c := New(core)

	req, err := c.newRequest(context.Background(), http.MethodGet, "/v0/tags", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ct := req.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	if ua := req.Header.Get("User-Agent"); ua == "" {
		t.Error("expected User-Agent header to be set")
	}
}

func TestTagClient_Tags(t *testing.T) {
	core := newTestCoreClient()
	c := New(core)

	svc := c.Tags()
	if svc == nil {
		t.Fatal("expected non-nil TagService")
	}
	if _, ok := svc.(*tagService); !ok {
		t.Error("expected TagService to be of type *tagService")
	}
}
