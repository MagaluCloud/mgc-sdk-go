# Network

Package network provides a client for interacting with the Magalu Cloud Network
API. This package allows you to manage VPCs, subnets, ports, security groups,
rules, public IPs, subnet pools, and NAT gateways.
```
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
```
```
const (
DefaultBasePath = "/network"
)


```
```
type BookCIDRRequest struct {
CIDR *string `json:"cidr,omitempty"`
Mask *int    `json:"mask,omitempty"`
}
BookCIDRRequest represents parameters for booking a CIDR range

```
```
type BookCIDRResponse struct {
CIDR string `json:"cidr"`
}
BookCIDRResponse represents the response after booking a CIDR range

```
```
type CreateNatGatewayRequest struct {
Name        string  `json:"name"`
Description *string `json:"description,omitempty"`
Zone        string  `json:"zone"`
VPCID       string  `json:"vpc_id"`
}
CreateNatGatewayRequest represents the parameters for creating a new NAT
Gateway

```
```
type CreateSubnetPoolRequest struct {
CIDR        *string `json:"cidr,omitempty"`
Name        string  `json:"name"`
Description string  `json:"description"`
Type        *string `json:"type,omitempty"`
}
CreateSubnetPoolRequest represents parameters for creating a new subnet pool

```
```
type CreateSubnetPoolResponse struct {
ID string `json:"id"`
}
CreateSubnetPoolResponse represents the response after creating a subnet
pool

```
```
type CreateVPCRequest struct {
Name        string  `json:"name"`
Description *string `json:"description,omitempty"`
}
CreateVPCRequest represents the parameters for creating a new VPC

```
```
type CreateVPCResponse struct {
ID     string `json:"id"`
Status string `json:"status"`
}
CreateVPCResponse represents the response after creating a VPC

```
```
type DHCPPoolStr struct {
Start string `json:"start"`
End   string `json:"end"`
}
DHCPPoolStr represents a DHCP pool configuration

```
```
type IPAddress struct {
IPAddress string  `json:"ip_address"`
SubnetID  string  `json:"subnet_id"`
EtherType *string `json:"ethertype,omitempty"`
}
IPAddress represents an IP address configuration

```
```
type IpAddress struct {
IPAddress string  `json:"ip_address"`
SubnetID  string  `json:"subnet_id"`
Ethertype *string `json:"ethertype,omitempty"`
}
IpAddress represents an IP address configuration for a port

```
```
type LinkModel struct {
Previous *string `json:"previous"`
Next     *string `json:"next"`
Self     string  `json:"self"`
}
LinkModel represents navigation links

```
```
type ListOptions struct {
Limit  *int
Offset *int
Sort   *string
}
ListOptions represents parameters for filtering and pagination

```
```
type ListSubnetPoolsOptions struct {
Limit  *int
Offset *int
Sort   *string
}
ListSubnetPoolsOptions represents parameters for filtering and pagination

```
```
type ListSubnetPoolsResponse struct {
Meta    MetaModel            `json:"meta"`
Results []SubnetPoolResponse `json:"results"`
}
ListSubnetPoolsResponse represents a list of subnet pools response

```
```
type ListSubnetsResponse struct {
Subnets []SubnetResponse `json:"subnets"`
}
ListSubnetsResponse represents a list of subnets response

```
```
type ListVPCsResponse struct {
VPCs []VPC `json:"vpcs"`
}
ListVPCsResponse represents a list of VPCs response

```
```
type Meta struct {
Page  MetaPageInfo `json:"page"`
Links MetaLinks    `json:"links"`
}
Meta represents pagination metadata

```
```
type MetaLinks struct {
Previous *string `json:"previous,omitempty"`
Next     *string `json:"next,omitempty"`
Self     string  `json:"self"`
}
MetaLinks represents navigation links

```
```
type MetaModel struct {
Page  PageModel `json:"page"`
Links LinkModel `json:"links"`
}
MetaModel represents pagination metadata

```
```
type MetaPageInfo struct {
Limit           int `json:"limit"`
Offset          int `json:"offset"`
Count           int `json:"count"`
Total           int `json:"total"`
MaxItemsPerPage int `json:"max_items_per_page"`
}
MetaPageInfo represents pagination information

```
```
type NatGatewayCreateResponse struct {
ID     string `json:"id"`
Status string `json:"status"`
}
NatGatewayCreateResponse represents the response after creating a NAT
Gateway

```
```
type NatGatewayDetailsResponse struct {
ID           *string                         `json:"id,omitempty"`
Name         *string                         `json:"name,omitempty"`
Description  *string                         `json:"description,omitempty"`
VPCID        *string                         `json:"vpc_id,omitempty"`
Zone         *string                         `json:"zone,omitempty"`
NatGatewayIP *string                         `json:"nat_gateway_ip,omitempty"`
CreatedAt    *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
Updated      *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
Status       string                          `json:"status"`
}
NatGatewayDetailsResponse represents detailed information about a NAT
Gateway

```
```
type NatGatewayListResponse struct {
Meta   Meta                 `json:"meta"`
Result []NatGatewayResponse `json:"result"`
}
NatGatewayListResponse represents a NAT Gateway listing response

```
```
type NatGatewayResponse struct {
ID           *string                         `json:"id,omitempty"`
Name         *string                         `json:"name,omitempty"`
Description  *string                         `json:"description,omitempty"`
VPCID        *string                         `json:"vpc_id,omitempty"`
Zone         *string                         `json:"zone,omitempty"`
NatGatewayIP *string                         `json:"nat_gateway_ip,omitempty"`
CreatedAt    *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
Updated      *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
Status       string                          `json:"status"`
}
NatGatewayResponse represents a NAT Gateway resource

```
```
type NatGatewayService interface {
Create(ctx context.Context, req CreateNatGatewayRequest) (string, error)
Delete(ctx context.Context, id string) error
Get(ctx context.Context, id string) (*NatGatewayDetailsResponse, error)
List(ctx context.Context, vpcID string, opts ListOptions) ([]NatGatewayResponse, error)
}
NatGatewayService provides operations for managing NAT Gateways

```
```
type NetworkClient struct {
*client.CoreClient
}
NetworkClient represents a client for interacting with the network services

```
```
func New(core *client.CoreClient) *NetworkClient
New creates a new NetworkClient instance

```
```
func (c *NetworkClient) NatGateways() NatGatewayService
NatGateways returns a service for managing NAT gateway resources

```
```
func (c *NetworkClient) Ports() PortService
Ports returns a service for managing port resources

```
```
func (c *NetworkClient) PublicIPs() PublicIPService
PublicIPs returns a service for managing public IP resources

```
```
func (c *NetworkClient) Rules() RuleService
Rules returns a service for managing security group rule resources

```
```
func (c *NetworkClient) SecurityGroups() SecurityGroupService
SecurityGroups returns a service for managing security group resources

```
```
func (c *NetworkClient) SubnetPools() SubnetPoolService
SubnetPools returns a service for managing subnet pool resources

```
```
func (c *NetworkClient) Subnets() SubnetService
Subnets returns a service for managing subnet resources

```
```
func (c *NetworkClient) VPCs() VPCService
VPCs returns a service for managing VPC resources

```
```
type PageModel struct {
Limit  *int `json:"limit,omitempty"`
Offset *int `json:"offset,omitempty"`
Count  int  `json:"count"`
Total  int  `json:"total"`
}
PageModel represents page information

```
```
type PortCreateOptions struct {
Zone *string `json:"zone,omitempty"`
}
PortCreateOptions represents additional options for port creation

```
```
type PortCreateRequest struct {
Name           string    `json:"name"`
HasPIP         *bool     `json:"has_pip,omitempty"`
HasSG          *bool     `json:"has_sg,omitempty"`
Subnets        *[]string `json:"subnets,omitempty"`
SecurityGroups *[]string `json:"security_groups_id,omitempty"`
}
PortCreateRequest represents the parameters for creating a port

```
```
type PortCreateResponse struct {
ID string `json:"id"`
}
PortCreateResponse represents the response after creating a port

```
```
type PortIPAddress struct {
Ethertype *string `json:"ethertype,omitempty"`
IPAddress string  `json:"ip_address"`
SubnetID  string  `json:"subnet_id"`
}
PortIPAddress represents an IP address configuration for a port

```
```
type PortListResponse struct {
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
PortListResponse represents a port list response

```
```
type PortNetworkResponse struct {
AvailabilityZone *string `json:"availability_zone,omitempty"`
ID               *string `json:"id,omitempty"`
Zone             *string `json:"zone,omitempty"`
}
PortNetworkResponse represents the AvailabilityZone associated with a port

```
```
type PortPublicIP struct {
PublicIP   *string `json:"public_ip,omitempty"`
PublicIPID *string `json:"public_ip_id,omitempty"`
}
PortPublicIP represents a public IP configuration for a port

```
```
type PortResponse struct {
ID                    *string                         `json:"id,omitempty"`
Name                  *string                         `json:"name,omitempty"`
Description           *string                         `json:"description,omitempty"`
IsAdminStateUp        *bool                           `json:"is_admin_state_up,omitempty"`
VPCID                 *string                         `json:"vpc_id,omitempty"`
IsPortSecurityEnabled *bool                           `json:"is_port_security_enabled,omitempty"`
SecurityGroups        *[]string                       `json:"security_groups"`
PublicIP              *[]PublicIpResponsePort         `json:"public_ip"`
IPAddress             *[]IpAddress                    `json:"ip_address"`
CreatedAt             *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
Updated               *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
Network               *PortNetworkResponse            `json:"network,omitempty"`
}
PortResponse represents a network port resource

```
```
type PortService interface {
List(ctx context.Context) ([]PortResponse, error)
Get(ctx context.Context, id string) (*PortResponse, error)
Delete(ctx context.Context, id string) error
Update(ctx context.Context, id string, req PortUpdateRequest) error
AttachSecurityGroup(ctx context.Context, portID string, securityGroupID string) error
DetachSecurityGroup(ctx context.Context, portID string, securityGroupID string) error
}
PortService provides operations for managing network ports

```
```
type PortSimpleResponse struct {
ID        *string         `json:"id,omitempty"`
IPAddress []PortIPAddress `json:"ip_address,omitempty"`
}
PortSimpleResponse represents a simplified port response

```
```
type PortUpdateRequest struct {
IPSpoofingGuard *bool `json:"ip_spoofing_guard,omitempty"`
}
PortUpdateRequest represents the fields available for update in a port
resource

```
```
type PortsList struct {
Ports           *[]PortResponse      `json:"ports,omitempty"`
PortsSimplified []PortSimpleResponse `json:"ports_simplified"`
}
PortsList represents a list of ports

```
```
type PublicIPCreateRequest struct {
Description *string `json:"description,omitempty"`
}
PublicIPCreateRequest represents the parameters for creating a public IP

```
```
type PublicIPCreateResponse struct {
ID string `json:"id"`
}
PublicIPCreateResponse represents the response after creating a public IP

```
```
type PublicIPDb struct {
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
PublicIPDb represents a public IP resource

```
```
type PublicIPListResponse struct {
PublicIPs []PublicIPResponse `json:"public_ips"`
}
PublicIPListResponse represents a list of public IPs response

```
```
type PublicIPResponse struct {
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
PublicIPResponse represents a public IP resource response

```
```
type PublicIPService interface {
List(ctx context.Context) ([]PublicIPResponse, error)
Get(ctx context.Context, id string) (*PublicIPResponse, error)
Delete(ctx context.Context, id string) error
AttachToPort(ctx context.Context, publicIPID string, portID string) error
DetachFromPort(ctx context.Context, publicIPID string, portID string) error
}
PublicIPService provides operations for managing Public IPs

```
```
type PublicIPsList struct {
PublicIPs []PublicIPDb `json:"public_ips"`
}
PublicIPsList represents a list of public IPs

```
```
type PublicIpResponsePort struct {
PublicIPID *string `json:"public_ip_id,omitempty"`
PublicIP   *string `json:"public_ip,omitempty"`
}
PublicIpResponsePort represents a public IP associated with a port

```
```
type RenameVPCRequest struct {
Name string `json:"name"`
}
RenameVPCRequest represents the parameters for renaming a VPC

```
```
type RuleCreateRequest struct {
Direction      *string `json:"direction,omitempty"`
PortRangeMin   *int    `json:"port_range_min,omitempty"`
PortRangeMax   *int    `json:"port_range_max,omitempty"`
Protocol       *string `json:"protocol,omitempty"`
RemoteIPPrefix *string `json:"remote_ip_prefix,omitempty"`
EtherType      string  `json:"ethertype"`
Description    *string `json:"description,omitempty"`
}
RuleCreateRequest represents the parameters for creating a new security
group rule

```
```
type RuleCreateResponse struct {
ID string `json:"id"`
}
RuleCreateResponse represents the response after creating a security group
rule

```
```
type RuleResponse struct {
ID              *string                         `json:"id,omitempty"`
ExternalID      *string                         `json:"external_id,omitempty"`
SecurityGroupID *string                         `json:"security_group_id,omitempty"`
Direction       *string                         `json:"direction,omitempty"`
PortRangeMin    *int                            `json:"port_range_min,omitempty"`
PortRangeMax    *int                            `json:"port_range_max,omitempty"`
Protocol        *string                         `json:"protocol,omitempty"`
RemoteIPPrefix  *string                         `json:"remote_ip_prefix,omitempty"`
EtherType       *string                         `json:"ethertype"`
CreatedAt       *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
Status          string                          `json:"status"`
Error           *string                         `json:"error,omitempty"`
Description     *string                         `json:"description,omitempty"`
}
RuleResponse represents a security group rule resource

```
```
type RuleService interface {
List(ctx context.Context, securityGroupID string) ([]RuleResponse, error)
Get(ctx context.Context, id string) (*RuleResponse, error)
Create(ctx context.Context, securityGroupID string, req RuleCreateRequest) (string, error)
Delete(ctx context.Context, id string) error
}
RuleService provides operations for managing security group rules

```
```
type RulesList struct {
Rules []RuleResponse `json:"rules"`
}
RulesList represents a list of security group rules

```
```
type SecurityGroupCreateRequest struct {
Name             string  `json:"name"`
Description      *string `json:"description,omitempty"`
SkipDefaultRules *bool   `json:"skip_default_rules,omitempty"`
}
SecurityGroupCreateRequest represents the parameters for creating a new
security group

```
```
type SecurityGroupCreateResponse struct {
ID string `json:"id"`
}
SecurityGroupCreateResponse represents the response after creating a
security group

```
```
type SecurityGroupDetailResponse struct {
SecurityGroupResponse
ExternalID *string         `json:"external_id,omitempty"`
Rules      *[]RuleResponse `json:"rules"`
}
SecurityGroupDetailResponse represents detailed information about a security
group

```
```
type SecurityGroupListResponse struct {
SecurityGroups []SecurityGroupResponse `json:"security_groups"`
}
SecurityGroupListResponse represents a list of security groups response

```
```
type SecurityGroupResponse struct {
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
SecurityGroupResponse represents a security group resource

```
```
type SecurityGroupService interface {
List(ctx context.Context) ([]SecurityGroupResponse, error)
Get(ctx context.Context, id string) (*SecurityGroupDetailResponse, error)
Create(ctx context.Context, req SecurityGroupCreateRequest) (string, error)
Delete(ctx context.Context, id string) error
}
SecurityGroupService provides operations for managing security groups

```
```
type SubnetCreateOptions struct {
Zone *string `json:"zone,omitempty"`
}
SubnetCreateOptions represents additional options for subnet creation

```
```
type SubnetCreateRequest struct {
Name           string    `json:"name"`
Description    *string   `json:"description,omitempty"`
CIDRBlock      string    `json:"cidr_block"`
IPVersion      int       `json:"ip_version"`
DNSNameservers *[]string `json:"dns_nameservers,omitempty"`
SubnetPoolID   *string   `json:"subnetpool_id,omitempty"`
}
SubnetCreateRequest represents parameters for creating a new subnet

```
```
type SubnetCreateResponse struct {
ID string `json:"id"`
}
SubnetCreateResponse represents the response after creating a subnet

```
```
type SubnetPatchRequest struct {
DNSNameservers *[]string `json:"dns_nameservers,omitempty"`
}
SubnetPatchRequest represents parameters for updating a subnet

```
```
type SubnetPoolDetailsResponse struct {
CIDR        *string                        `json:"cidr,omitempty"`
ID          string                         `json:"id"`
CreatedAt   utils.LocalDateTimeWithoutZone `json:"created_at"`
TenantID    string                         `json:"tenant_id"`
IPVersion   int                            `json:"ip_version"`
IsDefault   bool                           `json:"is_default"`
Name        string                         `json:"name"`
Description string                         `json:"description"`
}
SubnetPoolDetailsResponse represents detailed subnet pool information

```
```
type SubnetPoolResponse struct {
CIDR        *string `json:"cidr,omitempty"`
ID          string  `json:"id"`
Name        string  `json:"name"`
TenantID    string  `json:"tenant_id"`
Description *string `json:"description,omitempty"`
IsDefault   bool    `json:"is_default"`
}
SubnetPoolResponse represents a subnet pool resource response

```
```
type SubnetPoolService interface {
List(ctx context.Context, opts ListOptions) ([]SubnetPoolResponse, error)
Get(ctx context.Context, id string) (*SubnetPoolDetailsResponse, error)
Create(ctx context.Context, req CreateSubnetPoolRequest) (string, error)
Delete(ctx context.Context, id string) error
BookCIDR(ctx context.Context, id string, req BookCIDRRequest) (*BookCIDRResponse, error)
UnbookCIDR(ctx context.Context, id string, req UnbookCIDRRequest) error
}
SubnetPoolService provides operations for managing subnet pools

```
```
type SubnetResponse struct {
ID           string                          `json:"id"`
VPCID        string                          `json:"vpc_id"`
Name         *string                         `json:"name,omitempty"`
Description  *string                         `json:"description,omitempty"`
CIDRBlock    string                          `json:"cidr_block"`
SubnetPoolID string                          `json:"subnetpool_id"`
IPVersion    string                          `json:"ip_version"`
Zone         string                          `json:"zone"`
CreatedAt    *utils.LocalDateTimeWithoutZone `json:"created_at,omitempty"`
Updated      *utils.LocalDateTimeWithoutZone `json:"updated,omitempty"`
}
SubnetResponse represents a subnet resource response

```
```
type SubnetResponseDetail struct {
SubnetResponse
GatewayIP      string        `json:"gateway_ip"`
DNSNameservers []string      `json:"dns_nameservers"`
DHCPPools      []DHCPPoolStr `json:"dhcp_pools"`
}
SubnetResponseDetail represents a detailed subnet response

```
```
type SubnetResponseId struct {
ID string `json:"id"`
}
SubnetResponseId represents a subnet ID response

```
```
type SubnetService interface {
Get(ctx context.Context, id string) (*SubnetResponseDetail, error)
Delete(ctx context.Context, id string) error
Update(ctx context.Context, id string, req SubnetPatchRequest) (*SubnetResponseId, error)
}
SubnetService provides operations for managing subnets

```
```
type UnbookCIDRRequest struct {
CIDR string `json:"cidr"`
}
UnbookCIDRRequest represents parameters for unbooking a CIDR range

```
```
type VPC struct {
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
VPC represents a Virtual Private Cloud resource

```
```
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
VPCService provides operations for managing VPCs

```
```
type VPCStateV1 string
VPCStateV1 represents VPC states

```
```
const (
VPCStateNew      VPCStateV1 = "new"
VPCStateActive   VPCStateV1 = "active"
VPCStateInactive VPCStateV1 = "inactive"
VPCStateDeleted  VPCStateV1 = "deleted"
)
```
```
type VPCStatusV1 string
VPCStatusV1 represents VPC statuses

```
```
const (
VPCStatusProvisioning    VPCStatusV1 = "provisioning"
VPCStatusCreating        VPCStatusV1 = "creating"
VPCStatusCompleted       VPCStatusV1 = "completed"
VPCStatusDeletingPending VPCStatusV1 = "deleting_pending"
VPCStatusDeleting        VPCStatusV1 = "deleting"
VPCStatusDeleted         VPCStatusV1 = "deleted"
VPCStatusError           VPCStatusV1 = "error"
)

```

