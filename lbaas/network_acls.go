package lbaas

import (
	"context"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	CreateNetworkACLRequest struct {
		Name           *string `json:"name,omitempty"`
		Ethertype      string  `json:"ethertype"` // ipv4, ipv6
		LoadBalancerID string  `json:"load_balancer_id"`
		Action         string  `json:"action"`   // ALLOW, DENY, DENY_UNSPECIFIED
		Protocol       string  `json:"protocol"` // tcp, tls
		RemoteIPPrefix string  `json:"remote_ip_prefix"`
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
	// POST /v0beta1/network-load-balancers/{load_balancer_id}/acls
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/acls"

	body := CreateNetworkACLRequest{
		Name:           req.Name,
		Ethertype:      req.Ethertype,
		Protocol:       req.Protocol,
		RemoteIPPrefix: req.RemoteIPPrefix,
		Action:         req.Action,
	}

	httpReq, err := s.client.newRequest(ctx, "POST", path, body)
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

func (s *networkACLService) Delete(ctx context.Context, req DeleteNetworkACLRequest) error {
	// DELETE /v0beta1/network-load-balancers/{load_balancer_id}/acls/{acl_id}
	path := "/v0beta1/network-load-balancers/" + req.LoadBalancerID + "/acls/" + req.ID

	httpReq, err := s.client.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
