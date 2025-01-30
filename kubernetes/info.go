package kubernetes

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	InfoService interface {
		ListFlavors(ctx context.Context) (*FlavorsAvailable, error)
		ListVersions(ctx context.Context) ([]Version, error)
	}

	infoService struct {
		client *KubernetesClient
	}
)

func (s *infoService) ListFlavors(ctx context.Context) (*FlavorsAvailable, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v0/info/flavors", nil)
	if err != nil {
		return nil, err
	}

	var flavors FlavorsAvailable
	_, err = mgc_http.Do(s.client.GetConfig(), ctx, req, &flavors)
	if err != nil {
		return nil, fmt.Errorf("failed to list flavors: %w", err)
	}

	return &flavors, nil
}

func (s *infoService) ListVersions(ctx context.Context) ([]Version, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v0/info/versions", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Results []Version `json:"results"`
	}
	_, err = mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	return response.Results, nil
}
