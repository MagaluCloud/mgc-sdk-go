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
		ID           *string                         `json:"id,omitempty"`
		Name         *string                         `json:"name,omitempty"`
		Description  *string                         `json:"description,omitempty"`
		VPCID        *string                         `json:"vpc_id,omitempty"`
		Zone         *string                         `json:"zone,omitempty"`
		NatGatewayIP *string                         `json:"nat_gateway_ip,omitempty"`
		CreatedAt    *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated      *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Status       string                          `json:"status"`
	}

	// Meta represents pagination metadata
	Meta struct {
		Page  MetaPageInfo `json:"page"`
		Links MetaLinks    `json:"links"`
	}

	// MetaPageInfo represents pagination information
	MetaPageInfo struct {
		Limit           int `json:"limit"`
		Offset          int `json:"offset"`
		Count           int `json:"count"`
		Total           int `json:"total"`
		MaxItemsPerPage int `json:"max_items_per_page"`
	}

	// MetaLinks represents navigation links
	MetaLinks struct {
		Previous *string `json:"previous,omitempty"`
		Next     *string `json:"next,omitempty"`
		Self     string  `json:"self"`
	}

	// NatGatewayListResponse represents a NAT Gateway listing response
	NatGatewayListResponse struct {
		Meta   Meta                 `json:"meta"`
		Result []NatGatewayResponse `json:"result"`
	}

	// NatGatewayDetailsResponse represents detailed information about a NAT Gateway
	NatGatewayDetailsResponse struct {
		ID           *string                         `json:"id,omitempty"`
		Name         *string                         `json:"name,omitempty"`
		Description  *string                         `json:"description,omitempty"`
		VPCID        *string                         `json:"vpc_id,omitempty"`
		Zone         *string                         `json:"zone,omitempty"`
		NatGatewayIP *string                         `json:"nat_gateway_ip,omitempty"`
		CreatedAt    *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated      *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Status       string                          `json:"status"`
	}

	// CreateNatGatewayRequest represents the parameters for creating a new NAT Gateway
	CreateNatGatewayRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Zone        string  `json:"zone"`
		VPCID       string  `json:"vpc_id"`
	}

	// NatGatewayCreateResponse represents the response after creating a NAT Gateway
	NatGatewayCreateResponse struct {
		ID     string `json:"id"`
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
