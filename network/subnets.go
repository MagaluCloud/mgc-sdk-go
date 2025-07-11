package network

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	// ListSubnetsResponse represents a list of subnets response
	ListSubnetsResponse struct {
		// Subnets contains the list of subnet resources
		Subnets []SubnetResponse `json:"subnets"`
	}

	// SubnetResponse represents a subnet resource response
	SubnetResponse struct {
		// ID is the unique identifier of the subnet
		ID string `json:"id"`
		// VPCID is the VPC identifier
		VPCID string `json:"vpc_id"`
		// Name is the name of the subnet (optional)
		Name *string `json:"name,omitempty"`
		// Description is the description of the subnet (optional)
		Description *string `json:"description,omitempty"`
		// CIDRBlock is the CIDR block of the subnet
		CIDRBlock string `json:"cidr_block"`
		// SubnetPoolID is the subnet pool identifier
		SubnetPoolID string `json:"subnetpool_id"`
		// IPVersion is the IP version (IPv4 or IPv6)
		IPVersion string `json:"ip_version"`
		// Zone is the availability zone
		Zone string `json:"zone"`
		// CreatedAt is the creation timestamp
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Updated is the last update timestamp
		Updated *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
	}

	// SubnetResponseDetail represents a detailed subnet response
	SubnetResponseDetail struct {
		SubnetResponse
		// GatewayIP is the gateway IP address
		GatewayIP string `json:"gateway_ip"`
		// DNSNameservers contains the DNS nameserver addresses
		DNSNameservers []string `json:"dns_nameservers"`
		// DHCPPools contains the DHCP pool configurations
		DHCPPools []DHCPPoolStr `json:"dhcp_pools"`
	}

	// DHCPPoolStr represents a DHCP pool configuration
	DHCPPoolStr struct {
		// Start is the starting IP address of the DHCP pool
		Start string `json:"start"`
		// End is the ending IP address of the DHCP pool
		End string `json:"end"`
	}

	// SubnetCreateRequest represents parameters for creating a new subnet
	SubnetCreateRequest struct {
		// Name is the name of the subnet
		Name string `json:"name"`
		// Description is the description of the subnet (optional)
		Description *string `json:"description,omitempty"`
		// CIDRBlock is the CIDR block for the subnet
		CIDRBlock string `json:"cidr_block"`
		// IPVersion is the IP version (4 for IPv4, 6 for IPv6)
		IPVersion int `json:"ip_version"`
		// DNSNameservers contains the DNS nameserver addresses (optional)
		DNSNameservers *[]string `json:"dns_nameservers,omitempty"`
		// SubnetPoolID is the subnet pool identifier (optional)
		SubnetPoolID *string `json:"subnetpool_id,omitempty"`
	}

	// SubnetCreateOptions represents additional options for subnet creation
	SubnetCreateOptions struct {
		// Zone specifies the availability zone (optional)
		Zone *string `json:"zone,omitempty"`
	}

	// SubnetPatchRequest represents parameters for updating a subnet
	SubnetPatchRequest struct {
		// DNSNameservers contains the new DNS nameserver addresses (optional)
		DNSNameservers *[]string `json:"dns_nameservers,omitempty"`
	}

	// SubnetCreateResponse represents the response after creating a subnet
	SubnetCreateResponse struct {
		// ID is the unique identifier of the created subnet
		ID string `json:"id"`
	}

	// SubnetResponseId represents a subnet ID response
	SubnetResponseId struct {
		// ID is the unique identifier of the subnet
		ID string `json:"id"`
	}
)

// SubnetService provides operations for managing subnets
type SubnetService interface {
	// Get retrieves details about a specific subnet
	Get(ctx context.Context, id string) (*SubnetResponseDetail, error)

	// Delete removes a subnet
	Delete(ctx context.Context, id string) error

	// Update modifies subnet properties
	Update(ctx context.Context, id string, req SubnetPatchRequest) (*SubnetResponseId, error)
}

// subnetService implements the SubnetService interface
type subnetService struct {
	client *NetworkClient
}

// Get retrieves details about a specific subnet
func (s *subnetService) Get(ctx context.Context, id string) (*SubnetResponseDetail, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SubnetResponseDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/subnets/%s", id),
		nil,
		nil,
	)
}

// Delete removes a subnet
func (s *subnetService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/subnets/%s", id),
		nil,
		nil,
	)
}

// Update modifies subnet properties
func (s *subnetService) Update(ctx context.Context, id string, req SubnetPatchRequest) (*SubnetResponseId, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SubnetResponseId](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v0/subnets/%s", id),
		req,
		nil,
	)
}
