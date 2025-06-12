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

type DBaaSClient struct {
	*client.CoreClient
}

type ClientOption func(*DBaaSClient)

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

func (c *DBaaSClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

func (c *DBaaSClient) Engines() EngineService {
	return &engineService{client: c}
}

func (c *DBaaSClient) InstanceTypes() InstanceTypeService {
	return &instanceTypeService{client: c}
}

func (c *DBaaSClient) Instances() InstanceService {
	return &instanceService{client: c}
}

func (c *DBaaSClient) Replicas() ReplicaService {
	return &replicaService{client: c}
}

func (c *DBaaSClient) ParametersGroup() ParameterGroupService {
	return &parameterGroupService{client: c}
}

func (c *DBaaSClient) Parameters() ParameterService {
	return &parameterService{client: c}
}

func (c *DBaaSClient) Clusters() ClusterService {
	return &clusterService{client: c}
}
