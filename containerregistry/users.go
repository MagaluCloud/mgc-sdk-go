package containerregistry

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// UsersService manages the container registry user lifecycle.
	UsersService interface {
		Create(ctx context.Context) (*UserResponse, error)
		Get(ctx context.Context) (*UserResponse, error)
		Delete(ctx context.Context, userID string) error
	}

	// UserResponse represents a container registry user.
	UserResponse struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
	}

	usersService struct {
		client *ContainerRegistryClient
	}
)

// Create provisions the authenticated caller as a container registry user.
func (s *usersService) Create(ctx context.Context) (*UserResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[UserResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodPost, "/v0/users", nil, nil,
	)
}

// Get returns the authenticated container registry user.
func (s *usersService) Get(ctx context.Context) (*UserResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[UserResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodGet, "/v0/users", nil, nil,
	)
}

// Delete removes a container registry user by ID.
func (s *usersService) Delete(ctx context.Context, userID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodDelete, fmt.Sprintf("/v0/users/%s", userID), nil, nil,
	)
}
