package kubernetes

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
		call func(c *KubernetesClient, opts client.CallOptions) error
	}{
		{"Clusters.List", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Clusters().List(ctx, opts, ListOptions{})
			return err
		}},
		{"Clusters.Create", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Clusters().Create(ctx, opts, ClusterRequest{})
			return err
		}},
		{"Clusters.Get", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Clusters().Get(ctx, opts, "cluster-1")
			return err
		}},
		{"Clusters.Delete", func(c *KubernetesClient, opts client.CallOptions) error {
			return c.Clusters().Delete(ctx, opts, "cluster-1")
		}},
		{"Clusters.Update", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Clusters().Update(ctx, opts, "cluster-1", PatchClusterRequest{})
			return err
		}},
		{"Clusters.GetKubeConfig", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Clusters().GetKubeConfig(ctx, opts, "cluster-1")
			return err
		}},
		{"Nodepools.Nodes", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Nodepools().Nodes(ctx, opts, "cluster-1", "pool-1")
			return err
		}},
		{"Nodepools.List", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Nodepools().List(ctx, opts, "cluster-1", ListOptions{})
			return err
		}},
		{"Nodepools.Create", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Nodepools().Create(ctx, opts, "cluster-1", CreateNodePoolRequest{})
			return err
		}},
		{"Nodepools.Get", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Nodepools().Get(ctx, opts, "cluster-1", "pool-1")
			return err
		}},
		{"Nodepools.Update", func(c *KubernetesClient, opts client.CallOptions) error {
			_, err := c.Nodepools().Update(ctx, opts, "cluster-1", "pool-1", PatchNodePoolRequest{})
			return err
		}},
		{"Nodepools.Delete", func(c *KubernetesClient, opts client.CallOptions) error {
			return c.Nodepools().Delete(ctx, opts, "cluster-1", "pool-1")
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

				if err := tt.call(testClient(server.URL), opts); err != nil {
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
