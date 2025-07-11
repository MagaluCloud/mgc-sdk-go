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
		// PublicIPID is the unique identifier of the public IP
		PublicIPID *string `json:"public_ip_id,omitempty"`
		// PublicIP is the public IP address
		PublicIP *string `json:"public_ip,omitempty"`
	}

	// IpAddress represents an IP address configuration for a port
	IpAddress struct {
		// IPAddress is the IP address
		IPAddress string `json:"ip_address"`
		// SubnetID is the subnet identifier
		SubnetID string `json:"subnet_id"`
		// Ethertype is the ethernet type (optional)
		Ethertype *string `json:"ethertype,omitempty"`
	}

	// PortResponse represents a network port resource
	PortResponse struct {
		// ID is the unique identifier of the port
		ID *string `json:"id,omitempty"`
		// Name is the name of the port (optional)
		Name *string `json:"name,omitempty"`
		// Description is the description of the port (optional)
		Description *string `json:"description,omitempty"`
		// IsAdminStateUp indicates if the port is administratively up (optional)
		IsAdminStateUp *bool `json:"is_admin_state_up,omitempty"`
		// VPCID is the VPC identifier (optional)
		VPCID *string `json:"vpc_id,omitempty"`
		// IsPortSecurityEnabled indicates if port security is enabled (optional)
		IsPortSecurityEnabled *bool `json:"is_port_security_enabled,omitempty"`
		// SecurityGroups contains the list of security group IDs (optional)
		SecurityGroups *[]string `json:"security_groups"`
		// PublicIP contains the public IP configurations (optional)
		PublicIP *[]PublicIpResponsePort `json:"public_ip"`
		// IPAddress contains the IP address configurations (optional)
		IPAddress *[]IpAddress `json:"ip_address"`
		// CreatedAt is the creation timestamp (optional)
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Updated is the last update timestamp (optional)
		Updated *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		// Network contains the network information (optional)
		Network *PortNetworkResponse `json:"network,omitempty"`
	}

	// PortNetworkResponse represents the AvailabilityZone associated with a port
	PortNetworkResponse struct {
		// AvailabilityZone is the availability zone (optional)
		AvailabilityZone *string `json:"availability_zone,omitempty"`
		// ID is the network identifier (optional)
		ID *string `json:"id,omitempty"`
		// Zone is the zone identifier (optional)
		Zone *string `json:"zone,omitempty"`
	}

	// PortUpdateRequest represents the fields available for update in a port resource
	PortUpdateRequest struct {
		// IPSpoofingGuard allows spoofed packets to enter a port (optional)
		IPSpoofingGuard *bool `json:"ip_spoofing_guard,omitempty"`
	}
)

// PortService provides operations for managing network ports
type PortService interface {
	// List retrieves all ports for the current tenant
	List(ctx context.Context) ([]PortResponse, error)

	// Get retrieves details of a specific port by its ID
	Get(ctx context.Context, id string) (*PortResponse, error)

	// Delete removes a port by its ID
	Delete(ctx context.Context, id string) error

	// Update updates a port
	Update(ctx context.Context, id string, req PortUpdateRequest) error

	// AttachSecurityGroup associates a security group with a specific port
	AttachSecurityGroup(ctx context.Context, portID string, securityGroupID string) error

	// DetachSecurityGroup removes the association between a security group and a port
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
