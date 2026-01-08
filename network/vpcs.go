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

const (
	SecurityGroupsExpand   = "security_groups"
	SubnetsExpand          = "subnets"
	PortStatusProvisioning = "provisioning"
	PortStatusActive       = "active"
	PortStatusError        = "error"
	PublicIPStatusCreated  = "created"
	PublicIPStatusPending  = "pending"
	PublicIPStatusError    = "error"
)

type (
	// ListVPCsResponse represents a list of VPCs response
	ListVPCsResponse struct {
		VPCs []VPC `json:"vpcs"`
	}

	// VPC represents a Virtual Private Cloud resource
	VPC struct {
		ID              *string                         `json:"id,omitempty"`
		TenantID        *string                         `json:"tenant_id,omitempty"`
		Name            *string                         `json:"name,omitempty"`
		Description     *string                         `json:"description,omitempty"`
		Status          string                          `json:"status"`
		RouterID        *string                         `json:"router_id,omitempty"`
		ExternalNetwork *string                         `json:"external_network,omitempty"`
		NetworkID       *string                         `json:"network_id,omitempty"`
		Subnets         *[]string                       `json:"subnets,omitempty"`
		SecurityGroups  *[]string                       `json:"security_groups,omitempty"`
		CreatedAt       *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated         *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		IsDefault       *bool                           `json:"is_default,omitempty"`
	}

	// CreateVPCRequest represents the parameters for creating a new VPC
	CreateVPCRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
	}

	// RenameVPCRequest represents the parameters for renaming a VPC
	RenameVPCRequest struct {
		Name string `json:"name"`
	}

	// ListOptions represents parameters for filtering and pagination
	ListOptions struct {
		Limit  *int
		Offset *int
		Sort   *string
	}

	// CreateVPCResponse represents the response after creating a VPC
	CreateVPCResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	// PortCreateRequest represents the parameters for creating a port
	PortCreateRequest struct {
		Name           string    `json:"name"`
		HasPIP         *bool     `json:"has_pip,omitempty"`
		HasSG          *bool     `json:"has_sg,omitempty"`
		Subnets        *[]string `json:"subnets,omitempty"`
		SecurityGroups *[]string `json:"security_groups_id,omitempty"`
		IPAddress      *string   `json:"ip_address,omitempty"`
	}

	// PortCreateOptions represents additional options for port creation
	PortCreateOptions struct {
		Zone *string `json:"zone,omitempty"`
	}

	// PublicIPCreateRequest represents the parameters for creating a public IP
	PublicIPCreateRequest struct {
		Description *string `json:"description,omitempty"`
	}

	// PublicIPCreateResponse represents the response after creating a public IP
	PublicIPCreateResponse struct {
		ID string `json:"id"`
	}

	// PortCreateResponse represents the response after creating a port
	PortCreateResponse struct {
		ID string `json:"id"`
	}

	// PublicIPsList represents a list of public IPs
	PublicIPsList struct {
		PublicIPs []PublicIPDb `json:"public_ips"`
	}

	// IPAddress represents an IP address configuration
	IPAddress struct {
		IPAddress string  `json:"ip_address"`
		SubnetID  string  `json:"subnet_id"`
		EtherType *string `json:"ethertype,omitempty"`
	}

	// PublicIPDb represents a public IP resource
	PublicIPDb struct {
		ID          *string                         `json:"id,omitempty"`
		ExternalID  *string                         `json:"external_id,omitempty"`
		VPCID       *string                         `json:"vpc_id,omitempty"`
		TenantID    *string                         `json:"tenant_id,omitempty"`
		ProjectType *string                         `json:"project_type,omitempty"`
		Description *string                         `json:"description,omitempty"`
		PublicIP    *string                         `json:"public_ip,omitempty"`
		PortID      *string                         `json:"port_id,omitempty"`
		CreatedAt   *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		Updated     *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		Status      *string                         `json:"status,omitempty"`
		Error       *string                         `json:"error,omitempty"`
	}

	// PortListResponse represents a port list response
	PortListResponse struct {
		CreatedAt             *string         `json:"created_at,omitempty"`
		Description           *string         `json:"description,omitempty"`
		ID                    *string         `json:"id,omitempty"`
		IPAddress             []PortIPAddress `json:"ip_address,omitempty"`
		IsAdminStateUp        *bool           `json:"is_admin_state_up,omitempty"`
		IsPortSecurityEnabled *bool           `json:"is_port_security_enabled,omitempty"`
		Name                  *string         `json:"name,omitempty"`
		PublicIP              []PortPublicIP  `json:"public_ip,omitempty"`
		SecurityGroups        []string        `json:"security_groups,omitempty"`
		Updated               *string         `json:"updated,omitempty"`
		VPCID                 *string         `json:"vpc_id,omitempty"`
	}

	// PortIPAddress represents an IP address configuration for a port
	PortIPAddress struct {
		Ethertype *string `json:"ethertype,omitempty"`
		IPAddress string  `json:"ip_address"`
		SubnetID  string  `json:"subnet_id"`
	}

	// PortPublicIP represents a public IP configuration for a port
	PortPublicIP struct {
		PublicIP   *string `json:"public_ip,omitempty"`
		PublicIPID *string `json:"public_ip_id,omitempty"`
	}

	// PortsList represents a list of ports
	PortsList struct {
		Ports           *[]PortResponse      `json:"ports,omitempty"`
		PortsSimplified []PortSimpleResponse `json:"ports_simplified"`
	}

	// PortSimpleResponse represents a simplified port response
	PortSimpleResponse struct {
		ID        *string         `json:"id,omitempty"`
		IPAddress []PortIPAddress `json:"ip_address,omitempty"`
	}
)

