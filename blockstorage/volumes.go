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
	VolumeTypeExpand   VolumeExpand = "volume_type"
	VolumeAttachExpand VolumeExpand = "attachment"
)

type VolumeExpand = string

// ListVolumesResponse represents the response from listing volumes.
// This structure encapsulates the API response format for volumes.
type ListVolumesResponse struct {
	Meta    Metadata `json:"meta"`
	Volumes []Volume `json:"volumes"`
}

// Iops represents the input/output operations per second specifications for a volume.
// IOPS defines the performance characteristics in terms of read/write operations.
type Iops struct {
	Read  int `json:"read"`
	Write int `json:"write"`
	Total int `json:"total"`
}

// Type represents the volume type information.
// Contains details about the volume type including IOPS specifications.
type Type struct {
	Iops     *Iops   `json:"iops,omitempty"`
	ID       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	DiskType *string `json:"disk_type,omitempty"`
	Status   *string `json:"status,omitempty"`
}

// Volume represents a block storage volume.
// A volume is a persistent block storage device that can be attached to instances.
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

// VolumeError represents error information for a volume operation.
type VolumeError struct {
	Slug    string `json:"slug"`
	Message string `json:"message"`
}

// VolumeAttachment represents the attachment of a volume to an instance.
type VolumeAttachment struct {
	Instance   AttachmentInstance `json:"instance"`
	AttachedAt time.Time          `json:"attached_at"`
	Device     *string            `json:"device,omitempty"`
}

// AttachmentInstance represents information about an instance attached to a volume.
type AttachmentInstance struct {
	ID        *string    `json:"id"`
	Name      *string    `json:"name"`
	Status    *string    `json:"status"`
	State     *string    `json:"state"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// CreateVolumeRequest represents the request to create a new volume.
type CreateVolumeRequest struct {
	AvailabilityZone *string   `json:"availability_zone,omitempty"`
	Name             string    `json:"name"`
	Size             int       `json:"size"`
	Type             IDOrName  `json:"type"`
	Snapshot         *IDOrName `json:"snapshot,omitempty"`
	Encrypted        *bool     `json:"encrypted"`
}

// ExtendVolumeRequest represents the request to extend a volume.
type ExtendVolumeRequest struct {
	Size int `json:"size"`
}

// RetypeVolumeRequest represents the request to change a volume's type.
type RetypeVolumeRequest struct {
	NewType IDOrName `json:"new_type"`
}

// RenameVolumeRequest represents the request to rename a volume.
type RenameVolumeRequest struct {
	Name string `json:"name"`
}

// IDOrName represents a reference that can be either an ID or a name.
// This structure is used when an API can accept either an ID or a name as a parameter.
type IDOrName struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ListOptions contains options for listing volumes.
// All fields are optional and allow controlling pagination and expansion.
type ListOptions struct {
	Limit  *int
	Offset *int
	Sort   *string
	Expand []VolumeExpand
}

// VolumeStateV1 represents the possible states of a volume.
// The state indicates the lifecycle stage of the volume.
type VolumeStateV1 string

const (
	VolumeStateNew       VolumeStateV1 = "new"
	VolumeStateAvailable VolumeStateV1 = "available"
	VolumeStateInUse     VolumeStateV1 = "in-use"
	VolumeStateDeleted   VolumeStateV1 = "deleted"
	VolumeStateLegacy    VolumeStateV1 = "legacy"
)

// VolumeStatusV1 represents the possible statuses of a volume.
// The status provides more detailed information about the volume's current condition.
type VolumeStatusV1 string

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

// VolumeService defines the interface for volume operations.
// This interface provides methods for managing block storage volumes.
type VolumeService interface {
	List(ctx context.Context, opts ListOptions) (*ListVolumesResponse, error)
	ListAll(ctx context.Context, expand []VolumeExpand) ([]Volume, error)
	Create(ctx context.Context, req CreateVolumeRequest) (string, error)
	Get(ctx context.Context, id string, expand []string) (*Volume, error)
	Delete(ctx context.Context, id string) error
	Rename(ctx context.Context, id string, newName string) error
	Extend(ctx context.Context, id string, req ExtendVolumeRequest) error
	Retype(ctx context.Context, id string, req RetypeVolumeRequest) error
	Attach(ctx context.Context, volumeID string, instanceID string) error
	Detach(ctx context.Context, volumeID string) error
}

// volumeService implements the VolumeService interface.
// This is an internal implementation that should not be used directly.
type volumeService struct {
	client *BlockStorageClient
}

// List retrieves a paginated list of volumes.
// This method makes an HTTP request to get the list of volumes
// and applies the filters specified in the options.
func (s *volumeService) List(ctx context.Context, opts ListOptions) (*ListVolumesResponse, error) {
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
	return result, nil
}

// ListAll retrieves all volumes by fetching all pages.
// This method repeatedly calls List to get all available volumes.
func (s *volumeService) ListAll(ctx context.Context, expand []VolumeExpand) ([]Volume, error) {
	var allVolumes []Volume
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		opts := ListOptions{
			Offset: &currentOffset,
			Limit:  &currentLimit,
			Expand: expand,
		}

		resp, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allVolumes = append(allVolumes, resp.Volumes...)

		if len(resp.Volumes) < limit {
			break
		}

		offset += limit
	}

	return allVolumes, nil
}

// Create provisions a new volume.
// This method makes an HTTP request to create a new volume
// and returns the ID of the created volume.
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
