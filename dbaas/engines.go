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
		List(ctx context.Context, opts ListEngineOptions) ([]EngineDetail, error)
		Get(ctx context.Context, id string) (*EngineDetail, error)
		ListEngineParameters(ctx context.Context, engineID string, opts ListEngineParametersOptions) ([]EngineParameterDetail, error)
	}

	// engineService implements the EngineService interface
	engineService struct {
		client *DBaaSClient
	}

	// ListEnginesResponse represents the response when listing engines
	ListEnginesResponse struct {
		Meta    MetaResponse   `json:"meta"`
		Results []EngineDetail `json:"results"`
	}

	// MetaResponse contains metadata about the response
	MetaResponse struct {
		Page    PageResponse       `json:"page"`
		Filters []FieldValueFilter `json:"filters"`
	}

	// PageResponse contains pagination details
	PageResponse struct {
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
		Count    int `json:"count"`
		Total    int `json:"total"`
		MaxLimit int `json:"max_limit"`
	}

	// FieldValueFilter represents a filter applied to the results
	FieldValueFilter struct {
		Field string `json:"field"`
		Value string `json:"value"`
	}

	// EngineDetail represents a database engine
	EngineDetail struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}

	// ListEngineOptions provides options for listing engines
	ListEngineOptions struct {
		Offset *int
		Limit  *int
		Status *string
	}

	// ListEngineParametersOptions provides options for listing engine parameters
	ListEngineParametersOptions struct {
		Offset     *int
		Limit      *int
		Dynamic    *bool
		Modifiable *bool
	}

	// EngineParametersResponse represents the response when listing engine parameters
	EngineParametersResponse struct {
		Results []EngineParameterDetail `json:"results"`
		Meta    MetaResponse            `json:"meta"`
	}

	// EngineParameterDetail represents a parameter of a database engine
	EngineParameterDetail struct {
		AllowedValues []string `json:"allowed_values"`
		DataType      string   `json:"data_type"`
		DefaultValue  string   `json:"default_value"`
		Description   string   `json:"description"`
		Dynamic       bool     `json:"dynamic"`
		EngineID      string   `json:"engine_id"`
		Modifiable    bool     `json:"modifiable"`
		Name          string   `json:"name"`
		ParameterName string   `json:"parameter_name"`
		RangedValue   bool     `json:"ranged_value"`
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
