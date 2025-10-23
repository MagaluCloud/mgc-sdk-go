package compute

import (
	"context"
	"net/http"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// Meta contains pagination metadata for API responses.
// This structure provides information about the current page and total results.
type Meta struct {
	Page Page `json:"page"`
}

// Page contains pagination information
type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

// InstanceType represents a virtual machine instance type configuration.
// Each instance type defines the hardware specifications for virtual machines.
type InstanceType struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	VCPUs             int       `json:"vcpus"`
	RAM               int       `json:"ram"`
	Disk              int       `json:"disk"`
	GPU               *int      `json:"gpu,omitempty"`
	Status            string    `json:"status"`
	AvailabilityZones *[]string `json:"availability_zones,omitempty"`
}

// InstanceTypeList represents the response from listing instance types.
// This structure encapsulates the API response format for instance types.
type InstanceTypeList struct {
	InstanceTypes []InstanceType `json:"instance_types"`
	Meta          Meta           `json:"meta"`
}

// InstanceTypeService provides operations for querying available machine types.
// This interface allows listing instance types with optional filtering.
type InstanceTypeService interface {
	List(ctx context.Context, opts InstanceTypeListOptions) (*InstanceTypeList, error)
	ListAll(ctx context.Context, opts InstanceTypeFilterOptions) ([]InstanceType, error)
}

// instanceTypeService implements the InstanceTypeService interface.
// This is an internal implementation that should not be used directly.
type instanceTypeService struct {
	client *VirtualMachineClient
}

// InstanceTypeListOptions defines parameters for filtering and pagination of machine type lists.
// All fields are optional and allow controlling the listing behavior.
type InstanceTypeListOptions struct {
	Limit            *int    `json:"_limit,omitempty"`
	Offset           *int    `json:"_offset,omitempty"`
	Sort             *string `json:"_sort,omitempty"`
	AvailabilityZone string  `json:"availability-zone,omitempty"`
}

// InstanceTypeFilterOptions defines filtering options for ListAll (without pagination).
type InstanceTypeFilterOptions struct {
	Sort             *string `json:"_sort,omitempty"`
	AvailabilityZone string  `json:"availability-zone,omitempty"`
}

// List retrieves instance types with pagination metadata.
// This method makes an HTTP request to get the list of instance types
// and applies the filters specified in the options.
func (s *instanceTypeService) List(ctx context.Context, opts InstanceTypeListOptions) (*InstanceTypeList, error) {
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

	response := &InstanceTypeList{}
	_, err = mgc_http.Do(s.client.GetConfig(), ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListAll retrieves all instance types across all pages with optional filtering.
// This method automatically handles pagination and returns all results.
func (s *instanceTypeService) ListAll(ctx context.Context, opts InstanceTypeFilterOptions) ([]InstanceType, error) {
	var allInstanceTypes []InstanceType
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		listOpts := InstanceTypeListOptions{
			Offset:           &currentOffset,
			Limit:            &currentLimit,
			Sort:             opts.Sort,
			AvailabilityZone: opts.AvailabilityZone,
		}

		response, err := s.List(ctx, listOpts)
		if err != nil {
			return nil, err
		}

		allInstanceTypes = append(allInstanceTypes, response.InstanceTypes...)

		// Check if we've retrieved all results
		if len(response.InstanceTypes) < limit {
			break
		}

		offset += limit
	}

	return allInstanceTypes, nil
}
