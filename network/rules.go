package network

import (
	"context"
	"fmt"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

type (
	RulesList struct {
		Rules []RuleResponse `json:"rules"`
	}

	RuleResponse struct {
		ID              *string                         `json:"id,omitempty"`
		ExternalID      *string                         `json:"external_id,omitempty"`
		SecurityGroupID *string                         `json:"security_group_id,omitempty"`
		Direction       *string                         `json:"direction,omitempty"`
		PortRangeMin    *int                            `json:"port_range_min,omitempty"`
		PortRangeMax    *int                            `json:"port_range_max,omitempty"`
		Protocol        *string                         `json:"protocol,omitempty"`
		RemoteIPPrefix  *string                         `json:"remote_ip_prefix,omitempty"`
		RemoteGroupID   *string                         `json:"remote_group_id,omitempty"`
		EtherType       *string                         `json:"ethertype"`
		CreatedAt       *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Status          string                          `json:"status"`
		Error           *string                         `json:"error,omitempty"`
		Description     *string                         `json:"description,omitempty"`
	}

	RuleCreateRequest struct {
		Direction      *string `json:"direction,omitempty"`
		PortRangeMin   *int    `json:"port_range_min,omitempty"`
		PortRangeMax   *int    `json:"port_range_max,omitempty"`
		Protocol       *string `json:"protocol,omitempty"`
		RemoteIPPrefix *string `json:"remote_ip_prefix,omitempty"`
		RemoteGroupID  *string `json:"remote_group_id,omitempty"`
		EtherType      string  `json:"ethertype"`
		Description    *string `json:"description,omitempty"`
	}

	RuleCreateResponse struct {
		ID string `json:"id"`
	}
)

type RuleService interface {
	List(ctx context.Context, securityGroupID string) ([]RuleResponse, error)
	Get(ctx context.Context, id string) (*RuleResponse, error)
	Create(ctx context.Context, securityGroupID string, req RuleCreateRequest) (string, error)
	Delete(ctx context.Context, id string) error
}

type ruleService struct {
	client *NetworkClient
}

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
