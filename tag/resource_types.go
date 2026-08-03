package tag

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	// ResourceType represents a resource type that supports tagging.
	ResourceType struct {
		Name      ResourceTypeName                `json:"name"`
		Product   Product                         `json:"product"`
		CreatedAt utils.LocalDateTimeWithoutZone  `json:"created_at"`
		UpdatedAt *utils.LocalDateTimeWithoutZone `json:"updated_at"`
	}

	// ListResourceTypesResponse represents the response of a resource type listing
	ListResourceTypesResponse struct {
		Results []ResourceType `json:"results"`
	}

	// ListResourceTypesOptions defines the filters and pagination accepted when
	// listing resource types.
	ListResourceTypesOptions struct {
		Name    *ResourceTypeName
		Product *Product
		Limit   *int
		Offset  *int
		Sort    *string
	}
)

// ResourceTypeService provides methods for listing the resource types that support tagging.
type ResourceTypeService interface {
	// List retrieves the resource types and the product that owns each one. The API
	// returns at most 100 items per call and defaults to 20, so use Limit and Offset
	// to page through the results.
	List(ctx context.Context, opts ListResourceTypesOptions) ([]ResourceType, error)
}

// resourceTypeService implements the ResourceTypeService interface
type resourceTypeService struct {
	client *TagClient
}

func (s *resourceTypeService) List(ctx context.Context, opts ListResourceTypesOptions) ([]ResourceType, error) {
	query := makeListQuery(listOptions{Limit: opts.Limit, Offset: opts.Offset, Sort: opts.Sort})

	if opts.Name != nil {
		query.Set("name", string(*opts.Name))
	}
	if opts.Product != nil {
		query.Set("product", string(*opts.Product))
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListResourceTypesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/resource-types",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}
