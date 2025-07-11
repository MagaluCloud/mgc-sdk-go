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
	// ListInstanceTypesResponse represents the response when listing instance types
	ListInstanceTypesResponse struct {
		// Meta contains pagination and filter information
		Meta MetaResponse `json:"meta"`
		// Results is the list of instance types
		Results []InstanceType `json:"results"`
	}

	// InstanceType represents a database instance type
	InstanceType struct {
		// ID is the unique identifier of the instance type
		ID string `json:"id"`
		// Name is the name of the instance type
		Name string `json:"name"`
		// Label is the display label of the instance type
		Label string `json:"label"`
		// VCPU is the number of virtual CPUs
		VCPU string `json:"vcpu"`
		// RAM is the amount of RAM in GB
		RAM string `json:"ram"`
		// FamilyDescription is the description of the instance family
		FamilyDescription string `json:"family_description"`
		// FamilySlug is the slug identifier of the instance family
		FamilySlug string `json:"family_slug"`
		// Size is the size category of the instance type
		Size string `json:"size"`
		// CompatibleProduct is the compatible database product
		CompatibleProduct string `json:"compatible_product"`
	}
)

type (
	// InstanceTypeService provides methods for managing database instance types
	InstanceTypeService interface {
		// List returns all available instance types
		List(ctx context.Context, opts ListInstanceTypeOptions) ([]InstanceType, error)

		// Get retrieves detailed information about a specific instance type
		Get(ctx context.Context, id string) (*InstanceType, error)
	}

	// instanceTypeService implements the InstanceTypeService interface
	instanceTypeService struct {
		client *DBaaSClient
	}

	// ListInstanceTypeOptions provides options for listing instance types
	ListInstanceTypeOptions struct {
		// Offset is the number of instance types to skip
		Offset *int `json:"offset,omitempty"`
		// Limit is the maximum number of instance types to return
		Limit *int `json:"limit,omitempty"`
		// Status filters instance types by status
		Status *string `json:"status,omitempty"`
		// EngineID filters instance types by compatible engine
		EngineID *string `json:"engine_id,omitempty"`
	}
)

// List returns all available instance types
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

// Get retrieves detailed information about a specific instance type
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
