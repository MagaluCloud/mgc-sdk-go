package network

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	// PublicIPResponse represents a public IP resource response
	PublicIPResponse struct {
		ID          string                         `json:"id,omitempty"`
		ExternalID  string                         `json:"external_id,omitempty"`
		VPCID       string                         `json:"vpc_id,omitempty"`
		TenantID    string                         `json:"tenant_id,omitempty"`
		ProjectType string                         `json:"project_type,omitempty"`
		Description string                         `json:"description,omitempty"`
		PublicIP    string                         `json:"public_ip,omitempty"`
		PortID      string                         `json:"port_id,omitempty"`
		CreatedAt   utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated     utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Status      string                         `json:"status,omitempty"`
		Error       string                         `json:"error,omitempty"`
	}

	// PublicIPListResponse represents a list of public IPs response
	PublicIPListResponse struct {
		PublicIPs []PublicIPResponse `json:"public_ips"`
	}
)

// PublicIPService provides operations for managing Public IPs
type PublicIPService interface {
	// List retrieves all public IPs for the current tenant
	List(ctx context.Context) ([]PublicIPResponse, error)

	// Get retrieves details of a specific public IP by its ID
	Get(ctx context.Context, id string) (*PublicIPResponse, error)

	// Delete removes a public IP by its ID
	Delete(ctx context.Context, id string) error

	// AttachToPort associates a public IP with a specific port
	AttachToPort(ctx context.Context, publicIPID string, portID string) error

	// DetachFromPort removes the association between a public IP and a port
	DetachFromPort(ctx context.Context, publicIPID string, portID string) error
}

type publicIPService struct {
	client *NetworkClient
}

// List retrieves all public IPs for the current tenant
func (s *publicIPService) List(ctx context.Context) ([]PublicIPResponse, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[PublicIPListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/public_ips",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result.PublicIPs, nil
}

// Get retrieves details of a specific public IP by its ID
func (s *publicIPService) Get(ctx context.Context, id string) (*PublicIPResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[PublicIPResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/public_ips/%s", id),
		nil,
		nil,
	)
}

// Delete removes a public IP by its ID
func (s *publicIPService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/public_ips/%s", id),
		nil,
		nil,
	)
}

// AttachToPort associates a public IP with a specific port
func (s *publicIPService) AttachToPort(ctx context.Context, publicIPID string, portID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/public_ips/%s/attach/%s", publicIPID, portID),
		nil,
		nil,
	)
}

// DetachFromPort removes the association between a public IP and a port
func (s *publicIPService) DetachFromPort(ctx context.Context, publicIPID string, portID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/public_ips/%s/detach/%s", publicIPID, portID),
		nil,
		nil,
	)
}
