# Compute

Package compute provides functionality to interact with the MagaluCloud compute
service. This package allows managing virtual machine instances, images,
instance types, and snapshots.
```
const (
InstanceImageExpand       = "image"
InstanceMachineTypeExpand = "machine-type"
InstanceNetworkExpand     = "network"
)
Constants for expanding related resources in instance responses.

```
```
const (
VmInstanceHeaderVersionName = "x-api-version"
VmInstanceHeaderVersion     = "1.1"
)
Constants for API version headers.

```
```
const (
// SnapshotImageExpand is used to include image information in snapshot responses
SnapshotImageExpand = "image"
// SnapshotMachineTypeExpand is used to include machine type information in snapshot responses
SnapshotMachineTypeExpand = "machine-type"
)
Constants for expanding related resources in snapshot responses.

```
```
const (
DefaultBasePath = "/compute"
)


```
```
type ClientOption func(*VirtualMachineClient)
ClientOption allows customizing the virtual machine client configuration.

```
```
type CopySnapshotRequest struct {
// DestinationRegion is the region where the snapshot should be copied
DestinationRegion string `json:"destination_region"`
}
CopySnapshotRequest represents the request to copy a snapshot to another
region.

```
```
type CreateParametersNetwork struct {
AssociatePublicIp *bool                             `json:"associate_public_ip,omitempty"`
Interface         *CreateParametersNetworkInterface `json:"interface,omitempty"`
Vpc               *IDOrName                         `json:"vpc,omitempty"`
}
CreateParametersNetwork represents network configuration for instance
creation.

```
```
type CreateParametersNetworkInterface struct {
Interface      *IDOrName                                             `json:"interface,omitempty"`
SecurityGroups *[]CreateParametersNetworkInterfaceSecurityGroupsItem `json:"security_groups,omitempty"`
}
CreateParametersNetworkInterface represents network interface configuration.

```
```
type CreateParametersNetworkInterfaceSecurityGroupsItem struct {
Id string `json:"id"`
}
CreateParametersNetworkInterfaceSecurityGroupsItem represents a security
group item.

```
```
type CreateRequest struct {
AvailabilityZone *string                  `json:"availability_zone,omitempty"`
Image            IDOrName                 `json:"image"`
Labels           *[]string                `json:"labels,omitempty"`
MachineType      IDOrName                 `json:"machine_type"`
Name             string                   `json:"name"`
Network          *CreateParametersNetwork `json:"network,omitempty"`
SshKeyName       *string                  `json:"ssh_key_name,omitempty"`
UserData         *string                  `json:"user_data,omitempty"`
}
CreateRequest represents the request to create a new instance.

```
```
type CreateSnapshotRequest struct {
Name     string   `json:"name"`
Instance IDOrName `json:"instance"`
}
CreateSnapshotRequest represents the request to create a new snapshot.

```
```
type Error struct {
Message string `json:"message"`
Slug    string `json:"slug"`
}
Error represents an error that occurred with an instance.

```
```
type IDOrName struct {
ID   *string `json:"id,omitempty,omitzero"`
Name *string `json:"name,omitempty,omitzero"`
}
IDOrName represents a resource that can be identified by ID or name.

```
```
type Image struct {
ID                   string              `json:"id"`
Name                 string              `json:"name"`
Status               ImageStatus         `json:"status"`
Version              *string             `json:"version,omitempty"`
Platform             *string             `json:"platform,omitempty"`
ReleaseAt            *string             `json:"release_at,omitempty"`
EndStandardSupportAt *string             `json:"end_standard_support_at,omitempty"`
EndLifeAt            *string             `json:"end_life_at,omitempty"`
MinimumRequirements  MinimumRequirements `json:"minimum_requirements"`
Labels               *[]string           `json:"labels,omitempty"`
AvailabilityZones    *[]string           `json:"availability_zones,omitempty"`
}
Image represents a virtual machine image. An image is a template that
contains the operating system and software for creating instances.

```
```
type ImageList struct {
Images []Image `json:"images"`
}
ImageList represents the response from listing images. This structure
encapsulates the API response format for images.

```
```
type ImageListOptions struct {
Limit            *int
Offset           *int
Sort             *string
Labels           []string
AvailabilityZone *string
}
ImageListOptions defines the parameters for filtering and pagination of
image lists. All fields are optional and allow controlling the listing
behavior.

```
```
type ImageService interface {
List(ctx context.Context, opts ImageListOptions) ([]Image, error)
}
ImageService provides operations for managing virtual machine images.
This interface allows listing available images with optional filtering.

```
```
type ImageStatus string
ImageStatus represents the current state of an image. The status indicates
the lifecycle stage and availability of the image.

```
```
const (
ImageStatusActive        ImageStatus = "active"
ImageStatusDeprecated    ImageStatus = "deprecated"
ImageStatusDeleted       ImageStatus = "deleted"
ImageStatusPending       ImageStatus = "pending"
ImageStatusCreating      ImageStatus = "creating"
ImageStatusImporting     ImageStatus = "importing"
ImageStatusError         ImageStatus = "error"
ImageStatusDeletingError ImageStatus = "deleting_error"
)
```
```
type InitLogResponse struct {
Logs []string `json:"logs"`
}
InitLogResponse represents the response from getting instance initialization
logs.

```
```
type Instance struct {
ID               string         `json:"id"`
Name             *string        `json:"name,omitempty"`
MachineType      *InstanceTypes `json:"machine_type"`
Image            *VmImage       `json:"image"`
Status           string         `json:"status"`
State            string         `json:"state"`
CreatedAt        time.Time      `json:"created_at"`
UpdatedAt        *time.Time     `json:"updated_at,omitempty,omitzero"`
SSHKeyName       *string        `json:"ssh_key_name,omitempty"`
AvailabilityZone *string        `json:"availability_zone,omitempty,omitzero"`
Network          *Network       `json:"network"`
UserData         *string        `json:"user_data,omitempty"`
Labels           *[]string      `json:"labels"`
Error            *Error         `json:"error,omitempty"`
}
Instance represents a virtual machine instance.

```
```
type InstanceService interface {
List(ctx context.Context, opts ListOptions) ([]Instance, error)
Create(ctx context.Context, req CreateRequest) (string, error)
Get(ctx context.Context, id string, expand []string) (*Instance, error)
Delete(ctx context.Context, id string, deletePublicIP bool) error
Rename(ctx context.Context, id string, newName string) error
Retype(ctx context.Context, id string, req RetypeRequest) error
Start(ctx context.Context, id string) error
Stop(ctx context.Context, id string) error
Suspend(ctx context.Context, id string) error
GetFirstWindowsPassword(ctx context.Context, id string) (*WindowsPasswordResponse, error)
AttachNetworkInterface(ctx context.Context, req NICRequest) error
DetachNetworkInterface(ctx context.Context, req NICRequest) error
InitLog(ctx context.Context, id string, maxLines *int) (*InitLogResponse, error)
}
InstanceService provides operations for managing virtual machine instances.

```
```
type InstanceType struct {
ID                string    `json:"id"`
Name              string    `json:"name"`
VCPUs             int       `json:"vcpus"`
RAM               int       `json:"ram"`
Disk              int       `json:"disk"`
GPU               *int      `json:"gpu,omitempty"`
Status            string    `json:"status"`
AvailabilityZones *[]string `json:"availability_zones,omitempty"`
}
InstanceType represents a virtual machine instance type configuration.
Each instance type defines the hardware specifications for virtual machines.

```
```
type InstanceTypeList struct {
InstanceTypes []InstanceType `json:"instance_types"`
Meta          Meta           `json:"meta"`
}
InstanceTypeList represents the response from listing instance types.
This structure encapsulates the API response format for instance types.

```
```
type InstanceTypeListOptions struct {
Limit            *int    `url:"_limit,omitempty"`
Offset           *int    `url:"_offset,omitempty"`
Sort             *string `url:"_sort,omitempty"`
AvailabilityZone string  `url:"availability-zone,omitempty"`
}
InstanceTypeListOptions defines parameters for filtering and pagination
of machine type lists. All fields are optional and allow controlling the
listing behavior.

```
```
type InstanceTypeService interface {
List(ctx context.Context, opts InstanceTypeListOptions) ([]InstanceType, error)
}
InstanceTypeService provides operations for querying available machine
types. This interface allows listing instance types with optional filtering.

```
```
type InstanceTypes struct {
ID    string  `json:"id"`
Name  *string `json:"name"`
Vcpus *int    `json:"vcpus"`
Ram   *int    `json:"ram"`
Disk  *int    `json:"disk"`
}
InstanceTypes represents the machine type configuration of an instance.

```
```
type IpAddressNewExpand struct {
PrivateIpv4 string `json:"private_ipv4"`
PublicIpv6  string `json:"public_ipv6,omitempty"`
}
IpAddressNewExpand represents IP address information for network interfaces.

```
```
type ListInstancesResponse struct {
Instances []Instance `json:"instances"`
}
ListInstancesResponse represents the response from listing instances.

```
```
type ListOptions struct {
Limit  *int
Offset *int
Sort   *string
Expand []string
Name   *string
}
ListOptions defines the parameters for filtering and pagination of instance
lists.

```
```
type ListSnapshotsResponse struct {
Snapshots []Snapshot `json:"snapshots"`
}
ListSnapshotsResponse represents the response from listing snapshots.
This structure encapsulates the API response format for snapshots.

```
```
type Meta struct {
Limit  int `json:"limit"`
Offset int `json:"offset"`
Count  int `json:"count"`
Total  int `json:"total"`
}
Meta contains pagination metadata for API responses. This structure provides
information about the current page and total results.

```
```
type MinimumRequirements struct {
VCPU int `json:"vcpu"`
RAM  int `json:"ram"`
Disk int `json:"disk"`
}
MinimumRequirements represents the minimum hardware requirements for an
image. These requirements must be met by the instance type when creating
instances from this image.

```
```
type NICRequest struct {
Instance IDOrName            `json:"instance"`
Network  NICRequestInterface `json:"network"`
}
NICRequest represents the request to attach or detach a network interface.

```
```
type NICRequestInterface struct {
Interface IDOrName `json:"interface"`
}
NICRequestInterface represents network interface configuration for NIC
operations.

```
```
type Network struct {
Vpc        *IDOrName           `json:"vpc,omitempty"`
Interfaces *[]NetworkInterface `json:"interfaces,omitempty"`
}
Network represents the network configuration of an instance.

```
```
type NetworkInterface struct {
ID                   string             `json:"id"`
Name                 string             `json:"name"`
SecurityGroups       *[]string          `json:"security_groups"`
Primary              *bool              `json:"primary"`
AssociatedPublicIpv4 *string            `json:"associated_public_ipv4,omitempty"`
IpAddresses          IpAddressNewExpand `json:"ip_addresses"`
}
NetworkInterface represents a network interface attached to an instance.

```
```
type RestoreSnapshotRequest struct {
Name             string                   `json:"name"`
MachineType      IDOrName                 `json:"machine_type"`
SSHKeyName       *string                  `json:"ssh_key_name,omitempty"`
AvailabilityZone *string                  `json:"availability_zone,omitempty"`
Network          *CreateParametersNetwork `json:"network,omitempty"`
UserData         *string                  `json:"user_data,omitempty"`
}
RestoreSnapshotRequest represents the request to restore an instance from a
snapshot.

```
```
type RetypeRequest struct {
MachineType IDOrName `json:"machine_type"`
}
RetypeRequest represents the request to change an instance's machine type.

```
```
type Snapshot struct {
ID        string            `json:"id"`
Name      string            `json:"name,omitempty"`
Status    string            `json:"status"`
State     string            `json:"state"`
CreatedAt time.Time         `json:"created_at"`
UpdatedAt *time.Time        `json:"updated_at,omitempty"`
Size      int               `json:"size"`
Instance  *SnapshotInstance `json:"instance"`
}
Snapshot represents an instance snapshot. A snapshot is a point-in-time copy
of an instance that can be used for backup or to create new instances.

```
```
type SnapshotInstance struct {
ID          string    `json:"id"`
Image       *IDOrName `json:"image,omitempty"`
MachineType *IDOrName `json:"machine_type,omitempty"`
}
SnapshotInstance represents information about the instance that was
snapshotted.

```
```
type SnapshotService interface {
List(ctx context.Context, opts ListOptions) ([]Snapshot, error)
Create(ctx context.Context, req CreateSnapshotRequest) (string, error)
Get(ctx context.Context, id string, expand []string) (*Snapshot, error)
Delete(ctx context.Context, id string) error
Rename(ctx context.Context, id string, newName string) error
Restore(ctx context.Context, id string, req RestoreSnapshotRequest) (string, error)
Copy(ctx context.Context, id string, req CopySnapshotRequest) error
}
SnapshotService provides operations for managing snapshots. This interface
allows creating, listing, retrieving, and managing instance snapshots.

```
```
type UpdateNameRequest struct {
Name string `json:"name"`
}
UpdateNameRequest represents the request to update an instance name.

```
```
type VirtualMachineClient struct {
*client.CoreClient
}
VirtualMachineClient represents a client for the compute service.
It encapsulates functionality to access instances, images, instance types,
and snapshots.

```
```
func New(core *client.CoreClient, opts ...ClientOption) *VirtualMachineClient
New creates a new instance of VirtualMachineClient. If the core client is
nil, returns nil.

```
```
func (c *VirtualMachineClient) Images() ImageService
Images returns a service to manage virtual machine images. This method
allows access to functionality such as listing available images.

```
```
func (c *VirtualMachineClient) InstanceTypes() InstanceTypeService
InstanceTypes returns a service to manage instance types. This method allows
access to functionality such as listing available machine types.

```
```
func (c *VirtualMachineClient) Instances() InstanceService
Instances returns a service to manage virtual machine instances. This method
allows access to functionality such as creating, listing, and managing
instances.

```
```
func (c *VirtualMachineClient) Snapshots() SnapshotService
Snapshots returns a service to manage instance snapshots. This method allows
access to functionality such as creating, listing, and managing snapshots.

```
```
type VmImage struct {
ID       string  `json:"id"`
Name     *string `json:"name"`
Platform *string `json:"platform,omitempty,omitzero"`
}
VmImage represents the image configuration of an instance.

```
```
type WindowsPasswordInstance struct {
ID        string    `json:"id"`
Password  string    `json:"password"`
CreatedAt time.Time `json:"created_at"`
User      string    `json:"user,omitempty"`
}
WindowsPasswordInstance represents Windows password information.

```
```
type WindowsPasswordResponse struct {
Instance WindowsPasswordInstance `json:"instance"`
}
WindowsPasswordResponse represents the response from getting Windows
password.


```

