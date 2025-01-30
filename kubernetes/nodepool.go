package kubernetes

import (
	"context"
	"fmt"
	"net/http"
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

	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1alpha0/clusters/%s/node_pools", clusterID), nil)
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
	req.URL.RawQuery = q.Encode()

	var response struct {
		Results []NodePool `json:"results"`
	}
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func (s *nodePoolService) Create(ctx context.Context, clusterID string, req CreateNodePoolRequest) (*NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: "cannot be empty"}
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v0/clusters/%s/node_pools", clusterID),
		req)
	if err != nil {
		return nil, err
	}

	var nodePool NodePool
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &nodePool)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *nodePoolService) Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error) {
	if clusterID == "" || nodePoolID == "" {
		return nil, &client.ValidationError{Field: "clusterID/nodePoolID", Message: "cannot be empty"}
	}

	req, err := s.client.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/v0/clusters/%s/node_pools/%s", clusterID, nodePoolID),
		nil)
	if err != nil {
		return nil, err
	}

	var nodePool NodePool
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &nodePool)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *nodePoolService) Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error) {
	if clusterID == "" || nodePoolID == "" {
		return nil, &client.ValidationError{Field: "clusterID/nodePoolID", Message: "cannot be empty"}
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/v0/clusters/%s/node_pools/%s", clusterID, nodePoolID),
		req)
	if err != nil {
		return nil, err
	}

	var updatedNodePool NodePool
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &updatedNodePool)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *nodePoolService) Delete(ctx context.Context, clusterID, nodePoolID string) error {
	if clusterID == "" || nodePoolID == "" {
		return &client.ValidationError{Field: "clusterID/nodePoolID", Message: "cannot be empty"}
	}

	req, err := s.client.newRequest(ctx, http.MethodDelete,
		fmt.Sprintf("/v0/clusters/%s/node_pools/%s", clusterID, nodePoolID),
		nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	return err
}
