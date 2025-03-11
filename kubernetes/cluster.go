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
	ClusterService interface {
		List(ctx context.Context, opts ListOptions) ([]ClusterList, error)
		Create(ctx context.Context, req ClusterRequest) (*CreateClusterResponse, error)
		Get(ctx context.Context, clusterID string, expand []string) (*Cluster, error)
		Delete(ctx context.Context, clusterID string) error
		Update(ctx context.Context, clusterID string, req AllowedCIDRsUpdateRequest) (*Cluster, error)
		GetKubeConfig(ctx context.Context, clusterID string) (*KubeConfig, error)
	}

	Network struct {
		UUID     string  `json:"uuid"`
		CIDR     string  `json:"cidr"`
		Name     *string `json:"name,omitempty"`
		SubnetID string  `json:"subnet_id"`
	}

	Addons struct {
		Loadbalance *string `json:"loadbalance,omitempty"`
		Volume      *string `json:"volume,omitempty"`
		Secrets     *string `json:"secrets,omitempty"`
	}

	KubeApiServer struct {
		DisableApiServerFip *bool   `json:"disable_api_server_fip,omitempty"`
		FixedIp             *string `json:"fixed_ip,omitempty"`
		FloatingIp          *string `json:"floating_ip,omitempty"`
		Port                *int    `json:"port,omitempty"`
	}

	AutoScaleResponse struct {
		MinReplicas *int `json:"min_replicas,omitempty"`
		MaxReplicas *int `json:"max_replicas,omitempty"`
	}

	ClusterListResponse struct {
		Results []ClusterList `json:"results"`
	}

	ClusterList struct {
		Description   *string        `json:"description,omitempty"`
		ID            string         `json:"id"`
		KubeApiServer *KubeApiServer `json:"kube_api_server,omitempty"`
		Name          string         `json:"name"`
		Region        *string        `json:"region,omitempty"`
		Status        *MessageState  `json:"status,omitempty"`
		Version       *string        `json:"version,omitempty"`
	}

	MessageState struct {
		State   string `json:"state"`
		Message string `json:"message"`
	}

	Cluster struct {
		Name          string         `json:"name"`
		ID            string         `json:"id"`
		Status        *Status        `json:"status"`
		Version       string         `json:"version"`
		Description   *string        `json:"description,omitempty"`
		Region        *string        `json:"region,omitempty"`
		CreatedAt     *time.Time     `json:"created_at"`
		UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
		Network       *Network       `json:"network,omitempty"`
		ControlPlane  *NodePool      `json:"controlplane,omitempty"`
		KubeApiServer *KubeApiServer `json:"kube_api_server,omitempty"`
		NodePools     *[]NodePool    `json:"node_pools,omitempty"`
		Addons        *Addons        `json:"addons,omitempty"`
		AllowedCIDRs  *[]string      `json:"allowed_cidrs,omitempty"`
	}

	Controlplane struct {
		AutoScale        AutoScale        `json:"auto_scale"`
		CreatedAt        *string          `json:"created_at,omitempty"`
		Id               string           `json:"id"`
		InstanceTemplate InstanceTemplate `json:"instance_template"`
		Labels           []string         `json:"labels"`
		Name             string           `json:"name"`
		Replicas         int              `json:"replicas"`
		SecurityGroups   *[]string        `json:"securityGroups,omitempty"`
		Status           *Status          `json:"status"`
		Tags             *[]string        `json:"tags,omitempty"`
		Taints           *[]Taint         `json:"taints,omitempty"`
		UpdatedAt        *string          `json:"updated_at,omitempty"`
		Zone             *[]string        `json:"zone"`
	}

	CreateClusterResponse struct {
		ID           string       `json:"id"`
		Name         string       `json:"name"`
		Status       MessageState `json:"status"`
		AllowedCidrs *[]string    `json:"allowed_cidrs,omitempty"`
	}

	ClusterRequest struct {
		Name               string                   `json:"name"`
		Version            *string                  `json:"version,omitempty"`
		Description        *string                  `json:"description,omitempty"`
		EnabledServerGroup *bool                    `json:"enabled_server_group,omitempty"`
		NodePools          *[]CreateNodePoolRequest `json:"node_pools,omitempty"`
		AllowedCIDRs       *[]string                `json:"allowed_cidrs,omitempty"`
	}

	AllowedCIDRsUpdateRequest struct {
		AllowedCIDRs []string `json:"allowed_cidrs"`
	}

	Status struct {
		State    string   `json:"state"`
		Messages []string `json:"messages,omitempty"`
	}

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

	clusterService struct {
		client *KubernetesClient
	}
)

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

func (s *clusterService) Create(ctx context.Context, req ClusterRequest) (*CreateClusterResponse, error) {
	response, err := mgc_http.ExecuteSimpleRequestWithRespBody[CreateClusterResponse](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodPost, "/v0/clusters", req, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *clusterService) Get(ctx context.Context, clusterID string, expand []string) (*Cluster, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	getClusterURL := fmt.Sprintf(clusterUrlWithID, clusterID)

	if len(expand) > 0 {
		q := url.Values{}
		q.Add("expand", strings.Join(expand, ","))
		getClusterURL += "?" + q.Encode()
	}

	cluster, err := mgc_http.ExecuteSimpleRequestWithRespBody[Cluster](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet, getClusterURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (s *clusterService) Delete(ctx context.Context, clusterID string) error {
	if clusterID == "" {
		return &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	err := mgc_http.ExecuteSimpleRequest(ctx, s.client.newRequest, s.client.GetConfig(), http.MethodDelete, fmt.Sprintf(clusterUrlWithID, clusterID), nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *clusterService) Update(ctx context.Context, clusterID string, req AllowedCIDRsUpdateRequest) (*Cluster, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[Cluster](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodPatch, fmt.Sprintf(clusterUrlWithID, clusterID), req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *clusterService) GetKubeConfig(ctx context.Context, clusterID string) (*KubeConfig, error) {
	if clusterID == "" {
		return nil, &client.ValidationError{Field: "clusterID", Message: utils.CannotBeEmpty}
	}

	kubeConfig, err := mgc_http.ExecuteSimpleRequestWithRespBody[KubeConfig](ctx, s.client.newRequest, s.client.GetConfig(), http.MethodGet, fmt.Sprintf("/v0/clusters/%s/kubeconfig", clusterID), nil, nil)
	if err != nil {
		return nil, err
	}

	return kubeConfig, nil
}
