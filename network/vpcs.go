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
	// SecurityGroupsExpand is used to expand security groups in responses
	SecurityGroupsExpand = "security_groups"
	// SubnetsExpand is used to expand subnets in responses
	SubnetsExpand = "subnets"
	// PortStatusProvisioning indicates a port is being provisioned
	PortStatusProvisioning = "provisioning"
	// PortStatusActive indicates a port is active and ready
	PortStatusActive = "active"
	// PortStatusError indicates a port has encountered an error
	PortStatusError = "error"
	// PublicIPStatusCreated indicates a public IP has been created
	PublicIPStatusCreated = "created"
	// PublicIPStatusPending indicates a public IP is pending creation
	PublicIPStatusPending = "pending"
	// PublicIPStatusError indicates a public IP has encountered an error
	PublicIPStatusError = "error"
)

type (
	// ListVPCsResponse represents a list of VPCs response
	ListVPCsResponse struct {
		// VPCs contains the list of VPC resources
		VPCs []VPC `json:"vpcs"`
	}

	// VPC represents a Virtual Private Cloud resource
	VPC struct {
		// ID is the unique identifier of the VPC
		ID *string `json:"id,omitempty"`
		// TenantID is the tenant identifier
		TenantID *string `json:"tenant_id,omitempty"`
		// Name is the display name of the VPC
		Name *string `json:"name,omitempty"`
		// Description is the description of the VPC
		Description *string `json:"description,omitempty"`
		// Status is the current status of the VPC
		Status string `json:"status"`
		// RouterID is the ID of the associated router
		RouterID *string `json:"router_id,omitempty"`
		// ExternalNetwork is the external network identifier
		ExternalNetwork *string `json:"external_network,omitempty"`
		// NetworkID is the network identifier
		NetworkID *string `json:"network_id,omitempty"`
		// Subnets contains the list of subnet IDs
		Subnets *[]string `json:"subnets,omitempty"`
		// SecurityGroups contains the list of security group IDs
		SecurityGroups *[]string `json:"security_groups,omitempty"`
		// CreatedAt is the creation timestamp
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Updated is the last update timestamp
		Updated *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		// IsDefault indicates if this is the default VPC
		IsDefault *bool `json:"is_default,omitempty"`
	}

	// CreateVPCRequest represents the parameters for creating a new VPC
	CreateVPCRequest struct {
		// Name is the name of the VPC
		Name string `json:"name"`
		// Description is the description of the VPC (optional)
		Description *string `json:"description,omitempty"`
	}

	// RenameVPCRequest represents the parameters for renaming a VPC
	RenameVPCRequest struct {
		// Name is the new name for the VPC
		Name string `json:"name"`
	}

	// ListOptions represents parameters for filtering and pagination
	ListOptions struct {
		// Limit specifies the maximum number of items to return
		Limit *int
		// Offset specifies the number of items to skip
		Offset *int
		// Sort specifies the field and direction for sorting results
		Sort *string
	}

	// CreateVPCResponse represents the response after creating a VPC
	CreateVPCResponse struct {
		// ID is the unique identifier of the created VPC
		ID string `json:"id"`
		// Status is the status of the created VPC
		Status string `json:"status"`
	}

	// PortCreateRequest represents the parameters for creating a port
	PortCreateRequest struct {
		// Name is the name of the port
		Name string `json:"name"`
		// HasPIP indicates if the port should have a public IP (optional)
		HasPIP *bool `json:"has_pip,omitempty"`
		// HasSG indicates if the port should have security groups (optional)
		HasSG *bool `json:"has_sg,omitempty"`
		// Subnets contains the list of subnet IDs (optional)
		Subnets *[]string `json:"subnets,omitempty"`
		// SecurityGroups contains the list of security group IDs (optional)
		SecurityGroups *[]string `json:"security_groups_id,omitempty"`
	}

	// PortCreateOptions represents additional options for port creation
	PortCreateOptions struct {
		// Zone specifies the availability zone (optional)
		Zone *string `json:"zone,omitempty"`
	}

	// PublicIPCreateRequest represents the parameters for creating a public IP
	PublicIPCreateRequest struct {
		// Description is the description of the public IP (optional)
		Description *string `json:"description,omitempty"`
	}

	// PublicIPCreateResponse represents the response after creating a public IP
	PublicIPCreateResponse struct {
		// ID is the unique identifier of the created public IP
		ID string `json:"id"`
	}

	// PortCreateResponse represents the response after creating a port
	PortCreateResponse struct {
		// ID is the unique identifier of the created port
		ID string `json:"id"`
	}

	// PublicIPsList represents a list of public IPs
	PublicIPsList struct {
		// PublicIPs contains the list of public IP resources
		PublicIPs []PublicIPDb `json:"public_ips"`
	}

	// IPAddress represents an IP address configuration
	IPAddress struct {
		// IPAddress is the IP address
		IPAddress string `json:"ip_address"`
		// SubnetID is the subnet identifier
		SubnetID string `json:"subnet_id"`
		// EtherType is the ethernet type (optional)
		EtherType *string `json:"ethertype,omitempty"`
	}

	// PublicIPDb represents a public IP resource
	PublicIPDb struct {
		// ID is the unique identifier of the public IP
		ID *string `json:"id,omitempty"`
		// ExternalID is the external identifier
		ExternalID *string `json:"external_id,omitempty"`
		// VPCID is the VPC identifier
		VPCID *string `json:"vpc_id,omitempty"`
		// TenantID is the tenant identifier
		TenantID *string `json:"tenant_id,omitempty"`
		// ProjectType is the project type
		ProjectType *string `json:"project_type,omitempty"`
		// Description is the description of the public IP
		Description *string `json:"description,omitempty"`
		// PublicIP is the public IP address
		PublicIP *string `json:"public_ip,omitempty"`
		// PortID is the associated port identifier
		PortID *string `json:"port_id,omitempty"`
		// CreatedAt is the creation timestamp
		CreatedAt *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
		// Updated is the last update timestamp
		Updated *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
		// Status is the current status of the public IP
		Status *string `json:"status,omitempty"`
		// Error contains error information if any
		Error *string `json:"error,omitempty"`
	}

	// PortListResponse represents a port list response
	PortListResponse struct {
		// CreatedAt is the creation timestamp
		CreatedAt *string `json:"created_at,omitempty"`
		// Description is the description of the port
		Description *string `json:"description,omitempty"`
		// ID is the unique identifier of the port
		ID *string `json:"id,omitempty"`
		// IPAddress contains the IP address configurations
		IPAddress []PortIPAddress `json:"ip_address,omitempty"`
		// IsAdminStateUp indicates if the port is administratively up
		IsAdminStateUp *bool `json:"is_admin_state_up,omitempty"`
		// IsPortSecurityEnabled indicates if port security is enabled
		IsPortSecurityEnabled *bool `json:"is_port_security_enabled,omitempty"`
		// Name is the name of the port
		Name *string `json:"name,omitempty"`
		// PublicIP contains the public IP configurations
		PublicIP []PortPublicIP `json:"public_ip,omitempty"`
		// SecurityGroups contains the list of security group IDs
		SecurityGroups []string `json:"security_groups,omitempty"`
		// Updated is the last update timestamp
		Updated *string `json:"updated,omitempty"`
		// VPCID is the VPC identifier
		VPCID *string `json:"vpc_id,omitempty"`
	}

	// PortIPAddress represents an IP address configuration for a port
	PortIPAddress struct {
		// Ethertype is the ethernet type (optional)
		Ethertype *string `json:"ethertype,omitempty"`
		// IPAddress is the IP address
		IPAddress string `json:"ip_address"`
		// SubnetID is the subnet identifier
		SubnetID string `json:"subnet_id"`
	}

	// PortPublicIP represents a public IP configuration for a port
	PortPublicIP struct {
		// PublicIP is the public IP address (optional)
		PublicIP *string `json:"public_ip,omitempty"`
		// PublicIPID is the public IP identifier (optional)
		PublicIPID *string `json:"public_ip_id,omitempty"`
	}

	// PortsList represents a list of ports
	PortsList struct {
		// Ports contains the detailed port responses (optional)
		Ports *[]PortResponse `json:"ports,omitempty"`
		// PortsSimplified contains the simplified port responses
		PortsSimplified []PortSimpleResponse `json:"ports_simplified"`
	}

	// PortSimpleResponse represents a simplified port response
	PortSimpleResponse struct {
		// ID is the unique identifier of the port
		ID *string `json:"id,omitempty"`
		// IPAddress contains the IP address configurations
		IPAddress []PortIPAddress `json:"ip_address,omitempty"`
	}
)

