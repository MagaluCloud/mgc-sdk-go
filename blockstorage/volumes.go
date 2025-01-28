package blockstorage

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	VolumeTypeExpand   = "volume_type"
	VolumeAttachExpand = "attachment"
)

type (
	ListVolumesResponse struct {
		Volumes []Volume `json:"volumes"`
	}

	Volume struct {
		ID                string            `json:"id"`
		Name              string            `json:"name"`
		Size              int               `json:"size"`
		Status            string            `json:"status"`
		State             string            `json:"state"`
		CreatedAt         time.Time         `json:"created_at"`
		UpdatedAt         time.Time         `json:"updated_at"`
		Type              IDOrName          `json:"type"`
		Error             *VolumeError      `json:"error,omitempty"`
		Attachment        *VolumeAttachment `json:"attachment,omitempty"`
		AvailabilityZone  string            `json:"availability_zone"`
		AvailabilityZones []string          `json:"availability_zones"`
		Encrypted         bool              `json:"encrypted"`
	}

	VolumeError struct {
		Slug    string `json:"slug"`
		Message string `json:"message"`
	}

	VolumeAttachment struct {
		Instance   AttachmentInstance `json:"instance"`
		AttachedAt time.Time          `json:"attached_at"`
		Device     *string            `json:"device,omitempty"`
	}

	AttachmentInstance struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Status    string    `json:"status"`
		State     string    `json:"state"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	CreateVolumeRequest struct {
		AvailabilityZone *string   `json:"availability_zone,omitempty"`
		Name             string    `json:"name"`
		Size             int       `json:"size"`
		Type             IDOrName  `json:"type"`
		Snapshot         *IDOrName `json:"snapshot,omitempty"`
		Encrypted        bool      `json:"encrypted"`
	}

	ExtendVolumeRequest struct {
		Size int `json:"size"`
	}

	RetypeVolumeRequest struct {
		NewType IDOrName `json:"new_type"`
	}

	RenameVolumeRequest struct {
		Name string `json:"name"`
	}
)

type IDOrName struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type ListOptions struct {
	Limit  *int
	Offset *int
	Sort   *string
	Expand []string
}

// VolumeStateV1 represents volume states
type VolumeStateV1 string

const (
	VolumeStateNew       VolumeStateV1 = "new"
	VolumeStateAvailable VolumeStateV1 = "available"
	VolumeStateInUse     VolumeStateV1 = "in-use"
	VolumeStateDeleted   VolumeStateV1 = "deleted"
	VolumeStateLegacy    VolumeStateV1 = "legacy"
)

// VolumeStatusV1 represents volume statuses
type VolumeStatusV1 string

const (
	VolumeStatusProvisioning       VolumeStatusV1 = "provisioning"
	VolumeStatusCreating           VolumeStatusV1 = "creating"
	VolumeStatusCreatingError      VolumeStatusV1 = "creating_error"
	VolumeStatusCreatingErrorQuota VolumeStatusV1 = "creating_error_quota"
	VolumeStatusCompleted          VolumeStatusV1 = "completed"
	VolumeStatusExtendPending      VolumeStatusV1 = "extend_pending"
	VolumeStatusExtending          VolumeStatusV1 = "extending"
	VolumeStatusExtendError        VolumeStatusV1 = "extend_error"
	VolumeStatusExtendErrorQuota   VolumeStatusV1 = "extend_error_quota"
	VolumeStatusAttachingPending   VolumeStatusV1 = "attaching_pending"
	VolumeStatusAttachingError     VolumeStatusV1 = "attaching_error"
	VolumeStatusAttaching          VolumeStatusV1 = "attaching"
	VolumeStatusDetachingPending   VolumeStatusV1 = "detaching_pending"
	VolumeStatusDetachingError     VolumeStatusV1 = "detaching_error"
	VolumeStatusDetaching          VolumeStatusV1 = "detaching"
	VolumeStatusRetypePending      VolumeStatusV1 = "retype_pending"
	VolumeStatusRetyping           VolumeStatusV1 = "retyping"
	VolumeStatusRetypeError        VolumeStatusV1 = "retype_error"
	VolumeStatusRetypeErrorQuota   VolumeStatusV1 = "retype_error_quota"
	VolumeStatusDeletingPending    VolumeStatusV1 = "deleting_pending"
	VolumeStatusDeleting           VolumeStatusV1 = "deleting"
	VolumeStatusDeleted            VolumeStatusV1 = "deleted"
	VolumeStatusDeletedError       VolumeStatusV1 = "deleted_error"
)

// VolumeService provides operations for managing block storage volumes
type VolumeService interface {
	// List returns a slice of volumes based on the provided listing options.
	// Use ListOptions to control pagination, sorting, and expansion of related resources.
	List(ctx context.Context, opts ListOptions) ([]Volume, error)

	// Create provisions a new volume with the specified configuration.
	// Returns the ID of the newly created volume.
	Create(ctx context.Context, req CreateVolumeRequest) (string, error)

	// Get retrieves detailed information about a specific volume.
	// The expand parameter allows fetching related resources in the same request.
	Get(ctx context.Context, id string, expand []string) (*Volume, error)

	// Delete removes a volume.
	// The volume must be detached from any instances before it can be deleted.
	Delete(ctx context.Context, id string) error

	// Rename updates the display name of an existing volume.
	// Returns an error if the operation fails or if the volume ID is invalid.
	Rename(ctx context.Context, id string, newName string) error

	// Extend increases the size of an existing volume.
	// The volume must be detached or the attached instance must be stopped.
	Extend(ctx context.Context, id string, req ExtendVolumeRequest) error

	// Retype changes the volume type.
	// The volume must be detached or the attached instance must be stopped.
	Retype(ctx context.Context, id string, req RetypeVolumeRequest) error

	// Attach connects a volume to an instance.
	// Returns an error if the volume is already attached or if either ID is invalid.
	Attach(ctx context.Context, volumeID string, instanceID string) error

	// Detach disconnects a volume from an instance.
	// Returns an error if the volume is not attached or if the operation fails.
	Detach(ctx context.Context, volumeID string) error
}

type volumeService struct {
	client *BlockStorageClient
}

// List retrieves all volumes
func (s *volumeService) List(ctx context.Context, opts ListOptions) ([]Volume, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/volumes", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if opts.Limit != nil {
		q.Add("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Add("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Add("_sort", *opts.Sort)
	}
	if len(opts.Expand) > 0 {
		q.Add("expand", strings.Join(opts.Expand, ","))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &ListVolumesResponse{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}

	return resp.Volumes, nil
}

// Create provisions a new volume
func (s *volumeService) Create(ctx context.Context, req CreateVolumeRequest) (string, error) {
	var result struct {
		ID string `json:"id"`
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, "/v1/volumes", req)
	if err != nil {
		return "", err
	}

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &result)
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", fmt.Errorf("empty response")
	}

	return result.ID, nil
}

// Get retrieves a specific volume
func (s *volumeService) Get(ctx context.Context, id string, expand []string) (*Volume, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/volumes/%s", id), nil)
	if err != nil {
		return nil, err
	}

	if len(expand) > 0 {
		q := req.URL.Query()
		q.Add("expand", strings.Join(expand, ","))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &Volume{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete removes a volume
func (s *volumeService) Delete(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/volumes/%s", id), nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// Rename changes the volume name
func (s *volumeService) Rename(ctx context.Context, id string, newName string) error {
	req, err := s.client.newRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/v1/volumes/%s/rename", id),
		RenameVolumeRequest{Name: newName})
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// Extend increases the volume size
func (s *volumeService) Extend(ctx context.Context, id string, req ExtendVolumeRequest) error {
	httpReq, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/extend", id), req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	if err != nil {
		return err
	}

	return nil
}

// Retype changes the volume type
func (s *volumeService) Retype(ctx context.Context, id string, req RetypeVolumeRequest) error {
	httpReq, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/retype", id), req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	if err != nil {
		return err
	}

	return nil
}

// Attach connects a volume to an instance
func (s *volumeService) Attach(ctx context.Context, volumeID string, instanceID string) error {
	req, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/attach/%s", volumeID, instanceID), nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// Detach disconnects a volume from an instance
func (s *volumeService) Detach(ctx context.Context, volumeID string) error {
	req, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/volumes/%s/detach", volumeID), nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
