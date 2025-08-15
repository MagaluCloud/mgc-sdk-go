package lbaas

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

//
// Package overview
//
// The lbaas package provides a client for Magalu Cloud's Network Load Balancer service.
// This file contains request/response models aligned with the OpenAPI v0beta1 specification.
//

// NetworkListenerRequest represents a listener configuration for load balancer creation.
// Name, BackendName, Protocol, and Port are required.
// For TLS protocol, set TLSCertificateName to reference a certificate by name.
type (
	// NetworkListenerRequest represents a listener configuration for load balancer creation
	NetworkListenerRequest struct {
		BackendName string `json:"backend_name"`
		CreateNetworkListenerRequest
	}

	// CreateNetworkLoadBalancerRequest is the request to create a Network Load Balancer.
	// Name, Visibility, VPCID, Listeners, and Backends are required.
	// Cross-references: BackendName must match a backend Name, TLSCertificateName must match a certificate Name.
	CreateNetworkLoadBalancerRequest struct {
		Name            string                            `json:"name"`
		Description     *string                           `json:"description,omitempty"`
		Type            *string                           `json:"type,omitempty"`
		Visibility      LoadBalancerVisibility            `json:"visibility"`
		Listeners       []NetworkListenerRequest          `json:"listeners"`
		Backends        []CreateNetworkBackendRequest     `json:"backends"`
		HealthChecks    []CreateNetworkHealthCheckRequest `json:"health_checks,omitempty"`
		TLSCertificates []CreateNetworkCertificateRequest `json:"tls_certificates,omitempty"`
		ACLs            []CreateNetworkACLRequest         `json:"acls,omitempty"`
		VPCID           string                            `json:"vpc_id"`
		SubnetPoolID    *string                           `json:"subnet_pool_id,omitempty"`
		PublicIPID      *string                           `json:"public_ip_id,omitempty"`
	}

	// DeleteNetworkLoadBalancerRequest controls deletion behavior.
	// Set DeletePublicIP=true to also delete the associated public IP.
	DeleteNetworkLoadBalancerRequest struct {
		DeletePublicIP *bool `json:"-"`
	}

	// ListNetworkLoadBalancerRequest defines pagination and sorting for listing.
	// All fields are optional and map to query parameters.
	// Sort format: "field:direction" (e.g., "created_at:desc").
	ListNetworkLoadBalancerRequest struct {
		Offset *int    `json:"-"`
		Limit  *int    `json:"-"`
		Sort   *string `json:"-"`
	}

	// UpdateNetworkLoadBalancerRequest updates a Network Load Balancer.
	// Only Name and Description can be updated. All fields are optional.
	// Sub-resource updates require separate API calls.
	UpdateNetworkLoadBalancerRequest struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// NetworkGenericCreationResponse represents a generic creation/update response
	NetworkGenericCreationResponse struct {
		ID string `json:"id"`
	}

	// NetworkPublicIPResponse contains information about a public IP associated
	// with a Load Balancer.
	NetworkPublicIPResponse struct {
		ID         string  `json:"id"`
		IPAddress  *string `json:"ip_address,omitempty"`
		ExternalID string  `json:"external_id"`
	}

	// NetworkLoadBalancerResponse describes a Load Balancer and its
	// sub-resources as returned by the API.
	//
	// Fields include status and timestamps that can help track provisioning
	// progress. The optional `LastOperationStatus` may surface the last
	// lifecycle operation outcome.
	NetworkLoadBalancerResponse struct {
		ID                  string                          `json:"id"`
		Name                string                          `json:"name"`
		ProjectType         *string                         `json:"project_type,omitempty"`
		Description         *string                         `json:"description,omitempty"`
		Type                string                          `json:"type"`
		Visibility          LoadBalancerVisibility          `json:"visibility"`
		Status              string                          `json:"status"`
		Listeners           []NetworkListenerResponse       `json:"listeners"`
		Backends            []NetworkBackendResponse        `json:"backends"`
		HealthChecks        []NetworkHealthCheckResponse    `json:"health_checks"`
		PublicIPs           []NetworkPublicIPResponse       `json:"public_ips"`
		TLSCertificates     []NetworkTLSCertificateResponse `json:"tls_certificates"`
		ACLs                []NetworkAclResponse            `json:"acls"`
		IPAddress           *string                         `json:"ip_address,omitempty"`
		VPCID               string                          `json:"vpc_id"`
		SubnetPoolID        *string                         `json:"subnet_pool_id,omitempty"`
		CreatedAt           time.Time                       `json:"created_at"`
		UpdatedAt           time.Time                       `json:"updated_at"`
		LastOperationStatus *string                         `json:"last_operation_status,omitempty"`
	}

	// PaginationLinks provides navigation links for pagination
	PaginationLinks struct {
		Next     *string `json:"next,omitempty"`
		Previous *string `json:"previous,omitempty"`
		Self     string  `json:"self"`
	}

	// PaginationPage contains pagination metadata
	PaginationPage struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	// PaginationMeta combines links and page information
	PaginationMeta struct {
		Links PaginationLinks `json:"links"`
		Page  PaginationPage  `json:"page"`
	}

	// NetworkLBPaginatedResponse is the paginated response for listing LBs.
	NetworkLBPaginatedResponse struct {
		Meta    PaginationMeta                `json:"meta"`
		Results []NetworkLoadBalancerResponse `json:"results"`
	}

	// NetworkLoadBalancerService provides CRUD operations for Network Load Balancers.
	NetworkLoadBalancerService interface {
		Create(ctx context.Context, create CreateNetworkLoadBalancerRequest) (string, error)
		Delete(ctx context.Context, id string, options DeleteNetworkLoadBalancerRequest) error
		Get(ctx context.Context, id string) (NetworkLoadBalancerResponse, error)
		List(ctx context.Context, options ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error)
		Update(ctx context.Context, id string, loadBalancer UpdateNetworkLoadBalancerRequest) (string, error)
	}

	// networkLoadBalancerService implements the NetworkLoadBalancerService interface.
	// It is typically constructed by `(*LbaasClient).NetworkLoadBalancers()`.
	networkLoadBalancerService struct {
		client *LbaasClient
	}
)

