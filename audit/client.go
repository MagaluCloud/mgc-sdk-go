package audit

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/audit"
)

type AuditClient struct {
	*client.CoreClient
}

func New(core *client.CoreClient) *AuditClient {
	if core == nil {
		return nil
	}
	return &AuditClient{
		CoreClient: core,
	}
}

func (c *AuditClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

func (c *AuditClient) Events() EventService {
	return &eventService{client: c}
}

func (c *AuditClient) EventTypes() EventTypeService {
	return &eventTypeService{client: c}
}
