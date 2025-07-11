package lbaas

// AclActionType represents the action type for ACL rules
type AclActionType string

const (
	// AclActionTypeAllow allows traffic that matches the rule
	AclActionTypeAllow AclActionType = "ALLOW"
	// AclActionTypeDeny denies traffic that matches the rule
	AclActionTypeDeny AclActionType = "DENY"
	// AclActionTypeDenyUnspecified denies unspecified traffic
	AclActionTypeDenyUnspecified AclActionType = "DENY_UNSPECIFIED"
)

// AclEtherType represents the ethernet type for ACL rules
type AclEtherType string

const (
	// AclEtherTypeIPv4 represents IPv4 traffic
	AclEtherTypeIPv4 AclEtherType = "IPv4"
	// AclEtherTypeIPv6 represents IPv6 traffic
	AclEtherTypeIPv6 AclEtherType = "IPv6"
)

// BackendBalanceAlgorithm represents the load balancing algorithm for backends
type BackendBalanceAlgorithm string

const (
	// BackendBalanceAlgorithmRoundRobin distributes requests in round-robin fashion
	BackendBalanceAlgorithmRoundRobin BackendBalanceAlgorithm = "round_robin"
)

// BackendType represents the type of backend targets
type BackendType string

const (
	// BackendTypeInstance represents instance-based backend targets
	BackendTypeInstance BackendType = "instance"
	// BackendTypeRaw represents raw IP/port backend targets
	BackendTypeRaw BackendType = "raw"
)

// HealthCheckProtocol represents the protocol for health checks
type HealthCheckProtocol string

const (
	// HealthCheckProtocolTCP uses TCP for health checks
	HealthCheckProtocolTCP HealthCheckProtocol = "tcp"
	// HealthCheckProtocolHTTP uses HTTP for health checks
	HealthCheckProtocolHTTP HealthCheckProtocol = "http"
)

// LoadBalancerStatus represents the status of a load balancer
type LoadBalancerStatus string

const (
	// LoadBalancerStatusCreating indicates the load balancer is being created
	LoadBalancerStatusCreating LoadBalancerStatus = "creating"
	// LoadBalancerStatusUpdating indicates the load balancer is being updated
	LoadBalancerStatusUpdating LoadBalancerStatus = "updating"
	// LoadBalancerStatusDeleting indicates the load balancer is being deleted
	LoadBalancerStatusDeleting LoadBalancerStatus = "deleting"
	// LoadBalancerStatusRunning indicates the load balancer is running normally
	LoadBalancerStatusRunning LoadBalancerStatus = "running"
	// LoadBalancerStatusFailed indicates the load balancer creation/update failed
	LoadBalancerStatusFailed LoadBalancerStatus = "failed"
	// LoadBalancerStatusCanceled indicates the load balancer operation was canceled
	LoadBalancerStatusCanceled LoadBalancerStatus = "canceled"
	// LoadBalancerStatusDeleted indicates the load balancer has been deleted
	LoadBalancerStatusDeleted LoadBalancerStatus = "deleted"
	// LoadBalancerStatusInactive indicates the load balancer is inactive
	LoadBalancerStatusInactive LoadBalancerStatus = "inactive"
)

// LoadBalancerVisibility represents the visibility of a load balancer
type LoadBalancerVisibility string

const (
	// LoadBalancerVisibilityInternal indicates internal load balancer visibility
	LoadBalancerVisibilityInternal LoadBalancerVisibility = "internal"
	// LoadBalancerVisibilityExternal indicates external load balancer visibility
	LoadBalancerVisibilityExternal LoadBalancerVisibility = "external"
)

// AclProtocol represents the protocol for ACL rules
type AclProtocol string

const (
	// AclProtocolTCP represents TCP protocol
	AclProtocolTCP AclProtocol = "tcp"
	// AclProtocolTLS represents TLS protocol
	AclProtocolTLS AclProtocol = "tls"
)

// ListenerProtocol represents the protocol for listeners
type ListenerProtocol string

const (
	// ListenerProtocolTCP represents TCP protocol
	ListenerProtocolTCP ListenerProtocol = "tcp"
	// ListenerProtocolTLS represents TLS protocol
	ListenerProtocolTLS ListenerProtocol = "tls"
)
