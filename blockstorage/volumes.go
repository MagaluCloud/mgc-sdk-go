// Package blockstorage provides functionality to interact with the MagaluCloud block storage service.
// This package allows managing volumes, volume types, and snapshots.
package blockstorage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// VolumeTypeExpand is a constant used for expanding volume type information in volume responses.
const (
	VolumeTypeExpand   = "volume_type"
	VolumeAttachExpand = "attachment"
)

// ListVolumesResponse represents the response from listing volumes.
// This structure encapsulates the API response format for volumes.
type ListVolumesResponse struct {
	// Volumes contains the list of volumes
	Volumes []Volume `json:"volumes"`
}

// Iops represents the input/output operations per second specifications for a volume.
// IOPS defines the performance characteristics in terms of read/write operations.
type Iops struct {
	// Read specifies the read IOPS limit
	Read int `json:"read"`
	// Write specifies the write IOPS limit
	Write int `json:"write"`
	// Total specifies the total IOPS limit
	Total int `json:"total"`
}

// Type represents the volume type information.
// Contains details about the volume type including IOPS specifications.
type Type struct {
	// IOPS contains the input/output operations per second specifications
	Iops *Iops `json:"iops,omitempty"`
	// ID is the unique identifier of the volume type
	ID string `json:"id"`
	// Name is the display name of the volume type
	Name *string `json:"name,omitempty"`
	// DiskType specifies the physical disk type (e.g., nvme, hdd)
	DiskType *string `json:"disk_type,omitempty"`
	// Status indicates the current status of the volume type
	Status *string `json:"status,omitempty"`
}

// Volume represents a block storage volume.
// A volume is a persistent block storage device that can be attached to instances.
type Volume struct {
	// ID is the unique identifier of the volume
	ID string `json:"id"`
	// Name is the display name of the volume
	Name string `json:"name"`
	// Size is the size of the volume in gigabytes
	Size int `json:"size"`
	// Status indicates the current status of the volume
	Status string `json:"status"`
	// State indicates the current state of the volume
	State string `json:"state"`
	// CreatedAt is the timestamp when the volume was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the volume was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// Type contains information about the volume type
	Type Type `json:"type"`
	// Error contains error information if the volume operation failed
	Error *VolumeError `json:"error,omitempty"`
	// Attachment contains information about the current attachment to an instance
	Attachment *VolumeAttachment `json:"attachment,omitempty"`
	// AvailabilityZone specifies the availability zone where the volume is located
	AvailabilityZone string `json:"availability_zone"`
	// AvailabilityZones lists all availability zones where this volume is available
	AvailabilityZones []string `json:"availability_zones"`
	// Encrypted indicates whether the volume is encrypted
	Encrypted *bool `json:"encrypted,omitempty"`
}

// VolumeError represents error information for a volume operation.
type VolumeError struct {
	// Slug is a short error identifier
	Slug string `json:"slug"`
	// Message provides a detailed error description
	Message string `json:"message"`
}

// VolumeAttachment represents the attachment of a volume to an instance.
type VolumeAttachment struct {
	// Instance contains information about the attached instance
	Instance AttachmentInstance `json:"instance"`
	// AttachedAt is the timestamp when the volume was attached
	AttachedAt time.Time `json:"attached_at"`
	// Device specifies the device name on the instance (optional)
	Device *string `json:"device,omitempty"`
}

// AttachmentInstance represents information about an instance attached to a volume.
type AttachmentInstance struct {
	// ID is the unique identifier of the instance
	ID *string `json:"id"`
	// Name is the display name of the instance
	Name *string `json:"name"`
	// Status indicates the current status of the instance
	Status *string `json:"status"`
	// State indicates the current state of the instance
	State *string `json:"state"`
	// CreatedAt is the timestamp when the instance was created
	CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the instance was last updated
	UpdatedAt *time.Time `json:"updated_at"`
}

