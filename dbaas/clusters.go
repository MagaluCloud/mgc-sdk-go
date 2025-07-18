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

// ClusterStatus represents the possible states of a cluster
type ClusterStatus string

const (
	ClusterStatusActive        ClusterStatus = "ACTIVE"
	ClusterStatusError         ClusterStatus = "ERROR"
	ClusterStatusPending       ClusterStatus = "PENDING"
	ClusterStatusCreating      ClusterStatus = "CREATING"
	ClusterStatusDeleting      ClusterStatus = "DELETING"
	ClusterStatusDeleted       ClusterStatus = "DELETED"
	ClusterStatusErrorDeleting ClusterStatus = "ERROR_DELETING"
	ClusterStatusStarting      ClusterStatus = "STARTING"
	ClusterStatusStopping      ClusterStatus = "STOPPING"
	ClusterStatusStopped       ClusterStatus = "STOPPED"
	ClusterStatusBackingUp     ClusterStatus = "BACKING_UP"
)

type (
	// ListClustersOptions provides options for listing clusters
	ListClustersOptions struct {
		Offset           *int
		Limit            *int
		Status           *ClusterStatus
		EngineID         *string
		VolumeSize       *int
		VolumeSizeGt     *int
		VolumeSizeGte    *int
		VolumeSizeLt     *int
		VolumeSizeLte    *int
		ParameterGroupID *string
	}

	// ClustersResponse represents the response when listing clusters
	ClustersResponse struct {
		Results []ClusterDetailResponse `json:"results"`
	}

	// ClusterVolumeRequest represents volume configuration for cluster creation
	ClusterVolumeRequest struct {
		Size int     `json:"size"`
		Type *string `json:"type,omitempty"`
	}

	// ClusterCreateRequest represents the request payload for creating a cluster
	ClusterCreateRequest struct {
		Name                string               `json:"name"`
		EngineID            string               `json:"engine_id"`
		InstanceTypeID      string               `json:"instance_type_id"`
		User                string               `json:"user"`
		Password            string               `json:"password"`
		Volume              ClusterVolumeRequest `json:"volume"`
		ParameterGroupID    *string              `json:"parameter_group_id,omitempty"`
		BackupRetentionDays *int                 `json:"backup_retention_days,omitempty"`
		BackupStartAt       *string              `json:"backup_start_at,omitempty"`
	}

	// ClusterResponse represents the response when creating a cluster
	ClusterResponse struct {
		ID string `json:"id"`
	}

	// ClusterVolumeResponse represents volume information for a cluster
	ClusterVolumeResponse struct {
		Size int    `json:"size"`
		Type string `json:"type"`
	}

	// LoadBalancerAddress represents a load balancer address
	LoadBalancerAddress struct {
		Access  AddressAccess `json:"access"`
		Type    AddressType   `json:"type,omitempty"`
		Address string        `json:"address,omitempty"`
		Port    string        `json:"port,omitempty"`
	}

	// ClusterDetailResponse represents detailed information about a cluster
	ClusterDetailResponse struct {
		ID                     string                `json:"id"`
		Name                   string                `json:"name"`
		EngineID               string                `json:"engine_id"`
		InstanceTypeID         string                `json:"instance_type_id"`
		ParameterGroupID       string                `json:"parameter_group_id"`
		Volume                 ClusterVolumeResponse `json:"volume"`
		Status                 ClusterStatus         `json:"status"`
		Addresses              []LoadBalancerAddress `json:"addresses"`
		ApplyParametersPending bool                  `json:"apply_parameters_pending"`
		BackupRetentionDays    int                   `json:"backup_retention_days"`
		BackupStartAt          string                `json:"backup_start_at"`
		CreatedAt              time.Time             `json:"created_at"`
		UpdatedAt              *time.Time            `json:"updated_at,omitempty"`
		StartedAt              *string               `json:"started_at,omitempty"`
		FinishedAt             *string               `json:"finished_at,omitempty"`
	}

	// ClusterUpdateRequest represents the request payload for updating a cluster
	ClusterUpdateRequest struct {
		ParameterGroupID    *string `json:"parameter_group_id,omitempty"`
		BackupRetentionDays *int    `json:"backup_retention_days,omitempty"`
		BackupStartAt       *string `json:"backup_start_at,omitempty"`
	}
)

// ClusterService provides methods for managing database clusters
type ClusterService interface {
	List(ctx context.Context, opts ListClustersOptions) ([]ClusterDetailResponse, error)
	Create(ctx context.Context, req ClusterCreateRequest) (*ClusterResponse, error)
	Get(ctx context.Context, ID string) (*ClusterDetailResponse, error)
	Update(ctx context.Context, ID string, req ClusterUpdateRequest) (*ClusterDetailResponse, error)
	Delete(ctx context.Context, ID string) error
	Start(ctx context.Context, ID string) (*ClusterDetailResponse, error)
	Stop(ctx context.Context, ID string) (*ClusterDetailResponse, error)
}

// clusterService implements the ClusterService interface
type clusterService struct {
	client *DBaaSClient
}

// List implements the ClusterService interface
func (s *clusterService) List(ctx context.Context, opts ListClustersOptions) ([]ClusterDetailResponse, error) {
	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Status != nil {
		query.Set("status", string(*opts.Status))
	}
	if opts.EngineID != nil {
		query.Set("engine_id", *opts.EngineID)
	}
	if opts.VolumeSize != nil {
		query.Set("volume.size", strconv.Itoa(*opts.VolumeSize))
	}
	if opts.VolumeSizeGt != nil {
		query.Set("volume.size__gt", strconv.Itoa(*opts.VolumeSizeGt))
	}
	if opts.VolumeSizeGte != nil {
		query.Set("volume.size__gte", strconv.Itoa(*opts.VolumeSizeGte))
	}
	if opts.VolumeSizeLt != nil {
		query.Set("volume.size__lt", strconv.Itoa(*opts.VolumeSizeLt))
	}
	if opts.VolumeSizeLte != nil {
		query.Set("volume.size__lte", strconv.Itoa(*opts.VolumeSizeLte))
	}
	if opts.ParameterGroupID != nil {
		query.Set("parameter_group_id", *opts.ParameterGroupID)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ClustersResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v2/clusters",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

// Create implements the ClusterService interface
func (s *clusterService) Create(ctx context.Context, req ClusterCreateRequest) (*ClusterResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v2/clusters",
		req,
		nil,
	)
}

// Get implements the ClusterService interface
func (s *clusterService) Get(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v2/clusters/%s", ID),
		nil,
		nil,
	)
}

// Update implements the ClusterService interface
func (s *clusterService) Update(ctx context.Context, ID string, req ClusterUpdateRequest) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v2/clusters/%s", ID),
		req,
		nil,
	)
}

// Delete implements the ClusterService interface
func (s *clusterService) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v2/clusters/%s", ID),
		nil,
		nil,
	)
}

// Start implements the ClusterService interface
func (s *clusterService) Start(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v2/clusters/%s/start", ID),
		nil,
		nil,
	)
}

// Stop implements the ClusterService interface
func (s *clusterService) Stop(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v2/clusters/%s/stop", ID),
		nil,
		nil,
	)
}
