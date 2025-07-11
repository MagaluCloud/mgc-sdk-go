package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

const (
	// clusterUrlWithID is the URL template for cluster operations with ID
	clusterUrlWithID = "/v0/clusters/%s"
)

type (
	// ClusterService provides methods for managing Kubernetes clusters
	ClusterService interface {
		// List returns a list of Kubernetes clusters with optional filtering and pagination
		List(ctx context.Context, opts ListOptions) ([]ClusterList, error)
		// Create creates a new Kubernetes cluster
		Create(ctx context.Context, req ClusterRequest) (*CreateClusterResponse, error)
		// Get retrieves detailed information about a specific cluster
		Get(ctx context.Context, clusterID string) (*Cluster, error)
		// Delete removes a Kubernetes cluster
		Delete(ctx context.Context, clusterID string) error
		// Update updates the allowed CIDRs for a cluster
		Update(ctx context.Context, clusterID string, req AllowedCIDRsUpdateRequest) (*Cluster, error)
		// GetKubeConfig retrieves the kubeconfig for a cluster
		GetKubeConfig(ctx context.Context, clusterID string) (*KubeConfig, error)
	}

	// Network represents network configuration for a cluster
	Network struct {
		// UUID is the network UUID
		UUID string `json:"uuid"`
		// CIDR is the network CIDR
		CIDR string `json:"cidr"`
		// Name is the network name (optional)
		Name *string `json:"name,omitempty"`
		// SubnetID is the subnet identifier
		SubnetID string `json:"subnet_id"`
	}

	// Addons represents cluster addons configuration
	Addons struct {
		// Loadbalance is the load balancer addon (optional)
		Loadbalance *string `json:"loadbalance,omitempty"`
		// Volume is the volume addon (optional)
		Volume *string `json:"volume,omitempty"`
		// Secrets is the secrets addon (optional)
		Secrets *string `json:"secrets,omitempty"`
	}

	// KubeApiServer represents Kubernetes API server configuration
	KubeApiServer struct {
		// DisableApiServerFip disables floating IP for API server (optional)
		DisableApiServerFip *bool `json:"disable_api_server_fip,omitempty"`
		// FixedIp is the fixed IP for API server (optional)
		FixedIp *string `json:"fixed_ip,omitempty"`
		// FloatingIp is the floating IP for API server (optional)
		FloatingIp *string `json:"floating_ip,omitempty"`
		// Port is the API server port (optional)
		Port *int `json:"port,omitempty"`
	}

	// AutoScaleResponse represents autoscaling configuration
	AutoScaleResponse struct {
		// MinReplicas is the minimum number of replicas (optional)
		MinReplicas *int `json:"min_replicas,omitempty"`
		// MaxReplicas is the maximum number of replicas (optional)
		MaxReplicas *int `json:"max_replicas,omitempty"`
	}

	// ClusterListResponse represents the response when listing clusters
	ClusterListResponse struct {
		// Results is the list of clusters
		Results []ClusterList `json:"results"`
	}

	// ClusterList represents a cluster in the list view
	ClusterList struct {
		// Description is the cluster description (optional)
		Description *string `json:"description,omitempty"`
		// ID is the unique identifier of the cluster
		ID string `json:"id"`
		// KubeApiServer contains API server configuration (optional)
		KubeApiServer *KubeApiServer `json:"kube_api_server,omitempty"`
		// Name is the name of the cluster
		Name string `json:"name"`
		// Region is the cluster region (optional)
		Region *string `json:"region,omitempty"`
		// Status contains the cluster status (optional)
		Status *MessageState `json:"status,omitempty"`
		// Version is the Kubernetes version (optional)
		Version *string `json:"version,omitempty"`
	}

	// MessageState represents a status message
	MessageState struct {
		// State is the current state
		State string `json:"state"`
		// Message is the status message
		Message string `json:"message"`
	}

	// Cluster represents detailed information about a Kubernetes cluster
	Cluster struct {
		// Name is the name of the cluster
		Name string `json:"name"`
		// ID is the unique identifier of the cluster
		ID string `json:"id"`
		// Status contains the cluster status
		Status *MessageState `json:"status"`
		// Version is the Kubernetes version
		Version string `json:"version"`
		// Description is the cluster description (optional)
		Description *string `json:"description,omitempty"`
		// Region is the cluster region (optional)
		Region *string `json:"region,omitempty"`
		// CreatedAt is the timestamp when the cluster was created
		CreatedAt *time.Time `json:"created_at"`
		// UpdatedAt is the timestamp when the cluster was last updated (optional)
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
		// Network contains network configuration (optional)
		Network *Network `json:"network,omitempty"`
		// ControlPlane contains control plane configuration (optional)
		ControlPlane *NodePool `json:"controlplane,omitempty"`
		// KubeApiServer contains API server configuration (optional)
		KubeApiServer *KubeApiServer `json:"kube_api_server,omitempty"`
		// NodePools contains the list of node pools (optional)
		NodePools *[]NodePool `json:"node_pools,omitempty"`
		// Addons contains addon configuration (optional)
		Addons *Addons `json:"addons,omitempty"`
		// AllowedCIDRs contains allowed CIDR ranges (optional)
		AllowedCIDRs *[]string `json:"allowed_cidrs,omitempty"`
		// ServicesIpV4CIDR is the services IPv4 CIDR (optional)
		ServicesIpV4CIDR *string `json:"services_ipv4_cidr,omitempty"`
		// ClusterIPv4CIDR is the cluster IPv4 CIDR (optional)
		ClusterIPv4CIDR *string `json:"cluster_ipv4_cidr,omitempty"`
	}

	// Controlplane represents control plane configuration
	Controlplane struct {
		// AutoScale contains autoscaling configuration
		AutoScale AutoScale `json:"auto_scale"`
		// CreatedAt is the creation timestamp (optional)
		CreatedAt *string `json:"created_at,omitempty"`
		// Id is the control plane identifier
		Id string `json:"id"`
		// InstanceTemplate contains instance template configuration
		InstanceTemplate InstanceTemplate `json:"instance_template"`
		// Labels contains node labels
		Labels []string `json:"labels"`
		// Name is the control plane name
		Name string `json:"name"`
		// Replicas is the number of replicas
		Replicas int `json:"replicas"`
		// SecurityGroups contains security group IDs (optional)
		SecurityGroups *[]string `json:"securityGroups,omitempty"`
		// Status contains the control plane status
		Status *Status `json:"status"`
		// Tags contains node tags (optional)
		Tags *[]string `json:"tags,omitempty"`
		// Taints contains node taints (optional)
		Taints *[]Taint `json:"taints,omitempty"`
		// UpdatedAt is the last update timestamp (optional)
		UpdatedAt *string `json:"updated_at,omitempty"`
		// Zone contains availability zones
		Zone *[]string `json:"zone"`
	}

	// CreateClusterResponse represents the response when creating a cluster
	CreateClusterResponse struct {
		// ID is the unique identifier of the created cluster
		ID string `json:"id"`
		// Name is the name of the created cluster
		Name string `json:"name"`
		// Status contains the cluster status
		Status MessageState `json:"status"`
		// AllowedCidrs contains allowed CIDR ranges (optional)
		AllowedCidrs *[]string `json:"allowed_cidrs,omitempty"`
	}

	// ClusterRequest represents the request payload for creating a cluster
	ClusterRequest struct {
		// Name is the name of the cluster
		Name string `json:"name"`
		// Version is the Kubernetes version (optional)
		Version *string `json:"version,omitempty"`
		// Description is the cluster description (optional)
		Description *string `json:"description,omitempty"`
		// EnabledServerGroup enables server groups (optional)
		EnabledServerGroup *bool `json:"enabled_server_group,omitempty"`
		// NodePools contains node pool configurations (optional)
		NodePools *[]CreateNodePoolRequest `json:"node_pools,omitempty"`
		// AllowedCIDRs contains allowed CIDR ranges (optional)
		AllowedCIDRs *[]string `json:"allowed_cidrs,omitempty"`
		// ServicesIpV4CIDR is the services IPv4 CIDR (optional)
		ServicesIpV4CIDR *string `json:"services_ipv4_cidr,omitempty"`
		// ClusterIPv4CIDR is the cluster IPv4 CIDR (optional)
		ClusterIPv4CIDR *string `json:"cluster_ipv4_cidr,omitempty"`
	}

	// AllowedCIDRsUpdateRequest represents the request payload for updating allowed CIDRs
	AllowedCIDRsUpdateRequest struct {
		// AllowedCIDRs is the list of allowed CIDR ranges
		AllowedCIDRs []string `json:"allowed_cidrs"`
	}

	// Status represents a status with messages
	Status struct {
		// State is the current state
		State string `json:"state"`
		// Messages contains status messages (optional)
		Messages []string `json:"messages,omitempty"`
	}

	// KubeConfig represents a Kubernetes configuration file
	KubeConfig struct {
		// APIVersion is the API version
		APIVersion string `yaml:"apiVersion"`
		// Clusters contains cluster configurations
		Clusters []struct {
			Cluster struct {
				// CertificateAuthorityData contains the CA certificate data
				CertificateAuthorityData string `yaml:"certificate-authority-data"`
				// Server is the cluster server URL
				Server string `yaml:"server"`
			} `yaml:"cluster"`
			// Name is the cluster name
			Name string `yaml:"name"`
		} `yaml:"clusters"`
		// Contexts contains context configurations
		Contexts []struct {
			Context struct {
				// Cluster is the cluster name
				Cluster string `yaml:"cluster"`
				// Namespace is the default namespace
				Namespace string `yaml:"namespace"`
				// User is the user name
				User string `yaml:"user"`
			} `yaml:"context"`
			// Name is the context name
			Name string `yaml:"name"`
		} `yaml:"contexts"`
		// CurrentContext is the current context name
		CurrentContext string `yaml:"current-context"`
		// Kind is the resource kind
		Kind string `yaml:"kind"`
		// Users contains user configurations
		Users []struct {
			// Name is the user name
			Name string `yaml:"name"`
			User struct {
				// ClientCertificateData contains the client certificate data
				ClientCertificateData string `yaml:"client-certificate-data"`
				// ClientKeyData contains the client key data
				ClientKeyData string `yaml:"client-key-data"`
			} `yaml:"user"`
		} `yaml:"users"`
	}

	// clusterService implements the ClusterService interface
	clusterService struct {
		client *KubernetesClient
	}
)

