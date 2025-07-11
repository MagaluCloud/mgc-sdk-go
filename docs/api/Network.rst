Network
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 7

Package Documentation
-------------------

.. code-block:: go

   package network // import "github.com/MagaluCloud/mgc-sdk-go/network"
   
   Package network provides a client for interacting with the Magalu Cloud Network
   API. This package allows you to manage VPCs, subnets, ports, security groups,
   rules, public IPs, subnet pools, and NAT gateways.
   
   const SecurityGroupsExpand = "security_groups" ...
   const DefaultBasePath = "/network"
   type BookCIDRRequest struct{ ... }
   type BookCIDRResponse struct{ ... }
   type CreateNatGatewayRequest struct{ ... }
   type CreateSubnetPoolRequest struct{ ... }
   type CreateSubnetPoolResponse struct{ ... }
   type CreateVPCRequest struct{ ... }
   type CreateVPCResponse struct{ ... }
   type DHCPPoolStr struct{ ... }
   type IPAddress struct{ ... }
   type IpAddress struct{ ... }
   type LinkModel struct{ ... }
   type ListOptions struct{ ... }


Types
-----

- :type:`BookCIDRRequest`
- :type:`BookCIDRResponse`
- :type:`CreateNatGatewayRequest`
- :type:`CreateSubnetPoolRequest`
- :type:`CreateSubnetPoolResponse`
- :type:`CreateVPCRequest`
- :type:`CreateVPCResponse`
- :type:`DHCPPoolStr`
- :type:`IPAddress`
- :type:`IpAddress`
- :type:`LinkModel`
- :type:`ListOptions`
- :type:`ListSubnetPoolsOptions`
- :type:`ListSubnetPoolsResponse`
- :type:`ListSubnetsResponse`
- :type:`ListVPCsResponse`
- :type:`Meta`
- :type:`MetaLinks`
- :type:`MetaModel`
- :type:`MetaPageInfo`
- :type:`NatGatewayCreateResponse`
- :type:`NatGatewayDetailsResponse`
- :type:`NatGatewayListResponse`
- :type:`NatGatewayResponse`
- :type:`NatGatewayService`
- :type:`NetworkClient`
- :type:`PageModel`
- :type:`PortCreateOptions`
- :type:`PortCreateRequest`
- :type:`PortCreateResponse`
- :type:`PortIPAddress`
- :type:`PortListResponse`
- :type:`PortNetworkResponse`
- :type:`PortPublicIP`
- :type:`PortResponse`
- :type:`PortService`
- :type:`PortSimpleResponse`
- :type:`PortUpdateRequest`
- :type:`PortsList`
- :type:`PublicIPCreateRequest`
- :type:`PublicIPCreateResponse`
- :type:`PublicIPDb`
- :type:`PublicIPListResponse`
- :type:`PublicIPResponse`
- :type:`PublicIPService`
- :type:`PublicIPsList`
- :type:`PublicIpResponsePort`
- :type:`RenameVPCRequest`
- :type:`RuleCreateRequest`
- :type:`RuleCreateResponse`
- :type:`RuleResponse`
- :type:`RuleService`
- :type:`RulesList`
- :type:`SecurityGroupCreateRequest`
- :type:`SecurityGroupCreateResponse`
- :type:`SecurityGroupDetailResponse`
- :type:`SecurityGroupListResponse`
- :type:`SecurityGroupResponse`
- :type:`SecurityGroupService`
- :type:`SubnetCreateOptions`
- :type:`SubnetCreateRequest`
- :type:`SubnetCreateResponse`
- :type:`SubnetPatchRequest`
- :type:`SubnetPoolDetailsResponse`
- :type:`SubnetPoolResponse`
- :type:`SubnetPoolService`
- :type:`SubnetResponse`
- :type:`SubnetResponseDetail`
- :type:`SubnetResponseId`
- :type:`SubnetService`
- :type:`UnbookCIDRRequest`
- :type:`VPC`
- :type:`VPCService`
- :type:`VPCStateV1`
- :type:`VPCStatusV1`

Constants
---------

- :const:`SecurityGroupsExpand`
- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/network"

   // Use the Network package
   // See the examples directory for complete examples

