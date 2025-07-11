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
	// nodePoolIdField is the field name for node pool ID validation
	nodePoolIdField = "nodePoolID"
	// clusterIdField is the field name for cluster ID validation
	clusterIdField = "clusterID"
	// clusterNodepoolURL is the URL template for cluster node pool operations
	clusterNodepoolURL = "/v0/clusters/%s/node_pools/%s"
)

type (
	// ListOptions provides options for listing resources
	ListOptions struct {
		// Limit is the maximum number of items to return
		Limit *int
		// Offset is the number of items to skip
		Offset *int
		// Sort is the field to sort by
		Sort *string
		// Expand contains fields to expand in the response
		Expand []string
	}

	// NodePoolService provides methods for managing Kubernetes node pools
	NodePoolService interface {
		// Nodes returns a list of nodes in a specific node pool
		Nodes(ctx context.Context, clusterID, nodePoolID string) ([]Node, error)
		// List returns a list of node pools in a cluster with optional filtering and pagination
		List(ctx context.Context, clusterID string, opts ListOptions) ([]NodePool, error)
		// Create creates a new node pool in a cluster
		Create(ctx context.Context, clusterID string, req CreateNodePoolRequest) (*NodePool, error)
		// Get retrieves detailed information about a specific node pool
		Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error)
		// Update updates a node pool's properties
		Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error)
		// Delete removes a node pool from a cluster
		Delete(ctx context.Context, clusterID, nodePoolID string) error
	}

	// NodePoolList represents the response when listing node pools
	NodePoolList struct {
		// Results is the list of node pools
		Results []NodePool `json:"results"`
	}

	// InstanceTemplate represents the template for node instances
	InstanceTemplate struct {
		// Flavor contains the flavor configuration
		Flavor Flavor `json:"flavor"`
		// NodeImage is the image used for nodes
		NodeImage string `json:"node_image"`
		// DiskSize is the disk size in GB
		DiskSize int `json:"disk_size"`
		// DiskType is the type of disk
		DiskType string `json:"disk_type"`
	}

	// NodePool represents a Kubernetes node pool
	NodePool struct {
		// ID is the unique identifier of the node pool
		ID string `json:"id"`
		// Name is the name of the node pool
		Name string `json:"name"`
		// InstanceTemplate contains the instance template configuration
		InstanceTemplate InstanceTemplate `json:"instance_template"`
		// Replicas is the number of replicas
		Replicas int `json:"replicas"`
		// Zone contains availability zones (optional)
		Zone *[]string `json:"zone,omitempty"`
		// Tags contains node tags (optional)
		Tags *[]string `json:"tags"`
		// Labels contains node labels (optional)
		Labels map[string]string `json:"labels,omitempty"`
		// Taints contains node taints (optional)
		Taints *[]Taint `json:"taints,omitempty"`
		// SecurityGroups contains security group IDs (optional)
		SecurityGroups *[]string `json:"security_groups,omitempty"`
		// CreatedAt is the timestamp when the node pool was created
		CreatedAt *time.Time `json:"created_at"`
		// UpdatedAt is the timestamp when the node pool was last updated (optional)
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
		// AutoScale contains autoscaling configuration (optional)
		AutoScale *AutoScale `json:"auto_scale,omitempty"`
		// Status contains the node pool status
		Status Status `json:"status"`
		// Flavor is the flavor identifier
		Flavor string `json:"flavor"`
		// MaxPodsPerNode is the maximum number of pods per node (optional)
		MaxPodsPerNode *int `json:"max_pods_per_node,omitempty"`
		// AvailabilityZones contains availability zones (optional)
		AvailabilityZones *[]string `json:"availability_zones,omitempty"`
	}

	// Addresses represents network addresses
	Addresses struct {
		// Address is the IP address
		Address string `json:"address"`
		// Type is the address type
		Type string `json:"type"`
	}

	// Allocatable represents allocatable resources
	Allocatable struct {
		// CPU is the allocatable CPU
		CPU string `json:"cpu"`
		// EphemeralStorage is the allocatable ephemeral storage
		EphemeralStorage string `json:"ephemeral_storage"`
		// Hugepages1Gi is the allocatable 1Gi hugepages
		Hugepages1Gi string `json:"hugepages_1Gi"`
		// Hugepages2Mi is the allocatable 2Mi hugepages
		Hugepages2Mi string `json:"hugepages_2Mi"`
		// Memory is the allocatable memory
		Memory string `json:"memory"`
		// Pods is the allocatable number of pods
		Pods string `json:"pods"`
	}

	// Capacity represents total capacity
	Capacity struct {
		// CPU is the total CPU
		CPU string `json:"cpu"`
		// EphemeralStorage is the total ephemeral storage
		EphemeralStorage string `json:"ephemeral_storage"`
		// Hugepages1Gi is the total 1Gi hugepages
		Hugepages1Gi string `json:"hugepages_1Gi"`
		// Hugepages2Mi is the total 2Mi hugepages
		Hugepages2Mi string `json:"hugepages_2Mi"`
		// Memory is the total memory
		Memory string `json:"memory"`
		// Pods is the total number of pods
		Pods string `json:"pods"`
	}

	// Infrastructure represents node infrastructure information
	Infrastructure struct {
		// Allocatable contains allocatable resources
		Allocatable Allocatable `json:"allocatable"`
		// Architecture is the node architecture
		Architecture string `json:"architecture"`
		// Capacity contains total capacity
		Capacity Capacity `json:"capacity"`
		// ContainerRuntimeVersion is the container runtime version
		ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
		// KernelVersion is the kernel version
		KernelVersion string `json:"kernelVersion"`
		// KubeProxyVersion is the kube-proxy version
		KubeProxyVersion string `json:"kubeProxyVersion"`
		// KubeletVersion is the kubelet version
		KubeletVersion string `json:"kubeletVersion"`
		// OperatingSystem is the operating system
		OperatingSystem string `json:"operatingSystem"`
		// OsImage is the OS image
		OsImage string `json:"osImage"`
	}

	// Node represents a Kubernetes node
	Node struct {
		// Addresses contains network addresses
		Addresses []Addresses `json:"addresses"`
		// Annotations contains node annotations
		Annotations map[string]string `json:"annotations"`
		// ClusterName is the name of the cluster
		ClusterName string `json:"cluster_name"`
		// CreatedAt is the timestamp when the node was created
		CreatedAt time.Time `json:"created_at"`
		// Flavor is the flavor identifier
		Flavor string `json:"flavor"`
		// ID is the unique identifier of the node
		ID string `json:"id"`
		// Infrastructure contains infrastructure information
		Infrastructure Infrastructure `json:"infrastructure"`
		// Labels contains node labels
		Labels map[string]string `json:"labels"`
		// Name is the name of the node
		Name string `json:"name"`
		// Namespace is the namespace
		Namespace string `json:"namespace"`
		// NodeImage is the node image
		NodeImage string `json:"node_image"`
		// NodepoolName is the name of the node pool
		NodepoolName string `json:"nodepool_name"`
		// Status contains the node status
		Status MessageState `json:"status"`
		// Taints contains node taints (optional)
		Taints *[]Taint `json:"taints,omitempty"`
		// Zone is the availability zone (optional)
		Zone *string `json:"zone,omitempty"`
	}

	// CreateNodePoolRequest represents the request payload for creating a node pool
	CreateNodePoolRequest struct {
		// Name is the name of the node pool
		Name string `json:"name"`
		// Flavor is the flavor identifier
		Flavor string `json:"flavor"`
		// Replicas is the number of replicas
		Replicas int `json:"replicas"`
		// Tags contains node tags (optional)
		Tags *[]string `json:"tags,omitempty"`
		// Taints contains node taints (optional)
		Taints *[]Taint `json:"taints,omitempty"`
		// AutoScale contains autoscaling configuration (optional)
		AutoScale *AutoScale `json:"auto_scale,omitempty"`
		// MaxPodsPerNode is the maximum number of pods per node (optional)
		MaxPodsPerNode *int `json:"max_pods_per_node,omitempty"`
		// AvailabilityZones contains availability zones (optional)
		AvailabilityZones *[]string `json:"availability_zones,omitempty"`
	}

	// PatchNodePoolRequest represents the request payload for updating a node pool
	PatchNodePoolRequest struct {
		// Replicas is the new number of replicas (optional)
		Replicas *int `json:"replicas,omitempty"`
		// AutoScale contains new autoscaling configuration (optional)
		AutoScale *AutoScale `json:"auto_scale,omitempty"`
	}

	// Taint represents a node taint
	Taint struct {
		// Key is the taint key
		Key string `json:"key"`
		// Value is the taint value
		Value string `json:"value"`
		// Effect is the taint effect
		Effect string `json:"effect"`
	}

	// AutoScale represents autoscaling configuration
	AutoScale struct {
		// MinReplicas is the minimum number of replicas (optional)
		MinReplicas *int `json:"min_replicas"`
		// MaxReplicas is the maximum number of replicas (optional)
		MaxReplicas *int `json:"max_replicas"`
	}

	// nodePoolService implements the NodePoolService interface
	nodePoolService struct {
		client *KubernetesClient
	}
)

