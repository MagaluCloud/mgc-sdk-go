# Kubernetes

Package kubernetes provides a client for interacting with the Magalu Cloud
Kubernetes API. This package allows you to manage Kubernetes clusters, node
pools, flavors, and versions.
```
const (
DefaultBasePath = "/kubernetes"
)


```
```
type Addons struct {
Loadbalance *string `json:"loadbalance,omitempty"`
Volume      *string `json:"volume,omitempty"`
Secrets     *string `json:"secrets,omitempty"`
}
Addons represents cluster addons configuration

```
```
type Addresses struct {
Address string `json:"address"`
Type    string `json:"type"`
}
Addresses represents network addresses

```
```
type Allocatable struct {
CPU              string `json:"cpu"`
EphemeralStorage string `json:"ephemeral_storage"`
Hugepages1Gi     string `json:"hugepages_1Gi"`
Hugepages2Mi     string `json:"hugepages_2Mi"`
Memory           string `json:"memory"`
Pods             string `json:"pods"`
}
Allocatable represents allocatable resources

```
```
type AllowedCIDRsUpdateRequest struct {
AllowedCIDRs []string `json:"allowed_cidrs"`
}
AllowedCIDRsUpdateRequest represents the request payload for updating
allowed CIDRs

```
```
type AutoScale struct {
MinReplicas *int `json:"min_replicas"`
MaxReplicas *int `json:"max_replicas"`
}
AutoScale represents autoscaling configuration

```
```
type AutoScaleResponse struct {
MinReplicas *int `json:"min_replicas,omitempty"`
MaxReplicas *int `json:"max_replicas,omitempty"`
}
AutoScaleResponse represents autoscaling configuration

```
```
type Capacity struct {
CPU              string `json:"cpu"`
EphemeralStorage string `json:"ephemeral_storage"`
Hugepages1Gi     string `json:"hugepages_1Gi"`
Hugepages2Mi     string `json:"hugepages_2Mi"`
Memory           string `json:"memory"`
Pods             string `json:"pods"`
}
Capacity represents total capacity

```
```
type ClientOption func(*KubernetesClient)
ClientOption is a function type for configuring KubernetesClient options

```
```
type Cluster struct {
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
}
Cluster represents detailed information about a Kubernetes cluster

```
```
type ClusterList struct {
Description   *string        `json:"description,omitempty"`
ID            string         `json:"id"`
KubeApiServer *KubeApiServer `json:"kube_api_server,omitempty"`
Name          string         `json:"name"`
Region        *string        `json:"region,omitempty"`
Status        *MessageState  `json:"status,omitempty"`
Version       *string        `json:"version,omitempty"`
}
ClusterList represents a cluster in the list view

```
```
type ClusterListResponse struct {
Results []ClusterList `json:"results"`
}
ClusterListResponse represents the response when listing clusters

```
```
type ClusterRequest struct {
Name               string                   `json:"name"`
Version            *string                  `json:"version,omitempty"`
Description        *string                  `json:"description,omitempty"`
EnabledServerGroup *bool                    `json:"enabled_server_group,omitempty"`
NodePools          *[]CreateNodePoolRequest `json:"node_pools,omitempty"`
AllowedCIDRs       *[]string                `json:"allowed_cidrs,omitempty"`
ServicesIpV4CIDR   *string                  `json:"services_ipv4_cidr,omitempty"`
ClusterIPv4CIDR    *string                  `json:"cluster_ipv4_cidr,omitempty"`
}
ClusterRequest represents the request payload for creating a cluster

```
```
type ClusterService interface {
List(ctx context.Context, opts ListOptions) ([]ClusterList, error)
Create(ctx context.Context, req ClusterRequest) (*CreateClusterResponse, error)
Get(ctx context.Context, clusterID string) (*Cluster, error)
Delete(ctx context.Context, clusterID string) error
Update(ctx context.Context, clusterID string, req AllowedCIDRsUpdateRequest) (*Cluster, error)
GetKubeConfig(ctx context.Context, clusterID string) (*KubeConfig, error)
}
ClusterService provides methods for managing Kubernetes clusters

```
```
type Controlplane struct {
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
Controlplane represents control plane configuration

```
```
type CreateClusterResponse struct {
ID           string       `json:"id"`
Name         string       `json:"name"`
Status       MessageState `json:"status"`
AllowedCidrs *[]string    `json:"allowed_cidrs,omitempty"`
}
CreateClusterResponse represents the response when creating a cluster

```
```
type CreateNodePoolRequest struct {
Name              string     `json:"name"`
Flavor            string     `json:"flavor"`
Replicas          int        `json:"replicas"`
Tags              *[]string  `json:"tags,omitempty"`
Taints            *[]Taint   `json:"taints,omitempty"`
AutoScale         *AutoScale `json:"auto_scale,omitempty"`
MaxPodsPerNode    *int       `json:"max_pods_per_node,omitempty"`
AvailabilityZones *[]string  `json:"availability_zones,omitempty"`
}
CreateNodePoolRequest represents the request payload for creating a node
pool

```
```
type Flavor struct {
Name string `json:"name"`
ID   string `json:"id"`
VCPU int    `json:"vcpu"`
RAM  int    `json:"ram"`
Size int    `json:"size"`
}
Flavor represents a Kubernetes flavor (instance type)

```
```
type FlavorList struct {
Results []FlavorsAvailable `json:"results"`
}
FlavorList represents the response when listing flavors

```
```
type FlavorService interface {
List(ctx context.Context, opts ListOptions) (*FlavorsAvailable, error)
}
FlavorService provides methods for managing Kubernetes flavors

```
```
type FlavorsAvailable struct {
NodePool     []Flavor `json:"nodepool"`
ControlPlane []Flavor `json:"controlplane"`
}
FlavorsAvailable represents available flavors for different components

```
```
type Infrastructure struct {
Allocatable             Allocatable `json:"allocatable"`
Architecture            string      `json:"architecture"`
Capacity                Capacity    `json:"capacity"`
ContainerRuntimeVersion string      `json:"containerRuntimeVersion"`
KernelVersion           string      `json:"kernelVersion"`
KubeProxyVersion        string      `json:"kubeProxyVersion"`
KubeletVersion          string      `json:"kubeletVersion"`
OperatingSystem         string      `json:"operatingSystem"`
OsImage                 string      `json:"osImage"`
}
Infrastructure represents node infrastructure information

```
```
type InstanceTemplate struct {
Flavor    Flavor `json:"flavor"`
NodeImage string `json:"node_image"`
DiskSize  int    `json:"disk_size"`
DiskType  string `json:"disk_type"`
}
InstanceTemplate represents the template for node instances

```
```
type KubeApiServer struct {
DisableApiServerFip *bool   `json:"disable_api_server_fip,omitempty"`
FixedIp             *string `json:"fixed_ip,omitempty"`
FloatingIp          *string `json:"floating_ip,omitempty"`
Port                *int    `json:"port,omitempty"`
}
KubeApiServer represents Kubernetes API server configuration

```
```
type KubeConfig struct {
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
KubeConfig represents a Kubernetes configuration file

```
```
type KubernetesClient struct {
*client.CoreClient
}
KubernetesClient represents a client for the Kubernetes service

```
```
func New(core *client.CoreClient, opts ...ClientOption) *KubernetesClient
New creates a new KubernetesClient instance with the provided core client
and options

```
```
func (c *KubernetesClient) Clusters() ClusterService
Clusters returns a service for managing Kubernetes clusters

```
```
func (c *KubernetesClient) Flavors() FlavorService
Flavors returns a service for managing Kubernetes flavors

```
```
func (c *KubernetesClient) Nodepools() NodePoolService
Nodepools returns a service for managing Kubernetes node pools

```
```
func (c *KubernetesClient) Versions() VersionService
Versions returns a service for managing Kubernetes versions

```
```
type ListOptions struct {
Limit  *int
Offset *int
Sort   *string
Expand []string
}
ListOptions provides options for listing resources

```
```
type MessageState struct {
State   string `json:"state"`
Message string `json:"message"`
}
MessageState represents a status message

```
```
type Network struct {
UUID     string  `json:"uuid"`
CIDR     string  `json:"cidr"`
Name     *string `json:"name,omitempty"`
SubnetID string  `json:"subnet_id"`
}
Network represents network configuration for a cluster

```
```
type Node struct {
Addresses      []Addresses       `json:"addresses"`
Annotations    map[string]string `json:"annotations"`
ClusterName    string            `json:"cluster_name"`
CreatedAt      time.Time         `json:"created_at"`
Flavor         string            `json:"flavor"`
ID             string            `json:"id"`
Infrastructure Infrastructure    `json:"infrastructure"`
Labels         map[string]string `json:"labels"`
Name           string            `json:"name"`
Namespace      string            `json:"namespace"`
NodeImage      string            `json:"node_image"`
NodepoolName   string            `json:"nodepool_name"`
Status         MessageState      `json:"status"`
Taints         *[]Taint          `json:"taints,omitempty"`
Zone           *string           `json:"zone,omitempty"`
}
Node represents a Kubernetes node

```
```
type NodePool struct {
ID                string            `json:"id"`
Name              string            `json:"name"`
InstanceTemplate  InstanceTemplate  `json:"instance_template"`
Replicas          int               `json:"replicas"`
Zone              *[]string         `json:"zone,omitempty"`
Tags              *[]string         `json:"tags"`
Labels            map[string]string `json:"labels,omitempty"`
Taints            *[]Taint          `json:"taints,omitempty"`
SecurityGroups    *[]string         `json:"security_groups,omitempty"`
CreatedAt         *time.Time        `json:"created_at"`
UpdatedAt         *time.Time        `json:"updated_at,omitempty"`
AutoScale         *AutoScale        `json:"auto_scale,omitempty"`
Status            Status            `json:"status"`
Flavor            string            `json:"flavor"`
MaxPodsPerNode    *int              `json:"max_pods_per_node,omitempty"`
AvailabilityZones *[]string         `json:"availability_zones,omitempty"`
}
NodePool represents a Kubernetes node pool

```
```
type NodePoolList struct {
Results []NodePool `json:"results"`
}
NodePoolList represents the response when listing node pools

```
```
type NodePoolService interface {
Nodes(ctx context.Context, clusterID, nodePoolID string) ([]Node, error)
List(ctx context.Context, clusterID string, opts ListOptions) ([]NodePool, error)
Create(ctx context.Context, clusterID string, req CreateNodePoolRequest) (*NodePool, error)
Get(ctx context.Context, clusterID, nodePoolID string) (*NodePool, error)
Update(ctx context.Context, clusterID, nodePoolID string, req PatchNodePoolRequest) (*NodePool, error)
Delete(ctx context.Context, clusterID, nodePoolID string) error
}
NodePoolService provides methods for managing Kubernetes node pools

```
```
type PatchNodePoolRequest struct {
Replicas  *int       `json:"replicas,omitempty"`
AutoScale *AutoScale `json:"auto_scale,omitempty"`
}
PatchNodePoolRequest represents the request payload for updating a node pool

```
```
type Status struct {
State    string   `json:"state"`
Messages []string `json:"messages,omitempty"`
}
Status represents a status with messages

```
```
type Taint struct {
Key    string `json:"key"`
Value  string `json:"value"`
Effect string `json:"effect"`
}
Taint represents a node taint

```
```
type Version struct {
Version    string `json:"version"`
Deprecated bool   `json:"deprecated"`
}
Version represents a Kubernetes version

```
```
type VersionList struct {
Results []Version `json:"results"`
}
VersionList represents the response when listing versions

```
```
type VersionService interface {
List(ctx context.Context) ([]Version, error)
}
VersionService provides methods for managing Kubernetes versions


```

