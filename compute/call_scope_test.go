package compute

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
		call func(c *VirtualMachineClient, opts client.CallOptions) error
	}{
		{"Instances.List", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Instances().List(ctx, opts, ListOptions{})
			return err
		}},
		{"Instances.ListAll", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Instances().ListAll(ctx, opts, InstanceFilterOptions{})
			return err
		}},
		{"Instances.Create", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Instances().Create(ctx, opts, CreateRequest{})
			return err
		}},
		{"Instances.Get", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Instances().Get(ctx, opts, "inst-1", nil)
			return err
		}},
		{"Instances.Delete", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().Delete(ctx, opts, "inst-1", false)
		}},
		{"Instances.Rename", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().Rename(ctx, opts, "inst-1", "new-name")
		}},
		{"Instances.Retype", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().Retype(ctx, opts, "inst-1", RetypeRequest{})
		}},
		{"Instances.Start", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().Start(ctx, opts, "inst-1")
		}},
		{"Instances.Stop", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().Stop(ctx, opts, "inst-1")
		}},
		{"Instances.Suspend", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().Suspend(ctx, opts, "inst-1")
		}},
		{"Instances.GetFirstWindowsPassword", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Instances().GetFirstWindowsPassword(ctx, opts, "inst-1")
			return err
		}},
		{"Instances.AttachNetworkInterface", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().AttachNetworkInterface(ctx, opts, NICRequest{})
		}},
		{"Instances.DetachNetworkInterface", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Instances().DetachNetworkInterface(ctx, opts, NICRequest{})
		}},
		{"Instances.InitLog", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Instances().InitLog(ctx, opts, "inst-1", nil)
			return err
		}},
		{"Images.CreateCustom", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Images().CreateCustom(ctx, opts, CreateCustomImageRequest{})
			return err
		}},
		{"Images.GetCustom", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Images().GetCustom(ctx, opts, "img-1")
			return err
		}},
		{"Images.ListCustom", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Images().ListCustom(ctx, opts, CustomImageListOptions{})
			return err
		}},
		{"Images.DeleteCustom", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Images().DeleteCustom(ctx, opts, "img-1")
		}},
		{"Images.UpdateCustom", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Images().UpdateCustom(ctx, opts, "img-1", UpdateCustomImageRequest{})
		}},
		{"Snapshots.List", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Snapshots().List(ctx, opts, SnapshotListOptions{})
			return err
		}},
		{"Snapshots.ListAll", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Snapshots().ListAll(ctx, opts, SnapshotFilterOptions{})
			return err
		}},
		{"Snapshots.Create", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Snapshots().Create(ctx, opts, CreateSnapshotRequest{})
			return err
		}},
		{"Snapshots.Get", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Snapshots().Get(ctx, opts, "snap-1", nil)
			return err
		}},
		{"Snapshots.Delete", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Snapshots().Delete(ctx, opts, "snap-1")
		}},
		{"Snapshots.Rename", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Snapshots().Rename(ctx, opts, "snap-1", "new-name")
		}},
		{"Snapshots.Restore", func(c *VirtualMachineClient, opts client.CallOptions) error {
			_, err := c.Snapshots().Restore(ctx, opts, "snap-1", RestoreSnapshotRequest{})
			return err
		}},
		{"Snapshots.Copy", func(c *VirtualMachineClient, opts client.CallOptions) error {
			return c.Snapshots().Copy(ctx, opts, "snap-1", CopySnapshotRequest{})
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