// VPCStateV1 represents VPC states
type VPCStateV1 string

const (
	VPCStateNew      VPCStateV1 = "new"
	VPCStateActive   VPCStateV1 = "active"
	VPCStateInactive VPCStateV1 = "inactive"
	VPCStateDeleted  VPCStateV1 = "deleted"
)

// VPCStatusV1 represents VPC statuses
type VPCStatusV1 string

const (
	VPCStatusProvisioning    VPCStatusV1 = "provisioning"
	VPCStatusCreating        VPCStatusV1 = "creating"
	VPCStatusCompleted       VPCStatusV1 = "completed"
	VPCStatusDeletingPending VPCStatusV1 = "deleting_pending"
	VPCStatusDeleting        VPCStatusV1 = "deleting"
	VPCStatusDeleted         VPCStatusV1 = "deleted"
	VPCStatusError           VPCStatusV1 = "error"
)

// VPCService provides operations for managing VPCs
type VPCService interface {
	List(ctx context.Context) ([]VPC, error)
	Get(ctx context.Context, id string) (*VPC, error)
	Create(ctx context.Context, req CreateVPCRequest) (string, error)
	Delete(ctx context.Context, id string) error
	Rename(ctx context.Context, id string, newName string) error
	ListPorts(ctx context.Context, vpcID string, detailed bool, opts ListOptions) (*PortsList, error)
	CreatePort(ctx context.Context, vpcID string, req PortCreateRequest, opts PortCreateOptions) (string, error)
	ListPublicIPs(ctx context.Context, vpcID string) ([]PublicIPDb, error)
	CreatePublicIP(ctx context.Context, vpcID string, req PublicIPCreateRequest) (string, error)
	ListSubnets(ctx context.Context, vpcID string) ([]SubnetResponse, error)
	CreateSubnet(ctx context.Context, vpcID string, req SubnetCreateRequest, opts SubnetCreateOptions) (string, error)
}

// vpcService implements the VPCService interface
type vpcService struct {
	client *NetworkClient
}

// List returns a slice of VPCs based on the provided listing options
func (s *vpcService) List(ctx context.Context) ([]VPC, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListVPCsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		"/v0/vpcs",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result.VPCs, nil
}

// Get retrieves detailed information about a specific VPC
func (s *vpcService) Get(ctx context.Context, id string) (*VPC, error) {
	path := fmt.Sprintf("/v0/vpcs/%s", id)
	return mgc_http.ExecuteSimpleRequestWithRespBody[VPC](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		nil,
	)
}

// Create provisions a new VPC
func (s *vpcService) Create(ctx context.Context, req CreateVPCRequest) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[CreateVPCResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/vpcs",
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Delete removes a VPC
func (s *vpcService) Delete(ctx context.Context, id string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodDelete,
		fmt.Sprintf("/v0/vpcs/%s", id),
		nil,
		nil,
	)
}

// Rename updates the display name of an existing VPC
func (s *vpcService) Rename(ctx context.Context, id string, newName string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		fmt.Sprintf("/v0/vpcs/%s/rename", id),
		RenameVPCRequest{Name: newName},
		nil,
	)
}

// ListPorts returns all ports for a VPC
func (s *vpcService) ListPorts(ctx context.Context, vpcID string, detailed bool, opts ListOptions) (*PortsList, error) {
	query := makeListOptionsQuery(opts)
	query.Set("detailed", strconv.FormatBool(detailed))

	path := fmt.Sprintf("/v0/vpcs/%s/ports", vpcID)

	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[PortsList](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		query,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreatePort creates a new port in a VPC
func (s *vpcService) CreatePort(ctx context.Context, vpcID string, req PortCreateRequest, opts PortCreateOptions) (string, error) {
	nreq, err := s.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/vpcs/%s/ports", vpcID), req)
	if err != nil {
		return "", err
	}

	if opts.Zone != nil {
		nreq.Header.Set("x-zone", *opts.Zone)
	}

	result, err := mgc_http.Do(s.client.GetConfig(), ctx, nreq, &PortCreateResponse{})
	if err != nil {
		return "", err
	}

	return result.ID, nil
}

// ListPublicIPs returns all public IPs for a VPC
func (s *vpcService) ListPublicIPs(ctx context.Context, vpcID string) ([]PublicIPDb, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[PublicIPsList](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/vpcs/%s/public_ips", vpcID),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result.PublicIPs, nil
}

// CreatePublicIP creates a new public IP in a VPC
func (s *vpcService) CreatePublicIP(ctx context.Context, vpcID string, req PublicIPCreateRequest) (string, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[PublicIPCreateResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		fmt.Sprintf("/v0/vpcs/%s/public_ips", vpcID),
		req,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// ListSubnets returns all subnets in a VPC
func (s *vpcService) ListSubnets(ctx context.Context, vpcID string) ([]SubnetResponse, error) {
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[ListSubnetsResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		fmt.Sprintf("/v0/vpcs/%s/subnets", vpcID),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result.Subnets, nil
}

// CreateSubnet creates a new subnet in a VPC
func (s *vpcService) CreateSubnet(ctx context.Context, vpcID string, req SubnetCreateRequest, opts SubnetCreateOptions) (string, error) {
	nreq, err := s.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/vpcs/%s/subnets", vpcID), req)
	if err != nil {
		return "", err
	}

	if opts.Zone != nil {
		nreq.Header.Set("x-zone", *opts.Zone)
	}

	result, err := mgc_http.Do(s.client.GetConfig(), ctx, nreq, &SubnetCreateResponse{})
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// makeListOptionsQuery converts ListOptions to URL query parameters
func makeListOptionsQuery(opts ListOptions) url.Values {
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
	return query
}
