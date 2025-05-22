package lbaas

import "context"

type (
	CreateNetworkCertificateRequest struct{}
	DeleteNetworkCertificateRequest struct{}
	GetNetworkCertificateRequest    struct{}
	ListNetworkCertificateRequest   struct{}
	UpdateNetworkCertificateRequest struct{}

	NetworkCertificateService interface {
		Create(ctx context.Context, req CreateNetworkCertificateRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkCertificateRequest) error
		Get(ctx context.Context, req GetNetworkCertificateRequest) (*any, error)
		List(ctx context.Context, req ListNetworkCertificateRequest) ([]any, error)
		Update(ctx context.Context, req UpdateNetworkCertificateRequest) error
	}

	networkCertificateService struct {
		client *LbaasClient
	}
)

func (s *networkCertificateService) Create(ctx context.Context, req CreateNetworkCertificateRequest) (string, error) {
	panic("not implemented")
}

func (s *networkCertificateService) Delete(ctx context.Context, req DeleteNetworkCertificateRequest) error {
	panic("not implemented")
}

func (s *networkCertificateService) Get(ctx context.Context, req GetNetworkCertificateRequest) (*any, error) {
	panic("not implemented")
}

func (s *networkCertificateService) List(ctx context.Context, req ListNetworkCertificateRequest) ([]any, error) {
	panic("not implemented")
}

func (s *networkCertificateService) Update(ctx context.Context, req UpdateNetworkCertificateRequest) error {
	panic("not implemented")
}
