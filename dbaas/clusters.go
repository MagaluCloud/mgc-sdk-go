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
	// ClusterStatusActive indicates the cluster is running normally
	ClusterStatusActive ClusterStatus = "ACTIVE"
	// ClusterStatusError indicates the cluster is in an error state
	ClusterStatusError ClusterStatus = "ERROR"
	// ClusterStatusPending indicates the cluster is waiting to be processed
	ClusterStatusPending ClusterStatus = "PENDING"
	// ClusterStatusCreating indicates the cluster is being created
	ClusterStatusCreating ClusterStatus = "CREATING"
	// ClusterStatusDeleting indicates the cluster is being deleted
	ClusterStatusDeleting ClusterStatus = "DELETING"
	// ClusterStatusDeleted indicates the cluster has been deleted
	ClusterStatusDeleted ClusterStatus = "DELETED"
	// ClusterStatusErrorDeleting indicates an error occurred during deletion
	ClusterStatusErrorDeleting ClusterStatus = "ERROR_DELETING"
	// ClusterStatusStarting indicates the cluster is starting up
	ClusterStatusStarting ClusterStatus = "STARTING"
	// ClusterStatusStopping indicates the cluster is shutting down
	ClusterStatusStopping ClusterStatus = "STOPPING"
	// ClusterStatusStopped indicates the cluster is stopped
	ClusterStatusStopped ClusterStatus = "STOPPED"
	// ClusterStatusBackingUp indicates the cluster is performing a backup
	ClusterStatusBackingUp ClusterStatus = "BACKING_UP"
)

type (
	// ListClustersOptions provides options for listing clusters
	ListClustersOptions struct {
		// Offset is the number of clusters to skip
		Offset *int
		// Limit is the maximum number of clusters to return
		Limit *int
		// Status filters clusters by status
		Status *ClusterStatus
		// EngineID filters clusters by database engine
		EngineID *string
		// VolumeSize filters clusters by exact volume size
		VolumeSize *int
		// VolumeSizeGt filters clusters by volume size greater than
		VolumeSizeGt *int
		// VolumeSizeGte filters clusters by volume size greater than or equal
		VolumeSizeGte *int
		// VolumeSizeLt filters clusters by volume size less than
		VolumeSizeLt *int
		// VolumeSizeLte filters clusters by volume size less than or equal
		VolumeSizeLte *int
		// ParameterGroupID filters clusters by parameter group
		ParameterGroupID *string
	}

	// ClustersResponse represents the response when listing clusters
	ClustersResponse struct {
		// Results is the list of clusters
		Results []ClusterDetailResponse `json:"results"`
	}

	// ClusterVolumeRequest represents volume configuration for cluster creation
	ClusterVolumeRequest struct {
		// Size is the volume size in GB
		Size int `json:"size"`
		// Type is the volume type (optional)
		Type *string `json:"type,omitempty"`
	}

	// ClusterCreateRequest represents the request payload for creating a cluster
	ClusterCreateRequest struct {
		// Name is the name of the cluster
		Name string `json:"name"`
		// EngineID is the ID of the database engine
		EngineID string `json:"engine_id"`
		// InstanceTypeID is the ID of the instance type
		InstanceTypeID string `json:"instance_type_id"`
		// User is the database user
		User string `json:"user"`
		// Password is the database password
		Password string `json:"password"`
		// Volume contains volume configuration
		Volume ClusterVolumeRequest `json:"volume"`
		// ParameterGroupID is the ID of the parameter group (optional)
		ParameterGroupID *string `json:"parameter_group_id,omitempty"`
		// BackupRetentionDays is the number of days to retain backups (optional)
		BackupRetentionDays *int `json:"backup_retention_days,omitempty"`
		// BackupStartAt is the backup start time (optional)
		BackupStartAt *string `json:"backup_start_at,omitempty"`
	}

	// ClusterResponse represents the response when creating a cluster
	ClusterResponse struct {
		// ID is the unique identifier of the created cluster
		ID string `json:"id"`
	}

	// ClusterVolumeResponse represents volume information for a cluster
	ClusterVolumeResponse struct {
		// Size is the volume size in GB
		Size int `json:"size"`
		// Type is the volume type
		Type string `json:"type"`
	}

	// LoadBalancerAddress represents a load balancer address
	LoadBalancerAddress struct {
		// Access indicates the access type (PRIVATE or PUBLIC)
		Access AddressAccess `json:"access"`
		// Type indicates the address type (IPv4 or IPv6)
		Type AddressType `json:"type,omitempty"`
		// Address is the IP address
		Address string `json:"address,omitempty"`
		// Port is the port number
		Port string `json:"port,omitempty"`
	}

	// ClusterDetailResponse represents detailed information about a cluster
	ClusterDetailResponse struct {
		// ID is the unique identifier of the cluster
		ID string `json:"id"`
		// Name is the name of the cluster
		Name string `json:"name"`
		// EngineID is the ID of the database engine
		EngineID string `json:"engine_id"`
		// InstanceTypeID is the ID of the instance type
		InstanceTypeID string `json:"instance_type_id"`
		// ParameterGroupID is the ID of the parameter group
		ParameterGroupID string `json:"parameter_group_id"`
		// Volume contains volume information
		Volume ClusterVolumeResponse `json:"volume"`
		// Status is the current status of the cluster
		Status ClusterStatus `json:"status"`
		// Addresses contains load balancer addresses
		Addresses []LoadBalancerAddress `json:"addresses"`
		// ApplyParametersPending indicates if parameter changes are pending
		ApplyParametersPending bool `json:"apply_parameters_pending"`
		// BackupRetentionDays is the number of days to retain backups
		BackupRetentionDays int `json:"backup_retention_days"`
		// BackupStartAt is the backup start time
		BackupStartAt string `json:"backup_start_at"`
		// CreatedAt is the timestamp when the cluster was created
		CreatedAt time.Time `json:"created_at"`
		// UpdatedAt is the timestamp when the cluster was last updated
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
		// StartedAt is the timestamp when the cluster was started
		StartedAt *string `json:"started_at,omitempty"`
		// FinishedAt is the timestamp when the cluster operation finished
		FinishedAt *string `json:"finished_at,omitempty"`
	}

	// ClusterUpdateRequest represents the request payload for updating a cluster
	ClusterUpdateRequest struct {
		// ParameterGroupID is the new parameter group ID (optional)
		ParameterGroupID *string `json:"parameter_group_id,omitempty"`
		// BackupRetentionDays is the new backup retention days (optional)
		BackupRetentionDays *int `json:"backup_retention_days,omitempty"`
		// BackupStartAt is the new backup start time (optional)
		BackupStartAt *string `json:"backup_start_at,omitempty"`
	}
)

// ClusterService provides methods for managing database clusters
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
