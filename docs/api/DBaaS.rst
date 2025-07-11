DBaaS
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 5

Package Documentation
-------------------

.. code-block:: go

   package dbaas // import "github.com/MagaluCloud/mgc-sdk-go/dbaas"
   
   Package dbaas provides a client for interacting with the Magalu Cloud Database
   as a Service (DBaaS) API. This package allows you to manage database instances,
   clusters, replicas, engines, instance types, and parameters.
   
   const InstancePath = "/v2/instances" ...
   const ParameterGroupTypeSystem ParameterGroupType = "SYSTEM" ...
   const DefaultBasePath = "/database"
   type Address struct{ ... }
   type AddressAccess string
       const AddressAccessPrivate AddressAccess = "PRIVATE" ...
   type AddressType string
       const AddressTypeIPv4 AddressType = "IPv4" ...
   type ClientOption func(*DBaaSClient)
   type ClusterCreateRequest struct{ ... }
   type ClusterDetailResponse struct{ ... }
   type ClusterResponse struct{ ... }
   type ClusterService interface{ ... }
   type ClusterStatus string


Types
-----

- :type:`Address`
- :type:`AddressAccess`
- :type:`AddressType`
- :type:`ClientOption`
- :type:`ClusterCreateRequest`
- :type:`ClusterDetailResponse`
- :type:`ClusterResponse`
- :type:`ClusterService`
- :type:`ClusterStatus`
- :type:`ClusterUpdateRequest`
- :type:`ClusterVolumeRequest`
- :type:`ClusterVolumeResponse`
- :type:`ClustersResponse`
- :type:`DBaaSClient`
- :type:`DatabaseInstanceUpdateRequest`
- :type:`EngineDetail`
- :type:`EngineParameterDetail`
- :type:`EngineParametersResponse`
- :type:`EngineService`
- :type:`FieldValueFilter`
- :type:`GetInstanceOptions`
- :type:`InstanceCreateRequest`
- :type:`InstanceDetail`
- :type:`InstanceParametersRequest`
- :type:`InstanceParametersResponse`
- :type:`InstanceResizeRequest`
- :type:`InstanceResponse`
- :type:`InstanceService`
- :type:`InstanceStatus`
- :type:`InstanceStatusUpdate`
- :type:`InstanceType`
- :type:`InstanceTypeService`
- :type:`InstanceVolumeRequest`
- :type:`InstanceVolumeResizeRequest`
- :type:`InstancesResponse`
- :type:`ListClustersOptions`
- :type:`ListEngineOptions`
- :type:`ListEngineParametersOptions`
- :type:`ListEnginesResponse`
- :type:`ListInstanceOptions`
- :type:`ListInstanceTypeOptions`
- :type:`ListInstanceTypesResponse`
- :type:`ListParameterGroupsOptions`
- :type:`ListParametersOptions`
- :type:`ListReplicaOptions`
- :type:`ListSnapshotOptions`
- :type:`LoadBalancerAddress`
- :type:`MetaResponse`
- :type:`PageResponse`
- :type:`ParameterCreateRequest`
- :type:`ParameterDetailResponse`
- :type:`ParameterGroupCreateRequest`
- :type:`ParameterGroupDetailResponse`
- :type:`ParameterGroupResponse`
- :type:`ParameterGroupService`
- :type:`ParameterGroupType`
- :type:`ParameterGroupUpdateRequest`
- :type:`ParameterGroupsResponse`
- :type:`ParameterResponse`
- :type:`ParameterService`
- :type:`ParameterUpdateRequest`
- :type:`ParametersResponse`
- :type:`ReplicaAddressResponse`
- :type:`ReplicaCreateRequest`
- :type:`ReplicaDetailResponse`
- :type:`ReplicaResizeRequest`
- :type:`ReplicaResponse`
- :type:`ReplicaService`
- :type:`ReplicasResponse`
- :type:`RestoreSnapshotRequest`
- :type:`SnapshotCreateRequest`
- :type:`SnapshotDetailResponse`
- :type:`SnapshotInstanceDetailResponse`
- :type:`SnapshotResponse`
- :type:`SnapshotStatus`
- :type:`SnapshotType`
- :type:`SnapshotUpdateRequest`
- :type:`SnapshotsResponse`
- :type:`Volume`

Constants
---------

- :const:`InstancePath`
- :const:`ParameterGroupTypeSystem`
- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/dbaas"

   // Use the DBaaS package
   // See the examples directory for complete examples

