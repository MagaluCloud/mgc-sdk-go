package lbaas

import "context"

type (
	CreateNetworkHealthCheckRequest struct{}
	DeleteNetworkHealthCheckRequest struct{}
	GetNetworkHealthCheckRequest    struct{}
	ListNetworkHealthCheckRequest   struct{}
	UpdateNetworkHealthCheckRequest struct{}

	NetworkHealthCheckService interface {
		Create(ctx context.Context, req CreateNetworkHealthCheckRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkHealthCheckRequest) error
		Get(ctx context.Context, req GetNetworkHealthCheckRequest) (*any, error)
		List(ctx context.Context, req ListNetworkHealthCheckRequest) ([]any, error)
		Update(ctx context.Context, req UpdateNetworkHealthCheckRequest) error
	}

	networkHealthCheckService struct {
		client *LbaasClient
	}
)

func (s *networkHealthCheckService) Create(ctx context.Context, req CreateNetworkHealthCheckRequest) (string, error) {
	panic("not implemented")
}

func (s *networkHealthCheckService) Delete(ctx context.Context, req DeleteNetworkHealthCheckRequest) error {
	panic("not implemented")
}

func (s *networkHealthCheckService) Get(ctx context.Context, req GetNetworkHealthCheckRequest) (*any, error) {
	panic("not implemented")
}

func (s *networkHealthCheckService) List(ctx context.Context, req ListNetworkHealthCheckRequest) ([]any, error) {
	panic("not implemented")
}

func (s *networkHealthCheckService) Update(ctx context.Context, req UpdateNetworkHealthCheckRequest) error {
	panic("not implemented")
}
