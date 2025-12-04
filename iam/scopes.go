package iam

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	scopesPath = "/scopes"
)

// ScopeService provides methods for managing scopes
type ScopeService interface {
	GroupsAndProductsAndScopes(ctx context.Context) ([]ScopeGroup, error)
}

// scopeService implements the ScopeService interface
type scopeService struct {
	client *IAMClient
}

// GetGroupsAndProductsAndScopes returns groups, products, and scopes
func (s *scopeService) GroupsAndProductsAndScopes(ctx context.Context) ([]ScopeGroup, error) {
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]ScopeGroup](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		scopesPath,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
