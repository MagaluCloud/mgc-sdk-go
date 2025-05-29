package lbaas

import (
	"context"
	"net/http"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	CreateNetworkHealthCheckRequest struct {
		LoadBalancerID          string              `json:"-"`
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

	DeleteNetworkHealthCheckRequest struct {
		LoadBalancerID string `json:"-"`
		HealthCheckID  string `json:"-"`
	}

	GetNetworkHealthCheckRequest struct {
		LoadBalancerID string `json:"-"`
		HealthCheckID  string `json:"-"`
	}

	ListNetworkHealthCheckRequest struct {
		LoadBalancerID string  `json:"-"`
		Offset         *int    `json:"-"`
		Limit          *int    `json:"-"`
		Sort           *string `json:"-"`
	}

	UpdateNetworkHealthCheckRequest struct {
		LoadBalancerID          string              `json:"-"`
		HealthCheckID           string              `json:"-"`
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

	NetworkHealthCheckResponse struct {
		ID                      string              `json:"id"`
		Name                    string              `json:"name"`
		Description             *string             `json:"description,omitempty"`
		Protocol                HealthCheckProtocol `json:"protocol"`
		Path                    *string             `json:"path,omitempty"`
		Port                    int                 `json:"port"`
		HealthyStatusCode       int                 `json:"healthy_status_code"`
		IntervalSeconds         int                 `json:"interval_seconds"`
		TimeoutSeconds          int                 `json:"timeout_seconds"`
		InitialDelaySeconds     int                 `json:"initial_delay_seconds"`
		HealthyThresholdCount   int                 `json:"healthy_threshold_count"`
		UnhealthyThresholdCount int                 `json:"unhealthy_threshold_count"`
		CreatedAt               string              `json:"created_at"`
		UpdatedAt               string              `json:"updated_at"`
	}

	NetworkPaginatedHealthCheckResponse struct {
		Meta    interface{}                  `json:"meta"`
		Results []NetworkHealthCheckResponse `json:"results"`
	}

	NetworkHealthCheckService interface {
		Create(ctx context.Context, req CreateNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error)
		Delete(ctx context.Context, req DeleteNetworkHealthCheckRequest) error
		Get(ctx context.Context, req GetNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error)
		List(ctx context.Context, req ListNetworkHealthCheckRequest) ([]NetworkHealthCheckResponse, error)
		Update(ctx context.Context, req UpdateNetworkHealthCheckRequest) error
	}

	networkHealthCheckService struct {
		client *LbaasClient
	}
)

func (s *networkHealthCheckService) Create(ctx context.Context, req CreateNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, "health-checks")

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}

	var resp NetworkHealthCheckResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *networkHealthCheckService) Delete(ctx context.Context, req DeleteNetworkHealthCheckRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, "health-checks", req.HealthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

func (s *networkHealthCheckService) Get(ctx context.Context, req GetNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, "health-checks", req.HealthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkHealthCheckResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *networkHealthCheckService) List(ctx context.Context, req ListNetworkHealthCheckRequest) ([]NetworkHealthCheckResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, "health-checks")

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
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

	var resp NetworkPaginatedHealthCheckResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (s *networkHealthCheckService) Update(ctx context.Context, req UpdateNetworkHealthCheckRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, "health-checks", req.HealthCheckID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
