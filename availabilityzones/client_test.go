package availabilityzones

import (
	"context"
	"net/http"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func newTestCoreClient() *client.CoreClient {
	httpClient := &http.Client{}
	return client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl("http://test-api.com")),
		client.WithHTTPClient(httpClient))
}

func TestNewClient(t *testing.T) {
	core := newTestCoreClient()
	azClient := New(core)

	if azClient == nil {
		t.Error("expected azClient to not be nil")
		return
	}
	if azClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, azClient.CoreClient)
	}
}

func TestClient_newRequest(t *testing.T) {
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
			path:    "/v0/availability-zones",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/v0/availability-zones",
			body:    map[string]string{"name": "test-az"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/v0/availability-zones",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newTestCoreClient()
			azClient := New(core)

			req, err := azClient.newRequest(context.Background(), tt.method, tt.path, tt.body)

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

func TestNewClient_WithOptions(t *testing.T) {
	core := newTestCoreClient()
	customEndpoint := client.MgcUrl("http://custom-endpoint")

	azClient := New(core, WithGlobalBasePath(customEndpoint))

	if azClient == nil {
		t.Fatal("expected azClient to not be nil")
	}

	if azClient.GetConfig().BaseURL != customEndpoint {
		t.Errorf("expected BaseURL to be %s, got %s", customEndpoint, azClient.GetConfig().BaseURL)
	}
}

func TestClient_newRequest_Headers(t *testing.T) {
	core := newTestCoreClient()
	azClient := New(core)

	ctx := context.Background()
	req, err := azClient.newRequest(ctx, http.MethodGet, "/v0/availability-zones", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type header to be application/json, got %s", contentType)
	}

	if userAgent := req.Header.Get("User-Agent"); userAgent == "" {
		t.Error("expected User-Agent header to be set")
	}
}

func TestNewClient_WithNilCore(t *testing.T) {
	azClient := New(nil)
	if azClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestClient_AvailabilityZones(t *testing.T) {
	core := newTestCoreClient()
	azClient := New(core)

	t.Run("AvailabilityZones", func(t *testing.T) {
		svc := azClient.AvailabilityZones()
		if svc == nil {
			t.Error("expected AvailabilityZones service to not be nil")
		}
		if _, ok := svc.(*service); !ok {
			t.Error("expected AvailabilityZones service to be of type *service")
		}
	})
}

func TestClient_DefaultBasePath(t *testing.T) {
	core := newTestCoreClient()
	azClient := New(core)

	req, err := azClient.newRequest(context.Background(), http.MethodGet, "/v0/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := DefaultBasePath + "/v0/test"
	if req.URL.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, req.URL.Path)
	}
}
