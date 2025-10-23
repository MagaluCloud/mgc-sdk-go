package audit

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// EventType represents an audit event type.
// Contains information about the category or classification of an event.
type EventType struct {
	Type string `json:"type"`
}

// EventTypeFilterParams defines filtering parameters for ListAll (without pagination).
type EventTypeFilterParams struct {
	TenantID *string `json:"X-Tenant-ID,omitempty"`
}

// ListEventTypesParams defines parameters for listing event types.
// It extends EventTypeFilterParams by adding pagination fields.
type ListEventTypesParams struct {
	EventTypeFilterParams
	Limit  *int `json:"_limit,omitempty"`
	Offset *int `json:"_offset,omitempty"`
}

// PaginatedMeta contains metadata about the paginated response.
// Provides information about pagination and result counting.
type PaginatedMeta = helpers.AuditPaginatedMeta

// PaginatedResponse represents a generic paginated response.
// Used to encapsulate paginated results of different types.
type PaginatedResponse[T any] = helpers.AuditPaginatedResponse[T]

// EventTypeService defines the interface for audit event type operations.
// This interface allows listing available event types.
type EventTypeService interface {
	List(ctx context.Context, params *ListEventTypesParams) (*PaginatedResponse[EventType], error)
	ListAll(ctx context.Context, params *EventTypeFilterParams) ([]EventType, error)
}

// eventTypeService implements the EventTypeService interface.
// This is an internal implementation that should not be used directly.
type eventTypeService struct {
	client *AuditClient
}

// List retrieves event types with pagination metadata.
// This method makes an HTTP request to get the list of event types
// and applies the filters specified in the parameters.
func (s *eventTypeService) List(ctx context.Context, params *ListEventTypesParams) (*PaginatedResponse[EventType], error) {
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
	return result, nil
}

// ListAll retrieves all event types across all pages with optional filtering.
// This method automatically handles pagination and returns all results.
func (s *eventTypeService) ListAll(ctx context.Context, params *EventTypeFilterParams) ([]EventType, error) {
	var allEventTypes []EventType
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		listParams := &ListEventTypesParams{
			Offset: &currentOffset,
			Limit:  &currentLimit,
		}

		if params != nil {
			listParams.EventTypeFilterParams = *params
		}

		response, err := s.List(ctx, listParams)
		if err != nil {
			return nil, err
		}

		allEventTypes = append(allEventTypes, response.Results...)

		// Check if we've retrieved all results
		if len(response.Results) < limit {
			break
		}

		offset += limit
	}

	return allEventTypes, nil
}
