package lbaas

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const listeners = "listeners"

type (
	// CreateNetworkListenerRequest represents the request payload for creating a network listener
	CreateNetworkListenerRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// BackendID is the ID of the backend to associate with the listener
		BackendID string `json:"-"`
		// TLSCertificateID is the ID of the TLS certificate (optional)
		TLSCertificateID *string `json:"tls_certificate_id,omitempty"`
		// Name is the name of the listener
		Name string `json:"name"`
		// Description is the description of the listener (optional)
		Description *string `json:"description,omitempty"`
		// Protocol is the protocol for the listener
		Protocol ListenerProtocol `json:"protocol"`
		// Port is the port number for the listener
		Port int `json:"port"`
	}

	// DeleteNetworkListenerRequest represents the request payload for deleting a network listener
	DeleteNetworkListenerRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// ListenerID is the ID of the listener to delete
		ListenerID string `json:"-"`
	}

	// GetNetworkListenerRequest represents the request payload for getting a network listener
	GetNetworkListenerRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// ListenerID is the ID of the listener to retrieve
		ListenerID string `json:"-"`
	}

	// ListNetworkListenerRequest represents the request payload for listing network listeners
	ListNetworkListenerRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// Offset is the number of listeners to skip
		Offset *int `json:"-"`
		// Limit is the maximum number of listeners to return
		Limit *int `json:"-"`
		// Sort is the field to sort by
		Sort *string `json:"-"`
	}

	// UpdateNetworkListenerRequest represents the request payload for updating a network listener
	UpdateNetworkListenerRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// ListenerID is the ID of the listener to update
		ListenerID string `json:"-"`
		// TLSCertificateID is the new TLS certificate ID (optional)
		TLSCertificateID *string `json:"tls_certificate_id,omitempty"`
	}

	// NetworkListenerResponse represents a network listener response
	NetworkListenerResponse struct {
		// ID is the unique identifier of the listener
		ID string `json:"id"`
		// TLSCertificateID is the ID of the associated TLS certificate (optional)
		TLSCertificateID *string `json:"tls_certificate_id,omitempty"`
		// BackendID is the ID of the associated backend
		BackendID string `json:"backend_id"`
		// Name is the name of the listener
		Name string `json:"name"`
		// Description is the description of the listener (optional)
		Description *string `json:"description,omitempty"`
		// Protocol is the protocol for the listener
		Protocol ListenerProtocol `json:"protocol"`
		// Port is the port number for the listener
		Port int `json:"port"`
		// CreatedAt is the creation timestamp
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the last update timestamp
		UpdatedAt string `json:"updated_at"`
	}

	// NetworkPaginatedListenerResponse represents a paginated listener response
	NetworkPaginatedListenerResponse struct {
		// Meta contains pagination metadata
		Meta interface{} `json:"meta"`
		// Results contains the list of listeners
		Results []NetworkListenerResponse `json:"results"`
	}

	// NetworkListenerService provides methods for managing network listeners
	NetworkListenerService interface {
		// Create creates a new network listener
		Create(ctx context.Context, req CreateNetworkListenerRequest) (*NetworkListenerResponse, error)
		// Delete removes a network listener
		Delete(ctx context.Context, req DeleteNetworkListenerRequest) error
		// Get retrieves detailed information about a specific listener
		Get(ctx context.Context, req GetNetworkListenerRequest) (*NetworkListenerResponse, error)
		// List returns a list of network listeners with optional filtering and pagination
		List(ctx context.Context, req ListNetworkListenerRequest) ([]NetworkListenerResponse, error)
		// Update updates a network listener's properties
		Update(ctx context.Context, req UpdateNetworkListenerRequest) error
	}

	// networkListenerService implements the NetworkListenerService interface
	networkListenerService struct {
		client *LbaasClient
	}
)

// Create creates a new network listener
func (s *networkListenerService) Create(ctx context.Context, req CreateNetworkListenerRequest) (*NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}

	// Add backend_id as required query parameter
	query := httpReq.URL.Query()
	query.Set("backend_id", req.BackendID)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete removes a network listener
func (s *networkListenerService) Delete(ctx context.Context, req DeleteNetworkListenerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners, req.ListenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific listener
func (s *networkListenerService) Get(ctx context.Context, req GetNetworkListenerRequest) (*NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners, req.ListenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List returns a list of network listeners with optional filtering and pagination
func (s *networkListenerService) List(ctx context.Context, req ListNetworkListenerRequest) ([]NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", req.Offset)
	query.AddReflect("_limit", req.Limit)
	query.Add("_sort", req.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update updates a network listener's properties
func (s *networkListenerService) Update(ctx context.Context, req UpdateNetworkListenerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners, req.ListenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
