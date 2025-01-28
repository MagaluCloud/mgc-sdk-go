package sshkeys

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

func TestNewSSHKeyClient(t *testing.T) {
	core := newTestCoreClient()
	sshClient := New(core)

	if sshClient == nil {
		t.Error("expected sshClient to not be nil")
		return
	}
	if sshClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, sshClient.CoreClient)
	}
}

func TestSSHKeyClient_newRequest(t *testing.T) {
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
			path:    "/v0/ssh-keys",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/v0/ssh-keys",
			body:    map[string]string{"name": "test-key"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/v0/ssh-keys",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newTestCoreClient()
			sshClient := New(core)

			req, err := sshClient.newRequest(context.Background(), tt.method, tt.path, tt.body)

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

func TestNewSSHKeyClient_WithOptions(t *testing.T) {
	core := newTestCoreClient()
	customEndpoint := client.MgcUrl("http://custom-endpoint")

	sshClient := New(core, WithGlobalBasePath(customEndpoint))

	if sshClient == nil {
		t.Fatal("expected sshClient to not be nil")
	}

	if sshClient.GetConfig().BaseURL != customEndpoint {
		t.Errorf("expected BaseURL to be %s, got %s", customEndpoint, sshClient.GetConfig().BaseURL)
	}
}

func TestSSHKeyClient_newRequest_Headers(t *testing.T) {
	core := newTestCoreClient()
	sshClient := New(core)

	ctx := context.Background()
	req, err := sshClient.newRequest(ctx, http.MethodGet, "/v0/ssh-keys", nil)

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

func TestNewSSHKeyClient_WithNilCore(t *testing.T) {
	sshClient := New(nil)
	if sshClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestSSHKeyClient_Services(t *testing.T) {
	core := newTestCoreClient()
	sshClient := New(core)

	t.Run("Keys", func(t *testing.T) {
		svc := sshClient.Keys()
		if svc == nil {
			t.Error("expected KeyService to not be nil")
		}
		if _, ok := svc.(*keyService); !ok {
			t.Error("expected KeyService to be of type *keyService")
		}
	})
}

func TestSSHKeyClient_DefaultBasePath(t *testing.T) {
	core := newTestCoreClient()
	sshClient := New(core)

	req, err := sshClient.newRequest(context.Background(), http.MethodGet, "/v0/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := DefaultBasePath + "/v0/test"
	if req.URL.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, req.URL.Path)
	}
}