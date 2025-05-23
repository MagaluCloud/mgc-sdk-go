package lbaas

import "context"

type (
	CreateNetworkListenerRequest struct{}
	DeleteNetworkListenerRequest struct{}
	GetNetworkListenerRequest    struct{}
	ListNetworkListenerRequest   struct{}
	UpdateNetworkListenerRequest struct{}

	NetworkListenerService interface {
		Create(ctx context.Context, req CreateNetworkListenerRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkListenerRequest) error
		Get(ctx context.Context, req GetNetworkListenerRequest) (*any, error)
		List(ctx context.Context, req ListNetworkListenerRequest) ([]any, error)
		Update(ctx context.Context, req UpdateNetworkListenerRequest) error
	}

	networkListenerService struct {
		client *LbaasClient
	}
)

func (s *networkListenerService) Create(ctx context.Context, req CreateNetworkListenerRequest) (string, error) {
	// POST /v0beta1/network-load-balancers/{load_balancer_id}/listeners
	panic("not implemented")
}

func (s *networkListenerService) Delete(ctx context.Context, req DeleteNetworkListenerRequest) error {
	// DELETE /v0beta1/network-load-balancers/{load_balancer_id}/listeners/{listener_id}
	panic("not implemented")
}

func (s *networkListenerService) Get(ctx context.Context, req GetNetworkListenerRequest) (*any, error) {
	// GET /v0beta1/network-load-balancers/{load_balancer_id}/listeners/{listener_id}
	panic("not implemented")
}

func (s *networkListenerService) List(ctx context.Context, req ListNetworkListenerRequest) ([]any, error) {
	// GET /v0beta1/network-load-balancers/{load_balancer_id}/listeners
	panic("not implemented")
}

func (s *networkListenerService) Update(ctx context.Context, req UpdateNetworkListenerRequest) error {
	// PUT /v0beta1/network-load-balancers/{load_balancer_id}/listeners/{listener_id}
	panic("not implemented")
}
