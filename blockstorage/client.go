// Package blockstorage provides functionality to interact with the MagaluCloud block storage service.
// This package allows managing volumes, volume types, and snapshots.
package blockstorage

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// DefaultBasePath defines the default base path for block storage APIs.
const (
	DefaultBasePath = "/volume"
)

// BlockStorageClient represents a client for the block storage service.
// It encapsulates functionality to access volumes, volume types, and snapshots.
type BlockStorageClient struct {
	*client.CoreClient
}

// ClientOption allows customizing the block storage client configuration.
type ClientOption func(*BlockStorageClient)

// New creates a new instance of BlockStorageClient.
// If the core client is nil, returns nil.
func New(core *client.CoreClient, opts ...ClientOption) *BlockStorageClient {
	if core == nil {
		return nil
	}
	bsClient := &BlockStorageClient{
		CoreClient: core,
	}
	for _, opt := range opts {
		opt(bsClient)
	}
	return bsClient
}

// newRequest creates a new HTTP request for the block storage service.
// This method is internal and should not be called directly by SDK users.
func (c *BlockStorageClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Volumes returns a service to manage block storage volumes.
// This method allows access to functionality such as creating, listing, and managing volumes.
func (c *BlockStorageClient) Volumes() VolumeService {
	return &volumeService{client: c}
}

// VolumeTypes returns a service to manage volume types.
// This method allows access to functionality such as listing available volume types.
func (c *BlockStorageClient) VolumeTypes() VolumeTypeService {
	return &volumeTypeService{client: c}
}

// Snapshots returns a service to manage volume snapshots.
// This method allows access to functionality such as creating, listing, and managing snapshots.
func (c *BlockStorageClient) Snapshots() SnapshotService {
	return &snapshotService{client: c}
}
