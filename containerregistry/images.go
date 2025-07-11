package containerregistry

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// ImagesService provides methods for managing images within repositories
	ImagesService interface {
		// List retrieves a list of images within a repository with optional filtering
		List(ctx context.Context, registryID, repositoryName string, opts ListOptions) (*ImagesResponse, error)
		// Delete removes an image from a repository by digest or tag
		Delete(ctx context.Context, registryID, repositoryName, digestOrTag string) error
		// Get retrieves a specific image from a repository by digest or tag
		Get(ctx context.Context, registryID, repositoryName, digestOrTag string) (*ImageResponse, error)
	}

	// ImagesResponse represents the response when listing images
	ImagesResponse struct {
		// Results is the list of images
		Results []ImageResponse `json:"results"`
	}

	// ImageResponse represents a container image
	ImageResponse struct {
		// Digest is the SHA256 digest of the image
		Digest string `json:"digest"`
		// SizeBytes is the size of the image in bytes
		SizeBytes int `json:"size_bytes"`
		// PushedAt is the timestamp when the image was pushed
		PushedAt string `json:"pushed_at"`
		// PulledAt is the timestamp when the image was last pulled
		PulledAt string `json:"pulled_at"`
		// ManifestMediaType is the media type of the image manifest
		ManifestMediaType string `json:"manifest_media_type"`
		// MediaType is the media type of the image
		MediaType string `json:"media_type"`
		// Tags is the list of tags associated with the image
		Tags []string `json:"tags"`
		// TagsDetails contains detailed information about each tag
		TagsDetails []ImageTagResponse `json:"tags_details"`
		// ExtraAttr contains additional attributes of the image
		ExtraAttr string `json:"extra_attr"`
	}

	// ImageTagResponse represents detailed information about an image tag
	ImageTagResponse struct {
		// Name is the name of the tag
		Name string `json:"name"`
		// PushedAt is the timestamp when the tag was pushed
		PushedAt string `json:"pushed_at"`
		// PulledAt is the timestamp when the tag was last pulled
		PulledAt string `json:"pulled_at"`
		// Signed indicates whether the tag is signed
		Signed bool `json:"signed"`
	}

	// imagesService implements the ImagesService interface
	imagesService struct {
		client *ContainerRegistryClient
	}
)

// List retrieves a list of images within a repository with optional filtering
func (c *imagesService) List(ctx context.Context, registryID, repositoryName string, opts ListOptions) (*ImagesResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s/images", registryID, repositoryName)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[ImagesResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
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
