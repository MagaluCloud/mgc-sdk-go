package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

	ClustersResponse struct {
		Results []ClusterDetailResponse `json:"results"`
	}

	ClusterVolumeRequest struct {
		Size int     `json:"size"`
		Type *string `json:"type,omitempty"`
	}

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

	ClusterResponse struct {
		ID string `json:"id"`
	}

	ClusterVolumeResponse struct {
		Size int    `json:"size"`
		Type string `json:"type"`
	}

	LoadBalancerAddress struct {
		Access  AddressAccess `json:"access"`
		Type    AddressType   `json:"type,omitempty"`
		Address string        `json:"address,omitempty"`
		Port    string        `json:"port,omitempty"`
	}

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
		CreatedAt              string                `json:"created_at"`
		UpdatedAt              *string               `json:"updated_at,omitempty"`
		StartedAt              *string               `json:"started_at,omitempty"`
		FinishedAt             *string               `json:"finished_at,omitempty"`
	}

	ClusterUpdateRequest struct {
		ParameterGroupID    *string `json:"parameter_group_id,omitempty"`
		BackupRetentionDays *int    `json:"backup_retention_days,omitempty"`
		BackupStartAt       *string `json:"backup_start_at,omitempty"`
	}
)

type ClusterService interface {
	// List retrieves a list of database clusters for the tenant
	List(ctx context.Context, opts ListClustersOptions) ([]ClusterDetailResponse, error)

	// Create creates a new database high availability cluster asynchronously
	Create(ctx context.Context, req ClusterCreateRequest) (*ClusterResponse, error)

	// Get retrieves details of a specific cluster by its ID
	Get(ctx context.Context, ID string) (*ClusterDetailResponse, error)

	// Update updates properties of a database cluster
	Update(ctx context.Context, ID string, req ClusterUpdateRequest) (*ClusterDetailResponse, error)

	// Delete deletes a database cluster asynchronously
	Delete(ctx context.Context, ID string) error

	// Start starts a database cluster
	Start(ctx context.Context, ID string) (*ClusterDetailResponse, error)

	// Stop stops a database cluster
	Stop(ctx context.Context, ID string) (*ClusterDetailResponse, error)
}

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
	path := fmt.Sprintf("/v2/clusters/%s", ID)

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		nil,
	)
}

// Update implements the ClusterService interface
func (s *clusterService) Update(ctx context.Context, ID string, req ClusterUpdateRequest) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	path := fmt.Sprintf("/v2/clusters/%s", ID)

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		path,
		req,
		nil,
	)
}

// Delete implements the ClusterService interface
func (s *clusterService) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	path := fmt.Sprintf("/v2/clusters/%s", ID)

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		path,
		nil,
		nil,
	)
}

// Start implements the ClusterService interface
func (s *clusterService) Start(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	path := fmt.Sprintf("/v2/clusters/%s/start", ID)

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		path,
		nil,
		nil,
	)
}

// Stop implements the ClusterService interface
func (s *clusterService) Stop(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	path := fmt.Sprintf("/v2/clusters/%s/stop", ID)

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		path,
		nil,
		nil,
	)
}
