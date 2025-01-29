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
		ID      string       `json:"id"`
		Name    string       `json:"name"`
		Version string       `json:"version"`
		Status  EngineStatus `json:"status"`
	}
)

// EngineStatus represents the status of a database engine
type EngineStatus string

const (
	EngineStatusActive     EngineStatus = "ACTIVE"
	EngineStatusDeprecated EngineStatus = "DEPRECATED"
)

type (
	EngineService interface {
		// List returns all available database engines
		List(ctx context.Context, opts ListEngineOptions) ([]EngineDetail, error)

		// Get retrieves detailed information about a specific engine
		Get(ctx context.Context, id string) (*EngineDetail, error)
	}

	engineService struct {
		client *DBaaSClient
	}

	ListEngineOptions struct {
		Offset *int
		Limit  *int
		Status *EngineStatus
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
		"/v1/engines",
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
		fmt.Sprintf("/v1/engines/%s", id),
		nil,
		nil,
	)
}
