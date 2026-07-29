package network

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

// VpcsPeeringStatus represents the lifecycle status of a VPC peering.
type VpcsPeeringStatus string

const (
	VpcsPeeringStatusPending    VpcsPeeringStatus = "pending"
	VpcsPeeringStatusProcessing VpcsPeeringStatus = "processing"
	VpcsPeeringStatusCreated    VpcsPeeringStatus = "created"
	VpcsPeeringStatusUpdating   VpcsPeeringStatus = "updating"
	VpcsPeeringStatusDeleting   VpcsPeeringStatus = "deleting"
	VpcsPeeringStatusDeleted    VpcsPeeringStatus = "deleted"
	VpcsPeeringStatusError      VpcsPeeringStatus = "error"
)

// VpcsPeeringDirectRole represents the side a VPC takes in a peering.
type VpcsPeeringDirectRole string

const (
	VpcsPeeringDirectRoleRequester VpcsPeeringDirectRole = "requester"
	VpcsPeeringDirectRoleAccepter  VpcsPeeringDirectRole = "accepter"
)

type (
	// VpcsPeeringMember represents one of the VPCs connected by a peering.
	VpcsPeeringMember struct {
		ID         string                `json:"id"`
		VpcID      string                `json:"vpc_id"`
		DirectRole VpcsPeeringDirectRole `json:"direct_role"`
	}

	// VpcsPeering represents a peering connection between two VPCs.
	VpcsPeering struct {
		ID          string                          `json:"vpc_peering_id"`
		Name        string                          `json:"name"`
		Description *string                         `json:"description,omitempty"`
		Status      VpcsPeeringStatus               `json:"status"`
		CreatedAt   *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated     *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Members     []VpcsPeeringMember             `json:"members"`
	}

	// VpcsPeeringMembers represents the members of a peering and its current status.
	VpcsPeeringMembers struct {
		ID      string              `json:"vpc_peering_id"`
		Status  VpcsPeeringStatus   `json:"status"`
		Members []VpcsPeeringMember `json:"members"`
	}

	// ListVpcsPeeringsOptions represents the filters accepted when listing peerings.
	ListVpcsPeeringsOptions struct {
		// VpcID restricts the result to the peerings a given VPC takes part in.
		VpcID string
	}

	// ListVpcsPeeringsResponse represents a peering listing response.
	ListVpcsPeeringsResponse struct {
		Meta   Meta          `json:"meta"`
		Result []VpcsPeering `json:"result"`
	}

	// VpcsPeeringsCreateVpcs identifies the two VPCs taking part in a peering.
	// The requester asks for the connection and the accepter receives the invitation.
	VpcsPeeringsCreateVpcs struct {
		RequesterVpcID string `json:"requester_vpc_id"`
		AccepterVpcID  string `json:"accepter_vpc_id"`
	}

	// VpcsPeeringsCreateRequest represents the parameters for creating a new VPC peering.
	VpcsPeeringsCreateRequest struct {
		Name        string                 `json:"name"`
		Description *string                `json:"description,omitempty"`
		VPCs        VpcsPeeringsCreateVpcs `json:"vpcs"`
	}

	// VpcsPeeringsCreateResponse represents the response after requesting a VPC peering.
	VpcsPeeringsCreateResponse struct {
		ID     string            `json:"id"`
		Status VpcsPeeringStatus `json:"status"`
	}
)

// VpcsPeeringsService defines operations for managing VPC peerings.
type VpcsPeeringsService interface {
	// List retrieves the peerings of the current tenant, optionally filtered by VPC.
	//
	// The API does not accept pagination parameters on this endpoint, so the result is
	// whatever page the server chooses to return: check Meta.Page to detect truncation.
	List(ctx context.Context, opts *ListVpcsPeeringsOptions) (*ListVpcsPeeringsResponse, error)
	// GetMembers retrieves the members of a peering and its current status.
	//
	// The API has no endpoint returning a single peering in full: name, description and
	// timestamps are only available through List.
	GetMembers(ctx context.Context, peeringID string) (*VpcsPeeringMembers, error)
	// Create requests a new peering between two VPCs.
	Create(ctx context.Context, req VpcsPeeringsCreateRequest) (*VpcsPeeringsCreateResponse, error)
	// Delete removes a peering by its ID.
	Delete(ctx context.Context, peeringID string) error
}

type vpcsPeeringsService struct {
	client *NetworkClient
}

func (s *vpcsPeeringsService) List(ctx context.Context, opts *ListVpcsPeeringsOptions) (*ListVpcsPeeringsResponse, error) {
	query := make(url.Values)

	if opts != nil && opts.VpcID != "" {
		query.Set("vpc_id", opts.VpcID)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ListVpcsPeeringsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v1/vpcs_peerings",
		nil,
		query,
	)
}

func (s *vpcsPeeringsService) GetMembers(ctx context.Context, peeringID string) (*VpcsPeeringMembers, error) {
	if peeringID == "" {
		return nil, fmt.Errorf("vpc_peering_id cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[VpcsPeeringMembers](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v1/vpcs_peerings/%s", peeringID),
		nil,
		nil,
	)
}

func (s *vpcsPeeringsService) Create(ctx context.Context, req VpcsPeeringsCreateRequest) (*VpcsPeeringsCreateResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if req.VPCs.RequesterVpcID == "" {
		return nil, fmt.Errorf("requester_vpc_id cannot be empty")
	}
	if req.VPCs.AccepterVpcID == "" {
		return nil, fmt.Errorf("accepter_vpc_id cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[VpcsPeeringsCreateResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/vpcs_peerings",
		req,
		nil,
	)
}

func (s *vpcsPeeringsService) Delete(ctx context.Context, peeringID string) error {
	if peeringID == "" {
		return fmt.Errorf("vpc_peering_id cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v1/vpcs_peerings/%s", peeringID),
		nil,
		nil,
	)
}
