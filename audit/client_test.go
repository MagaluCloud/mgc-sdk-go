package audit

import (
	"context"
	"net/http"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func newTestCoreClient() *client.CoreClient {
	httpClient := &http.Client{}
	return client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl("http://test-api.com")),
		client.WithHTTPClient(httpClient))
}

func TestNewAuditClient(t *testing.T) {
	core := newTestCoreClient()
	auditClient := New(core)

	if auditClient == nil {
		t.Error("expected auditClient to not be nil")
		return
	}
	if auditClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, auditClient.CoreClient)
	}
}

func TestAuditClient_newRequest(t *testing.T) {
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
			path:    "/events",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/events",
			body:    map[string]string{"event": "test"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/events",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newTestCoreClient()
			auditClient := New(core)

			req, err := auditClient.newRequest(context.Background(), tt.method, tt.path, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if req == nil {
				t.Error("expected request to not be nil")
				return
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

func TestAuditClient_newRequest_Headers(t *testing.T) {
	core := newTestCoreClient()
	auditClient := New(core)

	ctx := context.Background()
	req, err := auditClient.newRequest(ctx, http.MethodGet, "/test", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check content type header
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type header to be application/json, got %s", contentType)
	}

	// Check user agent
	if userAgent := req.Header.Get("User-Agent"); userAgent == "" {
		t.Error("expected User-Agent header to be set")
	}
}

func TestNewAuditClient_WithNilCore(t *testing.T) {
	auditClient := New(nil)
	if auditClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestAuditClient_Services(t *testing.T) {
	core := newTestCoreClient()
	auditClient := New(core)

	t.Run("Events", func(t *testing.T) {
		svc := auditClient.Events()
		if svc == nil {
			t.Error("expected EventService to not be nil")
		}
		if _, ok := svc.(*eventService); !ok {
			t.Error("expected EventService to be of type *eventService")
		}
	})

	t.Run("EventTypes", func(t *testing.T) {
		svc := auditClient.EventTypes()
		if svc == nil {
			t.Error("expected EventTypeService to not be nil")
		}
		if _, ok := svc.(*eventTypeService); !ok {
			t.Error("expected EventTypeService to be of type *eventTypeService")
		}
	})
}

func TestAuditClient_DefaultBasePath(t *testing.T) {
	core := newTestCoreClient()
	auditClient := New(core)

	req, err := auditClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := DefaultBasePath + "/test"
	if req.URL.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, req.URL.Path)
	}
}
