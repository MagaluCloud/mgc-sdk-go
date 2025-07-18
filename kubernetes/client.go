// Package kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
// This package allows you to manage Kubernetes clusters, node pools, flavors, and versions.
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

// KubernetesClient represents a client for the Kubernetes service
type KubernetesClient struct {
	*client.CoreClient
}

// ClientOption is a function type for configuring KubernetesClient options
type ClientOption func(*KubernetesClient)

// New creates a new KubernetesClient instance with the provided core client and options
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

// newRequest creates a new HTTP request for the Kubernetes API
func (c *KubernetesClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Clusters returns a service for managing Kubernetes clusters
func (c *KubernetesClient) Clusters() ClusterService {
	return &clusterService{client: c}
}

// Flavors returns a service for managing Kubernetes flavors
func (c *KubernetesClient) Flavors() FlavorService {
	return &flavorService{client: c}
}

// Nodepools returns a service for managing Kubernetes node pools
func (c *KubernetesClient) Nodepools() NodePoolService {
	return &nodePoolService{client: c}
}

// Versions returns a service for managing Kubernetes versions
func (c *KubernetesClient) Versions() VersionService {
	return &versionService{client: c}
}
