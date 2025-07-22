// Package dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
// This package allows you to manage database instances, clusters, replicas, engines, instance types, and parameters.
package dbaas

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/database"
)

// DBaaSClient represents a client for the Database as a Service
type DBaaSClient struct {
	*client.CoreClient
}

// ClientOption is a function type for configuring DBaaSClient options
type ClientOption func(*DBaaSClient)

// New creates a new DBaaSClient instance with the provided core client and options
func New(core *client.CoreClient, opts ...ClientOption) *DBaaSClient {
	if core == nil {
		return nil
	}

	client := &DBaaSClient{
		CoreClient: core,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// newRequest creates a new HTTP request for the database API
func (c *DBaaSClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Engines returns a service for managing database engines
func (c *DBaaSClient) Engines() EngineService {
	return &engineService{client: c}
}

// InstanceTypes returns a service for managing database instance types
func (c *DBaaSClient) InstanceTypes() InstanceTypeService {
	return &instanceTypeService{client: c}
}

// Instances returns a service for managing database instances
func (c *DBaaSClient) Instances() InstanceService {
	return &instanceService{client: c}
}

// Replicas returns a service for managing database replicas
func (c *DBaaSClient) Replicas() ReplicaService {
	return &replicaService{client: c}
}

// ParametersGroup returns a service for managing parameter groups
func (c *DBaaSClient) ParametersGroup() ParameterGroupService {
	return &parameterGroupService{client: c}
}

// Parameters returns a service for managing parameters within parameter groups
func (c *DBaaSClient) Parameters() ParameterService {
	return &parameterService{client: c}
}

// Clusters returns a service for managing database clusters
func (c *DBaaSClient) Clusters() ClusterService {
	return &clusterService{client: c}
}
