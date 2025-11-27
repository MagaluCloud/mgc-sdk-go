package containerregistry

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ProxyCache represents a proxy-cache.
// A proxy-cache is an intermediate mirror that fetches images from an external registry and stores them locally
type ProxyCache struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Provider  string    `json:"provider"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListProxyCachesResponse represents the response from listing proxy-caches.
// This structure encapsulates the API response format for proxy-caches.
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

// ProxyCachesService provides operations for managing proxy-caches.
// This interface allows creating, listing, deleting, and managing proxy-caches.
type ProxyCachesService interface {
	List(ctx context.Context, opts ProxyCacheListOptions) (*ListProxyCachesResponse, error)
}

// proxyCachesService implements the ProxyCachesService interface.
// This is an internal implementation that should not be used directly.
type proxyCachesService struct {
	client *ContainerRegistryClient
}

// List returns a paginated list of proxy-caches.
// This method makes an HTTP request to get the list of proxy-caches and applies the filters specified in the options.
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

	path := "/v0/proxy-caches"

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListProxyCachesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		q,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
