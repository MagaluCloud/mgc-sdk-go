package lbaas

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const tls_certificates = "tls-certificates"

type (
	// CreateNetworkCertificateRequest represents the request payload for creating a network TLS certificate
	CreateNetworkCertificateRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Certificate string  `json:"certificate"`
		PrivateKey  string  `json:"private_key"`
	}

	// UpdateNetworkCertificateRequest represents the request payload for updating a network TLS certificate
	UpdateNetworkCertificateRequest struct {
		Certificate string `json:"certificate"`
		PrivateKey  string `json:"private_key"`
	}

	// NetworkTLSCertificateResponse represents a network TLS certificate response
	NetworkTLSCertificateResponse struct {
		ID             string     `json:"id"`
		Name           string     `json:"name"`
		ExpirationDate *time.Time `json:"expiration_date,omitempty"`
		Description    *string    `json:"description,omitempty"`
		CreatedAt      time.Time  `json:"created_at"`
		UpdatedAt      time.Time  `json:"updated_at"`
	}

	// NetworkPaginatedTLSCertificateResponse represents a paginated TLS certificate response
	NetworkPaginatedTLSCertificateResponse struct {
		Meta    PaginationMeta                  `json:"meta"`
		Results []NetworkTLSCertificateResponse `json:"results"`
	}

	// NetworkCertificateService provides methods for managing network TLS certificates
	NetworkCertificateService interface {
		Create(ctx context.Context, lbID string, req CreateNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error)
		Delete(ctx context.Context, lbID, certicateID string) error
		Get(ctx context.Context, lbID, certicateID string) (*NetworkTLSCertificateResponse, error)
		List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) (NetworkPaginatedTLSCertificateResponse, error)
		ListAll(ctx context.Context, lbID string) ([]NetworkTLSCertificateResponse, error)
		Update(ctx context.Context, lbID, certicateID string, req UpdateNetworkCertificateRequest) error
	}

	// networkCertificateService implements the NetworkCertificateService interface
	networkCertificateService struct {
		client *LbaasClient
	}
)

// Create creates a new network TLS certificate
func (s *networkCertificateService) Create(ctx context.Context, lbID string, req CreateNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, tls_certificates)

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
func (s *networkCertificateService) Delete(ctx context.Context, lbID, certicateID string) error {
	path := urlNetworkLoadBalancer(&lbID, tls_certificates, certicateID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a specific TLS certificate
func (s *networkCertificateService) Get(ctx context.Context, lbID, certicateID string) (*NetworkTLSCertificateResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, tls_certificates, certicateID)

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

// List returns a paginated list of network TLS certificates with optional filtering and pagination
func (s *networkCertificateService) List(ctx context.Context, lbID string, options ListNetworkLoadBalancerRequest) (NetworkPaginatedTLSCertificateResponse, error) {
	path := urlNetworkLoadBalancer(&lbID, tls_certificates)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return NetworkPaginatedTLSCertificateResponse{}, err
	}
	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", options.Offset)
	query.AddReflect("_limit", options.Limit)
	query.Add("_sort", options.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedTLSCertificateResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return NetworkPaginatedTLSCertificateResponse{}, err
	}
	return *result, nil
}

// ListAll retrieves all network TLS certificates by fetching all pages
func (s *networkCertificateService) ListAll(ctx context.Context, lbID string) ([]NetworkTLSCertificateResponse, error) {
	var allCertificates []NetworkTLSCertificateResponse
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

		allCertificates = append(allCertificates, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allCertificates, nil
}

// Update updates a network TLS certificate's properties
func (s *networkCertificateService) Update(ctx context.Context, lbID, certicateID string, req UpdateNetworkCertificateRequest) error {
	path := urlNetworkLoadBalancer(&lbID, tls_certificates, certicateID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
