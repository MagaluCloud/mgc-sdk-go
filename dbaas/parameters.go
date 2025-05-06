package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type ParameterCreateRequest struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type ParameterResponse struct {
	ID string `json:"id"`
}

type ParameterUpdateRequest struct {
	Value any `json:"value"`
}

type ParameterDetailResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type ParametersResponse struct {
	Meta    MetaResponse              `json:"meta"`
	Results []ParameterDetailResponse `json:"results"`
}

type ListParametersOptions struct {
	ParameterGroupID string
	Offset           *int
	Limit            *int
}

type ParameterService interface {
	List(ctx context.Context, opts ListParametersOptions) ([]ParameterDetailResponse, error)
	Create(ctx context.Context, groupID string, req ParameterCreateRequest) (*ParameterResponse, error)
	Update(ctx context.Context, groupID, parameterID string, req ParameterUpdateRequest) (*ParameterDetailResponse, error)
	Delete(ctx context.Context, groupID, parameterID string) error
}

type parameterService struct {
	client *DBaaSClient
}

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
