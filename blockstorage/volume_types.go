// Package blockstorage provides functionality to interact with the MagaluCloud block storage service.
// This package allows managing volumes, volume types, and snapshots.
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
	// Types contains the list of available volume types
	Types []VolumeType `json:"types"`
}

// VolumeType represents a block storage volume type.
// Each volume type defines the characteristics and capabilities of volumes created with it.
type VolumeType struct {
	// ID is the unique identifier of the volume type
	ID string `json:"id"`
	// Name is the display name of the volume type
	Name string `json:"name"`
	// DiskType specifies the physical disk type (e.g., nvme, hdd)
	DiskType string `json:"disk_type"`
	// Status indicates the current status of the volume type
	Status string `json:"status"`
	// IOPS contains the input/output operations per second specifications
	IOPS VolumeTypeIOPS `json:"iops"`
	// AvailabilityZones lists the availability zones where this volume type is available
	AvailabilityZones []string `json:"availability_zones"`
	// AllowsEncryption indicates whether volumes of this type can be encrypted
	AllowsEncryption bool `json:"allows_encryption"`
}

// VolumeTypeIOPS represents the IOPS specifications for a volume type.
// IOPS defines the performance characteristics in terms of read/write operations.
type VolumeTypeIOPS struct {
	// Read specifies the read IOPS limit
	Read int `json:"read"`
	// Write specifies the write IOPS limit
	Write int `json:"write"`
	// Total specifies the total IOPS limit
	Total int `json:"total"`
}

// DiskType represents the physical disk type used for storage.
// Different disk types offer different performance characteristics and costs.
type DiskType string

const (
	// DiskTypeNVMe represents NVMe (Non-Volatile Memory Express) storage
	DiskTypeNVMe DiskType = "nvme"
	// DiskTypeHDD represents traditional hard disk drive storage
	DiskTypeHDD DiskType = "hdd"
)

// ListVolumeTypesOptions contains the options for listing volume types.
// All fields are optional and allow filtering the results.
type ListVolumeTypesOptions struct {
	// AvailabilityZone filters volume types by availability zone
	AvailabilityZone string
	// Name filters volume types by name
	Name string
	// AllowsEncryption filters volume types by encryption support
	AllowsEncryption *bool
}

// VolumeTypeService provides operations for managing volume types.
// This interface allows listing available volume types with optional filtering.
type VolumeTypeService interface {
	// List returns all available volume types.
	// Use options to filter by availability zone, encryption support, or name.
	//
	// Parameters:
	//   - ctx: Request context
	//   - opts: Options to filter the volume types
	//
	// Returns:
	//   - []VolumeType: List of available volume types
	//   - error: Error if there's a failure in the request
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
//
// Parameters:
//   - ctx: Request context
//   - opts: Options to filter the volume types
//
// Returns:
//   - []VolumeType: List of available volume types
//   - error: Error if there's a failure in the request
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
