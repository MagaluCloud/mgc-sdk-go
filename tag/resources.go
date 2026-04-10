package tag

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// RegionEnum represents a cloud region
type RegionEnum string

const (
	RegionBrSe1 RegionEnum = "br-se1"
	RegionBrNe1 RegionEnum = "br-ne1"
)

type (
	// TagValueResource represents the link between a tag value and a cloud resource
	TagValueResource struct {
		ResourceID     string       `json:"resource_id"`
		ResourceTypeID string       `json:"resource_type_id"`
		Region         []RegionEnum `json:"region"`
		CreatedAt      string       `json:"created_at"`
		UpdatedAt      *string      `json:"updated_at"`
	}

	// TagValueResourceListResponse represents a list of tag value resources
	TagValueResourceListResponse struct {
		Results []TagValueResource `json:"results"`
	}

	// CreateTagValueResourceRequest represents the parameters for linking a resource to a tag value
	CreateTagValueResourceRequest struct {
		ResourceID     string       `json:"resource_id"`
		ResourceTypeID string       `json:"resource_type_id"`
		Region         []RegionEnum `json:"region"`
	}

	// ListTagValueResourcesOptions defines parameters for filtering and paginating tag value resources
	ListTagValueResourcesOptions struct {
		ResourceTypeID *string
		Region         *RegionEnum
		ResourceID     *string
		Limit          *int
		Offset         *int
		Sort           *string
	}
)

// TagValueResourceService provides methods for managing the links between tag values and cloud resources.
type TagValueResourceService interface {
	List(ctx context.Context, tagName, valueName string, opts ListTagValueResourcesOptions) ([]TagValueResource, error)
	Get(ctx context.Context, tagName, valueName, resourceID string) (*TagValueResource, error)
	Create(ctx context.Context, tagName, valueName string, req CreateTagValueResourceRequest) (*TagValueResource, error)
	Delete(ctx context.Context, tagName, valueName, resourceID string) (*TagValueResource, error)
}

// tagValueResourceService implements TagValueResourceService
type tagValueResourceService struct {
	client *TagClient
}

// List returns all resources linked to a specific tag value
func (s *tagValueResourceService) List(ctx context.Context, tagName, valueName string, opts ListTagValueResourcesOptions) ([]TagValueResource, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return nil, &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}

	query := make(url.Values)
	if opts.ResourceTypeID != nil {
		query.Set("resource_type_id", *opts.ResourceTypeID)
	}
	if opts.Region != nil {
		query.Set("region", string(*opts.Region))
	}
	if opts.ResourceID != nil {
		query.Set("resource_id", *opts.ResourceID)
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[TagValueResourceListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/tags/%s/values/%s/resources", tagName, valueName),
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Get retrieves a specific resource link by resource ID
func (s *tagValueResourceService) Get(ctx context.Context, tagName, valueName, resourceID string) (*TagValueResource, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return nil, &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}
	if resourceID == "" {
		return nil, &client.ValidationError{Field: "resource_id", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValueResource](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/tags/%s/values/%s/resources/%s", tagName, valueName, resourceID),
		nil,
		nil,
	)
}

// Create links a cloud resource to a tag value
func (s *tagValueResourceService) Create(ctx context.Context, tagName, valueName string, req CreateTagValueResourceRequest) (*TagValueResource, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return nil, &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}
	if req.ResourceID == "" {
		return nil, &client.ValidationError{Field: "resource_id", Message: "cannot be empty"}
	}
	if req.ResourceTypeID == "" {
		return nil, &client.ValidationError{Field: "resource_type_id", Message: "cannot be empty"}
	}
	if len(req.Region) == 0 {
		return nil, &client.ValidationError{Field: "region", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValueResource](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/tags/%s/values/%s/resources", tagName, valueName),
		req,
		nil,
	)
}

// Delete removes the link between a cloud resource and a tag value
func (s *tagValueResourceService) Delete(ctx context.Context, tagName, valueName, resourceID string) (*TagValueResource, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return nil, &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}
	if resourceID == "" {
		return nil, &client.ValidationError{Field: "resource_id", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValueResource](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/tags/%s/values/%s/resources/%s", tagName, valueName, resourceID),
		nil,
		nil,
	)
}
