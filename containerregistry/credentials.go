package containerregistry

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	CredentialsService interface {
		Get(ctx context.Context) (*CredentialsResponse, error)
		ResetPassword(ctx context.Context) (*CredentialsResponse, error)
	}

	credentialsService struct {
		client *ContainerRegistryClient
	}

	CredentialsResponse struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
)

func (c *credentialsService) Get(ctx context.Context) (*CredentialsResponse, error) {
	path := "/v0/credentials"

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[CredentialsResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *credentialsService) ResetPassword(ctx context.Context) (*CredentialsResponse, error) {
	path := "/v0/credentials/password"

	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[CredentialsResponse](ctx, c.client.newRequest, c.client.GetConfig(), http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
