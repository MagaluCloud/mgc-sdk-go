package sshkeys

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// ListSSHKeysResponse represents a list of SSH keys response
	ListSSHKeysResponse struct {
		Results []SSHKey `json:"results"`
	}

	// SSHKey represents an SSH key resource
	SSHKey struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Key     string `json:"key"`
		KeyType string `json:"key_type"`
	}

	// CreateSSHKeyRequest represents the parameters for creating a new SSH key
	CreateSSHKeyRequest struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}

	// ListOptions defines parameters for filtering and paginating SSH key lists
	ListOptions struct {
		Limit  *int
		Offset *int
		Sort   *string
	}
)

// KeyService provides methods for managing SSH keys.
// All operations in this service are performed against the global endpoint,
// as SSH keys are not region-specific resources.
type KeyService interface {
	List(ctx context.Context, opts ListOptions) ([]SSHKey, error)
	Create(ctx context.Context, req CreateSSHKeyRequest) (*SSHKey, error)
	Get(ctx context.Context, keyID string) (*SSHKey, error)
	Delete(ctx context.Context, keyID string) (*SSHKey, error)
}

// keyService implements the KeyService interface
type keyService struct {
	client *SSHKeyClient
}

// List returns all SSH keys for the tenant
func (s *keyService) List(ctx context.Context, opts ListOptions) ([]SSHKey, error) {
	query := make(url.Values)

	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListSSHKeysResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/ssh-keys",
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Create registers a new SSH key globally
func (s *keyService) Create(ctx context.Context, req CreateSSHKeyRequest) (*SSHKey, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SSHKey](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v0/ssh-keys",
		req,
		nil,
	)
}

// Get retrieves a specific SSH key by ID
func (s *keyService) Get(ctx context.Context, keyID string) (*SSHKey, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SSHKey](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/ssh-keys/%s", keyID),
		nil,
		nil,
	)
}

// Delete removes an SSH key globally
func (s *keyService) Delete(ctx context.Context, keyID string) (*SSHKey, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SSHKey](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/ssh-keys/%s", keyID),
		nil,
		nil,
	)
}
