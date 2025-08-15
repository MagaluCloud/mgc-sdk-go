package lbaas

import (
	"context"
	"net/http"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const backends = "backends"

type (
	// NetworkBackendInstanceRequest represents an instance-based backend target
	NetworkBackendInstanceRequest struct {
		NicID string `json:"nic_id"`
		Port  int    `json:"port"`
	}

	// NetworkBackendRawTargetRequest represents a raw IP/port backend target
	NetworkBackendRawTargetRequest struct {
		IPAddress string `json:"ip_address"`
		Port      int    `json:"port"`
	}

	// TargetsRawOrInstancesRequest represents backend targets that can be either instances or raw IPs
	TargetsRawOrInstancesRequest struct {
		TargetsInstances []NetworkBackendInstanceRequest  `json:"-"`
		TargetsRaw       []NetworkBackendRawTargetRequest `json:"-"`
	}

	// NetworkBackendInstanceUpdateRequest represents an instance-based backend target for updates
	NetworkBackendInstanceUpdateRequest struct {
		NicID string `json:"nic_id"`
		Port  int    `json:"port"`
	}

	// NetworkBackendRawTargetUpdateRequest represents a raw IP/port backend target for updates
	NetworkBackendRawTargetUpdateRequest struct {
		IPAddress string `json:"ip_address"`
		Port      int    `json:"port"`
	}

	// TargetsRawOrInstancesUpdateRequest represents backend targets for updates
	TargetsRawOrInstancesUpdateRequest struct {
		TargetsInstances []NetworkBackendInstanceUpdateRequest  `json:"-"`
		TargetsRaw       []NetworkBackendRawTargetUpdateRequest `json:"-"`
	}

	// CreateNetworkBackendRequest represents the request payload for creating a network backend
	CreateNetworkBackendRequest struct {
		LoadBalancerID   string                        `json:"-"`
		Name             string                        `json:"name"`
		Description      *string                       `json:"description,omitempty"`
		BalanceAlgorithm BackendBalanceAlgorithm       `json:"balance_algorithm"`
		TargetsType      BackendType                   `json:"targets_type"`
		Targets          *TargetsRawOrInstancesRequest `json:"targets,omitempty"`
		HealthCheckID    *string                       `json:"health_check_id,omitempty"`
	}

	// DeleteNetworkBackendRequest represents the request payload for deleting a network backend
	DeleteNetworkBackendRequest struct {
		LoadBalancerID string `json:"-"`
		BackendID      string `json:"-"`
	}

	// GetNetworkBackendRequest represents the request payload for getting a network backend
	GetNetworkBackendRequest struct {
		LoadBalancerID string `json:"-"`
		BackendID      string `json:"-"`
	}

	// ListNetworkBackendRequest represents the request payload for listing network backends
	ListNetworkBackendRequest struct {
		LoadBalancerID string `json:"-"`
	}

	// UpdateNetworkBackendRequest represents the request payload for updating a network backend
	UpdateNetworkBackendRequest struct {
		LoadBalancerID   string                              `json:"-"`
		BackendID        string                              `json:"-"`
		Name             *string                             `json:"name,omitempty"`
		Description      *string                             `json:"description,omitempty"`
		BalanceAlgorithm *BackendBalanceAlgorithm            `json:"balance_algorithm,omitempty"`
		TargetsType      *BackendType                        `json:"targets_type,omitempty"`
		Targets          *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
		TargetsInstances *[]NetworkBackendInstanceRequest    `json:"targets_instances,omitempty"`
		TargetsRaw       *[]NetworkBackendRawTargetRequest   `json:"targets_raw,omitempty"`
		HealthCheckID    *string                             `json:"health_check_id,omitempty"`
	}

	// NetworkBackendInstanceResponse represents an instance-based backend target response
	NetworkBackendInstanceResponse struct {
		ID        string    `json:"id"`
		IPAddress *string   `json:"ip_address,omitempty"`
		NicID     string    `json:"nic_id,omitempty"`
		Port      int       `json:"port"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// NetworkBackendRawTargetResponse represents a raw IP/port backend target response
	NetworkBackendRawTargetResponse struct {
		ID        string    `json:"id"`
		IPAddress *string   `json:"ip_address,omitempty"`
		Port      int       `json:"port"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// NetworkBackendResponse represents a network backend response
	NetworkBackendResponse struct {
		ID               string                         `json:"id"`
		HealthCheckID    *string                        `json:"health_check_id,omitempty"`
		Name             string                         `json:"name"`
		Description      *string                        `json:"description,omitempty"`
		BalanceAlgorithm BackendBalanceAlgorithm        `json:"balance_algorithm"`
		TargetsType      BackendType                    `json:"targets_type"`
		Targets          []NetworkBackendTargetResponse `json:"targets"`
		CreatedAt        time.Time                      `json:"created_at"`
		UpdatedAt        time.Time                      `json:"updated_at"`
	}

	NetworkBackendTargetResponse struct {
		ID        string    `json:"id"`
		IPAddress string    `json:"ip_address"`
		Port      int       `json:"port"`
		NicID     *string   `json:"nic_id,omitempty"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// NetworkPaginatedBackendResponse represents a paginated backend response
	NetworkPaginatedBackendResponse struct {
		Meta    any                      `json:"meta"`
		Results []NetworkBackendResponse `json:"results"`
	}

	// NetworkBackendService provides methods for managing network backends
	NetworkBackendService interface {
		Create(ctx context.Context, req CreateNetworkBackendRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkBackendRequest) error
		Get(ctx context.Context, req GetNetworkBackendRequest) (*NetworkBackendResponse, error)
		List(ctx context.Context, req ListNetworkBackendRequest) ([]NetworkBackendResponse, error)
		Update(ctx context.Context, req UpdateNetworkBackendRequest) error
		Targets() *networkBackendTargetService
	}

	// networkBackendService implements the NetworkBackendService interface
	networkBackendService struct {
		client *LbaasClient
	}
)

// Targets returns a service for managing backend targets
func (s *networkBackendService) Targets() *networkBackendTargetService {
	return &networkBackendTargetService{client: s.client}
}

// Create creates a new network backend
func (s *networkBackendService) Create(ctx context.Context, req CreateNetworkBackendRequest) (string, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, backends)

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
func (s *networkBackendService) Delete(ctx context.Context, req DeleteNetworkBackendRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, backends, req.BackendID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific backend
func (s *networkBackendService) Get(ctx context.Context, req GetNetworkBackendRequest) (*NetworkBackendResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, backends, req.BackendID)

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

// List returns a list of network backends
func (s *networkBackendService) List(ctx context.Context, req ListNetworkBackendRequest) ([]NetworkBackendResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, backends)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkPaginatedBackendResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update updates a network backend's properties
func (s *networkBackendService) Update(ctx context.Context, req UpdateNetworkBackendRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, backends, req.BackendID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
