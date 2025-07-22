// Package availabilityzones provides functionality to interact with the MagaluCloud availability zones service.
// This package allows listing availability zones across different regions.
package availabilityzones

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// DefaultBasePath defines the default base path for availability zones APIs.
const (
	DefaultBasePath = "/profile"
)

// Client handles operations on availability zones in the Magalu Cloud platform.
// Availability zones are managed as a global service, meaning they are not bound to any specific region.
// By default, the service uses the global endpoint.
type Client struct {
	*client.CoreClient
}

// ClientOption allows customizing the availability zones client configuration
type ClientOption func(*Client)

// WithGlobalBasePath allows overriding the default global endpoint for availability zones service.
// This is rarely needed as availability zones are managed globally, but provided for flexibility.
//
// Example:
//
//	client := availabilityzones.New(core, availabilityzones.WithGlobalBasePath("custom-endpoint"))
func WithGlobalBasePath(basePath client.MgcUrl) ClientOption {
	return func(c *Client) {
		c.GetConfig().BaseURL = basePath
	}
}

// New creates a new availability zones client using the provided core client.
// The availability zones service operates globally and is not region-specific.
// By default, it uses the global endpoint (api.magalu.cloud).
//
// To customize the endpoint, use WithGlobalBasePath option.
func New(core *client.CoreClient, opts ...ClientOption) *Client {
	if core == nil {
		return nil
	}
	azClient := &Client{
		CoreClient: core,
	}

	azClient.GetConfig().BaseURL = client.Global

	for _, opt := range opts {
		opt(azClient)
	}
	return azClient
}

// newRequest creates a new HTTP request for the availability zones service.
// This method is internal and should not be called directly by SDK users.
func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// AvailabilityZones returns a service to manage availability zones.
// This method allows access to functionality such as listing availability zones.
func (c *Client) AvailabilityZones() Service {
	return &service{client: c}
}
