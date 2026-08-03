package tag

import (
	"context"
	"net/http"
	"net/url"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// ListTagValuesResponse represents the response of a tag value listing
	ListTagValuesResponse struct {
		Results []TagValue `json:"results"`
	}

	// ListTagValuesOptions defines the filters and pagination accepted when listing tag values
	ListTagValuesOptions struct {
		Name   *string
		Limit  *int
		Offset *int
		Sort   *string
	}

	// UpdateTagValueRequest represents the parameters for updating a tag value.
	// A nil Description clears the current one.
	UpdateTagValueRequest struct {
		Description *string `json:"description"`
	}
)

// TagValueService provides methods for managing the values of a tag.
type TagValueService interface {
	// List retrieves the values of a tag. The API returns at most 100 items per
	// call and defaults to 20, so use Limit and Offset to page through the results.
	List(ctx context.Context, tagName string, opts ListTagValuesOptions) ([]TagValue, error)
	// Get retrieves a value of a tag by name.
	Get(ctx context.Context, tagName, valueName string) (*TagValue, error)
	// Create adds a new value to an existing tag.
	Create(ctx context.Context, tagName string, req CreateTagValueRequest) (*TagValue, error)
	// Update changes the description of a value.
	Update(ctx context.Context, tagName, valueName string, req UpdateTagValueRequest) (*TagValue, error)
	// Delete removes a value from a tag.
	Delete(ctx context.Context, tagName, valueName string) error
}

// tagValueService implements the TagValueService interface
type tagValueService struct {
	client *TagClient
}

func (s *tagValueService) List(ctx context.Context, tagName string, opts ListTagValuesOptions) ([]TagValue, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}

	query := makeListQuery(listOptions{Limit: opts.Limit, Offset: opts.Offset, Sort: opts.Sort})
	if opts.Name != nil {
		query.Set("name", *opts.Name)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListTagValuesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		valuesPath(tagName),
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (s *tagValueService) Get(ctx context.Context, tagName, valueName string) (*TagValue, error) {
	if err := validateValueNames(tagName, valueName); err != nil {
		return nil, err
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValue](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		valuePath(tagName, valueName),
		nil,
		nil,
	)
}

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
		valuesPath(tagName),
		req,
		nil,
	)
}

func (s *tagValueService) Update(ctx context.Context, tagName, valueName string, req UpdateTagValueRequest) (*TagValue, error) {
	if err := validateValueNames(tagName, valueName); err != nil {
		return nil, err
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[TagValue](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		valuePath(tagName, valueName),
		req,
		nil,
	)
}

func (s *tagValueService) Delete(ctx context.Context, tagName, valueName string) error {
	if err := validateValueNames(tagName, valueName); err != nil {
		return err
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		valuePath(tagName, valueName),
		nil,
		nil,
	)
}

func validateValueNames(tagName, valueName string) error {
	if tagName == "" {
		return &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}
	if valueName == "" {
		return &client.ValidationError{Field: "value_name", Message: "cannot be empty"}
	}
	return nil
}

func valuesPath(tagName string) string {
	return tagPath(tagName) + "/values"
}

func valuePath(tagName, valueName string) string {
	return valuesPath(tagName) + "/" + url.PathEscape(valueName)
}
