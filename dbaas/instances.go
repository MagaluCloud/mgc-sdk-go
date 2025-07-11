package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	// InstancePath is the base path for instance operations
	InstancePath = "/v2/instances"
	// InstancePathID is the path template for instance operations with ID
	InstancePathID = InstancePath + "/%s"
	// SnapshotPath is the path template for snapshot operations
	SnapshotPath = InstancePathID + "/snapshots"
	// SnapshotPathID is the path template for snapshot operations with ID
	SnapshotPathID = SnapshotPath + "/%s"
)

type (
	// InstanceStatus represents the possible states of a database instance
	InstanceStatus string
	// AddressAccess represents the access type for network addresses
	AddressAccess string
	// AddressType represents the type of network address
	AddressType string
	// SnapshotType represents the type of snapshot
	SnapshotType string
	// SnapshotStatus represents the possible states of a snapshot
	SnapshotStatus string
)

const (
	// InstanceStatusCreating indicates the instance is being created
	InstanceStatusCreating InstanceStatus = "CREATING"
	// InstanceStatusError indicates the instance is in an error state
	InstanceStatusError InstanceStatus = "ERROR"
	// InstanceStatusStopped indicates the instance is stopped
	InstanceStatusStopped InstanceStatus = "STOPPED"
	// InstanceStatusReboot indicates the instance is rebooting
	InstanceStatusReboot InstanceStatus = "REBOOT"
	// InstanceStatusPending indicates the instance is waiting to be processed
	InstanceStatusPending InstanceStatus = "PENDING"
	// InstanceStatusResizing indicates the instance is being resized
	InstanceStatusResizing InstanceStatus = "RESIZING"
	// InstanceStatusDeleted indicates the instance has been deleted
	InstanceStatusDeleted InstanceStatus = "DELETED"
	// InstanceStatusActive indicates the instance is running normally
	InstanceStatusActive InstanceStatus = "ACTIVE"
	// InstanceStatusStarting indicates the instance is starting up
	InstanceStatusStarting InstanceStatus = "STARTING"
	// InstanceStatusStopping indicates the instance is shutting down
	InstanceStatusStopping InstanceStatus = "STOPPING"
	// InstanceStatusBackingUp indicates the instance is performing a backup
	InstanceStatusBackingUp InstanceStatus = "BACKING_UP"
	// InstanceStatusDeleting indicates the instance is being deleted
	InstanceStatusDeleting InstanceStatus = "DELETING"
	// InstanceStatusRestoring indicates the instance is being restored
	InstanceStatusRestoring InstanceStatus = "RESTORING"
	// InstanceStatusErrorDeleting indicates an error occurred during deletion
	InstanceStatusErrorDeleting InstanceStatus = "ERROR_DELETING"
	// InstanceStatusMaintenance indicates the instance is under maintenance
	InstanceStatusMaintenance InstanceStatus = "MAINTENANCE"
	// InstanceStatusMaintenanceError indicates an error occurred during maintenance
	InstanceStatusMaintenanceError InstanceStatus = "MAINTENANCE_ERROR"
)

const (
	// AddressAccessPrivate indicates private network access
	AddressAccessPrivate AddressAccess = "PRIVATE"
	// AddressAccessPublic indicates public network access
	AddressAccessPublic AddressAccess = "PUBLIC"
)

const (
	// AddressTypeIPv4 indicates IPv4 address type
	AddressTypeIPv4 AddressType = "IPv4"
	// AddressTypeIPv6 indicates IPv6 address type
	AddressTypeIPv6 AddressType = "IPv6"
)

const (
	// SnapshotTypeOnDemand indicates a manually created snapshot
	SnapshotTypeOnDemand SnapshotType = "ON_DEMAND"
	// SnapshotTypeAutomated indicates an automatically created snapshot
	SnapshotTypeAutomated SnapshotType = "AUTOMATED"
)

