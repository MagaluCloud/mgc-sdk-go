Audit
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 5

Package Documentation
-------------------

.. code-block:: go

   package audit // import "github.com/MagaluCloud/mgc-sdk-go/audit"
   
   Package audit provides functionality to interact with the MagaluCloud audit
   service. This package allows listing audit events and event types.
   
   const DefaultBasePath = "/audit"
   type AuditClient struct{ ... }
       func New(core *client.CoreClient) *AuditClient
   type Event struct{ ... }
   type EventService interface{ ... }
   type EventType struct{ ... }
   type EventTypeService interface{ ... }
   type ListEventTypesParams struct{ ... }
   type ListEventsParams struct{ ... }
   type PaginatedMeta struct{ ... }
   type PaginatedResponse[T any] struct{ ... }


Types
-----

- :type:`AuditClient`
- :type:`Event`
- :type:`EventService`
- :type:`EventType`
- :type:`EventTypeService`
- :type:`ListEventTypesParams`
- :type:`ListEventsParams`
- :type:`PaginatedMeta`
- :type:`PaginatedResponse[T`

Constants
---------

- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/audit"

   // Use the Audit package
   // See the examples directory for complete examples

