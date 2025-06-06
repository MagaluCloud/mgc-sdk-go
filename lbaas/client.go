package lbaas

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/load-balancer"
)

type LbaasClient struct {
	*client.CoreClient
}

func New(core *client.CoreClient) *LbaasClient {
	if core == nil {
		return nil
	}
	return &LbaasClient{CoreClient: core}
}

func (c *LbaasClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

func (c *LbaasClient) NetworkACLs() NetworkACLService {
	return &networkACLService{client: c}
}

func (c *LbaasClient) NetworkBackends() NetworkBackendService {
	return &networkBackendService{client: c}
}

func (c *LbaasClient) NetworkCertificates() NetworkCertificateService {
	return &networkCertificateService{client: c}
}

func (c *LbaasClient) NetworkHealthChecks() NetworkHealthCheckService {
	return &networkHealthCheckService{client: c}
}

func (c *LbaasClient) NetworkListeners() NetworkListenerService {
	return &networkListenerService{client: c}
}

func (c *LbaasClient) NetworkLoadBalancers() NetworkLoadBalancerService {
	return &networkLoadBalancerService{client: c}
}
