package compute

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// MachineType represents a virtual machine instance type configuration
type MachineType struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	VCPUs             int       `json:"vcpus"`
	RAM               int       `json:"ram"`
	Disk              int       `json:"disk"`
	GPU               *int      `json:"gpu,omitempty"`
	Status            string    `json:"status"`
	SKU               *string   `json:"sku,omitempty"`
	AvailabilityZones *[]string `json:"availability_zones,omitempty"`
}

// MachineTypeService provides operations for querying available machine types
type MachineTypeService interface {
	// List returns all available machine types with optional filtering
	List(ctx context.Context, opts MachineTypeListOptions) ([]MachineType, error)
}

type machineTypeService struct {
	client *VirtualMachineClient
}

// MachineTypeListOptions defines parameters for filtering and pagination of machine type lists
type MachineTypeListOptions struct {
	Limit            *int    `url:"_limit,omitempty"`
	Offset           *int    `url:"_offset,omitempty"`
	Sort             *string `url:"_sort,omitempty"`
	AvailabilityZone string  `url:"availability-zone,omitempty"`
}

// List retrieves all available machine types
func (s *machineTypeService) List(ctx context.Context, opts MachineTypeListOptions) ([]MachineType, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/machine-types", nil)
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

	var response struct {
		MachineTypes []MachineType `json:"machine_types"`
	}
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}

	return resp.MachineTypes, nil

}
