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
	ID          string                         `json:"id"`
	Source      string                         `json:"source"`
	Type        string                         `json:"type"`
	SpecVersion string                         `json:"specversion"`
	Subject     string                         `json:"subject"`
	Time        utils.LocalDateTimeWithoutZone `json:"time"`
	AuthID      string                         `json:"authid"`
	AuthType    string                         `json:"authtype"`
	Product     string                         `json:"product"`
	Region      *string                        `json:"region,omitempty"`
	TenantID    string                         `json:"tenantid"`
	Data        json.RawMessage                `json:"data"`
}

// EventFilterParams defines filtering parameters for ListAll (without pagination).
type EventFilterParams struct {
	ID          *string           `json:"id,omitempty"`
	SourceLike  *string           `json:"source__like,omitempty"`
	Time        *time.Time        `json:"time,omitempty"`
	TypeLike    *string           `json:"type__like,omitempty"`
	ProductLike *string           `json:"product__like,omitempty"`
	AuthID      *string           `json:"authid,omitempty"`
	TenantID    *string           `json:"X-Tenant-ID,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
}

// ListEventsParams defines parameters for listing audit events.
// It extends EventFilterParams by adding pagination fields.
type ListEventsParams struct {
	EventFilterParams
	Limit  *int `json:"_limit,omitempty"`
	Offset *int `json:"_offset,omitempty"`
}

// EventService defines the interface for audit event operations.
// This interface allows listing events with different filters and pagination options.
type EventService interface {
	List(ctx context.Context, params *ListEventsParams) (*PaginatedResponse[Event], error)
	ListAll(ctx context.Context, params *EventFilterParams) ([]Event, error)
}

// eventService implements the EventService interface.
// This is an internal implementation that should not be used directly.
type eventService struct {
	client *AuditClient
}

// List retrieves audit events with pagination metadata.
// This method makes an HTTP request to get the list of audit events
// and applies the filters specified in the parameters.
func (s *eventService) List(ctx context.Context, params *ListEventsParams) (*PaginatedResponse[Event], error) {
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
	return result, nil
}

// ListAll retrieves all audit events across all pages with optional filtering.
// This method automatically handles pagination and returns all results.
func (s *eventService) ListAll(ctx context.Context, params *EventFilterParams) ([]Event, error) {
	var allEvents []Event
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		listParams := &ListEventsParams{
			Offset: &currentOffset,
			Limit:  &currentLimit,
		}

		if params != nil {
			listParams.EventFilterParams = *params
		}

		response, err := s.List(ctx, listParams)
		if err != nil {
			return nil, err
		}

		allEvents = append(allEvents, response.Results...)

		// Check if we've retrieved all results
		if len(response.Results) < limit {
			break
		}

		offset += limit
	}

	return allEvents, nil
}
