package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
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

	replicaService struct {
		client *DBaaSClient
	}

	ListReplicaOptions struct {
		Offset   *int
		Limit    *int
		SourceID *string
	}

	ReplicasResponse struct {
		Meta    MetaResponse            `json:"meta"`
		Results []ReplicaDetailResponse `json:"results"`
	}

	ReplicaCreateRequest struct {
		SourceID       string `json:"source_id"`
		Name          string `json:"name"`
		FlavorID      string `json:"flavor_id,omitempty"`
		InstanceTypeID string `json:"instance_type_id,omitempty"`
	}

	ReplicaResizeRequest struct {
		InstanceTypeID string `json:"instance_type_id,omitempty"`
		FlavorID      string `json:"flavor_id,omitempty"`
	}

	ReplicaResponse struct {
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
		"/v1/replicas",
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
		fmt.Sprintf("/v1/replicas/%s", id),
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
		"/v1/replicas",
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
		fmt.Sprintf("/v1/replicas/%s", id),
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
		fmt.Sprintf("/v1/replicas/%s/resize", id),
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
		fmt.Sprintf("/v1/replicas/%s/start", id),
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
		fmt.Sprintf("/v1/replicas/%s/stop", id),
		nil,
		nil,
	)
}
