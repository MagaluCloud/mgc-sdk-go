SSHKeys
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 7

Package Documentation
-------------------

.. code-block:: go

   package sshkeys // import "github.com/MagaluCloud/mgc-sdk-go/sshkeys"
   
   Package sshkeys provides client implementation for managing SSH keys in the
   Magalu Cloud platform. SSH keys are managed as a global service, meaning they
   are not bound to any specific region. By default, the service uses the global
   endpoint, but this can be overridden if needed.
   
   const DefaultBasePath = "/profile"
   type ClientOption func(*SSHKeyClient)
       func WithGlobalBasePath(basePath client.MgcUrl) ClientOption
   type CreateSSHKeyRequest struct{ ... }
   type KeyService interface{ ... }
   type ListOptions struct{ ... }
   type ListSSHKeysResponse struct{ ... }
   type SSHKey struct{ ... }
   type SSHKeyClient struct{ ... }
       func New(core *client.CoreClient, opts ...ClientOption) *SSHKeyClient


Types
-----

- :type:`ClientOption`
- :type:`CreateSSHKeyRequest`
- :type:`KeyService`
- :type:`ListOptions`
- :type:`ListSSHKeysResponse`
- :type:`SSHKey`
- :type:`SSHKeyClient`

Constants
---------

- :const:`DefaultBasePath`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/sshkeys"

   // Use the SSHKeys package
   // See the examples directory for complete examples

