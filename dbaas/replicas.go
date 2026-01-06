package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// ReplicaService provides methods for managing database replicas
	ReplicaService interface {
		List(ctx context.Context, opts ListReplicaOptions) (*ReplicasResponse, error)
		ListAll(ctx context.Context, opts ReplicaFilterOptions) ([]ReplicaDetailResponse, error)
		Get(ctx context.Context, id string) (*ReplicaDetailResponse, error)
		Create(ctx context.Context, req ReplicaCreateRequest) (*ReplicaResponse, error)
		Delete(ctx context.Context, id string) error
		Resize(ctx context.Context, id string, req ReplicaResizeRequest) (*ReplicaDetailResponse, error)
		Start(ctx context.Context, id string) (*ReplicaDetailResponse, error)
		Stop(ctx context.Context, id string) (*ReplicaDetailResponse, error)
	}

	// replicaService implements the ReplicaService interface
	replicaService struct {
		client *DBaaSClient
	}

	// ListReplicaOptions provides options for listing replicas
	ListReplicaOptions struct {
		Offset   *int
		Limit    *int
		SourceID *string
	}

	// ReplicaFilterOptions provides filtering options for ListAll
	ReplicaFilterOptions struct {
		SourceID *string
	}

	// ReplicaDetailResponse represents detailed information about a replica
	ReplicaDetailResponse struct {
		ID                     string                   `json:"id"`
		SourceID               string                   `json:"source_id"`
		Name                   string                   `json:"name"`
		EngineID               string                   `json:"engine_id"`
		InstanceTypeID         string                   `json:"instance_type_id"`
		Volume                 Volume                   `json:"volume"`
		Addresses              []ReplicaAddressResponse `json:"addresses"`
		Status                 InstanceStatus           `json:"status"`
		Generation             string                   `json:"generation"`
		CreatedAt              time.Time                `json:"created_at"`
		UpdatedAt              *time.Time               `json:"updated_at,omitempty"`
		StartedAt              *string                  `json:"started_at,omitempty"`
		FinishedAt             *string                  `json:"finished_at,omitempty"`
		MaintenanceScheduledAt *string                  `json:"maintenance_scheduled_at,omitempty"`
	}

	// ReplicasResponse represents the response when listing replicas
	ReplicasResponse struct {
		Meta    MetaResponse            `json:"meta"`
		Results []ReplicaDetailResponse `json:"results"`
	}

	// ReplicaCreateRequest represents the request payload for creating a replica
	ReplicaCreateRequest struct {
		SourceID       string  `json:"source_id"`
		Name           string  `json:"name"`
		InstanceTypeID *string `json:"instance_type_id,omitempty"`
	}

	// ReplicaResizeRequest represents the request payload for resizing a replica
	ReplicaResizeRequest struct {
		InstanceTypeID *string                      `json:"instance_type_id,omitempty"`
		Volume         *InstanceVolumeResizeRequest `json:"volume,omitempty"`
	}

	// ReplicaResponse represents the response when creating a replica
	ReplicaResponse struct {
		ID string `json:"id"`
	}
)

// List returns a paginated list of database replicas with optional source_id filter
func (s *replicaService) List(ctx context.Context, opts ListReplicaOptions) (*ReplicasResponse, error) {
	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.SourceID != nil {
		query.Set("source_id", *opts.SourceID)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ReplicasResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v2/replicas",
		nil,
		query,
	)
}

// ListAll retrieves all replicas across all pages with optional filtering
func (s *replicaService) ListAll(ctx context.Context, opts ReplicaFilterOptions) ([]ReplicaDetailResponse, error) {
	var allResults []ReplicaDetailResponse
	offset := 0
	limit := 25

	for {
		currentOffset := offset
		currentLimit := limit
		listOpts := ListReplicaOptions{
			Offset:   &currentOffset,
			Limit:    &currentLimit,
			SourceID: opts.SourceID,
		}

		resp, err := s.List(ctx, listOpts)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allResults, nil
}

// Get retrieves details of a specific replica instance
func (s *replicaService) Get(ctx context.Context, id string) (*ReplicaDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ReplicaDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v2/replicas/%s", id),
		nil,
		nil,
	)
}

// Create initiates the creation of a new replica instance
func (s *replicaService) Create(ctx context.Context, req ReplicaCreateRequest) (*ReplicaResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ReplicaResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v2/replicas",
		req,
		nil,
	)
}

// Delete initiates the deletion of a replica instance
func (s *replicaService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v2/replicas/%s", id),
		nil,
		nil,
	)
}

// Resize changes the instance type of a replica
func (s *replicaService) Resize(ctx context.Context, id string, req ReplicaResizeRequest) (*ReplicaDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ReplicaDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v2/replicas/%s/resize", id),
		req,
		nil,
	)
}

// Start initiates a stopped replica instance
func (s *replicaService) Start(ctx context.Context, id string) (*ReplicaDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ReplicaDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v2/replicas/%s/start", id),
		nil,
		nil,
	)
}

// Stop stops a running replica instance
func (s *replicaService) Stop(ctx context.Context, id string) (*ReplicaDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ReplicaDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v2/replicas/%s/stop", id),
		nil,
		nil,
	)
}