// List returns a list of Kubernetes clusters with optional filtering and pagination
func (s *clusterService) List(ctx context.Context, opts ListOptions) ([]ClusterList, error) {
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

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[ClusterListResponse](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet, "/v0/clusters", nil, query)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

// Create creates a new Kubernetes cluster
func (s *clusterService) Create(ctx context.Context, req ClusterRequest) (*CreateClusterResponse, error) {
	response, err := mgc_http.ExecuteSimpleRequestWithRespBody[CreateClusterResponse](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodPost, "/v0/clusters", req, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Get retrieves detailed information about a specific cluster
func (s *clusterService) Get(ctx context.Context, clusterID string) (*Cluster, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Cluster](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet, fmt.Sprintf(clusterUrlWithID, clusterID), nil, nil)
}

// Delete removes a Kubernetes cluster
func (s *clusterService) Delete(ctx context.Context, clusterID string) error {
	if clusterID == "" {
		return &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(ctx, s.client.newRequest, s.client.GetConfig(), http.MethodDelete, fmt.Sprintf(clusterUrlWithID, clusterID), nil, nil)
}

// Update updates the allowed CIDRs for a cluster
func (s *clusterService) Update(ctx context.Context, clusterID string, req AllowedCIDRsUpdateRequest) (*Cluster, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Cluster](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodPatch, fmt.Sprintf(clusterUrlWithID, clusterID), req, nil)
}

// GetKubeConfig retrieves the kubeconfig for a cluster
func (s *clusterService) GetKubeConfig(ctx context.Context, clusterID string) (*KubeConfig, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[KubeConfig](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet, fmt.Sprintf(clusterUrlWithID+"/kubeconfig", clusterID), nil, nil)
}
