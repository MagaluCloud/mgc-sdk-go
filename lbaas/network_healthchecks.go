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
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// Name is the name of the health check
		Name string `json:"name"`
		// Description is the description of the health check (optional)
		Description *string `json:"description,omitempty"`
		// Protocol is the protocol for the health check
		Protocol HealthCheckProtocol `json:"protocol"`
		// Path is the path for HTTP health checks (optional)
		Path *string `json:"path,omitempty"`
		// Port is the port for the health check
		Port int `json:"port"`
		// HealthyStatusCode is the expected status code for healthy responses (optional)
		HealthyStatusCode *int `json:"healthy_status_code,omitempty"`
		// IntervalSeconds is the interval between health checks in seconds (optional)
		IntervalSeconds *int `json:"interval_seconds,omitempty"`
		// TimeoutSeconds is the timeout for health checks in seconds (optional)
		TimeoutSeconds *int `json:"timeout_seconds,omitempty"`
		// InitialDelaySeconds is the initial delay before starting health checks (optional)
		InitialDelaySeconds *int `json:"initial_delay_seconds,omitempty"`
		// HealthyThresholdCount is the number of consecutive successful checks required (optional)
		HealthyThresholdCount *int `json:"healthy_threshold_count,omitempty"`
		// UnhealthyThresholdCount is the number of consecutive failed checks required (optional)
		UnhealthyThresholdCount *int `json:"unhealthy_threshold_count,omitempty"`
	}

	// DeleteNetworkHealthCheckRequest represents the request payload for deleting a network health check
	DeleteNetworkHealthCheckRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// HealthCheckID is the ID of the health check to delete
		HealthCheckID string `json:"-"`
	}

	// GetNetworkHealthCheckRequest represents the request payload for getting a network health check
	GetNetworkHealthCheckRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// HealthCheckID is the ID of the health check to retrieve
		HealthCheckID string `json:"-"`
	}

	// ListNetworkHealthCheckRequest represents the request payload for listing network health checks
	ListNetworkHealthCheckRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// Offset is the number of health checks to skip
		Offset *int `json:"-"`
		// Limit is the maximum number of health checks to return
		Limit *int `json:"-"`
		// Sort is the field to sort by
		Sort *string `json:"-"`
	}

	// UpdateNetworkHealthCheckRequest represents the request payload for updating a network health check
	UpdateNetworkHealthCheckRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// HealthCheckID is the ID of the health check to update
		HealthCheckID string `json:"-"`
		// Protocol is the protocol for the health check
		Protocol HealthCheckProtocol `json:"protocol"`
		// Path is the path for HTTP health checks (optional)
		Path *string `json:"path,omitempty"`
		// Port is the port for the health check
		Port int `json:"port"`
		// HealthyStatusCode is the expected status code for healthy responses (optional)
		HealthyStatusCode *int `json:"healthy_status_code,omitempty"`
		// IntervalSeconds is the interval between health checks in seconds (optional)
		IntervalSeconds *int `json:"interval_seconds,omitempty"`
		// TimeoutSeconds is the timeout for health checks in seconds (optional)
		TimeoutSeconds *int `json:"timeout_seconds,omitempty"`
		// InitialDelaySeconds is the initial delay before starting health checks (optional)
		InitialDelaySeconds *int `json:"initial_delay_seconds,omitempty"`
		// HealthyThresholdCount is the number of consecutive successful checks required (optional)
		HealthyThresholdCount *int `json:"healthy_threshold_count,omitempty"`
		// UnhealthyThresholdCount is the number of consecutive failed checks required (optional)
		UnhealthyThresholdCount *int `json:"unhealthy_threshold_count,omitempty"`
	}

	// NetworkHealthCheckResponse represents a network health check response
	NetworkHealthCheckResponse struct {
		// ID is the unique identifier of the health check
		ID string `json:"id"`
		// Name is the name of the health check
		Name string `json:"name"`
		// Description is the description of the health check (optional)
		Description *string `json:"description,omitempty"`
		// Protocol is the protocol for the health check
		Protocol HealthCheckProtocol `json:"protocol"`
		// Path is the path for HTTP health checks (optional)
		Path *string `json:"path,omitempty"`
		// Port is the port for the health check
		Port int `json:"port"`
		// HealthyStatusCode is the expected status code for healthy responses
		HealthyStatusCode int `json:"healthy_status_code"`
		// IntervalSeconds is the interval between health checks in seconds
		IntervalSeconds int `json:"interval_seconds"`
		// TimeoutSeconds is the timeout for health checks in seconds
		TimeoutSeconds int `json:"timeout_seconds"`
		// InitialDelaySeconds is the initial delay before starting health checks
		InitialDelaySeconds int `json:"initial_delay_seconds"`
		// HealthyThresholdCount is the number of consecutive successful checks required
		HealthyThresholdCount int `json:"healthy_threshold_count"`
		// UnhealthyThresholdCount is the number of consecutive failed checks required
		UnhealthyThresholdCount int `json:"unhealthy_threshold_count"`
		// CreatedAt is the creation timestamp
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the last update timestamp
		UpdatedAt string `json:"updated_at"`
	}

	// NetworkPaginatedHealthCheckResponse represents a paginated health check response
	NetworkPaginatedHealthCheckResponse struct {
		// Meta contains pagination metadata
		Meta interface{} `json:"meta"`
		// Results contains the list of health checks
		Results []NetworkHealthCheckResponse `json:"results"`
	}

	// NetworkHealthCheckService provides methods for managing network health checks
	NetworkHealthCheckService interface {
		// Create creates a new network health check
		Create(ctx context.Context, req CreateNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error)
		// Delete removes a network health check
		Delete(ctx context.Context, req DeleteNetworkHealthCheckRequest) error
		// Get retrieves detailed information about a specific health check
		Get(ctx context.Context, req GetNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error)
		// List returns a list of network health checks with optional filtering and pagination
		List(ctx context.Context, req ListNetworkHealthCheckRequest) ([]NetworkHealthCheckResponse, error)
		// Update updates a network health check's properties
		Update(ctx context.Context, req UpdateNetworkHealthCheckRequest) error
	}

	// networkHealthCheckService implements the NetworkHealthCheckService interface
	networkHealthCheckService struct {
		client *LbaasClient
	}
)

// Create creates a new network health check
func (s *networkHealthCheckService) Create(ctx context.Context, req CreateNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, health_checks)

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
func (s *networkHealthCheckService) Delete(ctx context.Context, req DeleteNetworkHealthCheckRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, health_checks, req.HealthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific health check
func (s *networkHealthCheckService) Get(ctx context.Context, req GetNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, health_checks, req.HealthCheckID)

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
func (s *networkHealthCheckService) List(ctx context.Context, req ListNetworkHealthCheckRequest) ([]NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, health_checks)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", req.Offset)
	query.AddReflect("_limit", req.Limit)
	query.Add("_sort", req.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedHealthCheckResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update updates a network health check's properties
func (s *networkHealthCheckService) Update(ctx context.Context, req UpdateNetworkHealthCheckRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, health_checks, req.HealthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
