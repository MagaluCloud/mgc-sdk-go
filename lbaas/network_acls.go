package lbaas

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const acls = "acls"

type (
	CreateNetworkACLRequest struct {
		Name           *string       `json:"name,omitempty"`
		Ethertype      AclEtherType  `json:"ethertype"` // ipv4, ipv6
		LoadBalancerID string        `json:"load_balancer_id"`
		Action         AclActionType `json:"action"`   // ALLOW, DENY, DENY_UNSPECIFIED
		Protocol       AclProtocol   `json:"protocol"` // tcp, tls
		RemoteIPPrefix string        `json:"remote_ip_prefix"`
	}

	GetNetworkACLRequest struct {
		LoadBalancerID string `json:"-"`
		NetworkACLID   string `json:"-"`
	}

	ListNetworkACLRequest struct {
		LoadBalancerID string `json:"-"`
	}

	DeleteNetworkACLRequest struct {
		LoadBalancerID string `json:"load_balancer_id"`
		ID             string `json:"id"`
	}

	NetworkACLService interface {
		Create(ctx context.Context, req CreateNetworkACLRequest) (string, error)
		// Get(ctx context.Context, req GetNetworkACLRequest) (*NetworkAclResponse, error)
		// List(ctx context.Context, req ListNetworkACLRequest) ([]NetworkAclResponse, error)
		Delete(ctx context.Context, req DeleteNetworkACLRequest) error
	}

	networkACLService struct {
		client *LbaasClient
	}
)

func (s *networkACLService) Create(ctx context.Context, req CreateNetworkACLRequest) (string, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, acls)
	body := CreateNetworkACLRequest{
		Name:           req.Name,
		Ethertype:      req.Ethertype,
		Protocol:       req.Protocol,
		RemoteIPPrefix: req.RemoteIPPrefix,
		Action:         req.Action,
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return "", err
	}

	var resp struct {
		ID string `json:"id"`
	}

	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// func (s *networkACLService) Get(ctx context.Context, req GetNetworkACLRequest) (*NetworkAclResponse, error) {
// 	// GET /v0beta1/network-load-balancers/{load_balancer_id}/acls/{acl_id}
// 	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/acls/" + req.NetworkACLID

// 	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &NetworkAclResponse{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func (s *networkACLService) List(ctx context.Context, req ListNetworkACLRequest) ([]NetworkAclResponse, error) {
// 	// GET /v0beta1/network-load-balancers/{load_balancer_id}/acls
// 	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/acls"

// 	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &[]NetworkAclResponse{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return *result, nil
// }

func (s *networkACLService) Delete(ctx context.Context, req DeleteNetworkACLRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, acls, req.ID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
