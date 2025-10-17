package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ImageExpand represents the expand options for image responses.
type ImageExpand string

// Constants for expanding related resources in image responses.
const (
	ImageTagsDetailsExpand       ImageExpand = "tags_details"
	ImageExtraAttrExpand         ImageExpand = "extra_attr"
	ImageManifestMediaTypeExpand ImageExpand = "manifest_media_type"
	ImageMediaTypeExpand         ImageExpand = "media_type"
)

type (
	// ImagesService provides methods for managing images within repositories
	ImagesService interface {
		List(ctx context.Context, registryID, repositoryName string, opts ImageListOptions) (*ImagesResponse, error)
		ListAll(ctx context.Context, registryID, repositoryName string, filterOpts ImageFilterOptions) ([]ImageResponse, error)
		Delete(ctx context.Context, registryID, repositoryName, digestOrTag string) error
		Get(ctx context.Context, registryID, repositoryName, digestOrTag string) (*ImageResponse, error)
	}

	// ImageListOptions provides options for listing images with pagination
	ImageListOptions struct {
		Offset *int
		Limit  *int
		ImageFilterOptions
	}

	// ImageFilterOptions provides filtering options for images
	ImageFilterOptions struct {
		Sort   *string
		Expand []ImageExpand
	}

	// ImagesResponse represents the response when listing images
	ImagesResponse = helpers.PaginatedResponse[ImageResponse]

	// ImageResponse represents a container image
	ImageResponse struct {
		Digest            string             `json:"digest"`
		SizeBytes         int                `json:"size_bytes"`
		PushedAt          string             `json:"pushed_at"`
		PulledAt          string             `json:"pulled_at"`
		ManifestMediaType string             `json:"manifest_media_type"`
		MediaType         string             `json:"media_type"`
		Tags              []string           `json:"tags"`
		TagsDetails       []ImageTagResponse `json:"tags_details"`
		ExtraAttr         string             `json:"extra_attr"`
	}

	// ImageTagResponse represents detailed information about an image tag
	ImageTagResponse struct {
		Name     string `json:"name"`
		PushedAt string `json:"pushed_at"`
		PulledAt string `json:"pulled_at"`
		Signed   bool   `json:"signed"`
	}

	// imagesService implements the ImagesService interface
	imagesService struct {
		client *ContainerRegistryClient
	}
)

// List retrieves a list of images within a repository with optional filtering and pagination
func (c *imagesService) List(ctx context.Context, registryID, repositoryName string, opts ImageListOptions) (*ImagesResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s/images", registryID, repositoryName)
	query := c.createImageQueryParams(opts)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[ImagesResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListAll retrieves all images across all pages with optional filtering
func (c *imagesService) ListAll(ctx context.Context, registryID, repositoryName string, filterOpts ImageFilterOptions) ([]ImageResponse, error) {
	var allImages []ImageResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		opts := ImageListOptions{
			Offset:             &currentOffset,
			Limit:              &currentLimit,
			ImageFilterOptions: filterOpts,
		}

		result, err := c.List(ctx, registryID, repositoryName, opts)
		if err != nil {
			return nil, err
		}

		allImages = append(allImages, result.Results...)

		if len(result.Results) < limit {
			break
		}

		offset += limit
	}

	return allImages, nil
}

// createImageQueryParams creates URL query parameters from ImageListOptions
func (c *imagesService) createImageQueryParams(opts ImageListOptions) url.Values {
	query := make(url.Values)

	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}
	if len(opts.Expand) > 0 {
		expandStrs := make([]string, len(opts.Expand))
		for i, exp := range opts.Expand {
			expandStrs[i] = string(exp)
		}
		query.Set("_expand", strings.Join(expandStrs, ","))
	}

	return query
}

// Delete removes an image from a repository by digest or tag
func (c *imagesService) Delete(ctx context.Context, registryID, repositoryName, digestOrTag string) error {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s/images/%s", registryID, repositoryName, digestOrTag)

	err := mgc_http.ExecuteSimpleRequest(ctx, c.client.newRequest, c.client.GetConfig(), http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a specific image from a repository by digest or tag
func (c *imagesService) Get(ctx context.Context, registryID, repositoryName, digestOrTag string) (*ImageResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s/images/%s", registryID, repositoryName, digestOrTag)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[ImageResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
