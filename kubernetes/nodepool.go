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
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

const (
	nodePoolIdField    = "nodePoolID"
	clusterIdField     = "clusterID"
	clusterNodepoolURL = "/v0/clusters/%s/node_pools/%s"
)

type (
	// NodePoolService provides methods for managing Kubernetes node pools
	NodePoolService interface {
		Nodes(ctx context.Context, clusterID, nodePoolID string) ([]NodeResponse, error)
		List(ctx context.Context, clusterID string, opts ListOptions) ([]NodePool, error)
		Create(ctx context.Context, clusterID string, req CreateNodePoolRequest) (*NodePool, error)
		Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error)
		Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error)
		Delete(ctx context.Context, clusterID, nodePoolID string) error
	}

	// NodePoolList represents the response when listing node pools
	NodePoolList struct {
		Results []NodePool `json:"results"`
	}

	// NodeAddress represents network addresses
	NodeAddress struct {
		Address string `json:"address"`
		Type    string `json:"type"`
	}

	// NodeResources represents node resources (used for both allocatable and capacity)
	NodeResources struct {
		CPU              string `json:"cpu"`
		EphemeralStorage string `json:"ephemeral_storage"`
		Hugepages1Gi     string `json:"hugepages_1Gi"`
		Hugepages2Mi     string `json:"hugepages_2Mi"`
		Memory           string `json:"memory"`
		Pods             string `json:"pods"`
	}

	// NodeInfrastructure represents node infrastructure information
	NodeInfrastructure struct {
		Architecture            string        `json:"architecture"`
		ContainerRuntimeVersion string        `json:"containerRuntimeVersion"`
		KernelVersion           string        `json:"kernelVersion"`
		KubeProxyVersion        string        `json:"kubeProxyVersion"`
		KubeletVersion          string        `json:"kubeletVersion"`
		OperatingSystem         string        `json:"operatingSystem"`
		OsImage                 string        `json:"osImage"`
		Allocatable             NodeResources `json:"allocatable"`
		Capacity                NodeResources `json:"capacity"`
	}

	// NodeResponse represents a Kubernetes node
	NodeResponse struct {
		ID             string             `json:"id"`
		Name           string             `json:"name"`
		Namespace      string             `json:"namespace"`
		ClusterName    string             `json:"cluster_name"`
		NodepoolName   string             `json:"nodepool_name"`
		CreatedAt      time.Time          `json:"created_at"`
		Annotations    map[string]string  `json:"annotations"`
		Labels         map[string]string  `json:"labels"`
		Taints         *[]Taint           `json:"taints,omitempty"`
		Addresses      []NodeAddress      `json:"addresses"`
		Flavor         string             `json:"flavor"`
		Infrastructure NodeInfrastructure `json:"infrastructure"`
		Status         MessageState       `json:"status"`
	}

	// NodesResponse represents the response when listing nodes
	NodesResponse struct {
		Results []NodeResponse `json:"results"`
	}

	// PatchNodePoolRequest represents the request payload for updating a node pool
	PatchNodePoolRequest struct {
		Replicas  *int       `json:"replicas,omitempty"`
		AutoScale *AutoScale `json:"auto_scale,omitempty"`
	}

	// nodePoolService implements the NodePoolService interface
	nodePoolService struct {
		client *KubernetesClient
	}
)

// Nodes returns a list of nodes in a specific node pool
func (s *nodePoolService) Nodes(ctx context.Context, clusterID, nodePoolID string) ([]NodeResponse, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: clusterIdField, Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return nil, &client.ValidationError{Field: nodePoolIdField, Message: utils.CannotBeEmpty}
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[NodesResponse](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodGet,
		fmt.Sprintf(clusterNodepoolURL+"/nodes", clusterID, nodePoolID), nil, nil)

	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

// List returns a list of node pools in a cluster with optional filtering and pagination
func (s *nodePoolService) List(ctx context.Context, clusterID string, opts ListOptions) ([]NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: clusterIdField, Message: utils.CannotBeEmpty}
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

// Create creates a new node pool in a cluster
func (s *nodePoolService) Create(ctx context.Context, clusterID string, req CreateNodePoolRequest) (*NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: clusterIdField, Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodPost,
		fmt.Sprintf("/v0/clusters/%s/node_pools", clusterID), req, nil)
}

// Get retrieves detailed information about a specific node pool
func (s *nodePoolService) Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: clusterIdField, Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return nil, &client.ValidationError{Field: nodePoolIdField, Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodGet,
		fmt.Sprintf(clusterNodepoolURL, clusterID, nodePoolID), nil, nil)
}

// Update updates a node pool's properties
func (s *nodePoolService) Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: clusterIdField, Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return nil, &client.ValidationError{Field: nodePoolIdField, Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodPatch,
		fmt.Sprintf(clusterNodepoolURL, clusterID, nodePoolID), req, nil)
}

// Delete removes a node pool from a cluster
func (s *nodePoolService) Delete(ctx context.Context, clusterID, nodePoolID string) error {
	if clusterID == "" {
		return &client.ValidationError{Field: clusterIdField, Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return &client.ValidationError{Field: nodePoolIdField, Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodDelete,
		fmt.Sprintf(clusterNodepoolURL, clusterID, nodePoolID), nil, nil)
}
