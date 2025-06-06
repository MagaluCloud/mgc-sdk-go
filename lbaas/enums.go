package lbaas

type AclActionType string

const (
	AclActionTypeAllow           AclActionType = "ALLOW"
	AclActionTypeDeny            AclActionType = "DENY"
	AclActionTypeDenyUnspecified AclActionType = "DENY_UNSPECIFIED"
)

type AclEtherType string

const (
	AclEtherTypeIPv4 AclEtherType = "IPv4"
	AclEtherTypeIPv6 AclEtherType = "IPv6"
)

type BackendBalanceAlgorithm string

const (
	BackendBalanceAlgorithmRoundRobin BackendBalanceAlgorithm = "round_robin"
)

type BackendType string

const (
	BackendTypeInstance BackendType = "instance"
	BackendTypeRaw      BackendType = "raw"
)

type HealthCheckProtocol string

const (
	HealthCheckProtocolTCP  HealthCheckProtocol = "tcp"
	HealthCheckProtocolHTTP HealthCheckProtocol = "http"
)

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

type LoadBalancerVisibility string

const (
	LoadBalancerVisibilityInternal LoadBalancerVisibility = "internal"
	LoadBalancerVisibilityExternal LoadBalancerVisibility = "external"
)

type AclProtocol string

const (
	AclProtocolTCP AclProtocol = "tcp"
	AclProtocolTLS AclProtocol = "tls"
)

type ListenerProtocol string

const (
	ListenerProtocolTCP ListenerProtocol = "tcp"
	ListenerProtocolTLS ListenerProtocol = "tls"
)
