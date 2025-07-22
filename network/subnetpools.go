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
		Limit  *int
		Offset *int
		Sort   *string
	}

	// ListSubnetPoolsResponse represents a list of subnet pools response
	ListSubnetPoolsResponse struct {
		Meta    MetaModel            `json:"meta"`
		Results []SubnetPoolResponse `json:"results"`
	}

	// MetaModel represents pagination metadata
	MetaModel struct {
		Page  PageModel `json:"page"`
		Links LinkModel `json:"links"`
	}

	// PageModel represents page information
	PageModel struct {
		Limit  *int `json:"limit,omitempty"`
		Offset *int `json:"offset,omitempty"`
		Count  int  `json:"count"`
		Total  int  `json:"total"`
	}

	// LinkModel represents navigation links
	LinkModel struct {
		Previous *string `json:"previous"`
		Next     *string `json:"next"`
		Self     string  `json:"self"`
	}

	// SubnetPoolResponse represents a subnet pool resource response
	SubnetPoolResponse struct {
		CIDR        *string `json:"cidr,omitempty"`
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		TenantID    string  `json:"tenant_id"`
		Description *string `json:"description,omitempty"`
		IsDefault   bool    `json:"is_default"`
	}

	// SubnetPoolDetailsResponse represents detailed subnet pool information
	SubnetPoolDetailsResponse struct {
		CIDR        *string                        `json:"cidr,omitempty"`
		ID          string                         `json:"id"`
		CreatedAt   utils.LocalDateTimeWithoutZone `json:"created_at"`
		TenantID    string                         `json:"tenant_id"`
		IPVersion   int                            `json:"ip_version"`
		IsDefault   bool                           `json:"is_default"`
		Name        string                         `json:"name"`
		Description string                         `json:"description"`
	}

	// CreateSubnetPoolRequest represents parameters for creating a new subnet pool
	CreateSubnetPoolRequest struct {
		CIDR        *string `json:"cidr,omitempty"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Type        *string `json:"type,omitempty"`
	}

	// BookCIDRRequest represents parameters for booking a CIDR range
	BookCIDRRequest struct {
		CIDR *string `json:"cidr,omitempty"`
		Mask *int    `json:"mask,omitempty"`
	}

	// BookCIDRResponse represents the response after booking a CIDR range
	BookCIDRResponse struct {
		CIDR string `json:"cidr"`
	}

	// UnbookCIDRRequest represents parameters for unbooking a CIDR range
	UnbookCIDRRequest struct {
		CIDR string `json:"cidr"`
	}

	// CreateSubnetPoolResponse represents the response after creating a subnet pool
	CreateSubnetPoolResponse struct {
		ID string `json:"id"`
	}
)

// SubnetPoolService provides operations for managing subnet pools
type SubnetPoolService interface {
	List(ctx context.Context, opts ListOptions) ([]SubnetPoolResponse, error)
	Get(ctx context.Context, id string) (*SubnetPoolDetailsResponse, error)
	Create(ctx context.Context, req CreateSubnetPoolRequest) (string, error)
	Delete(ctx context.Context, id string) error
	BookCIDR(ctx context.Context, id string, req BookCIDRRequest) (*BookCIDRResponse, error)
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
