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
		Name           *string       `json:"name,omitempty"`
		Ethertype      AclEtherType  `json:"ethertype"`
		Action         AclActionType `json:"action"`
		Protocol       AclProtocol   `json:"protocol"`
		RemoteIPPrefix string        `json:"remote_ip_prefix"`
	}

	// NetworkAclResponse is the persisted form of an ACL rule, as returned by
	// the API.
	NetworkAclResponse struct {
		ID             string       `json:"id"`
		Name           *string      `json:"name,omitempty"`
		Ethertype      AclEtherType `json:"ethertype"`
		Protocol       AclProtocol  `json:"protocol"`
		RemoteIPPrefix string       `json:"remote_ip_prefix"`
		Action         string       `json:"action"`
	}

	UpdateNetworkACLRequest struct {
		Acls []CreateNetworkACLRequest `json:"acls"`
	}

	// NetworkACLService provides methods for managing network ACL rules
	NetworkACLService interface {
		Create(ctx context.Context, lbID string, req CreateNetworkACLRequest) (string, error)
		Delete(ctx context.Context, lbID, aclID string) error
		Replace(ctx context.Context, lbID string, req UpdateNetworkACLRequest) error
	}

	// networkACLService implements the NetworkACLService interface
	networkACLService struct {
		client *LbaasClient
	}
)

// Create creates a new network ACL rule
func (s *networkACLService) Create(ctx context.Context, lbID string, req CreateNetworkACLRequest) (string, error) {
	path := urlNetworkLoadBalancer(&lbID, acls)
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
func (s *networkACLService) Delete(ctx context.Context, lbID, aclID string) error {
	path := urlNetworkLoadBalancer(&lbID, acls, aclID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Replace updates the network ACL rules for a load balancer
func (s *networkACLService) Replace(ctx context.Context, lbID string, req UpdateNetworkACLRequest) error {
	path := urlNetworkLoadBalancer(&lbID, acls)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
