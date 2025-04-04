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
	// ParameterGroupTypeSystem represents a system parameter group
	ParameterGroupTypeSystem ParameterGroupType = "SYSTEM"
	// ParameterGroupTypeUser represents a user-defined parameter group
	ParameterGroupTypeUser ParameterGroupType = "USER"
)

type (
	// ListParameterGroupsOptions defines query parameters for listing parameter groups
	ListParameterGroupsOptions struct {
		Offset   *int
		Limit    *int
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
		Description string             `json:"description"`
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
		// ListParameterGroups retrieves a list of parameter groups for the tenant.
		List(ctx context.Context, opts ListParameterGroupsOptions) ([]ParameterGroupDetailResponse, error)

		// CreateParameterGroup creates a new custom parameter group.
		Create(ctx context.Context, req ParameterGroupCreateRequest) (*ParameterGroupResponse, error)

		// GetParameterGroup retrieves details of a specific parameter group by its ID.
		Get(ctx context.Context, ID string) (*ParameterGroupDetailResponse, error)

		// UpdateParameterGroup updates the name or description of a parameter group.
		Update(ctx context.Context, ID string, req ParameterGroupUpdateRequest) (*ParameterGroupDetailResponse, error)

		// DeleteParameterGroup deletes a custom parameter group.
		Delete(ctx context.Context, ID string) error
	}

	// parameterGroupService implements the ParameterGroupService interface
	parameterGroupService struct {
		client *DBaaSClient
	}
)

// List retrieves a list of parameter groups for the tenant.
func (s *parameterGroupService) List(ctx context.Context, opts ListParameterGroupsOptions) ([]ParameterGroupDetailResponse, error) {
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
		"/v2/parameter-groups",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

// CreateParameterGroup creates a new custom parameter group.
func (s *parameterGroupService) Create(ctx context.Context, req ParameterGroupCreateRequest) (*ParameterGroupResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterGroupResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v2/parameter-groups",
		req,
		nil,
	)
}

// Get retrieves details of a specific parameter group by its ID.
func (s *parameterGroupService) Get(ctx context.Context, ID string) (*ParameterGroupDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterGroupDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v2/parameter-groups/%s", ID),
		nil,
		nil,
	)
}

// Update updates the name or description of a parameter group.
func (s *parameterGroupService) Update(ctx context.Context, ID string, req ParameterGroupUpdateRequest) (*ParameterGroupDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterGroupDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v2/parameter-groups/%s", ID),
		req,
		nil,
	)
}

// Delete deletes a custom parameter group.
func (s *parameterGroupService) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v2/parameter-groups/%s", ID),
		nil,
		nil,
	)
}
