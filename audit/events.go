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

// ListEventsParams defines parameters for listing audit events.
// All fields are optional and allow filtering results in different ways.
type ListEventsParams struct {
	Limit       *int              `url:"_limit,omitempty"`
	Offset      *int              `url:"_offset,omitempty"`
	ID          *string           `url:"id,omitempty"`
	SourceLike  *string           `url:"source__like,omitempty"`
	Time        *time.Time        `url:"time,omitempty"`
	TypeLike    *string           `url:"type__like,omitempty"`
	ProductLike *string           `url:"product__like,omitempty"`
	AuthID      *string           `url:"authid,omitempty"`
	TenantID    *string           `url:"X-Tenant-ID,omitempty"`
	Data        map[string]string `url:"data,omitempty"`
}

// PaginatedMeta contains metadata about the paginated response.
// Provides information about pagination and result counting.
type PaginatedMeta struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

// PaginatedResponse represents a generic paginated response.
// Used to encapsulate paginated results of different types.
type PaginatedResponse[T any] struct {
	Meta    PaginatedMeta `json:"meta"`
	Results []T           `json:"results"`
}

// EventService defines the interface for audit event operations.
// This interface allows listing events with different filters and pagination options.
type EventService interface {
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
