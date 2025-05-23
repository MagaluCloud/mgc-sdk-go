package lbaas

import "context"

type (
	CreateNetworkACLRequest struct {
		Name           string `json:"name"`
		Ethertype      string `json:"ethertype"`
		LoadBalancerID string `json:"load_balancer_id"`
		Action         string `json:"action"`
		Protocol       string `json:"protocol"`
		RemoteIPPrefix string `json:"remote_ip_prefix"`
	}

	DeleteNetworkACLRequest struct {
		LoadBalancerID string `json:"load_balancer_id"`
		ID             string `json:"id"`
	}

	NetworkACLService interface {
		Create(ctx context.Context, req CreateNetworkACLRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkACLRequest) error
	}

	networkACLService struct {
		client *LbaasClient
	}
)

func (s *networkACLService) Create(ctx context.Context, req CreateNetworkACLRequest) (string, error) {
	// GET /v0beta1/network-load-balancers/{load_balancer_id}/acls
	panic("not implemented")
}

func (s *networkACLService) Delete(ctx context.Context, req DeleteNetworkACLRequest) error {
	// DELETE /v0beta1/network-load-balancers/{load_balancer_id}/acls/{acl_id}
	panic("not implemented")
}
