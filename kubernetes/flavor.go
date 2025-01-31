package kubernetes

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	FlavorService interface {
		List(ctx context.Context, opts ListOptions) (*FlavorList, error)
	}

	FlavorList struct {
		Results []FlavorsAvailable `json:"results"`
	}

	FlavorsAvailable struct {
		NodePool     []FlavorWithSku `json:"nodepool"`
		ControlPlane []FlavorWithSku `json:"controlplane"`
	}

	FlavorWithSku struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		VCPU int    `json:"vcpu"`
		RAM  int    `json:"ram"`
		Size int    `json:"size"`
		SKU  string `json:"sku"`
	}

	flavorService struct {
		client *KubernetesClient
	}
)

func (s *flavorService) List(ctx context.Context, opts ListOptions) (*FlavorList, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/flavors", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if opts.Limit != nil {
		q.Add("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Add("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Add("_sort", *opts.Sort)
	}
	if len(opts.Expand) > 0 {
		q.Add("expand", strings.Join(opts.Expand, ","))
	}
	req.URL.RawQuery = q.Encode()

	var response FlavorList
	_, err = mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