// Create creates a new Network Load Balancer and returns its ID.
func (s *networkLoadBalancerService) Create(ctx context.Context, create CreateNetworkLoadBalancerRequest) (string, error) {
	path := urlNetworkLoadBalancer(nil)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, create)
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

// Delete removes a Network Load Balancer by ID.
// Set DeletePublicIP=true in options to also remove the public IP.
func (s *networkLoadBalancerService) Delete(ctx context.Context, id string, options DeleteNetworkLoadBalancerRequest) error {
	path := urlNetworkLoadBalancer(&id)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	if options.DeletePublicIP != nil {
		query := httpReq.URL.Query()
		query.Set("delete_public_ip", strconv.FormatBool(*options.DeletePublicIP))
		httpReq.URL.RawQuery = query.Encode()
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Get retrieves detailed information about a Load Balancer by ID.
func (s *networkLoadBalancerService) Get(ctx context.Context, id string) (NetworkLoadBalancerResponse, error) {
	path := urlNetworkLoadBalancer(&id)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return NetworkLoadBalancerResponse{}, err
	}

	var resp NetworkLoadBalancerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return NetworkLoadBalancerResponse{}, err
	}
	if result == nil {
		return NetworkLoadBalancerResponse{}, errors.New("load balancer found but response is nil")
	}

	return *result, nil
}

// List returns Network Load Balancers with optional pagination and sorting.
func (s *networkLoadBalancerService) List(ctx context.Context, options ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error) {
	path := urlNetworkLoadBalancer(nil)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := helpers.NewQueryParams(httpReq)
	query.AddReflect("_offset", options.Offset)
	query.AddReflect("_limit", options.Limit)
	query.Add("_sort", options.Sort)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkLBPaginatedResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Update modifies a Network Load Balancer's name and/or description.
// Returns the ID of the updated load balancer.
// Only name and description can be updated via this endpoint.
func (s *networkLoadBalancerService) Update(ctx context.Context, id string, loadBalancer UpdateNetworkLoadBalancerRequest) (string, error) {
	path := urlNetworkLoadBalancer(&id)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, loadBalancer)
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
