package network

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

func TestNewNetworkClient(t *testing.T) {
	t.Parallel()
	core := newTestCoreClient()
	networkClient := New(core)

	if networkClient == nil {
		t.Error("expected networkClient to not be nil")
		return
	}
	if networkClient.CoreClient != core {
		t.Errorf("expected CoreClient to be %v, got %v", core, networkClient.CoreClient)
	}
}

func TestNetworkClient_newRequest(t *testing.T) {
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
			path:    "/v1/vpcs",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "valid POST request with body",
			method:  http.MethodPost,
			path:    "/v1/vpcs",
			body:    map[string]string{"name": "test-vpc"},
			wantErr: false,
		},
		{
			name:    "invalid body",
			method:  http.MethodPost,
			path:    "/v1/vpcs",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			core := newTestCoreClient()
			networkClient := New(core)

			req, err := networkClient.newRequest(context.Background(), tt.method, tt.path, tt.body)

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

func TestNetworkClient_newRequest_Headers(t *testing.T) {
	t.Parallel()
	core := newTestCoreClient()
	networkClient := New(core)

	ctx := context.Background()
	req, err := networkClient.newRequest(ctx, http.MethodGet, "/v1/vpcs", nil)

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

func TestNewNetworkClient_WithNilCore(t *testing.T) {
	t.Parallel()
	networkClient := New(nil)
	if networkClient != nil {
		t.Error("expected nil client when core is nil")
	}
}

func TestNetworkClient_Services(t *testing.T) {
	t.Parallel()
	core := newTestCoreClient()
	networkClient := New(core)

	t.Run("VPCs", func(t *testing.T) {
		t.Parallel()
		svc := networkClient.VPCs()
		if svc == nil {
			t.Error("expected VPCService to not be nil")
		}
		if _, ok := svc.(*vpcService); !ok {
			t.Error("expected VPCService to be of type *vpcService")
		}
	})

	t.Run("Subnets", func(t *testing.T) {
		t.Parallel()
		svc := networkClient.Subnets()
		if svc == nil {
			t.Error("expected SubnetService to not be nil")
		}
		if _, ok := svc.(*subnetService); !ok {
			t.Error("expected SubnetService to be of type *subnetService")
		}
	})

	t.Run("Ports", func(t *testing.T) {
		t.Parallel()
		svc := networkClient.Ports()
		if svc == nil {
			t.Error("expected PortService to not be nil")
		}
		if _, ok := svc.(*portService); !ok {
			t.Error("expected PortService to be of type *portService")
		}
	})

	t.Run("SecurityGroups", func(t *testing.T) {
		t.Parallel()
		svc := networkClient.SecurityGroups()
		if svc == nil {
			t.Error("expected SecurityGroupService to not be nil")
		}
		if _, ok := svc.(*securityGroupService); !ok {
			t.Error("expected SecurityGroupService to be of type *securityGroupService")
		}
	})

	t.Run("Rules", func(t *testing.T) {
		t.Parallel()
		svc := networkClient.Rules()
		if svc == nil {
			t.Error("expected RuleService to not be nil")
		}
		if _, ok := svc.(*ruleService); !ok {
			t.Error("expected RuleService to be of type *ruleService")
		}
	})

	t.Run("PublicIPs", func(t *testing.T) {
		t.Parallel()
		svc := networkClient.PublicIPs()
		if svc == nil {
			t.Error("expected PublicIPService to not be nil")
		}
		if _, ok := svc.(*publicIPService); !ok {
			t.Error("expected PublicIPService to be of type *publicIPService")
		}
	})

	t.Run("SubnetPools", func(t *testing.T) {
		t.Parallel()
		svc := networkClient.SubnetPools()
		if svc == nil {
			t.Error("expected SubnetPoolService to not be nil")
		}
		if _, ok := svc.(*subnetPoolService); !ok {
			t.Error("expected SubnetPoolService to be of type *subnetPoolService")
		}
	})
}

func TestNetworkClient_DefaultBasePath(t *testing.T) {
	t.Parallel()
	core := newTestCoreClient()
	networkClient := New(core)

	req, err := networkClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := DefaultBasePath + "/test"
	if req.URL.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, req.URL.Path)
	}
}
