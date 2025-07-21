# Lbaas

Package lbaas provides a client for interacting with the Magalu Cloud Load
Balancer as a Service (LBaaS) API. This package allows you to manage network
load balancers, listeners, backends, health checks, certificates, and ACLs.
```
const (
DefaultBasePath = "/load-balancer"
)


```
```
type AclActionType string
AclActionType represents the action type for ACL rules

```
```
const (
AclActionTypeAllow           AclActionType = "ALLOW"
AclActionTypeDeny            AclActionType = "DENY"
AclActionTypeDenyUnspecified AclActionType = "DENY_UNSPECIFIED"
)
```
```
type AclEtherType string
AclEtherType represents the ethernet type for ACL rules

```
```
const (
AclEtherTypeIPv4 AclEtherType = "IPv4"
AclEtherTypeIPv6 AclEtherType = "IPv6"
)
```
```
type AclProtocol string
AclProtocol represents the protocol for ACL rules

```
```
const (
AclProtocolTCP AclProtocol = "tcp"
AclProtocolTLS AclProtocol = "tls"
)
```
```
type BackendBalanceAlgorithm string
BackendBalanceAlgorithm represents the load balancing algorithm for backends

```
```
const (
BackendBalanceAlgorithmRoundRobin BackendBalanceAlgorithm = "round_robin"
)
```
```
type BackendType string
BackendType represents the type of backend targets

```
```
const (
BackendTypeInstance BackendType = "instance"
BackendTypeRaw      BackendType = "raw"
)
```
```
type CreateNetworkACLRequest struct {
Name           *string       `json:"name,omitempty"`
Ethertype      AclEtherType  `json:"ethertype"`
LoadBalancerID string        `json:"load_balancer_id"`
Action         AclActionType `json:"action"`
Protocol       AclProtocol   `json:"protocol"`
RemoteIPPrefix string        `json:"remote_ip_prefix"`
}
CreateNetworkACLRequest represents the request payload for creating a
network ACL rule

```
```
type CreateNetworkBackendRequest struct {
LoadBalancerID   string                        `json:"-"`
Name             string                        `json:"name"`
Description      *string                       `json:"description,omitempty"`
BalanceAlgorithm BackendBalanceAlgorithm       `json:"balance_algorithm"`
TargetsType      BackendType                   `json:"targets_type"`
Targets          *TargetsRawOrInstancesRequest `json:"targets,omitempty"`
HealthCheckID    *string                       `json:"health_check_id,omitempty"`
}
CreateNetworkBackendRequest represents the request payload for creating a
network backend

```
```
type CreateNetworkBackendTargetRequest struct {
LoadBalancerID   string      `json:"-"`
NetworkBackendID string      `json:"-"`
TargetsID        []string    `json:"targets_id"`
TargetsType      BackendType `json:"targets_type"`
}
CreateNetworkBackendTargetRequest represents the request payload for
creating backend targets

```
```
type CreateNetworkCertificateRequest struct {
LoadBalancerID string  `json:"-"`
Name           string  `json:"name"`
Description    *string `json:"description,omitempty"`
Certificate    string  `json:"certificate"`
PrivateKey     string  `json:"private_key"`
}
CreateNetworkCertificateRequest represents the request payload for creating
a network TLS certificate

```
```
type CreateNetworkHealthCheckRequest struct {
LoadBalancerID          string              `json:"-"`
Name                    string              `json:"name"`
Description             *string             `json:"description,omitempty"`
Protocol                HealthCheckProtocol `json:"protocol"`
Path                    *string             `json:"path,omitempty"`
Port                    int                 `json:"port"`
HealthyStatusCode       *int                `json:"healthy_status_code,omitempty"`
IntervalSeconds         *int                `json:"interval_seconds,omitempty"`
TimeoutSeconds          *int                `json:"timeout_seconds,omitempty"`
InitialDelaySeconds     *int                `json:"initial_delay_seconds,omitempty"`
HealthyThresholdCount   *int                `json:"healthy_threshold_count,omitempty"`
UnhealthyThresholdCount *int                `json:"unhealthy_threshold_count,omitempty"`
}
CreateNetworkHealthCheckRequest represents the request payload for creating
a network health check

```
```
type CreateNetworkListenerRequest struct {
LoadBalancerID   string           `json:"-"`
BackendID        string           `json:"-"`
TLSCertificateID *string          `json:"tls_certificate_id,omitempty"`
Name             string           `json:"name"`
Description      *string          `json:"description,omitempty"`
Protocol         ListenerProtocol `json:"protocol"`
Port             int              `json:"port"`
}
CreateNetworkListenerRequest represents the request payload for creating a
network listener

```
```
type CreateNetworkLoadBalancerRequest struct {
Name            string                         `json:"name"`
Description     *string                        `json:"description,omitempty"`
Type            *string                        `json:"type,omitempty"`
Visibility      LoadBalancerVisibility         `json:"visibility"`
Listeners       []NetworkListenerRequest       `json:"listeners"`
Backends        []NetworkBackendRequest        `json:"backends"`
HealthChecks    []NetworkHealthCheckRequest    `json:"health_checks,omitempty"`
TLSCertificates []NetworkTLSCertificateRequest `json:"tls_certificates,omitempty"`
ACLs            []NetworkAclRequest            `json:"acls,omitempty"`
VPCID           string                         `json:"vpc_id"`
SubnetPoolID    *string                        `json:"subnet_pool_id,omitempty"`
PublicIPID      *string                        `json:"public_ip_id,omitempty"`
PanicThreshold  *int                           `json:"panic_threshold,omitempty"`
}
CreateNetworkLoadBalancerRequest represents the request payload for creating
a load balancer

```
```
type DeleteNetworkACLRequest struct {
LoadBalancerID string `json:"load_balancer_id"`
ID             string `json:"id"`
}
DeleteNetworkACLRequest represents the request payload for deleting a
network ACL rule

```
```
type DeleteNetworkBackendRequest struct {
LoadBalancerID string `json:"-"`
BackendID      string `json:"-"`
}
DeleteNetworkBackendRequest represents the request payload for deleting a
network backend

```
```
type DeleteNetworkBackendTargetRequest struct {
LoadBalancerID   string `json:"-"`
NetworkBackendID string `json:"-"`
TargetID         string `json:"-"`
}
DeleteNetworkBackendTargetRequest represents the request payload for
deleting a backend target

```
```
type DeleteNetworkCertificateRequest struct {
LoadBalancerID   string `json:"-"`
TLSCertificateID string `json:"-"`
}
DeleteNetworkCertificateRequest represents the request payload for deleting
a network TLS certificate

```
```
type DeleteNetworkHealthCheckRequest struct {
LoadBalancerID string `json:"-"`
HealthCheckID  string `json:"-"`
}
DeleteNetworkHealthCheckRequest represents the request payload for deleting
a network health check

```
```
type DeleteNetworkListenerRequest struct {
LoadBalancerID string `json:"-"`
ListenerID     string `json:"-"`
}
DeleteNetworkListenerRequest represents the request payload for deleting a
network listener

```
```
type DeleteNetworkLoadBalancerRequest struct {
LoadBalancerID string `json:"-"`
DeletePublicIP *bool  `json:"-"`
}
DeleteNetworkLoadBalancerRequest represents the request payload for deleting
a load balancer

```
```
type GetNetworkACLRequest struct {
LoadBalancerID string `json:"-"`
NetworkACLID   string `json:"-"`
}
GetNetworkACLRequest represents the request payload for getting a network
ACL rule

```
```
type GetNetworkBackendRequest struct {
LoadBalancerID string `json:"-"`
BackendID      string `json:"-"`
}
GetNetworkBackendRequest represents the request payload for getting a
network backend

```
```
type GetNetworkCertificateRequest struct {
LoadBalancerID   string `json:"-"`
TLSCertificateID string `json:"-"`
}
GetNetworkCertificateRequest represents the request payload for getting a
network TLS certificate

```
```
type GetNetworkHealthCheckRequest struct {
LoadBalancerID string `json:"-"`
HealthCheckID  string `json:"-"`
}
GetNetworkHealthCheckRequest represents the request payload for getting a
network health check

```
```
type GetNetworkListenerRequest struct {
LoadBalancerID string `json:"-"`
ListenerID     string `json:"-"`
}
GetNetworkListenerRequest represents the request payload for getting a
network listener

```
```
type GetNetworkLoadBalancerRequest struct {
LoadBalancerID string `json:"-"`
}
GetNetworkLoadBalancerRequest represents the request payload for getting a
load balancer

```
```
type HealthCheckProtocol string
HealthCheckProtocol represents the protocol for health checks

```
```
const (
HealthCheckProtocolTCP  HealthCheckProtocol = "tcp"
HealthCheckProtocolHTTP HealthCheckProtocol = "http"
)
```
```
type LbaasClient struct {
*client.CoreClient
}
LbaasClient represents a client for the Load Balancer as a Service

```
```
func New(core *client.CoreClient) *LbaasClient
New creates a new LbaasClient instance with the provided core client

```
```
func (c *LbaasClient) NetworkACLs() NetworkACLService
NetworkACLs returns a service for managing network ACLs

```
```
func (c *LbaasClient) NetworkBackends() NetworkBackendService
NetworkBackends returns a service for managing network backends

```
```
func (c *LbaasClient) NetworkCertificates() NetworkCertificateService
NetworkCertificates returns a service for managing network certificates

```
```
func (c *LbaasClient) NetworkHealthChecks() NetworkHealthCheckService
NetworkHealthChecks returns a service for managing network health checks

```
```
func (c *LbaasClient) NetworkListeners() NetworkListenerService
NetworkListeners returns a service for managing network listeners

```
```
func (c *LbaasClient) NetworkLoadBalancers() NetworkLoadBalancerService
NetworkLoadBalancers returns a service for managing network load balancers

```
```
type ListNetworkACLRequest struct {
LoadBalancerID string `json:"-"`
}
ListNetworkACLRequest represents the request payload for listing network ACL
rules

```
```
type ListNetworkBackendRequest struct {
LoadBalancerID string `json:"-"`
}
ListNetworkBackendRequest represents the request payload for listing network
backends

```
```
type ListNetworkCertificateRequest struct {
LoadBalancerID string  `json:"-"`
Offset         *int    `json:"-"`
Limit          *int    `json:"-"`
Sort           *string `json:"-"`
}
ListNetworkCertificateRequest represents the request payload for listing
network TLS certificates

```
```
type ListNetworkHealthCheckRequest struct {
LoadBalancerID string  `json:"-"`
Offset         *int    `json:"-"`
Limit          *int    `json:"-"`
Sort           *string `json:"-"`
}
ListNetworkHealthCheckRequest represents the request payload for listing
network health checks

```
```
type ListNetworkListenerRequest struct {
LoadBalancerID string  `json:"-"`
Offset         *int    `json:"-"`
Limit          *int    `json:"-"`
Sort           *string `json:"-"`
}
ListNetworkListenerRequest represents the request payload for listing
network listeners

```
```
type ListNetworkLoadBalancerRequest struct {
Offset *int    `json:"-"`
Limit  *int    `json:"-"`
Sort   *string `json:"-"`
}
ListNetworkLoadBalancerRequest represents the request payload for listing
load balancers

```
```
type ListenerProtocol string
ListenerProtocol represents the protocol for listeners

```
```
const (
ListenerProtocolTCP ListenerProtocol = "tcp"
ListenerProtocolTLS ListenerProtocol = "tls"
)
```
```
type LoadBalancerStatus string
LoadBalancerStatus represents the status of a load balancer

```
```
const (
LoadBalancerStatusCreating LoadBalancerStatus = "creating"
LoadBalancerStatusUpdating LoadBalancerStatus = "updating"
LoadBalancerStatusDeleting LoadBalancerStatus = "deleting"
LoadBalancerStatusRunning  LoadBalancerStatus = "running"
LoadBalancerStatusFailed   LoadBalancerStatus = "failed"
LoadBalancerStatusCanceled LoadBalancerStatus = "canceled"
LoadBalancerStatusDeleted  LoadBalancerStatus = "deleted"
LoadBalancerStatusInactive LoadBalancerStatus = "inactive"
)
```
```
type LoadBalancerVisibility string
LoadBalancerVisibility represents the visibility of a load balancer

```
```
const (
LoadBalancerVisibilityInternal LoadBalancerVisibility = "internal"
LoadBalancerVisibilityExternal LoadBalancerVisibility = "external"
)
```
```
type NetworkACLService interface {
Create(ctx context.Context, req CreateNetworkACLRequest) (string, error)
Delete(ctx context.Context, req DeleteNetworkACLRequest) error
}
NetworkACLService provides methods for managing network ACL rules

```
```
type NetworkAclRequest struct {
Name           *string       `json:"name,omitempty"`
Ethertype      AclEtherType  `json:"ethertype"`
Protocol       AclProtocol   `json:"protocol"`
RemoteIPPrefix string        `json:"remote_ip_prefix"`
Action         AclActionType `json:"action"`
}
NetworkAclRequest represents an ACL rule configuration for load balancer
creation

```
```
type NetworkAclResponse struct {
ID             string       `json:"id"`
Name           *string      `json:"name,omitempty"`
Ethertype      AclEtherType `json:"ethertype"`
Protocol       AclProtocol  `json:"protocol"`
RemoteIPPrefix string       `json:"remote_ip_prefix"`
Action         string       `json:"action"`
}
NetworkAclResponse represents an ACL rule response

```
```
type NetworkBackendInstanceRequest struct {
NicID string `json:"nic_id"`
Port  int    `json:"port"`
}
NetworkBackendInstanceRequest represents an instance-based backend target

```
```
type NetworkBackendInstanceResponse struct {
ID        string  `json:"id"`
IPAddress *string `json:"ip_address,omitempty"`
NicID     string  `json:"nic_id,omitempty"`
Port      int     `json:"port"`
CreatedAt string  `json:"created_at"`
UpdatedAt string  `json:"updated_at"`
}
NetworkBackendInstanceResponse represents an instance-based backend target
response

```
```
type NetworkBackendInstanceUpdateRequest struct {
NicID string `json:"nic_id"`
Port  int    `json:"port"`
}
NetworkBackendInstanceUpdateRequest represents an instance-based backend
target for updates

```
```
type NetworkBackendRawTargetRequest struct {
IPAddress string `json:"ip_address"`
Port      int    `json:"port"`
}
NetworkBackendRawTargetRequest represents a raw IP/port backend target

```
```
type NetworkBackendRawTargetResponse struct {
ID        string  `json:"id"`
IPAddress *string `json:"ip_address,omitempty"`
Port      int     `json:"port"`
CreatedAt string  `json:"created_at"`
UpdatedAt string  `json:"updated_at"`
}
NetworkBackendRawTargetResponse represents a raw IP/port backend target
response

```
```
type NetworkBackendRawTargetUpdateRequest struct {
IPAddress string `json:"ip_address"`
Port      int    `json:"port"`
}
NetworkBackendRawTargetUpdateRequest represents a raw IP/port backend target
for updates

```
```
type NetworkBackendRequest struct {
HealthCheckName  *string                       `json:"health_check_name,omitempty"`
Name             string                        `json:"name"`
Description      *string                       `json:"description,omitempty"`
BalanceAlgorithm BackendBalanceAlgorithm       `json:"balance_algorithm"`
TargetsType      BackendType                   `json:"targets_type"`
Targets          *TargetsRawOrInstancesRequest `json:"targets,omitempty"`
}
NetworkBackendRequest represents a backend configuration for load balancer
creation

```
```
type NetworkBackendResponse struct {
ID               string                  `json:"id"`
HealthCheckID    *string                 `json:"health_check_id,omitempty"`
Name             string                  `json:"name"`
Description      *string                 `json:"description,omitempty"`
BalanceAlgorithm BackendBalanceAlgorithm `json:"balance_algorithm"`
TargetsType      BackendType             `json:"targets_type"`
Targets          interface{}             `json:"targets"`
CreatedAt        string                  `json:"created_at"`
UpdatedAt        string                  `json:"updated_at"`
}
NetworkBackendResponse represents a network backend response

```
```
type NetworkBackendService interface {
Create(ctx context.Context, req CreateNetworkBackendRequest) (string, error)
Delete(ctx context.Context, req DeleteNetworkBackendRequest) error
Get(ctx context.Context, req GetNetworkBackendRequest) (*NetworkBackendResponse, error)
List(ctx context.Context, req ListNetworkBackendRequest) ([]NetworkBackendResponse, error)
Update(ctx context.Context, req UpdateNetworkBackendRequest) error
Targets() *networkBackendTargetService
}
NetworkBackendService provides methods for managing network backends

```
```
type NetworkBackendTargetService interface {
Create(ctx context.Context, req CreateNetworkBackendTargetRequest) (string, error)
Delete(ctx context.Context, req DeleteNetworkBackendTargetRequest) error
}
NetworkBackendTargetService provides methods for managing backend targets

```
```
type NetworkBackendUpdateRequest struct {
ID            string                              `json:"id"`
HealthCheckID *string                             `json:"health_check_id,omitempty"`
TargetsType   BackendType                         `json:"targets_type"`
Targets       *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
}
NetworkBackendUpdateRequest represents a backend update configuration

```
```
type NetworkCertificateService interface {
Create(ctx context.Context, req CreateNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error)
Delete(ctx context.Context, req DeleteNetworkCertificateRequest) error
Get(ctx context.Context, req GetNetworkCertificateRequest) (*NetworkTLSCertificateResponse, error)
List(ctx context.Context, req ListNetworkCertificateRequest) ([]NetworkTLSCertificateResponse, error)
Update(ctx context.Context, req UpdateNetworkCertificateRequest) error
}
NetworkCertificateService provides methods for managing network TLS
certificates

```
```
type NetworkHealthCheckRequest struct {
Name                    string              `json:"name"`
Description             *string             `json:"description,omitempty"`
Protocol                HealthCheckProtocol `json:"protocol"`
Path                    *string             `json:"path,omitempty"`
Port                    int                 `json:"port"`
HealthyStatusCode       *int                `json:"healthy_status_code,omitempty"`
IntervalSeconds         *int                `json:"interval_seconds,omitempty"`
TimeoutSeconds          *int                `json:"timeout_seconds,omitempty"`
InitialDelaySeconds     *int                `json:"initial_delay_seconds,omitempty"`
HealthyThresholdCount   *int                `json:"healthy_threshold_count,omitempty"`
UnhealthyThresholdCount *int                `json:"unhealthy_threshold_count,omitempty"`
}
NetworkHealthCheckRequest represents a health check configuration for load
balancer creation

```
```
type NetworkHealthCheckResponse struct {
ID                      string              `json:"id"`
Name                    string              `json:"name"`
Description             *string             `json:"description,omitempty"`
Protocol                HealthCheckProtocol `json:"protocol"`
Path                    *string             `json:"path,omitempty"`
Port                    int                 `json:"port"`
HealthyStatusCode       int                 `json:"healthy_status_code"`
IntervalSeconds         int                 `json:"interval_seconds"`
TimeoutSeconds          int                 `json:"timeout_seconds"`
InitialDelaySeconds     int                 `json:"initial_delay_seconds"`
HealthyThresholdCount   int                 `json:"healthy_threshold_count"`
UnhealthyThresholdCount int                 `json:"unhealthy_threshold_count"`
CreatedAt               string              `json:"created_at"`
UpdatedAt               string              `json:"updated_at"`
}
NetworkHealthCheckResponse represents a network health check response

```
```
type NetworkHealthCheckService interface {
Create(ctx context.Context, req CreateNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error)
Delete(ctx context.Context, req DeleteNetworkHealthCheckRequest) error
Get(ctx context.Context, req GetNetworkHealthCheckRequest) (*NetworkHealthCheckResponse, error)
List(ctx context.Context, req ListNetworkHealthCheckRequest) ([]NetworkHealthCheckResponse, error)
Update(ctx context.Context, req UpdateNetworkHealthCheckRequest) error
}
NetworkHealthCheckService provides methods for managing network health
checks

```
```
type NetworkHealthCheckUpdateRequest struct {
ID                      string              `json:"id"`
Protocol                HealthCheckProtocol `json:"protocol"`
Path                    *string             `json:"path,omitempty"`
Port                    int                 `json:"port"`
HealthyStatusCode       *int                `json:"healthy_status_code,omitempty"`
IntervalSeconds         *int                `json:"interval_seconds,omitempty"`
TimeoutSeconds          *int                `json:"timeout_seconds,omitempty"`
InitialDelaySeconds     *int                `json:"initial_delay_seconds,omitempty"`
HealthyThresholdCount   *int                `json:"healthy_threshold_count,omitempty"`
UnhealthyThresholdCount *int                `json:"unhealthy_threshold_count,omitempty"`
}
NetworkHealthCheckUpdateRequest represents a health check update
configuration

```
```
type NetworkLBPaginatedResponse struct {
Results []NetworkLoadBalancerResponse `json:"results"`
}
NetworkLBPaginatedResponse represents a paginated load balancer response

```
```
type NetworkListenerRequest struct {
TLSCertificateName *string          `json:"tls_certificate_name,omitempty"`
Name               string           `json:"name"`
Description        *string          `json:"description,omitempty"`
BackendName        string           `json:"backend_name"`
Protocol           ListenerProtocol `json:"protocol"`
Port               int              `json:"port"`
}
NetworkListenerRequest represents a listener configuration for load balancer
creation

```
```
type NetworkListenerResponse struct {
ID               string           `json:"id"`
TLSCertificateID *string          `json:"tls_certificate_id,omitempty"`
BackendID        string           `json:"backend_id"`
Name             string           `json:"name"`
Description      *string          `json:"description,omitempty"`
Protocol         ListenerProtocol `json:"protocol"`
Port             int              `json:"port"`
CreatedAt        string           `json:"created_at"`
UpdatedAt        string           `json:"updated_at"`
}
NetworkListenerResponse represents a network listener response

```
```
type NetworkListenerService interface {
Create(ctx context.Context, req CreateNetworkListenerRequest) (*NetworkListenerResponse, error)
Delete(ctx context.Context, req DeleteNetworkListenerRequest) error
Get(ctx context.Context, req GetNetworkListenerRequest) (*NetworkListenerResponse, error)
List(ctx context.Context, req ListNetworkListenerRequest) ([]NetworkListenerResponse, error)
Update(ctx context.Context, req UpdateNetworkListenerRequest) error
}
NetworkListenerService provides methods for managing network listeners

```
```
type NetworkLoadBalancerResponse struct {
ID                  string                          `json:"id"`
Name                string                          `json:"name"`
ProjectType         *string                         `json:"project_type,omitempty"`
Description         *string                         `json:"description,omitempty"`
Type                string                          `json:"type"`
Visibility          LoadBalancerVisibility          `json:"visibility"`
Status              string                          `json:"status"`
Listeners           []NetworkListenerResponse       `json:"listeners"`
Backends            []NetworkBackendResponse        `json:"backends"`
HealthChecks        []NetworkHealthCheckResponse    `json:"health_checks"`
PublicIPs           []NetworkPublicIPResponse       `json:"public_ips"`
TLSCertificates     []NetworkTLSCertificateResponse `json:"tls_certificates"`
ACLs                []NetworkAclResponse            `json:"acls"`
IPAddress           *string                         `json:"ip_address,omitempty"`
Port                *string                         `json:"port,omitempty"`
VPCID               string                          `json:"vpc_id"`
SubnetPoolID        *string                         `json:"subnet_pool_id,omitempty"`
CreatedAt           string                          `json:"created_at"`
UpdatedAt           string                          `json:"updated_at"`
LastOperationStatus *string                         `json:"last_operation_status,omitempty"`
}
NetworkLoadBalancerResponse represents a load balancer response

```
```
type NetworkLoadBalancerService interface {
Create(ctx context.Context, req CreateNetworkLoadBalancerRequest) (string, error)
Delete(ctx context.Context, req DeleteNetworkLoadBalancerRequest) error
Get(ctx context.Context, req GetNetworkLoadBalancerRequest) (*NetworkLoadBalancerResponse, error)
List(ctx context.Context, req ListNetworkLoadBalancerRequest) ([]NetworkLoadBalancerResponse, error)
Update(ctx context.Context, req UpdateNetworkLoadBalancerRequest) error
}
NetworkLoadBalancerService provides methods for managing network load
balancers

```
```
type NetworkPaginatedBackendResponse struct {
Meta    interface{}              `json:"meta"`
Results []NetworkBackendResponse `json:"results"`
}
NetworkPaginatedBackendResponse represents a paginated backend response

```
```
type NetworkPaginatedHealthCheckResponse struct {
Meta    interface{}                  `json:"meta"`
Results []NetworkHealthCheckResponse `json:"results"`
}
NetworkPaginatedHealthCheckResponse represents a paginated health check
response

```
```
type NetworkPaginatedListenerResponse struct {
Meta    interface{}               `json:"meta"`
Results []NetworkListenerResponse `json:"results"`
}
NetworkPaginatedListenerResponse represents a paginated listener response

```
```
type NetworkPaginatedTLSCertificateResponse struct {
Meta    interface{}                     `json:"meta"`
Results []NetworkTLSCertificateResponse `json:"results"`
}
NetworkPaginatedTLSCertificateResponse represents a paginated TLS
certificate response

```
```
type NetworkPublicIPResponse struct {
ID         string  `json:"id"`
IPAddress  *string `json:"ip_address,omitempty"`
ExternalID string  `json:"external_id"`
}
NetworkPublicIPResponse represents a public IP response

```
```
type NetworkTLSCertificateRequest struct {
Name        string  `json:"name"`
Description *string `json:"description,omitempty"`
Certificate string  `json:"certificate"`
PrivateKey  string  `json:"private_key"`
}
NetworkTLSCertificateRequest represents a TLS certificate configuration for
load balancer creation

```
```
type NetworkTLSCertificateResponse struct {
ID             string  `json:"id"`
Name           string  `json:"name"`
ExpirationDate *string `json:"expiration_date,omitempty"`
Description    *string `json:"description,omitempty"`
CreatedAt      string  `json:"created_at"`
UpdatedAt      string  `json:"updated_at"`
}
NetworkTLSCertificateResponse represents a network TLS certificate response

```
```
type NetworkTLSCertificateUpdateRequest struct {
ID          string `json:"id"`
Certificate string `json:"certificate"`
PrivateKey  string `json:"private_key"`
}
NetworkTLSCertificateUpdateRequest represents a TLS certificate update
configuration

```
```
type TargetsRawOrInstancesRequest struct {
TargetsInstances []NetworkBackendInstanceRequest  `json:"-"`
TargetsRaw       []NetworkBackendRawTargetRequest `json:"-"`
}
TargetsRawOrInstancesRequest represents backend targets that can be either
instances or raw IPs

```
```
func (t *TargetsRawOrInstancesRequest) MarshalJSON() ([]byte, error)
MarshalJSON implements custom JSON marshaling for
TargetsRawOrInstancesRequest Prioritizes TargetsInstances over TargetsRaw
when marshaling

```
```
func (t *TargetsRawOrInstancesRequest) UnmarshalJSON(data []byte) error
UnmarshalJSON implements custom JSON unmarshaling for
TargetsRawOrInstancesRequest This allows flexible handling of different
target types in the JSON payload

```
```
type TargetsRawOrInstancesUpdateRequest struct {
TargetsInstances []NetworkBackendInstanceUpdateRequest  `json:"-"`
TargetsRaw       []NetworkBackendRawTargetUpdateRequest `json:"-"`
}
TargetsRawOrInstancesUpdateRequest represents backend targets for updates

```
```
func (t *TargetsRawOrInstancesUpdateRequest) MarshalJSON() ([]byte, error)
MarshalJSON implements custom JSON marshaling for
TargetsRawOrInstancesUpdateRequest Prioritizes TargetsInstances over
TargetsRaw when marshaling

```
```
func (t *TargetsRawOrInstancesUpdateRequest) UnmarshalJSON(data []byte) error
UnmarshalJSON implements custom JSON unmarshaling for
TargetsRawOrInstancesUpdateRequest This allows flexible handling of
different target types in the JSON payload

```
```
type UpdateNetworkBackendRequest struct {
LoadBalancerID   string                              `json:"-"`
BackendID        string                              `json:"-"`
Name             *string                             `json:"name,omitempty"`
Description      *string                             `json:"description,omitempty"`
BalanceAlgorithm *BackendBalanceAlgorithm            `json:"balance_algorithm,omitempty"`
TargetsType      *BackendType                        `json:"targets_type,omitempty"`
Targets          *TargetsRawOrInstancesUpdateRequest `json:"targets,omitempty"`
TargetsInstances *[]NetworkBackendInstanceRequest    `json:"targets_instances,omitempty"`
TargetsRaw       *[]NetworkBackendRawTargetRequest   `json:"targets_raw,omitempty"`
HealthCheckID    *string                             `json:"health_check_id,omitempty"`
}
UpdateNetworkBackendRequest represents the request payload for updating a
network backend

```
```
type UpdateNetworkCertificateRequest struct {
LoadBalancerID   string `json:"-"`
TLSCertificateID string `json:"-"`
Certificate      string `json:"certificate"`
PrivateKey       string `json:"private_key"`
}
UpdateNetworkCertificateRequest represents the request payload for updating
a network TLS certificate

```
```
type UpdateNetworkHealthCheckRequest struct {
LoadBalancerID          string              `json:"-"`
HealthCheckID           string              `json:"-"`
Protocol                HealthCheckProtocol `json:"protocol"`
Path                    *string             `json:"path,omitempty"`
Port                    int                 `json:"port"`
HealthyStatusCode       *int                `json:"healthy_status_code,omitempty"`
IntervalSeconds         *int                `json:"interval_seconds,omitempty"`
TimeoutSeconds          *int                `json:"timeout_seconds,omitempty"`
InitialDelaySeconds     *int                `json:"initial_delay_seconds,omitempty"`
HealthyThresholdCount   *int                `json:"healthy_threshold_count,omitempty"`
UnhealthyThresholdCount *int                `json:"unhealthy_threshold_count,omitempty"`
}
UpdateNetworkHealthCheckRequest represents the request payload for updating
a network health check

```
```
type UpdateNetworkListenerRequest struct {
LoadBalancerID   string  `json:"-"`
ListenerID       string  `json:"-"`
TLSCertificateID *string `json:"tls_certificate_id,omitempty"`
}
UpdateNetworkListenerRequest represents the request payload for updating a
network listener

```
```
type UpdateNetworkLoadBalancerRequest struct {
LoadBalancerID  string                               `json:"-"`
Name            *string                              `json:"name,omitempty"`
Description     *string                              `json:"description,omitempty"`
Backends        []NetworkBackendUpdateRequest        `json:"backends,omitempty"`
HealthChecks    []NetworkHealthCheckUpdateRequest    `json:"health_checks,omitempty"`
TLSCertificates []NetworkTLSCertificateUpdateRequest `json:"tls_certificates,omitempty"`
PanicThreshold  *int                                 `json:"panic_threshold,omitempty"`
}
UpdateNetworkLoadBalancerRequest represents the request payload for updating
a load balancer


```

