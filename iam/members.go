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
	membersPath            = "/members"
	memberPathWithUUID     = "/members/%s"
	memberGrantsPath       = "/members/%s/grants"
	membersGrantsBatchPath = "/members/grants/batch"
)

// MemberService provides methods for managing IAM members
type MemberService interface {
	List(ctx context.Context, email *string) ([]Member, error)
	Create(ctx context.Context, req CreateMember) (*Member, error)
	Delete(ctx context.Context, uuid string) error
	Grants(ctx context.Context, uuid string) (*Privileges, error)
	AddGrants(ctx context.Context, uuid string, req EditGrant) error
	BatchUpdate(ctx context.Context, req BatchUpdateMembers) error
}

// memberService implements the MemberService interface
type memberService struct {
	client *IAMClient
}

// List returns a list of members with optional email filter
func (s *memberService) List(ctx context.Context, email *string) ([]Member, error) {
	query := url.Values{}
	if email != nil {
		query.Add("email", *email)
	}

	resp, err := mgc_http.ExecuteSimpleRequestWithRespBody[[]Member](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		membersPath,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Create creates a new member
func (s *memberService) Create(ctx context.Context, req CreateMember) (*Member, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[Member](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		membersPath,
		req,
		nil,
	)
}

// Delete removes a member
func (s *memberService) Delete(ctx context.Context, uuid string) error {
	if uuid == "" {
		return &client.ValidationError{Field: "uuid", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf(memberPathWithUUID, uuid),
		nil,
		nil,
	)
}

// GetGrants retrieves the grants (roles and permissions) for a member
func (s *memberService) Grants(ctx context.Context, uuid string) (*Privileges, error) {
	if uuid == "" {
		return nil, &client.ValidationError{Field: "uuid", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequestWithRespBody[Privileges](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf(memberGrantsPath, uuid),
		nil,
		nil,
	)
}

// AddGrants adds or removes grants (roles/permissions) for a member
func (s *memberService) AddGrants(ctx context.Context, uuid string, req EditGrant) error {
	if uuid == "" {
		return &client.ValidationError{Field: "uuid", Message: utils.CannotBeEmpty}
	}

	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf(memberGrantsPath, uuid),
		req,
		nil,
	)
}

// BatchUpdate updates multiple members in batch
func (s *memberService) BatchUpdate(ctx context.Context, req BatchUpdateMembers) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		membersGrantsBatchPath,
		req,
		nil,
	)
}
