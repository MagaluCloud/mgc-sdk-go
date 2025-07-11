MGC Go SDK Documentation
========================

.. image:: https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=alert_status
   :target: https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go
   :alt: Quality Gate Status

.. image:: https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=coverage
   :target: https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go
   :alt: Coverage

.. image:: https://pkg.go.dev/badge/github.com/MagaluCloud/mgc-sdk-go
   :target: https://pkg.go.dev/github.com/MagaluCloud/mgc-sdk-go
   :alt: Go Reference

.. image:: https://img.shields.io/badge/License-Apache%202.0-blue.svg
   :target: https://opensource.org/licenses/Apache-2.0
   :alt: License

The MGC Go SDK provides a convenient way to interact with the Magalu Cloud API from Go applications. This SDK is actively maintained and supports a wide range of Magalu Cloud services.

For more information about Magalu Cloud, visit:

- **Website**: https://magalu.cloud/
- **Documentation**: https://docs.magalu.cloud/

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   installation
   authentication
   configuration
   services/index
   api/index
   examples
   error-handling
   contributing

Quick Start
----------

Install the SDK:

.. code-block:: bash

   go get github.com/MagaluCloud/mgc-sdk-go

Basic usage:

.. code-block:: go

   package main

   import (
       "context"
       "log"
       "os"

       "github.com/MagaluCloud/mgc-sdk-go/client"
       "github.com/MagaluCloud/mgc-sdk-go/compute"
   )

   func main() {
       // Initialize the client
       apiToken := os.Getenv("MGC_API_TOKEN")
       c := client.NewMgcClient(apiToken)
       computeClient := compute.New(c)

       // List instances
       instances, err := computeClient.Instances().List(context.Background(), compute.ListOptions{})
       if err != nil {
           log.Fatal(err)
       }

       for _, instance := range instances {
           log.Printf("Instance: %s", instance.Name)
       }
   }

Supported Services
-----------------

The MGC Go SDK supports the following Magalu Cloud services:

**Compute Services**
- Virtual Machines (Instances)
- Machine Types
- Images
- Snapshots

**Storage Services**
- Block Storage (Volumes)
- Volume Snapshots
- Volume Types

**Network Services**
- VPCs (Virtual Private Clouds)
- Subnets
- Security Groups
- Public IPs
- NAT Gateways
- Subnet Pools

**Database Services**
- Database as a Service (DBaaS)
- Instances
- Instance Types
- Snapshots
- Replicas
- Engines
- Clusters
- Parameters

**Container Services**
- Container Registry
- Repositories
- Registries
- Images
- Credentials

**Kubernetes Services**
- Clusters
- Flavors
- Node Pools
- Version Information

**Global Services**
- SSH Keys Management
- Availability Zones

**Monitoring & Audit**
- Audit Events
- Event Types

Key Features
-----------

- **Region Support**: Full support for Magalu Cloud regions (BR-SE1, BR-NE1)
- **Global Services**: Support for global services like SSH Keys
- **Error Handling**: Comprehensive error handling with detailed error types
- **Retry Logic**: Built-in retry mechanism with exponential backoff
- **Request Tracking**: Support for request IDs for tracing
- **Logging**: Configurable logging with structured logging support
- **Type Safety**: Full Go type safety with comprehensive structs
- **Documentation**: Complete API documentation with examples

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search` 