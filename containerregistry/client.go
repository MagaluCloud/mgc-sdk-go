// Package containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
// This package allows you to manage container registries, repositories, images, and credentials.
package containerregistry

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/container-registry"
)

// ContainerRegistryClient represents a client for the Container Registry service
type ContainerRegistryClient struct {
	*client.CoreClient
}

// ClientOption is a function type for configuring ContainerRegistryClient options
type ClientOption func(*ContainerRegistryClient)

// New creates a new ContainerRegistryClient instance with the provided core client and options
func New(core *client.CoreClient, opts ...ClientOption) *ContainerRegistryClient {
	if core == nil {
		return nil
	}
	crClient := &ContainerRegistryClient{
		CoreClient: core,
	}
	for _, opt := range opts {
		opt(crClient)
	}
	return crClient
}

// newRequest creates a new HTTP request for the container registry API
func (c *ContainerRegistryClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Credentials returns a service for managing container registry credentials
func (c *ContainerRegistryClient) Credentials() CredentialsService {
	return &credentialsService{client: c}
}

// Registries returns a service for managing container registries
func (c *ContainerRegistryClient) Registries() RegistriesService {
	return &registriesService{client: c}
}

// Repositories returns a service for managing repositories within registries
func (c *ContainerRegistryClient) Repositories() RepositoriesService {
	return &repositoriesService{client: c}
}

// Images returns a service for managing images within repositories
func (c *ContainerRegistryClient) Images() ImagesService {
	return &imagesService{client: c}
}

// ProxyCaches returns a service for managing proxy-caches
func (c *ContainerRegistryClient) ProxyCaches() ProxyCachesService {
	return &proxyCachesService{client: c}
}
