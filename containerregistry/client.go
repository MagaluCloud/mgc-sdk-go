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

type ContainerRegistryClient struct {
	*client.CoreClient
}

type ClientOption func(*ContainerRegistryClient)

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

func (c *ContainerRegistryClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

func (c *ContainerRegistryClient) Credentials() CredentialsService {
	return &credentialsService{client: c}
}

func (c *ContainerRegistryClient) Registries() RegistriesService {
	return &registriesService{client: c}
}

func (c *ContainerRegistryClient) Repositories() RepositoriesService {
	return &repositoriesService{client: c}
}

func (c *ContainerRegistryClient) Images() ImagesService {
	return &imagesService{client: c}
}
