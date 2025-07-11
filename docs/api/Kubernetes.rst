Kubernetes
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 10

Package Documentation
-------------------

.. code-block:: go

   package kubernetes // import "github.com/MagaluCloud/mgc-sdk-go/kubernetes"
   
   Package kubernetes provides a client for interacting with the Magalu Cloud
   Kubernetes API. This package allows you to manage Kubernetes clusters, node
   pools, flavors, and versions.
   
   const DefaultBasePath = "/kubernetes"
   type Addons struct{ ... }
   type Addresses struct{ ... }
   type Allocatable struct{ ... }
   type AllowedCIDRsUpdateRequest struct{ ... }
   type AutoScale struct{ ... }
   type AutoScaleResponse struct{ ... }
   type Capacity struct{ ... }
   type ClientOption func(*KubernetesClient)
   type Cluster struct{ ... }
   type ClusterList struct{ ... }
   type ClusterListResponse struct{ ... }
   type ClusterRequest struct{ ... }
   type ClusterService interface{ ... }


Types
-----

- :type:`Addons`
- :type:`Addresses`
- :type:`Allocatable`
- :type:`AllowedCIDRsUpdateRequest`
- :type:`AutoScale`
- :type:`AutoScaleResponse`
- :type:`Capacity`
- :type:`ClientOption`
- :type:`Cluster`
- :type:`ClusterList`
- :type:`ClusterListResponse`
- :type:`ClusterRequest`
- :type:`ClusterService`
- :type:`Controlplane`
- :type:`CreateClusterResponse`
- :type:`CreateNodePoolRequest`
- :type:`Flavor`
- :type:`FlavorList`
- :type:`FlavorService`
- :type:`FlavorsAvailable`
- :type:`Infrastructure`
- :type:`InstanceTemplate`
- :type:`KubeApiServer`
- :type:`KubeConfig`
- :type:`KubernetesClient`
- :type:`ListOptions`
- :type:`MessageState`
- :type:`Network`
- :type:`Node`
- :type:`NodePool`
- :type:`NodePoolList`
- :type:`NodePoolService`
- :type:`PatchNodePoolRequest`
- :type:`Status`
- :type:`Taint`
- :type:`Version`
- :type:`VersionList`
- :type:`VersionService`

Constants
---------

- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/kubernetes"

   // Use the Kubernetes package
   // See the examples directory for complete examples

