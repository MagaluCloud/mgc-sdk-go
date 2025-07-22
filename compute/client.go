// Package compute provides functionality to interact with the MagaluCloud compute service.
// This package allows managing virtual machine instances, images, instance types, and snapshots.
package compute

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/compute"
)

// VirtualMachineClient represents a client for the compute service.
// It encapsulates functionality to access instances, images, instance types, and snapshots.
type VirtualMachineClient struct {
	*client.CoreClient
}

// ClientOption allows customizing the virtual machine client configuration.
type ClientOption func(*VirtualMachineClient)

// New creates a new instance of VirtualMachineClient.
// If the core client is nil, returns nil.
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

// newRequest creates a new HTTP request for the compute service.
// This method is internal and should not be called directly by SDK users.
func (c *VirtualMachineClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Instances returns a service to manage virtual machine instances.
// This method allows access to functionality such as creating, listing, and managing instances.
func (c *VirtualMachineClient) Instances() InstanceService {
	return &instanceService{client: c}
}

// Images returns a service to manage virtual machine images.
// This method allows access to functionality such as listing available images.
func (c *VirtualMachineClient) Images() ImageService {
	return &imageService{client: c}
}

// InstanceTypes returns a service to manage instance types.
// This method allows access to functionality such as listing available machine types.
func (c *VirtualMachineClient) InstanceTypes() InstanceTypeService {
	return &instanceTypeService{client: c}
}

// Snapshots returns a service to manage instance snapshots.
// This method allows access to functionality such as creating, listing, and managing snapshots.
func (c *VirtualMachineClient) Snapshots() SnapshotService {
	return &snapshotService{client: c}
}
