// Package blockstorage provides functionality to interact with the MagaluCloud block storage service.
// This package allows managing volumes, volume types, and snapshots.
package blockstorage

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

// SnapshotVolumeExpand is a constant used for expanding volume information in snapshot responses.
const (
	SnapshotVolumeExpand = "volume"
)

// ListSnapshotsResponse represents the response from listing snapshots.
// This structure encapsulates the API response format for snapshots.
type ListSnapshotsResponse struct {
	// Snapshots contains the list of snapshots
	Snapshots []Snapshot `json:"snapshots"`
}

// Snapshot represents a volume snapshot.
// A snapshot is a point-in-time copy of a volume that can be used for backup or to create new volumes.
type Snapshot struct {
	// ID is the unique identifier of the snapshot
	ID string `json:"id"`
	// Name is the display name of the snapshot
	Name string `json:"name"`
	// Size is the size of the snapshot in gigabytes
	Size int `json:"size"`
	// Description provides additional information about the snapshot
	Description *string `json:"description"`
	// State indicates the current state of the snapshot
	State SnapshotStateV1 `json:"state"`
	// Status indicates the current status of the snapshot
	Status SnapshotStatusV1 `json:"status"`
	// CreatedAt is the timestamp when the snapshot was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the snapshot was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// Volume contains information about the source volume (optional)
	Volume *IDOrName `json:"volume,omitempty"`
	// Error contains error information if the snapshot creation failed
	Error *SnapshotError `json:"error,omitempty"`
	// AvailabilityZones lists the availability zones where this snapshot is available
	AvailabilityZones []string `json:"availability_zones"`
	// Type specifies the type of snapshot
	Type string `json:"type"`
}

// SnapshotError represents error information for a snapshot operation.
type SnapshotError struct {
	// Slug is a short error identifier
	Slug string `json:"slug"`
	// Message provides a detailed error description
	Message string `json:"message"`
}

// CreateSnapshotRequest represents the request to create a new snapshot.
type CreateSnapshotRequest struct {
	// Name is the display name for the new snapshot
	Name string `json:"name"`
	// Volume specifies the volume to snapshot (optional if source snapshot is provided)
	Volume *IDOrName `json:"volume,omitempty"`
	// Description provides additional information about the snapshot
	Description *string `json:"description"`
	// Type specifies the type of snapshot to create
	Type *string `json:"type"`
	// SourceSnapshot specifies an existing snapshot to use as source (optional)
	SourceSnapshot *IDOrName `json:"source_snapshot,omitempty"`
}

// RenameSnapshotRequest represents the request to rename a snapshot.
type RenameSnapshotRequest struct {
	// Name is the new name for the snapshot
	Name string `json:"name"`
}

// SnapshotStateV1 represents the possible states of a snapshot.
// The state indicates the lifecycle stage of the snapshot.
type SnapshotStateV1 string

const (
	// SnapshotStateNew indicates the snapshot is newly created
	SnapshotStateNew SnapshotStateV1 = "new"
	// SnapshotStateAvailable indicates the snapshot is ready for use
	SnapshotStateAvailable SnapshotStateV1 = "available"
	// SnapshotStateDeleted indicates the snapshot has been deleted
	SnapshotStateDeleted SnapshotStateV1 = "deleted"
)

// SnapshotStatusV1 represents the possible statuses of a snapshot.
// The status provides more detailed information about the snapshot's current condition.
type SnapshotStatusV1 string

const (
	// SnapshotStatusProvisioning indicates the snapshot is being provisioned
	SnapshotStatusProvisioning SnapshotStatusV1 = "provisioning"
	// SnapshotStatusCreating indicates the snapshot is being created
	SnapshotStatusCreating SnapshotStatusV1 = "creating"
	// SnapshotStatusCreatingError indicates an error occurred during creation
	SnapshotStatusCreatingError SnapshotStatusV1 = "creating_error"
	// SnapshotStatusCreatingErrorQuota indicates a quota error during creation
	SnapshotStatusCreatingErrorQuota SnapshotStatusV1 = "creating_error_quota"
	// SnapshotStatusCompleted indicates the snapshot creation is complete
	SnapshotStatusCompleted SnapshotStatusV1 = "completed"
	// SnapshotStatusDeleting indicates the snapshot is being deleted
	SnapshotStatusDeleting SnapshotStatusV1 = "deleting"
	// SnapshotStatusDeleted indicates the snapshot has been deleted
	SnapshotStatusDeleted SnapshotStatusV1 = "deleted"
	// SnapshotStatusDeletedError indicates an error occurred during deletion
	SnapshotStatusDeletedError SnapshotStatusV1 = "deleted_error"
	// SnapshotStatusReplicating indicates the snapshot is being replicated
	SnapshotStatusReplicating SnapshotStatusV1 = "replicating"
	// SnapshotStatusReplicatingError indicates an error occurred during replication
	SnapshotStatusReplicatingError SnapshotStatusV1 = "replicating_error"
	// SnapshotStatusRestoring indicates the snapshot is being restored
	SnapshotStatusRestoring SnapshotStatusV1 = "restoring"
	// SnapshotStatusRestoringError indicates an error occurred during restoration
	SnapshotStatusRestoringError SnapshotStatusV1 = "restoring_error"
	// SnapshotStatusReserved indicates the snapshot is reserved
	SnapshotStatusReserved SnapshotStatusV1 = "reserved"
)