// CreateVolumeRequest represents the request to create a new volume.
type CreateVolumeRequest struct {
	// AvailabilityZone specifies the availability zone where the volume should be created
	AvailabilityZone *string `json:"availability_zone,omitempty"`
	// Name is the display name for the new volume
	Name string `json:"name"`
	// Size is the size of the volume in gigabytes
	Size int `json:"size"`
	// Type specifies the volume type to use
	Type IDOrName `json:"type"`
	// Snapshot specifies an existing snapshot to use as source (optional)
	Snapshot *IDOrName `json:"snapshot,omitempty"`
	// Encrypted indicates whether the volume should be encrypted
	Encrypted *bool `json:"encrypted"`
}

// ExtendVolumeRequest represents the request to extend a volume.
type ExtendVolumeRequest struct {
	// Size is the new size of the volume in gigabytes
	Size int `json:"size"`
}

// RetypeVolumeRequest represents the request to change a volume's type.
type RetypeVolumeRequest struct {
	// NewType specifies the new volume type to use
	NewType IDOrName `json:"new_type"`
}

// RenameVolumeRequest represents the request to rename a volume.
type RenameVolumeRequest struct {
	// Name is the new name for the volume
	Name string `json:"name"`
}

// IDOrName represents a reference that can be either an ID or a name.
// This structure is used when an API can accept either an ID or a name as a parameter.
type IDOrName struct {
	// ID is the unique identifier (optional if Name is provided)
	ID *string `json:"id,omitempty"`
	// Name is the display name (optional if ID is provided)
	Name *string `json:"name,omitempty"`
}

// ListOptions contains options for listing volumes.
// All fields are optional and allow controlling pagination and expansion.
type ListOptions struct {
	// Limit defines the maximum number of results to be returned
	Limit *int
	// Offset defines the number of results to be skipped (for pagination)
	Offset *int
	// Sort specifies the sorting criteria
	Sort *string
	// Expand specifies which related resources to include in the response
	Expand []string
}

// VolumeStateV1 represents the possible states of a volume.
// The state indicates the lifecycle stage of the volume.
type VolumeStateV1 string

const (
	// VolumeStateNew indicates the volume is newly created
	VolumeStateNew VolumeStateV1 = "new"
	// VolumeStateAvailable indicates the volume is ready for use
	VolumeStateAvailable VolumeStateV1 = "available"
	// VolumeStateInUse indicates the volume is attached to an instance
	VolumeStateInUse VolumeStateV1 = "in-use"
	// VolumeStateDeleted indicates the volume has been deleted
	VolumeStateDeleted VolumeStateV1 = "deleted"
	// VolumeStateLegacy indicates the volume is in legacy state
	VolumeStateLegacy VolumeStateV1 = "legacy"
)

// VolumeStatusV1 represents the possible statuses of a volume.
// The status provides more detailed information about the volume's current condition.
type VolumeStatusV1 string

