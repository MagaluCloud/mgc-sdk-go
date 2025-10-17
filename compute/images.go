package compute

import (
	"context"
	"net/http"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ImageList represents the response from listing images.
// This structure encapsulates the API response format for images.
type ImageList struct {
	Meta   Meta    `json:"meta"`
	Images []Image `json:"images"`
}

// Image represents a virtual machine image.
// An image is a template that contains the operating system and software for creating instances.
type Image struct {
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

// MinimumRequirements represents the minimum hardware requirements for an image.
// These requirements must be met by the instance type when creating instances from this image.
type MinimumRequirements struct {
	VCPU int `json:"vcpu"`
	RAM  int `json:"ram"`
	Disk int `json:"disk"`
}

// ImageStatus represents the current state of an image.
// The status indicates the lifecycle stage and availability of the image.
type ImageStatus string

const (
	ImageStatusActive        ImageStatus = "active"
	ImageStatusDeprecated    ImageStatus = "deprecated"
	ImageStatusDeleted       ImageStatus = "deleted"
	ImageStatusPending       ImageStatus = "pending"
	ImageStatusCreating      ImageStatus = "creating"
	ImageStatusImporting     ImageStatus = "importing"
	ImageStatusError         ImageStatus = "error"
	ImageStatusDeletingError ImageStatus = "deleting_error"
)

// ImageService provides operations for managing virtual machine images.
// This interface allows listing available images with optional filtering.
type ImageService interface {
	List(ctx context.Context, opts ImageListOptions) (*ImageList, error)
	ListAll(ctx context.Context, opts ImageFilterOptions) ([]Image, error)
}

// imageService implements the ImageService interface.
// This is an internal implementation that should not be used directly.
type imageService struct {
	client *VirtualMachineClient
}

// ImageListOptions defines the parameters for filtering and pagination of image lists.
// All fields are optional and allow controlling the listing behavior.
type ImageListOptions struct {
	Limit            *int
	Offset           *int
	Sort             *string
	AvailabilityZone *string
}

// ImageFilterOptions defines filtering options for ListAll (without pagination)
type ImageFilterOptions struct {
	Sort             *string
	AvailabilityZone *string
}

// List retrieves images matching the provided options with pagination metadata.
// This method makes an HTTP request to get the list of images
// and applies the filters specified in the options.
func (s *imageService) List(ctx context.Context, opts ImageListOptions) (*ImageList, error) {
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
	if opts.AvailabilityZone != nil {
		q.Add("availability-zone", *opts.AvailabilityZone)
	}
	req.URL.RawQuery = q.Encode()

	response := &ImageList{}

	_, err = mgc_http.Do(s.client.GetConfig(), ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListAll retrieves all images across all pages with optional filtering.
// This method automatically handles pagination and returns all results.
func (s *imageService) ListAll(ctx context.Context, opts ImageFilterOptions) ([]Image, error) {
	var allImages []Image
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		listOpts := ImageListOptions{
			Offset:           &currentOffset,
			Limit:            &currentLimit,
			Sort:             opts.Sort,
			AvailabilityZone: opts.AvailabilityZone,
		}

		response, err := s.List(ctx, listOpts)
		if err != nil {
			return nil, err
		}

		allImages = append(allImages, response.Images...)

		// Check if we've retrieved all results
		if len(response.Images) < limit {
			break
		}

		offset += limit
	}

	return allImages, nil
}
