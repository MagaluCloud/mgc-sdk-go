package lbaas

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const acls = "acls"

type (
	// CreateNetworkACLRequest represents the request payload for creating a network ACL rule
	CreateNetworkACLRequest struct {
		// Name is the name of the ACL rule (optional)
		Name *string `json:"name,omitempty"`
		// Ethertype is the ethernet type for the ACL rule (IPv4 or IPv6)
		Ethertype AclEtherType `json:"ethertype"`
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"load_balancer_id"`
		// Action is the action to take for matching traffic (ALLOW, DENY, DENY_UNSPECIFIED)
		Action AclActionType `json:"action"`
		// Protocol is the protocol for the ACL rule (TCP or TLS)
		Protocol AclProtocol `json:"protocol"`
		// RemoteIPPrefix is the remote IP prefix for the ACL rule
		RemoteIPPrefix string `json:"remote_ip_prefix"`
	}

	// GetNetworkACLRequest represents the request payload for getting a network ACL rule
	GetNetworkACLRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
		// NetworkACLID is the ID of the ACL rule to retrieve
		NetworkACLID string `json:"-"`
	}

	// ListNetworkACLRequest represents the request payload for listing network ACL rules
	ListNetworkACLRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"-"`
	}

	// DeleteNetworkACLRequest represents the request payload for deleting a network ACL rule
	DeleteNetworkACLRequest struct {
		// LoadBalancerID is the ID of the load balancer
		LoadBalancerID string `json:"load_balancer_id"`
		// ID is the ID of the ACL rule to delete
		ID string `json:"id"`
	}

	// NetworkACLService provides methods for managing network ACL rules
	NetworkACLService interface {
		// Create creates a new network ACL rule
		Create(ctx context.Context, req CreateNetworkACLRequest) (string, error)
		// Delete removes a network ACL rule
		Delete(ctx context.Context, req DeleteNetworkACLRequest) error
	}

	// networkACLService implements the NetworkACLService interface
	networkACLService struct {
		client *LbaasClient
	}
)

// Create creates a new network ACL rule
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

// Delete removes a network ACL rule
func (s *networkACLService) Delete(ctx context.Context, req DeleteNetworkACLRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, acls, req.ID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
