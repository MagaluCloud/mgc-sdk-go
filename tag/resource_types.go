package tag

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ResourceEnum represents a supported cloud resource type name
type ResourceEnum string

const (
	ResourceBucket               ResourceEnum = "bucket"
	ResourceObjects              ResourceEnum = "objects"
	ResourceNatGateways          ResourceEnum = "nat-gateways"
	ResourcePorts                ResourceEnum = "ports"
	ResourcePublicIPs            ResourceEnum = "public-ips"
	ResourceRules                ResourceEnum = "rules"
	ResourceSecurityGroups       ResourceEnum = "security-groups"
	ResourceSubnets              ResourceEnum = "subnets"
	ResourceVPCs                 ResourceEnum = "vpcs"
	ResourceImages               ResourceEnum = "images"
	ResourceInstancesVM          ResourceEnum = "instances-vm"
	ResourceSnapshotsVM          ResourceEnum = "snapshots-vm"
	ResourceSnapshotsBS          ResourceEnum = "snapshots-bs"
	ResourceVolumes              ResourceEnum = "volumes"
	ResourceCluster              ResourceEnum = "cluster"
	ResourceNodepool             ResourceEnum = "nodepool"
	ResourceRegistries           ResourceEnum = "registries"
	ResourceRepositories         ResourceEnum = "repositories"
	ResourceClusters             ResourceEnum = "clusters"
	ResourceInstancesDB          ResourceEnum = "instances-db"
	ResourceParameterGroups      ResourceEnum = "parameter-groups"
	ResourceReplicas             ResourceEnum = "replicas"
	ResourceSnapshotsDB          ResourceEnum = "snapshots-db"
	ResourceNetworkACLs          ResourceEnum = "network-acls"
	ResourceNetworkBackends      ResourceEnum = "network-backends"
	ResourceNetworkCertificates  ResourceEnum = "network-certificates"
	ResourceNetworkHealthcheck   ResourceEnum = "network-healthcheck"
	ResourceNetworkListeners     ResourceEnum = "network-listeners"
	ResourceNetworkLoadbalancers ResourceEnum = "network-loadbalancers"
)

// ProductEnum represents a Magalu Cloud product
type ProductEnum string

const (
	ProductObjectStorage     ProductEnum = "object-storage"
	ProductNetwork           ProductEnum = "network"
	ProductVirtualMachine    ProductEnum = "virtual-machine"
	ProductBlockStorage      ProductEnum = "block-storage"
	ProductKubernetes        ProductEnum = "kubernetes"
	ProductContainerRegistry ProductEnum = "container-registry"
	ProductDatabase          ProductEnum = "database"
	ProductLoadBalancer      ProductEnum = "load-balancer"
)

type (
	// ResourceType represents a cloud resource type
	ResourceType struct {
		Name      []ResourceEnum `json:"name"`
		Product   []ProductEnum  `json:"product"`
		CreatedAt string         `json:"created_at"`
		UpdatedAt *string        `json:"updated_at"`
	}

	// ResourceTypeListResponse represents a list of resource types
	ResourceTypeListResponse struct {
		Results []ResourceType `json:"results"`
	}

	// ListResourceTypesOptions defines parameters for filtering and paginating resource types
	ListResourceTypesOptions struct {
		Name    *ResourceEnum
		Product *ProductEnum
		Limit   *int
		Offset  *int
		Sort    *string
	}
)

// ResourceTypeService provides methods for listing available cloud resource types.
type ResourceTypeService interface {
	List(ctx context.Context, opts ListResourceTypesOptions) ([]ResourceType, error)
}

// resourceTypeService implements ResourceTypeService
type resourceTypeService struct {
	client *TagClient
}

// List returns all available cloud resource types
func (s *resourceTypeService) List(ctx context.Context, opts ListResourceTypesOptions) ([]ResourceType, error) {
	query := make(url.Values)
	if opts.Name != nil {
		query.Set("name", string(*opts.Name))
	}
	if opts.Product != nil {
		query.Set("product", string(*opts.Product))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ResourceTypeListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/resource-types",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}
