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
	clusterUrlWithID = "/v0/clusters/%s"
)

type (
	// ClusterService provides methods for managing Kubernetes clusters
	ClusterService interface {
		List(ctx context.Context, opts ListOptions) ([]ClusterList, error)
		Create(ctx context.Context, req ClusterRequest) (*CreateClusterResponse, error)
		Get(ctx context.Context, clusterID string) (*Cluster, error)
		Delete(ctx context.Context, clusterID string) error
		Update(ctx context.Context, clusterID string, req PatchClusterRequest) (*PatchClusterResponse, error)
		GetKubeConfig(ctx context.Context, clusterID string) (*KubeConfig, error)
	}

	// Network represents network configuration for a cluster
	Network struct {
		UUID     string `json:"uuid"`
		CIDR     string `json:"cidr"`
		Name     string `json:"name"`
		SubnetID string `json:"subnet_id"`
	}

	// Addons represents cluster addons configuration
	Addons struct {
		Loadbalance string `json:"loadbalance"`
		Volume      string `json:"volume"`
		Secrets     string `json:"secrets"`
	}

	// KubeApiServer represents Kubernetes API server configuration
	KubeApiServer struct {
		DisableApiServerFip *bool   `json:"disable_api_server_fip,omitempty"`
		FixedIp             *string `json:"fixed_ip,omitempty"`
		FloatingIp          *string `json:"floating_ip,omitempty"`
		Port                *int    `json:"port,omitempty"`
	}

	// ClusterListResponse represents the response when listing clusters
	ClusterListResponse struct {
		Results []ClusterList `json:"results"`
	}

	// ClusterList represents a cluster in the list view
	ClusterList struct {
		Description   *string        `json:"description,omitempty"`
		ID            string         `json:"id"`
		KubeApiServer *KubeApiServer `json:"kube_api_server,omitempty"`
		Name          string         `json:"name"`
		Region        *string        `json:"region,omitempty"`
		Status        *MessageState  `json:"status,omitempty"`
		Version       *string        `json:"version,omitempty"`

		CreatedAt          *time.Time          `json:"created_at,omitempty"`
		MachineTypesSource *MachineTypesSource `json:"machine_types_source,omitempty"`
		ClusterIPv4CIDR    *string             `json:"cluster_ipv4_cidr,omitempty"`
		ServicesIpV4CIDR   *string             `json:"services_ipv4_cidr,omitempty"`
		Platform           *Platform           `json:"platform,omitempty"`
	}

	// MessageState represents a status message
	MessageState struct {
		State   string `json:"state"`
		Message string `json:"message"`
	}

	// Cluster represents detailed information about a Kubernetes cluster
	Cluster struct {
		Name             string         `json:"name"`
		ID               string         `json:"id"`
		Status           *MessageState  `json:"status"`
		Version          string         `json:"version"`
		Description      *string        `json:"description,omitempty"`
		Region           *string        `json:"region,omitempty"`
		CreatedAt        *time.Time     `json:"created_at"`
		UpdatedAt        *time.Time     `json:"updated_at,omitempty"`
		Network          *Network       `json:"network,omitempty"`
		ControlPlane     *NodePool      `json:"controlplane,omitempty"`
		KubeApiServer    *KubeApiServer `json:"kube_api_server,omitempty"`
		NodePools        *[]NodePool    `json:"node_pools,omitempty"`
		Addons           *Addons        `json:"addons,omitempty"`
		AllowedCIDRs     *[]string      `json:"allowed_cidrs,omitempty"`
		ServicesIpV4CIDR *string        `json:"services_ipv4_cidr,omitempty"`
		ClusterIPv4CIDR  *string        `json:"cluster_ipv4_cidr,omitempty"`

		MachineTypesSource *MachineTypesSource `json:"machine_types_source,omitempty"`
		Platform           *Platform           `json:"platform,omitempty"`
	}

	// CreateClusterResponse represents the response when creating a cluster
	CreateClusterResponse struct {
		ID               string       `json:"id"`
		Name             string       `json:"name"`
		Status           MessageState `json:"status"`
		AllowedCidrs     *[]string    `json:"allowed_cidrs,omitempty"`
		ClusterIPv4CIDR  *string      `json:"cluster_ipv4_cidr,omitempty"`
		ServicesIpV4CIDR *string      `json:"services_ipv4_cidr,omitempty"`
	}

	// ClusterRequest represents the request payload for creating a cluster
	ClusterRequest struct {
		Name               string                   `json:"name"`
		Version            *string                  `json:"version,omitempty"`
		Description        *string                  `json:"description,omitempty"`
		EnabledServerGroup *bool                    `json:"enabled_server_group,omitempty"`
		NodePools          *[]CreateNodePoolRequest `json:"node_pools,omitempty"`
		AllowedCIDRs       *[]string                `json:"allowed_cidrs,omitempty"`
		ServicesIpV4CIDR   *string                  `json:"services_ipv4_cidr,omitempty"`
		ClusterIPv4CIDR    *string                  `json:"cluster_ipv4_cidr,omitempty"`
	}

	// PatchClusterRequest represents the request payload for patching a cluster
	PatchClusterRequest struct {
		AllowedCIDRs *[]string `json:"allowed_cidrs,omitempty"`
	}

	// PatchClusterResponse represents the response when patching a cluster
	PatchClusterResponse struct {
		AllowedCIDRs *[]string `json:"allowed_cidrs,omitempty"`
	}

	// MachineTypesSource represents the source of machine types
	MachineTypesSource string

	// Platform represents platform information
	Platform struct {
		Version string `json:"version"`
	}

	// Status represents a status with messages
	Status struct {
		State    string   `json:"state"`
		Messages []string `json:"messages,omitempty"`
	}

	// KubeConfig represents a Kubernetes configuration file
	KubeConfig struct {
		APIVersion string `yaml:"apiVersion"`
		Clusters   []struct {
			Cluster struct {
				CertificateAuthorityData string `yaml:"certificate-authority-data"`
				Server                   string `yaml:"server"`
			} `yaml:"cluster"`
			Name string `yaml:"name"`
		} `yaml:"clusters"`
		Contexts []struct {
			Context struct {
				Cluster   string `yaml:"cluster"`
				Namespace string `yaml:"namespace"`
				User      string `yaml:"user"`
			} `yaml:"context"`
			Name string `yaml:"name"`
		} `yaml:"contexts"`
		CurrentContext string `yaml:"current-context"`
		Kind           string `yaml:"kind"`
		Users          []struct {
			Name string `yaml:"name"`
			User struct {
				ClientCertificateData string `yaml:"client-certificate-data"`
				ClientKeyData         string `yaml:"client-key-data"`
			} `yaml:"user"`
		} `yaml:"users"`
	}

	// clusterService implements the ClusterService interface
	clusterService struct {
		client *KubernetesClient
	}
)

// Constants for MachineTypesSource
const (
	MachineTypesSourceExternal MachineTypesSource = "external"
	MachineTypesSourceInternal MachineTypesSource = "internal"
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
func (s *clusterService) Update(ctx context.Context, clusterID string, req PatchClusterRequest) (*PatchClusterResponse, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[PatchClusterResponse](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodPatch, fmt.Sprintf(clusterUrlWithID, clusterID), req, nil)
}

// GetKubeConfig retrieves the kubeconfig for a cluster
func (s *clusterService) GetKubeConfig(ctx context.Context, clusterID string) (*KubeConfig, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[KubeConfig](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet, fmt.Sprintf(clusterUrlWithID+"/kubeconfig", clusterID), nil, nil)
}
