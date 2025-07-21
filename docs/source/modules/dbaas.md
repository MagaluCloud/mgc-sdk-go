# Dbaas

Package dbaas provides a client for interacting with the Magalu Cloud Database
as a Service (DBaaS) API. This package allows you to manage database instances,
clusters, replicas, engines, instance types, and parameters.
```
const (
InstancePath   = "/v2/instances"
InstancePathID = InstancePath + "/%s"
SnapshotPath   = InstancePathID + "/snapshots"
SnapshotPathID = SnapshotPath + "/%s"
)
```
```
const (
ParameterGroupTypeSystem ParameterGroupType = "SYSTEM"
ParameterGroupTypeUser   ParameterGroupType = "USER"

ErrorIDEmpty        = "ID cannot be empty"
PathParametersGroup = "/v2/parameter-groups"
)
```
```
const (
DefaultBasePath = "/database"
)


```
```
type Address struct {
Access  AddressAccess `json:"access"`
Type    *AddressType  `json:"type,omitempty"`
Address *string       `json:"address,omitempty"`
}
Address represents a network address for an instance

```
```
type AddressAccess string

```
```
const (
AddressAccessPrivate AddressAccess = "PRIVATE"
AddressAccessPublic  AddressAccess = "PUBLIC"
)
```
```
type AddressType string

```
```
const (
AddressTypeIPv4 AddressType = "IPv4"
AddressTypeIPv6 AddressType = "IPv6"
)
```
```
type ClientOption func(*DBaaSClient)
ClientOption is a function type for configuring DBaaSClient options

```
```
type ClusterCreateRequest struct {
Name                string               `json:"name"`
EngineID            string               `json:"engine_id"`
InstanceTypeID      string               `json:"instance_type_id"`
User                string               `json:"user"`
Password            string               `json:"password"`
Volume              ClusterVolumeRequest `json:"volume"`
ParameterGroupID    *string              `json:"parameter_group_id,omitempty"`
BackupRetentionDays *int                 `json:"backup_retention_days,omitempty"`
BackupStartAt       *string              `json:"backup_start_at,omitempty"`
}
ClusterCreateRequest represents the request payload for creating a cluster

```
```
type ClusterDetailResponse struct {
ID                     string                `json:"id"`
Name                   string                `json:"name"`
EngineID               string                `json:"engine_id"`
InstanceTypeID         string                `json:"instance_type_id"`
ParameterGroupID       string                `json:"parameter_group_id"`
Volume                 ClusterVolumeResponse `json:"volume"`
Status                 ClusterStatus         `json:"status"`
Addresses              []LoadBalancerAddress `json:"addresses"`
ApplyParametersPending bool                  `json:"apply_parameters_pending"`
BackupRetentionDays    int                   `json:"backup_retention_days"`
BackupStartAt          string                `json:"backup_start_at"`
CreatedAt              time.Time             `json:"created_at"`
UpdatedAt              *time.Time            `json:"updated_at,omitempty"`
StartedAt              *string               `json:"started_at,omitempty"`
FinishedAt             *string               `json:"finished_at,omitempty"`
}
ClusterDetailResponse represents detailed information about a cluster

```
```
type ClusterResponse struct {
ID string `json:"id"`
}
ClusterResponse represents the response when creating a cluster

```
```
type ClusterService interface {
List(ctx context.Context, opts ListClustersOptions) ([]ClusterDetailResponse, error)
Create(ctx context.Context, req ClusterCreateRequest) (*ClusterResponse, error)
Get(ctx context.Context, ID string) (*ClusterDetailResponse, error)
Update(ctx context.Context, ID string, req ClusterUpdateRequest) (*ClusterDetailResponse, error)
Delete(ctx context.Context, ID string) error
Start(ctx context.Context, ID string) (*ClusterDetailResponse, error)
Stop(ctx context.Context, ID string) (*ClusterDetailResponse, error)
}
ClusterService provides methods for managing database clusters

```
```
type ClusterStatus string
ClusterStatus represents the possible states of a cluster

```
```
const (
ClusterStatusActive        ClusterStatus = "ACTIVE"
ClusterStatusError         ClusterStatus = "ERROR"
ClusterStatusPending       ClusterStatus = "PENDING"
ClusterStatusCreating      ClusterStatus = "CREATING"
ClusterStatusDeleting      ClusterStatus = "DELETING"
ClusterStatusDeleted       ClusterStatus = "DELETED"
ClusterStatusErrorDeleting ClusterStatus = "ERROR_DELETING"
ClusterStatusStarting      ClusterStatus = "STARTING"
ClusterStatusStopping      ClusterStatus = "STOPPING"
ClusterStatusStopped       ClusterStatus = "STOPPED"
ClusterStatusBackingUp     ClusterStatus = "BACKING_UP"
)
```
```
type ClusterUpdateRequest struct {
ParameterGroupID    *string `json:"parameter_group_id,omitempty"`
BackupRetentionDays *int    `json:"backup_retention_days,omitempty"`
BackupStartAt       *string `json:"backup_start_at,omitempty"`
}
ClusterUpdateRequest represents the request payload for updating a cluster

```
```
type ClusterVolumeRequest struct {
Size int     `json:"size"`
Type *string `json:"type,omitempty"`
}
ClusterVolumeRequest represents volume configuration for cluster creation

```
```
type ClusterVolumeResponse struct {
Size int    `json:"size"`
Type string `json:"type"`
}
ClusterVolumeResponse represents volume information for a cluster

```
```
type ClustersResponse struct {
Results []ClusterDetailResponse `json:"results"`
}
ClustersResponse represents the response when listing clusters

```
```
type DBaaSClient struct {
*client.CoreClient
}
DBaaSClient represents a client for the Database as a Service

```
```
func New(core *client.CoreClient, opts ...ClientOption) *DBaaSClient
New creates a new DBaaSClient instance with the provided core client and
options

```
```
func (c *DBaaSClient) Clusters() ClusterService
Clusters returns a service for managing database clusters

```
```
func (c *DBaaSClient) Engines() EngineService
Engines returns a service for managing database engines

```
```
func (c *DBaaSClient) InstanceTypes() InstanceTypeService
InstanceTypes returns a service for managing database instance types

```
```
func (c *DBaaSClient) Instances() InstanceService
Instances returns a service for managing database instances

```
```
func (c *DBaaSClient) Parameters() ParameterService
Parameters returns a service for managing parameters within parameter groups

```
```
func (c *DBaaSClient) ParametersGroup() ParameterGroupService
ParametersGroup returns a service for managing parameter groups

```
```
func (c *DBaaSClient) Replicas() ReplicaService
Replicas returns a service for managing database replicas

```
```
type DatabaseInstanceUpdateRequest struct {
BackupRetentionDays *int    `json:"backup_retention_days,omitempty"`
BackupStartAt       *string `json:"backup_start_at,omitempty"`
}
DatabaseInstanceUpdateRequest represents the request payload for updating an
instance

```
```
type EngineDetail struct {
ID      string `json:"id"`
Name    string `json:"name"`
Version string `json:"version"`
Status  string `json:"status"`
}
EngineDetail represents a database engine

```
```
type EngineParameterDetail struct {
AllowedValues []string `json:"allowed_values"`
DataType      string   `json:"data_type"`
DefaultValue  string   `json:"default_value"`
Description   string   `json:"description"`
Dynamic       bool     `json:"dynamic"`
EngineID      string   `json:"engine_id"`
Modifiable    bool     `json:"modifiable"`
Name          string   `json:"name"`
ParameterName string   `json:"parameter_name"`
RangedValue   bool     `json:"ranged_value"`
}
EngineParameterDetail represents a parameter of a database engine

```
```
type EngineParametersResponse struct {
Results []EngineParameterDetail `json:"results"`
Meta    MetaResponse            `json:"meta"`
}
EngineParametersResponse represents the response when listing engine
parameters

```
```
type EngineService interface {
List(ctx context.Context, opts ListEngineOptions) ([]EngineDetail, error)
Get(ctx context.Context, id string) (*EngineDetail, error)
ListEngineParameters(ctx context.Context, engineID string, opts ListEngineParametersOptions) ([]EngineParameterDetail, error)
}
EngineService provides methods for managing database engines

```
```
type FieldValueFilter struct {
Field string `json:"field"`
Value string `json:"value"`
}
FieldValueFilter represents a filter applied to the results

```
```
type GetInstanceOptions struct {
ExpandedFields []string
}
GetInstanceOptions provides options for getting instance details

```
```
type InstanceCreateRequest struct {
Name                string                `json:"name"`
EngineID            *string               `json:"engine_id,omitempty"`
InstanceTypeID      *string               `json:"instance_type_id,omitempty"`
User                string                `json:"user"`
Password            string                `json:"password"`
Volume              InstanceVolumeRequest `json:"volume"`
ParameterGroupID    *string               `json:"parameter_group_id,omitempty"`
AvailabilityZone    *string               `json:"availability_zone,omitempty"`
BackupStartAt       *string               `json:"backup_start_at,omitempty"`
BackupRetentionDays *int                  `json:"backup_retention_days,omitempty"`
}
InstanceCreateRequest represents the request payload for creating an
instance

```
```
type InstanceDetail struct {
ID                     string                  `json:"id"`
Name                   string                  `json:"name"`
EngineID               string                  `json:"engine_id"`
DatastoreID            string                  `json:"datastore_id"`
FlavorID               string                  `json:"flavor_id"`
InstanceTypeID         string                  `json:"instance_type_id"`
Volume                 Volume                  `json:"volume"`
Addresses              []Address               `json:"addresses"`
Status                 InstanceStatus          `json:"status"`
Generation             string                  `json:"generation"`
ApplyParametersPending bool                    `json:"apply_parameters_pending"`
ParameterGroupID       string                  `json:"parameter_group_id"`
AvailabilityZone       string                  `json:"availability_zone"`
BackupRetentionDays    int                     `json:"backup_retention_days"`
BackupStartAt          string                  `json:"backup_start_at"`
CreatedAt              string                  `json:"created_at"`
UpdatedAt              *string                 `json:"updated_at,omitempty"`
StartedAt              *string                 `json:"started_at,omitempty"`
FinishedAt             *string                 `json:"finished_at,omitempty"`
MaintenanceScheduledAt *string                 `json:"maintenance_scheduled_at,omitempty"`
Replicas               []ReplicaDetailResponse `json:"replicas,omitempty"`
}
InstanceDetail represents detailed information about an instance

```
```
type InstanceParametersRequest struct {
Name  string `json:"name"`
Value string `json:"value"`
}
InstanceParametersRequest represents a parameter request

```
```
type InstanceParametersResponse struct {
Name  string `json:"name"`
Value string `json:"value"`
}
InstanceParametersResponse represents a parameter response

```
```
type InstanceResizeRequest struct {
InstanceTypeID *string                      `json:"instance_type_id,omitempty"`
Volume         *InstanceVolumeResizeRequest `json:"volume,omitempty"`
}
InstanceResizeRequest represents the request payload for resizing an
instance

```
```
type InstanceResponse struct {
ID string `json:"id"`
}
InstanceResponse represents the response when creating an instance

```
```
type InstanceService interface {
List(ctx context.Context, opts ListInstanceOptions) ([]InstanceDetail, error)
Get(ctx context.Context, id string, opts GetInstanceOptions) (*InstanceDetail, error)
Create(ctx context.Context, req InstanceCreateRequest) (*InstanceResponse, error)
Delete(ctx context.Context, id string) error
Update(ctx context.Context, id string, req DatabaseInstanceUpdateRequest) (*InstanceDetail, error)
Resize(ctx context.Context, id string, req InstanceResizeRequest) (*InstanceDetail, error)
Start(ctx context.Context, id string) (*InstanceDetail, error)
Stop(ctx context.Context, id string) (*InstanceDetail, error)
ListSnapshots(ctx context.Context, instanceID string, opts ListSnapshotOptions) ([]SnapshotDetailResponse, error)
CreateSnapshot(ctx context.Context, instanceID string, req SnapshotCreateRequest) (*SnapshotResponse, error)
GetSnapshot(ctx context.Context, instanceID, snapshotID string) (*SnapshotDetailResponse, error)
UpdateSnapshot(ctx context.Context, instanceID, snapshotID string, req SnapshotUpdateRequest) (*SnapshotDetailResponse, error)
DeleteSnapshot(ctx context.Context, instanceID, snapshotID string) error
RestoreSnapshot(ctx context.Context, instanceID, snapshotID string, req RestoreSnapshotRequest) (*InstanceResponse, error)
}
InstanceService provides methods for managing database instances

```
```
type InstanceStatus string

```
```
const (
InstanceStatusCreating         InstanceStatus = "CREATING"
InstanceStatusError            InstanceStatus = "ERROR"
InstanceStatusStopped          InstanceStatus = "STOPPED"
InstanceStatusReboot           InstanceStatus = "REBOOT"
InstanceStatusPending          InstanceStatus = "PENDING"
InstanceStatusResizing         InstanceStatus = "RESIZING"
InstanceStatusDeleted          InstanceStatus = "DELETED"
InstanceStatusActive           InstanceStatus = "ACTIVE"
InstanceStatusStarting         InstanceStatus = "STARTING"
InstanceStatusStopping         InstanceStatus = "STOPPING"
InstanceStatusBackingUp        InstanceStatus = "BACKING_UP"
InstanceStatusDeleting         InstanceStatus = "DELETING"
InstanceStatusRestoring        InstanceStatus = "RESTORING"
InstanceStatusErrorDeleting    InstanceStatus = "ERROR_DELETING"
InstanceStatusMaintenance      InstanceStatus = "MAINTENANCE"
InstanceStatusMaintenanceError InstanceStatus = "MAINTENANCE_ERROR"
)
```
```
type InstanceStatusUpdate string
InstanceStatusUpdate represents the status update for an instance

```
```
const (
InstanceStatusUpdateActive  InstanceStatusUpdate = "ACTIVE"
InstanceStatusUpdateStopped InstanceStatusUpdate = "STOPPED"
)
```
```
type InstanceType struct {
ID                string `json:"id"`
Name              string `json:"name"`
Label             string `json:"label"`
VCPU              string `json:"vcpu"`
RAM               string `json:"ram"`
FamilyDescription string `json:"family_description"`
FamilySlug        string `json:"family_slug"`
Size              string `json:"size"`
CompatibleProduct string `json:"compatible_product"`
}
InstanceType represents a database instance type

```
```
type InstanceTypeService interface {
List(ctx context.Context, opts ListInstanceTypeOptions) ([]InstanceType, error)
Get(ctx context.Context, id string) (*InstanceType, error)
}
InstanceTypeService provides methods for managing database instance types

```
```
type InstanceVolumeRequest struct {
Size int    `json:"size"`
Type string `json:"type,omitempty"`
}
InstanceVolumeRequest represents volume configuration for instance creation

```
```
type InstanceVolumeResizeRequest struct {
Size int    `json:"size"`
Type string `json:"type,omitempty"`
}
InstanceVolumeResizeRequest represents volume configuration for instance
resizing

```
```
type InstancesResponse struct {
Meta    MetaResponse     `json:"meta"`
Results []InstanceDetail `json:"results"`
}
InstancesResponse represents the response when listing instances

```
```
type ListClustersOptions struct {
Offset           *int
Limit            *int
Status           *ClusterStatus
EngineID         *string
VolumeSize       *int
VolumeSizeGt     *int
VolumeSizeGte    *int
VolumeSizeLt     *int
VolumeSizeLte    *int
ParameterGroupID *string
}
ListClustersOptions provides options for listing clusters

```
```
type ListEngineOptions struct {
Offset *int
Limit  *int
Status *string
}
ListEngineOptions provides options for listing engines

```
```
type ListEngineParametersOptions struct {
Offset     *int
Limit      *int
Dynamic    *bool
Modifiable *bool
}
ListEngineParametersOptions provides options for listing engine parameters

```
```
type ListEnginesResponse struct {
Meta    MetaResponse   `json:"meta"`
Results []EngineDetail `json:"results"`
}
ListEnginesResponse represents the response when listing engines

```
```
type ListInstanceOptions struct {
Offset         *int
Limit          *int
Status         *InstanceStatus
EngineID       *string
VolumeSize     *int
VolumeSizeGt   *int
VolumeSizeGte  *int
VolumeSizeLt   *int
VolumeSizeLte  *int
ExpandedFields []string
}
ListInstanceOptions provides options for listing instances

```
```
type ListInstanceTypeOptions struct {
Offset   *int    `json:"offset,omitempty"`
Limit    *int    `json:"limit,omitempty"`
Status   *string `json:"status,omitempty"`
EngineID *string `json:"engine_id,omitempty"`
}
ListInstanceTypeOptions provides options for listing instance types

```
```
type ListInstanceTypesResponse struct {
Meta    MetaResponse   `json:"meta"`
Results []InstanceType `json:"results"`
}
ListInstanceTypesResponse represents the response when listing instance
types

```
```
type ListParameterGroupsOptions struct {
Offset   *int
Limit    *int
Type     *ParameterGroupType
EngineID *string
}
ListParameterGroupsOptions defines query parameters for listing parameter
groups

```
```
type ListParametersOptions struct {
ParameterGroupID string
Offset           *int
Limit            *int
}
ListParametersOptions provides options for listing parameters

```
```
type ListReplicaOptions struct {
Offset   *int
Limit    *int
SourceID *string
}
ListReplicaOptions provides options for listing replicas

```
```
type ListSnapshotOptions struct {
Offset *int            `json:"offset,omitempty"`
Limit  *int            `json:"limit,omitempty"`
Type   *SnapshotType   `json:"type,omitempty"`
Status *SnapshotStatus `json:"status,omitempty"`
}
ListSnapshotOptions provides options for listing snapshots

```
```
type LoadBalancerAddress struct {
Access  AddressAccess `json:"access"`
Type    AddressType   `json:"type,omitempty"`
Address string        `json:"address,omitempty"`
Port    string        `json:"port,omitempty"`
}
LoadBalancerAddress represents a load balancer address

```
```
type MetaResponse struct {
Page    PageResponse       `json:"page"`
Filters []FieldValueFilter `json:"filters"`
}
MetaResponse contains metadata about the response

```
```
type PageResponse struct {
Offset   int `json:"offset"`
Limit    int `json:"limit"`
Count    int `json:"count"`
Total    int `json:"total"`
MaxLimit int `json:"max_limit"`
}
PageResponse contains pagination details

```
```
type ParameterCreateRequest struct {
Name  string `json:"name"`
Value any    `json:"value"`
}
ParameterCreateRequest represents the request payload for creating a
parameter

```
```
type ParameterDetailResponse struct {
ID    string `json:"id"`
Name  string `json:"name"`
Value any    `json:"value"`
}
ParameterDetailResponse represents detailed information about a parameter

```
```
type ParameterGroupCreateRequest struct {
Name        string  `json:"name"`
EngineID    string  `json:"engine_id"`
Description *string `json:"description,omitempty"`
}
ParameterGroupCreateRequest contains the data for creating a new parameter
group

```
```
type ParameterGroupDetailResponse struct {
ID          string             `json:"id"`
Name        string             `json:"name"`
Description *string            `json:"description,omitempty"`
Type        ParameterGroupType `json:"type"`
EngineID    string             `json:"engine_id"`
}
ParameterGroupDetailResponse represents the detailed view of a parameter
group

```
```
type ParameterGroupResponse struct {
ID string `json:"id"`
}
ParameterGroupResponse contains the ID of a newly created parameter group

```
```
type ParameterGroupService interface {
List(ctx context.Context, opts ListParameterGroupsOptions) ([]ParameterGroupDetailResponse, error)
Create(ctx context.Context, req ParameterGroupCreateRequest) (*ParameterGroupResponse, error)
Get(ctx context.Context, ID string) (*ParameterGroupDetailResponse, error)
Update(ctx context.Context, ID string, req ParameterGroupUpdateRequest) (*ParameterGroupDetailResponse, error)
Delete(ctx context.Context, ID string) error
}
ParameterGroupService defines the interface for parameter group operations

```
```
type ParameterGroupType string
ParameterGroupType defines the type of parameter group

```
```
type ParameterGroupUpdateRequest struct {
Name        *string `json:"name,omitempty"`
Description *string `json:"description,omitempty"`
}
ParameterGroupUpdateRequest contains the fields that can be updated for a
parameter group

```
```
type ParameterGroupsResponse struct {
Meta    MetaResponse                   `json:"meta"`
Results []ParameterGroupDetailResponse `json:"results"`
}
ParameterGroupsResponse represents the API response for multiple parameter
groups

```
```
type ParameterResponse struct {
ID string `json:"id"`
}
ParameterResponse represents the response when creating a parameter

```
```
type ParameterService interface {
List(ctx context.Context, opts ListParametersOptions) ([]ParameterDetailResponse, error)
Create(ctx context.Context, groupID string, req ParameterCreateRequest) (*ParameterResponse, error)
Update(ctx context.Context, groupID, parameterID string, req ParameterUpdateRequest) (*ParameterDetailResponse, error)
Delete(ctx context.Context, groupID, parameterID string) error
}
ParameterService provides methods for managing parameters within parameter
groups

```
```
type ParameterUpdateRequest struct {
Value any `json:"value"`
}
ParameterUpdateRequest represents the request payload for updating a
parameter

```
```
type ParametersResponse struct {
Meta    MetaResponse              `json:"meta"`
Results []ParameterDetailResponse `json:"results"`
}
ParametersResponse represents the response when listing parameters

```
```
type ReplicaAddressResponse struct {
Access  AddressAccess `json:"access"`
Type    *AddressType  `json:"type,omitempty"`
Address *string       `json:"address,omitempty"`
}
ReplicaAddressResponse represents a replica address

```
```
type ReplicaCreateRequest struct {
SourceID       string  `json:"source_id"`
Name           string  `json:"name"`
InstanceTypeID *string `json:"instance_type_id,omitempty"`
}
ReplicaCreateRequest represents the request payload for creating a replica

```
```
type ReplicaDetailResponse struct {
ID                     string                   `json:"id"`
SourceID               string                   `json:"source_id"`
Name                   string                   `json:"name"`
EngineID               string                   `json:"engine_id"`
InstanceTypeID         string                   `json:"instance_type_id"`
Volume                 Volume                   `json:"volume"`
Addresses              []ReplicaAddressResponse `json:"addresses"`
Status                 InstanceStatus           `json:"status"`
Generation             string                   `json:"generation"`
CreatedAt              time.Time                `json:"created_at"`
UpdatedAt              *time.Time               `json:"updated_at,omitempty"`
StartedAt              *string                  `json:"started_at,omitempty"`
FinishedAt             *string                  `json:"finished_at,omitempty"`
MaintenanceScheduledAt *string                  `json:"maintenance_scheduled_at,omitempty"`
}
ReplicaDetailResponse represents detailed information about a replica

```
```
type ReplicaResizeRequest struct {
InstanceTypeID string `json:"instance_type_id,omitempty"`
}
ReplicaResizeRequest represents the request payload for resizing a replica

```
```
type ReplicaResponse struct {
ID string `json:"id"`
}
ReplicaResponse represents the response when creating a replica

```
```
type ReplicaService interface {
List(ctx context.Context, opts ListReplicaOptions) ([]ReplicaDetailResponse, error)
Get(ctx context.Context, id string) (*ReplicaDetailResponse, error)
Create(ctx context.Context, req ReplicaCreateRequest) (*ReplicaResponse, error)
Delete(ctx context.Context, id string) error
Resize(ctx context.Context, id string, req ReplicaResizeRequest) (*ReplicaDetailResponse, error)
Start(ctx context.Context, id string) (*ReplicaDetailResponse, error)
Stop(ctx context.Context, id string) (*ReplicaDetailResponse, error)
}
ReplicaService provides methods for managing database replicas

```
```
func NewReplicaService(client *DBaaSClient) ReplicaService
NewReplicaService creates a new replica service instance

```
```
type ReplicasResponse struct {
Meta    MetaResponse            `json:"meta"`
Results []ReplicaDetailResponse `json:"results"`
}
ReplicasResponse represents the response when listing replicas

```
```
type RestoreSnapshotRequest struct {
Name                string                 `json:"name"`
InstanceTypeID      string                 `json:"instance_type_id"`
Volume              *InstanceVolumeRequest `json:"volume,omitempty"`
BackupRetentionDays *int                   `json:"backup_retention_days,omitempty"`
BackupStartAt       *string                `json:"backup_start_at,omitempty"`
}
RestoreSnapshotRequest represents the request payload for restoring from a
snapshot

```
```
type SnapshotCreateRequest struct {
Name        string  `json:"name"`
Description *string `json:"description,omitempty"`
}
SnapshotCreateRequest represents the request payload for creating a snapshot

```
```
type SnapshotDetailResponse struct {
ID               string                         `json:"id"`
Instance         SnapshotInstanceDetailResponse `json:"instance"`
Name             string                         `json:"name"`
Description      string                         `json:"description"`
Type             SnapshotType                   `json:"type"`
Status           SnapshotStatus                 `json:"status"`
AllocatedSize    int                            `json:"allocated_size"`
CreatedAt        time.Time                      `json:"created_at"`
StartedAt        *string                        `json:"started_at,omitempty"`
AvailabilityZone string                         `json:"availability_zone"`
FinishedAt       *string                        `json:"finished_at,omitempty"`
UpdatedAt        *time.Time                     `json:"updated_at,omitempty"`
}
SnapshotDetailResponse represents detailed information about a snapshot

```
```
type SnapshotInstanceDetailResponse struct {
ID   string `json:"id"`
Name string `json:"name"`
}
SnapshotInstanceDetailResponse represents instance information in a snapshot

```
```
type SnapshotResponse struct {
ID string `json:"id"`
}
SnapshotResponse represents the response when creating a snapshot

```
```
type SnapshotStatus string

```
```
const (
SnapshotStatusPending   SnapshotStatus = "PENDING"
SnapshotStatusCreating  SnapshotStatus = "CREATING"
SnapshotStatusAvailable SnapshotStatus = "AVAILABLE"
SnapshotStatusRestoring SnapshotStatus = "RESTORING"
SnapshotStatusError     SnapshotStatus = "ERROR"
SnapshotStatusDeleting  SnapshotStatus = "DELETING"
SnapshotStatusDeleted   SnapshotStatus = "DELETED"
)
```
```
type SnapshotType string

```
```
const (
SnapshotTypeOnDemand  SnapshotType = "ON_DEMAND"
SnapshotTypeAutomated SnapshotType = "AUTOMATED"
)
```
```
type SnapshotUpdateRequest struct {
Name        string  `json:"name,omitempty"`
Description *string `json:"description,omitempty"`
}
SnapshotUpdateRequest represents the request payload for updating a snapshot

```
```
type SnapshotsResponse struct {
Meta    MetaResponse             `json:"meta"`
Results []SnapshotDetailResponse `json:"results"`
}
SnapshotsResponse represents the response when listing snapshots

```
```
type Volume struct {
Size int    `json:"size"`
Type string `json:"type"`
}
Volume represents volume information for an instance


```

