package lbaas

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

//
// Package overview
//
// The lbaas package exposes a typed, ergonomic client for Magalu Cloud's
// Load Balancer as a Service (LBaaS). This file focuses on the Network
// Load Balancer (Layer 4 - TCP/TLS) service and its request/response
// models.
//
// Key concepts:
//
//   - Load Balancer: The top-level resource (internal or external visibility).
//   - Listener: Entry point for traffic (port + protocol). For TLS listeners,
//     you can reference a TLS certificate uploaded with the LB.
//   - Backend: A set of targets (instances or raw IPs) and a balance algorithm.
//   - Health Check: Periodic checks (TCP/HTTP) that keep unhealthy targets
//     out of rotation.
//   - ACL: IP-based allow/deny rules enforced at the LB edge.
//   - TLS Certificates: PEM-encoded certificate/key pairs used by TLS listeners.
//
// Typical workflow:
//
//   1. Create the LB with listeners, backends, (optional) health checks,
//      ACLs and TLS certificates.
//   2. Associate a public or private address (according to `visibility`).
//   3. Point your DNS to the LB’s address.
//   4. Keep the LB and its components up to date via `Update`.
//   5. Delete the LB when you no longer need it (optionally removing
//      the public IP together).
//
// API endpoints (as of v0beta1):
//
//   - POST   /v0beta1/network-load-balancers
//   - GET    /v0beta1/network-load-balancers/{id}
//   - GET    /v0beta1/network-load-balancers?_offset=…&_limit=…&_sort=…
//   - PUT    /v0beta1/network-load-balancers/{id}
//   - DELETE /v0beta1/network-load-balancers/{id}?delete_public_ip=…
//

