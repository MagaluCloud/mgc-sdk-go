// Package tag provides client implementation for managing Tags in the Magalu Cloud platform.
// Tags are managed as a global service, meaning they are not bound to any specific region.
// By default, the service uses the global endpoint, but this can be overridden if needed.
package tag

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	// DefaultBasePath is the default API base path for tag operations
	DefaultBasePath = ""
)

// TagClient represents a client for interacting with the tags service
type TagClient struct {
	*client.CoreClient
}

// ClientOption allows customizing the tag client configuration.
type ClientOption func(*TagClient)

// WithBasePath allows overriding the default base path for the tags service.
func WithBasePath(basePath client.MgcUrl) ClientOption {
	return func(c *TagClient) {
		c.GetConfig().BaseURL = basePath
	}
}

// New creates a new tag client using the provided core client.
// The tags service operates globally and is not region-specific.
// By default, it uses the global endpoint (api.magalu.cloud).
//
// To customize the endpoint, use WithBasePath option.
func New(core *client.CoreClient, opts ...ClientOption) *TagClient {
	if core == nil {
		return nil
	}
	tagClient := &TagClient{
		CoreClient: core,
	}

	tagClient.GetConfig().BaseURL = client.Global

	for _, opt := range opts {
		opt(tagClient)
	}
	return tagClient
}

// newRequest creates a new HTTP request with the tags API base path
func (c *TagClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Tags returns a service for managing tag resources
func (c *TagClient) Tags() TagService {
	return &tagService{client: c}
}

// Values returns a service for managing tag value resources
func (c *TagClient) Values() TagValueService {
	return &tagValueService{client: c}
}

// Resources returns a service for managing links between tag values and cloud resources
func (c *TagClient) Resources() TagValueResourceService {
	return &tagValueResourceService{client: c}
}

// ResourceTypes returns a service for listing available cloud resource types
func (c *TagClient) ResourceTypes() ResourceTypeService {
	return &resourceTypeService{client: c}
}