// VPCStateV1 represents VPC states
type VPCStateV1 string

const (
	// VPCStateNew indicates a VPC is in new state
	VPCStateNew VPCStateV1 = "new"
	// VPCStateActive indicates a VPC is active
	VPCStateActive VPCStateV1 = "active"
	// VPCStateInactive indicates a VPC is inactive
	VPCStateInactive VPCStateV1 = "inactive"
	// VPCStateDeleted indicates a VPC has been deleted
	VPCStateDeleted VPCStateV1 = "deleted"
)

// VPCStatusV1 represents VPC statuses
type VPCStatusV1 string

const (
	// VPCStatusProvisioning indicates a VPC is being provisioned
	VPCStatusProvisioning VPCStatusV1 = "provisioning"
	// VPCStatusCreating indicates a VPC is being created
	VPCStatusCreating VPCStatusV1 = "creating"
	// VPCStatusCompleted indicates a VPC creation is completed
	VPCStatusCompleted VPCStatusV1 = "completed"
	// VPCStatusDeletingPending indicates a VPC deletion is pending
	VPCStatusDeletingPending VPCStatusV1 = "deleting_pending"
	// VPCStatusDeleting indicates a VPC is being deleted
	VPCStatusDeleting VPCStatusV1 = "deleting"
	// VPCStatusDeleted indicates a VPC has been deleted
	VPCStatusDeleted VPCStatusV1 = "deleted"
	// VPCStatusError indicates a VPC has encountered an error
	VPCStatusError VPCStatusV1 = "error"
)

