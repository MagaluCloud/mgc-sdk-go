// Package compute provides functionality to interact with the MagaluCloud compute service.
// This package allows managing virtual machine instances, images, instance types, and snapshots.
package compute

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// Meta contains pagination metadata for API responses.
// This structure provides information about the current page and total results.
type Meta struct {
	// Limit is the maximum number of results per page
	Limit int `json:"limit"`
	// Offset is the number of results skipped
	Offset int `json:"offset"`
	// Count is the number of results in the current page
	Count int `json:"count"`
	// Total is the total number of available results
	Total int `json:"total"`
}

// InstanceType represents a virtual machine instance type configuration.
// Each instance type defines the hardware specifications for virtual machines.
type InstanceType struct {
	// ID is the unique identifier of the instance type
	ID string `json:"id"`
	// Name is the display name of the instance type
	Name string `json:"name"`
	// VCPUs is the number of virtual CPUs
	VCPUs int `json:"vcpus"`
	// RAM is the amount of RAM in MB
	RAM int `json:"ram"`
	// Disk is the disk size in GB
	Disk int `json:"disk"`
	// GPU is the number of GPUs (optional)
	GPU *int `json:"gpu,omitempty"`
	// Status indicates the current status of the instance type
	Status string `json:"status"`
	// AvailabilityZones lists the availability zones where this instance type is available
	AvailabilityZones *[]string `json:"availability_zones,omitempty"`
}

// InstanceTypeList represents the response from listing instance types.
// This structure encapsulates the API response format for instance types.
type InstanceTypeList struct {
	// InstanceTypes contains the list of available instance types
	InstanceTypes []InstanceType `json:"instance_types"`
	// Meta contains pagination information
	Meta Meta `json:"meta"`
}

// InstanceTypeService provides operations for querying available machine types.
// This interface allows listing instance types with optional filtering.
type InstanceTypeService interface {
	// List returns all available machine types with optional filtering.
	//
	// Parameters:
	//   - ctx: Request context
	//   - opts: Options to filter and paginate the instance types
	//
	// Returns:
	//   - []InstanceType: List of available instance types
	//   - error: Error if there's a failure in the request
	List(ctx context.Context, opts InstanceTypeListOptions) ([]InstanceType, error)
}

// instanceTypeService implements the InstanceTypeService interface.
// This is an internal implementation that should not be used directly.
type instanceTypeService struct {
	client *VirtualMachineClient
}

// InstanceTypeListOptions defines parameters for filtering and pagination of machine type lists.
// All fields are optional and allow controlling the listing behavior.
type InstanceTypeListOptions struct {
	// Limit defines the maximum number of results to be returned
	Limit *int `url:"_limit,omitempty"`
	// Offset defines the number of results to be skipped (for pagination)
	Offset *int `url:"_offset,omitempty"`
	// Sort specifies the sorting criteria
	Sort *string `url:"_sort,omitempty"`
	// AvailabilityZone filters instance types by availability zone
	AvailabilityZone string `url:"availability-zone,omitempty"`
}

// List retrieves all available machine types.
// This method makes an HTTP request to get the list of instance types
// and applies the filters specified in the options.
//
// Parameters:
//   - ctx: Request context
//   - opts: Options to filter and paginate the instance types
//
// Returns:
//   - []InstanceType: List of available instance types
//   - error: Error if there's a failure in the request
func (s *instanceTypeService) List(ctx context.Context, opts InstanceTypeListOptions) ([]InstanceType, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/instance-types", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if opts.Limit != nil {
		q.Add("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Add("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Add("_sort", *opts.Sort)
	}
	if opts.AvailabilityZone != "" {
		q.Add("availability-zone", opts.AvailabilityZone)
	}
	req.URL.RawQuery = q.Encode()

	var response InstanceTypeList
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}

	return response.InstanceTypes, nil

}
