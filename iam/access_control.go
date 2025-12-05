package iam

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const (
	accessControlPath = "/access-control"
)

// AccessControlService provides methods for managing access control
type AccessControlService interface {
	Get(ctx context.Context) (*AccessControl, error)
	Create(ctx context.Context, req AccessControlCreate) (*AccessControl, error)
	Update(ctx context.Context, req AccessControlStatus) (*AccessControl, error)
}

// accessControlService implements the AccessControlService interface
type accessControlService struct {
	client *IAMClient
}

// Get retrieves the access control configuration
func (s *accessControlService) Get(ctx context.Context) (*AccessControl, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[AccessControl](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		accessControlPath,
		nil,
		nil,
	)
}

// Create creates a new access control configuration
func (s *accessControlService) Create(ctx context.Context, req AccessControlCreate) (*AccessControl, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[AccessControl](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		accessControlPath,
		req,
		nil,
	)
}

// Update updates the access control status
func (s *accessControlService) Update(ctx context.Context, req AccessControlStatus) (*AccessControl, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[AccessControl](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		accessControlPath,
		req,
		nil,
	)
}
