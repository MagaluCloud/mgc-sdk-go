ContainerRegistry
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 17

Package Documentation
-------------------

.. code-block:: go

   package containerregistry // import "github.com/MagaluCloud/mgc-sdk-go/containerregistry"
   
   Package containerregistry provides a client for interacting with the Magalu
   Cloud Container Registry API. This package allows you to manage container
   registries, repositories, images, and credentials.
   
   const DefaultBasePath = "/container-registry"
   type AmountRepositoryResponse struct{ ... }
   type ClientOption func(*ContainerRegistryClient)
   type ContainerRegistryClient struct{ ... }
       func New(core *client.CoreClient, opts ...ClientOption) *ContainerRegistryClient
   type CredentialsResponse struct{ ... }
   type CredentialsService interface{ ... }
   type ImageResponse struct{ ... }
   type ImageTagResponse struct{ ... }
   type ImagesResponse struct{ ... }
   type ImagesService interface{ ... }
   type ListOptions struct{ ... }
   type ListRegistriesResponse struct{ ... }
   type RegistriesService interface{ ... }


Types
-----

- :type:`AmountRepositoryResponse`
- :type:`ClientOption`
- :type:`ContainerRegistryClient`
- :type:`CredentialsResponse`
- :type:`CredentialsService`
- :type:`ImageResponse`
- :type:`ImageTagResponse`
- :type:`ImagesResponse`
- :type:`ImagesService`
- :type:`ListOptions`
- :type:`ListRegistriesResponse`
- :type:`RegistriesService`
- :type:`RegistryRequest`
- :type:`RegistryResponse`
- :type:`RepositoriesResponse`
- :type:`RepositoriesService`
- :type:`RepositoryResponse`

Constants
---------

- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/containerregistry"

   // Use the ContainerRegistry package
   // See the examples directory for complete examples

