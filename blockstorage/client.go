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
//
// Parameters:
//   - core: The core client that will be used for HTTP requests
//   - opts: Optional configuration options for the client
//
// Returns:
//   - *BlockStorageClient: A new instance of the block storage client or nil if core is nil
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
//
// Parameters:
//   - ctx: Request context
//   - method: HTTP method (GET, POST, etc.)
//   - path: API path (will be concatenated with DefaultBasePath)
//   - body: Request body (can be nil)
//
// Returns:
//   - *http.Request: The created HTTP request
//   - error: Error if there's a failure creating the request
func (c *BlockStorageClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Volumes returns a service to manage block storage volumes.
// This method allows access to functionality such as creating, listing, and managing volumes.
//
// Returns:
//   - VolumeService: Interface for volume operations
func (c *BlockStorageClient) Volumes() VolumeService {
	return &volumeService{client: c}
}

// VolumeTypes returns a service to manage volume types.
// This method allows access to functionality such as listing available volume types.
//
// Returns:
//   - VolumeTypeService: Interface for volume type operations
func (c *BlockStorageClient) VolumeTypes() VolumeTypeService {
	return &volumeTypeService{client: c}
}

// Snapshots returns a service to manage volume snapshots.
// This method allows access to functionality such as creating, listing, and managing snapshots.
//
// Returns:
//   - SnapshotService: Interface for snapshot operations
func (c *BlockStorageClient) Snapshots() SnapshotService {
	return &snapshotService{client: c}
}
