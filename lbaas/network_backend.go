package lbaas

import "context"

type (
	CreateNetworkBackendRequest struct {
		BalanceAlgorithm string `json:"balance_algorithm"`
		Description      string `json:"description"`
		HealthCheckID    string `json:"health_check_id"`
		LoadBalancerID   string `json:"load_balancer_id"`
		Name             string `json:"name"`
	}

	DeleteNetworkBackendRequest struct{}
	GetNetworkBackendRequest    struct{}
	ListNetworkBackendRequest   struct{}
	UpdateNetworkBackendRequest struct{}

	NetworkBackendService interface {
		Create(ctx context.Context, req CreateNetworkBackendRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkBackendRequest) error
		Get(ctx context.Context, req GetNetworkBackendRequest) (*any, error)
		List(ctx context.Context, req ListNetworkBackendRequest) ([]any, error)
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
	// POST /v0beta1/network-load-balancers/{load_balancer_id}/backends
	panic("not implemented")
}

func (s *networkBackendService) Delete(ctx context.Context, req DeleteNetworkBackendRequest) error {
	// DELETE /v0beta1/network-load-balancers/{load_balancer_id}/backends/{backend_id}
	panic("not implemented")
}

func (s *networkBackendService) Get(ctx context.Context, req GetNetworkBackendRequest) (*any, error) {
	// GET /v0beta1/network-load-balancers/{load_balancer_id}/backends/{backend_id}
	panic("not implemented")
}

func (s *networkBackendService) List(ctx context.Context, req ListNetworkBackendRequest) ([]any, error) {
	// GET /v0beta1/network-load-balancers/{load_balancer_id}/backends
	panic("not implemented")
}

func (s *networkBackendService) Update(ctx context.Context, req UpdateNetworkBackendRequest) error {
	// PUT /v0beta1/network-load-balancers/{load_balancer_id}/backends/{backend_id}
	panic("not implemented")
}
