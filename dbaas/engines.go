package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// EngineService provides methods for managing database engines
	EngineService interface {
		// List returns all available database engines
		List(ctx context.Context, opts ListEngineOptions) ([]EngineDetail, error)

		// Get retrieves detailed information about a specific engine
		Get(ctx context.Context, id string) (*EngineDetail, error)

		// ListEngineParameters retrieves parameters for a specific engine
		ListEngineParameters(ctx context.Context, engineID string, opts ListEngineParametersOptions) ([]EngineParameterDetail, error)
	}

	// engineService implements the EngineService interface
	engineService struct {
		client *DBaaSClient
	}

	// ListEnginesResponse represents the response when listing engines
	ListEnginesResponse struct {
		// Meta contains pagination and filter information
		Meta MetaResponse `json:"meta"`
		// Results is the list of engines
		Results []EngineDetail `json:"results"`
	}

	// MetaResponse contains metadata about the response
	MetaResponse struct {
		// Page contains pagination information
		Page PageResponse `json:"page"`
		// Filters contains applied filters
		Filters []FieldValueFilter `json:"filters"`
	}

	// PageResponse contains pagination details
	PageResponse struct {
		// Offset is the number of items skipped
		Offset int `json:"offset"`
		// Limit is the maximum number of items returned
		Limit int `json:"limit"`
		// Count is the number of items in the current page
		Count int `json:"count"`
		// Total is the total number of items available
		Total int `json:"total"`
		// MaxLimit is the maximum allowed limit
		MaxLimit int `json:"max_limit"`
	}

	// FieldValueFilter represents a filter applied to the results
	FieldValueFilter struct {
		// Field is the field being filtered
		Field string `json:"field"`
		// Value is the filter value
		Value string `json:"value"`
	}

	// EngineDetail represents a database engine
	EngineDetail struct {
		// ID is the unique identifier of the engine
		ID string `json:"id"`
		// Name is the name of the engine
		Name string `json:"name"`
		// Version is the version of the engine
		Version string `json:"version"`
		// Status is the current status of the engine
		Status string `json:"status"`
	}

	// ListEngineOptions provides options for listing engines
	ListEngineOptions struct {
		// Offset is the number of engines to skip
		Offset *int
		// Limit is the maximum number of engines to return
		Limit *int
		// Status filters engines by status
		Status *string
	}

	// ListEngineParametersOptions provides options for listing engine parameters
	ListEngineParametersOptions struct {
		// Offset is the number of parameters to skip
		Offset *int
		// Limit is the maximum number of parameters to return
		Limit *int
		// Dynamic filters parameters by dynamic flag
		Dynamic *bool
		// Modifiable filters parameters by modifiable flag
		Modifiable *bool
	}

	// EngineParametersResponse represents the response when listing engine parameters
	EngineParametersResponse struct {
		// Results is the list of engine parameters
		Results []EngineParameterDetail `json:"results"`
		// Meta contains pagination and filter information
		Meta MetaResponse `json:"meta"`
	}

	// EngineParameterDetail represents a parameter of a database engine
	EngineParameterDetail struct {
		// AllowedValues contains the allowed values for this parameter
		AllowedValues []string `json:"allowed_values"`
		// DataType is the data type of the parameter
		DataType string `json:"data_type"`
		// DefaultValue is the default value of the parameter
		DefaultValue string `json:"default_value"`
		// Description is the description of the parameter
		Description string `json:"description"`
		// Dynamic indicates if the parameter can be changed dynamically
		Dynamic bool `json:"dynamic"`
		// EngineID is the ID of the engine this parameter belongs to
		EngineID string `json:"engine_id"`
		// Modifiable indicates if the parameter can be modified
		Modifiable bool `json:"modifiable"`
		// Name is the display name of the parameter
		Name string `json:"name"`
		// ParameterName is the internal name of the parameter
		ParameterName string `json:"parameter_name"`
		// RangedValue indicates if the parameter accepts a range of values
		RangedValue bool `json:"ranged_value"`
	}
)

// List returns all available database engines
func (s *engineService) List(ctx context.Context, opts ListEngineOptions) ([]EngineDetail, error) {
	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Status != nil {
		query.Set("status", string(*opts.Status))
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListEnginesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v2/engines",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

// Get retrieves detailed information about a specific engine
func (s *engineService) Get(ctx context.Context, id string) (*EngineDetail, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[EngineDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v2/engines/%s", id),
		nil,
		nil,
	)
}

// ListEngineParameters retrieves parameters for a specific engine
func (s *engineService) ListEngineParameters(ctx context.Context, engineID string, opts ListEngineParametersOptions) ([]EngineParameterDetail, error) {
	if engineID == "" {
		return nil, fmt.Errorf("engineID cannot be empty")
	}
	path := fmt.Sprintf("/v2/engines/%s/parameters", engineID)

	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Dynamic != nil {
		query.Set("dynamic", strconv.FormatBool(*opts.Dynamic))
	}
	if opts.Modifiable != nil {
		query.Set("modifiable", strconv.FormatBool(*opts.Modifiable))
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[EngineParametersResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}
