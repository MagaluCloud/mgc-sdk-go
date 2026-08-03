package tag

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

var colorRegexp = regexp.MustCompile(`^[0-9a-f]{6}$`)

type (
	// TagSummary identifies the tag that owns a value.
	TagSummary struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	// TagValue represents a value of a tag.
	TagValue struct {
		Name        string                          `json:"name"`
		Description *string                         `json:"description"`
		CreatedAt   utils.LocalDateTimeWithoutZone  `json:"created_at"`
		UpdatedAt   *utils.LocalDateTimeWithoutZone `json:"updated_at"`
		// Tag is nil when the value comes embedded in a Tag.
		Tag *TagSummary `json:"tag,omitempty"`
	}

	// Tag represents a tag and the values defined for it.
	Tag struct {
		Name        string                          `json:"name"`
		Description *string                         `json:"description"`
		Color       *string                         `json:"color"`
		Kinds       []TagKind                       `json:"kinds"`
		Values      []TagValue                      `json:"values"`
		CreatedAt   utils.LocalDateTimeWithoutZone  `json:"created_at"`
		UpdatedAt   *utils.LocalDateTimeWithoutZone `json:"updated_at"`
	}

	// ListTagsResponse represents the response of a tag listing
	ListTagsResponse struct {
		Results []Tag `json:"results"`
	}

	// ListTagsOptions defines the filters and pagination accepted when listing tags
	ListTagsOptions struct {
		Name   *string
		Color  *string
		Kinds  []TagKind
		Limit  *int
		Offset *int
		Sort   *string
	}

	// CreateTagValueRequest represents the parameters for creating a value of a tag
	CreateTagValueRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
	}

	// CreateTagRequest represents the parameters for creating a new tag
	CreateTagRequest struct {
		Name        string                  `json:"name"`
		Description *string                 `json:"description"`
		Color       *string                 `json:"color,omitempty"`
		Kinds       []TagKind               `json:"kinds,omitempty"`
		Values      []CreateTagValueRequest `json:"values,omitempty"`
	}

	// UpdateTagRequest represents the parameters for updating a tag.
	// Fields left nil are not sent, and the API keeps their current value.
	UpdateTagRequest struct {
		Description *string    `json:"description"`
		Color       *string    `json:"color,omitempty"`
		Kinds       *[]TagKind `json:"kinds,omitempty"`
	}
)

// TagService provides methods for managing tags.
// All operations in this service are performed against the global endpoint,
// as tags are not region-specific resources.
type TagService interface {
	// List retrieves the tags of the tenant. The API returns at most 100 items per
	// call and defaults to 20, so use Limit and Offset to page through the results.
	List(ctx context.Context, opts ListTagsOptions) ([]Tag, error)
	// Get retrieves a tag by name, along with its values.
	Get(ctx context.Context, tagName string) (*Tag, error)
	// Create registers a new tag, optionally with its first values.
	Create(ctx context.Context, req CreateTagRequest) (*Tag, error)
	// Update changes the description, color or kinds of a tag.
	Update(ctx context.Context, tagName string, req UpdateTagRequest) (*Tag, error)
	// Delete removes a tag and every value defined for it.
	Delete(ctx context.Context, tagName string) error
}

// tagService implements the TagService interface
type tagService struct {
	client *TagClient
}

// normalizeColor lowercases the color and validates it as a 6-digit hex RGB code
// without the '#' prefix.
func normalizeColor(color string) (string, error) {
	lower := strings.ToLower(color)
	if colorRegexp.MatchString(lower) {
		return lower, nil
	}
	return "", &client.ValidationError{
		Field:   "color",
		Message: "must be a 6-character hexadecimal string representing RGB (e.g. f54927)",
	}
}

func (s *tagService) List(ctx context.Context, opts ListTagsOptions) ([]Tag, error) {
	query := makeListQuery(listOptions{Limit: opts.Limit, Offset: opts.Offset, Sort: opts.Sort})

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
	for _, kind := range opts.Kinds {
		query.Add("kinds", string(kind))
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

func (s *tagService) Get(ctx context.Context, tagName string) (*Tag, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Tag](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		tagPath(tagName),
		nil,
		nil,
	)
}

func (s *tagService) Create(ctx context.Context, req CreateTagRequest) (*Tag, error) {
	if req.Name == "" {
		return nil, &client.ValidationError{Field: "name", Message: "cannot be empty"}
	}

	if req.Color != nil {
		normalized, err := normalizeColor(*req.Color)
		if err != nil {
			return nil, err
		}
		req.Color = &normalized
	}

	for i, value := range req.Values {
		if value.Name == "" {
			return nil, &client.ValidationError{
				Field:   fmt.Sprintf("values[%d].name", i),
				Message: "cannot be empty",
			}
		}
	}

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

func (s *tagService) Update(ctx context.Context, tagName string, req UpdateTagRequest) (*Tag, error) {
	if tagName == "" {
		return nil, &client.ValidationError{Field: "name", Message: "cannot be empty"}
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
		tagPath(tagName),
		req,
		nil,
	)
}

func (s *tagService) Delete(ctx context.Context, tagName string) error {
	if tagName == "" {
		return &client.ValidationError{Field: "name", Message: "cannot be empty"}
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		tagPath(tagName),
		nil,
		nil,
	)
}

// tagPath builds the path of a tag. Tag names accept spaces and brackets, so they
// have to be escaped.
func tagPath(tagName string) string {
	return "/v0/tags/" + url.PathEscape(tagName)
}
