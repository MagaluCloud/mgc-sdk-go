package lbaas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

const scopedProjectID = "550e8400-e29b-41d4-a716-446655440000"

func TestServices_ProjectHeader(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		call func(c *LbaasClient, opts client.CallOptions) error
	}{
		{"NetworkACLs.Create", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkACLs().Create(ctx, opts, "lb-1", CreateNetworkACLRequest{})
			return err
		}},
		{"NetworkACLs.Delete", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkACLs().Delete(ctx, opts, "lb-1", "acl-1")
		}},
		{"NetworkACLs.Replace", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkACLs().Replace(ctx, opts, "lb-1", UpdateNetworkACLRequest{})
		}},
		{"NetworkBackends.Create", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkBackends().Create(ctx, opts, "lb-1", CreateBackendRequest{})
			return err
		}},
		{"NetworkBackends.Delete", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkBackends().Delete(ctx, opts, "lb-1", "be-1")
		}},
		{"NetworkBackends.Get", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkBackends().Get(ctx, opts, "lb-1", "be-1")
			return err
		}},
		{"NetworkBackends.List", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkBackends().List(ctx, opts, "lb-1", ListNetworkLoadBalancerRequest{})
			return err
		}},
		{"NetworkBackends.ListAll", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkBackends().ListAll(ctx, opts, "lb-1")
			return err
		}},
		{"NetworkBackends.Update", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkBackends().Update(ctx, opts, "lb-1", "be-1", UpdateNetworkBackendRequest{})
			return err
		}},
		{"NetworkBackendTargets.Create", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkBackendTargets().Create(ctx, opts, "lb-1", "be-1", CreateNetworkBackendTargetRequest{})
			return err
		}},
		{"NetworkBackendTargets.Replace", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkBackendTargets().Replace(ctx, opts, "lb-1", "be-1", CreateNetworkBackendTargetRequest{})
			return err
		}},
		{"NetworkBackendTargets.Delete", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkBackendTargets().Delete(ctx, opts, "lb-1", "be-1", "target-1")
		}},
		{"NetworkCertificates.Create", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkCertificates().Create(ctx, opts, "lb-1", CreateNetworkCertificateRequest{})
			return err
		}},
		{"NetworkCertificates.Delete", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkCertificates().Delete(ctx, opts, "lb-1", "cert-1")
		}},
		{"NetworkCertificates.Get", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkCertificates().Get(ctx, opts, "lb-1", "cert-1")
			return err
		}},
		{"NetworkCertificates.List", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkCertificates().List(ctx, opts, "lb-1", ListNetworkLoadBalancerRequest{})
			return err
		}},
		{"NetworkCertificates.ListAll", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkCertificates().ListAll(ctx, opts, "lb-1")
			return err
		}},
		{"NetworkCertificates.Update", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkCertificates().Update(ctx, opts, "lb-1", "cert-1", UpdateNetworkCertificateRequest{})
		}},
		{"NetworkHealthChecks.Create", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkHealthChecks().Create(ctx, opts, "lb-1", CreateNetworkHealthCheckRequest{})
			return err
		}},
		{"NetworkHealthChecks.Delete", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkHealthChecks().Delete(ctx, opts, "lb-1", "hc-1")
		}},
		{"NetworkHealthChecks.Get", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkHealthChecks().Get(ctx, opts, "lb-1", "hc-1")
			return err
		}},
		{"NetworkHealthChecks.List", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkHealthChecks().List(ctx, opts, "lb-1", ListNetworkLoadBalancerRequest{})
			return err
		}},
		{"NetworkHealthChecks.ListAll", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkHealthChecks().ListAll(ctx, opts, "lb-1")
			return err
		}},
		{"NetworkHealthChecks.Update", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkHealthChecks().Update(ctx, opts, "lb-1", "hc-1", UpdateNetworkHealthCheckRequest{})
		}},
		{"NetworkListeners.Create", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkListeners().Create(ctx, opts, "lb-1", "be-1", CreateNetworkListenerRequest{})
			return err
		}},
		{"NetworkListeners.Delete", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkListeners().Delete(ctx, opts, "lb-1", "listener-1")
		}},
		{"NetworkListeners.Get", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkListeners().Get(ctx, opts, "lb-1", "listener-1")
			return err
		}},
		{"NetworkListeners.List", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkListeners().List(ctx, opts, "lb-1", ListNetworkLoadBalancerRequest{})
			return err
		}},
		{"NetworkListeners.ListAll", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkListeners().ListAll(ctx, opts, "lb-1")
			return err
		}},
		{"NetworkListeners.Update", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkListeners().Update(ctx, opts, "lb-1", "listener-1", UpdateNetworkListenerRequest{})
		}},
		{"NetworkLoadBalancers.Create", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkLoadBalancers().Create(ctx, opts, CreateNetworkLoadBalancerRequest{})
			return err
		}},
		{"NetworkLoadBalancers.Delete", func(c *LbaasClient, opts client.CallOptions) error {
			return c.NetworkLoadBalancers().Delete(ctx, opts, "lb-1", DeleteNetworkLoadBalancerRequest{})
		}},
		{"NetworkLoadBalancers.Get", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkLoadBalancers().Get(ctx, opts, "lb-1")
			return err
		}},
		{"NetworkLoadBalancers.List", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkLoadBalancers().List(ctx, opts, ListNetworkLoadBalancerRequest{})
			return err
		}},
		{"NetworkLoadBalancers.ListAll", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkLoadBalancers().ListAll(ctx, opts)
			return err
		}},
		{"NetworkLoadBalancers.Update", func(c *LbaasClient, opts client.CallOptions) error {
			_, err := c.NetworkLoadBalancers().Update(ctx, opts, "lb-1", UpdateNetworkLoadBalancerRequest{})
			return err
		}},
	}

	scopes := map[string]client.CallOptions{
		"scoped":          {ProjectID: scopedProjectID},
		"default project": {},
	}

	for _, tt := range tests {
		for scope, opts := range scopes {
			t.Run(tt.name+"/"+scope, func(t *testing.T) {
				probe := &projectHeaderProbe{}
				server := httptest.NewServer(probe)
				defer server.Close()

				core := client.NewMgcClient(
					client.WithAPIKey("test-api-key"),
					client.WithBaseURL(client.MgcUrl(server.URL)),
				)

				if err := tt.call(New(core), opts); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				probe.assertEvery(t, opts)
			})
		}
	}
}

// projectHeaderProbe records the project header of every request it serves.
type projectHeaderProbe struct {
	mu       sync.Mutex
	requests []projectHeader
}

type projectHeader struct {
	value   string
	present bool
}

func (p *projectHeaderProbe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values, present := r.Header[client.ProjectIDHeader]
	seen := projectHeader{present: present}
	if present {
		seen.value = values[0]
	}

	p.mu.Lock()
	p.requests = append(p.requests, seen)
	p.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"id":"abc"}`))
}

func (p *projectHeaderProbe) assertEvery(t *testing.T, opts client.CallOptions) {
	t.Helper()

	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.requests) == 0 {
		t.Fatal("no request reached the server")
	}

	want := projectHeader{value: opts.ProjectID, present: opts.ProjectID != ""}
	for i, got := range p.requests {
		if got != want {
			t.Errorf("request %d: %s = %q (present %v), want %q (present %v)",
				i+1, client.ProjectIDHeader, got.value, got.present, want.value, want.present)
		}
	}
}
