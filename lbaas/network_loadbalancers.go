package lbaas

import (
	"context"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// Structs auxiliares para criação de Load Balancer
	NetworkListenerRequest struct {
		TLSCertificateName *string `json:"tls_certificate_name,omitempty"`
		Name               string  `json:"name"`
		Description        *string `json:"description,omitempty"`
		BackendName        string  `json:"backend_name"`
		Protocol           string  `json:"protocol"`
		Port               int     `json:"port"`
	}

	NetworkBackendRequest struct {
		HealthCheckName  *string                `json:"health_check_name,omitempty"`
		Name             string                 `json:"name"`
		Description      *string                `json:"description,omitempty"`
		BalanceAlgorithm string                 `json:"balance_algorithm"`
		TargetsType      string                 `json:"targets_type"`
		Targets          *TargetsRawOrInstances `json:"targets,omitempty"`
	}

	NetworkHealthCheckRequest struct {
		Name                    string  `json:"name"`
		Description             *string `json:"description,omitempty"`
		Protocol                string  `json:"protocol"`
		Path                    *string `json:"path,omitempty"`
		Port                    int     `json:"port"`
		HealthyStatusCode       *int    `json:"healthy_status_code,omitempty"`
		IntervalSeconds         *int    `json:"interval_seconds,omitempty"`
		TimeoutSeconds          *int    `json:"timeout_seconds,omitempty"`
		InitialDelaySeconds     *int    `json:"initial_delay_seconds,omitempty"`
		HealthyThresholdCount   *int    `json:"healthy_threshold_count,omitempty"`
		UnhealthyThresholdCount *int    `json:"unhealthy_threshold_count,omitempty"`
	}

	NetworkTLSCertificateRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Certificate string  `json:"certificate"`
		PrivateKey  string  `json:"private_key"`
	}

	NetworkAclRequest struct {
		Name           *string `json:"name,omitempty"`
		Ethertype      string  `json:"ethertype"`
		Protocol       string  `json:"protocol"`
		RemoteIPPrefix string  `json:"remote_ip_prefix"`
		Action         string  `json:"action"`
	}

	CreateNetworkLoadBalancerRequest struct {
		Name            string                         `json:"name"`
		Description     *string                        `json:"description,omitempty"`
		Type            *string                        `json:"type,omitempty"`
		Visibility      string                         `json:"visibility"`
		Listeners       []NetworkListenerRequest       `json:"listeners"`
		Backends        []NetworkBackendRequest        `json:"backends"`
		HealthChecks    []NetworkHealthCheckRequest    `json:"health_checks,omitempty"`
		TLSCertificates []NetworkTLSCertificateRequest `json:"tls_certificates,omitempty"`
		ACLs            []NetworkAclRequest            `json:"acls,omitempty"`
		VPCID           string                         `json:"vpc_id"`
		SubnetPoolID    *string                        `json:"subnet_pool_id,omitempty"`
		PublicIPID      *string                        `json:"public_ip_id,omitempty"`
		PanicThreshold  *int                           `json:"panic_threshold,omitempty"`
	}

	DeleteNetworkLoadBalancerRequest struct {
		LoadBalancerID string `json:"-"`
		DeletePublicIP *bool  `json:"-"`
	}

	GetNetworkLoadBalancerRequest struct {
		LoadBalancerID string `json:"-"`
	}

	ListNetworkLoadBalancerRequest struct {
		Offset *int    `json:"-"`
		Limit  *int    `json:"-"`
		Sort   *string `json:"-"`
	}

	NetworkBackendUpdateRequest struct {
		ID            string      `json:"id"`
		HealthCheckID *string     `json:"health_check_id,omitempty"`
		TargetsType   string      `json:"targets_type"`
		Targets       interface{} `json:"targets,omitempty"`
	}

	UpdateNetworkLoadBalancerRequest struct {
		LoadBalancerID  string                               `json:"-"`
		Name            *string                              `json:"name,omitempty"`
		Description     *string                              `json:"description,omitempty"`
		Backends        []NetworkBackendUpdateRequest        `json:"backends,omitempty"`
		HealthChecks    []NetworkHealthCheckUpdateRequest    `json:"health_checks,omitempty"`
		TLSCertificates []NetworkTLSCertificateUpdateRequest `json:"tls_certificates,omitempty"`
		PanicThreshold  *int                                 `json:"panic_threshold,omitempty"`
	}

	NetworkHealthCheckUpdateRequest struct {
		ID                      string  `json:"id"`
		Protocol                string  `json:"protocol"`
		Path                    *string `json:"path,omitempty"`
		Port                    int     `json:"port"`
		HealthyStatusCode       *int    `json:"healthy_status_code,omitempty"`
		IntervalSeconds         *int    `json:"interval_seconds,omitempty"`
		TimeoutSeconds          *int    `json:"timeout_seconds,omitempty"`
		InitialDelaySeconds     *int    `json:"initial_delay_seconds,omitempty"`
		HealthyThresholdCount   *int    `json:"healthy_threshold_count,omitempty"`
		UnhealthyThresholdCount *int    `json:"unhealthy_threshold_count,omitempty"`
	}

	NetworkTLSCertificateUpdateRequest struct {
		ID          string `json:"id"`
		Certificate string `json:"certificate"`
		PrivateKey  string `json:"private_key"`
	}

	// Response structs
	NetworkPublicIPResponse struct {
		ID         string  `json:"id"`
		IPAddress  *string `json:"ip_address,omitempty"`
		ExternalID string  `json:"external_id"`
	}

	NetworkAclResponse struct {
		ID             string  `json:"id"`
		Name           *string `json:"name,omitempty"`
		Ethertype      string  `json:"ethertype"`
		Protocol       string  `json:"protocol"`
		RemoteIPPrefix string  `json:"remote_ip_prefix"`
		Action         string  `json:"action"`
	}

	NetworkLoadBalancerResponse struct {
		ID                  string                          `json:"id"`
		Name                string                          `json:"name"`
		ProjectType         *string                         `json:"project_type,omitempty"`
		Description         *string                         `json:"description,omitempty"`
		Type                string                          `json:"type"`
		Visibility          string                          `json:"visibility"`
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

	NetworkLBPaginatedResponse struct {
		Results []NetworkLoadBalancerResponse `json:"results"`
	}

	NetworkLoadBalancerService interface {
		Create(ctx context.Context, req CreateNetworkLoadBalancerRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error
		Get(ctx context.Context, req GetNetworkLoadBalancerRequest) (*NetworkLoadBalancerResponse, error)
		List(ctx context.Context, req ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error)
		Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error
	}

	networkLoadBalancerService struct {
		client *LbaasClient
	}
)

func (s *networkLoadBalancerService) Create(ctx context.Context, req CreateNetworkLoadBalancerRequest) (string, error) {
	path := "/v0beta1/network-load-balancers"

	httpReq, err := s.client.newRequest(ctx, "POST", path, req)
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

func (s *networkLoadBalancerService) Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID

	httpReq, err := s.client.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	// Adicionar query parameter delete_public_ip se fornecido
	if req.DeletePublicIP != nil {
		query := httpReq.URL.Query()
		query.Set("delete_public_ip", strconv.FormatBool(*req.DeletePublicIP))
		httpReq.URL.RawQuery = query.Encode()
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

func (s *networkLoadBalancerService) Get(ctx context.Context, req GetNetworkLoadBalancerRequest) (*NetworkLoadBalancerResponse, error) {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID

	httpReq, err := s.client.newRequest(ctx, "GET", path, nil)
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

func (s *networkLoadBalancerService) List(ctx context.Context, req ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error) {
	path := "/v0beta1/network-load-balancers"

	httpReq, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Adicionar query parameters se fornecidos
	query := httpReq.URL.Query()
	if req.Offset != nil {
		query.Set("_offset", strconv.Itoa(*req.Offset))
	}
	if req.Limit != nil {
		query.Set("_limit", strconv.Itoa(*req.Limit))
	}
	if req.Sort != nil {
		query.Set("_sort", *req.Sort)
	}
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkLBPaginatedResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (s *networkLoadBalancerService) Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID

	httpReq, err := s.client.newRequest(ctx, "PUT", path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
