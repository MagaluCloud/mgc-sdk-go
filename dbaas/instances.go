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
	InstancePath   = "/v2/instances"
	InstancePathID = InstancePath + "/%s"
	SnapshotPath   = InstancePathID + "/snapshots"
	SnapshotPathID = SnapshotPath + "/%s"
)

type (
	InstanceStatus string
	AddressAccess  string
	AddressType    string
	SnapshotType   string
	SnapshotStatus string
)

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

const (
	AddressAccessPrivate AddressAccess = "PRIVATE"
	AddressAccessPublic  AddressAccess = "PUBLIC"
)

const (
	AddressTypeIPv4 AddressType = "IPv4"
	AddressTypeIPv6 AddressType = "IPv6"
)

const (
	SnapshotTypeOnDemand  SnapshotType = "ON_DEMAND"
	SnapshotTypeAutomated SnapshotType = "AUTOMATED"
)

const (
	SnapshotStatusPending   SnapshotStatus = "PENDING"
	SnapshotStatusCreating  SnapshotStatus = "CREATING"
	SnapshotStatusAvailable SnapshotStatus = "AVAILABLE"
	SnapshotStatusRestoring SnapshotStatus = "RESTORING"
	SnapshotStatusError     SnapshotStatus = "ERROR"
	SnapshotStatusDeleting  SnapshotStatus = "DELETING"
	SnapshotStatusDeleted   SnapshotStatus = "DELETED"
)

type (
	Volume struct {
		Size int    `json:"size"`
		Type string `json:"type"`
	}

	Address struct {
		Access  AddressAccess `json:"access"`
		Type    *AddressType  `json:"type,omitempty"`
		Address *string       `json:"address,omitempty"`
	}

	InstanceParametersResponse struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	InstanceParametersRequest struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	InstanceVolumeRequest struct {
		Size int    `json:"size"`
		Type string `json:"type,omitempty"`
	}

	InstanceVolumeResizeRequest struct {
		Size int    `json:"size"`
		Type string `json:"type,omitempty"`
	}

	InstanceCreateRequest struct {
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

	InstanceResizeRequest struct {
		InstanceTypeID *string                      `json:"instance_type_id,omitempty"`
		Volume         *InstanceVolumeResizeRequest `json:"volume,omitempty"`
	}

	DatabaseInstanceUpdateRequest struct {
		BackupRetentionDays *int    `json:"backup_retention_days,omitempty"`
		BackupStartAt       *string `json:"backup_start_at,omitempty"`
	}

	ReplicaAddressResponse struct {
		Access  AddressAccess `json:"access"`
		Type    *AddressType  `json:"type,omitempty"`
		Address *string       `json:"address,omitempty"`
	}

	SnapshotsResponse struct {
		Meta    MetaResponse             `json:"meta"`
		Results []SnapshotDetailResponse `json:"results"`
	}

	SnapshotDetailResponse struct {
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

	SnapshotInstanceDetailResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	SnapshotCreateRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
	}

	SnapshotUpdateRequest struct {
		Name        string  `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	SnapshotResponse struct {
		ID string `json:"id"`
	}

	RestoreSnapshotRequest struct {
		Name                string                 `json:"name"`
		InstanceTypeID      string                 `json:"instance_type_id"`
		Volume              *InstanceVolumeRequest `json:"volume,omitempty"`
		BackupRetentionDays *int                   `json:"backup_retention_days,omitempty"`
		BackupStartAt       *string                `json:"backup_start_at,omitempty"`
	}

	ListSnapshotOptions struct {
		Offset *int            `json:"offset,omitempty"`
		Limit  *int            `json:"limit,omitempty"`
		Type   *SnapshotType   `json:"type,omitempty"`
		Status *SnapshotStatus `json:"status,omitempty"`
	}
)

type InstanceStatusUpdate string

const (
	InstanceStatusUpdateActive  InstanceStatusUpdate = "ACTIVE"
	InstanceStatusUpdateStopped InstanceStatusUpdate = "STOPPED"
)

type (
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

	instanceService struct {
		client *DBaaSClient
	}

	ListInstanceOptions struct {
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

	GetInstanceOptions struct {
		ExpandedFields []string
	}

	InstancesResponse struct {
		Meta    MetaResponse     `json:"meta"`
		Results []InstanceDetail `json:"results"`
	}

	InstanceResponse struct {
		ID string `json:"id"`
	}

	InstanceDetail struct {
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
