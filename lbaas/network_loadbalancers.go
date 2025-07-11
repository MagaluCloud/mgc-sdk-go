package lbaas

import (
	"context"
	"net/http"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// NetworkListenerRequest represents a listener configuration for load balancer creation
	NetworkListenerRequest struct {
		// TLSCertificateName is the name of the TLS certificate (optional)
		TLSCertificateName *string `json:"tls_certificate_name,omitempty"`
		// Name is the name of the listener
		Name string `json:"name"`
		// Description is the description of the listener (optional)
		Description *string `json:"description,omitempty"`
		// BackendName is the name of the associated backend
		BackendName string `json:"backend_name"`
		// Protocol is the protocol for the listener
		Protocol ListenerProtocol `json:"protocol"`
		// Port is the port number for the listener
		Port int `json:"port"`
	}

	// NetworkBackendRequest represents a backend configuration for load balancer creation
	NetworkBackendRequest struct {
		// HealthCheckName is the name of the health check (optional)
		HealthCheckName *string `json:"health_check_name,omitempty"`
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
	}

	// NetworkHealthCheckRequest represents a health check configuration for load balancer creation
	NetworkHealthCheckRequest struct {
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

	// NetworkTLSCertificateRequest represents a TLS certificate configuration for load balancer creation
	NetworkTLSCertificateRequest struct {
		// Name is the name of the TLS certificate
		Name string `json:"name"`
		// Description is the description of the TLS certificate (optional)
		Description *string `json:"description,omitempty"`
		// Certificate is the TLS certificate content
		Certificate string `json:"certificate"`
		// PrivateKey is the private key content
		PrivateKey string `json:"private_key"`
	}

	// NetworkAclRequest represents an ACL rule configuration for load balancer creation
	NetworkAclRequest struct {
		// Name is the name of the ACL rule (optional)
		Name *string `json:"name,omitempty"`
		// Ethertype is the ethernet type for the ACL rule
		Ethertype AclEtherType `json:"ethertype"`
		// Protocol is the protocol for the ACL rule
		Protocol AclProtocol `json:"protocol"`
		// RemoteIPPrefix is the remote IP prefix for the ACL rule
		RemoteIPPrefix string `json:"remote_ip_prefix"`
		// Action is the action to take for matching traffic
		Action AclActionType `json:"action"`
	}

	// CreateNetworkLoadBalancerRequest represents the request payload for creating a load balancer
	CreateNetworkLoadBalancerRequest struct {
		// Name is the name of the load balancer
		Name string `json:"name"`
		// Description is the description of the load balancer (optional)
		Description *string `json:"description,omitempty"`
		// Type is the type of the load balancer (optional)
		Type *string `json:"type,omitempty"`
		// Visibility is the visibility of the load balancer
		Visibility LoadBalancerVisibility `json:"visibility"`
		// Listeners contains the listener configurations
		Listeners []NetworkListenerRequest `json:"listeners"`
		// Backends contains the backend configurations
		Backends []NetworkBackendRequest `json:"backends"`
		// HealthChecks contains the health check configurations (optional)
		HealthChecks []NetworkHealthCheckRequest `json:"health_checks,omitempty"`
		// TLSCertificates contains the TLS certificate configurations (optional)
		TLSCertificates []NetworkTLSCertificateRequest `json:"tls_certificates,omitempty"`
		// ACLs contains the ACL rule configurations (optional)
		ACLs []NetworkAclRequest `json:"acls,omitempty"`
		// VPCID is the ID of the VPC
		VPCID string `json:"vpc_id"`
		// SubnetPoolID is the ID of the subnet pool (optional)
		SubnetPoolID *string `json:"subnet_pool_id,omitempty"`
		// PublicIPID is the ID of the public IP (optional)
		PublicIPID *string `json:"public_ip_id,omitempty"`
		// PanicThreshold is the panic threshold for the load balancer (optional)
		PanicThreshold *int `json:"panic_threshold,omitempty"`
	}

	// DeleteNetworkLoadBalancerRequest represents the request payload for deleting a load balancer
	DeleteNetworkLoadBalancerRequest struct {
		// LoadBalancerID is the ID of the load balancer to delete
		LoadBalancerID string `json:"-"`
		// DeletePublicIP indicates whether to delete the associated public IP
		DeletePublicIP *bool `json:"-"`
	}

	// GetNetworkLoadBalancerRequest represents the request payload for getting a load balancer
	GetNetworkLoadBalancerRequest struct {
		// LoadBalancerID is the ID of the load balancer to retrieve
		LoadBalancerID string `json:"-"`
	}

	// ListNetworkLoadBalancerRequest represents the request payload for listing load balancers
	ListNetworkLoadBalancerRequest struct {
		// Offset is the number of load balancers to skip
		Offset *int `json:"-"`
		// Limit is the maximum number of load balancers to return
		Limit *int `json:"-"`
		// Sort is the field to sort by
		Sort *string `json:"-"`
	}

	// NetworkBackendUpdateRequest represents a backend update configuration
	NetworkBackendUpdateRequest struct {
		// ID is the ID of the backend to update
		ID string `json:"id"`
		// HealthCheckID is the ID of the health check (optional)
		HealthCheckID *string `json:"health_check_id,omitempty"`
		// TargetsType is the type of backend targets
		TargetsType BackendType `json:"targets_type"`
		// Targets contains the backend targets (optional)
		Targets *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
	}

	// UpdateNetworkLoadBalancerRequest represents the request payload for updating a load balancer
	UpdateNetworkLoadBalancerRequest struct {
		// LoadBalancerID is the ID of the load balancer to update
		LoadBalancerID string `json:"-"`
		// Name is the new name of the load balancer (optional)
		Name *string `json:"name,omitempty"`
		// Description is the new description of the load balancer (optional)
		Description *string `json:"description,omitempty"`
		// Backends contains the backend update configurations (optional)
		Backends []NetworkBackendUpdateRequest `json:"backends,omitempty"`
		// HealthChecks contains the health check update configurations (optional)
		HealthChecks []NetworkHealthCheckUpdateRequest `json:"health_checks,omitempty"`
		// TLSCertificates contains the TLS certificate update configurations (optional)
		TLSCertificates []NetworkTLSCertificateUpdateRequest `json:"tls_certificates,omitempty"`
		// PanicThreshold is the new panic threshold (optional)
		PanicThreshold *int `json:"panic_threshold,omitempty"`
	}

	// NetworkHealthCheckUpdateRequest represents a health check update configuration
	NetworkHealthCheckUpdateRequest struct {
		// ID is the ID of the health check to update
		ID string `json:"id"`
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

	// NetworkTLSCertificateUpdateRequest represents a TLS certificate update configuration
	NetworkTLSCertificateUpdateRequest struct {
		// ID is the ID of the TLS certificate to update
		ID string `json:"id"`
		// Certificate is the new TLS certificate content
		Certificate string `json:"certificate"`
		// PrivateKey is the new private key content
		PrivateKey string `json:"private_key"`
	}

	// NetworkPublicIPResponse represents a public IP response
	NetworkPublicIPResponse struct {
		// ID is the unique identifier of the public IP
		ID string `json:"id"`
		// IPAddress is the IP address (optional)
		IPAddress *string `json:"ip_address,omitempty"`
		// ExternalID is the external identifier
		ExternalID string `json:"external_id"`
	}

	// NetworkAclResponse represents an ACL rule response
	NetworkAclResponse struct {
		// ID is the unique identifier of the ACL rule
		ID string `json:"id"`
		// Name is the name of the ACL rule (optional)
		Name *string `json:"name,omitempty"`
		// Ethertype is the ethernet type for the ACL rule
		Ethertype AclEtherType `json:"ethertype"`
		// Protocol is the protocol for the ACL rule
		Protocol AclProtocol `json:"protocol"`
		// RemoteIPPrefix is the remote IP prefix for the ACL rule
		RemoteIPPrefix string `json:"remote_ip_prefix"`
		// Action is the action taken for matching traffic
		Action string `json:"action"`
	}

	// NetworkLoadBalancerResponse represents a load balancer response
	NetworkLoadBalancerResponse struct {
		// ID is the unique identifier of the load balancer
		ID string `json:"id"`
		// Name is the name of the load balancer
		Name string `json:"name"`
		// ProjectType is the project type (optional)
		ProjectType *string `json:"project_type,omitempty"`
		// Description is the description of the load balancer (optional)
		Description *string `json:"description,omitempty"`
		// Type is the type of the load balancer
		Type string `json:"type"`
		// Visibility is the visibility of the load balancer
		Visibility LoadBalancerVisibility `json:"visibility"`
		// Status is the status of the load balancer
		Status string `json:"status"`
		// Listeners contains the listener responses
		Listeners []NetworkListenerResponse `json:"listeners"`
		// Backends contains the backend responses
		Backends []NetworkBackendResponse `json:"backends"`
		// HealthChecks contains the health check responses
		HealthChecks []NetworkHealthCheckResponse `json:"health_checks"`
		// PublicIPs contains the public IP responses
		PublicIPs []NetworkPublicIPResponse `json:"public_ips"`
		// TLSCertificates contains the TLS certificate responses
		TLSCertificates []NetworkTLSCertificateResponse `json:"tls_certificates"`
		// ACLs contains the ACL rule responses
		ACLs []NetworkAclResponse `json:"acls"`
		// IPAddress is the IP address of the load balancer (optional)
		IPAddress *string `json:"ip_address,omitempty"`
		// Port is the port of the load balancer (optional)
		Port *string `json:"port,omitempty"`
		// VPCID is the ID of the VPC
		VPCID string `json:"vpc_id"`
		// SubnetPoolID is the ID of the subnet pool (optional)
		SubnetPoolID *string `json:"subnet_pool_id,omitempty"`
		// CreatedAt is the creation timestamp
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the last update timestamp
		UpdatedAt string `json:"updated_at"`
		// LastOperationStatus is the status of the last operation (optional)
		LastOperationStatus *string `json:"last_operation_status,omitempty"`
	}

	// NetworkLBPaginatedResponse represents a paginated load balancer response
	NetworkLBPaginatedResponse struct {
		// Results contains the list of load balancers
		Results []NetworkLoadBalancerResponse `json:"results"`
	}

	// NetworkLoadBalancerService provides methods for managing network load balancers
	NetworkLoadBalancerService interface {
		// Create creates a new network load balancer
		Create(ctx context.Context, req CreateNetworkLoadBalancerRequest) (string, error)
		// Delete removes a network load balancer
		Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error
		// Get retrieves detailed information about a specific load balancer
		Get(ctx context.Context, req GetNetworkLoadBalancerRequest) (*NetworkLoadBalancerResponse, error)
		// List returns a list of network load balancers with optional filtering and pagination
		List(ctx context.Context, req ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error)
		// Update updates a network load balancer's properties
		Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error
	}

	// networkLoadBalancerService implements the NetworkLoadBalancerService interface
	networkLoadBalancerService struct {
		client *LbaasClient
	}
)

// Create creates a new network load balancer
func (s *networkLoadBalancerService) Create(ctx context.Context, req CreateNetworkLoadBalancerRequest) (string, error) {
	path := urlNetworkLoadBalancer(nil)

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

// Delete removes a network load balancer
func (s *networkLoadBalancerService) Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	if req.DeletePublicIP != nil && *req.DeletePublicIP {
		query := httpReq.URL.Query()
		query.Set("delete_public_ip", strconv.FormatBool(*req.DeletePublicIP))
		httpReq.URL.RawQuery = query.Encode()
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific load balancer
func (s *networkLoadBalancerService) Get(ctx context.Context, req GetNetworkLoadBalancerRequest) (*NetworkLoadBalancerResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkLoadBalancerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List returns a list of network load balancers with optional filtering and pagination
func (s *networkLoadBalancerService) List(ctx context.Context, req ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error) {
	path := urlNetworkLoadBalancer(nil)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", req.Offset)
	query.AddReflect("_limit", req.Limit)
	query.Add("_sort", req.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkLBPaginatedResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update updates a network load balancer's properties
func (s *networkLoadBalancerService) Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
