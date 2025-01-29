package network

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/network"
)

// NetworkClient represents a client for interacting with the network services
type NetworkClient struct {
	*client.CoreClient
}

// New creates a new NetworkClient instance
func New(core *client.CoreClient) *NetworkClient {
	if core == nil {
		return nil
	}
	client := &NetworkClient{
		CoreClient: core,
	}
	return client
}

// newRequest creates a new HTTP request with the network API base path
func (c *NetworkClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// VPCs returns a service for managing VPC resources
func (c *NetworkClient) VPCs() VPCService {
	return &vpcService{client: c}
}

// Subnets returns a service for managing subnet resources
func (c *NetworkClient) Subnets() SubnetService {
	return &subnetService{client: c}
}

// Ports returns a service for managing port resources
func (c *NetworkClient) Ports() PortService {
	return &portService{client: c}
}

// SecurityGroups returns a service for managing security group resources
func (c *NetworkClient) SecurityGroups() SecurityGroupService {
	return &securityGroupService{client: c}
}

// Rules returns a service for managing security group rule resources
func (c *NetworkClient) Rules() RuleService {
	return &ruleService{client: c}
}

// PublicIPs returns a service for managing public IP resources
func (c *NetworkClient) PublicIPs() PublicIPService {
	return &publicIPService{client: c}
}

// SubnetPools returns a service for managing subnet pool resources
func (c *NetworkClient) SubnetPools() SubnetPoolService {
	return &subnetPoolService{client: c}
}
