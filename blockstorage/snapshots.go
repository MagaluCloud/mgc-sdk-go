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

const (
	SnapshotVolumeExpand = "volume"
)

type (
	ListSnapshotsResponse struct {
		Snapshots []Snapshot `json:"snapshots"`
	}

	Snapshot struct {
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

	SnapshotError struct {
		Slug    string `json:"slug"`
		Message string `json:"message"`
	}

	CreateSnapshotRequest struct {
		Name           string    `json:"name"`
		Volume         *IDOrName `json:"volume,omitempty"`
		Description    *string   `json:"description"`
		Type           *string   `json:"type"`
		SourceSnapshot *IDOrName `json:"source_snapshot,omitempty"`
	}

	RenameSnapshotRequest struct {
		Name string `json:"name"`
	}
)

// SnapshotStateV1 represents snapshot states
type SnapshotStateV1 string

const (
	SnapshotStateNew       SnapshotStateV1 = "new"
	SnapshotStateAvailable SnapshotStateV1 = "available"
	SnapshotStateDeleted   SnapshotStateV1 = "deleted"
)

// SnapshotStatusV1 represents snapshot statuses
type SnapshotStatusV1 string

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

// SnapshotService provides operations for managing snapshots
type SnapshotService interface {
	// List returns a slice of snapshots based on the provided listing options
	List(ctx context.Context, opts ListOptions) ([]Snapshot, error)

	// Create creates a new snapshot from a volume
	Create(ctx context.Context, req CreateSnapshotRequest) (string, error)

	// Get retrieves a specific snapshot
	Get(ctx context.Context, id string, expand []string) (*Snapshot, error)

	// Delete removes a snapshot
	Delete(ctx context.Context, id string) error

	// Rename updates the name of a snapshot
	Rename(ctx context.Context, id string, newName string) error
}

type snapshotService struct {
	client *BlockStorageClient
}

// List returns all snapshots
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

// Create provisions a new snapshot
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

// Get retrieves a specific snapshot
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

// Delete removes a snapshot
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

// Rename updates the name of a snapshot
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
