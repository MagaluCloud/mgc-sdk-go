// Package tag provides client implementation for managing tags in the Magalu Cloud platform.
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
	DefaultBasePath = "/tags"
)

// TagClient represents a client for interacting with the tags service
type TagClient struct {
	*client.CoreClient
	baseURL  client.MgcUrl
	tenantID string
}

// ClientOption allows customizing the tag client configuration.
type ClientOption func(*TagClient)

// WithGlobalBasePath allows overriding the default global endpoint for the tags service.
// This is rarely needed as tags are managed globally, but provided for flexibility.
//
// Example:
//
//	client := tag.New(core, tag.WithGlobalBasePath("https://tags.example.com"))
func WithGlobalBasePath(basePath client.MgcUrl) ClientOption {
	return func(c *TagClient) {
		c.baseURL = basePath
	}
}

// WithTenantID sets the x-tenant-id header sent by this client, for accounts that
// own more than one tenant. It takes precedence over a header of the same name
// configured in the core client.
func WithTenantID(tenantID string) ClientOption {
	return func(c *TagClient) {
		c.tenantID = tenantID
	}
}

// New creates a new tag client using the provided core client.
// The tags service operates globally and is not region-specific, so the endpoint
// configured in the core client is not used and is left untouched.
// By default, it uses the global endpoint (api.magalu.cloud).
//
// To customize the endpoint, use the WithGlobalBasePath option.
func New(core *client.CoreClient, opts ...ClientOption) *TagClient {
	if core == nil {
		return nil
	}
	tagClient := &TagClient{
		CoreClient: core,
		baseURL:    client.Global,
	}

	for _, opt := range opts {
		opt(tagClient)
	}
	return tagClient
}

// newRequest creates a new HTTP request with the tags API base path.
// The endpoint is applied to a copy of the configuration, so the core client,
// which is shared with the other services, keeps the endpoint of its own region.
func (c *TagClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	config := *c.GetConfig()
	config.BaseURL = c.baseURL

	req, err := mgc_http.NewRequest(&config, ctx, method, DefaultBasePath+path, &body)
	if err != nil {
		return nil, err
	}

	if c.tenantID != "" {
		req.Header.Set("x-tenant-id", c.tenantID)
	}
	return req, nil
}

// Tags returns a service for managing tag resources
func (c *TagClient) Tags() TagService {
	return &tagService{client: c}
}

// Values returns a service for managing the values of a tag
func (c *TagClient) Values() TagValueService {
	return &tagValueService{client: c}
}

// ResourceTypes returns a service for listing the resource types that support tagging
func (c *TagClient) ResourceTypes() ResourceTypeService {
	return &resourceTypeService{client: c}
}

// Resources returns a service for managing the tags attached to resources
func (c *TagClient) Resources() ResourceService {
	return &resourceService{client: c}
}
