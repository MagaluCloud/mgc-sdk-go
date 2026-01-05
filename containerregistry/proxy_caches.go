package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ProxyCache represents a proxy-cache.
// A proxy-cache is an intermediate mirror that fetches images from an external registry and stores them locally
type ProxyCache struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListProxyCachesResponse struct {
	Meta    Meta         `json:"meta"`
	Results []ProxyCache `json:"results"`
}

// ProxyCacheListOptions contains options for listing proxy-caches.
// All fields are optional and allow controlling pagination.
type ProxyCacheListOptions struct {
	Limit  *int
	Offset *int
	Sort   *string
}

// ProxyCacheListAllOptions provides options for ListAll (without pagination)
type ProxyCacheListAllOptions struct {
	Sort *string
}

type CreateProxyCacheRequest struct {
	Name        string  `json:"name"`
	Provider    string  `json:"provider"`
	URL         string  `json:"url"`
	Description *string `json:"description"`
}

type CreateProxyCacheResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateProxyCacheRequest struct {
	Name        *string `json:"name"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
}

type GetProxyCacheResponse struct {
	ProxyCache
	Description string `json:"description"`
}

type ListProxyCacheStatusResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type CreateProxyCacheStatusRequest struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
}

type CreateProxyCacheStatusResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// This interface allows creating, listing, deleting, and managing proxy-caches.
type ProxyCachesService interface {
	List(ctx context.Context, opts ProxyCacheListOptions) (*ListProxyCachesResponse, error)
	ListAll(ctx context.Context, opts ProxyCacheListAllOptions) ([]ProxyCache, error)
	Create(ctx context.Context, req CreateProxyCacheRequest) (*CreateProxyCacheResponse, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, req UpdateProxyCacheRequest) (*ProxyCache, error)
	Get(ctx context.Context, id string) (*GetProxyCacheResponse, error)
	ListStatus(ctx context.Context, id string) (*ListProxyCacheStatusResponse, error)
	CreateStatus(ctx context.Context, req CreateProxyCacheStatusRequest) (*CreateProxyCacheStatusResponse, error)
}

// proxyCachesService implements the ProxyCachesService interface.
// This is an internal implementation that should not be used directly.
type proxyCachesService struct {
	client *ContainerRegistryClient
}

const basePath = "/v0/proxy-caches"
const pathWithID = basePath + "/%s"

// This method makes a HTTP request to get the list of proxy-caches and applies the filters specified in the options.
func (s *proxyCachesService) List(ctx context.Context, opts ProxyCacheListOptions) (*ListProxyCachesResponse, error) {
	q := url.Values{}

	if opts.Limit != nil {
		q.Add("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Add("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Add("_sort", *opts.Sort)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListProxyCachesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		basePath,
		nil,
		q,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListAll retrieves all proxy-caches by fetching all pages with optional filtering
func (s *proxyCachesService) ListAll(ctx context.Context, opts ProxyCacheListAllOptions) ([]ProxyCache, error) {
	var allProxyCaches []ProxyCache
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit

		opts := ProxyCacheListOptions{
			Offset: &currentOffset,
			Limit:  &currentLimit,
			Sort:   opts.Sort,
		}

		resp, err := s.List(ctx, opts)

		if err != nil {
			return nil, err
		}

		allProxyCaches = append(allProxyCaches, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allProxyCaches, nil
}

// This method makes a HTTP request to create a new proxy-cache and returns the ID and name of the created proxy-cache.
func (s *proxyCachesService) Create(ctx context.Context, req CreateProxyCacheRequest) (*CreateProxyCacheResponse, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[CreateProxyCacheResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		basePath,
		req,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// This method makes a HTTP request to delete a proxy-cache permanently.
func (s *proxyCachesService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(pathWithID, id),
		nil,
		nil,
	)
}

// This method makes a HTTP request to update an existing proxy-cache.
func (s *proxyCachesService) Update(ctx context.Context, id string, req UpdateProxyCacheRequest) (*ProxyCache, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ProxyCache](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(pathWithID, id),
		req,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// This method makes a HTTP request to get detailed informations about a proxy-cache.
func (s *proxyCachesService) Get(ctx context.Context, id string) (*GetProxyCacheResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[GetProxyCacheResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(pathWithID, id),
		nil,
		nil,
	)
}

// This method makes a HTTP request to get the status of a proxy-cache.
func (s *proxyCachesService) ListStatus(ctx context.Context, id string) (*ListProxyCacheStatusResponse, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListProxyCacheStatusResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/proxy-caches/%s/status", id),
		nil,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// This method makes a HTTP request to validates the provided credentials and endpoint information.
func (s *proxyCachesService) CreateStatus(ctx context.Context, req CreateProxyCacheStatusRequest) (*CreateProxyCacheStatusResponse, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[CreateProxyCacheStatusResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v0/proxy-caches/status",
		req,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
