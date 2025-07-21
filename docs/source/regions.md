---
title: Global Services
---

# Global Services

Some Magalu Cloud services operate globally and use a dedicated global endpoint (api.magalu.cloud). These global services are:

- SSH Keys Management
- Availability Zones

When using global services, any region configuration set on the core client will be overridden. The service will automatically use the global endpoint.

```go
// Even if core client has a region set
core := client.NewMgcClient(apiToken, client.WithBaseURL(client.BrSe1))

// Global services will ignore the region and use global endpoint
sshClient := sshkeys.New(core) // Uses api.magalu.cloud

// To use a custom endpoint for a global service, use the service's specific option
sshClient := sshkeys.New(core, sshkeys.WithGlobalBasePath("custom-endpoint"))
```

> **Note**: Setting a region on the core client only affects regional services. Global services will always use their global endpoint unless explicitly configured otherwise using their specific options.
