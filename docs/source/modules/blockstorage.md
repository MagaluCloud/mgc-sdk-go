# Blockstorage

Package blockstorage provides functionality to interact with the MagaluCloud
block storage service. This package allows managing volumes, volume types,
and snapshots.
```
const (
VolumeTypeExpand   = "volume_type"
VolumeAttachExpand = "attachment"
)
VolumeTypeExpand is a constant used for expanding volume type information in
volume responses.

```
```
const (
DefaultBasePath = "/volume"
)
DefaultBasePath defines the default base path for block storage APIs.

```
```
const (
SnapshotVolumeExpand = "volume"
)
SnapshotVolumeExpand is a constant used for expanding volume information in
snapshot responses.



```
```
type AttachmentInstance struct {
ID        *string    `json:"id"`
Name      *string    `json:"name"`
Status    *string    `json:"status"`
State     *string    `json:"state"`
CreatedAt *time.Time `json:"created_at"`
UpdatedAt *time.Time `json:"updated_at"`
}
AttachmentInstance represents information about an instance attached to a
volume.

```
```
type BlockStorageClient struct {
*client.CoreClient
}
BlockStorageClient represents a client for the block storage service. It
encapsulates functionality to access volumes, volume types, and snapshots.

```
```
func New(core *client.CoreClient, opts ...ClientOption) *BlockStorageClient
New creates a new instance of BlockStorageClient. If the core client is nil,
returns nil.

```
```
func (c *BlockStorageClient) Snapshots() SnapshotService
Snapshots returns a service to manage volume snapshots. This method allows
access to functionality such as creating, listing, and managing snapshots.

```
```
func (c *BlockStorageClient) VolumeTypes() VolumeTypeService
VolumeTypes returns a service to manage volume types. This method allows
access to functionality such as listing available volume types.

```
```
func (c *BlockStorageClient) Volumes() VolumeService
Volumes returns a service to manage block storage volumes. This method
allows access to functionality such as creating, listing, and managing
volumes.

```
```
type ClientOption func(*BlockStorageClient)
ClientOption allows customizing the block storage client configuration.

```
```
type CreateSnapshotRequest struct {
Name           string    `json:"name"`
Volume         *IDOrName `json:"volume,omitempty"`
Description    *string   `json:"description"`
Type           *string   `json:"type"`
SourceSnapshot *IDOrName `json:"source_snapshot,omitempty"`
}
CreateSnapshotRequest represents the request to create a new snapshot.

```
```
type CreateVolumeRequest struct {
AvailabilityZone *string   `json:"availability_zone,omitempty"`
Name             string    `json:"name"`
Size             int       `json:"size"`
Type             IDOrName  `json:"type"`
Snapshot         *IDOrName `json:"snapshot,omitempty"`
Encrypted        *bool     `json:"encrypted"`
}
CreateVolumeRequest represents the request to create a new volume.

```
```
type DiskType string
DiskType represents the physical disk type used for storage. Different disk
types offer different performance characteristics and costs.

```
```
const (
DiskTypeNVMe DiskType = "nvme"
DiskTypeHDD  DiskType = "hdd"
)
```
```
type ExtendVolumeRequest struct {
Size int `json:"size"`
}
ExtendVolumeRequest represents the request to extend a volume.

```
```
type IDOrName struct {
ID   *string `json:"id,omitempty"`
Name *string `json:"name,omitempty"`
}
IDOrName represents a reference that can be either an ID or a name.
This structure is used when an API can accept either an ID or a name as a
parameter.

```
```
type Iops struct {
Read  int `json:"read"`
Write int `json:"write"`
Total int `json:"total"`
}
Iops represents the input/output operations per second specifications for a
volume. IOPS defines the performance characteristics in terms of read/write
operations.

```
```
type ListOptions struct {
Limit  *int
Offset *int
Sort   *string
Expand []string
}
ListOptions contains options for listing volumes. All fields are optional
and allow controlling pagination and expansion.

```
```
type ListSnapshotsResponse struct {
Snapshots []Snapshot `json:"snapshots"`
}
ListSnapshotsResponse represents the response from listing snapshots.
This structure encapsulates the API response format for snapshots.

```
```
type ListVolumeTypesOptions struct {
AvailabilityZone string
Name             string
AllowsEncryption *bool
}
ListVolumeTypesOptions contains the options for listing volume types.
All fields are optional and allow filtering the results.

```
```
type ListVolumeTypesResponse struct {
Types []VolumeType `json:"types"`
}
ListVolumeTypesResponse represents the response from listing volume types.
This structure encapsulates the API response format for volume types.

```
```
type ListVolumesResponse struct {
Volumes []Volume `json:"volumes"`
}
ListVolumesResponse represents the response from listing volumes. This
structure encapsulates the API response format for volumes.

```
```
type RenameSnapshotRequest struct {
Name string `json:"name"`
}
RenameSnapshotRequest represents the request to rename a snapshot.

```
```
type RenameVolumeRequest struct {
Name string `json:"name"`
}
RenameVolumeRequest represents the request to rename a volume.

```
```
type RetypeVolumeRequest struct {
NewType IDOrName `json:"new_type"`
}
RetypeVolumeRequest represents the request to change a volume's type.

```
```
type Snapshot struct {
ID                string           `json:"id"`
Name              string           `json:"name"`
Size              int              `json:"size"`
Description       *string          `json:"description"`
State             SnapshotStateV1  `json:"state"`
Status            SnapshotStatusV1 `json:"status"`
CreatedAt         time.Time        `json:"created_at"`
UpdatedAt         time.Time        `json:"updated_at"`
Volume            *IDOrName        `json:"volume,omitempty"`
Error             *SnapshotError   `json:"error,omitempty"`
AvailabilityZones []string         `json:"availability_zones"`
Type              string           `json:"type"`
}
Snapshot represents a volume snapshot. A snapshot is a point-in-time copy of
a volume that can be used for backup or to create new volumes.

```
```
type SnapshotError struct {
Slug    string `json:"slug"`
Message string `json:"message"`
}
SnapshotError represents error information for a snapshot operation.

```
```
type SnapshotService interface {
List(ctx context.Context, opts ListOptions) ([]Snapshot, error)
Create(ctx context.Context, req CreateSnapshotRequest) (string, error)
Get(ctx context.Context, id string, expand []string) (*Snapshot, error)
Delete(ctx context.Context, id string) error
Rename(ctx context.Context, id string, newName string) error
}
SnapshotService provides operations for managing volume snapshots. This
interface allows creating, listing, retrieving, and managing snapshots.

```
```
type SnapshotStateV1 string
SnapshotStateV1 represents the possible states of a snapshot. The state
indicates the lifecycle stage of the snapshot.

```
```
const (
SnapshotStateNew       SnapshotStateV1 = "new"
SnapshotStateAvailable SnapshotStateV1 = "available"
SnapshotStateDeleted   SnapshotStateV1 = "deleted"
)
```
```
type SnapshotStatusV1 string
SnapshotStatusV1 represents the possible statuses of a snapshot. The status
provides more detailed information about the snapshot's current condition.

```
```
const (
SnapshotStatusProvisioning       SnapshotStatusV1 = "provisioning"
SnapshotStatusCreating           SnapshotStatusV1 = "creating"
SnapshotStatusCreatingError      SnapshotStatusV1 = "creating_error"
SnapshotStatusCreatingErrorQuota SnapshotStatusV1 = "creating_error_quota"
SnapshotStatusCompleted          SnapshotStatusV1 = "completed"
SnapshotStatusDeleting           SnapshotStatusV1 = "deleting"
SnapshotStatusDeleted            SnapshotStatusV1 = "deleted"
SnapshotStatusDeletedError       SnapshotStatusV1 = "deleted_error"
SnapshotStatusReplicating        SnapshotStatusV1 = "replicating"
SnapshotStatusReplicatingError   SnapshotStatusV1 = "replicating_error"
SnapshotStatusRestoring          SnapshotStatusV1 = "restoring"
SnapshotStatusRestoringError     SnapshotStatusV1 = "restoring_error"
SnapshotStatusReserved           SnapshotStatusV1 = "reserved"
)
```
```
type Type struct {
Iops     *Iops   `json:"iops,omitempty"`
ID       string  `json:"id"`
Name     *string `json:"name,omitempty"`
DiskType *string `json:"disk_type,omitempty"`
Status   *string `json:"status,omitempty"`
}
Type represents the volume type information. Contains details about the
volume type including IOPS specifications.

```
```
type Volume struct {
ID                string            `json:"id"`
Name              string            `json:"name"`
Size              int               `json:"size"`
Status            string            `json:"status"`
State             string            `json:"state"`
CreatedAt         time.Time         `json:"created_at"`
UpdatedAt         time.Time         `json:"updated_at"`
Type              Type              `json:"type"`
Error             *VolumeError      `json:"error,omitempty"`
Attachment        *VolumeAttachment `json:"attachment,omitempty"`
AvailabilityZone  string            `json:"availability_zone"`
AvailabilityZones []string          `json:"availability_zones"`
Encrypted         *bool             `json:"encrypted,omitempty"`
}
Volume represents a block storage volume. A volume is a persistent block
storage device that can be attached to instances.

```
```
type VolumeAttachment struct {
Instance   AttachmentInstance `json:"instance"`
AttachedAt time.Time          `json:"attached_at"`
Device     *string            `json:"device,omitempty"`
}
VolumeAttachment represents the attachment of a volume to an instance.

```
```
type VolumeError struct {
Slug    string `json:"slug"`
Message string `json:"message"`
}
VolumeError represents error information for a volume operation.

```
```
type VolumeService interface {
List(ctx context.Context, opts ListOptions) ([]Volume, error)
Create(ctx context.Context, req CreateVolumeRequest) (string, error)
Get(ctx context.Context, id string, expand []string) (*Volume, error)
Delete(ctx context.Context, id string) error
Rename(ctx context.Context, id string, newName string) error
Extend(ctx context.Context, id string, req ExtendVolumeRequest) error
Retype(ctx context.Context, id string, req RetypeVolumeRequest) error
Attach(ctx context.Context, volumeID string, instanceID string) error
Detach(ctx context.Context, volumeID string) error
}
VolumeService defines the interface for volume operations. This interface
provides methods for managing block storage volumes.

```
```
type VolumeStateV1 string
VolumeStateV1 represents the possible states of a volume. The state
indicates the lifecycle stage of the volume.

```
```
const (
VolumeStateNew       VolumeStateV1 = "new"
VolumeStateAvailable VolumeStateV1 = "available"
VolumeStateInUse     VolumeStateV1 = "in-use"
VolumeStateDeleted   VolumeStateV1 = "deleted"
VolumeStateLegacy    VolumeStateV1 = "legacy"
)
```
```
type VolumeStatusV1 string
VolumeStatusV1 represents the possible statuses of a volume. The status
provides more detailed information about the volume's current condition.

```
```
const (
VolumeStatusProvisioning VolumeStatusV1 = "provisioning"
VolumeStatusCreating     VolumeStatusV1 = "creating"
VolumeStatusAvailable    VolumeStatusV1 = "available"
VolumeStatusAttaching    VolumeStatusV1 = "attaching"
VolumeStatusInUse        VolumeStatusV1 = "in-use"
VolumeStatusDetaching    VolumeStatusV1 = "detaching"
VolumeStatusDeleting     VolumeStatusV1 = "deleting"
VolumeStatusError        VolumeStatusV1 = "error"
VolumeStatusLegacy       VolumeStatusV1 = "legacy"
)
```
```
type VolumeType struct {
ID                string         `json:"id"`
Name              string         `json:"name"`
DiskType          string         `json:"disk_type"`
Status            string         `json:"status"`
IOPS              VolumeTypeIOPS `json:"iops"`
AvailabilityZones []string       `json:"availability_zones"`
AllowsEncryption  bool           `json:"allows_encryption"`
}
VolumeType represents a block storage volume type. Each volume type defines
the characteristics and capabilities of volumes created with it.

```
```
type VolumeTypeIOPS struct {
Read  int `json:"read"`
Write int `json:"write"`
Total int `json:"total"`
}
VolumeTypeIOPS represents the IOPS specifications for a volume type. IOPS
defines the performance characteristics in terms of read/write operations.

```
```
type VolumeTypeService interface {
List(ctx context.Context, opts ListVolumeTypesOptions) ([]VolumeType, error)
}
VolumeTypeService provides operations for managing volume types. This
interface allows listing available volume types with optional filtering.


```