// NetworkListenerRequest represents a Listener to be created with a Load Balancer.
//
// Required fields:
//   - Name
//   - BackendName
//   - Protocol (TCP or TLS; see `ListenerProtocol`)
//   - Port
//
// TLS-specific notes:
//   - For `Protocol == TLS`, set `TLSCertificateName` to reference a certificate
//     provided in `CreateNetworkLoadBalancerRequest.TLSCertificates` (by `Name`).
//     The certificate/key must be PEM-encoded.
//
// Routing notes:
//   - `BackendName` must match the `Name` of one of the backends supplied in
//     `CreateNetworkLoadBalancerRequest.Backends`.
type (
	// NetworkListenerRequest represents a listener configuration for load balancer creation
	NetworkListenerRequest struct {
		TLSCertificateName *string          `json:"tls_certificate_name,omitempty"`
		Name               string           `json:"name"`
		Description        *string          `json:"description,omitempty"`
		BackendName        string           `json:"backend_name"`
		Protocol           ListenerProtocol `json:"protocol"`
		Port               int              `json:"port"`
	}

	// NetworkBackendRequest represents a Backend to be created with a Load Balancer.
	//
	// Fields:
	//   - BalanceAlgorithm: The target distribution strategy. Common values:
	//     "round_robin". Future algorithms may be available depending on region.
	//   - TargetsType: The target category, typically "instances" or "raw_ips".
	//   - Targets: The actual list of targets to register (shape depends on TargetsType).
	//
	// Health checks:
	//   - Optionally link the backend to a health check by name using
	//     `HealthCheckName`. Only healthy targets receive traffic.
	NetworkBackendRequest struct {
		HealthCheckName  *string                       `json:"health_check_name,omitempty"`
		Name             string                        `json:"name"`
		Description      *string                       `json:"description,omitempty"`
		BalanceAlgorithm BackendBalanceAlgorithm       `json:"balance_algorithm"`
		TargetsType      BackendType                   `json:"targets_type"`
		Targets          *TargetsRawOrInstancesRequest `json:"targets,omitempty"`
	}

	// NetworkHealthCheckRequest describes a Health Check to be created with a Load Balancer.
	//
	// Supported protocols:
	//   - TCP: a successful TCP connect indicates health.
	//   - HTTP: HTTP request to `Path` must return a 2xx/3xx (configurable via `HealthyStatusCode`).
	//
	// Timing and thresholds:
	//   - `IntervalSeconds`, `TimeoutSeconds`, and threshold fields control how aggressively
	//     health is evaluated. Choose conservative defaults to avoid flapping.
	//   - `InitialDelaySeconds` can be used to give newly started targets time to warm up.
	NetworkHealthCheckRequest struct {
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

	// NetworkTLSCertificateRequest uploads a TLS certificate to be used by TLS listeners.
	//
	// Requirements:
	//   - `Certificate`: PEM-encoded certificate (can include chain).
	//   - `PrivateKey`: PEM-encoded private key that matches the certificate.
	//
	// Notes:
	//   - The certificate is referenced by `Name` from `NetworkListenerRequest.TLSCertificateName`.
	NetworkTLSCertificateRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Certificate string  `json:"certificate"`
		PrivateKey  string  `json:"private_key"`
	}

	// NetworkAclRequest defines an Access Control List (ACL) rule at the LB level.
	//
	// Rules are evaluated to allow or deny traffic by IP prefix and protocol.
	// Common use cases: allow office ranges, deny known bad networks, or lock down
	// internal LBs to specific subnets.
	//
	// Fields:
	//   - Ethertype: IPv4 or IPv6.
	//   - Protocol: typically TCP for L4 LB.
	//   - RemoteIPPrefix: CIDR, e.g., "203.0.113.0/24".
	//   - Action: ALLOW or DENY (see `AclActionType`).
	NetworkAclRequest struct {
		Name           *string       `json:"name,omitempty"`
		Ethertype      AclEtherType  `json:"ethertype"`
		Protocol       AclProtocol   `json:"protocol"`
		RemoteIPPrefix string        `json:"remote_ip_prefix"`
		Action         AclActionType `json:"action"`
	}

	// CreateNetworkLoadBalancerRequest is the payload to create a Network LB.
	//
	// Minimal required fields:
	//   - Name
	//   - Visibility: "internal" or "external" (see `LoadBalancerVisibility`)
	//   - VPCID
	//   - Listeners (at least one)
	//   - Backends (at least one)
	//
	// Optional:
	//   - Description, Type (e.g., "proxy"), SubnetPoolID, PublicIPID,
	//     HealthChecks, TLSCertificates and ACLs.
	//
	// Visibility and addressing:
	//   - `visibility == external`: associate a public IP via `PublicIPID`.
	//   - `visibility == internal`: reachable only inside the VPC.
	//
	// Cross-references:
	//   - `NetworkListenerRequest.BackendName` must exist in `Backends`.
	//   - For TLS listeners, `TLSCertificateName` must exist in `TLSCertificates`.
	CreateNetworkLoadBalancerRequest struct {
		Name            string                         `json:"name"`
		Description     *string                        `json:"description,omitempty"`
		Type            *string                        `json:"type,omitempty"`
		Visibility      LoadBalancerVisibility         `json:"visibility"`
		Listeners       []NetworkListenerRequest       `json:"listeners"`
		Backends        []NetworkBackendRequest        `json:"backends"`
		HealthChecks    []NetworkHealthCheckRequest    `json:"health_checks,omitempty"`
		TLSCertificates []NetworkTLSCertificateRequest `json:"tls_certificates,omitempty"`
		ACLs            []NetworkAclRequest            `json:"acls,omitempty"`
		VPCID           string                         `json:"vpc_id"`
		SubnetPoolID    *string                        `json:"subnet_pool_id,omitempty"`
		PublicIPID      *string                        `json:"public_ip_id,omitempty"`
	}

	// DeleteNetworkLoadBalancerRequest controls deletion behavior.
	//
	// Fields:
	//   - DeletePublicIP: when true, also deletes the associated public IP
	//     (if any). This translates to the query parameter `delete_public_ip`
	//     on the DELETE request. Default is to keep the IP (nil or false).
	DeleteNetworkLoadBalancerRequest struct {
		DeletePublicIP *bool `json:"-"`
	}

	// GetNetworkLoadBalancerRequest describes a "get" by id request.
	//
	// Note: The `Get` method below receives the id directly; this struct is
	// provided for symmetry and potential future expansion.
	GetNetworkLoadBalancerRequest struct {
		LoadBalancerID string `json:"-"`
	}

	// ListNetworkLoadBalancerRequest defines pagination and sorting for List.
	//
	// Fields map to query params:
	//   - Offset -> _offset
	//   - Limit  -> _limit
	//   - Sort   -> _sort (e.g., "created_at" or "name"; check availability).
	ListNetworkLoadBalancerRequest struct {
		Offset *int    `json:"-"`
		Limit  *int    `json:"-"`
		Sort   *string `json:"-"`
	}

	// NetworkBackendUpdateRequest updates an existing Backend.
	//
	// Fields:
	//   - ID (required): the backend identifier to update.
	//   - TargetsType/Targets: provide to replace or modify targets (semantics
	//     depend on the API; typically, the provided set becomes the desired state).
	//   - HealthCheckID: re-associate the backend with a different health check.
	NetworkBackendUpdateRequest struct {
		ID            string                              `json:"id"`
		HealthCheckID *string                             `json:"health_check_id,omitempty"`
		TargetsType   BackendType                         `json:"targets_type"`
		Targets       *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
	}

	// UpdateNetworkLoadBalancerRequest partially updates a Network LB and its
	// sub-resources (backends, health checks, TLS certificates).
	//
	// Notes:
	//   - Only fields provided are considered for update; omitted fields are left
	//     unchanged.
	//   - For sub-resources, elements are matched by their `ID`.
	//   - The exact replace/merge semantics for slice fields are defined by the API.
	//     In general, provide the new desired state for the items you want to change.
	UpdateNetworkLoadBalancerRequest struct {
		Name            *string                              `json:"name,omitempty"`
		Description     *string                              `json:"description,omitempty"`
		Backends        []NetworkBackendUpdateRequest        `json:"backends,omitempty"`
		HealthChecks    []NetworkHealthCheckUpdateRequest    `json:"health_checks,omitempty"`
		TLSCertificates []NetworkTLSCertificateUpdateRequest `json:"tls_certificates,omitempty"`
	}

	// NetworkHealthCheckUpdateRequest updates an existing Health Check.
	//
	// Required:
	//   - ID, Protocol, Port
	//
	// Optional:
	//   - Path (HTTP), timing, status/threshold parameters. Choose values that
	//     fit your application behavior to avoid false positives/negatives.
	NetworkHealthCheckUpdateRequest struct {
		ID                      string              `json:"id"`
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

	// NetworkTLSCertificateUpdateRequest replaces the PEM materials of an
	// existing TLS certificate.
	//
	// Required:
	//   - ID: certificate identifier to update.
	//   - Certificate, PrivateKey: PEM-encoded contents.
	NetworkTLSCertificateUpdateRequest struct {
		ID          string `json:"id"`
		Certificate string `json:"certificate"`
		PrivateKey  string `json:"private_key"`
	}

	// NetworkPublicIPResponse contains information about a public IP associated
	// with a Load Balancer.
	NetworkPublicIPResponse struct {
		ID         string  `json:"id"`
		IPAddress  *string `json:"ip_address,omitempty"`
		ExternalID string  `json:"external_id"`
	}

	// NetworkAclResponse is the persisted form of an ACL rule, as returned by
	// the API.
	NetworkAclResponse struct {
		ID             string       `json:"id"`
		Name           *string      `json:"name,omitempty"`
		Ethertype      AclEtherType `json:"ethertype"`
		Protocol       AclProtocol  `json:"protocol"`
		RemoteIPPrefix string       `json:"remote_ip_prefix"`
		Action         string       `json:"action"`
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
		Port                *string                         `json:"port,omitempty"`
		VPCID               string                          `json:"vpc_id"`
		SubnetPoolID        *string                         `json:"subnet_pool_id,omitempty"`
		CreatedAt           string                          `json:"created_at"`
		UpdatedAt           string                          `json:"updated_at"`
		LastOperationStatus *string                         `json:"last_operation_status,omitempty"`
	}

	// NetworkLBPaginatedResponse is the paginated response for listing LBs.
	NetworkLBPaginatedResponse struct {
		Results []NetworkLoadBalancerResponse `json:"results"`
	}

	// NetworkLoadBalancerService provides CRUD operations for Network LBs.
	//
	// Methods map to the REST API as follows:
	//   - Create -> POST   /v0beta1/network-load-balancers
	//   - Get    -> GET    /v0beta1/network-load-balancers/{id}
	//   - List   -> GET    /v0beta1/network-load-balancers
	//   - Update -> PUT    /v0beta1/network-load-balancers/{id}
	//   - Delete -> DELETE /v0beta1/network-load-balancers/{id}
	//
	// Example usage:
	//
	//   ctx := context.Background()
	//   svc := client.NetworkLoadBalancers()
	//
	//   id, err := svc.Create(ctx, createReq)
	//   if err != nil { /* handle */ }
	//
	//   lb, err := svc.Get(ctx, id)
	//   if err != nil { /* handle */ }
	//
	//   err = svc.Update(ctx, id, updateReq)
	//   if err != nil { /* handle */ }
	//
	//   err = svc.Delete(ctx, id, lbaas.DeleteNetworkLoadBalancerRequest{DeletePublicIP: ptr(true)})
	//   if err != nil { /* handle */ }
	NetworkLoadBalancerService interface {
		Create(ctx context.Context, create CreateNetworkLoadBalancerRequest) (string, error)
		Delete(ctx context.Context, id string, options DeleteNetworkLoadBalancerRequest) error
		Get(ctx context.Context, id string) (NetworkLoadBalancerResponse, error)
		List(ctx context.Context, options ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error)
		Update(ctx context.Context, id string, loadBalancer UpdateNetworkLoadBalancerRequest) error
	}

	// networkLoadBalancerService implements the NetworkLoadBalancerService interface.
	// It is typically constructed by `(*LbaasClient).NetworkLoadBalancers()`.
	networkLoadBalancerService struct {
		client *LbaasClient
	}
)

// Create creates a new Network Load Balancer and returns its ID.
//
// API: POST /v0beta1/network-load-balancers
//
// Notes:
//   - The operation may be asynchronous; use `Get` to poll `Status` until the
//     LB is fully provisioned.
//   - Validate cross-references between listeners/backends/certificates in the
//     request to avoid server-side rejections.
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
//
// API: DELETE /v0beta1/network-load-balancers/{id}?delete_public_ip=…
//
// Behavior:
//   - If `options.DeletePublicIP` is true, the request includes
//     `delete_public_ip=true`, instructing the API to remove the public IP
//     (if any) along with the LB.
//   - Deleting a non-existent LB should return an error from the server.
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

// Get retrieves detailed information about a specific Load Balancer.
//
// API: GET /v0beta1/network-load-balancers/{id}
//
// Returns:
//   - Full state of the LB, including listeners, backends, health checks,
//     ACLs, TLS certificates, and addressing details.
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
//
// API: GET /v0beta1/network-load-balancers?_offset=…&_limit=…&_sort=…
//
// Parameters:
//   - Offset: zero-based index for pagination.
//   - Limit: page size.
//   - Sort: field to sort by (e.g., "created_at", "name"; availability depends
//     on the environment).
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

// Update modifies a Network Load Balancer and/or its sub-resources.
//
// API: PUT /v0beta1/network-load-balancers/{id}
//
// Notes:
//   - Provide only the fields you intend to change.
//   - For sub-resources, include entries with their IDs.
//   - The server defines merge vs. replace semantics for slice fields; in
//     general, provide the desired state for the items you want to update.
func (s *networkLoadBalancerService) Update(ctx context.Context, id string, loadBalancer UpdateNetworkLoadBalancerRequest) error {
	path := urlNetworkLoadBalancer(&id)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, loadBalancer)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
