package availabilityzones

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// ListOptions contains the options for listing availability zones
type ListOptions struct {
	ShowBlocked bool
}

// Service defines the interface for availability zone operations
type Service interface {
	List(ctx context.Context, opts ListOptions) ([]Region, error)
}

// service implements the Service interface.
// This is an internal implementation that should not be used directly.
type service struct {
	client *Client
}

// BlockType represents the possible blocking states of an availability zone.
// An availability zone can be in different states that affect its usability.
type BlockType string

const (
	BlockTypeNone     BlockType = "none"
	BlockTypeTotal    BlockType = "total"
	BlockTypeReadOnly BlockType = "read-only"
)

// AvailabilityZone represents a single availability zone within a region.
// Each availability zone has a unique identifier and can have different blocking states.
type AvailabilityZone struct {
	ID        string    `json:"az_id"`
	BlockType BlockType `json:"block_type"`
}

// Region represents a region and its associated availability zones.
// A region contains multiple availability zones that can be used for resource deployment.
type Region struct {
	ID                string             `json:"region_id"`
	AvailabilityZones []AvailabilityZone `json:"availability_zones"`
}

// ListResponse represents the response from listing availability zones.
// This structure encapsulates the API response format.
type ListResponse struct {
	Results []Region `json:"results"`
}

// List retrieves all availability zones across all regions.
// This method makes an HTTP request to get the list of availability zones
// and applies the filters specified in the options.
func (s *service) List(ctx context.Context, opts ListOptions) ([]Region, error) {
	query := url.Values{}
	query.Set("show_is_blocked", strconv.FormatBool(opts.ShowBlocked))

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/availability-zones",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}
