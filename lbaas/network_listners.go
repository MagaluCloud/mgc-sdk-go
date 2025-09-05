package lbaas

import (
	"context"
	"net/http"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const listeners = "listeners"

type (
	// CreateNetworkListenerRequest represents the request payload for creating a network listener
	CreateNetworkListenerRequest struct {
		TLSCertificateID *string          `json:"tls_certificate_id,omitempty"`
		Name             string           `json:"name"`
		Description      *string          `json:"description,omitempty"`
		Protocol         ListenerProtocol `json:"protocol"`
		Port             int              `json:"port"`
	}

	// UpdateNetworkListenerRequest represents the request payload for updating a network listener
	UpdateNetworkListenerRequest struct {
		TLSCertificateID *string `json:"tls_certificate_id,omitempty"`
		Name             *string `json:"name,omitempty"`
	}

	// NetworkListenerResponse represents a network listener response
	NetworkListenerResponse struct {
		ID               string           `json:"id"`
		TLSCertificateID *string          `json:"tls_certificate_id,omitempty"`
		BackendID        string           `json:"backend_id"`
		Name             string           `json:"name"`
		Description      *string          `json:"description,omitempty"`
		Protocol         ListenerProtocol `json:"protocol"`
		Port             int              `json:"port"`
		CreatedAt        time.Time        `json:"created_at"`
		UpdatedAt        time.Time        `json:"updated_at"`
	}

	// NetworkPaginatedListenerResponse represents a paginated listener response
	NetworkPaginatedListenerResponse struct {
		Meta    PaginationMeta            `json:"meta"`
		Results []NetworkListenerResponse `json:"results"`
	}

	// NetworkListenerService provides methods for managing network listeners
	NetworkListenerService interface {
		Create(ctx context.Context, lbID, backendID string, req CreateNetworkListenerRequest) (*NetworkListenerResponse, error)
		Delete(ctx context.Context, lbID, listenerID string) error
		Get(ctx context.Context, lbID, listenerID string) (*NetworkListenerResponse, error)
		List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) ([]NetworkListenerResponse, error)
		Update(ctx context.Context, lbID, listenerID string, req UpdateNetworkListenerRequest) error
	}

	// networkListenerService implements the NetworkListenerService interface
	networkListenerService struct {
		client *LbaasClient
	}
)

// Create creates a new network listener
func (s *networkListenerService) Create(ctx context.Context, lbID, backendID string, req CreateNetworkListenerRequest) (*NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, listeners)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}

	// Add backend_id as required query parameter
	query := httpReq.URL.Query()
	query.Set("backend_id", backendID)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete removes a network listener
func (s *networkListenerService) Delete(ctx context.Context, lbID, listenerID string) error {
	path := urlNetworkLoadBalancer(&lbID, listeners, listenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific listener
func (s *networkListenerService) Get(ctx context.Context, lbID, listenerID string) (*NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, listeners, listenerID)

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
func (s *networkListenerService) List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) ([]NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, listeners)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", options.Offset)
	query.AddReflect("_limit", options.Limit)
	query.AddReflect("_sort", options.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update updates a network listener's properties
func (s *networkListenerService) Update(ctx context.Context, lbID, listenerID string, req UpdateNetworkListenerRequest) error {
	path := urlNetworkLoadBalancer(&lbID, listeners, listenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
