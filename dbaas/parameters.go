package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ParameterCreateRequest represents the request payload for creating a parameter
type ParameterCreateRequest struct {
	// Name is the name of the parameter
	Name string `json:"name"`
	// Value is the value of the parameter
	Value any `json:"value"`
}

// ParameterResponse represents the response when creating a parameter
type ParameterResponse struct {
	// ID is the unique identifier of the created parameter
	ID string `json:"id"`
}

// ParameterUpdateRequest represents the request payload for updating a parameter
type ParameterUpdateRequest struct {
	// Value is the new value of the parameter
	Value any `json:"value"`
}

// ParameterDetailResponse represents detailed information about a parameter
type ParameterDetailResponse struct {
	// ID is the unique identifier of the parameter
	ID string `json:"id"`
	// Name is the name of the parameter
	Name string `json:"name"`
	// Value is the current value of the parameter
	Value any `json:"value"`
}

// ParametersResponse represents the response when listing parameters
type ParametersResponse struct {
	// Meta contains pagination and filter information
	Meta MetaResponse `json:"meta"`
	// Results is the list of parameters
	Results []ParameterDetailResponse `json:"results"`
}

// ListParametersOptions provides options for listing parameters
type ListParametersOptions struct {
	// ParameterGroupID is the ID of the parameter group to list parameters from
	ParameterGroupID string
	// Offset is the number of parameters to skip
	Offset *int
	// Limit is the maximum number of parameters to return
	Limit *int
}

// ParameterService provides methods for managing parameters within parameter groups
type ParameterService interface {
	// List returns a list of parameters within a parameter group
	List(ctx context.Context, opts ListParametersOptions) ([]ParameterDetailResponse, error)
	// Create creates a new parameter within a parameter group
	Create(ctx context.Context, groupID string, req ParameterCreateRequest) (*ParameterResponse, error)
	// Update updates an existing parameter within a parameter group
	Update(ctx context.Context, groupID, parameterID string, req ParameterUpdateRequest) (*ParameterDetailResponse, error)
	// Delete removes a parameter from a parameter group
	Delete(ctx context.Context, groupID, parameterID string) error
}

// parameterService implements the ParameterService interface
type parameterService struct {
	client *DBaaSClient
}

// List returns a list of parameters within a parameter group
func (s *parameterService) List(ctx context.Context, opts ListParametersOptions) ([]ParameterDetailResponse, error) {
	q := make(url.Values)
	if opts.Offset != nil {
		q.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		q.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[ParametersResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v2/parameter-groups/%s/parameters", opts.ParameterGroupID),
		nil,
		q,
	)
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// Create creates a new parameter within a parameter group
func (s *parameterService) Create(ctx context.Context, groupID string, req ParameterCreateRequest) (*ParameterResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v2/parameter-groups/%s/parameters", groupID),
		req,
		nil,
	)
}

// Update updates an existing parameter within a parameter group
func (s *parameterService) Update(ctx context.Context, groupID, parameterID string, req ParameterUpdateRequest) (*ParameterDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ParameterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v2/parameter-groups/%s/parameters/%s", groupID, parameterID),
		req,
		nil,
	)
}

// Delete removes a parameter from a parameter group
func (s *parameterService) Delete(ctx context.Context, groupID, parameterID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v2/parameter-groups/%s/parameters/%s", groupID, parameterID),
		nil,
		nil,
	)
}
