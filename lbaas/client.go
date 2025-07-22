// Package lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
// This package allows you to manage network load balancers, listeners, backends, health checks, certificates, and ACLs.
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

// LbaasClient represents a client for the Load Balancer as a Service
type LbaasClient struct {
	*client.CoreClient
}

// New creates a new LbaasClient instance with the provided core client
func New(core *client.CoreClient) *LbaasClient {
	if core == nil {
		return nil
	}
	return &LbaasClient{CoreClient: core}
}

// newRequest creates a new HTTP request for the load balancer API
func (c *LbaasClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// NetworkACLs returns a service for managing network ACLs
func (c *LbaasClient) NetworkACLs() NetworkACLService {
	return &networkACLService{client: c}
}

// NetworkBackends returns a service for managing network backends
func (c *LbaasClient) NetworkBackends() NetworkBackendService {
	return &networkBackendService{client: c}
}

// NetworkCertificates returns a service for managing network certificates
func (c *LbaasClient) NetworkCertificates() NetworkCertificateService {
	return &networkCertificateService{client: c}
}

// NetworkHealthChecks returns a service for managing network health checks
func (c *LbaasClient) NetworkHealthChecks() NetworkHealthCheckService {
	return &networkHealthCheckService{client: c}
}

// NetworkListeners returns a service for managing network listeners
func (c *LbaasClient) NetworkListeners() NetworkListenerService {
	return &networkListenerService{client: c}
}

// NetworkLoadBalancers returns a service for managing network load balancers
func (c *LbaasClient) NetworkLoadBalancers() NetworkLoadBalancerService {
	return &networkLoadBalancerService{client: c}
}
