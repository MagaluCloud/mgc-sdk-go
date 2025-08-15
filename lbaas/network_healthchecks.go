package lbaas

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const health_checks = "health-checks"

type (
	// CreateNetworkHealthCheckRequest represents the request payload for creating a network health check
	CreateNetworkHealthCheckRequest struct {
		Name                    string              `json:"name"`
		Description             *string             `json:"description,omitempty"`
		Protocol                HealthCheckProtocol `json:"protocol"`
		Path                    *string             `json:"path,omitempty"`
		Port                    int                 `json:"port"`
		HealthyStatusCode       *int                `json:"healthy_status_code,omitempty"`
		IntervalSeconds         *int                `json:"interval_seconds,omitempty"`
		TimeoutSeconds          *int                `json:"timeout_seconds,omitempty"`
		InitialDelaySeconds     *int                `json:"initial_delay_seconds,omitempty"`
		HealthyThresholdCount   *int                `json:"healthy_threshold_count,omitempty"`
		UnhealthyThresholdCount *int                `json:"unhealthy_threshold_count,omitempty"`
	}

	// UpdateNetworkHealthCheckRequest represents the request payload for updating a network health check
	UpdateNetworkHealthCheckRequest struct {
		Protocol                HealthCheckProtocol `json:"protocol"`
		Path                    *string             `json:"path,omitempty"`
		Port                    int                 `json:"port"`
		HealthyStatusCode       *int                `json:"healthy_status_code,omitempty"`
		IntervalSeconds         *int                `json:"interval_seconds,omitempty"`
		TimeoutSeconds          *int                `json:"timeout_seconds,omitempty"`
		InitialDelaySeconds     *int                `json:"initial_delay_seconds,omitempty"`
		HealthyThresholdCount   *int                `json:"healthy_threshold_count,omitempty"`
		UnhealthyThresholdCount *int                `json:"unhealthy_threshold_count,omitempty"`
	}

	// NetworkHealthCheckResponse represents a network health check response
	NetworkHealthCheckResponse struct {
		ID                      string              `json:"id"`
		Name                    string              `json:"name"`
		Description             *string             `json:"description,omitempty"`
		Protocol                HealthCheckProtocol `json:"protocol"`
		Path                    *string             `json:"path,omitempty"`
		Port                    int                 `json:"port"`
		HealthyStatusCode       int                 `json:"healthy_status_code"`
		IntervalSeconds         int                 `json:"interval_seconds"`
		TimeoutSeconds          int                 `json:"timeout_seconds"`
		InitialDelaySeconds     int                 `json:"initial_delay_seconds"`
		HealthyThresholdCount   int                 `json:"healthy_threshold_count"`
		UnhealthyThresholdCount int                 `json:"unhealthy_threshold_count"`
		CreatedAt               string              `json:"created_at"`
		UpdatedAt               string              `json:"updated_at"`
	}

	// NetworkPaginatedHealthCheckResponse represents a paginated health check response
	NetworkPaginatedHealthCheckResponse struct {
		Meta    PaginationMeta               `json:"meta"`
		Results []NetworkHealthCheckResponse `json:"results"`
	}

	// NetworkHealthCheckService provides methods for managing network health checks
	NetworkHealthCheckService interface {
		Create(ctx context.Context, lbID string, req CreateNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error)
		Delete(ctx context.Context, lbID, healthCheckID string) error
		Get(ctx context.Context, lbID, healthCheckID string) (*NetworkHealthCheckResponse, error)
		List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) ([]NetworkHealthCheckResponse, error)
		Update(ctx context.Context, lbID, healthCheckID string, req UpdateNetworkHealthCheckRequest) error
	}

	// networkHealthCheckService implements the NetworkHealthCheckService interface
	networkHealthCheckService struct {
		client *LbaasClient
	}
)

// Create creates a new network health check
func (s *networkHealthCheckService) Create(ctx context.Context, lbID string, req CreateNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, health_checks)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}

	var resp NetworkHealthCheckResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete removes a network health check
func (s *networkHealthCheckService) Delete(ctx context.Context, lbID, healthCheckID string) error {
	path := urlNetworkLoadBalancer(&lbID, health_checks, healthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific health check
func (s *networkHealthCheckService) Get(ctx context.Context, lbID, healthCheckID string) (*NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, health_checks, healthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkHealthCheckResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List returns a list of network health checks with optional filtering and pagination
func (s *networkHealthCheckService) List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) ([]NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, health_checks)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", options.Offset)
	query.AddReflect("_limit", options.Limit)
	query.Add("_sort", options.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedHealthCheckResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update updates a network health check's properties
func (s *networkHealthCheckService) Update(ctx context.Context, lbID, healthCheckID string, req UpdateNetworkHealthCheckRequest) error {
	path := urlNetworkLoadBalancer(&lbID, health_checks, healthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
