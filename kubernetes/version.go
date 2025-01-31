package kubernetes

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	VersionService interface {
		List(ctx context.Context) ([]Version, error)
	}

	VersionList struct {
		Results []Version `json:"results"`
	}

	Version struct {
		Version    string `json:"version"`
		Deprecated bool   `json:"deprecated"`
	}

	versionService struct {
		client *KubernetesClient
	}
)

func (s *versionService) List(ctx context.Context) ([]Version, error) {
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[VersionList](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodGet, "/v1/versions", nil, nil)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
