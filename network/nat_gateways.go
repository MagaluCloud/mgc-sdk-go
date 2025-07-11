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

type (
	// NatGatewayResponse represents a NAT Gateway resource
	NatGatewayResponse struct {
		// ID is the unique identifier of the NAT Gateway
		ID *string `json:"id,omitempty"`
		// Name is the name of the NAT Gateway (optional)
		Name *string `json:"name,omitempty"`
		// Description is the description of the NAT Gateway (optional)
		Description *string `json:"description,omitempty"`
		// VPCID is the VPC identifier (optional)
		VPCID *string `json:"vpc_id,omitempty"`
		// Zone is the availability zone (optional)
		Zone *string `json:"zone,omitempty"`
		// NatGatewayIP is the IP address of the NAT Gateway (optional)
		NatGatewayIP *string `json:"nat_gateway_ip,omitempty"`
		// CreatedAt is the creation timestamp (optional)
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Updated is the last update timestamp (optional)
		Updated *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		// Status is the current status of the NAT Gateway
		Status string `json:"status"`
	}

	// Meta represents pagination metadata
	Meta struct {
		// Page contains page information
		Page MetaPageInfo `json:"page"`
		// Links contains navigation links
		Links MetaLinks `json:"links"`
	}

	// MetaPageInfo represents pagination information
	MetaPageInfo struct {
		// Limit is the maximum number of items per page
		Limit int `json:"limit"`
		// Offset is the number of items skipped
		Offset int `json:"offset"`
		// Count is the number of items in the current page
		Count int `json:"count"`
		// Total is the total number of items
		Total int `json:"total"`
		// MaxItemsPerPage is the maximum number of items allowed per page
		MaxItemsPerPage int `json:"max_items_per_page"`
	}

	// MetaLinks represents navigation links
	MetaLinks struct {
		// Previous is the link to the previous page (optional)
		Previous *string `json:"previous,omitempty"`
		// Next is the link to the next page (optional)
		Next *string `json:"next,omitempty"`
		// Self is the link to the current page
		Self string `json:"self"`
	}

	// NatGatewayListResponse represents a NAT Gateway listing response
	NatGatewayListResponse struct {
		// Meta contains pagination metadata
		Meta Meta `json:"meta"`
		// Result contains the list of NAT Gateway resources
		Result []NatGatewayResponse `json:"result"`
	}

	// NatGatewayDetailsResponse represents detailed information about a NAT Gateway
	NatGatewayDetailsResponse struct {
		// ID is the unique identifier of the NAT Gateway
		ID *string `json:"id,omitempty"`
		// Name is the name of the NAT Gateway (optional)
		Name *string `json:"name,omitempty"`
		// Description is the description of the NAT Gateway (optional)
		Description *string `json:"description,omitempty"`
		// VPCID is the VPC identifier (optional)
		VPCID *string `json:"vpc_id,omitempty"`
		// Zone is the availability zone (optional)
		Zone *string `json:"zone,omitempty"`
		// NatGatewayIP is the IP address of the NAT Gateway (optional)
		NatGatewayIP *string `json:"nat_gateway_ip,omitempty"`
		// CreatedAt is the creation timestamp (optional)
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Updated is the last update timestamp (optional)
		Updated *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		// Status is the current status of the NAT Gateway
		Status string `json:"status"`
	}

	// CreateNatGatewayRequest represents the parameters for creating a new NAT Gateway
	CreateNatGatewayRequest struct {
		// Name is the name of the NAT Gateway
		Name string `json:"name"`
		// Description is the description of the NAT Gateway (optional)
		Description *string `json:"description,omitempty"`
		// Zone is the availability zone
		Zone string `json:"zone"`
		// VPCID is the VPC identifier
		VPCID string `json:"vpc_id"`
	}

	// NatGatewayCreateResponse represents the response after creating a NAT Gateway
	NatGatewayCreateResponse struct {
		// ID is the unique identifier of the created NAT Gateway
		ID string `json:"id"`
		// Status is the status of the created NAT Gateway
		Status string `json:"status"`
	}
)

// NatGatewayService provides operations for managing NAT Gateways
type NatGatewayService interface {
	// Create creates a new NAT Gateway with the provided configuration
	Create(ctx context.Context, req CreateNatGatewayRequest) (string, error)
	// Delete removes a NAT Gateway by its ID
	Delete(ctx context.Context, id string) error
	// Get retrieves details of a specific NAT Gateway by its ID
	Get(ctx context.Context, id string) (*NatGatewayDetailsResponse, error)
	// List retrieves all NAT Gateways for a specific VPC
	List(ctx context.Context, vpcID string, opts ListOptions) ([]NatGatewayResponse, error)
}

// natGatewayService implements the NatGatewayService interface
type natGatewayService struct {
	client *NetworkClient
}

// Create creates a new NAT Gateway with the provided configuration
func (s *natGatewayService) Create(ctx context.Context, req CreateNatGatewayRequest) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[NatGatewayCreateResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/nat_gateways",
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Delete removes a NAT Gateway by its ID
func (s *natGatewayService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v1/nat_gateways/%s", id),
		nil,
		nil,
	)
}

// Get retrieves details of a specific NAT Gateway by its ID
func (s *natGatewayService) Get(ctx context.Context, id string) (*NatGatewayDetailsResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[NatGatewayDetailsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v1/nat_gateways/%s", id),
		nil,
		nil,
	)
}

// List retrieves all NAT Gateways for a specific VPC
func (s *natGatewayService) List(ctx context.Context, vpcID string, opts ListOptions) ([]NatGatewayResponse, error) {
	queryParams := url.Values{}
	queryParams.Add("vpc_id", vpcID)

	if opts.Sort != nil {
		queryParams.Add("sort", *opts.Sort)
	}

	if opts.Limit != nil {
		queryParams.Add("items_per_page", strconv.Itoa(*opts.Limit))
	}

	if opts.Offset != nil {
		// Calculate the page based on offset and limit
		page := 1
		if opts.Limit != nil && *opts.Limit > 0 {
			page = (*opts.Offset / *opts.Limit) + 1
		}
		queryParams.Add("page", strconv.Itoa(page))
	} else {
		queryParams.Add("page", "1") // Default page
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[NatGatewayListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v1/nat_gateways",
		nil,
		queryParams,
	)
	if err != nil {
		return nil, err
	}

	return result.Result, nil
}
