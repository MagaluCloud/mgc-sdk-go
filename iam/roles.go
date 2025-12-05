package iam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

const (
	rolesPath           = "/roles"
	rolePathWithName    = "/roles/%s"
	rolePermissionsPath = "/roles/%s/permissions"
	roleMembersPath     = "/roles/%s/members"
)

// RoleService provides methods for managing IAM roles
type RoleService interface {
	List(ctx context.Context, roleName *string) ([]Role, error)
	Create(ctx context.Context, req CreateRole) ([]Role, error)
	Delete(ctx context.Context, roleName string) error
	Permissions(ctx context.Context, roleName string) (*RolePermissions, error)
	EditPermissions(ctx context.Context, roleName string, req EditPermissions) ([]Role, error)
	Members(ctx context.Context, roleName string) ([]RolesMember, error)
}

// roleService implements the RoleService interface
type roleService struct {
	client *IAMClient
}

// List returns a list of roles with optional role name filter
func (s *roleService) List(ctx context.Context, roleName *string) ([]Role, error) {
	query := url.Values{}
	if roleName != nil {
		query.Add("role_name", *roleName)
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]Role](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		rolesPath,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Create creates a new role
func (s *roleService) Create(ctx context.Context, req CreateRole) ([]Role, error) {
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]Role](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		rolesPath,
		req,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Delete removes a role
func (s *roleService) Delete(ctx context.Context, roleName string) error {
	if roleName == "" {
		return &client.ValidationError{Field: "roleName", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(rolePathWithName, roleName),
		nil,
		nil,
	)
}

// GetPermissions retrieves the permissions for a role
func (s *roleService) Permissions(ctx context.Context, roleName string) (*RolePermissions, error) {
	if roleName == "" {
		return nil, &client.ValidationError{Field: "roleName", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[RolePermissions](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(rolePermissionsPath, roleName),
		nil,
		nil,
	)
}

// EditPermissions adds or removes permissions for a role
func (s *roleService) EditPermissions(ctx context.Context, roleName string, req EditPermissions) ([]Role, error) {
	if roleName == "" {
		return nil, &client.ValidationError{Field: "roleName", Message: utils.CannotBeEmpty}
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]Role](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(rolePermissionsPath, roleName),
		req,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetMembers retrieves the members that have a specific role
func (s *roleService) Members(ctx context.Context, roleName string) ([]RolesMember, error) {
	if roleName == "" {
		return nil, &client.ValidationError{Field: "roleName", Message: utils.CannotBeEmpty}
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]RolesMember](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(roleMembersPath, roleName),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
