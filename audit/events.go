package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

// Event represents an audit event.
// Contains detailed information about an action or operation performed in the system.
type Event struct {
	// ID is the unique identifier of the event
	ID string `json:"id"`
	// Source identifies the origin or source of the event
	Source string `json:"source"`
	// Type defines the type/category of the event
	Type string `json:"type"`
	// SpecVersion specifies the version of the event specification
	SpecVersion string `json:"specversion"`
	// Subject describes the subject or object of the event
	Subject string `json:"subject"`
	// Time records when the event occurred
	Time utils.LocalDateTimeWithoutZone `json:"time"`
	// AuthID identifies the user or entity that performed the action
	AuthID string `json:"authid"`
	// AuthType specifies the type of authentication used
	AuthType string `json:"authtype"`
	// Product identifies the product or service related to the event
	Product string `json:"product"`
	// Region specifies the region where the event occurred (optional)
	Region *string `json:"region,omitempty"`
	// TenantID identifies the tenant associated with the event
	TenantID string `json:"tenantid"`
	// Data contains additional event-specific data in JSON format
	Data json.RawMessage `json:"data"`
}

// ListEventsParams defines parameters for listing audit events.
// All fields are optional and allow filtering results in different ways.
type ListEventsParams struct {
	// Limit defines the maximum number of results to be returned
	Limit *int `url:"_limit,omitempty"`
	// Offset defines the number of results to be skipped (for pagination)
	Offset *int `url:"_offset,omitempty"`
	// ID filters events by specific ID
	ID *string `url:"id,omitempty"`
	// SourceLike filters events by source using similarity search
	SourceLike *string `url:"source__like,omitempty"`
	// Time filters events by specific date/time
	Time *time.Time `url:"time,omitempty"`
	// TypeLike filters events by type using similarity search
	TypeLike *string `url:"type__like,omitempty"`
	// ProductLike filters events by product using similarity search
	ProductLike *string `url:"product__like,omitempty"`
	// AuthID filters events by authentication ID
	AuthID *string `url:"authid,omitempty"`
	// TenantID filters events by tenant ID
	TenantID *string `url:"X-Tenant-ID,omitempty"`
	// Data allows filtering by specific fields within the event data
	Data map[string]string `url:"data,omitempty"`
}

// PaginatedMeta contains metadata about the paginated response.
// Provides information about pagination and result counting.
type PaginatedMeta struct {
	// Limit is the maximum number of results per page
	Limit int `json:"limit,omitempty"`
	// Offset is the number of results skipped
	Offset int `json:"offset,omitempty"`
	// Count is the number of results in the current page
	Count int `json:"count"`
	// Total is the total number of available results
	Total int `json:"total"`
}

// PaginatedResponse represents a generic paginated response.
// Used to encapsulate paginated results of different types.
type PaginatedResponse[T any] struct {
	// Meta contains pagination information
	Meta PaginatedMeta `json:"meta"`
	// Results contains the list of results for the current page
	Results []T `json:"results"`
}

// EventService defines the interface for audit event operations.
// This interface allows listing events with different filters and pagination options.
type EventService interface {
	// List returns a list of audit events.
	// Results can be filtered using the provided parameters.
	//
	// Parameters:
	//   - ctx: Request context
	//   - params: Optional parameters to filter and paginate results
	//
	// Returns:
	//   - []Event: List of audit events
	//   - error: Error if there's a failure in the request
	List(ctx context.Context, params *ListEventsParams) ([]Event, error)
}

// eventService implements the EventService interface.
// This is an internal implementation that should not be used directly.
type eventService struct {
	client *AuditClient
}

// List implements the List method of the EventService interface.
// This method makes an HTTP request to get the list of audit events
// and applies the filters specified in the parameters.
//
// Parameters:
//   - ctx: Request context
//   - params: Optional parameters to filter and paginate results
//
// Returns:
//   - []Event: List of audit events
//   - error: Error if there's a failure in the request
func (s *eventService) List(ctx context.Context, params *ListEventsParams) ([]Event, error) {
	query := make(url.Values)

	if params != nil {
		if params.Limit != nil {
			query.Set("_limit", strconv.Itoa(*params.Limit))
		}
		if params.Offset != nil {
			query.Set("_offset", strconv.Itoa(*params.Offset))
		}
		if params.ID != nil {
			query.Set("id", *params.ID)
		}
		if params.SourceLike != nil {
			query.Set("source__like", *params.SourceLike)
		}
		if params.TypeLike != nil {
			query.Set("type__like", *params.TypeLike)
		}
		if params.ProductLike != nil {
			query.Set("product__like", *params.ProductLike)
		}
		if params.AuthID != nil {
			query.Set("authid", *params.AuthID)
		}
		if params.TenantID != nil {
			query.Set("X-Tenant-ID", *params.TenantID)
		}
		for k, v := range params.Data {
			query.Set("data."+k, v)
		}
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[PaginatedResponse[Event]](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/events",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}
