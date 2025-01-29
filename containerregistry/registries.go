package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	RegistriesService interface {
		Create(ctx context.Context, request *RegistryRequest) (*RegistryResponse, error)
		List(ctx context.Context, opts ListOptions) (*ListRegistriesResponse, error)
		Get(ctx context.Context, registryID string) (*RegistryResponse, error)
		Delete(ctx context.Context, registryID string) error
	}

	RegistryRequest struct {
		Name string `json:"name"`
	}

	RegistryResponse struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Storage   int    `json:"storage_usage_bytes"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	ListOptions struct {
		Limit  *int
		Offset *int
		Sort   *string
		Expand []string
	}

	ListRegistriesResponse struct {
		Registries []RegistryResponse `json:"results"`
	}

	registriesService struct {
		client *ContainerRegistryClient
	}
)

func (c *registriesService) Create(ctx context.Context, request *RegistryRequest) (*RegistryResponse, error) {
	path := "/v0/registries"

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RegistryResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodPost, path, request, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *registriesService) List(ctx context.Context, opts ListOptions) (*ListRegistriesResponse, error) {
	path := "/v0/registries"
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

	if len(opts.Expand) > 0 {
		query.Set("_expand", strings.Join(opts.Expand, ","))
	}

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListRegistriesResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *registriesService) Get(ctx context.Context, registryID string) (*RegistryResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s", registryID)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RegistryResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *registriesService) Delete(ctx context.Context, registryID string) error {
	path := fmt.Sprintf("/v0/registries/%s", registryID)

	err := mgc_http.ExecuteSimpleRequest(ctx, c.client.newRequest, c.client.GetConfig(), http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
