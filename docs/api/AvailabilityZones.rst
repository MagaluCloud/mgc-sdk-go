AvailabilityZones
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 17

Package Documentation
-------------------

.. code-block:: go

   package availabilityzones // import "github.com/MagaluCloud/mgc-sdk-go/availabilityzones"
   
   Package availabilityzones provides functionality to interact with the
   MagaluCloud availability zones service. This package allows listing availability
   zones across different regions.
   
   const DefaultBasePath = "/profile"
   type AvailabilityZone struct{ ... }
   type BlockType string
       const BlockTypeNone BlockType = "none" ...
   type Client struct{ ... }
       func New(core *client.CoreClient, opts ...ClientOption) *Client
   type ClientOption func(*Client)
       func WithGlobalBasePath(basePath client.MgcUrl) ClientOption
   type ListOptions struct{ ... }
   type ListResponse struct{ ... }
   type Region struct{ ... }
   type Service interface{ ... }


Types
-----

- :type:`AvailabilityZone`
- :type:`BlockType`
- :type:`Client`
- :type:`ClientOption`
- :type:`ListOptions`
- :type:`ListResponse`
- :type:`Region`
- :type:`Service`

Constants
---------

- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/availabilityzones"

   // Use the AvailabilityZones package
   // See the examples directory for complete examples

