package containerregistry

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

func TestNewContainerRegistryClient(t *testing.T) {
	core := newTestCoreClient()
	crClient := New(core)

	if crClient == nil {
		t.Error("expected crClient to not be nil")
		return
	}
	if crClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, crClient.CoreClient)
	}
}

func TestContainerRegistryClient_newRequest(t *testing.T) {
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
			path:    "/v0/registries",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/v0/registries",
			body:    map[string]string{"name": "test-vol"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/v0/registries",
			body:    make(chan int),
			wantErr: true,
		},
		{
			name:    "valid GET request with query params",
			method:  http.MethodGet,
			path:    "/v0/registries",
			body:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newTestCoreClient()
			crClient := New(core)

			req, err := crClient.newRequest(context.Background(), tt.method, tt.path, tt.body)

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

func TestNewContainerRegistryClient_WithOptions(t *testing.T) {
	core := newTestCoreClient()

	var calledOpt bool
	testOpt := func(c *ContainerRegistryClient) {
		calledOpt = true
	}

	bsClient := New(core, testOpt)

	if !calledOpt {
		t.Error("expected option to be called")
	}

	if bsClient == nil {
		t.Error("expected client to not be nil")
	}
}

func TestContainerRegistryClient_newRequest_Headers(t *testing.T) {
	core := newTestCoreClient()
	crClient := New(core)

	ctx := context.Background()
	req, err := crClient.newRequest(ctx, http.MethodGet, "/v0/registries", nil)

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

func TestNewContainerRegistryClient_WithNilCore(t *testing.T) {
	crClient := New(nil)
	if crClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestContainerRegistryClient_Services(t *testing.T) {
	core := newTestCoreClient()
	crClient := New(core)

	t.Run("Registries", func(t *testing.T) {
		svc := crClient.Registries()
		if svc == nil {
			t.Error("expected RegistriesService to not be nil")
		}
		if _, ok := svc.(*registriesService); !ok {
			t.Error("expected RegistriesService to be of type *registriesService")
		}
	})
}

func TestContainerRegistryClient_DefaultBasePath(t *testing.T) {
	core := newTestCoreClient()
	crClient := New(core)

	req, err := crClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := DefaultBasePath + "/test"
	if req.URL.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, req.URL.Path)
	}
}

func TestNew_WithOptions(t *testing.T) {
	core := newTestCoreClient()

	var calledOpt bool
	testOpt := func(c *ContainerRegistryClient) {
		calledOpt = true
	}

	crClient := New(core, testOpt)

	if !calledOpt {
		t.Error("expected option to be called")
	}

	if crClient == nil {
		t.Error("expected client to not be nil")
	}
}
