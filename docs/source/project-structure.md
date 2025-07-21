# Project Structure

This section describes the organization of files and directories in the MGC Go SDK.

## Structure Overview

```
mgc-sdk-go/
├── client/         # Base client implementation and configuration
├── compute/        # Computing service API (instances, images, machine types)
├── blockstorage/   # Block storage service API
├── network/        # Network service API
├── kubernetes/     # Kubernetes service API
├── dbaas/          # Database as a Service API
├── containerregistry/ # Container Registry service API
├── sshkeys/        # SSH Keys service API
├── availabilityzones/ # Availability Zones service API
├── audit/          # Audit service API
├── lbaas/          # Load Balancer as a Service API
├── helpers/        # Utility functions
├── internal/       # Internal packages
└── cmd/            # Usage examples
```

## Module Descriptions

### client/
Contains the base HTTP client implementation and configurations for communicating with the Magalu Cloud API.

### compute/
Provides functionality to manage virtual instances, machine types, images, and snapshots.

### blockstorage/
Allows managing block storage volumes, snapshots, and volume types.

### network/
Provides functionality to manage VPCs, subnets, security groups, and other network resources.

### kubernetes/
Allows managing Kubernetes clusters, nodepools, and related configurations.

### dbaas/
Provides functionality to manage database instances, clusters, and configurations.

### containerregistry/
Allows managing container registries, repositories, and images.

### sshkeys/
Provides functionality to manage SSH keys.

### availabilityzones/
Allows querying available availability zones.

### audit/
Provides functionality to access audit logs and events.

### lbaas/
Allows managing load balancers and related configurations.

### helpers/
Contains reusable utility functions throughout the SDK.

### internal/
Contains internal packages not publicly exposed.

### cmd/
Contains practical examples of how to use each SDK module.
