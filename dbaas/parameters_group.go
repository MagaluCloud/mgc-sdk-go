package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ParameterGroupType defines the type of parameter group
type ParameterGroupType string

const (
	ParameterGroupTypeSystem ParameterGroupType = "SYSTEM"
	ParameterGroupTypeUser   ParameterGroupType = "USER"

	ErrorIDEmpty        = "ID cannot be empty"
	PathParametersGroup = "/v2/parameter-groups"
)

type (
	// ListParameterGroupsOptions defines query parameters for listing parameter groups
	ListParameterGroupsOptions struct {
		Offset   *int
		Limit    *int
		Type     *ParameterGroupType
		EngineID *string
	}

	// ParameterGroupFilterOptions provides filtering options for ListAll (without pagination)
	ParameterGroupFilterOptions struct {
		Type     *ParameterGroupType
		EngineID *string
	}

	// ParameterGroupsResponse represents the API response for multiple parameter groups
	ParameterGroupsResponse struct {
		Meta    MetaResponse                   `json:"meta"`
		Results []ParameterGroupDetailResponse `json:"results"`
	}

	// ParameterGroupCreateRequest contains the data for creating a new parameter group
	ParameterGroupCreateRequest struct {
		Name        string  `json:"name"`
		EngineID    string  `json:"engine_id"`
		Description *string `json:"description,omitempty"`
	}

	// ParameterGroupResponse contains the ID of a newly created parameter group
	ParameterGroupResponse struct {
		ID string `json:"id"`
	}

	// ParameterGroupDetailResponse represents the detailed view of a parameter group
	ParameterGroupDetailResponse struct {
		ID          string             `json:"id"`
		Name        string             `json:"name"`
		Description *string            `json:"description,omitempty"`
		Type        ParameterGroupType `json:"type"`
		EngineID    string             `json:"engine_id"`
	}

	// ParameterGroupUpdateRequest contains the fields that can be updated for a parameter group
	ParameterGroupUpdateRequest struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// ParameterGroupService defines the interface for parameter group operations
	ParameterGroupService interface {
		List(ctx context.Context, opts ListParameterGroupsOptions) (*ParameterGroupsResponse, error)
		ListAll(ctx context.Context, filterOpts ParameterGroupFilterOptions) ([]ParameterGroupDetailResponse, error)
		Create(ctx context.Context, req ParameterGroupCreateRequest) (*ParameterGroupResponse, error)
		Get(ctx context.Context, ID string) (*ParameterGroupDetailResponse, error)
		Update(ctx context.Context, ID string, req ParameterGroupUpdateRequest) (*ParameterGroupDetailResponse, error)
		Delete(ctx context.Context, ID string) error
	}

	// parameterGroupService implements the ParameterGroupService interface
	parameterGroupService struct {
		client *DBaaSClient
	}
)

// List retrieves a list of parameter groups for the tenant.
func (s *parameterGroupService) List(ctx context.Context, opts ListParameterGroupsOptions) (*ParameterGroupsResponse, error) {
	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Type != nil {
		query.Set("type", string(*opts.Type))
	}
	if opts.EngineID != nil {
		query.Set("engine_id", *opts.EngineID)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ParameterGroupsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		PathParametersGroup,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListAll retrieves all parameter groups by fetching all pages with optional filtering
func (s *parameterGroupService) ListAll(ctx context.Context, filterOpts ParameterGroupFilterOptions) ([]ParameterGroupDetailResponse, error) {
	var allGroups []ParameterGroupDetailResponse
	offset := 0
	limit := 25

	for {
		currentOffset := offset
		currentLimit := limit
		opts := ListParameterGroupsOptions{
			Offset:   &currentOffset,
			Limit:    &currentLimit,
			Type:     filterOpts.Type,
			EngineID: filterOpts.EngineID,
		}

		resp, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allGroups = append(allGroups, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allGroups, nil
}

// CreateParameterGroup creates a new custom parameter group.
func (s *parameterGroupService) Create(ctx context.Context, req ParameterGroupCreateRequest) (*ParameterGroupResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterGroupResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		PathParametersGroup,
		req,
		nil,
	)
}

// Get retrieves details of a specific parameter group by its ID.
func (s *parameterGroupService) Get(ctx context.Context, ID string) (*ParameterGroupDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(ErrorIDEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterGroupDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(PathParametersGroup+"/%s", ID),
		nil,
		nil,
	)
}

// Update updates the name or description of a parameter group.
func (s *parameterGroupService) Update(ctx context.Context, ID string, req ParameterGroupUpdateRequest) (*ParameterGroupDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(ErrorIDEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterGroupDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(PathParametersGroup+"/%s", ID),
		req,
		nil,
	)
}

// Delete deletes a custom parameter group.
func (s *parameterGroupService) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf(ErrorIDEmpty)
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(PathParametersGroup+"/%s", ID),
		nil,
		nil,
	)
}
