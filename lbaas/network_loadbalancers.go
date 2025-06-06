package lbaas

import (
	"context"
	"net/http"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// Structs auxiliares para criação de Load Balancer
	NetworkListenerRequest struct {
		TLSCertificateName *string          `json:"tls_certificate_name,omitempty"`
		Name               string           `json:"name"`
		Description        *string          `json:"description,omitempty"`
		BackendName        string           `json:"backend_name"`
		Protocol           ListenerProtocol `json:"protocol"`
		Port               int              `json:"port"`
	}

	NetworkBackendRequest struct {
		HealthCheckName  *string                       `json:"health_check_name,omitempty"`
		Name             string                        `json:"name"`
		Description      *string                       `json:"description,omitempty"`
		BalanceAlgorithm BackendBalanceAlgorithm       `json:"balance_algorithm"`
		TargetsType      BackendType                   `json:"targets_type"`
		Targets          *TargetsRawOrInstancesRequest `json:"targets,omitempty"`
	}

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

	NetworkTLSCertificateRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Certificate string  `json:"certificate"`
		PrivateKey  string  `json:"private_key"`
	}

	NetworkAclRequest struct {
		Name           *string       `json:"name,omitempty"`
		Ethertype      AclEtherType  `json:"ethertype"`
		Protocol       AclProtocol   `json:"protocol"`
		RemoteIPPrefix string        `json:"remote_ip_prefix"`
		Action         AclActionType `json:"action"`
	}

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
		ID            string                              `json:"id"`
		HealthCheckID *string                             `json:"health_check_id,omitempty"`
		TargetsType   BackendType                         `json:"targets_type"`
		Targets       *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
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
		ID             string       `json:"id"`
		Name           *string      `json:"name,omitempty"`
		Ethertype      AclEtherType `json:"ethertype"`
		Protocol       AclProtocol  `json:"protocol"`
		RemoteIPPrefix string       `json:"remote_ip_prefix"`
		Action         string       `json:"action"`
	}

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

func (s *networkLoadBalancerService) Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
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

func (s *networkLoadBalancerService) Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