// SnapshotService provides operations for managing volume snapshots.
// This interface allows creating, listing, retrieving, and managing snapshots.
type SnapshotService interface {
	// List returns a slice of snapshots based on the provided listing options.
	//
	// Parameters:
	//   - ctx: Request context
	//   - opts: Options to control pagination and expansion
	//
	// Returns:
	//   - []Snapshot: List of snapshots
	//   - error: Error if there's a failure in the request
	List(ctx context.Context, opts ListOptions) ([]Snapshot, error)

	// Create creates a new snapshot from a volume.
	//
	// Parameters:
	//   - ctx: Request context
	//   - req: Request containing snapshot creation parameters
	//
	// Returns:
	//   - string: ID of the created snapshot
	//   - error: Error if there's a failure in the request
	Create(ctx context.Context, req CreateSnapshotRequest) (string, error)

	// Get retrieves a specific snapshot.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the snapshot to retrieve
	//   - expand: List of related resources to expand in the response
	//
	// Returns:
	//   - *Snapshot: The requested snapshot
	//   - error: Error if there's a failure in the request
	Get(ctx context.Context, id string, expand []string) (*Snapshot, error)

	// Delete removes a snapshot.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the snapshot to delete
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Delete(ctx context.Context, id string) error

	// Rename updates the name of a snapshot.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the snapshot to rename
	//   - newName: New name for the snapshot
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Rename(ctx context.Context, id string, newName string) error
}

// snapshotService implements the SnapshotService interface.
// This is an internal implementation that should not be used directly.
type snapshotService struct {
	client *BlockStorageClient
}

// List returns all snapshots.
// This method makes an HTTP request to get the list of snapshots
// and applies the filters specified in the options.
//
// Parameters:
//   - ctx: Request context
//   - opts: Options to control pagination and expansion
//
// Returns:
//   - []Snapshot: List of snapshots
//   - error: Error if there's a failure in the request
func (s *snapshotService) List(ctx context.Context, opts ListOptions) ([]Snapshot, error) {
	q := url.Values{}
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

	path := "/v1/snapshots"

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListSnapshotsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		q,
	)
	if err != nil {
		return nil, err
	}

	return result.Snapshots, nil
}

// Create provisions a new snapshot.
// This method makes an HTTP request to create a new snapshot
// and returns the ID of the created snapshot.
//
// Parameters:
//   - ctx: Request context
//   - req: Request containing snapshot creation parameters
//
// Returns:
//   - string: ID of the created snapshot
//   - error: Error if there's a failure in the request
func (s *snapshotService) Create(ctx context.Context, req CreateSnapshotRequest) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[struct{ ID string }](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/snapshots",
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Get retrieves a specific snapshot.
// This method makes an HTTP request to get detailed information about a snapshot
// and optionally expands related resources.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the snapshot to retrieve
//   - expand: List of related resources to expand in the response
//
// Returns:
//   - *Snapshot: The requested snapshot
//   - error: Error if there's a failure in the request
func (s *snapshotService) Get(ctx context.Context, id string, expand []string) (*Snapshot, error) {
	path := fmt.Sprintf("/v1/snapshots/%s", id)
	if len(expand) > 0 {
		path += "?expand=" + strings.Join(expand, ",")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Snapshot](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		nil,
	)
}

// Delete removes a snapshot.
// This method makes an HTTP request to delete a snapshot permanently.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the snapshot to delete
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *snapshotService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v1/snapshots/%s", id),
		nil,
		nil,
	)
}

// Rename updates the name of a snapshot.
// This method makes an HTTP request to rename an existing snapshot.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the snapshot to rename
//   - newName: New name for the snapshot
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *snapshotService) Rename(ctx context.Context, id string, newName string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v1/snapshots/%s/rename", id),
		RenameSnapshotRequest{Name: newName},
		nil,
	)
}
