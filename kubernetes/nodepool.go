package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	ListOptions struct {
		Limit  *int
		Offset *int
		Sort   *string
		Expand []string
	}

	NodePoolService interface {
		List(ctx context.Context, clusterID string, opts ListOptions) ([]NodePool, error)
		Create(ctx context.Context, clusterID string, req CreateNodePoolRequest) (*NodePool, error)
		Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error)
		Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error)
		Delete(ctx context.Context, clusterID, nodePoolID string) error
	}

	NodePoolList struct {
		Results []NodePool `json:"results"`
	}

	NodePool struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		Flavor    string     `json:"flavor"`
		Replicas  int        `json:"replicas"`
		Tags      []string   `json:"tags"`
		Taints    []Taint    `json:"taints"`
		AutoScale *AutoScale `json:"auto_scale,omitempty"`
		Status    Status     `json:"status"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
	}

	CreateNodePoolRequest struct {
		Name      string     `json:"name"`
		Flavor    string     `json:"flavor"`
		Replicas  int        `json:"replicas"`
		Tags      []string   `json:"tags,omitempty"`
		Taints    []Taint    `json:"taints,omitempty"`
		AutoScale *AutoScale `json:"auto_scale,omitempty"`
	}

	PatchNodePoolRequest struct {
		Replicas  *int       `json:"replicas,omitempty"`
		AutoScale *AutoScale `json:"auto_scale,omitempty"`
	}

	Taint struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Effect string `json:"effect"`
	}

	AutoScale struct {
		MinReplicas int `json:"min_replicas"`
		MaxReplicas int `json:"max_replicas"`
	}

	nodePoolService struct {
		client *KubernetesClient
	}
)

func (s *nodePoolService) List(ctx context.Context, clusterID string, opts ListOptions) ([]NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: "cannot be empty"}
	}

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

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[NodePoolList](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodGet, fmt.Sprintf("/v1alpha0/clusters/%s/node-pools", clusterID), nil, query)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func (s *nodePoolService) Create(ctx context.Context, clusterID string, req CreateNodePoolRequest) (*NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: "cannot be empty"}
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodPost,
		fmt.Sprintf("/v0/clusters/%s/node_pools", clusterID), req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *nodePoolService) Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error) {
	if clusterID == "" || nodePoolID == "" {
		return nil, &client.ValidationError{Field: "clusterID/nodePoolID", Message: "cannot be empty"}
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet,
		fmt.Sprintf("/v0/clusters/%s/node_pools/%s", clusterID, nodePoolID), nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *nodePoolService) Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error) {
	if clusterID == "" || nodePoolID == "" {
		return nil, &client.ValidationError{Field: "clusterID/nodePoolID", Message: "cannot be empty"}
	}
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodPatch,
		fmt.Sprintf("/v0/clusters/%s/node_pools/%s", clusterID, nodePoolID), req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *nodePoolService) Delete(ctx context.Context, clusterID, nodePoolID string) error {
	if clusterID == "" || nodePoolID == "" {
		return &client.ValidationError{Field: "clusterID/nodePoolID", Message: "cannot be empty"}
	}

	err := mgc_http.ExecuteSimpleRequest(ctx, s.client.newRequest, s.client.GetConfig(), http.MethodDelete,
		fmt.Sprintf("/v0/clusters/%s/node_pools/%s", clusterID, nodePoolID), nil, nil)
	if err != nil {
		return err
	}

	return nil
}
