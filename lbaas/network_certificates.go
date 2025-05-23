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
	// POST /v0beta1/network-load-balancers/{load_balancer_id}/tls-certificates
	panic("not implemented")
}

func (s *networkCertificateService) Delete(ctx context.Context, req DeleteNetworkCertificateRequest) error {
	// DELETE /v0beta1/network-load-balancers/{load_balancer_id}/tls-certificates/{tls_certificate_id}
	panic("not implemented")
}

func (s *networkCertificateService) Get(ctx context.Context, req GetNetworkCertificateRequest) (*any, error) {
	// GET /v0beta1/network-load-balancers/{load_balancer_id}/tls-certificates/{tls_certificate_id}
	panic("not implemented")
}

func (s *networkCertificateService) List(ctx context.Context, req ListNetworkCertificateRequest) ([]any, error) {
	// GET /v0beta1/network-load-balancers/{load_balancer_id}/tls-certificates
	panic("not implemented")
}

func (s *networkCertificateService) Update(ctx context.Context, req UpdateNetworkCertificateRequest) error {
	// PUT /v0beta1/network-load-balancers/{load_balancer_id}/tls-certificates/{tls_certificate_id}
	panic("not implemented")
}
