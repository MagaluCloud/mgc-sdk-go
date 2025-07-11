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
		// List returns a list of database replicas.
		// Supports filtering by source_id and pagination.
		List(ctx context.Context, opts ListReplicaOptions) ([]ReplicaDetailResponse, error)

		// Get returns details of a specific replica by its ID.
		Get(ctx context.Context, id string) (*ReplicaDetailResponse, error)

		// Create creates a new replica for an instance asynchronously.
		Create(ctx context.Context, req ReplicaCreateRequest) (*ReplicaResponse, error)

		// Delete deletes a replica instance asynchronously.
		Delete(ctx context.Context, id string) error

		// Resize changes the instance type of a replica.
		Resize(ctx context.Context, id string, req ReplicaResizeRequest) (*ReplicaDetailResponse, error)

		// Start initiates a stopped replica instance.
		Start(ctx context.Context, id string) (*ReplicaDetailResponse, error)

		// Stop stops a running replica instance.
		Stop(ctx context.Context, id string) (*ReplicaDetailResponse, error)
	}

	// replicaService implements the ReplicaService interface
	replicaService struct {
		client *DBaaSClient
	}

	// ListReplicaOptions provides options for listing replicas
	ListReplicaOptions struct {
		// Offset is the number of replicas to skip
		Offset *int
		// Limit is the maximum number of replicas to return
		Limit *int
		// SourceID filters replicas by source instance ID
		SourceID *string
	}

	// ReplicaDetailResponse represents detailed information about a replica
	ReplicaDetailResponse struct {
		// ID is the unique identifier of the replica
		ID string `json:"id"`
		// SourceID is the ID of the source instance
		SourceID string `json:"source_id"`
		// Name is the name of the replica
		Name string `json:"name"`
		// EngineID is the ID of the database engine
		EngineID string `json:"engine_id"`
		// InstanceTypeID is the ID of the instance type
		InstanceTypeID string `json:"instance_type_id"`
		// Volume contains volume information
		Volume Volume `json:"volume"`
		// Addresses contains network addresses
		Addresses []ReplicaAddressResponse `json:"addresses"`
		// Status is the current status of the replica
		Status InstanceStatus `json:"status"`
		// Generation is the generation identifier
		Generation string `json:"generation"`
		// CreatedAt is the timestamp when the replica was created
		CreatedAt time.Time `json:"created_at"`
		// UpdatedAt is the timestamp when the replica was last updated
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
		// StartedAt is the timestamp when the replica was started
		StartedAt *string `json:"started_at,omitempty"`
		// FinishedAt is the timestamp when the replica operation finished
		FinishedAt *string `json:"finished_at,omitempty"`
		// MaintenanceScheduledAt is the scheduled maintenance timestamp
		MaintenanceScheduledAt *string `json:"maintenance_scheduled_at,omitempty"`
	}

	// ReplicasResponse represents the response when listing replicas
	ReplicasResponse struct {
		// Meta contains pagination and filter information
		Meta MetaResponse `json:"meta"`
		// Results is the list of replicas
		Results []ReplicaDetailResponse `json:"results"`
	}

	// ReplicaCreateRequest represents the request payload for creating a replica
	ReplicaCreateRequest struct {
		// SourceID is the ID of the source instance to replicate
		SourceID string `json:"source_id"`
		// Name is the name of the replica
		Name string `json:"name"`
		// InstanceTypeID is the ID of the instance type for the replica
		InstanceTypeID *string `json:"instance_type_id,omitempty"`
	}

	// ReplicaResizeRequest represents the request payload for resizing a replica
	ReplicaResizeRequest struct {
		// InstanceTypeID is the new instance type ID
		InstanceTypeID string `json:"instance_type_id,omitempty"`
	}

	// ReplicaResponse represents the response when creating a replica
	ReplicaResponse struct {
		// ID is the unique identifier of the created replica
		ID string `json:"id"`
	}
)

// NewReplicaService creates a new replica service instance
func NewReplicaService(client *DBaaSClient) ReplicaService {
	return &replicaService{client: client}
}

// List returns a paginated list of database replicas with optional source_id filter
func (s *replicaService) List(ctx context.Context, opts ListReplicaOptions) ([]ReplicaDetailResponse, error) {
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

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ReplicasResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v2/replicas",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
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
