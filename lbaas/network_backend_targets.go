package lbaas

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const targets = "targets"

type (
	// CreateNetworkBackendTargetRequest represents the request payload for creating backend targets
	CreateNetworkBackendTargetRequest struct {
		LoadBalancerID   string      `json:"-"`
		NetworkBackendID string      `json:"-"`
		TargetsID        []string    `json:"targets_id"`
		TargetsType      BackendType `json:"targets_type"`
	}

	// DeleteNetworkBackendTargetRequest represents the request payload for deleting a backend target
	DeleteNetworkBackendTargetRequest struct {
		LoadBalancerID   string `json:"-"`
		NetworkBackendID string `json:"-"`
		TargetID         string `json:"-"`
	}

	// NetworkBackendTargetService provides methods for managing backend targets
	NetworkBackendTargetService interface {
		Create(ctx context.Context, req CreateNetworkBackendTargetRequest) (string, error)
		Delete(ctx context.Context, req DeleteNetworkBackendTargetRequest) error
	}

	// networkBackendTargetService implements the NetworkBackendTargetService interface
	networkBackendTargetService struct {
		client *LbaasClient
	}
)

// Create adds new targets to a backend
func (s *networkBackendTargetService) Create(ctx context.Context, req CreateNetworkBackendTargetRequest) (string, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, backends, req.NetworkBackendID, targets)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
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

// Delete removes a target from a backend
func (s *networkBackendTargetService) Delete(ctx context.Context, req DeleteNetworkBackendTargetRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, backends, req.NetworkBackendID, targets, req.TargetID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
