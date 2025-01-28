package blockstorage

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	ListVolumeTypesResponse struct {
		Types []VolumeType `json:"types"`
	}

	VolumeType struct {
		ID                string         `json:"id"`
		Name              string         `json:"name"`
		DiskType          string         `json:"disk_type"`
		Status            string         `json:"status"`
		IOPS              VolumeTypeIOPS `json:"iops"`
		AvailabilityZones []string       `json:"availability_zones"`
		AllowsEncryption  bool           `json:"allows_encryption"`
	}

	VolumeTypeIOPS struct {
		Read  int `json:"read"`
		Write int `json:"write"`
		Total int `json:"total"`
	}
)

// DiskType represents the physical disk type
type DiskType string

const (
	DiskTypeNVMe DiskType = "nvme"
	DiskTypeHDD  DiskType = "hdd"
)

// ListVolumeTypesOptions contains the options for listing volume types
type ListVolumeTypesOptions struct {
	AvailabilityZone string
	Name             string
	AllowsEncryption *bool
}

// VolumeTypeService provides operations for managing volume types
type VolumeTypeService interface {
	// List returns all available volume types
	// Use options to filter by availability zone, encryption support, or name
	List(ctx context.Context, opts ListVolumeTypesOptions) ([]VolumeType, error)
}

type volumeTypeService struct {
	client *BlockStorageClient
}

// List retrieves all volume types with optional filtering
func (s *volumeTypeService) List(ctx context.Context, opts ListVolumeTypesOptions) ([]VolumeType, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/volume-types", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if opts.AvailabilityZone != "" {
		q.Add("availability-zone", opts.AvailabilityZone)
	}
	if opts.Name != "" {
		q.Add("name", opts.Name)
	}
	if opts.AllowsEncryption != nil {
		q.Add("allows-encryption", fmt.Sprintf("%v", *opts.AllowsEncryption))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &ListVolumeTypesResponse{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}

	return resp.Types, nil
}
