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

type (
	// ListTagValuesResponse represents a list of tag values response
	ListTagValuesResponse struct {
		Results []TagValue `json:"results"`
	}

	// ListTagValuesOptions defines parameters for filtering and paginating tag values
	ListTagValuesOptions struct {
		Name   *string
		Limit  *int
		Offset *int
		Sort   *string
	}

	// UpdateTagValueRequest represents the parameters for updating a tag value
	UpdateTagValueRequest struct {
		Description *string `json:"description,omitempty"`
	}
)

// TagValueService provides methods for managing values within a tag.
type TagValueService interface {
	List(ctx context.Context, tagName string, opts ListTagValuesOptions) ([]TagValue, error)
	Get(ctx context.Context, tagName, valueName string) (*TagValue, error)
	Create(ctx context.Context, tagName string, req CreateTagValueRequest) (*TagValue, error)
	Update(ctx context.Context, tagName, valueName string, req UpdateTagValueRequest) (*TagValue, error)
	Delete(ctx context.Context, tagName, valueName string) (*TagValue, error)
}

// tagValueService implements TagValueService
type tagValueService struct {
	client *TagClient
}

// List returns all values for a given tag
func (s *tagValueService) List(ctx context.Context, tagName string, opts ListTagValuesOptions) ([]TagValue, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}

	query := make(url.Values)
	if opts.Name != nil {
		query.Set("name", *opts.Name)
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

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListTagValuesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/tags/%s/values", tagName),
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Get retrieves a specific value from a tag by name
func (s *tagValueService) Get(ctx context.Context, tagName, valueName string) (*TagValue, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return nil, &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValue](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/tags/%s/values/%s", tagName, valueName),
		nil,
		nil,
	)
}

// Create adds a new value to an existing tag
func (s *tagValueService) Create(ctx context.Context, tagName string, req CreateTagValueRequest) (*TagValue, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if req.Name == "" {
		return nil, &client.ValidationError{Field: "name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValue](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/tags/%s/values", tagName),
		req,
		nil,
	)
}

// Update modifies an existing tag value
func (s *tagValueService) Update(ctx context.Context, tagName, valueName string, req UpdateTagValueRequest) (*TagValue, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return nil, &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValue](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v0/tags/%s/values/%s", tagName, valueName),
		req,
		nil,
	)
}

// Delete removes a value from a tag
func (s *tagValueService) Delete(ctx context.Context, tagName, valueName string) (*TagValue, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return nil, &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValue](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/tags/%s/values/%s", tagName, valueName),
		nil,
		nil,
	)
}
