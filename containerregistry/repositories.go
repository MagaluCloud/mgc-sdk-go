package containerregistry

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// RepositoriesService provides methods for managing repositories within container registries
	RepositoriesService interface {
		List(ctx context.Context, registryID string, opts ListOptions) (*RepositoriesResponse, error)
		Get(ctx context.Context, registryID, repositoryName string) (*RepositoryResponse, error)
		Delete(ctx context.Context, registryID, repositoryName string) error
	}

	// RepositoryResponse represents a repository within a container registry
	RepositoryResponse struct {
		RegistryName string `json:"registry_name"`
		Name         string `json:"name"`
		ImageCount   int    `json:"image_count"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}

	// AmountRepositoryResponse represents the total count of repositories
	AmountRepositoryResponse struct {
		Total int `json:"total"`
	}

	// RepositoriesResponse represents the response when listing repositories
	RepositoriesResponse struct {
		Meta    Meta                 `json:"meta"`
		Results []RepositoryResponse `json:"results"`
	}

	// repositoriesService implements the RepositoriesService interface
	repositoriesService struct {
		client *ContainerRegistryClient
	}
)

// List retrieves a list of repositories within a registry with optional filtering and pagination
func (c *repositoriesService) List(ctx context.Context, registryID string, opts ListOptions) (*RepositoriesResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories", registryID)
	query := CreatePaginationParams(opts)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RepositoriesResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Get retrieves a specific repository within a registry
func (c *repositoriesService) Get(ctx context.Context, registryID, repositoryName string) (*RepositoryResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s", registryID, repositoryName)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RepositoryResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete removes a repository within a registry
func (c *repositoriesService) Delete(ctx context.Context, registryID, repositoryName string) error {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s", registryID, repositoryName)

	err := mgc_http.ExecuteSimpleRequest(ctx, c.client.newRequest, c.client.GetConfig(), http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