const (
	// SnapshotStatusPending indicates the snapshot is waiting to be processed
	SnapshotStatusPending SnapshotStatus = "PENDING"
	// SnapshotStatusCreating indicates the snapshot is being created
	SnapshotStatusCreating SnapshotStatus = "CREATING"
	// SnapshotStatusAvailable indicates the snapshot is available for use
	SnapshotStatusAvailable SnapshotStatus = "AVAILABLE"
	// SnapshotStatusRestoring indicates the snapshot is being restored
	SnapshotStatusRestoring SnapshotStatus = "RESTORING"
	// SnapshotStatusError indicates the snapshot is in an error state
	SnapshotStatusError SnapshotStatus = "ERROR"
	// SnapshotStatusDeleting indicates the snapshot is being deleted
	SnapshotStatusDeleting SnapshotStatus = "DELETING"
	// SnapshotStatusDeleted indicates the snapshot has been deleted
	SnapshotStatusDeleted SnapshotStatus = "DELETED"
)

type (
	// Volume represents volume information for an instance
	Volume struct {
		// Size is the volume size in GB
		Size int `json:"size"`
		// Type is the volume type
		Type string `json:"type"`
	}

	// Address represents a network address for an instance
	Address struct {
		// Access indicates the access type (PRIVATE or PUBLIC)
		Access AddressAccess `json:"access"`
		// Type indicates the address type (IPv4 or IPv6)
		Type *AddressType `json:"type,omitempty"`
		// Address is the IP address
		Address *string `json:"address,omitempty"`
	}

	// InstanceParametersResponse represents a parameter response
	InstanceParametersResponse struct {
		// Name is the parameter name
		Name string `json:"name"`
		// Value is the parameter value
		Value string `json:"value"`
	}

	// InstanceParametersRequest represents a parameter request
	InstanceParametersRequest struct {
		// Name is the parameter name
		Name string `json:"name"`
		// Value is the parameter value
		Value string `json:"value"`
	}

	// InstanceVolumeRequest represents volume configuration for instance creation
	InstanceVolumeRequest struct {
		// Size is the volume size in GB
		Size int `json:"size"`
		// Type is the volume type (optional)
		Type string `json:"type,omitempty"`
	}

	// InstanceVolumeResizeRequest represents volume configuration for instance resizing
	InstanceVolumeResizeRequest struct {
		// Size is the new volume size in GB
		Size int `json:"size"`
		// Type is the volume type (optional)
		Type string `json:"type,omitempty"`
	}

	// InstanceCreateRequest represents the request payload for creating an instance
	InstanceCreateRequest struct {
		// Name is the name of the instance
		Name string `json:"name"`
		// EngineID is the ID of the database engine (optional)
		EngineID *string `json:"engine_id,omitempty"`
		// InstanceTypeID is the ID of the instance type (optional)
		InstanceTypeID *string `json:"instance_type_id,omitempty"`
		// User is the database user
		User string `json:"user"`
		// Password is the database password
		Password string `json:"password"`
		// Volume contains volume configuration
		Volume InstanceVolumeRequest `json:"volume"`
		// ParameterGroupID is the ID of the parameter group (optional)
		ParameterGroupID *string `json:"parameter_group_id,omitempty"`
		// AvailabilityZone is the availability zone (optional)
		AvailabilityZone *string `json:"availability_zone,omitempty"`
		// BackupStartAt is the backup start time (optional)
		BackupStartAt *string `json:"backup_start_at,omitempty"`
		// BackupRetentionDays is the number of days to retain backups (optional)
		BackupRetentionDays *int `json:"backup_retention_days,omitempty"`
	}

	// InstanceResizeRequest represents the request payload for resizing an instance
	InstanceResizeRequest struct {
		// InstanceTypeID is the new instance type ID (optional)
		InstanceTypeID *string `json:"instance_type_id,omitempty"`
		// Volume contains new volume configuration (optional)
		Volume *InstanceVolumeResizeRequest `json:"volume,omitempty"`
	}

	// DatabaseInstanceUpdateRequest represents the request payload for updating an instance
	DatabaseInstanceUpdateRequest struct {
		// BackupRetentionDays is the new backup retention days (optional)
		BackupRetentionDays *int `json:"backup_retention_days,omitempty"`
		// BackupStartAt is the new backup start time (optional)
		BackupStartAt *string `json:"backup_start_at,omitempty"`
	}

	// ReplicaAddressResponse represents a replica address
	ReplicaAddressResponse struct {
		// Access indicates the access type (PRIVATE or PUBLIC)
		Access AddressAccess `json:"access"`
		// Type indicates the address type (IPv4 or IPv6)
		Type *AddressType `json:"type,omitempty"`
		// Address is the IP address
		Address *string `json:"address,omitempty"`
	}

	// SnapshotsResponse represents the response when listing snapshots
	SnapshotsResponse struct {
		// Meta contains pagination and filter information
		Meta MetaResponse `json:"meta"`
		// Results is the list of snapshots
		Results []SnapshotDetailResponse `json:"results"`
	}

	// SnapshotDetailResponse represents detailed information about a snapshot
	SnapshotDetailResponse struct {
		// ID is the unique identifier of the snapshot
		ID string `json:"id"`
		// Instance contains information about the source instance
		Instance SnapshotInstanceDetailResponse `json:"instance"`
		// Name is the name of the snapshot
		Name string `json:"name"`
		// Description is the description of the snapshot
		Description string `json:"description"`
		// Type is the type of snapshot
		Type SnapshotType `json:"type"`
		// Status is the current status of the snapshot
		Status SnapshotStatus `json:"status"`
		// AllocatedSize is the allocated size in bytes
		AllocatedSize int `json:"allocated_size"`
		// CreatedAt is the timestamp when the snapshot was created
		CreatedAt time.Time `json:"created_at"`
		// StartedAt is the timestamp when the snapshot creation started
		StartedAt *string `json:"started_at,omitempty"`
		// AvailabilityZone is the availability zone of the snapshot
		AvailabilityZone string `json:"availability_zone"`
		// FinishedAt is the timestamp when the snapshot creation finished
		FinishedAt *string `json:"finished_at,omitempty"`
		// UpdatedAt is the timestamp when the snapshot was last updated
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
	}

	// SnapshotInstanceDetailResponse represents instance information in a snapshot
	SnapshotInstanceDetailResponse struct {
		// ID is the ID of the source instance
		ID string `json:"id"`
		// Name is the name of the source instance
		Name string `json:"name"`
	}

	// SnapshotCreateRequest represents the request payload for creating a snapshot
	SnapshotCreateRequest struct {
		// Name is the name of the snapshot
		Name string `json:"name"`
		// Description is the description of the snapshot (optional)
		Description *string `json:"description,omitempty"`
	}

	// SnapshotUpdateRequest represents the request payload for updating a snapshot
	SnapshotUpdateRequest struct {
		// Name is the new name of the snapshot (optional)
		Name string `json:"name,omitempty"`
		// Description is the new description of the snapshot (optional)
		Description *string `json:"description,omitempty"`
	}

	// SnapshotResponse represents the response when creating a snapshot
	SnapshotResponse struct {
		// ID is the unique identifier of the created snapshot
		ID string `json:"id"`
	}

	// RestoreSnapshotRequest represents the request payload for restoring from a snapshot
	RestoreSnapshotRequest struct {
		// Name is the name of the new instance
		Name string `json:"name"`
		// InstanceTypeID is the ID of the instance type for the new instance
		InstanceTypeID string `json:"instance_type_id"`
		// Volume contains volume configuration for the new instance (optional)
		Volume *InstanceVolumeRequest `json:"volume,omitempty"`
		// BackupRetentionDays is the number of days to retain backups (optional)
		BackupRetentionDays *int `json:"backup_retention_days,omitempty"`
		// BackupStartAt is the backup start time (optional)
		BackupStartAt *string `json:"backup_start_at,omitempty"`
	}

	// ListSnapshotOptions provides options for listing snapshots
	ListSnapshotOptions struct {
		// Offset is the number of snapshots to skip
		Offset *int `json:"offset,omitempty"`
		// Limit is the maximum number of snapshots to return
		Limit *int `json:"limit,omitempty"`
		// Type filters snapshots by type
		Type *SnapshotType `json:"type,omitempty"`
		// Status filters snapshots by status
		Status *SnapshotStatus `json:"status,omitempty"`
	}
)

