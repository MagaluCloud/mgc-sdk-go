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

// FilterMetadata represents a filter applied to a list request
type FilterMetadata struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// Metadata represents pagination metadata
type Metadata struct {
	Filters []FilterMetadata `json:"filters"`
	Page    PageMetadata     `json:"page"`
}

// PageMetadata contains pagination information
type PageMetadata struct {
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
	Count    int `json:"count"`
	Total    int `json:"total"`
	MaxLimit int `json:"max_limit"`
}

// ClusterStatus represents the possible states of a cluster
type ClusterStatus string

const (
	ClusterStatusActive             ClusterStatus = "ACTIVE"
	ClusterStatusError              ClusterStatus = "ERROR"
	ClusterStatusPending            ClusterStatus = "PENDING"
	ClusterStatusCreating           ClusterStatus = "CREATING"
	ClusterStatusDeleting           ClusterStatus = "DELETING"
	ClusterStatusDeleted            ClusterStatus = "DELETED"
	ClusterStatusErrorDeleting      ClusterStatus = "ERROR_DELETING"
	ClusterStatusStarting           ClusterStatus = "STARTING"
	ClusterStatusStopping           ClusterStatus = "STOPPING"
	ClusterStatusStopped            ClusterStatus = "STOPPED"
	ClusterStatusBackingUp          ClusterStatus = "BACKING_UP"
	ClusterStatusStartingImportMode ClusterStatus = "STARTING_IMPORT_MODE"
	ClusterStatusImportMode         ClusterStatus = "IMPORT_MODE"
	ClusterStatusStoppingImportMode ClusterStatus = "STOPPING_IMPORT_MODE"
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

	// ClusterFilterOptions provides filtering options for ListAll (without pagination)
	ClusterFilterOptions struct {
		Status           *ClusterStatus
		EngineID         *string
		VolumeSize       *int
		VolumeSizeGt     *int
		VolumeSizeGte    *int
		VolumeSizeLt     *int
		VolumeSizeLte    *int
		ParameterGroupID *string
	}

	// AddressPurpose represents the network address purpose on a cluster
	AddressPurpose string

	// ClustersResponse represents the response when listing clusters
	ClustersResponse struct {
		Meta    Metadata                `json:"meta"`
		Results []ClusterDetailResponse `json:"results"`
	}

	// ClusterVolumeRequest represents volume configuration for cluster creation
	ClusterVolumeRequest struct {
		Size int     `json:"size"`
		Type *string `json:"type,omitempty"`
	}

	ClusterVolumeResizeRequest struct {
		Size int    `json:"size"`
		Type string `json:"type,omitempty"`
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
		DeletionProtected   *bool                `json:"deletion_protected,omitempty"`
	}

	ClusterResizeRequest struct {
		InstanceTypeID *string                     `json:"instance_type_id,omitempty"`
		Volume         *ClusterVolumeResizeRequest `json:"volume,omitempty"`
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
		Access  AddressAccess  `json:"access"`
		Type    AddressType    `json:"type,omitempty"`
		Address string         `json:"address,omitempty"`
		Port    string         `json:"port,omitempty"`
		Purpose AddressPurpose `json:"purpose,omitempty"`
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
		DeletionProtected      bool                  `json:"deletion_protected"`
	}

	// ClusterUpdateRequest represents the request payload for updating a cluster
	ClusterUpdateRequest struct {
		ParameterGroupID    *string `json:"parameter_group_id,omitempty"`
		BackupRetentionDays *int    `json:"backup_retention_days,omitempty"`
		BackupStartAt       *string `json:"backup_start_at,omitempty"`
		DeletionProtected   *bool   `json:"deletion_protected,omitempty"`
	}

	// ListClusterSnapshotOptions provides options for listing cluster snapshots
	ListClusterSnapshotOptions struct {
		Offset *int
		Limit  *int
		Type   *SnapshotType
		Status *SnapshotStatus
	}

	// ClusterSnapshotFilterOptions provides filtering options for ListAllSnapshots (without pagination)
	ClusterSnapshotFilterOptions struct {
		Type   *SnapshotType
		Status *SnapshotStatus
	}

	// ClusterSnapshotInfo represents cluster information referenced by a cluster snapshot
	ClusterSnapshotInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// ClusterSnapshotDetailResponse represents detailed information about a cluster snapshot
	ClusterSnapshotDetailResponse struct {
		ID            string              `json:"id"`
		Cluster       ClusterSnapshotInfo `json:"cluster"`
		Name          string              `json:"name"`
		Description   string              `json:"description"`
		Type          SnapshotType        `json:"type"`
		Status        SnapshotStatus      `json:"status"`
		AllocatedSize int                 `json:"allocated_size"`
		CreatedAt     time.Time           `json:"created_at"`
		StartedAt     *string             `json:"started_at,omitempty"`
		FinishedAt    *string             `json:"finished_at,omitempty"`
		UpdatedAt     *time.Time          `json:"updated_at,omitempty"`
	}

	// ClusterSnapshotsResponse represents the response when listing cluster snapshots
	ClusterSnapshotsResponse struct {
		Meta    Metadata                        `json:"meta"`
		Results []ClusterSnapshotDetailResponse `json:"results"`
	}

	// ClusterSnapshotCreateRequest represents the request payload for creating a cluster snapshot
	ClusterSnapshotCreateRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
	}

	// ClusterSnapshotResponse represents the response when creating a cluster snapshot
	ClusterSnapshotResponse struct {
		ID string `json:"id"`
	}

	// ClusterSnapshotUpdateRequest represents the request payload for updating a cluster snapshot
	ClusterSnapshotUpdateRequest struct {
		Name        string  `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// ClusterRestoreRequest represents the request payload for restoring a cluster from a snapshot
	ClusterRestoreRequest struct {
		Name                string                `json:"name"`
		InstanceTypeID      string                `json:"instance_type_id"`
		Volume              *ClusterVolumeRequest `json:"volume,omitempty"`
		BackupRetentionDays *int                  `json:"backup_retention_days,omitempty"`
		BackupStartAt       *string               `json:"backup_start_at,omitempty"`
	}
)

// ClusterService provides methods for managing database clusters
type ClusterService interface {
	List(ctx context.Context, opts ListClustersOptions) (*ClustersResponse, error)
	ListAll(ctx context.Context, filterOpts ClusterFilterOptions) ([]ClusterDetailResponse, error)
	Create(ctx context.Context, req ClusterCreateRequest) (*ClusterResponse, error)
	Get(ctx context.Context, ID string) (*ClusterDetailResponse, error)
	Update(ctx context.Context, ID string, req ClusterUpdateRequest) (*ClusterDetailResponse, error)
	Resize(ctx context.Context, id string, req ClusterResizeRequest) (*ClusterDetailResponse, error)
	Delete(ctx context.Context, ID string) error
	Start(ctx context.Context, ID string) (*ClusterDetailResponse, error)
	Stop(ctx context.Context, ID string) (*ClusterDetailResponse, error)
	StartImportMode(ctx context.Context, ID string) (*ClusterDetailResponse, error)
	StopImportMode(ctx context.Context, ID string) (*ClusterDetailResponse, error)
	ListSnapshots(ctx context.Context, clusterID string, opts ListClusterSnapshotOptions) (*ClusterSnapshotsResponse, error)
	ListAllSnapshots(ctx context.Context, clusterID string, filterOpts ClusterSnapshotFilterOptions) ([]ClusterSnapshotDetailResponse, error)
	CreateSnapshot(ctx context.Context, clusterID string, req ClusterSnapshotCreateRequest) (*ClusterSnapshotResponse, error)
	GetSnapshot(ctx context.Context, clusterID, snapshotID string) (*ClusterSnapshotDetailResponse, error)
	UpdateSnapshot(ctx context.Context, clusterID, snapshotID string, req ClusterSnapshotUpdateRequest) (*ClusterSnapshotDetailResponse, error)
	DeleteSnapshot(ctx context.Context, clusterID, snapshotID string) error
	RestoreSnapshot(ctx context.Context, clusterID, snapshotID string, req ClusterRestoreRequest) (*ClusterDetailResponse, error)
}

// clusterService implements the ClusterService interface
type clusterService struct {
	client *DBaaSClient
}

const v2ClustersPath = "/v2/clusters"
const errIDCannotBeEmpty = "ID cannot be empty"

const (
	ClusterSnapshotPath   = v2ClustersPath + "/%s/snapshots"
	ClusterSnapshotPathID = ClusterSnapshotPath + "/%s"
)

// List implements the ClusterService interface and returns a paginated list of clusters
func (s *clusterService) List(ctx context.Context, opts ListClustersOptions) (*ClustersResponse, error) {
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
		v2ClustersPath,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListAll retrieves all clusters by fetching all pages with optional filtering
func (s *clusterService) ListAll(ctx context.Context, filterOpts ClusterFilterOptions) ([]ClusterDetailResponse, error) {
	var allClusters []ClusterDetailResponse
	offset := 0
	limit := 25

	for {
		currentOffset := offset
		currentLimit := limit
		opts := ListClustersOptions{
			Offset:           &currentOffset,
			Limit:            &currentLimit,
			Status:           filterOpts.Status,
			EngineID:         filterOpts.EngineID,
			VolumeSize:       filterOpts.VolumeSize,
			VolumeSizeGt:     filterOpts.VolumeSizeGt,
			VolumeSizeGte:    filterOpts.VolumeSizeGte,
			VolumeSizeLt:     filterOpts.VolumeSizeLt,
			VolumeSizeLte:    filterOpts.VolumeSizeLte,
			ParameterGroupID: filterOpts.ParameterGroupID,
		}

		resp, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allClusters = append(allClusters, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allClusters, nil
}

// Create implements the ClusterService interface
func (s *clusterService) Create(ctx context.Context, req ClusterCreateRequest) (*ClusterResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		v2ClustersPath,
		req,
		nil,
	)
}

// Get implements the ClusterService interface
func (s *clusterService) Get(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("%s/%s", v2ClustersPath, ID),
		nil,
		nil,
	)
}

// Update implements the ClusterService interface
func (s *clusterService) Update(ctx context.Context, ID string, req ClusterUpdateRequest) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("%s/%s", v2ClustersPath, ID),
		req,
		nil,
	)
}

// Resize changes the instance type and/or volume specifications of a database instance.
// Returns the instance details with the new specifications.
func (s *clusterService) Resize(ctx context.Context, ID string, req ClusterResizeRequest) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("%s/%s/resize", v2ClustersPath, ID),
		req,
		nil,
	)
}

// Delete implements the ClusterService interface
func (s *clusterService) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("%s/%s", v2ClustersPath, ID),
		nil,
		nil,
	)
}

// Start implements the ClusterService interface
func (s *clusterService) Start(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("%s/%s/start", v2ClustersPath, ID),
		nil,
		nil,
	)
}

// Stop implements the ClusterService interface
func (s *clusterService) Stop(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("%s/%s/stop", v2ClustersPath, ID),
		nil,
		nil,
	)
}

// StartImportMode enables optimized mode for bulk data import
func (s *clusterService) StartImportMode(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("%s/%s/start-import-mode", v2ClustersPath, ID),
		nil,
		nil,
	)
}

// StopImportMode disables optimized mode for bulk data import
func (s *clusterService) StopImportMode(ctx context.Context, ID string) (*ClusterDetailResponse, error) {
	if ID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("%s/%s/stop-import-mode", v2ClustersPath, ID),
		nil,
		nil,
	)
}

// ListSnapshots returns a list of snapshots for a specific cluster.
func (s *clusterService) ListSnapshots(ctx context.Context, clusterID string, opts ListClusterSnapshotOptions) (*ClusterSnapshotsResponse, error) {
	if clusterID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	query := make(url.Values)

	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Type != nil {
		query.Set("type", string(*opts.Type))
	}
	if opts.Status != nil {
		query.Set("status", string(*opts.Status))
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ClusterSnapshotsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(ClusterSnapshotPath, clusterID),
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListAllSnapshots retrieves all snapshots for a cluster by fetching all pages with optional filtering
func (s *clusterService) ListAllSnapshots(ctx context.Context, clusterID string, filterOpts ClusterSnapshotFilterOptions) ([]ClusterSnapshotDetailResponse, error) {
	var allSnapshots []ClusterSnapshotDetailResponse
	offset := 0
	limit := 25

	for {
		currentOffset := offset
		currentLimit := limit
		opts := ListClusterSnapshotOptions{
			Offset: &currentOffset,
			Limit:  &currentLimit,
			Type:   filterOpts.Type,
			Status: filterOpts.Status,
		}

		resp, err := s.ListSnapshots(ctx, clusterID, opts)
		if err != nil {
			return nil, err
		}

		allSnapshots = append(allSnapshots, resp.Results...)

		if len(resp.Results) < limit {
			break
		}

		offset += limit
	}

	return allSnapshots, nil
}

// CreateSnapshot creates a new snapshot for the specified cluster.
func (s *clusterService) CreateSnapshot(ctx context.Context, clusterID string, req ClusterSnapshotCreateRequest) (*ClusterSnapshotResponse, error) {
	if clusterID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterSnapshotResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(ClusterSnapshotPath, clusterID),
		req,
		nil,
	)
}

// GetSnapshot retrieves details of a specific cluster snapshot.
func (s *clusterService) GetSnapshot(ctx context.Context, clusterID, snapshotID string) (*ClusterSnapshotDetailResponse, error) {
	if clusterID == "" || snapshotID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterSnapshotDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(ClusterSnapshotPathID, clusterID, snapshotID),
		nil,
		nil,
	)
}

// UpdateSnapshot updates the properties of an existing cluster snapshot.
func (s *clusterService) UpdateSnapshot(ctx context.Context, clusterID, snapshotID string, req ClusterSnapshotUpdateRequest) (*ClusterSnapshotDetailResponse, error) {
	if clusterID == "" || snapshotID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterSnapshotDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(ClusterSnapshotPathID, clusterID, snapshotID),
		req,
		nil,
	)
}

// DeleteSnapshot deletes a cluster snapshot.
func (s *clusterService) DeleteSnapshot(ctx context.Context, clusterID, snapshotID string) error {
	if clusterID == "" || snapshotID == "" {
		return fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(ClusterSnapshotPathID, clusterID, snapshotID),
		nil,
		nil,
	)
}

// RestoreSnapshot creates a new cluster from a snapshot.
func (s *clusterService) RestoreSnapshot(ctx context.Context, clusterID, snapshotID string, req ClusterRestoreRequest) (*ClusterDetailResponse, error) {
	if clusterID == "" || snapshotID == "" {
		return nil, fmt.Errorf(errIDCannotBeEmpty)
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ClusterDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(ClusterSnapshotPathID+"/restore", clusterID, snapshotID),
		req,
		nil,
	)
}
