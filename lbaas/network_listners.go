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
	panic("not implemented")
}

func (s *networkListenerService) Delete(ctx context.Context, req DeleteNetworkListenerRequest) error {
	panic("not implemented")
}

func (s *networkListenerService) Get(ctx context.Context, req GetNetworkListenerRequest) (*any, error) {
	panic("not implemented")
}

func (s *networkListenerService) List(ctx context.Context, req ListNetworkListenerRequest) ([]any, error) {
	panic("not implemented")
}

func (s *networkListenerService) Update(ctx context.Context, req UpdateNetworkListenerRequest) error {
	panic("not implemented")
}
