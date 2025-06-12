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
	EngineService interface {
		// List returns all available database engines
		List(ctx context.Context, opts ListEngineOptions) ([]EngineDetail, error)

		// Get retrieves detailed information about a specific engine
		Get(ctx context.Context, id string) (*EngineDetail, error)

		// List Engine Parameters retrieves parameters for a specific engine
		ListEngineParameters(ctx context.Context, engineID string, opts ListEngineParametersOptions) ([]EngineParameterDetail, error)
	}

	engineService struct {
		client *DBaaSClient
	}

	ListEnginesResponse struct {
		Meta    MetaResponse   `json:"meta"`
		Results []EngineDetail `json:"results"`
	}

	MetaResponse struct {
		Page    PageResponse       `json:"page"`
		Filters []FieldValueFilter `json:"filters"`
	}

	PageResponse struct {
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
		Count    int `json:"count"`
		Total    int `json:"total"`
		MaxLimit int `json:"max_limit"`
	}

	FieldValueFilter struct {
		Field string `json:"field"`
		Value string `json:"value"`
	}

	EngineDetail struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}

	ListEngineOptions struct {
		Offset *int
		Limit  *int
		Status *string
	}

	ListEngineParametersOptions struct {
		Offset     *int
		Limit      *int
		Dynamic    *bool
		Modifiable *bool
	}

	EngineParametersResponse struct {
		Results []EngineParameterDetail `json:"results"`
		Meta    MetaResponse            `json:"meta"`
	}

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
