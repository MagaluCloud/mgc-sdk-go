package network

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

// VpcsPeeringStatus represents the lifecycle status of a VPC peering.
type VpcsPeeringStatus string

const (
	VpcsPeeringStatusPending           VpcsPeeringStatus = "pending"
	VpcsPeeringStatusPendingRouteTable VpcsPeeringStatus = "pending_route"
	VpcsPeeringStatusCreated           VpcsPeeringStatus = "created"
	VpcsPeeringStatusDeleted           VpcsPeeringStatus = "deleted"
	VpcsPeeringStatusError             VpcsPeeringStatus = "error"
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
		ID          string                          `json:"id"`
		Name        string                          `json:"name"`
		Description *string                         `json:"description,omitempty"`
		Status      VpcsPeeringStatus               `json:"status"`
		CreatedAt   *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated     *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Members     []VpcsPeeringMember             `json:"members"`
	}

	// ListVpcsPeeringsOptions represents the filters and pagination accepted when
	// listing peerings.
	ListVpcsPeeringsOptions struct {
		// VpcID restricts the result to the peerings a given VPC takes part in.
		VpcID string
		// Defines the sorting in the format field:asc|desc.
		//
		// Default value: name:asc.
		Sort string
		// Page defines the page number (1-based).
		//
		// Default value: 1. Minimum value: 1.
		Page *int
		// ItemsPerPage defines the maximum number of items returned per page.
		//
		// Default value: 10. Minimum value: 1. Maximum value: 100.
		ItemsPerPage *int
	}

	// ListAllVpcsPeeringsOptions represents the filters accepted when listing
	// all peerings.
	ListAllVpcsPeeringsOptions struct {
		// VpcID restricts the result to the peerings a given VPC takes part in.
		VpcID string
		// Defines the sorting in the format field:asc|desc.
		Sort string
	}

	// ListVpcsPeeringsResponse represents a peering listing response.
	ListVpcsPeeringsResponse struct {
		Meta   Meta          `json:"meta"`
		Result []VpcsPeering `json:"result"`
	}

	// VpcsPeeringsCreateVpcs identifies the two VPCs taking part in a peering.
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
	// List retrieves the peerings of the current tenant, optionally filtered by VPC
	// and pagination.
	List(ctx context.Context, opts *ListVpcsPeeringsOptions) (*ListVpcsPeeringsResponse, error)
	// ListAll retrieves all peerings of the current tenant, optionally filtered
	// by VPC, automatically handling pagination.
	ListAll(ctx context.Context, opts *ListAllVpcsPeeringsOptions) ([]VpcsPeering, error)
	// Get retrieves the members of a peering.
	Get(ctx context.Context, peeringID string) (*VpcsPeering, error)
	// Create requests a new peering between two VPCs.
	Create(ctx context.Context, req VpcsPeeringsCreateRequest) (*VpcsPeeringsCreateResponse, error)
	// Delete removes a peering by its ID.
	Delete(ctx context.Context, peeringID string) error
}

type vpcsPeeringsService struct {
	client *NetworkClient
}

func (s *vpcsPeeringsService) List(ctx context.Context, opts *ListVpcsPeeringsOptions) (*ListVpcsPeeringsResponse, error) {
	if opts == nil {
		opts = &ListVpcsPeeringsOptions{}
	}

	query := make(url.Values)

	if opts.VpcID != "" {
		query.Set("vpc_id", opts.VpcID)
	}
	if opts.Sort != "" {
		query.Set("sort", opts.Sort)
	}
	if opts.Page != nil {
		query.Set("page", strconv.Itoa(*opts.Page))
	}
	if opts.ItemsPerPage != nil {
		query.Set("items_per_page", strconv.Itoa(*opts.ItemsPerPage))
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

func (s *vpcsPeeringsService) ListAll(ctx context.Context, opts *ListAllVpcsPeeringsOptions) ([]VpcsPeering, error) {
	if opts == nil {
		opts = &ListAllVpcsPeeringsOptions{}
	}

	allPeerings := []VpcsPeering{}
	page := 1
	itemsPerPage := 100

	for {
		currentPage := page
		resp, err := s.List(ctx, &ListVpcsPeeringsOptions{
			VpcID:        opts.VpcID,
			Sort:         opts.Sort,
			Page:         &currentPage,
			ItemsPerPage: &itemsPerPage,
		})
		if err != nil {
			return nil, err
		}

		allPeerings = append(allPeerings, resp.Result...)

		if page*itemsPerPage >= resp.Meta.Page.Total {
			break
		}

		page++
	}

	return allPeerings, nil
}

func (s *vpcsPeeringsService) Get(ctx context.Context, peeringID string) (*VpcsPeering, error) {
	if peeringID == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[VpcsPeering](
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
		return fmt.Errorf("id cannot be empty")
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
