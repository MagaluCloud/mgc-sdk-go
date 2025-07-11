LoadBalancer
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 12

Package Documentation
-------------------

.. code-block:: go

   package lbaas // import "github.com/MagaluCloud/mgc-sdk-go/lbaas"
   
   Package lbaas provides a client for interacting with the Magalu Cloud Load
   Balancer as a Service (LBaaS) API. This package allows you to manage network
   load balancers, listeners, backends, health checks, certificates, and ACLs.
   
   const DefaultBasePath = "/load-balancer"
   type AclActionType string
       const AclActionTypeAllow AclActionType = "ALLOW" ...
   type AclEtherType string
       const AclEtherTypeIPv4 AclEtherType = "IPv4" ...
   type AclProtocol string
       const AclProtocolTCP AclProtocol = "tcp" ...
   type BackendBalanceAlgorithm string
       const BackendBalanceAlgorithmRoundRobin BackendBalanceAlgorithm = "round_robin"
   type BackendType string
       const BackendTypeInstance BackendType = "instance" ...
   type CreateNetworkACLRequest struct{ ... }
   type CreateNetworkBackendRequest struct{ ... }
   type CreateNetworkBackendTargetRequest struct{ ... }


Types
-----

- :type:`AclActionType`
- :type:`AclEtherType`
- :type:`AclProtocol`
- :type:`BackendBalanceAlgorithm`
- :type:`BackendType`
- :type:`CreateNetworkACLRequest`
- :type:`CreateNetworkBackendRequest`
- :type:`CreateNetworkBackendTargetRequest`
- :type:`CreateNetworkCertificateRequest`
- :type:`CreateNetworkHealthCheckRequest`
- :type:`CreateNetworkListenerRequest`
- :type:`CreateNetworkLoadBalancerRequest`
- :type:`DeleteNetworkACLRequest`
- :type:`DeleteNetworkBackendRequest`
- :type:`DeleteNetworkBackendTargetRequest`
- :type:`DeleteNetworkCertificateRequest`
- :type:`DeleteNetworkHealthCheckRequest`
- :type:`DeleteNetworkListenerRequest`
- :type:`DeleteNetworkLoadBalancerRequest`
- :type:`GetNetworkACLRequest`
- :type:`GetNetworkBackendRequest`
- :type:`GetNetworkCertificateRequest`
- :type:`GetNetworkHealthCheckRequest`
- :type:`GetNetworkListenerRequest`
- :type:`GetNetworkLoadBalancerRequest`
- :type:`HealthCheckProtocol`
- :type:`LbaasClient`
- :type:`ListNetworkACLRequest`
- :type:`ListNetworkBackendRequest`
- :type:`ListNetworkCertificateRequest`
- :type:`ListNetworkHealthCheckRequest`
- :type:`ListNetworkListenerRequest`
- :type:`ListNetworkLoadBalancerRequest`
- :type:`ListenerProtocol`
- :type:`LoadBalancerStatus`
- :type:`LoadBalancerVisibility`
- :type:`NetworkACLService`
- :type:`NetworkAclRequest`
- :type:`NetworkAclResponse`
- :type:`NetworkBackendInstanceRequest`
- :type:`NetworkBackendInstanceResponse`
- :type:`NetworkBackendInstanceUpdateRequest`
- :type:`NetworkBackendRawTargetRequest`
- :type:`NetworkBackendRawTargetResponse`
- :type:`NetworkBackendRawTargetUpdateRequest`
- :type:`NetworkBackendRequest`
- :type:`NetworkBackendResponse`
- :type:`NetworkBackendService`
- :type:`NetworkBackendTargetService`
- :type:`NetworkBackendUpdateRequest`
- :type:`NetworkCertificateService`
- :type:`NetworkHealthCheckRequest`
- :type:`NetworkHealthCheckResponse`
- :type:`NetworkHealthCheckService`
- :type:`NetworkHealthCheckUpdateRequest`
- :type:`NetworkLBPaginatedResponse`
- :type:`NetworkListenerRequest`
- :type:`NetworkListenerResponse`
- :type:`NetworkListenerService`
- :type:`NetworkLoadBalancerResponse`
- :type:`NetworkLoadBalancerService`
- :type:`NetworkPaginatedBackendResponse`
- :type:`NetworkPaginatedHealthCheckResponse`
- :type:`NetworkPaginatedListenerResponse`
- :type:`NetworkPaginatedTLSCertificateResponse`
- :type:`NetworkPublicIPResponse`
- :type:`NetworkTLSCertificateRequest`
- :type:`NetworkTLSCertificateResponse`
- :type:`NetworkTLSCertificateUpdateRequest`
- :type:`TargetsRawOrInstancesRequest`
- :type:`TargetsRawOrInstancesUpdateRequest`
- :type:`UpdateNetworkBackendRequest`
- :type:`UpdateNetworkCertificateRequest`
- :type:`UpdateNetworkHealthCheckRequest`
- :type:`UpdateNetworkListenerRequest`
- :type:`UpdateNetworkLoadBalancerRequest`

Constants
---------

- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/lbaas"

   // Use the LoadBalancer package
   // See the examples directory for complete examples

