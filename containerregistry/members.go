package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// MembersService manages the membership of users in a container registry.
	MembersService interface {
		Add(ctx context.Context, registryID string, request MemberRequest) (*MemberResponse, error)
		List(ctx context.Context, registryID string, opts MemberListOptions) (*ListMembersResponse, error)
		ListAll(ctx context.Context, registryID string, filterOpts MemberFilterOptions) ([]MemberResponse, error)
		Get(ctx context.Context, registryID, memberID string) (*MemberResponse, error)
		Update(ctx context.Context, registryID, memberID string, request MemberUpdateRequest) (*MemberResponse, error)
		Delete(ctx context.Context, registryID, memberID string) error
	}

	// MemberRequest is the payload to add a user as a registry member.
	MemberRequest struct {
		UserID string  `json:"user_id"`
		Role   *string `json:"role,omitempty"`
	}

	// MemberUpdateRequest is the payload to change a member's role.
	MemberUpdateRequest struct {
		Role string `json:"role"`
	}

	// MemberResponse represents a member relationship inside a registry.
	MemberResponse struct {
		ID         string `json:"id"`
		RegistryID string `json:"registry_id"`
		UserID     string `json:"user_id"`
		Role       string `json:"role"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	// MemberListOptions provides pagination and filter options for listing members.
	MemberListOptions struct {
		Offset *int
		Limit  *int
		MemberFilterOptions
	}

	// MemberFilterOptions provides filter options that survive across paginated calls.
	MemberFilterOptions struct {
		Sort *string
	}

	// ListMembersResponse is the paginated response of members.
	ListMembersResponse = helpers.PaginatedResponse[MemberResponse]

	membersService struct {
		client *ContainerRegistryClient
	}
)

// Add registers an existing user as a member of the given registry.
func (s *membersService) Add(ctx context.Context, registryID string, request MemberRequest) (*MemberResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[MemberResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodPost, fmt.Sprintf("/v0/registries/%s/members", registryID), request, nil,
	)
}

// List returns a page of members of the given registry.
func (s *membersService) List(ctx context.Context, registryID string, opts MemberListOptions) (*ListMembersResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ListMembersResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodGet, fmt.Sprintf("/v0/registries/%s/members", registryID), nil,
		memberListQuery(opts),
	)
}

// ListAll walks every page until the registry's full membership is retrieved.
func (s *membersService) ListAll(ctx context.Context, registryID string, filterOpts MemberFilterOptions) ([]MemberResponse, error) {
	var all []MemberResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		page, err := s.List(ctx, registryID, MemberListOptions{
			Offset:              &currentOffset,
			Limit:               &currentLimit,
			MemberFilterOptions: filterOpts,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, page.Results...)
		if len(page.Results) < limit {
			break
		}
		offset += limit
	}
	return all, nil
}

// Get returns a specific member relationship.
func (s *membersService) Get(ctx context.Context, registryID, memberID string) (*MemberResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[MemberResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodGet, fmt.Sprintf("/v0/registries/%s/members/%s", registryID, memberID), nil, nil,
	)
}

// Update changes the role assigned to a member.
func (s *membersService) Update(ctx context.Context, registryID, memberID string, request MemberUpdateRequest) (*MemberResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[MemberResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodPatch, fmt.Sprintf("/v0/registries/%s/members/%s", registryID, memberID), request, nil,
	)
}

// Delete removes a member from a registry.
func (s *membersService) Delete(ctx context.Context, registryID, memberID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodDelete, fmt.Sprintf("/v0/registries/%s/members/%s", registryID, memberID), nil, nil,
	)
}

func memberListQuery(opts MemberListOptions) url.Values {
	q := make(url.Values)
	if opts.Limit != nil {
		q.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Set("_sort", *opts.Sort)
	}
	return q
}
