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

func (s *networkBackendTargetService) Create(ctx context.Context, req CreateNetworkBackendTargetRequest) (string, error) {
	// POST /v0beta1/network-load-balancers/{load_balancer_id}/backends/{backend_id}/targets
	panic("not implemented")
}

func (s *networkBackendTargetService) Delete(ctx context.Context, req DeleteNetworkBackendTargetRequest) error {
	// DELETE /v0beta1/network-load-balancers/{load_balancer_id}/backends/{backend_id}/targets/{target_id}
	panic("not implemented")
}
