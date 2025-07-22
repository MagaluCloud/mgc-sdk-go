package lbaas

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const tls_certificates = "tls-certificates"

type (
	// CreateNetworkCertificateRequest represents the request payload for creating a network TLS certificate
	CreateNetworkCertificateRequest struct {
		LoadBalancerID string  `json:"-"`
		Name           string  `json:"name"`
		Description    *string `json:"description,omitempty"`
		Certificate    string  `json:"certificate"`
		PrivateKey     string  `json:"private_key"`
	}

	// DeleteNetworkCertificateRequest represents the request payload for deleting a network TLS certificate
	DeleteNetworkCertificateRequest struct {
		LoadBalancerID   string `json:"-"`
		TLSCertificateID string `json:"-"`
	}

	// GetNetworkCertificateRequest represents the request payload for getting a network TLS certificate
	GetNetworkCertificateRequest struct {
		LoadBalancerID   string `json:"-"`
		TLSCertificateID string `json:"-"`
	}

	// ListNetworkCertificateRequest represents the request payload for listing network TLS certificates
	ListNetworkCertificateRequest struct {
		LoadBalancerID string  `json:"-"`
		Offset         *int    `json:"-"`
		Limit          *int    `json:"-"`
		Sort           *string `json:"-"`
	}

	// UpdateNetworkCertificateRequest represents the request payload for updating a network TLS certificate
	UpdateNetworkCertificateRequest struct {
		LoadBalancerID   string `json:"-"`
		TLSCertificateID string `json:"-"`
		Certificate      string `json:"certificate"`
		PrivateKey       string `json:"private_key"`
	}

	// NetworkTLSCertificateResponse represents a network TLS certificate response
	NetworkTLSCertificateResponse struct {
		ID             string  `json:"id"`
		Name           string  `json:"name"`
		ExpirationDate *string `json:"expiration_date,omitempty"`
		Description    *string `json:"description,omitempty"`
		CreatedAt      string  `json:"created_at"`
		UpdatedAt      string  `json:"updated_at"`
	}

	// NetworkPaginatedTLSCertificateResponse represents a paginated TLS certificate response
	NetworkPaginatedTLSCertificateResponse struct {
		Meta    interface{}                     `json:"meta"`
		Results []NetworkTLSCertificateResponse `json:"results"`
	}

	// NetworkCertificateService provides methods for managing network TLS certificates
	NetworkCertificateService interface {
		Create(ctx context.Context, req CreateNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error)
		Delete(ctx context.Context, req DeleteNetworkCertificateRequest) error
		Get(ctx context.Context, req GetNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error)
		List(ctx context.Context, req ListNetworkCertificateRequest) ([]NetworkTLSCertificateResponse, error)
		Update(ctx context.Context, req UpdateNetworkCertificateRequest) error
	}

	// networkCertificateService implements the NetworkCertificateService interface
	networkCertificateService struct {
		client *LbaasClient
	}
)

// Create creates a new network TLS certificate
func (s *networkCertificateService) Create(ctx context.Context, req CreateNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, tls_certificates)

	// Validate if certificate and private key are base64 encoded
	if _, err := base64.StdEncoding.DecodeString(req.Certificate); err != nil {
		return nil, errors.New("certificate is not base64 encoded")
	}
	if _, err := base64.StdEncoding.DecodeString(req.PrivateKey); err != nil {
		return nil, errors.New("private key is not base64 encoded")
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}

	var resp NetworkTLSCertificateResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete removes a network TLS certificate
func (s *networkCertificateService) Delete(ctx context.Context, req DeleteNetworkCertificateRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, tls_certificates, req.TLSCertificateID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific TLS certificate
func (s *networkCertificateService) Get(ctx context.Context, req GetNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, tls_certificates, req.TLSCertificateID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkTLSCertificateResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List returns a list of network TLS certificates with optional filtering and pagination
func (s *networkCertificateService) List(ctx context.Context, req ListNetworkCertificateRequest) ([]NetworkTLSCertificateResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, tls_certificates)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", req.Offset)
	query.AddReflect("_limit", req.Limit)
	query.Add("_sort", req.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedTLSCertificateResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update updates a network TLS certificate's properties
func (s *networkCertificateService) Update(ctx context.Context, req UpdateNetworkCertificateRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, tls_certificates, req.TLSCertificateID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