// InstanceStatusUpdate represents the status update for an instance
type InstanceStatusUpdate string

const (
	// InstanceStatusUpdateActive indicates the instance should be started
	InstanceStatusUpdateActive InstanceStatusUpdate = "ACTIVE"
	// InstanceStatusUpdateStopped indicates the instance should be stopped
	InstanceStatusUpdateStopped InstanceStatusUpdate = "STOPPED"
)

type (
	// InstanceService provides methods for managing database instances
	InstanceService interface {
		// List returns a list of database instances for a x-tenant-id.
		// It supports pagination and filtering by status, engine_id, and volume size.
		List(ctx context.Context, opts ListInstanceOptions) ([]InstanceDetail, error)

		// Get returns a database instance detail by its ID.
		// Supports expanding additional fields through the options parameter.
		Get(ctx context.Context, id string, opts GetInstanceOptions) (*InstanceDetail, error)

		// Create creates a new database instance asynchronously for a tenant.
		// Returns the ID of the created instance.
		Create(ctx context.Context, req InstanceCreateRequest) (*InstanceResponse, error)

		// Delete deletes a database instance asynchronously.
		Delete(ctx context.Context, id string) error

		// Update updates a database instance's properties.
		// Supports updating status, backup retention days, and backup start time.
		Update(ctx context.Context, id string, req DatabaseInstanceUpdateRequest) (*InstanceDetail, error)

		// Resize changes the instance type and/or volume size of a database instance.
		Resize(ctx context.Context, id string, req InstanceResizeRequest) (*InstanceDetail, error)

		// Start initiates a stopped database instance.
		Start(ctx context.Context, id string) (*InstanceDetail, error)

		// Stop stops a running database instance.
		Stop(ctx context.Context, id string) (*InstanceDetail, error)

		// ListSnapshots returns a list of snapshots for a specific instance.
		ListSnapshots(ctx context.Context, instanceID string, opts ListSnapshotOptions) ([]SnapshotDetailResponse, error)

		// CreateSnapshot creates a new snapshot for the specified instance.
		CreateSnapshot(ctx context.Context, instanceID string, req SnapshotCreateRequest) (*SnapshotResponse, error)

		// GetSnapshot retrieves details of a specific snapshot.
		GetSnapshot(ctx context.Context, instanceID, snapshotID string) (*SnapshotDetailResponse, error)

		// UpdateSnapshot updates the properties of an existing snapshot.
		UpdateSnapshot(ctx context.Context, instanceID, snapshotID string, req SnapshotUpdateRequest) (*SnapshotDetailResponse, error)

		// DeleteSnapshot deletes a snapshot.
		DeleteSnapshot(ctx context.Context, instanceID, snapshotID string) error

		// RestoreSnapshot creates a new instance from a snapshot.
		RestoreSnapshot(ctx context.Context, instanceID, snapshotID string, req RestoreSnapshotRequest) (*InstanceResponse, error)
	}

	// instanceService implements the InstanceService interface
	instanceService struct {
		client *DBaaSClient
	}

	// ListInstanceOptions provides options for listing instances
	ListInstanceOptions struct {
		// Offset is the number of instances to skip
		Offset *int
		// Limit is the maximum number of instances to return
		Limit *int
		// Status filters instances by status
		Status *InstanceStatus
		// EngineID filters instances by database engine
		EngineID *string
		// VolumeSize filters instances by exact volume size
		VolumeSize *int
		// VolumeSizeGt filters instances by volume size greater than
		VolumeSizeGt *int
		// VolumeSizeGte filters instances by volume size greater than or equal
		VolumeSizeGte *int
		// VolumeSizeLt filters instances by volume size less than
		VolumeSizeLt *int
		// VolumeSizeLte filters instances by volume size less than or equal
		VolumeSizeLte *int
		// ExpandedFields contains fields to expand in the response
		ExpandedFields []string
	}

	// GetInstanceOptions provides options for getting instance details
	GetInstanceOptions struct {
		// ExpandedFields contains fields to expand in the response
		ExpandedFields []string
	}

	// InstancesResponse represents the response when listing instances
	InstancesResponse struct {
		// Meta contains pagination and filter information
		Meta MetaResponse `json:"meta"`
		// Results is the list of instances
		Results []InstanceDetail `json:"results"`
	}

	// InstanceResponse represents the response when creating an instance
	InstanceResponse struct {
		// ID is the unique identifier of the created instance
		ID string `json:"id"`
	}

	// InstanceDetail represents detailed information about an instance
	InstanceDetail struct {
		// ID is the unique identifier of the instance
		ID string `json:"id"`
		// Name is the name of the instance
		Name string `json:"name"`
		// EngineID is the ID of the database engine
		EngineID string `json:"engine_id"`
		// DatastoreID is the ID of the datastore
		DatastoreID string `json:"datastore_id"`
		// FlavorID is the ID of the flavor
		FlavorID string `json:"flavor_id"`
		// InstanceTypeID is the ID of the instance type
		InstanceTypeID string `json:"instance_type_id"`
		// Volume contains volume information
		Volume Volume `json:"volume"`
		// Addresses contains network addresses
		Addresses []Address `json:"addresses"`
		// Status is the current status of the instance
		Status InstanceStatus `json:"status"`
		// Generation is the generation identifier
		Generation string `json:"generation"`
		// ApplyParametersPending indicates if parameter changes are pending
		ApplyParametersPending bool `json:"apply_parameters_pending"`
		// ParameterGroupID is the ID of the parameter group
		ParameterGroupID string `json:"parameter_group_id"`
		// AvailabilityZone is the availability zone
		AvailabilityZone string `json:"availability_zone"`
		// BackupRetentionDays is the number of days to retain backups
		BackupRetentionDays int `json:"backup_retention_days"`
		// BackupStartAt is the backup start time
		BackupStartAt string `json:"backup_start_at"`
		// CreatedAt is the timestamp when the instance was created
		CreatedAt string `json:"created_at"`
		// UpdatedAt is the timestamp when the instance was last updated
		UpdatedAt *string `json:"updated_at,omitempty"`
		// StartedAt is the timestamp when the instance was started
		StartedAt *string `json:"started_at,omitempty"`
		// FinishedAt is the timestamp when the instance operation finished
		FinishedAt *string `json:"finished_at,omitempty"`
		// MaintenanceScheduledAt is the scheduled maintenance timestamp
		MaintenanceScheduledAt *string `json:"maintenance_scheduled_at,omitempty"`
		// Replicas contains information about associated replicas
		Replicas []ReplicaDetailResponse `json:"replicas,omitempty"`
	}
)

