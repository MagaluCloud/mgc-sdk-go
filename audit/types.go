package audit

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type EventType struct {
	Type string `json:"type"`
}

type ListEventTypesParams struct {
	Limit    *int    `url:"_limit,omitempty"`
	Offset   *int    `url:"_offset,omitempty"`
	TenantID *string `url:"X-Tenant-ID,omitempty"`
}

type EventTypeService interface {
	List(ctx context.Context, params *ListEventTypesParams) ([]EventType, error)
}

type eventTypeService struct {
	client *AuditClient
}

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
