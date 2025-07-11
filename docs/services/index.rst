Services
========

The MGC Go SDK provides access to all major Magalu Cloud services. Each service is organized into its own package with a dedicated client.

.. toctree::
   :maxdepth: 2
   :caption: Available Services:

   compute
   network
   blockstorage
   dbaas
   containerregistry
   kubernetes
   sshkeys
   availabilityzones
   audit

Service Overview
---------------

Compute Services
~~~~~~~~~~~~~~~

The compute service provides access to virtual machines and related resources:

- **Instances**: Create, manage, and monitor virtual machines
- **Machine Types**: List available machine configurations
- **Images**: Manage operating system images
- **Snapshots**: Create and manage instance snapshots

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/compute"

   computeClient := compute.New(coreClient)
   instances := computeClient.Instances().List(ctx, compute.ListOptions{})

Network Services
~~~~~~~~~~~~~~~

The network service manages networking infrastructure:

- **VPCs**: Virtual Private Clouds for network isolation
- **Subnets**: Network segments within VPCs
- **Security Groups**: Firewall rules for instances
- **Public IPs**: Public IP address management
- **NAT Gateways**: Network Address Translation
- **Subnet Pools**: IP address pool management

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/network"

   networkClient := network.New(coreClient)
   vpcs := networkClient.VPCs().List(ctx)

Block Storage Services
~~~~~~~~~~~~~~~~~~~~~

The block storage service manages persistent storage:

- **Volumes**: Block storage volumes for instances
- **Snapshots**: Volume snapshots for backup
- **Volume Types**: Different storage performance tiers

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/blockstorage"

   storageClient := blockstorage.New(coreClient)
   volumes := storageClient.Volumes().List(ctx, blockstorage.ListOptions{})

Database as a Service (DBaaS)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The DBaaS service provides managed database instances:

- **Instances**: Database server instances
- **Instance Types**: Database server configurations
- **Snapshots**: Database backups
- **Replicas**: Read replicas for scaling
- **Engines**: Database engines (MySQL, PostgreSQL, etc.)
- **Clusters**: Database clusters
- **Parameters**: Database configuration parameters

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/dbaas"

   dbaasClient := dbaas.New(coreClient)
   instances := dbaasClient.Instances().List(ctx, dbaas.ListOptions{})

Container Registry
~~~~~~~~~~~~~~~~~

The container registry service manages Docker images:

- **Repositories**: Docker image repositories
- **Registries**: Container registry instances
- **Images**: Docker images within repositories
- **Credentials**: Registry authentication

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/containerregistry"

   registryClient := containerregistry.New(coreClient)
   repositories := registryClient.Repositories().List(ctx, containerregistry.ListOptions{})

Kubernetes Services
~~~~~~~~~~~~~~~~~~

The Kubernetes service manages container orchestration:

- **Clusters**: Kubernetes clusters
- **Flavors**: Cluster configurations
- **Node Pools**: Worker node groups
- **Version**: Kubernetes version information

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/kubernetes"

   k8sClient := kubernetes.New(coreClient)
   clusters := k8sClient.Clusters().List(ctx, kubernetes.ListOptions{})

Global Services
~~~~~~~~~~~~~~~

Some services operate globally and are not region-specific:

SSH Keys
^^^^^^^^

Manage SSH keys for instance access:

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/sshkeys"

   sshClient := sshkeys.New(coreClient)
   keys := sshClient.Keys().List(ctx, sshkeys.ListOptions{})

Availability Zones
^^^^^^^^^^^^^^^^^

List available zones for resource placement:

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/availabilityzones"

   azClient := availabilityzones.New(coreClient)
   zones := azClient.List(ctx)

Audit Services
~~~~~~~~~~~~~

Monitor and audit system activities:

- **Events**: System audit events
- **Event Types**: Types of audit events

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/audit"

   auditClient := audit.New(coreClient)
   events := auditClient.Events().List(ctx, audit.ListOptions{})

Common Patterns
--------------

All services follow similar patterns for consistency:

Client Initialization
~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   // Create core client
   coreClient := client.NewMgcClient(apiToken, client.WithBaseURL(client.BrSe1))

   // Create service-specific client
   serviceClient := service.New(coreClient)

List Operations
~~~~~~~~~~~~~~

Most services support pagination and filtering:

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/helpers"

   opts := service.ListOptions{
       Limit:  helpers.IntPtr(10),
       Offset: helpers.IntPtr(0),
       Sort:   helpers.StrPtr("name"),
   }

   resources, err := serviceClient.Resources().List(ctx, opts)

Create Operations
~~~~~~~~~~~~~~~~

Create operations typically require a request struct:

.. code-block:: go

   req := service.CreateRequest{
       Name: "my-resource",
       // ... other fields
   }

   id, err := serviceClient.Resources().Create(ctx, req)

Get Operations
~~~~~~~~~~~~~

Retrieve specific resources by ID:

.. code-block:: go

   resource, err := serviceClient.Resources().Get(ctx, "resource-id")

Delete Operations
~~~~~~~~~~~~~~~~

Delete resources by ID:

.. code-block:: go

   err := serviceClient.Resources().Delete(ctx, "resource-id")

Error Handling
~~~~~~~~~~~~~

All operations return errors that should be handled:

.. code-block:: go

   resource, err := serviceClient.Resources().Get(ctx, "resource-id")
   if err != nil {
       if httpErr, ok := err.(*client.HTTPError); ok {
           switch httpErr.StatusCode {
           case 404:
               log.Println("Resource not found")
           case 403:
               log.Println("Permission denied")
           }
       }
       return err
   }

Service-Specific Features
------------------------

Each service may have additional features beyond the basic CRUD operations:

- **Compute**: Instance actions (start, stop, reboot), console access
- **Network**: Security group rules, port management
- **Block Storage**: Volume attachments, snapshot management
- **DBaaS**: Database operations, parameter management
- **Container Registry**: Image pushing/pulling, vulnerability scanning
- **Kubernetes**: Cluster scaling, node management

For detailed information about each service, see the individual service documentation pages. 