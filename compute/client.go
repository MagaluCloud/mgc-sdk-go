package compute

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

const (
	// DefaultBasePath is the default base path for the virtual machine API
	DefaultBasePath = "/compute"
)

type VirtualMachineClient struct {
	*client.CoreClient
}

type ClientOption func(*VirtualMachineClient)

func New(core *client.CoreClient, opts ...ClientOption) *VirtualMachineClient {
	if core == nil {
		return nil
	}
	vmClient := &VirtualMachineClient{
		CoreClient: core,
	}
	for _, opt := range opts {
		opt(vmClient)
	}
	return vmClient
}

func (c *VirtualMachineClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return c.CoreClient.NewRequest(ctx, method, DefaultBasePath+path, body)
}

func (c *VirtualMachineClient) Instances() InstanceService {
	return &instanceService{client: c}
}

func (c *VirtualMachineClient) Images() ImageService {
	return &imageService{client: c}
}

func (c *VirtualMachineClient) MachineTypes() MachineTypeService {
	return &machineTypeService{client: c}
}

func (c *VirtualMachineClient) Snapshots() SnapshotService {
	return &snapshotService{client: c}
}
