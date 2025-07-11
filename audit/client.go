// Package audit provides functionality to interact with the MagaluCloud audit service.
// This package allows listing audit events and event types.
package audit

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// DefaultBasePath defines the default base path for audit APIs.
const (
	DefaultBasePath = "/audit"
)

// AuditClient represents a client for the audit service.
// It encapsulates functionality to access events and event types.
type AuditClient struct {
	*client.CoreClient
}

// New creates a new instance of AuditClient.
// If the core client is nil, returns nil.
//
// Parameters:
//   - core: The core client that will be used for HTTP requests
//
// Returns:
//   - *AuditClient: A new instance of the audit client or nil if core is nil
func New(core *client.CoreClient) *AuditClient {
	if core == nil {
		return nil
	}
	return &AuditClient{
		CoreClient: core,
	}
}

// newRequest creates a new HTTP request for the audit service.
// This method is internal and should not be called directly by SDK users.
//
// Parameters:
//   - ctx: Request context
//   - method: HTTP method (GET, POST, etc.)
//   - path: API path (will be concatenated with DefaultBasePath)
//   - body: Request body (can be nil)
//
// Returns:
//   - *http.Request: The created HTTP request
//   - error: Error if there's a failure creating the request
func (c *AuditClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Events returns a service to manage audit events.
// This method allows access to functionality such as listing events.
//
// Returns:
//   - EventService: Interface for audit event operations
func (c *AuditClient) Events() EventService {
	return &eventService{client: c}
}

// EventTypes returns a service to manage audit event types.
// This method allows access to functionality such as listing event types.
//
// Returns:
//   - EventTypeService: Interface for audit event type operations
func (c *AuditClient) EventTypes() EventTypeService {
	return &eventTypeService{client: c}
}
