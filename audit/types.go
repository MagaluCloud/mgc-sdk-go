package audit

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// EventType represents an audit event type.
// Contains information about the category or classification of an event.
type EventType struct {
	Type string `json:"type"`
}

// ListEventTypesParams defines parameters for listing event types.
// All fields are optional and allow filtering the results.
type ListEventTypesParams struct {
	Limit    *int    `url:"_limit,omitempty"`
	Offset   *int    `url:"_offset,omitempty"`
	TenantID *string `url:"X-Tenant-ID,omitempty"`
}

// EventTypeService defines the interface for audit event type operations.
// This interface allows listing available event types.
type EventTypeService interface {
	List(ctx context.Context, params *ListEventTypesParams) ([]EventType, error)
}

// eventTypeService implements the EventTypeService interface.
// This is an internal implementation that should not be used directly.
type eventTypeService struct {
	client *AuditClient
}

// List implements the List method of the EventTypeService interface.
// This method makes an HTTP request to get the list of event types
// and applies the filters specified in the parameters.
func (s *eventTypeService) List(ctx context.Context, params *ListEventTypesParams) ([]EventType, error) {
	query := make(url.Values)

	if params != nil {
		if params.Limit != nil {
			query.Set("_limit", strconv.Itoa(*params.Limit))
		}
		if params.Offset != nil {
			query.Set("_offset", strconv.Itoa(*params.Offset))
		}
		if params.TenantID != nil {
			query.Set("X-Tenant-ID", *params.TenantID)
		}
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[PaginatedResponse[EventType]](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/event-types",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}
