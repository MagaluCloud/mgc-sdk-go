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
		HealthCheckID *string                               `json:"health_check_id,omitempty"`
		TargetsType   BackendType                           `json:"targets_type"`
		Targets       []NetworkBackendInstanceTargetRequest `json:"targets"`
	}

	// NetworkBackendTargetService provides methods for managing backend targets
	NetworkBackendTargetService interface {
		Create(ctx context.Context, lbID, backendID string, req CreateNetworkBackendTargetRequest) (string, error)
		Replace(ctx context.Context, lbID, backendID string, req CreateNetworkBackendTargetRequest) (string, error)
		Delete(ctx context.Context, lbID, backendID, targetID string) error
	}

	// networkBackendTargetService implements the NetworkBackendTargetService interface
	networkBackendTargetService struct {
		client *LbaasClient
	}
)

// Create adds new targets to a backend
func (s *networkBackendTargetService) Create(ctx context.Context, lbID, backendID string, req CreateNetworkBackendTargetRequest) (string, error) {
	path := urlNetworkLoadBalancer(&lbID, backends, backendID, targets)

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

// Replace replaces all targets in a backend
func (s *networkBackendTargetService) Replace(ctx context.Context, lbID, backendID string, req CreateNetworkBackendTargetRequest) (string, error) {
	path := urlNetworkLoadBalancer(&lbID, backends, backendID, targets)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
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
func (s *networkBackendTargetService) Delete(ctx context.Context, lbID, backendID, targetID string) error {
	path := urlNetworkLoadBalancer(&lbID, backends, backendID, targets, targetID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
