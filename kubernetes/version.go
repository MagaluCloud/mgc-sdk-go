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

	Version struct {
		Version    string `json:"version"`
		Deprecated bool   `json:"deprecated"`
	}

	versionService struct {
		client *KubernetesClient
	}
)

func (s *versionService) List(ctx context.Context) ([]Version, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/versions", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Results []Version `json:"results"`
	}
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
