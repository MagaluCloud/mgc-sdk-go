package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// RepositoriesService provides methods for managing repositories within container registries
	RepositoriesService interface {
		// List retrieves a list of repositories within a registry with optional filtering and pagination
		List(ctx context.Context, registryID string, opts ListOptions) (*RepositoriesResponse, error)
		// Get retrieves a specific repository within a registry
		Get(ctx context.Context, registryID, repositoryName string) (*RepositoryResponse, error)
		// Delete removes a repository within a registry
		Delete(ctx context.Context, registryID, repositoryName string) error
	}

	// RepositoryResponse represents a repository within a container registry
	RepositoryResponse struct {
		// RegistryName is the name of the parent registry
		RegistryName string `json:"registry_name"`
		// Name is the name of the repository
		Name string `json:"name"`
		// ImageCount is the number of images in the repository
		ImageCount int `json:"image_count"`
		// CreatedAt is the timestamp when the repository was created
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the timestamp when the repository was last updated
		UpdatedAt string `json:"updated_at"`
	}

	// AmountRepositoryResponse represents the total count of repositories
	AmountRepositoryResponse struct {
		// Total is the total number of repositories
		Total int `json:"total"`
	}

	// RepositoriesResponse represents the response when listing repositories
	RepositoriesResponse struct {
		// Goal contains the total count information
		Goal AmountRepositoryResponse `json:"goal"`
		// Results is the list of repositories
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

	query := make(url.Values)

	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}

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
