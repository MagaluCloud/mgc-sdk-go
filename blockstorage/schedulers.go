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

// SchedulerListResponse represents the response from listing schedulers.
// This structure encapsulates the API response format for schedulers.
type SchedulerListResponse struct {
	Meta       Metadata            `json:"meta"`
	Schedulers []SchedulerResponse `json:"schedulers"`
}

// Metadata represents pagination metadata.
type Metadata struct {
	Page PageMetadata `json:"page"`
}

// PageMetadata contains pagination information.
type PageMetadata struct {
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
	Count    int `json:"count"`
	Total    int `json:"total"`
	MaxLimit int `json:"max_limit"`
}

// DailyFrequency represents daily scheduling frequency.
type DailyFrequency struct {
	StartTime string `json:"start_time"`
}

// Frequency represents scheduling frequency configuration.
type Frequency struct {
	Daily DailyFrequency `json:"daily"`
}

// Policy represents the scheduler policy.
type Policy struct {
	RetentionInDays int       `json:"retention_in_days"`
	Frequency       Frequency `json:"frequency"`
}

// SnapshotConfig represents snapshot configuration.
type SnapshotConfig struct {
	Type string `json:"type"`
}

// SchedulerState represents the possible states of a scheduler.
type SchedulerState string

const (
	SchedulerStateAvailable SchedulerState = "available"
	SchedulerStateDeleted   SchedulerState = "deleted"
)

// ExpandSchedulers represents the possible expand options for schedulers.
type ExpandSchedulers string

const (
	ExpandSchedulersVolume ExpandSchedulers = "volume"
)

// SchedulerResponse represents a scheduler.
// A scheduler automates snapshot creation and retention for volumes.
type SchedulerResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Volumes     []string        `json:"volumes,omitempty"`
	Snapshot    *SnapshotConfig `json:"snapshot,omitempty"`
	State       SchedulerState  `json:"state"`
	Policy      Policy          `json:"policy"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// SchedulerPayload represents the request to create a new scheduler.
type SchedulerPayload struct {
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Snapshot    SnapshotConfig `json:"snapshot"`
	Policy      Policy         `json:"policy"`
}

// SchedulerVolumeIdentifierPayload represents the request to attach/detach a volume to/from a scheduler.
type SchedulerVolumeIdentifierPayload struct {
	Volume IDOrName `json:"volume"`
}

// SchedulerListOptions contains options for listing schedulers.
// All fields are optional and allow controlling pagination and expansion.
type SchedulerListOptions struct {
	Limit  *int
	Offset *int
	Sort   *string
	Expand []ExpandSchedulers
}

// SchedulerFilterOptions provides filtering options for ListAll (without pagination)
type SchedulerFilterOptions struct {
	Sort   *string
	Expand []ExpandSchedulers
}

// SchedulerService provides operations for managing volume schedulers.
// This interface allows creating, listing, retrieving, and managing schedulers.
type SchedulerService interface {
	List(ctx context.Context, opts SchedulerListOptions) (*SchedulerListResponse, error)
	ListAll(ctx context.Context, filterOpts SchedulerFilterOptions) ([]SchedulerResponse, error)
	Create(ctx context.Context, req SchedulerPayload) (string, error)
	Get(ctx context.Context, id string, expand []ExpandSchedulers) (*SchedulerResponse, error)
	Delete(ctx context.Context, id string) error
	AttachVolume(ctx context.Context, id string, req SchedulerVolumeIdentifierPayload) error
	DetachVolume(ctx context.Context, id string, req SchedulerVolumeIdentifierPayload) error
}

// schedulerService implements the SchedulerService interface.
// This is an internal implementation that should not be used directly.
type schedulerService struct {
	client *BlockStorageClient
}

// List returns all schedulers.
// This method makes an HTTP request to get the list of schedulers
// and applies the filters specified in the options.
func (s *schedulerService) List(ctx context.Context, opts SchedulerListOptions) (*SchedulerListResponse, error) {
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
		expandStrs := make([]string, len(opts.Expand))
		for i, expand := range opts.Expand {
			expandStrs[i] = string(expand)
		}
		for _, expand := range expandStrs {
			query.Add("expand", expand)
		}
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[SchedulerListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v1/schedulers",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListAll retrieves all schedulers by fetching all pages with optional filtering
func (s *schedulerService) ListAll(ctx context.Context, filterOpts SchedulerFilterOptions) ([]SchedulerResponse, error) {
	var allSchedulers []SchedulerResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		opts := SchedulerListOptions{
			Offset: &currentOffset,
			Limit:  &currentLimit,
			Sort:   filterOpts.Sort,
			Expand: filterOpts.Expand,
		}

		resp, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allSchedulers = append(allSchedulers, resp.Schedulers...)

		if len(resp.Schedulers) < limit {
			break
		}

		offset += limit
	}

	return allSchedulers, nil
}

// Create provisions a new scheduler.
// This method makes an HTTP request to create a new scheduler
// and returns the ID of the created scheduler.
func (s *schedulerService) Create(ctx context.Context, req SchedulerPayload) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[struct{ ID string }](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/schedulers",
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Get retrieves a specific scheduler.
// This method makes an HTTP request to get detailed information about a scheduler
// and optionally expands related resources.
func (s *schedulerService) Get(ctx context.Context, id string, expand []ExpandSchedulers) (*SchedulerResponse, error) {
	path := fmt.Sprintf("/v1/schedulers/%s", id)
	query := make(url.Values)

	if len(expand) > 0 {
		expandStrs := make([]string, len(expand))
		for i, expand := range expand {
			expandStrs[i] = string(expand)
		}
		for _, expand := range expandStrs {
			query.Add("expand", expand)
		}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[SchedulerResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		query,
	)
}

// Delete removes a scheduler.
// This method makes an HTTP request to delete a scheduler permanently.
func (s *schedulerService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v1/schedulers/%s", id),
		nil,
		nil,
	)
}

// AttachVolume attaches a volume to a scheduler.
// This method makes an HTTP request to attach a volume to an existing scheduler.
func (s *schedulerService) AttachVolume(ctx context.Context, id string, req SchedulerVolumeIdentifierPayload) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v1/schedulers/%s/attach", id),
		req,
		nil,
	)
}

// DetachVolume detaches a volume from a scheduler.
// This method makes an HTTP request to detach a volume from an existing scheduler.
func (s *schedulerService) DetachVolume(ctx context.Context, id string, req SchedulerVolumeIdentifierPayload) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v1/schedulers/%s/detach", id),
		req,
		nil,
	)
}
