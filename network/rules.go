package network

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	// RulesList represents a list of security group rules
	RulesList struct {
		// Rules contains the list of security group rule resources
		Rules []RuleResponse `json:"rules"`
	}

	// RuleResponse represents a security group rule resource
	RuleResponse struct {
		// ID is the unique identifier of the rule
		ID *string `json:"id,omitempty"`
		// ExternalID is the external identifier (optional)
		ExternalID *string `json:"external_id,omitempty"`
		// SecurityGroupID is the security group identifier (optional)
		SecurityGroupID *string `json:"security_group_id,omitempty"`
		// Direction is the traffic direction (ingress/egress) (optional)
		Direction *string `json:"direction,omitempty"`
		// PortRangeMin is the minimum port number (optional)
		PortRangeMin *int `json:"port_range_min,omitempty"`
		// PortRangeMax is the maximum port number (optional)
		PortRangeMax *int `json:"port_range_max,omitempty"`
		// Protocol is the network protocol (optional)
		Protocol *string `json:"protocol,omitempty"`
		// RemoteIPPrefix is the remote IP prefix (optional)
		RemoteIPPrefix *string `json:"remote_ip_prefix,omitempty"`
		// EtherType is the ethernet type (optional)
		EtherType *string `json:"ethertype"`
		// CreatedAt is the creation timestamp (optional)
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Status is the current status of the rule
		Status string `json:"status"`
		// Error contains error information if any (optional)
		Error *string `json:"error,omitempty"`
		// Description is the description of the rule (optional)
		Description *string `json:"description,omitempty"`
	}

	// RuleCreateRequest represents the parameters for creating a new security group rule
	RuleCreateRequest struct {
		// Direction is the traffic direction (ingress/egress) (optional)
		Direction *string `json:"direction,omitempty"`
		// PortRangeMin is the minimum port number (optional)
		PortRangeMin *int `json:"port_range_min,omitempty"`
		// PortRangeMax is the maximum port number (optional)
		PortRangeMax *int `json:"port_range_max,omitempty"`
		// Protocol is the network protocol (optional)
		Protocol *string `json:"protocol,omitempty"`
		// RemoteIPPrefix is the remote IP prefix (optional)
		RemoteIPPrefix *string `json:"remote_ip_prefix,omitempty"`
		// EtherType is the ethernet type
		EtherType string `json:"ethertype"`
		// Description is the description of the rule (optional)
		Description *string `json:"description,omitempty"`
	}

	// RuleCreateResponse represents the response after creating a security group rule
	RuleCreateResponse struct {
		// ID is the unique identifier of the created rule
		ID string `json:"id"`
	}
)

// RuleService provides operations for managing security group rules
type RuleService interface {
	// List retrieves all rules for a specific security group
	List(ctx context.Context, securityGroupID string) ([]RuleResponse, error)
	// Get retrieves details of a specific rule by its ID
	Get(ctx context.Context, id string) (*RuleResponse, error)
	// Create creates a new rule in a security group
	Create(ctx context.Context, securityGroupID string, req RuleCreateRequest) (string, error)
	// Delete removes a rule by its ID
	Delete(ctx context.Context, id string) error
}

// ruleService implements the RuleService interface
type ruleService struct {
	client *NetworkClient
}

// List retrieves all rules for a specific security group
func (s *ruleService) List(ctx context.Context, securityGroupID string) ([]RuleResponse, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[RulesList](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/security_groups/%s/rules", securityGroupID),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result.Rules, nil
}

// Get retrieves details of a specific rule by its ID
func (s *ruleService) Get(ctx context.Context, id string) (*RuleResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[RuleResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/rules/%s", id),
		nil,
		nil,
	)
}

// Create creates a new rule in a security group
func (s *ruleService) Create(ctx context.Context, securityGroupID string, req RuleCreateRequest) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[RuleCreateResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/security_groups/%s/rules", securityGroupID),
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Delete removes a rule by its ID
func (s *ruleService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/rules/%s", id),
		nil,
		nil,
	)
}
