# MGC Go SDK

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=bugs)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=coverage)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)
[![CodeQL](https://github.com/MagaluCloud/mgc-sdk-go/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/MagaluCloud/mgc-sdk-go/actions/workflows/github-code-scanning/codeql)
[![Unit Tests](https://github.com/MagaluCloud/mgc-sdk-go/actions/workflows/unit-tests.yml/badge.svg?branch=main&event=push)](https://github.com/MagaluCloud/mgc-sdk-go/actions/workflows/unit-tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/MagaluCloud/mgc-sdk-go)](https://goreportcard.com/report/github.com/MagaluCloud/mgc-sdk-go)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/MagaluCloud/mgc-sdk-go)](https://pkg.go.dev/github.com/MagaluCloud/mgc-sdk-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The MGC Go SDK provides a convenient way to interact with the Magalu Cloud API from Go applications.

This is an very active project and work in progress, with many products support comming soon.

For more information about Magalu Cloud, visit:
- Website: [https://magalu.cloud/](https://magalu.cloud/)
- Documentation: [https://docs.magalu.cloud/](https://docs.magalu.cloud/)

## Installation

```bash
go get github.com/MagaluCloud/mgc-sdk-go
```

## Supported Products

- Compute (Virtual Machines)
  - Instances
  - Machine Types
  - Images
  - Snapshots
- Block Storage
  - Volumes
  - Snapshots
  - Volume Types
- SSH Keys
- Availability Zones
- Audit
  - Events
  - Events Types

## Authentication

The SDK uses an API token for authentication.

```go
client := client.New("your-api-token")
```

Find more information about how to generate an API token in the [Magalu Cloud documentation](https://docs.magalu.cloud/docs/devops-tools/api-keys/overview).

## Global Services

Some Magalu Cloud services operate globally and use a dedicated global endpoint (api.magalu.cloud). These global services are:
- SSH Keys Management
- Availability Zones

When using global services, any region configuration set on the core client will be overridden. The service will automatically use the global endpoint.

```go
// Even if core client has a region set
core := client.NewMgcClient(apiToken, client.WithRegion(client.BrMgl1))

// Global services will ignore the region and use global endpoint
sshClient := sshkeys.New(core) // Uses api.magalu.cloud

// To use a custom endpoint for a global service, use the service's specific option
sshClient := sshkeys.New(core, sshkeys.WithGlobalBasePath("custom-endpoint"))
```

> **Note**: Setting a region on the core client only affects regional services. Global services will always use their global endpoint unless explicitly configured otherwise using their specific options.

## Project Structure

```
mgc-sdk-go/
├── client/         # Base client implementation and configuration
├── compute/        # Compute service API (instances, images, machine types)
├── helpers/        # Utility functions
├── internal/       # Internal packages
└── cmd/            # Examples
```

## Usage Examples

### Initializing the Client

```go
import (
    "github.com/MagaluCloud/mgc-sdk-go/client"
    "github.com/MagaluCloud/mgc-sdk-go/compute"
)

apiToken := os.Getenv("MGC_API_TOKEN")
c := client.NewMgcClient(apiToken)
computeClient := compute.New(c)
```

### Client Configuration Options

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
)
```

Available options:
- `WithTimeout`: Sets the client timeout for requests
- `WithUserAgent`: Sets a custom User-Agent header
- `WithLogger`: Configures a custom logger
- `WithRetryConfig`: Customizes the retry behavior
- `WithHTTPClient`: Uses a custom HTTP client
- `WithBaseURL`: Changes the API endpoint (useful for testing)

### Listing Instances

```go
instances, err := computeClient.Instances().List(context.Background(), compute.ListOptions{
    Limit:  helpers.IntPtr(10),
    Offset: helpers.IntPtr(0),
    Expand: []string{compute.InstanceMachineTypeExpand, compute.InstanceImageExpand},
})
```

### Creating an Instance

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

### Managing Machine Types

```go
machineTypes, err := computeClient.MachineTypes().List(context.Background(), compute.MachineTypeListOptions{})
```

### Managing Images

```go
images, err := computeClient.Images().List(context.Background(), compute.ImageListOptions{})
```

### Using Request IDs

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

## Error Handling

### HTTP Errors
```go
err := computeClient.Instances().Delete(ctx, id)
if httpErr, ok := err.(*client.HTTPError); ok {
    switch httpErr.StatusCode {
    case 404:
        log.Fatal("Instance not found")
    case 403:
        log.Fatal("Permission denied")
    case 429:
        log.Fatal("Rate limit exceeded")
    }
    log.Printf("Response body: %s", string(httpErr.Body))
}
```

### Validation Errors
```go
_, err := computeClient.Instances().Create(ctx, compute.CreateRequest{})
if validErr, ok := err.(*client.ValidationError); ok {
    log.Printf("Invalid field %s: %s", validErr.Field, validErr.Message)
}
```

### Error Types and Interfaces

The SDK provides these error types:

```go
// HTTPError contains details about API request failures
type HTTPError struct {
    StatusCode int         // HTTP status code
    Status     string      // Status message
    Body       []byte      // Raw response body
    Response   *http.Response
}

// ValidationError occurs when request parameters are invalid
type ValidationError struct {
    Field   string   // Which field failed validation
    Message string   // Why the validation failed
}
```

Common error handling patterns:

```go
// Check for specific error types
err := computeClient.Instances().Delete(ctx, id)
switch e := err.(type) {
case *client.HTTPError:
    // Handle HTTP errors (404, 403, etc)
    fmt.Printf("HTTP %d: %s\n", e.StatusCode, e.Body)
case *client.ValidationError:
    // Handle validation failures
    fmt.Printf("Invalid %s: %s\n", e.Field, e.Message)
default:
    // Handle other errors (context timeout, network issues)
    fmt.Printf("Error: %v\n", err)
}

// Check if error has additional details
if detailed, ok := err.(interface{ ErrorDetails() map[string]interface{} }); ok {
    details := detailed.ErrorDetails()
    fmt.Printf("Error details: %+v\n", details)
}
```

### Retries
The client automatically retries on network errors and 5xx responses:
```go
client := client.NewMgcClient(
    apiToken,
    client.WithRetryConfig(
        3, // maxAttempts
        1 * time.Second, // initialInterval
        30 * time.Second, // maxInterval
        1.5, // backoffFactor
    ),
)
```

## Full Example

Check the [cmd/examples](cmd/examples) directory for complete working examples of all SDK features.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
