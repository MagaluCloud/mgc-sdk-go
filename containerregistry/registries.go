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
	// RegistriesService provides methods for managing container registries
	RegistriesService interface {
		Create(ctx context.Context, request *RegistryRequest) (*RegistryResponse, error)
		List(ctx context.Context, opts RegistryListOptions) (*ListRegistriesResponse, error)
		ListAll(ctx context.Context, filterOpts RegistryFilterOptions) ([]RegistryResponse, error)
		Get(ctx context.Context, registryID string) (*RegistryResponse, error)
		Delete(ctx context.Context, registryID string) error
	}

	// RegistryRequest represents the request payload for creating a registry
	RegistryRequest struct {
		Name string `json:"name"`
	}

	// RegistryResponse represents a container registry
	RegistryResponse struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Storage   int    `json:"storage_usage_bytes"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	// RegistryListOptions provides options for listing registries with pagination
	RegistryListOptions struct {
		Offset *int
		Limit  *int
		RegistryFilterOptions
	}

	// RegistryFilterOptions provides filtering options for registries
	RegistryFilterOptions struct {
		Sort *string
	}

	// ListRegistriesResponse represents the response when listing registries
	ListRegistriesResponse = helpers.PaginatedResponse[RegistryResponse]

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
func (c *registriesService) List(ctx context.Context, opts RegistryListOptions) (*ListRegistriesResponse, error) {
	path := "/v0/registries"
	query := c.createRegistryQueryParams(opts)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListRegistriesResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListAll retrieves all registries across all pages with optional filtering
func (c *registriesService) ListAll(ctx context.Context, filterOpts RegistryFilterOptions) ([]RegistryResponse, error) {
	var allRegistries []RegistryResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		opts := RegistryListOptions{
			Offset:                &currentOffset,
			Limit:                 &currentLimit,
			RegistryFilterOptions: filterOpts,
		}

		result, err := c.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allRegistries = append(allRegistries, result.Results...)

		// Check if we've retrieved all results
		if len(result.Results) < limit {
			break
		}

		offset += limit
	}

	return allRegistries, nil
}

// createRegistryQueryParams creates URL query parameters from RegistryListOptions
func (c *registriesService) createRegistryQueryParams(opts RegistryListOptions) url.Values {
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
