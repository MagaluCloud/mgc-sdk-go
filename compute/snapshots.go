// Package compute provides functionality to interact with the MagaluCloud compute service.
// This package allows managing virtual machine instances, images, instance types, and snapshots.
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

// Constants for expanding related resources in snapshot responses.
const (
	// SnapshotImageExpand is used to include image information in snapshot responses
	SnapshotImageExpand = "image"
	// SnapshotMachineTypeExpand is used to include machine type information in snapshot responses
	SnapshotMachineTypeExpand = "machine-type"
)

// ListSnapshotsResponse represents the response from listing snapshots.
// This structure encapsulates the API response format for snapshots.
type ListSnapshotsResponse struct {
	// Snapshots contains the list of snapshots
	Snapshots []Snapshot `json:"snapshots"`
}

// Snapshot represents an instance snapshot.
// A snapshot is a point-in-time copy of an instance that can be used for backup or to create new instances.
type Snapshot struct {
	// ID is the unique identifier of the snapshot
	ID string `json:"id"`
	// Name is the display name of the snapshot
	Name string `json:"name,omitempty"`
	// Status indicates the current status of the snapshot
	Status string `json:"status"`
	// State indicates the current state of the snapshot
	State string `json:"state"`
	// CreatedAt is the timestamp when the snapshot was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the snapshot was last updated
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// Size is the size of the snapshot in GB
	Size int `json:"size"`
	// Instance contains information about the source instance
	Instance *SnapshotInstance `json:"instance"`
}

// SnapshotInstance represents information about the instance that was snapshotted.
type SnapshotInstance struct {
	// ID is the unique identifier of the instance
	ID string `json:"id"`
	// Image contains information about the instance's image
	Image *IDOrName `json:"image,omitempty"`
	// MachineType contains information about the instance's machine type
	MachineType *IDOrName `json:"machine_type,omitempty"`
}

// CreateSnapshotRequest represents the request to create a new snapshot.
type CreateSnapshotRequest struct {
	// Name is the display name for the new snapshot
	Name string `json:"name"`
	// Instance specifies the instance to snapshot
	Instance IDOrName `json:"instance"`
}

// RestoreSnapshotRequest represents the request to restore an instance from a snapshot.
type RestoreSnapshotRequest struct {
	// Name is the display name for the restored instance
	Name string `json:"name"`
	// MachineType specifies the machine type for the restored instance
	MachineType IDOrName `json:"machine_type"`
	// SSHKeyName specifies the SSH key to use (optional)
	SSHKeyName *string `json:"ssh_key_name,omitempty"`
	// AvailabilityZone specifies the availability zone for the restored instance (optional)
	AvailabilityZone *string `json:"availability_zone,omitempty"`
	// Network specifies the network configuration for the restored instance (optional)
	Network *CreateParametersNetwork `json:"network,omitempty"`
	// UserData specifies user data to pass to the restored instance (optional)
	UserData *string `json:"user_data,omitempty"`
}

// CopySnapshotRequest represents the request to copy a snapshot to another region.
type CopySnapshotRequest struct {
	// DestinationRegion is the region where the snapshot should be copied
	DestinationRegion string `json:"destination_region"`
}

// SnapshotService provides operations for managing snapshots.
// This interface allows creating, listing, retrieving, and managing instance snapshots.
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

	// Create creates a new snapshot from an instance.
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

	// Rename changes the name of a snapshot.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the snapshot to rename
	//   - newName: New name for the snapshot
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Rename(ctx context.Context, id string, newName string) error

	// Restore creates a new instance from a snapshot.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the snapshot to restore from
	//   - req: Request containing restoration parameters
	//
	// Returns:
	//   - string: ID of the created instance
	//   - error: Error if there's a failure in the request
	Restore(ctx context.Context, id string, req RestoreSnapshotRequest) (string, error)

	// Copy copies a snapshot to another region.
	//
	// Parameters:
	//   - ctx: Request context
	//   - id: ID of the snapshot to copy
	//   - req: Request containing copy parameters
	//
	// Returns:
	//   - error: Error if there's a failure in the request
	Copy(ctx context.Context, id string, req CopySnapshotRequest) error
}

// snapshotService implements the SnapshotService interface.
// This is an internal implementation that should not be used directly.
type snapshotService struct {
	client *VirtualMachineClient
}

// List returns a slice of snapshots based on the provided listing options.
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
		q.Add("expand", strings.Join(opts.Expand, ","))
	}
	req.URL.RawQuery = q.Encode()

	var response ListSnapshotsResponse
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return resp.Snapshots, nil
}

// Create creates a new snapshot from an instance.
// This method makes an HTTP request to create a new snapshot
// and returns the ID of the created snapshot.
//
// Parameters:
//   - ctx: Request context
//   - createReq: Request containing snapshot creation parameters
//
// Returns:
//   - string: ID of the created snapshot
//   - error: Error if there's a failure in the request
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
	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/snapshots/%s", id), nil)
	if err != nil {
		return nil, err
	}

	if len(expand) > 0 {
		q := req.URL.Query()
		q.Add("expand", strings.Join(expand, ","))
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
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the snapshot to delete
//
// Returns:
//   - error: Error if there's a failure in the request
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
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the snapshot to rename
//   - newName: New name for the snapshot
//
// Returns:
//   - error: Error if there's a failure in the request
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
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the snapshot to restore from
//   - restoreReq: Request containing restoration parameters
//
// Returns:
//   - string: ID of the created instance
//   - error: Error if there's a failure in the request
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
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the snapshot to copy
//   - copyReq: Request containing copy parameters
//
// Returns:
//   - error: Error if there's a failure in the request
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
