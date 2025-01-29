package kubernetes

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/kubernetes"
)

type KubernetesClient struct {
	*client.CoreClient
}

type ClientOption func(*KubernetesClient)

func New(core *client.CoreClient, opts ...ClientOption) *KubernetesClient {
	if core == nil {
		return nil
	}

	client := &KubernetesClient{
		CoreClient: core,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *KubernetesClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}
