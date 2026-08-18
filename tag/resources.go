package tag

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	// ResourceTypeSummary identifies the type of a tagged resource.
	ResourceTypeSummary struct {
		Name    ResourceTypeName `json:"name"`
		Product Product          `json:"product"`
	}

	// ResourceTag represents a tag attached to a resource and the value chosen for it.
	ResourceTag struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// Resource represents a resource of another product that can carry tags.
	// ExternalID is the ID the owning product uses for it, and is what identifies
	// the resource in every operation of this service.
	Resource struct {
		ID           string                          `json:"id"`
		ExternalID   string                          `json:"external_id"`
		ResourceType ResourceTypeSummary             `json:"resource_type"`
		Region       string                          `json:"region"`
		CreatedAt    utils.LocalDateTimeWithoutZone  `json:"created_at"`
		UpdatedAt    *utils.LocalDateTimeWithoutZone `json:"updated_at"`
		// LastTagAssociatedAt is nil while no tag was ever attached or detached.
		LastTagAssociatedAt *utils.LocalDateTimeWithoutZone `json:"last_tag_associated_at"`
		Tags                []ResourceTag                   `json:"tags"`
	}

	// ListResourcesResponse represents the response of a tagged resource listing
	ListResourcesResponse struct {
		Results []Resource `json:"results"`
	}

	// ListResourcesOptions defines the filters and pagination accepted when listing
	// tagged resources.
	ListResourcesOptions struct {
		ExternalID       *string
		ResourceTypeName *ResourceTypeName
		Region           *string
		Limit            *int
		Offset           *int
		Sort             *string
	}

	// AttachTag identifies a tag and the value to attach to a resource
	AttachTag struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// AttachTagsRequest represents the parameters for attaching tags to a resource
	AttachTagsRequest struct {
		Tags []AttachTag `json:"tags"`
	}
)

// ResourceService provides methods for managing the tags attached to the resources
// of other products.
type ResourceService interface {
	// List retrieves the resources that carry tags. The API returns at most 100 items
	// per call and defaults to 20, so use Limit and Offset to page through the results.
	List(ctx context.Context, opts ListResourcesOptions) ([]Resource, error)
	// Get retrieves a resource and the tags attached to it.
	Get(ctx context.Context, externalID string) (*Resource, error)
	// AttachTags attaches tags to a resource and returns the resource with every tag
	// it carries afterwards. A resource carries a single value per tag, so each tag
	// has to be paired with the value it takes.
	AttachTags(ctx context.Context, externalID string, req AttachTagsRequest) (*Resource, error)
	// DetachTag removes a tag, and the value it carried, from a resource.
	DetachTag(ctx context.Context, externalID, tagName string) error
}

// resourceService implements the ResourceService interface
type resourceService struct {
	client *TagClient
}

func (s *resourceService) List(ctx context.Context, opts ListResourcesOptions) ([]Resource, error) {
	query := makeListQuery(listOptions{Limit: opts.Limit, Offset: opts.Offset, Sort: opts.Sort})

	if opts.ExternalID != nil {
		query.Set("external_id", *opts.ExternalID)
	}
	if opts.ResourceTypeName != nil {
		query.Set("resource_type_name", string(*opts.ResourceTypeName))
	}
	if opts.Region != nil {
		query.Set("region", *opts.Region)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListResourcesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/resources",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (s *resourceService) Get(ctx context.Context, externalID string) (*Resource, error) {
	if externalID == "" {
		return nil, &client.ValidationError{Field: "external_id", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Resource](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		resourcePath(externalID),
		nil,
		nil,
	)
}

func (s *resourceService) AttachTags(ctx context.Context, externalID string, req AttachTagsRequest) (*Resource, error) {
	if externalID == "" {
		return nil, &client.ValidationError{Field: "external_id", Message: "cannot be empty"}
	}
	if len(req.Tags) == 0 {
		return nil, &client.ValidationError{Field: "tags", Message: "must have at least one tag"}
	}
	for i, tag := range req.Tags {
		if tag.Name == "" {
			return nil, &client.ValidationError{
				Field:   fmt.Sprintf("tags[%d].name", i),
				Message: "cannot be empty",
			}
		}
		if tag.Value == "" {
			return nil, &client.ValidationError{
				Field:   fmt.Sprintf("tags[%d].value", i),
				Message: "cannot be empty",
			}
		}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Resource](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		resourceTagsPath(externalID),
		req,
		nil,
	)
}

func (s *resourceService) DetachTag(ctx context.Context, externalID, tagName string) error {
	if externalID == "" {
		return &client.ValidationError{Field: "external_id", Message: "cannot be empty"}
	}
	if tagName == "" {
		return &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		resourceTagsPath(externalID)+"/"+url.PathEscape(tagName),
		nil,
		nil,
	)
}

// resourcePath builds the path of a resource. External IDs are minted by each
// product and may contain slashes, so they have to be escaped.
func resourcePath(externalID string) string {
	return "/v0/resources/" + url.PathEscape(externalID)
}

func resourceTagsPath(externalID string) string {
	return resourcePath(externalID) + "/tags"
}
