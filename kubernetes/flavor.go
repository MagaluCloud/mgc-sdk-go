package kubernetes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	FlavorService interface {
		List(ctx context.Context, opts ListOptions) (*[]FlavorsAvailable, error)
	}

	FlavorList struct {
		Results []FlavorsAvailable `json:"results"`
	}

	FlavorsAvailable struct {
		NodePool     []Flavor `json:"nodepool"`
		ControlPlane []Flavor `json:"controlplane"`
	}

	Flavor struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		VCPU int    `json:"vcpu"`
		RAM  int    `json:"ram"`
		Size int    `json:"size"`
	}

	flavorService struct {
		client *KubernetesClient
	}
)

func (s *flavorService) List(ctx context.Context, opts ListOptions) (*[]FlavorsAvailable, error) {
	query := url.Values{}
	if opts.Limit != nil {
		query.Add("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Add("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Add("_sort", *opts.Sort)
	}
	if len(opts.Expand) > 0 {
		query.Add("expand", strings.Join(opts.Expand, ","))
	}

	response, err := mgc_http.ExecuteSimpleRequestWithRespBody[FlavorList](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet, "/v1/flavors", nil, query)
	if err != nil {
		return nil, err
	}

	return &response.Results, nil
}
