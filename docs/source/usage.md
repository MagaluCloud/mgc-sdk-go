---
title: Usage Examples
---

# Usage Examples

## Initializing the Client

```go
import (
    "github.com/MagaluCloud/mgc-sdk-go/client"
    "github.com/MagaluCloud/mgc-sdk-go/compute"
)

apiToken := os.Getenv("MGC_API_TOKEN")
c := client.NewMgcClient(apiToken)
computeClient := compute.New(c)
```

## Client Configuration Options

You can customize the client behavior using options:

```go
import (
    "time"
    "log/slog"
)

// Configure client with multiple options
c := client.NewMgcClient(
    apiToken,
    client.WithTimeout(5 * time.Minute),
    client.WithUserAgent("my-app/1.0"),
    client.WithLogger(slog.Default().With("service", "mgc")),
    client.WithRetryConfig(
        3, // maxAttempts
        1 * time.Second, // initialInterval
        30 * time.Second, // maxInterval
        1.5, // backoffFactor
    ),
    client.WithBaseURL(client.BrSe1),
)
```

Available options:

- `WithTimeout`: Sets the client timeout for requests
- `WithUserAgent`: Sets a custom User-Agent header
- `WithLogger`: Configures a custom logger
- `WithRetryConfig`: Customizes the retry behavior
- `WithHTTPClient`: Uses a custom HTTP client
- `WithBaseURL`: Changes the API endpoint (useful for testing or setting a specific region to interact)
- `WithCustomHeader`: Adds custom headers to all requests

## Listing Instances

```go
instances, err := computeClient.Instances().List(context.Background(), compute.ListOptions{
    Limit:  helpers.IntPtr(10),
    Offset: helpers.IntPtr(0),
    Expand: []string{compute.InstanceMachineTypeExpand, compute.InstanceImageExpand},
})
```

## Creating an Instance

```go
createReq := compute.CreateRequest{
    Name: "my-test-vm",
    MachineType: compute.IDOrName{
        Name: helpers.StrPtr("BV1-1-40"),
    },
    Image: compute.IDOrName{
        Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
    },
    Network: &compute.CreateParametersNetwork{
        AssociatePublicIp: helpers.BoolPtr(false),
    },
    SshKeyName: helpers.StrPtr("my-ssh-key"),
}

id, err := computeClient.Instances().Create(context.Background(), createReq)
```

## Managing Machine Types

```go
machineTypes, err := computeClient.MachineTypes().List(context.Background(), compute.MachineTypeListOptions{})
```

## Managing Images

```go
images, err := computeClient.Images().List(context.Background(), compute.ImageListOptions{})
```

## Using Request IDs

You can track requests across systems by setting a request ID in the context. The request ID must be a valid UUIDv4 string:

```go
import (
    "context"
    "github.com/google/uuid"
    "github.com/MagaluCloud/mgc-sdk-go/client"
)

// Generate a valid UUIDv4 for the request
requestID := uuid.New().String()

// Create a context with request ID
ctx := context.WithValue(context.Background(), client.RequestIDKey, requestID)

// The client will automatically include the X-Request-ID header
instances, err := computeClient.Instances().List(ctx, compute.ListOptions{})
```

The request ID will be:

- Must be a valid UUIDv4 string (e.g. "123e4567-e89b-12d3-a456-426614174000")
- Included in the request as `X-Request-ID` header
- Logged in the client's logger
- Returned in the response headers for tracking
