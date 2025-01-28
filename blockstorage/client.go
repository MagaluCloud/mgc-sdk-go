package blockstorage

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/volume"
)

type BlockStorageClient struct {
	*client.CoreClient
}

type ClientOption func(*BlockStorageClient)

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

func (c *BlockStorageClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

func (c *BlockStorageClient) Volumes() VolumeService {
	return &volumeService{client: c}
}

func (c *BlockStorageClient) VolumeTypes() VolumeTypeService {
	return &volumeTypeService{client: c}
}

func (c *BlockStorageClient) Snapshots() SnapshotService {
	return &snapshotService{client: c}
}
