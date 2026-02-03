package network

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type RouteStatus string

const (
	RouteStatusProcessing RouteStatus = "processing"
	RouteStatusCreated    RouteStatus = "created"
	RouteStatusPending    RouteStatus = "pending"
	RouteStatusDeleting   RouteStatus = "deleting"
	RouteStatusDeleted    RouteStatus = "deleted"
	RouteStatusUpdating   RouteStatus = "updating"
	RouteStatusError      RouteStatus = "error"
)

type (
	RouteDetail struct {
		ID              string      `json:"id"`
		PortID          string      `json:"port_id"`
		CIDRDestination string      `json:"cidr_destination"`
		Description     string      `json:"description,omitempty"`
		NextHop         string      `json:"next_hop"`
		Type            string      `json:"type"`
		Status          RouteStatus `json:"status"`
	}

	Route struct {
		RouteDetail
		VpcID string `json:"vpc_id"`
	}

	CreateRequest struct {
		PortID          string  `json:"port_id"`
		CIDRDestination string  `json:"cidr_destination"`
		Description     *string `json:"description"`
	}

	CreateResponse struct {
		ID     string      `json:"id"`
		Status RouteStatus `json:"status"`
	}

	ListRouteOptions struct {
		// Zone filters routes by availability zone.
		Zone string
		// Defines the sorting in the format field:asc|desc.
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

	ListAllRoutesOptions struct {
		Zone string
		Sort string
	}

	ListLinks struct {
		Next     *string `json:"next,omitempty"`
		Previous *string `json:"previous,omitempty"`
		Self     string  `json:"self"`
	}

	ListPage struct {
		Count           int `json:"count"`
		Limit           int `json:"limit"`
		Offset          int `json:"offset"`
		Total           int `json:"total"`
		MaxItemsPerPage int `json:"max_items_per_page"`
	}

	ListMeta struct {
		Page  ListPage  `json:"page"`
		Links ListLinks `json:"links"`
	}

	ListResponse struct {
		Result []RouteDetail `json:"result"`
		Meta   ListMeta      `json:"meta"`
	}
)

// RouteService defines operations for managing VPC routes.
type RouteService interface {
	// List retrieves a paginated list of routes for a given VPC.
	List(ctx context.Context, vpcID string, opts *ListRouteOptions) (*ListResponse, error)
	// ListAll retrieves all routes for a given VPC, automatically handling pagination.
	ListAll(ctx context.Context, vpcID string, opts *ListAllRoutesOptions) ([]RouteDetail, error)
	// Get retrieves a single route by its ID.
	Get(ctx context.Context, vpcID, routeID string) (*Route, error)
	// Create creates a new route in the specified VPC.
	Create(ctx context.Context, vpcID string, req CreateRequest) (*CreateResponse, error)
	// Delete removes a route from the specified VPC.
	Delete(ctx context.Context, vpcID, routeID string) error
}

type routeService struct {
	client *NetworkClient
}

func (s *routeService) List(ctx context.Context, vpcID string, opts *ListRouteOptions) (*ListResponse, error) {
	query := make(url.Values)

	if opts == nil {
		opts = &ListRouteOptions{}
	}

	if opts.Zone != "" {
		query.Set("zone", opts.Zone)
	}
	if opts.Sort != "" {
		err := validateSortValue(opts.Sort)
		if err != nil {
			return nil, err
		}

		query.Set("sort", opts.Sort)
	}
	if opts.Page != nil {
		query.Set("page", strconv.Itoa(*opts.Page))
	}
	if opts.ItemsPerPage != nil {
		query.Set("items_per_page", strconv.Itoa(*opts.ItemsPerPage))
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v1/vpcs/%s/route_table/routes", vpcID),
		nil,
		query,
	)
}

func (s *routeService) ListAll(ctx context.Context, vpcID string, opts *ListAllRoutesOptions) ([]RouteDetail, error) {
	allRoutes := []RouteDetail{}
	page := 1
	itemsPerPage := 100

	if opts == nil {
		opts = &ListAllRoutesOptions{}
	}

	for {
		currentPage := page
		listOpts := ListRouteOptions{
			Page:         &currentPage,
			ItemsPerPage: &itemsPerPage,
			Sort:         opts.Sort,
			Zone:         opts.Zone,
		}

		resp, err := s.List(ctx, vpcID, &listOpts)
		if err != nil {
			return nil, err
		}

		allRoutes = append(allRoutes, resp.Result...)

		if page*itemsPerPage >= resp.Meta.Page.Total {
			break
		}

		page++
	}

	return allRoutes, nil
}

func (s *routeService) Get(ctx context.Context, vpcID, routeID string) (*Route, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[Route](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v1/vpcs/%s/route_table/%s", vpcID, routeID),
		nil,
		nil,
	)
}

func (s *routeService) Create(ctx context.Context, vpcID string, req CreateRequest) (*CreateResponse, error) {
	if req.PortID == "" {
		return nil, fmt.Errorf("port_id cannot be empty")
	}
	if req.CIDRDestination == "" {
		return nil, fmt.Errorf("cidr_destination cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[CreateResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v1/vpcs/%s/route_table/routes", vpcID),
		req,
		nil,
	)
}

func (s *routeService) Delete(ctx context.Context, vpcID, routeID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v1/vpcs/%s/route_table/%s", vpcID, routeID),
		nil,
		nil,
	)
}

func validateSortValue(sort string) error {
	allowedSortFields := []string{"id", "port_id", "vpc_id", "description", "cidr_destination", "type", "status"}

	parts := strings.Split(sort, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid sort format, expected field:asc|desc")
	}

	field := strings.ToLower(parts[0])
	direction := strings.ToLower(parts[1])

	if !slices.Contains(allowedSortFields, field) {
		return fmt.Errorf(
			"invalid sort field: %q, allowed fields are: %s",
			field,
			strings.Join(allowedSortFields, ", "),
		)
	}

	if direction != "asc" && direction != "desc" {
		return fmt.Errorf("invalid sort direction, expected asc or desc")
	}

	return nil
}