// Nodes returns a list of nodes in a specific node pool
func (s *nodePoolService) Nodes(ctx context.Context, clusterID, nodePoolID string) ([]Node, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return nil, &client.ValidationError{Field: "nodePoolID", Message: utils.CannotBeEmpty}
	}

	type NodeList struct {
		Results []Node `json:"results"`
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[NodeList](ctx, s.client.newRequest,
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
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
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
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodPost,
		fmt.Sprintf("/v1alpha0/clusters/%s/node-pools", clusterID), req, nil)
}

// Get retrieves detailed information about a specific node pool
func (s *nodePoolService) Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return nil, &client.ValidationError{Field: "nodePoolID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodGet,
		fmt.Sprintf(clusterNodepoolURL, clusterID, nodePoolID), nil, nil)
}

// Update updates a node pool's properties
func (s *nodePoolService) Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return nil, &client.ValidationError{Field: "nodePoolID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[NodePool](ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodPatch,
		fmt.Sprintf(clusterNodepoolURL, clusterID, nodePoolID), req, nil)
}

// Delete removes a node pool from a cluster
func (s *nodePoolService) Delete(ctx context.Context, clusterID, nodePoolID string) error {
	if clusterID == "" {
		return &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	if nodePoolID == "" {
		return &client.ValidationError{Field: "nodePoolID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(ctx, s.client.newRequest,
		s.client.GetConfig(), http.MethodDelete,
		fmt.Sprintf(clusterNodepoolURL, clusterID, nodePoolID), nil, nil)
}
