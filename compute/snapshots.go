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

const (
	SnapshotImageExpand       = "image"
	SnapshotMachineTypeExpand = "machine-type"
)

type (
	ListSnapshotsResponse struct {
		Snapshots []Snapshot `json:"snapshots"`
	}

	Snapshot struct {
		ID        string            `json:"id"`
		Name      string            `json:"name,omitempty"`
		Status    string            `json:"status"`
		State     string            `json:"state"`
		CreatedAt time.Time         `json:"created_at"`
		UpdatedAt *time.Time        `json:"updated_at,omitempty"`
		Size      int               `json:"size"`
		Instance  *SnapshotInstance `json:"instance"`
	}

	SnapshotInstance struct {
		ID          string    `json:"id"`
		Image       *IDOrName `json:"image,omitempty"`
		MachineType *IDOrName `json:"machine_type,omitempty"`
	}

	CreateSnapshotRequest struct {
		Name     string   `json:"name"`
		Instance IDOrName `json:"instance"`
	}

	RestoreSnapshotRequest struct {
		Name             string                   `json:"name"`
		MachineType      IDOrName                 `json:"machine_type"`
		SSHKeyName       *string                  `json:"ssh_key_name,omitempty"`
		AvailabilityZone *string                  `json:"availability_zone,omitempty"`
		Network          *CreateParametersNetwork `json:"network,omitempty"`
		UserData         *string                  `json:"user_data,omitempty"`
	}

	CopySnapshotRequest struct {
		DestinationRegion string `json:"destination_region"`
	}
)

// SnapshotService provides operations for managing snapshots
type SnapshotService interface {
	// List returns a slice of snapshots based on the provided listing options
	List(ctx context.Context, opts ListOptions) ([]Snapshot, error)

	// Create creates a new snapshot from an instance
	Create(ctx context.Context, req CreateSnapshotRequest) (string, error)

	// Get retrieves a specific snapshot
	Get(ctx context.Context, id string, expand []string) (*Snapshot, error)

	// Delete removes a snapshot
	Delete(ctx context.Context, id string) error

	// Rename changes the name of a snapshot
	Rename(ctx context.Context, id string, newName string) error

	// Restore creates a new instance from a snapshot
	Restore(ctx context.Context, id string, req RestoreSnapshotRequest) (string, error)

	// Copy copies a snapshot to another region
	Copy(ctx context.Context, id string, req CopySnapshotRequest) error
}

type snapshotService struct {
	client *VirtualMachineClient
}

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
	if _, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response); err != nil {
		return nil, err
	}

	return response.Snapshots, nil
}

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

func (s *snapshotService) Delete(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/snapshots/%s", id), nil)
	if err != nil {
		return err
	}

	resp, err := mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		return fmt.Errorf("unexpected response data for delete operation")
	}
	return nil
}

func (s *snapshotService) Rename(ctx context.Context, id string, newName string) error {
	req, err := s.client.newRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/v1/snapshots/%s/rename", id),
		UpdateNameRequest{Name: newName})
	if err != nil {
		return err
	}

	resp, err := mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		return fmt.Errorf("unexpected response data for rename operation")
	}
	return nil
}

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

func (s *snapshotService) Copy(ctx context.Context, id string, copyReq CopySnapshotRequest) error {
	req, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/snapshots/%s/copy", id),
		copyReq)
	if err != nil {
		return err
	}

	resp, err := mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		return fmt.Errorf("unexpected response data for copy operation")
	}
	return nil
}
