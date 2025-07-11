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
	// ListSubnetPoolsOptions represents parameters for filtering and pagination
	ListSubnetPoolsOptions struct {
		// Limit specifies the maximum number of items to return
		Limit *int
		// Offset specifies the number of items to skip
		Offset *int
		// Sort specifies the field and direction for sorting results
		Sort *string
	}

	// ListSubnetPoolsResponse represents a list of subnet pools response
	ListSubnetPoolsResponse struct {
		// Meta contains pagination metadata
		Meta MetaModel `json:"meta"`
		// Results contains the list of subnet pool resources
		Results []SubnetPoolResponse `json:"results"`
	}

	// MetaModel represents pagination metadata
	MetaModel struct {
		// Page contains page information
		Page PageModel `json:"page"`
		// Links contains navigation links
		Links LinkModel `json:"links"`
	}

	// PageModel represents page information
	PageModel struct {
		// Limit is the maximum number of items per page (optional)
		Limit *int `json:"limit,omitempty"`
		// Offset is the number of items skipped (optional)
		Offset *int `json:"offset,omitempty"`
		// Count is the number of items in the current page
		Count int `json:"count"`
		// Total is the total number of items
		Total int `json:"total"`
	}

	// LinkModel represents navigation links
	LinkModel struct {
		// Previous is the link to the previous page (optional)
		Previous *string `json:"previous"`
		// Next is the link to the next page (optional)
		Next *string `json:"next"`
		// Self is the link to the current page
		Self string `json:"self"`
	}

	// SubnetPoolResponse represents a subnet pool resource response
	SubnetPoolResponse struct {
		// CIDR is the CIDR block of the subnet pool (optional)
		CIDR *string `json:"cidr,omitempty"`
		// ID is the unique identifier of the subnet pool
		ID string `json:"id"`
		// Name is the name of the subnet pool
		Name string `json:"name"`
		// TenantID is the tenant identifier
		TenantID string `json:"tenant_id"`
		// Description is the description of the subnet pool (optional)
		Description *string `json:"description,omitempty"`
		// IsDefault indicates if this is the default subnet pool
		IsDefault bool `json:"is_default"`
	}

	// SubnetPoolDetailsResponse represents detailed subnet pool information
	SubnetPoolDetailsResponse struct {
		// CIDR is the CIDR block of the subnet pool (optional)
		CIDR *string `json:"cidr,omitempty"`
		// ID is the unique identifier of the subnet pool
		ID string `json:"id"`
		// CreatedAt is the creation timestamp
		CreatedAt utils.LocalDateTimeWithoutZone `json:"created_at"`
		// TenantID is the tenant identifier
		TenantID string `json:"tenant_id"`
		// IPVersion is the IP version (4 for IPv4, 6 for IPv6)
		IPVersion int `json:"ip_version"`
		// IsDefault indicates if this is the default subnet pool
		IsDefault bool `json:"is_default"`
		// Name is the name of the subnet pool
		Name string `json:"name"`
		// Description is the description of the subnet pool
		Description string `json:"description"`
	}

	// CreateSubnetPoolRequest represents parameters for creating a new subnet pool
	CreateSubnetPoolRequest struct {
		// CIDR is the CIDR block for the subnet pool (optional)
		CIDR *string `json:"cidr,omitempty"`
		// Name is the name of the subnet pool
		Name string `json:"name"`
		// Description is the description of the subnet pool
		Description string `json:"description"`
		// Type is the type of the subnet pool (optional)
		Type *string `json:"type,omitempty"`
	}

	// BookCIDRRequest represents parameters for booking a CIDR range
	BookCIDRRequest struct {
		// CIDR is the CIDR block to book (optional)
		CIDR *string `json:"cidr,omitempty"`
		// Mask is the subnet mask (optional)
		Mask *int `json:"mask,omitempty"`
	}

	// BookCIDRResponse represents the response after booking a CIDR range
	BookCIDRResponse struct {
		// CIDR is the booked CIDR block
		CIDR string `json:"cidr"`
	}

	// UnbookCIDRRequest represents parameters for unbooking a CIDR range
	UnbookCIDRRequest struct {
		// CIDR is the CIDR block to unbook
		CIDR string `json:"cidr"`
	}

	// CreateSubnetPoolResponse represents the response after creating a subnet pool
	CreateSubnetPoolResponse struct {
		// ID is the unique identifier of the created subnet pool
		ID string `json:"id"`
	}
)

// SubnetPoolService provides operations for managing subnet pools
type SubnetPoolService interface {
	// List returns all subnet pools
	List(ctx context.Context, opts ListOptions) ([]SubnetPoolResponse, error)

	// Get retrieves a specific subnet pool
	Get(ctx context.Context, id string) (*SubnetPoolDetailsResponse, error)

	// Create provisions a new subnet pool
	Create(ctx context.Context, req CreateSubnetPoolRequest) (string, error)

	// Delete removes a subnet pool
	Delete(ctx context.Context, id string) error

	// BookCIDR books a CIDR range from a subnet pool
	BookCIDR(ctx context.Context, id string, req BookCIDRRequest) (*BookCIDRResponse, error)

	// UnbookCIDR releases a CIDR range from a subnet pool
	UnbookCIDR(ctx context.Context, id string, req UnbookCIDRRequest) error
}

// subnetPoolService implements the SubnetPoolService interface
type subnetPoolService struct {
	client *NetworkClient
}

// List retrieves all subnet pools for the current tenant
func (s *subnetPoolService) List(ctx context.Context, opts ListOptions) ([]SubnetPoolResponse, error) {
	query := make(url.Values)
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListSubnetPoolsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/subnetpools",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Get retrieves details of a specific subnet pool by its ID
func (s *subnetPoolService) Get(ctx context.Context, id string) (*SubnetPoolDetailsResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SubnetPoolDetailsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/subnetpools/%s", id),
		nil,
		nil,
	)
}

// Create creates a new subnet pool with the provided configuration
func (s *subnetPoolService) Create(ctx context.Context, req CreateSubnetPoolRequest) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[CreateSubnetPoolResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v0/subnetpools",
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Delete removes a subnet pool by its ID
func (s *subnetPoolService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/subnetpools/%s", id),
		nil,
		nil,
	)
}

// BookCIDR books a CIDR range from a subnet pool
func (s *subnetPoolService) BookCIDR(ctx context.Context, id string, req BookCIDRRequest) (*BookCIDRResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[BookCIDRResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/subnetpools/%s/book_cidr", id),
		req,
		nil,
	)
}

// UnbookCIDR releases a CIDR range from a subnet pool
func (s *subnetPoolService) UnbookCIDR(ctx context.Context, id string, req UnbookCIDRRequest) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/subnetpools/%s/unbook_cidr", id),
		req,
		nil,
	)
}
