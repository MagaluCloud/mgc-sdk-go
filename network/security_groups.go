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
		// SecurityGroups contains the list of security group resources
		SecurityGroups []SecurityGroupResponse `json:"security_groups"`
	}

	// SecurityGroupResponse represents a security group resource
	SecurityGroupResponse struct {
		// ID is the unique identifier of the security group
		ID *string `json:"id,omitempty"`
		// VPCID is the VPC identifier (optional)
		VPCID *string `json:"vpc_id,omitempty"`
		// Name is the name of the security group (optional)
		Name *string `json:"name,omitempty"`
		// Description is the description of the security group (optional)
		Description *string `json:"description,omitempty"`
		// CreatedAt is the creation timestamp (optional)
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Updated is the last update timestamp (optional)
		Updated *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		// Status is the current status of the security group
		Status string `json:"status"`
		// Error contains error information if any (optional)
		Error *string `json:"error,omitempty"`
		// TenantID is the tenant identifier (optional)
		TenantID *string `json:"tenant_id,omitempty"`
		// ProjectType is the project type (optional)
		ProjectType *string `json:"project_type,omitempty"`
		// IsDefault indicates if this is the default security group (optional)
		IsDefault *bool `json:"is_default,omitempty"`
		// Ports contains the list of port IDs (optional)
		Ports *[]string `json:"ports,omitempty"`
	}

	// SecurityGroupDetailResponse represents detailed information about a security group
	SecurityGroupDetailResponse struct {
		SecurityGroupResponse
		// ExternalID is the external identifier (optional)
		ExternalID *string `json:"external_id,omitempty"`
		// Rules contains the security group rules (optional)
		Rules *[]RuleResponse `json:"rules"`
	}

	// SecurityGroupCreateRequest represents the parameters for creating a new security group
	SecurityGroupCreateRequest struct {
		// Name is the name of the security group
		Name string `json:"name"`
		// Description is the description of the security group (optional)
		Description *string `json:"description,omitempty"`
		// SkipDefaultRules indicates whether to skip default rules creation (optional)
		SkipDefaultRules *bool `json:"skip_default_rules,omitempty"`
	}

	// SecurityGroupCreateResponse represents the response after creating a security group
	SecurityGroupCreateResponse struct {
		// ID is the unique identifier of the created security group
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

// securityGroupService implements the SecurityGroupService interface
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
