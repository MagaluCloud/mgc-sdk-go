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
	Name  string `json:"name"`
	Value any    `json:"value"`
}

// ParameterResponse represents the response when creating a parameter
type ParameterResponse struct {
	ID string `json:"id"`
}

// ParameterUpdateRequest represents the request payload for updating a parameter
type ParameterUpdateRequest struct {
	Value any `json:"value"`
}

// ParameterDetailResponse represents detailed information about a parameter
type ParameterDetailResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}

// ParametersResponse represents the response when listing parameters
type ParametersResponse struct {
	Meta    MetaResponse              `json:"meta"`
	Results []ParameterDetailResponse `json:"results"`
}

// ListParametersOptions provides options for listing parameters
type ListParametersOptions struct {
	ParameterGroupID string
	Offset           *int
	Limit            *int
}

// ParameterFilterOptions provides filtering options for ListAll (without pagination)
type ParameterFilterOptions struct {
	ParameterGroupID string
}

// ParameterService provides methods for managing parameters within parameter groups
type ParameterService interface {
	List(ctx context.Context, opts ListParametersOptions) (*ParametersResponse, error)
	ListAll(ctx context.Context, filterOpts ParameterFilterOptions) ([]ParameterDetailResponse, error)
	Create(ctx context.Context, groupID string, req ParameterCreateRequest) (*ParameterResponse, error)
	Update(ctx context.Context, groupID, parameterID string, req ParameterUpdateRequest) (*ParameterDetailResponse, error)
	Delete(ctx context.Context, groupID, parameterID string) error
}

// parameterService implements the ParameterService interface
type parameterService struct {
	client *DBaaSClient
}

// List returns a list of parameters within a parameter group
func (s *parameterService) List(ctx context.Context, opts ListParametersOptions) (*ParametersResponse, error) {
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
	return resp, nil
}

// ListAll retrieves all parameters within a parameter group by fetching all pages
func (s *parameterService) ListAll(ctx context.Context, filterOpts ParameterFilterOptions) ([]ParameterDetailResponse, error) {
	var allParameters []ParameterDetailResponse
	offset := 0
	limit := 25

	for {
		currentOffset := offset
		currentLimit := limit
		opts := ListParametersOptions{
			ParameterGroupID: filterOpts.ParameterGroupID,
			Offset:           &currentOffset,
			Limit:            &currentLimit,
		}

		resp, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allParameters = append(allParameters, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allParameters, nil
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
