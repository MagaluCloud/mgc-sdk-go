package compute

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	ImageList struct {
		Images []Image `json:"images"`
	}

	// Image represents a virtual machine image
	Image struct {
		ID                   string              `json:"id"`
		Name                 string              `json:"name"`
		Status               ImageStatus         `json:"status"`
		Version              *string             `json:"version,omitempty"`
		Platform             *string             `json:"platform,omitempty"`
		ReleaseAt            *string             `json:"release_at,omitempty"`
		EndStandardSupportAt *string             `json:"end_standard_support_at,omitempty"`
		EndLifeAt            *string             `json:"end_life_at,omitempty"`
		MinimumRequirements  MinimumRequirements `json:"minimum_requirements"`
		Labels               *[]string           `json:"labels,omitempty"`
		AvailabilityZones    *[]string           `json:"availability_zones,omitempty"`
	}

	// MinimumRequirements represents the minimum hardware requirements for an image
	MinimumRequirements struct {
		VCPU int `json:"vcpu"`
		RAM  int `json:"ram"`
		Disk int `json:"disk"`
	}

	// ImageStatus represents the current state of an image
	ImageStatus string
)

const (
	ImageStatusActive     ImageStatus = "active"
	ImageStatusDeprecated ImageStatus = "deprecated"
	ImageStatusDeleted    ImageStatus = "deleted"
	ImageStatusPending    ImageStatus = "pending"
	ImageStatusCreating   ImageStatus = "creating"
	ImageStatusImporting  ImageStatus = "importing"
	ImageStatusError      ImageStatus = "error"
)

// ImageService provides operations for managing virtual machine images
type ImageService interface {
	// List returns a slice of images based on the provided listing options
	List(ctx context.Context, opts ImageListOptions) ([]Image, error)
}

type imageService struct {
	client *VirtualMachineClient
}

// ImageListOptions defines the parameters for filtering and pagination of image lists
type ImageListOptions struct {
	// Limit specifies the maximum number of results to return (default: 50)
	Limit *int
	// Offset specifies the number of results to skip for pagination
	Offset *int
	// Sort defines the field and direction for result ordering (default: "platform:asc,end_life_at:desc")
	Sort *string
	// Labels filters images by their labels
	Labels []string
	// AvailabilityZone filters images by availability zone
	AvailabilityZone *string
}

// List retrieves all images matching the provided options
func (s *imageService) List(ctx context.Context, opts ImageListOptions) ([]Image, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/images", nil)
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
	if len(opts.Labels) > 0 {
		q.Add("_labels", strings.Join(opts.Labels, ","))
	}
	if opts.AvailabilityZone != nil {
		q.Add("availability-zone", *opts.AvailabilityZone)
	}
	req.URL.RawQuery = q.Encode()

	response := &ImageList{}

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, response)
	if err != nil {
		return nil, err
	}

	return resp.Images, nil
}
