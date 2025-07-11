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
		// Results contains the list of SSH key resources
		Results []SSHKey `json:"results"`
	}

	// SSHKey represents an SSH key resource
	SSHKey struct {
		// ID is the unique identifier of the SSH key
		ID string `json:"id"`
		// Name is the name of the SSH key
		Name string `json:"name"`
		// Key is the SSH public key content
		Key string `json:"key"`
		// KeyType is the type of the SSH key (e.g., ssh-rsa, ssh-ed25519)
		KeyType string `json:"key_type"`
	}

	// CreateSSHKeyRequest represents the parameters for creating a new SSH key
	CreateSSHKeyRequest struct {
		// Name is the name for the SSH key
		Name string `json:"name"`
		// Key is the SSH public key content
		Key string `json:"key"`
	}

	// ListOptions defines parameters for filtering and paginating SSH key lists
	ListOptions struct {
		// Limit specifies the maximum number of items to return
		Limit *int
		// Offset specifies the number of items to skip
		Offset *int
		// Sort defines the sort order for the results
		Sort *string
	}
)

// KeyService provides methods for managing SSH keys.
// All operations in this service are performed against the global endpoint,
// as SSH keys are not region-specific resources.
type KeyService interface {
	// List returns all SSH keys for the tenant.
	// Use ListOptions to control pagination and sorting of the results.
	List(ctx context.Context, opts ListOptions) ([]SSHKey, error)

	// Create registers a new SSH key globally.
	// The key will be available across all regions.
	Create(ctx context.Context, req CreateSSHKeyRequest) (*SSHKey, error)

	// Get retrieves a specific SSH key by ID.
	// Since keys are global, the same key ID will return the same key
	// regardless of the region.
	Get(ctx context.Context, keyID string) (*SSHKey, error)

	// Delete removes an SSH key globally.
	// Once deleted, the key will be unavailable across all regions.
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