const (
	// VolumeStatusProvisioning indicates the volume is being provisioned
	VolumeStatusProvisioning VolumeStatusV1 = "provisioning"
	// VolumeStatusCreating indicates the volume is being created
	VolumeStatusCreating VolumeStatusV1 = "creating"
	// VolumeStatusCreatingError indicates an error occurred during creation
	VolumeStatusCreatingError VolumeStatusV1 = "creating_error"
	// VolumeStatusCreatingErrorQuota indicates a quota error during creation
	VolumeStatusCreatingErrorQuota VolumeStatusV1 = "creating_error_quota"
	// VolumeStatusCompleted indicates the volume creation is complete
	VolumeStatusCompleted VolumeStatusV1 = "completed"
	// VolumeStatusExtendPending indicates the volume extension is pending
	VolumeStatusExtendPending VolumeStatusV1 = "extend_pending"
	// VolumeStatusExtending indicates the volume is being extended
	VolumeStatusExtending VolumeStatusV1 = "extending"
	// VolumeStatusExtendError indicates an error occurred during extension
	VolumeStatusExtendError VolumeStatusV1 = "extend_error"
	// VolumeStatusExtendErrorQuota indicates a quota error during extension
	VolumeStatusExtendErrorQuota VolumeStatusV1 = "extend_error_quota"
	// VolumeStatusAttachingPending indicates the volume attachment is pending
	VolumeStatusAttachingPending VolumeStatusV1 = "attaching_pending"
	// VolumeStatusAttachingError indicates an error occurred during attachment
	VolumeStatusAttachingError VolumeStatusV1 = "attaching_error"
	// VolumeStatusAttaching indicates the volume is being attached
	VolumeStatusAttaching VolumeStatusV1 = "attaching"
	// VolumeStatusDetachingPending indicates the volume detachment is pending
	VolumeStatusDetachingPending VolumeStatusV1 = "detaching_pending"
	// VolumeStatusDetachingError indicates an error occurred during detachment
	VolumeStatusDetachingError VolumeStatusV1 = "detaching_error"
	// VolumeStatusDetaching indicates the volume is being detached
	VolumeStatusDetaching VolumeStatusV1 = "detaching"
	// VolumeStatusRetypePending indicates the volume retype is pending
	VolumeStatusRetypePending VolumeStatusV1 = "retype_pending"
	// VolumeStatusRetyping indicates the volume is being retyped
	VolumeStatusRetyping VolumeStatusV1 = "retyping"
	// VolumeStatusRetypeError indicates an error occurred during retype
	VolumeStatusRetypeError VolumeStatusV1 = "retype_error"
	// VolumeStatusRetypeErrorQuota indicates a quota error during retype
	VolumeStatusRetypeErrorQuota VolumeStatusV1 = "retype_error_quota"
	// VolumeStatusDeletingPending indicates the volume deletion is pending
	VolumeStatusDeletingPending VolumeStatusV1 = "deleting_pending"
	// VolumeStatusDeleting indicates the volume is being deleted
	VolumeStatusDeleting VolumeStatusV1 = "deleting"
	// VolumeStatusDeleted indicates the volume has been deleted
	VolumeStatusDeleted VolumeStatusV1 = "deleted"
	// VolumeStatusDeletedError indicates an error occurred during deletion
	VolumeStatusDeletedError VolumeStatusV1 = "deleted_error"
)

// VolumeService provides operations for managing block storage volumes.
// This interface allows creating, listing, retrieving, and managing volumes.
type VolumeService interface {
	// List returns a slice of volumes based on the provided listing options.
	// Use ListOptions to control pagination, sorting, and expansion of related resources.
	//
	// Parameters:
	//   - ctx: Request context
	//   - opts: Options to control pagination, sorting, and expansion
	//
	// Returns:
	//   - []Volume: List of volumes
	//   - error: Error if there's a failure in the request
	List(ctx context.Context, opts ListOptions) ([]Volume, error)

	// Create provisions a new volume with the specified configuration.
	// Returns the ID of the newly created volume.
	//
	// Parameters:
	//   - ctx: Request context
	//   - req: Request containing volume creation parameters
	//
	// Returns:
	//   - string: ID of the created volume
	//   - error: Error if there's a failure in the request
	Create(ctx context.Context, req CreateVolumeRequest) (string, error)

	// Get retrieves detailed information about a specific volume.
	// The expand parameter allows fetching related resources in the same request.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the volume to retrieve
	//   - expand: List of related resources to expand in the response
	//
	// Returns:
	//   - *Volume: The requested volume
	//   - error: Error if there's a failure in the request
	Get(ctx context.Context, id string, expand []string) (*Volume, error)

	// Delete removes a volume.
	// The volume must be detached from any instances before it can be deleted.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the volume to delete
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Delete(ctx context.Context, id string) error

	// Rename updates the display name of an existing volume.
	// Returns an error if the operation fails or if the volume ID is invalid.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the volume to rename
	//   - newName: New name for the volume
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Rename(ctx context.Context, id string, newName string) error

	// Extend increases the size of an existing volume.
	// The volume must be detached or the attached instance must be stopped.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the volume to extend
	//   - req: Request containing the new size
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Extend(ctx context.Context, id string, req ExtendVolumeRequest) error

	// Retype changes the volume type.
	// The volume must be detached or the attached instance must be stopped.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the volume to retype
	//   - req: Request containing the new volume type
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Retype(ctx context.Context, id string, req RetypeVolumeRequest) error

	// Attach connects a volume to an instance.
	// Returns an error if the volume is already attached or if either ID is invalid.
	//
	// Parameters:
	//   - ctx: Request context
	//   - volumeID: ID of the volume to attach
	//   - instanceID: ID of the instance to attach the volume to
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Attach(ctx context.Context, volumeID string, instanceID string) error

	// Detach disconnects a volume from an instance.
	// Returns an error if the volume is not attached or if the operation fails.
	//
	// Parameters:
	//   - ctx: Request context
	//   - volumeID: ID of the volume to detach
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Detach(ctx context.Context, volumeID string) error
}

