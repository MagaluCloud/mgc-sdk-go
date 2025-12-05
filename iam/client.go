// Package iam provides a client for interacting with the Magalu Cloud IAM API.
// This package allows you to manage members, roles, permissions, service accounts, and access control.
package iam

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	DefaultBasePath = "/iam/api/v1"
)

// IAMClient represents a client for the IAM service
type IAMClient struct {
	*client.CoreClient
}

// ClientOption is a function type for configuring IAMClient options
type ClientOption func(*IAMClient)

func WithGlobalBasePath(basePath client.MgcUrl) ClientOption {
	return func(c *IAMClient) {
		c.GetConfig().BaseURL = basePath
	}
}

// New creates a new IAMClient instance with the provided core client and options
func New(core *client.CoreClient, opts ...ClientOption) *IAMClient {
	if core == nil {
		return nil
	}

	iamClient := &IAMClient{
		CoreClient: core,
	}

	iamClient.GetConfig().BaseURL = client.Global

	for _, opt := range opts {
		opt(iamClient)
	}

	return iamClient
}

// newRequest creates a new HTTP request for the IAM API
func (c *IAMClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	return mgc_http.NewRequest(c.GetConfig(), ctx, method, DefaultBasePath+path, &body)
}

// Members returns a service for managing IAM members
func (c *IAMClient) Members() MemberService {
	return &memberService{client: c}
}

// Roles returns a service for managing IAM roles
func (c *IAMClient) Roles() RoleService {
	return &roleService{client: c}
}

// Permissions returns a service for managing IAM permissions
func (c *IAMClient) Permissions() PermissionService {
	return &permissionService{client: c}
}

// AccessControl returns a service for managing access control
func (c *IAMClient) AccessControl() AccessControlService {
	return &accessControlService{client: c}
}

// ServiceAccounts returns a service for managing service accounts
func (c *IAMClient) ServiceAccounts() ServiceAccountService {
	return &serviceAccountService{client: c}
}

// Scopes returns a service for managing scopes
func (c *IAMClient) Scopes() ScopeService {
	return &scopeService{client: c}
}