// List implements the List method of InstanceService.
// Returns a paginated list of database instances with optional filters.
func (s *instanceService) List(ctx context.Context, opts ListInstanceOptions) ([]InstanceDetail, error) {
	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Status != nil {
		query.Set("status", string(*opts.Status))
	}
	if opts.EngineID != nil {
		query.Set("engine_id", *opts.EngineID)
	}
	if opts.VolumeSize != nil {
		query.Set("volume.size", strconv.Itoa(*opts.VolumeSize))
	}
	if opts.VolumeSizeGt != nil {
		query.Set("volume.size__gt", strconv.Itoa(*opts.VolumeSizeGt))
	}
	if opts.VolumeSizeGte != nil {
		query.Set("volume.size__gte", strconv.Itoa(*opts.VolumeSizeGte))
	}
	if opts.VolumeSizeLt != nil {
		query.Set("volume.size__lt", strconv.Itoa(*opts.VolumeSizeLt))
	}
	if opts.VolumeSizeLte != nil {
		query.Set("volume.size__lte", strconv.Itoa(*opts.VolumeSizeLte))
	}
	if len(opts.ExpandedFields) > 0 {
		query.Set("_expand", strings.Join(opts.ExpandedFields, ","))
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[InstancesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		InstancePath,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

// Get retrieves details of a specific database instance.
// The instance_id parameter specifies which instance to retrieve.
func (s *instanceService) Get(ctx context.Context, id string, opts GetInstanceOptions) (*InstanceDetail, error) {
	query := make(url.Values)
	if len(opts.ExpandedFields) > 0 {
		query.Set("_expand", strings.Join(opts.ExpandedFields, ","))
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(InstancePathID, id),
		nil,
		query,
	)
}

// Create initiates the asynchronous creation of a new database instance.
// Returns a response containing the ID of the created instance.
func (s *instanceService) Create(ctx context.Context, req InstanceCreateRequest) (*InstanceResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		InstancePath,
		req,
		nil,
	)
}

// Delete initiates the asynchronous deletion of a database instance.
// The operation is considered successful when the deletion process begins.
func (s *instanceService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(InstancePathID, id),
		nil,
		nil,
	)
}

// Update modifies the properties of an existing database instance.
// Returns the updated instance details.
func (s *instanceService) Update(ctx context.Context, id string, req DatabaseInstanceUpdateRequest) (*InstanceDetail, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(InstancePathID, id),
		req,
		nil,
	)
}

