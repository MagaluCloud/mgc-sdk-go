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
	// RegistriesService provides methods for managing container registries
	RegistriesService interface {
		// Create creates a new container registry
		Create(ctx context.Context, request *RegistryRequest) (*RegistryResponse, error)
		// List retrieves a list of container registries with optional filtering and pagination
		List(ctx context.Context, opts ListOptions) (*ListRegistriesResponse, error)
		// Get retrieves a specific container registry by ID
		Get(ctx context.Context, registryID string) (*RegistryResponse, error)
		// Delete removes a container registry by ID
		Delete(ctx context.Context, registryID string) error
	}

	// RegistryRequest represents the request payload for creating a registry
	RegistryRequest struct {
		// Name is the name of the registry to create
		Name string `json:"name"`
	}

	// RegistryResponse represents a container registry
	RegistryResponse struct {
		// ID is the unique identifier of the registry
		ID string `json:"id"`
		// Name is the name of the registry
		Name string `json:"name"`
		// Storage is the storage usage in bytes
		Storage int `json:"storage_usage_bytes"`
		// CreatedAt is the timestamp when the registry was created
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the timestamp when the registry was last updated
		UpdatedAt string `json:"updated_at"`
	}

	// ListOptions provides options for listing registries
	ListOptions struct {
		// Limit is the maximum number of registries to return
		Limit *int
		// Offset is the number of registries to skip
		Offset *int
		// Sort is the field to sort by
		Sort *string
		// Expand is a list of fields to expand in the response
		Expand []string
	}

	// ListRegistriesResponse represents the response when listing registries
	ListRegistriesResponse struct {
		// Registries is the list of registries
		Registries []RegistryResponse `json:"results"`
	}

	// registriesService implements the RegistriesService interface
	registriesService struct {
		client *ContainerRegistryClient
	}
)

// Create creates a new container registry
func (c *registriesService) Create(ctx context.Context, request *RegistryRequest) (*RegistryResponse, error) {
	path := "/v0/registries"

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RegistryResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodPost, path, request, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// List retrieves a list of container registries with optional filtering and pagination
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

// Get retrieves a specific container registry by ID
func (c *registriesService) Get(ctx context.Context, registryID string) (*RegistryResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s", registryID)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[RegistryResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete removes a container registry by ID
func (c *registriesService) Delete(ctx context.Context, registryID string) error {
	path := fmt.Sprintf("/v0/registries/%s", registryID)

	err := mgc_http.ExecuteSimpleRequest(ctx, c.client.newRequest, c.client.GetConfig(), http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
