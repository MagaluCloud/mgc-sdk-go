package lbaas

import "context"

type (
	CreateNetworkBackendTargetRequest struct {
		LoadBalancerID   string   `json:"load_balancer_id"`
		NetworkBackendID string   `json:"network_backend_id"`
		TargetsID        []string `json:"targets_id"`
		TargetsType      string   `json:"targets_type"`
	}

	DeleteNetworkBackendTargetRequest struct {
		LoadBalancerID   string `json:"load_balancer_id"`
		NetworkBackendID string `json:"network_backend_id"`
		TargetID         string `json:"target_id"`
	}

	NetworkBackendTargetService interface {
		Create(ctx context.Context, req CreateNetworkBackendTargetRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkBackendTargetRequest) error
	}

	networkBackendTargetService struct {
		client *LbaasClient
	}
)
