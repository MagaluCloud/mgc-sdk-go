// Package compute provides functionality to interact with the MagaluCloud compute service.
// This package allows managing virtual machine instances, images, instance types, and snapshots.
package compute

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ImageList represents the response from listing images.
// This structure encapsulates the API response format for images.
type ImageList struct {
	// Images contains the list of available images
	Images []Image `json:"images"`
}

// Image represents a virtual machine image.
// An image is a template that contains the operating system and software for creating instances.
type Image struct {
	// ID is the unique identifier of the image
	ID string `json:"id"`
	// Name is the display name of the image
	Name string `json:"name"`
	// Status indicates the current status of the image
	Status ImageStatus `json:"status"`
	// Version specifies the version of the image
	Version *string `json:"version,omitempty"`
	// Platform indicates the operating system platform
	Platform *string `json:"platform,omitempty"`
	// ReleaseAt is the date when the image was released
	ReleaseAt *string `json:"release_at,omitempty"`
	// EndStandardSupportAt is the date when standard support ends
	EndStandardSupportAt *string `json:"end_standard_support_at,omitempty"`
	// EndLifeAt is the date when the image reaches end of life
	EndLifeAt *string `json:"end_life_at,omitempty"`
	// MinimumRequirements specifies the minimum hardware requirements
	MinimumRequirements MinimumRequirements `json:"minimum_requirements"`
	// Labels contains tags associated with the image
	Labels *[]string `json:"labels,omitempty"`
	// AvailabilityZones lists the availability zones where this image is available
	AvailabilityZones *[]string `json:"availability_zones,omitempty"`
}

// MinimumRequirements represents the minimum hardware requirements for an image.
// These requirements must be met by the instance type when creating instances from this image.
type MinimumRequirements struct {
	// VCPU is the minimum number of virtual CPUs required
	VCPU int `json:"vcpu"`
	// RAM is the minimum amount of RAM required in MB
	RAM int `json:"ram"`
	// Disk is the minimum disk size required in GB
	Disk int `json:"disk"`
}

// ImageStatus represents the current state of an image.
// The status indicates the lifecycle stage and availability of the image.
type ImageStatus string

const (
	// ImageStatusActive indicates the image is available for use
	ImageStatusActive ImageStatus = "active"
	// ImageStatusDeprecated indicates the image is deprecated but still available
	ImageStatusDeprecated ImageStatus = "deprecated"
	// ImageStatusDeleted indicates the image has been deleted
	ImageStatusDeleted ImageStatus = "deleted"
	// ImageStatusPending indicates the image is being processed
	ImageStatusPending ImageStatus = "pending"
	// ImageStatusCreating indicates the image is being created
	ImageStatusCreating ImageStatus = "creating"
	// ImageStatusImporting indicates the image is being imported
	ImageStatusImporting ImageStatus = "importing"
	// ImageStatusError indicates an error occurred during image processing
	ImageStatusError ImageStatus = "error"
	// ImageStatusDeletingError indicates an error occurred during image deletion
	ImageStatusDeletingError ImageStatus = "deleting_error"
)

// ImageService provides operations for managing virtual machine images.
// This interface allows listing available images with optional filtering.
type ImageService interface {
	// List returns a slice of images based on the provided listing options.
	//
	// Parameters:
	//   - ctx: Request context
	//   - opts: Options to filter and paginate the images
	//
	// Returns:
	//   - []Image: List of available images
	//   - error: Error if there's a failure in the request
	List(ctx context.Context, opts ImageListOptions) ([]Image, error)
}

// imageService implements the ImageService interface.
// This is an internal implementation that should not be used directly.
type imageService struct {
	client *VirtualMachineClient
}

// ImageListOptions defines the parameters for filtering and pagination of image lists.
// All fields are optional and allow controlling the listing behavior.
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

// List retrieves all images matching the provided options.
// This method makes an HTTP request to get the list of images
// and applies the filters specified in the options.
//
// Parameters:
//   - ctx: Request context
//   - opts: Options to filter and paginate the images
//
// Returns:
//   - []Image: List of available images
//   - error: Error if there's a failure in the request
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
