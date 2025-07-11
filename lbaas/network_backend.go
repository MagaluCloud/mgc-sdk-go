package lbaas

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const backends = "backends"

type (
	// NetworkBackendInstanceRequest represents an instance-based backend target
	NetworkBackendInstanceRequest struct {
		// NicID is the network interface ID of the instance
		NicID string `json:"nic_id"`
		// Port is the port number for the backend target
		Port int `json:"port"`
	}

	// NetworkBackendRawTargetRequest represents a raw IP/port backend target
	NetworkBackendRawTargetRequest struct {
		// IPAddress is the IP address of the backend target
		IPAddress string `json:"ip_address"`
		// Port is the port number for the backend target
		Port int `json:"port"`
	}

	// TargetsRawOrInstancesRequest represents backend targets that can be either instances or raw IPs
	TargetsRawOrInstancesRequest struct {
		// TargetsInstances contains instance-based backend targets
		TargetsInstances []NetworkBackendInstanceRequest `json:"-"`
		// TargetsRaw contains raw IP/port backend targets
		TargetsRaw []NetworkBackendRawTargetRequest `json:"-"`
	}

	// NetworkBackendInstanceUpdateRequest represents an instance-based backend target for updates
	NetworkBackendInstanceUpdateRequest struct {
		// NicID is the network interface ID of the instance
		NicID string `json:"nic_id"`
		// Port is the port number for the backend target
		Port int `json:"port"`
	}

	// NetworkBackendRawTargetUpdateRequest represents a raw IP/port backend target for updates
	NetworkBackendRawTargetUpdateRequest struct {
		// IPAddress is the IP address of the backend target
		IPAddress string `json:"ip_address"`
		// Port is the port number for the backend target
		Port int `json:"port"`
	}

	// TargetsRawOrInstancesUpdateRequest represents backend targets for updates
	TargetsRawOrInstancesUpdateRequest struct {
		// TargetsInstances contains instance-based backend targets
		TargetsInstances []NetworkBackendInstanceUpdateRequest `json:"-"`
		// TargetsRaw contains raw IP/port backend targets
		TargetsRaw []NetworkBackendRawTargetUpdateRequest `json:"-"`
	}

	// CreateNetworkBackendRequest represents the request payload for creating a network backend
	CreateNetworkBackendRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// Name is the name of the backend
		Name string `json:"name"`
		// Description is the description of the backend (optional)
		Description *string `json:"description,omitempty"`
		// BalanceAlgorithm is the load balancing algorithm
		BalanceAlgorithm BackendBalanceAlgorithm `json:"balance_algorithm"`
		// TargetsType is the type of backend targets
		TargetsType BackendType `json:"targets_type"`
		// Targets contains the backend targets (optional)
		Targets *TargetsRawOrInstancesRequest `json:"targets,omitempty"`
		// HealthCheckID is the ID of the health check (optional)
		HealthCheckID *string `json:"health_check_id,omitempty"`
	}

	// DeleteNetworkBackendRequest represents the request payload for deleting a network backend
	DeleteNetworkBackendRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// BackendID is the ID of the backend to delete
		BackendID string `json:"-"`
	}

	// GetNetworkBackendRequest represents the request payload for getting a network backend
	GetNetworkBackendRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// BackendID is the ID of the backend to retrieve
		BackendID string `json:"-"`
	}

	// ListNetworkBackendRequest represents the request payload for listing network backends
	ListNetworkBackendRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
	}

	// UpdateNetworkBackendRequest represents the request payload for updating a network backend
	UpdateNetworkBackendRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// BackendID is the ID of the backend to update
		BackendID string `json:"-"`
		// Name is the new name of the backend (optional)
		Name *string `json:"name,omitempty"`
		// Description is the new description of the backend (optional)
		Description *string `json:"description,omitempty"`
		// BalanceAlgorithm is the new load balancing algorithm (optional)
		BalanceAlgorithm *BackendBalanceAlgorithm `json:"balance_algorithm,omitempty"`
		// TargetsType is the new type of backend targets (optional)
		TargetsType *BackendType `json:"targets_type,omitempty"`
		// Targets contains the new backend targets (optional)
		Targets *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
		// TargetsInstances contains the new instance-based targets (optional)
		TargetsInstances *[]NetworkBackendInstanceRequest `json:"targets_instances,omitempty"`
		// TargetsRaw contains the new raw IP/port targets (optional)
		TargetsRaw *[]NetworkBackendRawTargetRequest `json:"targets_raw,omitempty"`
		// HealthCheckID is the new health check ID (optional)
		HealthCheckID *string `json:"health_check_id,omitempty"`
	}

	// NetworkBackendInstanceResponse represents an instance-based backend target response
	NetworkBackendInstanceResponse struct {
		// ID is the unique identifier of the backend target
		ID string `json:"id"`
		// IPAddress is the IP address of the instance (optional)
		IPAddress *string `json:"ip_address,omitempty"`
		// NicID is the network interface ID of the instance
		NicID string `json:"nic_id,omitempty"`
		// Port is the port number for the backend target
		Port int `json:"port"`
		// CreatedAt is the creation timestamp
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the last update timestamp
		UpdatedAt string `json:"updated_at"`
	}

	// NetworkBackendRawTargetResponse represents a raw IP/port backend target response
	NetworkBackendRawTargetResponse struct {
		// ID is the unique identifier of the backend target
		ID string `json:"id"`
		// IPAddress is the IP address of the backend target (optional)
		IPAddress *string `json:"ip_address,omitempty"`
		// Port is the port number for the backend target
		Port int `json:"port"`
		// CreatedAt is the creation timestamp
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the last update timestamp
		UpdatedAt string `json:"updated_at"`
	}

	// NetworkBackendResponse represents a network backend response
	NetworkBackendResponse struct {
		// ID is the unique identifier of the backend
		ID string `json:"id"`
		// HealthCheckID is the ID of the associated health check (optional)
		HealthCheckID *string `json:"health_check_id,omitempty"`
		// Name is the name of the backend
		Name string `json:"name"`
		// Description is the description of the backend (optional)
		Description *string `json:"description,omitempty"`
		// BalanceAlgorithm is the load balancing algorithm
		BalanceAlgorithm BackendBalanceAlgorithm `json:"balance_algorithm"`
		// TargetsType is the type of backend targets
		TargetsType BackendType `json:"targets_type"`
		// Targets contains the backend targets
		Targets interface{} `json:"targets"`
		// CreatedAt is the creation timestamp
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the last update timestamp
		UpdatedAt string `json:"updated_at"`
	}

	// NetworkPaginatedBackendResponse represents a paginated backend response
	NetworkPaginatedBackendResponse struct {
		// Meta contains pagination metadata
		Meta interface{} `json:"meta"`
		// Results contains the list of backends
		Results []NetworkBackendResponse `json:"results"`
	}

	// NetworkBackendService provides methods for managing network backends
	NetworkBackendService interface {
		// Create creates a new network backend
		Create(ctx context.Context, req CreateNetworkBackendRequest) (string, error)
		// Delete removes a network backend
		Delete(ctx context.Context, req DeleteNetworkBackendRequest) error
		// Get retrieves detailed information about a specific backend
		Get(ctx context.Context, req GetNetworkBackendRequest) (*NetworkBackendResponse, error)
		// List returns a list of network backends
		List(ctx context.Context, req ListNetworkBackendRequest) ([]NetworkBackendResponse, error)
		// Update updates a network backend's properties
		Update(ctx context.Context, req UpdateNetworkBackendRequest) error
		// Targets returns a service for managing backend targets
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
