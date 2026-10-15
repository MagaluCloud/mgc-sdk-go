package lbaas

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

func TestNewLbaasClient(t *testing.T) {
	t.Parallel()
	core := newTestCoreClient()
	lbaasClient := New(core)

	if lbaasClient == nil {
		t.Error("expected lbaasClient to not be nil")
		return
	}
	if lbaasClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, lbaasClient.CoreClient)
	}
}

func TestNewLbaasClient_WithNilCore(t *testing.T) {
	t.Parallel()
	lbaasClient := New(nil)
	if lbaasClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestLbaasClient_newRequest(t *testing.T) {
	t.Parallel()
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
			path:    "/v0beta1/network-load-balancers",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/v0beta1/network-load-balancers",
			body:    map[string]string{"name": "test-lb"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/v0beta1/network-load-balancers",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			core := newTestCoreClient()
			lbaasClient := New(core)

			req, err := lbaasClient.newRequest(context.Background(), tt.method, tt.path, tt.body)

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

func TestLbaasClient_newRequest_Headers(t *testing.T) {
	t.Parallel()
	core := newTestCoreClient()
	lbaasClient := New(core)

	ctx := context.Background()
	req, err := lbaasClient.newRequest(ctx, http.MethodGet, "/v0beta1/network-load-balancers", nil)

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

func TestLbaasClient_Services(t *testing.T) {
	t.Parallel()
	core := newTestCoreClient()
	lbaasClient := New(core)

	t.Run("NetworkACLs", func(t *testing.T) {
		t.Parallel()
		svc := lbaasClient.NetworkACLs()
		if svc == nil {
			t.Error("expected NetworkACLService to not be nil")
		}
		if _, ok := svc.(*networkACLService); !ok {
			t.Error("expected NetworkACLService to be of type *networkACLService")
		}
	})

	t.Run("NetworkBackends", func(t *testing.T) {
		t.Parallel()
		svc := lbaasClient.NetworkBackends()
		if svc == nil {
			t.Error("expected NetworkBackendService to not be nil")
		}
		if _, ok := svc.(*networkBackendService); !ok {
			t.Error("expected NetworkBackendService to be of type *networkBackendService")
		}
	})

	t.Run("NetworkCertificates", func(t *testing.T) {
		t.Parallel()
		svc := lbaasClient.NetworkCertificates()
		if svc == nil {
			t.Error("expected NetworkCertificateService to not be nil")
		}
		if _, ok := svc.(*networkCertificateService); !ok {
			t.Error("expected NetworkCertificateService to be of type *networkCertificateService")
		}
	})

	t.Run("NetworkHealthChecks", func(t *testing.T) {
		t.Parallel()
		svc := lbaasClient.NetworkHealthChecks()
		if svc == nil {
			t.Error("expected NetworkHealthCheckService to not be nil")
		}
		if _, ok := svc.(*networkHealthCheckService); !ok {
			t.Error("expected NetworkHealthCheckService to be of type *networkHealthCheckService")
		}
	})

	t.Run("NetworkListeners", func(t *testing.T) {
		t.Parallel()
		svc := lbaasClient.NetworkListeners()
		if svc == nil {
			t.Error("expected NetworkListenerService to not be nil")
		}
		if _, ok := svc.(*networkListenerService); !ok {
			t.Error("expected NetworkListenerService to be of type *networkListenerService")
		}
	})

	t.Run("NetworkLoadBalancers", func(t *testing.T) {
		t.Parallel()
		svc := lbaasClient.NetworkLoadBalancers()
		if svc == nil {
			t.Error("expected NetworkLoadBalancerService to not be nil")
		}
		if _, ok := svc.(*networkLoadBalancerService); !ok {
			t.Error("expected NetworkLoadBalancerService to be of type *networkLoadBalancerService")
		}
	})
}

func TestLbaasClient_DefaultBasePath(t *testing.T) {
	t.Parallel()
	expected := "/load-balancer"
	if DefaultBasePath != expected {
		t.Errorf("expected DefaultBasePath to be %s, got %s", expected, DefaultBasePath)
	}
}
