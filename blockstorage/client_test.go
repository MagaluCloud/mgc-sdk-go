package blockstorage

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

func TestNewBlockStorageClient(t *testing.T) {
	core := newTestCoreClient()
	bsClient := New(core)

	if bsClient == nil {
		t.Error("expected bsClient to not be nil")
		return
	}
	if bsClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, bsClient.CoreClient)
	}
}

func TestBlockStorageClient_newRequest(t *testing.T) {
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
			path:    "/v1/volumes",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/v1/volumes",
			body:    map[string]string{"name": "test-vol"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/v1/volumes",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newTestCoreClient()
			bsClient := New(core)

			req, err := bsClient.newRequest(context.Background(), tt.method, tt.path, tt.body)

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

func TestNewBlockStorageClient_WithOptions(t *testing.T) {
	core := newTestCoreClient()

	var calledOpt bool
	testOpt := func(c *BlockStorageClient) {
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

func TestBlockStorageClient_newRequest_Headers(t *testing.T) {
	core := newTestCoreClient()
	bsClient := New(core)

	ctx := context.Background()
	req, err := bsClient.newRequest(ctx, http.MethodGet, "/v1/volumes", nil)

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

func TestNewBlockStorageClient_WithNilCore(t *testing.T) {
	bsClient := New(nil)
	if bsClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestBlockStorageClient_Services(t *testing.T) {
	core := newTestCoreClient()
	bsClient := New(core)

	t.Run("Volumes", func(t *testing.T) {
		svc := bsClient.Volumes()
		if svc == nil {
			t.Error("expected VolumeService to not be nil")
		}
		if _, ok := svc.(*volumeService); !ok {
			t.Error("expected VolumeService to be of type *volumeService")
		}
	})

	t.Run("VolumeTypes", func(t *testing.T) {
		svc := bsClient.VolumeTypes()
		if svc == nil {
			t.Error("expected VolumeTypeService to not be nil")
		}
		if _, ok := svc.(*volumeTypeService); !ok {
			t.Error("expected VolumeTypeService to be of type *volumeTypeService")
		}
	})

	t.Run("Snapshots", func(t *testing.T) {
		svc := bsClient.Snapshots()
		if svc == nil {
			t.Error("expected SnapshotService to not be nil")
		}
		if _, ok := svc.(*snapshotService); !ok {
			t.Error("expected SnapshotService to be of type *snapshotService")
		}
	})
}

func TestBlockStorageClient_DefaultBasePath(t *testing.T) {
	core := newTestCoreClient()
	bsClient := New(core)

	req, err := bsClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := DefaultBasePath + "/test"
	if req.URL.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, req.URL.Path)
	}
}
