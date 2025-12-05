package iam

import (
	"context"
	"net/http"
	"net/url"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	permissionsPath = "/permissions"
)

// PermissionService provides methods for managing IAM permissions
type PermissionService interface {
	ProductsAndPermissions(ctx context.Context, productName *string) ([]Product, error)
}

// permissionService implements the PermissionService interface
type permissionService struct {
	client *IAMClient
}

// ProductsAndPermissions returns a list of products and their permissions with optional product name filter
func (s *permissionService) ProductsAndPermissions(ctx context.Context, productName *string) ([]Product, error) {
	query := url.Values{}
	if productName != nil {
		query.Add("product_name", *productName)
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]Product](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		permissionsPath,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