// volumeService implements the VolumeService interface.
// This is an internal implementation that should not be used directly.
type volumeService struct {
	client *BlockStorageClient
}

// List retrieves all volumes.
// This method makes an HTTP request to get the list of volumes
// and applies the filters specified in the options.
//
// Parameters:
//   - ctx: Request context
//   - opts: Options to control pagination, sorting, and expansion
//
// Returns:
//   - []Volume: List of volumes
//   - error: Error if there's a failure in the request
func (s *volumeService) List(ctx context.Context, opts ListOptions) ([]Volume, error) {
	path := "/v1/volumes"
	query := make(url.Values)

	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}
	if len(opts.Expand) > 0 {
		for _, expand := range opts.Expand {
			query.Add("expand", expand)
		}
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListVolumesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Volumes, nil
}

// Create provisions a new volume.
// This method makes an HTTP request to create a new volume
// and returns the ID of the created volume.
//
// Parameters:
//   - ctx: Request context
//   - req: Request containing volume creation parameters
//
// Returns:
//   - string: ID of the created volume
//   - error: Error if there's a failure in the request
func (s *volumeService) Create(ctx context.Context, req CreateVolumeRequest) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[struct{ ID string }](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/volumes",
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Get retrieves a specific volume.
// This method makes an HTTP request to get detailed information about a volume
// and optionally expands related resources.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the volume to retrieve
//   - expand: List of related resources to expand in the response
//
// Returns:
//   - *Volume: The requested volume
//   - error: Error if there's a failure in the request
func (s *volumeService) Get(ctx context.Context, id string, expand []string) (*Volume, error) {
	path := fmt.Sprintf("/v1/volumes/%s", id)
	query := make(url.Values)
	if len(expand) > 0 {
		for _, expand := range expand {
			query.Add("expand", expand)
		}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Volume](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		query,
	)
}

// Delete removes a volume.
// This method makes an HTTP request to delete a volume permanently.
// The volume must be detached from any instances before it can be deleted.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the volume to delete
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *volumeService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v1/volumes/%s", id),
		nil,
		nil,
	)
}

// Rename changes the volume name.
// This method makes an HTTP request to rename an existing volume.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the volume to rename
//   - newName: New name for the volume
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *volumeService) Rename(ctx context.Context, id string, newName string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v1/volumes/%s/rename", id),
		RenameVolumeRequest{Name: newName},
		nil,
	)
}

// Extend increases the volume size.
// This method makes an HTTP request to extend an existing volume.
// The volume must be detached or the attached instance must be stopped.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the volume to extend
//   - req: Request containing the new size
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *volumeService) Extend(ctx context.Context, id string, req ExtendVolumeRequest) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/extend", id),
		req,
		nil,
	)
}

// Retype changes the volume type.
// This method makes an HTTP request to change the type of an existing volume.
// The volume must be detached or the attached instance must be stopped.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the volume to retype
//   - req: Request containing the new volume type
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *volumeService) Retype(ctx context.Context, id string, req RetypeVolumeRequest) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/retype", id),
		req,
		nil,
	)
}

// Attach connects a volume to an instance.
// This method makes an HTTP request to attach a volume to an instance.
// Returns an error if the volume is already attached or if either ID is invalid.
//
// Parameters:
//   - ctx: Request context
//   - volumeID: ID of the volume to attach
//   - instanceID: ID of the instance to attach the volume to
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *volumeService) Attach(ctx context.Context, volumeID string, instanceID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/attach/%s", volumeID, instanceID),
		nil,
		nil,
	)
}

// Detach disconnects a volume from an instance.
// This method makes an HTTP request to detach a volume from an instance.
// Returns an error if the volume is not attached or if the operation fails.
//
// Parameters:
//   - ctx: Request context
//   - volumeID: ID of the volume to detach
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *volumeService) Detach(ctx context.Context, volumeID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/detach", volumeID),
		nil,
		nil,
	)
}
