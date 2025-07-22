package lbaas

// AclActionType represents the action type for ACL rules
type AclActionType string

const (
	AclActionTypeAllow           AclActionType = "ALLOW"
	AclActionTypeDeny            AclActionType = "DENY"
	AclActionTypeDenyUnspecified AclActionType = "DENY_UNSPECIFIED"
)

// AclEtherType represents the ethernet type for ACL rules
type AclEtherType string

const (
	AclEtherTypeIPv4 AclEtherType = "IPv4"
	AclEtherTypeIPv6 AclEtherType = "IPv6"
)

// BackendBalanceAlgorithm represents the load balancing algorithm for backends
type BackendBalanceAlgorithm string

const (
	BackendBalanceAlgorithmRoundRobin BackendBalanceAlgorithm = "round_robin"
)

// BackendType represents the type of backend targets
type BackendType string

const (
	BackendTypeInstance BackendType = "instance"
	BackendTypeRaw      BackendType = "raw"
)

// HealthCheckProtocol represents the protocol for health checks
type HealthCheckProtocol string

const (
	HealthCheckProtocolTCP  HealthCheckProtocol = "tcp"
	HealthCheckProtocolHTTP HealthCheckProtocol = "http"
)

// LoadBalancerStatus represents the status of a load balancer
type LoadBalancerStatus string

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

// LoadBalancerVisibility represents the visibility of a load balancer
type LoadBalancerVisibility string

const (
	LoadBalancerVisibilityInternal LoadBalancerVisibility = "internal"
	LoadBalancerVisibilityExternal LoadBalancerVisibility = "external"
)

// AclProtocol represents the protocol for ACL rules
type AclProtocol string

const (
	AclProtocolTCP AclProtocol = "tcp"
	AclProtocolTLS AclProtocol = "tls"
)

// ListenerProtocol represents the protocol for listeners
type ListenerProtocol string

const (
	ListenerProtocolTCP ListenerProtocol = "tcp"
	ListenerProtocolTLS ListenerProtocol = "tls"
)
