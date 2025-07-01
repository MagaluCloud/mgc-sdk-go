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
	ListInstanceTypesResponse struct {
		Meta    MetaResponse   `json:"meta"`
		Results []InstanceType `json:"results"`
	}

	InstanceType struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		Label             string `json:"label"`
		VCPU              string `json:"vcpu"`
		RAM               string `json:"ram"`
		FamilyDescription string `json:"family_description"`
		FamilySlug        string `json:"family_slug"`
		Size              string `json:"size"`
		CompatibleProduct string `json:"compatible_product"`
	}
)

type (
	InstanceTypeService interface {
		// List returns all available instance types
		List(ctx context.Context, opts ListInstanceTypeOptions) ([]InstanceType, error)

		// Get retrieves detailed information about a specific instance type
		Get(ctx context.Context, id string) (*InstanceType, error)
	}

	instanceTypeService struct {
		client *DBaaSClient
	}

	ListInstanceTypeOptions struct {
		Offset   *int    `json:"offset,omitempty"`
		Limit    *int    `json:"limit,omitempty"`
		Status   *string `json:"status,omitempty"`
		EngineID *string `json:"engine_id,omitempty"`
	}
)

func (s *instanceTypeService) List(ctx context.Context, opts ListInstanceTypeOptions) ([]InstanceType, error) {
	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Status != nil {
		query.Set("status", *opts.Status)
	}
	if opts.EngineID != nil {
		query.Set("engine_id", *opts.EngineID)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListInstanceTypesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v2/instance-types",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (s *instanceTypeService) Get(ctx context.Context, id string) (*InstanceType, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceType](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v2/instance-types/%s", id),
		nil,
		nil,
	)
}
