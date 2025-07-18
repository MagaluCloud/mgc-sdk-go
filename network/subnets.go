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
		Subnets []SubnetResponse `json:"subnets"`
	}

	// SubnetResponse represents a subnet resource response
	SubnetResponse struct {
		ID           string                          `json:"id"`
		VPCID        string                          `json:"vpc_id"`
		Name         *string                         `json:"name,omitempty"`
		Description  *string                         `json:"description,omitempty"`
		CIDRBlock    string                          `json:"cidr_block"`
		SubnetPoolID string                          `json:"subnetpool_id"`
		IPVersion    string                          `json:"ip_version"`
		Zone         string                          `json:"zone"`
		CreatedAt    *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated      *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
	}

	// SubnetResponseDetail represents a detailed subnet response
	SubnetResponseDetail struct {
		SubnetResponse
		GatewayIP      string        `json:"gateway_ip"`
		DNSNameservers []string      `json:"dns_nameservers"`
		DHCPPools      []DHCPPoolStr `json:"dhcp_pools"`
	}

	// DHCPPoolStr represents a DHCP pool configuration
	DHCPPoolStr struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	// SubnetCreateRequest represents parameters for creating a new subnet
	SubnetCreateRequest struct {
		Name           string    `json:"name"`
		Description    *string   `json:"description,omitempty"`
		CIDRBlock      string    `json:"cidr_block"`
		IPVersion      int       `json:"ip_version"`
		DNSNameservers *[]string `json:"dns_nameservers,omitempty"`
		SubnetPoolID   *string   `json:"subnetpool_id,omitempty"`
	}

	// SubnetCreateOptions represents additional options for subnet creation
	SubnetCreateOptions struct {
		Zone *string `json:"zone,omitempty"`
	}

	// SubnetPatchRequest represents parameters for updating a subnet
	SubnetPatchRequest struct {
		DNSNameservers *[]string `json:"dns_nameservers,omitempty"`
	}

	// SubnetCreateResponse represents the response after creating a subnet
	SubnetCreateResponse struct {
		ID string `json:"id"`
	}

	// SubnetResponseId represents a subnet ID response
	SubnetResponseId struct {
		ID string `json:"id"`
	}
)

// SubnetService provides operations for managing subnets
type SubnetService interface {
	Get(ctx context.Context, id string) (*SubnetResponseDetail, error)
	Delete(ctx context.Context, id string) error
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
