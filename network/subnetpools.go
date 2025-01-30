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
	// ListSubnetPoolsResponse represents a list of subnet pools response
	ListSubnetPoolsResponse struct {
		Meta    MetaModel            `json:"meta"`
		Results []SubnetPoolResponse `json:"results"`
	}

	MetaModel struct {
		Page  PageModel `json:"page"`
		Links LinkModel `json:"links"`
	}

	PageModel struct {
		Limit  int `json:"limit,omitempty"`
		Offset int `json:"offset,omitempty"`
		Count  int `json:"count"`
		Total  int `json:"total"`
	}

	LinkModel struct {
		Previous *string `json:"previous"`
		Next     *string `json:"next"`
		Self     string  `json:"self"`
	}

	// SubnetPoolResponse represents a subnet pool resource response
	SubnetPoolResponse struct {
		CIDR        string `json:"cidr,omitempty"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		TenantID    string `json:"tenant_id"`
		Description string `json:"description,omitempty"`
		IsDefault   bool   `json:"is_default"`
	}

	SubnetPoolDetailsResponse struct {
		CIDR        string                         `json:"cidr,omitempty"`
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
		CIDR        string `json:"cidr,omitempty"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type,omitempty"`
	}

	BookCIDRRequest struct {
		CIDR *string `json:"cidr,omitempty"`
		Mask *int    `json:"mask,omitempty"`
	}

	BookCIDRResponse struct {
		CIDR string `json:"cidr"`
	}

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
