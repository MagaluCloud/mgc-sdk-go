package kubernetes

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// FlavorService provides methods for managing Kubernetes flavors
	FlavorService interface {
		// List returns available flavors for node pools and control planes
		List(ctx context.Context, opts ListOptions) (*FlavorsAvailable, error)
	}

	// FlavorList represents the response when listing flavors
	FlavorList struct {
		// Results is the list of available flavors
		Results []FlavorsAvailable `json:"results"`
	}

	// FlavorsAvailable represents available flavors for different components
	FlavorsAvailable struct {
		// NodePool contains flavors available for node pools
		NodePool []Flavor `json:"nodepool"`
		// ControlPlane contains flavors available for control planes
		ControlPlane []Flavor `json:"controlplane"`
	}

	// Flavor represents a Kubernetes flavor (instance type)
	Flavor struct {
		// Name is the display name of the flavor
		Name string `json:"name"`
		// ID is the unique identifier of the flavor
		ID string `json:"id"`
		// VCPU is the number of virtual CPUs
		VCPU int `json:"vcpu"`
		// RAM is the amount of RAM in MB
		RAM int `json:"ram"`
		// Size is the size category of the flavor
		Size int `json:"size"`
	}

	// flavorService implements the FlavorService interface
	flavorService struct {
		client *KubernetesClient
	}
)

// List returns available flavors for node pools and control planes
func (s *flavorService) List(ctx context.Context, opts ListOptions) (*FlavorsAvailable, error) {
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

	if len(response.Results) == 0 {
		return nil, errors.New("no flavors available")
	}

	return &response.Results[0], nil
}
