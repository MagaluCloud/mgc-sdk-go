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

type service struct {
	client *Client
}

// BlockType represents the possible blocking states of an availability zone
type BlockType string

const (
	BlockTypeNone     BlockType = "none"
	BlockTypeTotal    BlockType = "total"
	BlockTypeReadOnly BlockType = "read-only"
)

// AvailabilityZone represents a single availability zone
type AvailabilityZone struct {
	ID        string    `json:"az_id"`
	BlockType BlockType `json:"block_type"`
}

// Region represents a region and its availability zones
type Region struct {
	ID                string             `json:"region_id"`
	AvailabilityZones []AvailabilityZone `json:"availability_zones"`
}

// ListResponse represents the response from listing availability zones
type ListResponse struct {
	Results []Region `json:"results"`
}

// List retrieves all availability zones
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
