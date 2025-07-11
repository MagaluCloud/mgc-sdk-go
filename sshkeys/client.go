// Package sshkeys provides client implementation for managing SSH keys in the Magalu Cloud platform.
// SSH keys are managed as a global service, meaning they are not bound to any specific region.
// By default, the service uses the global endpoint, but this can be overridden if needed.
package sshkeys

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	// DefaultBasePath is the default API base path for SSH key operations
	DefaultBasePath = "/profile"
)

// SSHKeyClient represents a client for interacting with the SSH keys service
type SSHKeyClient struct {
	*client.CoreClient
}

// ClientOption allows customizing the SSH key client configuration.
type ClientOption func(*SSHKeyClient)

// WithGlobalBasePath allows overriding the default global endpoint for SSH keys service.
// This is rarely needed as SSH keys are managed globally, but provided for flexibility.
//
// Example:
//
//	client := sshkeys.New(core, sshkeys.WithGlobalBasePath("custom-endpoint"))
func WithGlobalBasePath(basePath client.MgcUrl) ClientOption {
	return func(c *SSHKeyClient) {
		c.GetConfig().BaseURL = basePath
	}
}

// New creates a new SSH key client using the provided core client.
// The SSH keys service operates globally and is not region-specific.
// By default, it uses the global endpoint (api.magalu.cloud).
//
// To customize the endpoint, use WithGlobalBasePath option.
func New(core *client.CoreClient, opts ...ClientOption) *SSHKeyClient {
	if core == nil {
		return nil
	}
	sshClient := &SSHKeyClient{
		CoreClient: core,
	}

	sshClient.GetConfig().BaseURL = client.Global

	for _, opt := range opts {
		opt(sshClient)
	}
	return sshClient
}

// newRequest creates a new HTTP request with the SSH keys API base path
func (c *SSHKeyClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Keys returns a service for managing SSH key resources
func (c *SSHKeyClient) Keys() KeyService {
	return &keyService{client: c}
}
