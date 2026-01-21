package kubernetes

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// VersionService provides methods for managing Kubernetes versions
	VersionService interface {
		List(ctx context.Context, includeDeprecated bool) ([]Version, error)
	}

	// VersionList represents the response when listing versions
	VersionList struct {
		Results []Version `json:"results"`
	}

	// Version represents a Kubernetes version
	Version struct {
		Version    string `json:"version"`
		Deprecated bool   `json:"deprecated"`
	}

	// versionService implements the VersionService interface
	versionService struct {
		client *KubernetesClient
	}
)

// List returns all available Kubernetes versions
func (s *versionService) List(ctx context.Context, includeDeprecated bool) ([]Version, error) {
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[VersionList](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodGet, "/v1/versions", nil, nil)
	if err != nil {
		return nil, err
	}

	if includeDeprecated {
		return resp.Results, nil
	}

	filteredVersions := []Version{}

	for _, version := range resp.Results {
		if !version.Deprecated {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions, nil
}
