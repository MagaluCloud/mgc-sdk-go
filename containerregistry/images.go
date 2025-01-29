package containerregistry

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	ImagesService interface {
		List(ctx context.Context, registryID, repositoryName string, opts ListOptions) (*ImagesResponse, error)
		Delete(ctx context.Context, registryID, repositoryName, digestOrTag string) error
		Get(ctx context.Context, registryID, repositoryName, digestOrTag string) (*ImageResponse, error)
	}

	ImagesResponse struct {
		Results []ImageResponse `json:"results"`
	}

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

	ImageTagResponse struct {
		Name     string `json:"name"`
		PushedAt string `json:"pushed_at"`
		PulledAt string `json:"pulled_at"`
		Signed   bool   `json:"signed"`
	}

	imagesService struct {
		client *ContainerRegistryClient
	}
)

func (c *imagesService) List(ctx context.Context, registryID, repositoryName string, opts ListOptions) (*ImagesResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s/images", registryID, repositoryName)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[ImagesResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *imagesService) Delete(ctx context.Context, registryID, repositoryName, digestOrTag string) error {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s/images/%s", registryID, repositoryName, digestOrTag)

	err := mgc_http.ExecuteSimpleRequest(ctx, c.client.newRequest, c.client.GetConfig(), http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *imagesService) Get(ctx context.Context, registryID, repositoryName, digestOrTag string) (*ImageResponse, error) {
	path := fmt.Sprintf("/v0/registries/%s/repositories/%s/images/%s", registryID, repositoryName, digestOrTag)

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[ImageResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
