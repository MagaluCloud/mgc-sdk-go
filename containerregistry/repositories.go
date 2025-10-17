package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// RepositoriesService provides methods for managing repositories within container registries
	RepositoriesService interface {
		List(ctx context.Context, registryID string, opts RepositoryListOptions) (*RepositoriesResponse, error)
		ListAll(ctx context.Context, registryID string, filterOpts RepositoryFilterOptions) ([]RepositoryResponse, error)
		Get(ctx context.Context, registryID, repositoryName string) (*RepositoryResponse, error)
		Delete(ctx context.Context, registryID, repositoryName string) error
	}

	// RepositoryListOptions provides options for listing repositories with pagination
	RepositoryListOptions struct {
		Offset *int
		Limit  *int
		RepositoryFilterOptions
	}

	// RepositoryFilterOptions provides filtering options for repositories
	RepositoryFilterOptions struct {
		Sort *string
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
	RepositoriesResponse = helpers.PaginatedResponse[RepositoryResponse]

	// repositoriesService implements the RepositoriesService interface
	repositoriesService struct {
		client *ContainerRegistryClient
	}
)

// List retrieves a list of repositories within a registry with optional filtering and pagination
func (c *repositoriesService) List(ctx context.Context, registryID string, opts RepositoryListOptions) (*RepositoriesResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories", registryID)
	query := c.createRepositoryQueryParams(opts)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RepositoriesResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListAll retrieves all repositories within a registry across all pages with optional filtering
func (c *repositoriesService) ListAll(ctx context.Context, registryID string, filterOpts RepositoryFilterOptions) ([]RepositoryResponse, error) {
	var allRepositories []RepositoryResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		opts := RepositoryListOptions{
			Offset:                  &currentOffset,
			Limit:                   &currentLimit,
			RepositoryFilterOptions: filterOpts,
		}

		result, err := c.List(ctx, registryID, opts)
		if err != nil {
			return nil, err
		}

		allRepositories = append(allRepositories, result.Results...)

		if len(result.Results) < limit {
			break
		}

		offset += limit
	}

	return allRepositories, nil
}

// createRepositoryQueryParams creates URL query parameters from RepositoryListOptions
func (c *repositoriesService) createRepositoryQueryParams(opts RepositoryListOptions) url.Values {
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

	return query
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
