package network

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	// SecurityGroupListResponse represents a list of security groups response
	SecurityGroupListResponse struct {
		SecurityGroups []SecurityGroupResponse `json:"security_groups"`
	}

	// SecurityGroupResponse represents a security group resource
	SecurityGroupResponse struct {
		ID          *string                         `json:"id,omitempty"`
		VPCID       *string                         `json:"vpc_id,omitempty"`
		Name        *string                         `json:"name,omitempty"`
		Description *string                         `json:"description,omitempty"`
		CreatedAt   *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated     *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Status      string                          `json:"status"`
		Error       *string                         `json:"error,omitempty"`
		TenantID    *string                         `json:"tenant_id,omitempty"`
		ProjectType *string                         `json:"project_type,omitempty"`
		IsDefault   *bool                           `json:"is_default,omitempty"`
		Ports       *[]string                       `json:"ports,omitempty"`
	}

	// SecurityGroupDetailResponse represents detailed information about a security group
	SecurityGroupDetailResponse struct {
		SecurityGroupResponse
		ExternalID *string         `json:"external_id,omitempty"`
		Rules      *[]RuleResponse `json:"rules"`
	}

	// SecurityGroupCreateRequest represents the parameters for creating a new security group
	SecurityGroupCreateRequest struct {
		Name             string  `json:"name"`
		Description      *string `json:"description,omitempty"`
		SkipDefaultRules *bool   `json:"skip_default_rules,omitempty"`
	}

	// SecurityGroupCreateResponse represents the response after creating a security group
	SecurityGroupCreateResponse struct {
		ID string `json:"id"`
	}
)

// SecurityGroupService provides operations for managing security groups
type SecurityGroupService interface {
	// List retrieves all security groups for the current tenant
	List(ctx context.Context) ([]SecurityGroupResponse, error)

	// Get retrieves details of a specific security group by its ID
	Get(ctx context.Context, id string) (*SecurityGroupDetailResponse, error)

	// Create creates a new security group with the provided configuration
	Create(ctx context.Context, req SecurityGroupCreateRequest) (string, error)

	// Delete removes a security group by its ID
	Delete(ctx context.Context, id string) error
}

type securityGroupService struct {
	client *NetworkClient
}

// List retrieves all security groups for the current tenant
func (s *securityGroupService) List(ctx context.Context) ([]SecurityGroupResponse, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[SecurityGroupListResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/security_groups",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result.SecurityGroups, nil
}

// Get retrieves details of a specific security group by its ID
func (s *securityGroupService) Get(ctx context.Context, id string) (*SecurityGroupDetailResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[SecurityGroupDetailResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/security_groups/%s", id),
		nil,
		nil,
	)
}

// Create creates a new security group with the provided configuration
func (s *securityGroupService) Create(ctx context.Context, req SecurityGroupCreateRequest) (string, error) {
	queryParams := url.Values{}
	if req.SkipDefaultRules != nil {
		queryParams.Add("skip_default_rules", strconv.FormatBool(*req.SkipDefaultRules))
	}

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[SecurityGroupCreateResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v0/security_groups",
		req,
		queryParams,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Delete removes a security group by its ID
func (s *securityGroupService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/security_groups/%s", id),
		nil,
		nil,
	)
}
