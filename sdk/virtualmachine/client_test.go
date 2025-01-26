package virtualmachine

import (
	"context"
	"net/http"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func newTestCoreClient() *client.CoreClient {
	httpClient := &http.Client{}
	return client.New("test-api", client.WithBaseURL(client.MgcUrl("http://test-api.com")), client.WithHTTPClient(httpClient))
}

func TestNew(t *testing.T) {
	core := newTestCoreClient()
	vmClient := New(core)

	if vmClient == nil {
		t.Error("expected vmClient to not be nil")
		return
	}
	if vmClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, vmClient.CoreClient)
	}
}

func TestVirtualMachineClient_newRequest(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    interface{}
		wantErr bool
	}{
		{
			name:    "valid request",
			method:  http.MethodGet,
			path:    "/vms",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid request with body",
			method:  http.MethodPost,
			path:    "/vms",
			body:    map[string]string{"name": "test-vm"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/vms",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newTestCoreClient()
			vmClient := New(core)
			
			req, err := vmClient.newRequest(context.Background(), tt.method, tt.path, tt.body)
			
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
			
			if req.Method != tt.method {
				t.Errorf("expected method %s, got %s", tt.method, req.Method)
			}
		})
	}
}

func TestNew_WithOptions(t *testing.T) {
	core := newTestCoreClient()
	
	var calledOpt bool
	testOpt := func(c *VirtualMachineClient) {
		calledOpt = true
	}
	
	vmClient := New(core, testOpt)
	
	if !calledOpt {
		t.Error("expected option to be called")
	}
	
	if vmClient == nil {
		t.Error("expected client to not be nil")
	}
}

func TestVirtualMachineClient_newRequest_Headers(t *testing.T) {
	core := newTestCoreClient()
	vmClient := New(core)
	
	ctx := context.Background()
	req, err := vmClient.newRequest(ctx, http.MethodGet, "/vms", nil)
	
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

func TestVirtualMachineClient_NewWithNilCore(t *testing.T) {
	vmClient := New(nil)
	if vmClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestVirtualMachineClient_Instances(t *testing.T) {
	core := newTestCoreClient()
	vmClient := New(core)
	
	instanceSvc := vmClient.Instances()
	
	if instanceSvc == nil {
		t.Error("expected instanceSvc to not be nil")
	}
	
	_, ok := instanceSvc.(*instanceService)
	if !ok {
		t.Error("expected instanceSvc to be of type *instanceService")
	}
}