// VPCService provides operations for managing VPCs
type VPCService interface {
	// List returns a slice of VPCs based on the provided listing options
	List(ctx context.Context) ([]VPC, error)

	// Get retrieves detailed information about a specific VPC
	Get(ctx context.Context, id string) (*VPC, error)

	// Create provisions a new VPC
	Create(ctx context.Context, req CreateVPCRequest) (string, error)

	// Delete removes a VPC
	Delete(ctx context.Context, id string) error

	// Rename updates the display name of an existing VPC
	Rename(ctx context.Context, id string, newName string) error

	// ListPorts returns all ports for a VPC
	ListPorts(ctx context.Context, vpcID string, detailed bool, opts ListOptions) (*PortsList, error)

	// CreatePort creates a new port in a VPC
	CreatePort(ctx context.Context, vpcID string, req PortCreateRequest, opts PortCreateOptions) (string, error)

	// ListPublicIPs returns all public IPs for a VPC
	ListPublicIPs(ctx context.Context, vpcID string) ([]PublicIPDb, error)

	// CreatePublicIP creates a new public IP in a VPC
	CreatePublicIP(ctx context.Context, vpcID string, req PublicIPCreateRequest) (string, error)

	// ListSubnets returns all subnets in a VPC
	ListSubnets(ctx context.Context, vpcID string) ([]SubnetResponse, error)

	// CreateSubnet creates a new subnet in a VPC
	CreateSubnet(ctx context.Context, vpcID string, req SubnetCreateRequest, opts SubnetCreateOptions) (string, error)
}

