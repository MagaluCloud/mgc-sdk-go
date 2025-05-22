package lbaas

import "context"

type (
	CreateNetworkLoadBalancerRequest struct{}
	DeleteNetworkLoadBalancerRequest struct{}
	GetNetworkLoadBalancerRequest    struct{}
	ListNetworkLoadBalancerRequest   struct{}
	UpdateNetworkLoadBalancerRequest struct{}

	NetworkLoadBalancerService interface {
		Create(ctx context.Context, req CreateNetworkLoadBalancerRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error
		Get(ctx context.Context, req GetNetworkLoadBalancerRequest) (*any, error)
		List(ctx context.Context, req ListNetworkLoadBalancerRequest) ([]any, error)
		Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error
	}

	networkLoadBalancerService struct {
		client *LbaasClient
	}
)

func (s *networkLoadBalancerService) Create(ctx context.Context, req CreateNetworkLoadBalancerRequest) (string, error) {
	panic("not implemented")
}

func (s *networkLoadBalancerService) Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error {
	panic("not implemented")
}

func (s *networkLoadBalancerService) Get(ctx context.Context, req GetNetworkLoadBalancerRequest) (*any, error) {
	panic("not implemented")
}

func (s *networkLoadBalancerService) List(ctx context.Context, req ListNetworkLoadBalancerRequest) ([]any, error) {
	panic("not implemented")
}

func (s *networkLoadBalancerService) Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error {
	panic("not implemented")
}
