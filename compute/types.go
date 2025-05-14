package compute

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

// InstanceType represents a virtual machine instance type configuration
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

type InstanceTypeList struct {
	InstanceTypes []InstanceType `json:"instance_types"`
	Meta          Meta           `json:"meta"`
}

// InstanceTypeService provides operations for querying available machine types
type InstanceTypeService interface {
	// List returns all available machine types with optional filtering
	List(ctx context.Context, opts InstanceTypeListOptions) ([]InstanceType, error)
}

type instanceTypeService struct {
	client *VirtualMachineClient
}

// InstanceTypeListOptions defines parameters for filtering and pagination of machine type lists
type InstanceTypeListOptions struct {
	Limit            *int    `url:"_limit,omitempty"`
	Offset           *int    `url:"_offset,omitempty"`
	Sort             *string `url:"_sort,omitempty"`
	AvailabilityZone string  `url:"availability-zone,omitempty"`
}

// List retrieves all available machine types
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
