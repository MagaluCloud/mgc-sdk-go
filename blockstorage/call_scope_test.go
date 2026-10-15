package blockstorage

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
		call func(c *BlockStorageClient, opts client.CallOptions) error
	}{
		{"Volumes.List", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Volumes().List(ctx, opts, ListOptions{})
			return err
		}},
		{"Volumes.ListAll", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Volumes().ListAll(ctx, opts, VolumeFilterOptions{})
			return err
		}},
		{"Volumes.Create", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Volumes().Create(ctx, opts, CreateVolumeRequest{})
			return err
		}},
		{"Volumes.Get", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Volumes().Get(ctx, opts, "vol-1", nil)
			return err
		}},
		{"Volumes.Delete", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Volumes().Delete(ctx, opts, "vol-1")
		}},
		{"Volumes.Rename", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Volumes().Rename(ctx, opts, "vol-1", "new-name")
		}},
		{"Volumes.Extend", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Volumes().Extend(ctx, opts, "vol-1", ExtendVolumeRequest{})
		}},
		{"Volumes.Retype", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Volumes().Retype(ctx, opts, "vol-1", RetypeVolumeRequest{})
		}},
		{"Volumes.Attach", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Volumes().Attach(ctx, opts, "vol-1", "inst-1")
		}},
		{"Volumes.Detach", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Volumes().Detach(ctx, opts, "vol-1")
		}},
		{"Snapshots.List", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Snapshots().List(ctx, opts, SnaphotListOptions{})
			return err
		}},
		{"Snapshots.ListAll", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Snapshots().ListAll(ctx, opts, SnapshotFilterOptions{})
			return err
		}},
		{"Snapshots.Create", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Snapshots().Create(ctx, opts, CreateSnapshotRequest{})
			return err
		}},
		{"Snapshots.Get", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Snapshots().Get(ctx, opts, "snap-1", nil)
			return err
		}},
		{"Snapshots.Delete", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Snapshots().Delete(ctx, opts, "snap-1")
		}},
		{"Snapshots.Rename", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Snapshots().Rename(ctx, opts, "snap-1", "new-name")
		}},
		{"Snapshots.Copy", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Snapshots().Copy(ctx, opts, "snap-1", "br-se1")
		}},
		{"Schedulers.List", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Schedulers().List(ctx, opts, SchedulerListOptions{})
			return err
		}},
		{"Schedulers.ListAll", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Schedulers().ListAll(ctx, opts, SchedulerFilterOptions{})
			return err
		}},
		{"Schedulers.Create", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Schedulers().Create(ctx, opts, SchedulerPayload{})
			return err
		}},
		{"Schedulers.Get", func(c *BlockStorageClient, opts client.CallOptions) error {
			_, err := c.Schedulers().Get(ctx, opts, "sched-1", nil)
			return err
		}},
		{"Schedulers.Delete", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Schedulers().Delete(ctx, opts, "sched-1")
		}},
		{"Schedulers.AttachVolume", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Schedulers().AttachVolume(ctx, opts, "sched-1", SchedulerVolumeIdentifierPayload{})
		}},
		{"Schedulers.DetachVolume", func(c *BlockStorageClient, opts client.CallOptions) error {
			return c.Schedulers().DetachVolume(ctx, opts, "sched-1", SchedulerVolumeIdentifierPayload{})
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
