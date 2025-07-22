package blockstorage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ListVolumeTypesResponse represents the response from listing volume types.
// This structure encapsulates the API response format for volume types.
type ListVolumeTypesResponse struct {
	Types []VolumeType `json:"types"`
}

// VolumeType represents a block storage volume type.
// Each volume type defines the characteristics and capabilities of volumes created with it.
type VolumeType struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	DiskType          string         `json:"disk_type"`
	Status            string         `json:"status"`
	IOPS              VolumeTypeIOPS `json:"iops"`
	AvailabilityZones []string       `json:"availability_zones"`
	AllowsEncryption  bool           `json:"allows_encryption"`
}

// VolumeTypeIOPS represents the IOPS specifications for a volume type.
// IOPS defines the performance characteristics in terms of read/write operations.
type VolumeTypeIOPS struct {
	Read  int `json:"read"`
	Write int `json:"write"`
	Total int `json:"total"`
}

// DiskType represents the physical disk type used for storage.
// Different disk types offer different performance characteristics and costs.
type DiskType string

const (
	DiskTypeNVMe DiskType = "nvme"
	DiskTypeHDD  DiskType = "hdd"
)

// ListVolumeTypesOptions contains the options for listing volume types.
// All fields are optional and allow filtering the results.
type ListVolumeTypesOptions struct {
	AvailabilityZone string
	Name             string
	AllowsEncryption *bool
}

// VolumeTypeService provides operations for managing volume types.
// This interface allows listing available volume types with optional filtering.
type VolumeTypeService interface {
	List(ctx context.Context, opts ListVolumeTypesOptions) ([]VolumeType, error)
}

// volumeTypeService implements the VolumeTypeService interface.
// This is an internal implementation that should not be used directly.
type volumeTypeService struct {
	client *BlockStorageClient
}

// List retrieves all volume types with optional filtering.
// This method makes an HTTP request to get the list of volume types
// and applies the filters specified in the options.
func (s *volumeTypeService) List(ctx context.Context, opts ListVolumeTypesOptions) ([]VolumeType, error) {
	queryParams := make(url.Values)
	if opts.AvailabilityZone != "" {
		queryParams.Add("availability-zone", opts.AvailabilityZone)
	}
	if opts.Name != "" {
		queryParams.Add("name", opts.Name)
	}
	if opts.AllowsEncryption != nil {
		queryParams.Add("allows-encryption", fmt.Sprintf("%v", *opts.AllowsEncryption))
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListVolumeTypesResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v1/volume-types",
		nil,
		queryParams,
	)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}

	return resp.Types, nil
}