// vpcService implements the VPCService interface
type vpcService struct {
	client *NetworkClient
}

// List returns a slice of VPCs based on the provided listing options
func (s *vpcService) List(ctx context.Context) ([]VPC, error) {
	path := "/v1/vpcs"

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListVPCsResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.VPCs, nil
}

// Get retrieves detailed information about a specific VPC
func (s *vpcService) Get(ctx context.Context, id string) (*VPC, error) {
	path := fmt.Sprintf("/v1/vpcs/%s", id)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp VPC
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Create provisions a new VPC
func (s *vpcService) Create(ctx context.Context, req CreateVPCRequest) (string, error) {
	path := "/v1/vpcs"

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return "", err
	}

	var resp CreateVPCResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Delete removes a VPC
func (s *vpcService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/v1/vpcs/%s", id)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// Rename updates the display name of an existing VPC
func (s *vpcService) Rename(ctx context.Context, id string, newName string) error {
	path := fmt.Sprintf("/v1/vpcs/%s", id)
	req := RenameVPCRequest{Name: newName}

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

// ListPorts returns all ports for a VPC
func (s *vpcService) ListPorts(ctx context.Context, vpcID string, detailed bool, opts ListOptions) (*PortsList, error) {
	path := fmt.Sprintf("/v1/vpcs/%s/ports", vpcID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := makeListOptionsQuery(opts)
	if detailed {
		query.Set("detailed", "true")
	}
	httpReq.URL.RawQuery = query.Encode()

	var resp PortsList
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreatePort creates a new port in a VPC
func (s *vpcService) CreatePort(ctx context.Context, vpcID string, req PortCreateRequest, opts PortCreateOptions) (string, error) {
	path := fmt.Sprintf("/v1/vpcs/%s/ports", vpcID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return "", err
	}

	query := httpReq.URL.Query()
	if opts.Zone != nil {
		query.Set("zone", *opts.Zone)
	}
	httpReq.URL.RawQuery = query.Encode()

	var resp PortCreateResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// ListPublicIPs returns all public IPs for a VPC
func (s *vpcService) ListPublicIPs(ctx context.Context, vpcID string) ([]PublicIPDb, error) {
	path := fmt.Sprintf("/v1/vpcs/%s/public-ips", vpcID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp PublicIPsList
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.PublicIPs, nil
}

// CreatePublicIP creates a new public IP in a VPC
func (s *vpcService) CreatePublicIP(ctx context.Context, vpcID string, req PublicIPCreateRequest) (string, error) {
	path := fmt.Sprintf("/v1/vpcs/%s/public-ips", vpcID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return "", err
	}

	var resp PublicIPCreateResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// ListSubnets returns all subnets in a VPC
func (s *vpcService) ListSubnets(ctx context.Context, vpcID string) ([]SubnetResponse, error) {
	path := fmt.Sprintf("/v1/vpcs/%s/subnets", vpcID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Subnets []SubnetResponse `json:"subnets"`
	}
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Subnets, nil
}

// CreateSubnet creates a new subnet in a VPC
func (s *vpcService) CreateSubnet(ctx context.Context, vpcID string, req SubnetCreateRequest, opts SubnetCreateOptions) (string, error) {
	path := fmt.Sprintf("/v1/vpcs/%s/subnets", vpcID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return "", err
	}

	query := httpReq.URL.Query()
	if opts.Zone != nil {
		query.Set("zone", *opts.Zone)
	}
	httpReq.URL.RawQuery = query.Encode()

	var resp struct {
		ID string `json:"id"`
	}
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// makeListOptionsQuery converts ListOptions to URL query parameters
func makeListOptionsQuery(opts ListOptions) url.Values {
	query := url.Values{}
	if opts.Limit != nil {
		query.Set("limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("sort", *opts.Sort)
	}
	return query
}
