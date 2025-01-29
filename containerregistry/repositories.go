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
	RepositoriesService interface {
		List(ctx context.Context, registryID string, opts ListOptions) (*RepositoriesResponse, error)
		Get(ctx context.Context, registryID, repositoryName string) (*RepositoryResponse, error)
		Delete(ctx context.Context, registryID, repositoryName string) error
	}

	RepositoryResponse struct {
		RegistryName string `json:"registry_name"`
		Name         string `json:"name"`
		ImageCount   int    `json:"image_count"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}

	AmountRepositoryResponse struct {
		Total int `json:"total"`
	}

	RepositoriesResponse struct {
		Goal    AmountRepositoryResponse `json:"goal"`
		Results []RepositoryResponse     `json:"results"`
	}

	repositoriesService struct {
		client *ContainerRegistryClient
	}
)

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

func (c *repositoriesService) Get(ctx context.Context, registryID, repositoryName string) (*RepositoryResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s", registryID, repositoryName)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RepositoryResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *repositoriesService) Delete(ctx context.Context, registryID, repositoryName string) error {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s", registryID, repositoryName)

	err := mgc_http.ExecuteSimpleRequest(ctx, c.client.newRequest, c.client.GetConfig(), http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