// Resize changes the instance type and/or volume specifications of a database instance.
// Returns the instance details with the new specifications.
func (s *instanceService) Resize(ctx context.Context, id string, req InstanceResizeRequest) (*InstanceDetail, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(InstancePathID+"/resize", id),
		req,
		nil,
	)
}

// Start initiates the startup process of a stopped database instance.
// Returns the instance details with updated status.
func (s *instanceService) Start(ctx context.Context, id string) (*InstanceDetail, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(InstancePathID+"/start", id),
		nil,
		nil,
	)
}

// Stop initiates the shutdown process of a running database instance.
// Returns the instance details with updated status.
func (s *instanceService) Stop(ctx context.Context, id string) (*InstanceDetail, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(InstancePathID+"/stop", id),
		nil,
		nil,
	)
}

// ListSnapshots returns a list of snapshots for a specific instance.
func (s *instanceService) ListSnapshots(ctx context.Context, instanceID string, opts ListSnapshotOptions) ([]SnapshotDetailResponse, error) {
	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Type != nil {
		query.Set("type", string(*opts.Type))
	}
	if opts.Status != nil {
		query.Set("status", string(*opts.Status))
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[SnapshotsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(SnapshotPath, instanceID),
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

// CreateSnapshot creates a new snapshot for the specified instance.
func (s *instanceService) CreateSnapshot(ctx context.Context, instanceID string, req SnapshotCreateRequest) (*SnapshotResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SnapshotResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(SnapshotPath, instanceID),
		req,
		nil,
	)
}

// GetSnapshot retrieves details of a specific snapshot.
func (s *instanceService) GetSnapshot(ctx context.Context, instanceID, snapshotID string) (*SnapshotDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SnapshotDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(SnapshotPathID, instanceID, snapshotID),
		nil,
		nil,
	)
}

// UpdateSnapshot updates the properties of an existing snapshot.
func (s *instanceService) UpdateSnapshot(ctx context.Context, instanceID, snapshotID string, req SnapshotUpdateRequest) (*SnapshotDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SnapshotDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(SnapshotPathID, instanceID, snapshotID),
		req,
		nil,
	)
}

// DeleteSnapshot deletes a snapshot.
func (s *instanceService) DeleteSnapshot(ctx context.Context, instanceID, snapshotID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(SnapshotPathID, instanceID, snapshotID),
		nil,
		nil,
	)
}

// RestoreSnapshot creates a new instance from a snapshot.
func (s *instanceService) RestoreSnapshot(ctx context.Context, instanceID, snapshotID string, req RestoreSnapshotRequest) (*InstanceResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[InstanceResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(SnapshotPathID+"/restore", instanceID, snapshotID),
		req,
		nil,
	)
}
