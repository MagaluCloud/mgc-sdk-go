package iam

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

const (
	serviceAccountsPath               = "/service-accounts"
	serviceAccountPathWithUUID        = "/service-accounts/%s"
	serviceAccountAPIKeysPath         = "/service-accounts/%s/api-keys"
	serviceAccountAPIKeyPathWithUUIDs = "/service-accounts/%s/api-keys/%s"
)

// ServiceAccountService provides methods for managing service accounts
type ServiceAccountService interface {
	List(ctx context.Context) ([]ServiceAccountDetail, error)
	Create(ctx context.Context, req ServiceAccountCreate) (*ServiceAccountDetail, error)
	Delete(ctx context.Context, saUUID string) error
	Edit(ctx context.Context, saUUID string, req ServiceAccountEdit) (*ServiceAccountDetail, error)
	APIKeys(ctx context.Context, saUUID string) ([]APIKeyServiceAccountDetail, error)
	CreateAPIKey(ctx context.Context, saUUID string, req APIKeyServiceAccountCreate) (*APIKeyServiceAccountDetail, error)
	RevokeAPIKey(ctx context.Context, saUUID string, apikeyUUID string) error
	EditAPIKey(ctx context.Context, saUUID string, apikeyUUID string, req APIKeyServiceAccountEditInput) (*APIKeyServiceAccountDetail, error)
}

// serviceAccountService implements the ServiceAccountService interface
type serviceAccountService struct {
	client *IAMClient
}

// List returns a list of service accounts
func (s *serviceAccountService) List(ctx context.Context) ([]ServiceAccountDetail, error) {
	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]ServiceAccountDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		serviceAccountsPath,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Create creates a new service account
func (s *serviceAccountService) Create(ctx context.Context, req ServiceAccountCreate) (*ServiceAccountDetail, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ServiceAccountDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		serviceAccountsPath,
		req,
		nil,
	)
}

// Delete removes a service account
func (s *serviceAccountService) Delete(ctx context.Context, saUUID string) error {
	if saUUID == "" {
		return &client.ValidationError{Field: "saUUID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(serviceAccountPathWithUUID, saUUID),
		nil,
		nil,
	)
}

// Edit updates a service account
func (s *serviceAccountService) Edit(ctx context.Context, saUUID string, req ServiceAccountEdit) (*ServiceAccountDetail, error) {
	if saUUID == "" {
		return nil, &client.ValidationError{Field: "saUUID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[ServiceAccountDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(serviceAccountPathWithUUID, saUUID),
		req,
		nil,
	)
}

// GetAPIKeys returns a list of API keys for a service account
func (s *serviceAccountService) APIKeys(ctx context.Context, saUUID string) ([]APIKeyServiceAccountDetail, error) {
	if saUUID == "" {
		return nil, &client.ValidationError{Field: "saUUID", Message: utils.CannotBeEmpty}
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]APIKeyServiceAccountDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(serviceAccountAPIKeysPath, saUUID),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateAPIKey creates a new API key for a service account
func (s *serviceAccountService) CreateAPIKey(ctx context.Context, saUUID string, req APIKeyServiceAccountCreate) (*APIKeyServiceAccountDetail, error) {
	if saUUID == "" {
		return nil, &client.ValidationError{Field: "saUUID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[APIKeyServiceAccountDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf(serviceAccountAPIKeysPath, saUUID),
		req,
		nil,
	)
}

// RevokeAPIKey revokes an API key for a service account
func (s *serviceAccountService) RevokeAPIKey(ctx context.Context, saUUID string, apikeyUUID string) error {
	if saUUID == "" {
		return &client.ValidationError{Field: "saUUID", Message: utils.CannotBeEmpty}
	}
	if apikeyUUID == "" {
		return &client.ValidationError{Field: "apikeyUUID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(serviceAccountAPIKeyPathWithUUIDs, saUUID, apikeyUUID),
		nil,
		nil,
	)
}

// EditAPIKey updates an API key for a service account
func (s *serviceAccountService) EditAPIKey(ctx context.Context, saUUID string, apikeyUUID string, req APIKeyServiceAccountEditInput) (*APIKeyServiceAccountDetail, error) {
	if saUUID == "" {
		return nil, &client.ValidationError{Field: "saUUID", Message: utils.CannotBeEmpty}
	}
	if apikeyUUID == "" {
		return nil, &client.ValidationError{Field: "apikeyUUID", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[APIKeyServiceAccountDetail](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(serviceAccountAPIKeyPathWithUUIDs, saUUID, apikeyUUID),
		req,
		nil,
	)
}
