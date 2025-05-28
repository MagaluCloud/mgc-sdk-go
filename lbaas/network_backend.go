package lbaas

import (
	"context"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	NetworkBackendInstanceRequest struct {
		NicID string `json:"nic_id"`
		Port  int    `json:"port"`
	}

	NetworkBackendRawTargetRequest struct {
		IPAddress string `json:"ip_address"`
		Port      int    `json:"port"`
	}

	TargetsRawOrInstancesRequest struct {
		TargetsInstances []NetworkBackendInstanceRequest  `json:"-"`
		TargetsRaw       []NetworkBackendRawTargetRequest `json:"-"`
	}

	NetworkBackendInstanceUpdateRequest struct {
		NicID string `json:"nic_id"`
		Port  int    `json:"port"`
	}

	NetworkBackendRawTargetUpdateRequest struct {
		IPAddress string `json:"ip_address"`
		Port      int    `json:"port"`
	}

	TargetsRawOrInstancesUpdateRequest struct {
		TargetsInstances []NetworkBackendInstanceUpdateRequest  `json:"-"`
		TargetsRaw       []NetworkBackendRawTargetUpdateRequest `json:"-"`
	}

	CreateNetworkBackendRequest struct {
		LoadBalancerID   string                        `json:"-"`
		Name             string                        `json:"name"`
		Description      *string                       `json:"description,omitempty"`
		BalanceAlgorithm string                        `json:"balance_algorithm"`
		TargetsType      string                        `json:"targets_type"`
		Targets          *TargetsRawOrInstancesRequest `json:"targets,omitempty"`
		HealthCheckID    *string                       `json:"health_check_id,omitempty"`
	}

	DeleteNetworkBackendRequest struct {
		LoadBalancerID string `json:"-"`
		BackendID      string `json:"-"`
	}

	GetNetworkBackendRequest struct {
		LoadBalancerID string `json:"-"`
		BackendID      string `json:"-"`
	}

	ListNetworkBackendRequest struct {
		LoadBalancerID string `json:"-"`
	}

	UpdateNetworkBackendRequest struct {
		LoadBalancerID   string                              `json:"-"`
		BackendID        string                              `json:"-"`
		Name             *string                             `json:"name,omitempty"`
		Description      *string                             `json:"description,omitempty"`
		BalanceAlgorithm *string                             `json:"balance_algorithm,omitempty"`
		TargetsType      *string                             `json:"targets_type,omitempty"`
		Targets          *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
		TargetsInstances *[]NetworkBackendInstanceRequest    `json:"targets_instances,omitempty"`
		TargetsRaw       *[]NetworkBackendRawTargetRequest   `json:"targets_raw,omitempty"`
		HealthCheckID    *string                             `json:"health_check_id,omitempty"`
	}

	NetworkBackendInstanceResponse struct {
		ID        string  `json:"id"`
		IPAddress *string `json:"ip_address,omitempty"`
		NicID     string  `json:"nic_id,omitempty"`
		Port      int     `json:"port"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}

	NetworkBackendRawTargetResponse struct {
		ID        string  `json:"id"`
		IPAddress *string `json:"ip_address,omitempty"`
		Port      int     `json:"port"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}

	NetworkBackendResponse struct {
		ID               string      `json:"id"`
		HealthCheckID    *string     `json:"health_check_id,omitempty"`
		Name             string      `json:"name"`
		Description      *string     `json:"description,omitempty"`
		BalanceAlgorithm string      `json:"balance_algorithm"`
		TargetsType      string      `json:"targets_type"`
		Targets          interface{} `json:"targets"`
		CreatedAt        string      `json:"created_at"`
		UpdatedAt        string      `json:"updated_at"`
	}

	NetworkPaginatedBackendResponse struct {
		Results []NetworkBackendResponse `json:"results"`
	}

	NetworkBackendService interface {
		Create(ctx context.Context, req CreateNetworkBackendRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkBackendRequest) error
		Get(ctx context.Context, req GetNetworkBackendRequest) (*NetworkBackendResponse, error)
		List(ctx context.Context, req ListNetworkBackendRequest) ([]NetworkBackendResponse, error)
		Update(ctx context.Context, req UpdateNetworkBackendRequest) error
		Targets() *networkBackendTargetService
	}

	networkBackendService struct {
		client *LbaasClient
	}
)

func (s *networkBackendService) Targets() *networkBackendTargetService {
	return &networkBackendTargetService{client: s.client}
}

func (s *networkBackendService) Create(ctx context.Context, req CreateNetworkBackendRequest) (string, error) {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/backends"

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

func (s *networkBackendService) Delete(ctx context.Context, req DeleteNetworkBackendRequest) error {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/backends/" + req.BackendID

	httpReq, err := s.client.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

func (s *networkBackendService) Get(ctx context.Context, req GetNetworkBackendRequest) (*NetworkBackendResponse, error) {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/backends/" + req.BackendID

	httpReq, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkBackendResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *networkBackendService) List(ctx context.Context, req ListNetworkBackendRequest) ([]NetworkBackendResponse, error) {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/backends"

	httpReq, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkPaginatedBackendResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (s *networkBackendService) Update(ctx context.Context, req UpdateNetworkBackendRequest) error {
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/backends/" + req.BackendID

	httpReq, err := s.client.newRequest(ctx, "PUT", path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
