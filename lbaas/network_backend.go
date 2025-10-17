package lbaas

import (
	"context"
	"net/http"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const backends = "backends"

type (
	NetworkBackendInstanceTargetRequest struct {
		NicID     *string `json:"nic_id,omitempty"`
		IPAddress *string `json:"ip_address,omitempty"`
		Port      int64   `json:"port"`
	}

	CreateBackendRequest struct {
		HealthCheckName                     *string                                `json:"health_check_name,omitempty"`
		Name                                string                                 `json:"name"`
		Description                         *string                                `json:"description,omitempty"`
		BalanceAlgorithm                    BackendBalanceAlgorithm                `json:"balance_algorithm"`
		PanicThreshold                      *float64                               `json:"panic_threshold,omitempty"`
		TargetsType                         BackendType                            `json:"targets_type"`
		Targets                             *[]NetworkBackendInstanceTargetRequest `json:"targets,omitempty"`
		CloseConnectionsOnHostHealthFailure *bool                                  `json:"close_connections_on_host_health_failure,omitempty"`
		HealthCheckID                       *string                                `json:"health_check_id,omitempty"`
	}

	UpdateNetworkBackendRequest struct {
		HealthCheckID                       *string `json:"health_check_id,omitempty"`
		PanicThreshold                      *int    `json:"panic_threshold,omitempty"`
		CloseConnectionsOnHostHealthFailure *bool   `json:"close_connections_on_host_health_failure,omitempty"`
	}

	NetworkBackedTarget struct {
		ID        string    `json:"id"`
		IPAddress *string   `json:"ip_address,omitempty"`
		NicID     *string   `json:"nic_id,omitempty"`
		Port      *int64    `json:"port,omitempty"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	NetworkBackendResponse struct {
		ID                                  string                  `json:"id"`
		HealthCheckID                       *string                 `json:"health_check_id,omitempty"`
		Name                                string                  `json:"name"`
		Description                         *string                 `json:"description,omitempty"`
		BalanceAlgorithm                    BackendBalanceAlgorithm `json:"balance_algorithm"`
		PanicThreshold                      *int                    `json:"panic_threshold,omitempty"`
		CloseConnectionsOnHostHealthFailure *bool                   `json:"close_connections_on_host_health_failure,omitempty"`
		TargetsType                         BackendType             `json:"targets_type"`
		Targets                             []NetworkBackedTarget   `json:"targets"`
		CreatedAt                           time.Time               `json:"created_at"`
		UpdatedAt                           time.Time               `json:"updated_at"`
	}

	NetworkPaginatedBackendResponse struct {
		Meta    PaginationMeta           `json:"meta"`
		Results []NetworkBackendResponse `json:"results"`
	}

	NetworkBackendService interface {
		Create(ctx context.Context, lbID string, req CreateBackendRequest) (string, error)
		Delete(ctx context.Context, lbID, backendID string) error
		Get(ctx context.Context, lbID, backendID string) (*NetworkBackendResponse, error)
		List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) (NetworkPaginatedBackendResponse, error)
		ListAll(ctx context.Context, lbID string) ([]NetworkBackendResponse, error)
		Update(ctx context.Context, lbID, backendID string, req UpdateNetworkBackendRequest) (string, error)
	}

	// networkBackendService implements the NetworkBackendService interface
	networkBackendService struct {
		client *LbaasClient
	}
)

// Create creates a new network backend
func (s *networkBackendService) Create(ctx context.Context, lbID string, req CreateBackendRequest) (string, error) {
	path := urlNetworkLoadBalancer(&lbID, backends)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return "", err
	}

	var resp struct {
		ID string `json:"id"`
	}
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Delete removes a network backend
func (s *networkBackendService) Delete(ctx context.Context, lbID, backendID string) error {
	path := urlNetworkLoadBalancer(&lbID, backends, backendID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific backend
func (s *networkBackendService) Get(ctx context.Context, lbID, backendID string) (*NetworkBackendResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, backends, backendID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkBackendResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List returns a paginated list of network backends
func (s *networkBackendService) List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) (NetworkPaginatedBackendResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, backends)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return NetworkPaginatedBackendResponse{}, err
	}

	// Add query parameters for pagination and sorting
	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", options.Offset)
	query.AddReflect("_limit", options.Limit)
	query.Add("_sort", options.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedBackendResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return NetworkPaginatedBackendResponse{}, err
	}
	return *result, nil
}

// ListAll retrieves all network backends by fetching all pages
func (s *networkBackendService) ListAll(ctx context.Context, lbID string) ([]NetworkBackendResponse, error) {
	var allBackends []NetworkBackendResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		pageOptions := ListNetworkLoadBalancerRequest{
			Offset: &currentOffset,
			Limit:  &currentLimit,
		}

		resp, err := s.List(ctx, lbID, pageOptions)
		if err != nil {
			return nil, err
		}

		allBackends = append(allBackends, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allBackends, nil
}

// Update updates a network backend's properties and returns the backend ID
func (s *networkBackendService) Update(ctx context.Context, lbID, backendID string, req UpdateNetworkBackendRequest) (string, error) {
	path := urlNetworkLoadBalancer(&lbID, backends, backendID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return "", err
	}

	var resp NetworkGenericCreationResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}
