# Audit

Package audit provides functionality to interact with the MagaluCloud audit
service. This package allows listing audit events and event types.
```
const (
DefaultBasePath = "/audit"
)
DefaultBasePath defines the default base path for audit APIs.



```
```
type AuditClient struct {
*client.CoreClient
}
AuditClient represents a client for the audit service. It encapsulates
functionality to access events and event types.

```
```
func New(core *client.CoreClient) *AuditClient
New creates a new instance of AuditClient. If the core client is nil,
returns nil.

```
```
func (c *AuditClient) EventTypes() EventTypeService
EventTypes returns a service to manage audit event types. This method allows
access to functionality such as listing event types.

```
```
func (c *AuditClient) Events() EventService
Events returns a service to manage audit events. This method allows access
to functionality such as listing events.

```
```
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
Event represents an audit event. Contains detailed information about an
action or operation performed in the system.

```
```
type EventService interface {
List(ctx context.Context, params *ListEventsParams) ([]Event, error)
}
EventService defines the interface for audit event operations. This
interface allows listing events with different filters and pagination
options.

```
```
type EventType struct {
Type string `json:"type"`
}
EventType represents an audit event type. Contains information about the
category or classification of an event.

```
```
type EventTypeService interface {
List(ctx context.Context, params *ListEventTypesParams) ([]EventType, error)
}
EventTypeService defines the interface for audit event type operations.
This interface allows listing available event types.

```
```
type ListEventTypesParams struct {
Limit    *int    `url:"_limit,omitempty"`
Offset   *int    `url:"_offset,omitempty"`
TenantID *string `url:"X-Tenant-ID,omitempty"`
}
ListEventTypesParams defines parameters for listing event types. All fields
are optional and allow filtering the results.

```
```
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
ListEventsParams defines parameters for listing audit events. All fields are
optional and allow filtering results in different ways.

```
```
type PaginatedMeta struct {
Limit  int `json:"limit,omitempty"`
Offset int `json:"offset,omitempty"`
Count  int `json:"count"`
Total  int `json:"total"`
}
PaginatedMeta contains metadata about the paginated response. Provides
information about pagination and result counting.

```
```
type PaginatedResponse[T any] struct {
Meta    PaginatedMeta `json:"meta"`
Results []T           `json:"results"`
}
PaginatedResponse represents a generic paginated response. Used to
encapsulate paginated results of different types.


```

