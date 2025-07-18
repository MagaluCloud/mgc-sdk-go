package network

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	// PublicIpResponsePort represents a public IP associated with a port
	PublicIpResponsePort struct {
		PublicIPID *string `json:"public_ip_id,omitempty"`
		PublicIP   *string `json:"public_ip,omitempty"`
	}

	// IpAddress represents an IP address configuration for a port
	IpAddress struct {
		IPAddress string  `json:"ip_address"`
		SubnetID  string  `json:"subnet_id"`
		Ethertype *string `json:"ethertype,omitempty"`
	}

	// PortResponse represents a network port resource
	PortResponse struct {
		ID                    *string                         `json:"id,omitempty"`
		Name                  *string                         `json:"name,omitempty"`
		Description           *string                         `json:"description,omitempty"`
		IsAdminStateUp        *bool                           `json:"is_admin_state_up,omitempty"`
		VPCID                 *string                         `json:"vpc_id,omitempty"`
		IsPortSecurityEnabled *bool                           `json:"is_port_security_enabled,omitempty"`
		SecurityGroups        *[]string                       `json:"security_groups"`
		PublicIP              *[]PublicIpResponsePort         `json:"public_ip"`
		IPAddress             *[]IpAddress                    `json:"ip_address"`
		CreatedAt             *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated               *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Network               *PortNetworkResponse            `json:"network,omitempty"`
	}

	// PortNetworkResponse represents the AvailabilityZone associated with a port
	PortNetworkResponse struct {
		AvailabilityZone *string `json:"availability_zone,omitempty"`
		ID               *string `json:"id,omitempty"`
		Zone             *string `json:"zone,omitempty"`
	}

	// PortUpdateRequest represents the fields available for update in a port resource
	PortUpdateRequest struct {
		IPSpoofingGuard *bool `json:"ip_spoofing_guard,omitempty"`
	}
)

// PortService provides operations for managing network ports
type PortService interface {
	List(ctx context.Context) ([]PortResponse, error)
	Get(ctx context.Context, id string) (*PortResponse, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, req PortUpdateRequest) error
	AttachSecurityGroup(ctx context.Context, portID string, securityGroupID string) error
	DetachSecurityGroup(ctx context.Context, portID string, securityGroupID string) error
}

// portService implements the PortService interface
type portService struct {
	client *NetworkClient
}

// List retrieves all ports for the current tenant
func (s *portService) List(ctx context.Context) ([]PortResponse, error) {
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]PortResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/ports",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Get retrieves details of a specific port by its ID
func (s *portService) Get(ctx context.Context, id string) (*PortResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[PortResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/ports/%s", id),
		nil,
		nil,
	)
}

// Delete removes a port by its ID
func (s *portService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/ports/%s", id),
		nil,
		nil,
	)
}

// Update patches a port by its ID considering the desired fields
func (s *portService) Update(ctx context.Context, id string, req PortUpdateRequest) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v0/ports/%s", id),
		req,
		nil,
	)
}

// AttachSecurityGroup associates a security group with a specific port
func (s *portService) AttachSecurityGroup(ctx context.Context, portID string, securityGroupID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/ports/%s/attach/%s", portID, securityGroupID),
		nil,
		nil,
	)
}

// DetachSecurityGroup removes the association between a security group and a port
func (s *portService) DetachSecurityGroup(ctx context.Context, portID string, securityGroupID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/ports/%s/detach/%s", portID, securityGroupID),
		nil,
		nil,
	)
}
