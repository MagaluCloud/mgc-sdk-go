package tag

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

var colorRegexp = regexp.MustCompile(`^[0-9a-f]{8}$`)

type (
	// TagValue represents a value associated with a tag
	TagValue struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		CreatedAt   string  `json:"created_at"`
		UpdatedAt   *string `json:"updated_at"`
	}

	// Tag represents a tag resource
	Tag struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Color       string     `json:"color"`
		Kinds       []string   `json:"kinds"`
		Values      []TagValue `json:"values"`
		CreatedAt   string     `json:"created_at"`
		UpdatedAt   *string    `json:"updated_at"`
	}

	// ListTagsResponse represents a list of tags response
	ListTagsResponse struct {
		Results []Tag `json:"results"`
	}

	// ListOptions defines parameters for filtering and paginating tags
	ListOptions struct {
		Name   *string
		Color  *string
		Limit  *int
		Offset *int
		Sort   *string
	}

	// CreateTagValueRequest represents a value to be created alongside a tag
	CreateTagValueRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// CreateTagRequest represents the parameters for creating a new tag
	CreateTagRequest struct {
		Name        string                  `json:"name"`
		Description string                  `json:"description"`
		Color       string                  `json:"color"`
		Kinds       []string                `json:"kinds"`
		Values      []CreateTagValueRequest `json:"values,omitempty"`
	}

	// UpdateTagRequest represents the parameters for updating a tag
	UpdateTagRequest struct {
		Color       *string  `json:"color,omitempty"`
		Description *string  `json:"description,omitempty"`
		Kinds       []string `json:"kinds,omitempty"`
	}
)

// TagService provides methods for managing tags.
// All operations in this service are performed against the global endpoint,
// as tags are not region-specific resources.
type TagService interface {
	List(ctx context.Context, opts ListOptions) ([]Tag, error)
	Get(ctx context.Context, tagName string) (*Tag, error)
	Create(ctx context.Context, req CreateTagRequest) (*Tag, error)
	Update(ctx context.Context, tagName string, req UpdateTagRequest) (*Tag, error)
	Delete(ctx context.Context, tagName string) (*Tag, error)
}

// tagService implements the TagService interface
type tagService struct {
	client *TagClient
}

// normalizeColor converts color to lowercase and validates the format.
// Color must be an 8-character lowercase hexadecimal string representing RGBA, without '#' prefix.
func normalizeColor(color string) (string, error) {
	lower := strings.ToLower(color)
	if colorRegexp.MatchString(lower) {
		return lower, nil
	}
	return "", &client.ValidationError{
		Field:   "color",
		Message: "must be an 8-character hexadecimal string representing RGBA (e.g. ff0000ff)",
	}
}

// List returns all tags for the tenant, with optional filters
func (s *tagService) List(ctx context.Context, opts ListOptions) ([]Tag, error) {
	query := make(url.Values)

	if opts.Name != nil {
		query.Set("name", *opts.Name)
	}
	if opts.Color != nil {
		normalized, err := normalizeColor(*opts.Color)
		if err != nil {
			return nil, err
		}
		query.Set("color", normalized)
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

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListTagsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/tags",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Get retrieves a specific tag by name
func (s *tagService) Get(ctx context.Context, tagName string) (*Tag, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Tag](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/tags/%s", tagName),
		nil,
		nil,
	)
}

// Create registers a new tag globally
func (s *tagService) Create(ctx context.Context, req CreateTagRequest) (*Tag, error) {
	if req.Name == "" {
		return nil, &client.ValidationError{Field: "name", Message: "cannot be empty"}
	}

	for i, v := range req.Values {
		if v.Name == "" {
			return nil, &client.ValidationError{
				Field:   fmt.Sprintf("values[%d].name", i),
				Message: "cannot be empty",
			}
		}
	}

	normalized, err := normalizeColor(req.Color)
	if err != nil {
		return nil, err
	}
	req.Color = normalized

	return mgc_http.ExecuteSimpleRequestWithRespBody[Tag](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v0/tags",
		req,
		nil,
	)
}

// Update modifies an existing tag by name
func (s *tagService) Update(ctx context.Context, tagName string, req UpdateTagRequest) (*Tag, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}

	if req.Color != nil {
		normalized, err := normalizeColor(*req.Color)
		if err != nil {
			return nil, err
		}
		req.Color = &normalized
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Tag](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v0/tags/%s", tagName),
		req,
		nil,
	)
}

// Delete removes a tag by name
func (s *tagService) Delete(ctx context.Context, tagName string) (*Tag, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "tag_name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Tag](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/tags/%s", tagName),
		nil,
		nil,
	)
}
