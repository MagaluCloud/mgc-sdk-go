package compute

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// SnapshotExpand represents the expand options for snapshot responses.
type SnapshotExpand string

// Constants for expanding related resources in snapshot responses.
const (
	// SnapshotImageExpand is used to include image information in snapshot responses
	SnapshotImageExpand SnapshotExpand = "image"
	// SnapshotMachineTypeExpand is used to include machine type information in snapshot responses
	SnapshotMachineTypeExpand SnapshotExpand = "machine-type"
)

// ListSnapshotsResponse represents the response from listing snapshots.
// This structure encapsulates the API response format for snapshots with pagination metadata.
type ListSnapshotsResponse struct {
	Snapshots []Snapshot `json:"snapshots"`
	Meta      Meta       `json:"meta"`
}

// Snapshot represents an instance snapshot.
// A snapshot is a point-in-time copy of an instance that can be used for backup or to create new instances.
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

// SnapshotInstance represents information about the instance that was snapshotted.
type SnapshotInstance struct {
	ID          string    `json:"id"`
	Image       *IDOrName `json:"image,omitempty"`
	MachineType *IDOrName `json:"machine_type,omitempty"`
}

// CreateSnapshotRequest represents the request to create a new snapshot.
type CreateSnapshotRequest struct {
	Name     string   `json:"name"`
	Instance IDOrName `json:"instance"`
}

// RestoreSnapshotRequest represents the request to restore an instance from a snapshot.
type RestoreSnapshotRequest struct {
	Name             string                   `json:"name"`
	MachineType      IDOrName                 `json:"machine_type"`
	SSHKeyName       *string                  `json:"ssh_key_name,omitempty"`
	AvailabilityZone *string                  `json:"availability_zone,omitempty"`
	Network          *CreateParametersNetwork `json:"network,omitempty"`
	UserData         *string                  `json:"user_data,omitempty"`
}

// CopySnapshotRequest represents the request to copy a snapshot to another region.
type CopySnapshotRequest struct {
	// DestinationRegion is the region where the snapshot should be copied
	DestinationRegion string `json:"destination_region"`
}

// SnapshotListOptions defines the parameters for filtering and pagination of snapshot lists.
type SnapshotListOptions struct {
	Limit  *int
	Offset *int
	Sort   *string
	Expand []SnapshotExpand
}

// SnapshotFilterOptions defines filtering options for ListAll (without pagination).
type SnapshotFilterOptions struct {
	Sort   *string
	Expand []SnapshotExpand
}

// SnapshotService provides operations for managing snapshots.
// This interface allows creating, listing, retrieving, and managing instance snapshots.
type SnapshotService interface {
	List(ctx context.Context, opts SnapshotListOptions) (*ListSnapshotsResponse, error)
	ListAll(ctx context.Context, opts SnapshotFilterOptions) ([]Snapshot, error)
	Create(ctx context.Context, req CreateSnapshotRequest) (string, error)
	Get(ctx context.Context, id string, expand []SnapshotExpand) (*Snapshot, error)
	Delete(ctx context.Context, id string) error
	Rename(ctx context.Context, id string, newName string) error
	Restore(ctx context.Context, id string, req RestoreSnapshotRequest) (string, error)
	Copy(ctx context.Context, id string, req CopySnapshotRequest) error
}

// snapshotService implements the SnapshotService interface.
// This is an internal implementation that should not be used directly.
type snapshotService struct {
	client *VirtualMachineClient
}

// List retrieves snapshots with pagination metadata.
// This method makes an HTTP request to get the list of snapshots
// and applies the filters specified in the options.
func (s *snapshotService) List(ctx context.Context, opts SnapshotListOptions) (*ListSnapshotsResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/snapshots", nil)
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
		expandStrs := make([]string, len(opts.Expand))
		for i, e := range opts.Expand {
			expandStrs[i] = string(e)
		}
		q.Add("expand", strings.Join(expandStrs, ","))
	}
	req.URL.RawQuery = q.Encode()

	response := &ListSnapshotsResponse{}
	_, err = mgc_http.Do(s.client.GetConfig(), ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListAll retrieves all snapshots across all pages with optional filtering.
// This method automatically handles pagination and returns all results.
func (s *snapshotService) ListAll(ctx context.Context, opts SnapshotFilterOptions) ([]Snapshot, error) {
	var allSnapshots []Snapshot
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		listOpts := SnapshotListOptions{
			Offset: &currentOffset,
			Limit:  &currentLimit,
			Sort:   opts.Sort,
			Expand: opts.Expand,
		}

		response, err := s.List(ctx, listOpts)
		if err != nil {
			return nil, err
		}

		allSnapshots = append(allSnapshots, response.Snapshots...)

		// Check if we've retrieved all results
		if len(response.Snapshots) < limit {
			break
		}

		offset += limit
	}

	return allSnapshots, nil
}

// Create creates a new snapshot from an instance.
// This method makes an HTTP request to create a new snapshot
// and returns the ID of the created snapshot.
func (s *snapshotService) Create(ctx context.Context, createReq CreateSnapshotRequest) (string, error) {
	var result struct {
		ID string `json:"id"`
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "/v1/snapshots", createReq)
	if err != nil {
		return "", err
	}

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &result)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// Get retrieves a specific snapshot.
// This method makes an HTTP request to get detailed information about a snapshot
// and optionally expands related resources.
func (s *snapshotService) Get(ctx context.Context, id string, expand []SnapshotExpand) (*Snapshot, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/snapshots/%s", id), nil)
	if err != nil {
		return nil, err
	}

	if len(expand) > 0 {
		q := req.URL.Query()
		expandStrs := make([]string, len(expand))
		for i, e := range expand {
			expandStrs[i] = string(e)
		}
		q.Add("expand", strings.Join(expandStrs, ","))
		req.URL.RawQuery = q.Encode()
	}

	var snapshot Snapshot
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &snapshot)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete removes a snapshot.
// This method makes an HTTP request to delete a snapshot permanently.
func (s *snapshotService) Delete(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/snapshots/%s", id), nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}

// Rename changes the name of a snapshot.
// This method makes an HTTP request to rename an existing snapshot.
func (s *snapshotService) Rename(ctx context.Context, id string, newName string) error {
	req, err := s.client.newRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/v1/snapshots/%s/rename", id),
		UpdateNameRequest{Name: newName})
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}

// Restore creates a new instance from a snapshot.
// This method makes an HTTP request to restore an instance from a snapshot
// and returns the ID of the created instance.
func (s *snapshotService) Restore(ctx context.Context, id string, restoreReq RestoreSnapshotRequest) (string, error) {
	var result struct {
		ID string `json:"id"`
	}

	req, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/snapshots/%s", id),
		restoreReq)
	if err != nil {
		return "", err
	}

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &result)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// Copy copies a snapshot to another region.
// This method makes an HTTP request to copy a snapshot to a different region.
func (s *snapshotService) Copy(ctx context.Context, id string, copyReq CopySnapshotRequest) error {
	req, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/snapshots/%s/copy", id),
		copyReq)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}
