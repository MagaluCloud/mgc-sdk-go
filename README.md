# MGC Go SDK

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=bugs)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=coverage)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=MagaluCloud_mgc-sdk-go&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=MagaluCloud_mgc-sdk-go)

The MGC Go SDK provides a convenient way to interact with the Magalu Cloud API from Go applications.

For more information about Magalu Cloud, visit:
- Website: [https://magalu.cloud/](https://magalu.cloud/)
- Documentation: [https://docs.magalu.cloud/](https://docs.magalu.cloud/)

## Installation

```bash
go get github.com/MagaluCloud/mgc-sdk-go
```

## Authentication

The SDK uses an API token for authentication. 

```go
client := client.New("your-api-token")
```

Find more information about how to generate an API token in the [Magalu Cloud documentation](https://docs.magalu.cloud/docs/devops-tools/api-keys/overview).

## Project Structure

```
mgc-sdk-go/
├── client/         # Base client implementation and configuration
├── compute/        # Compute service API (instances, images, machine types)
├── helpers/        # Utility functions
├── internal/       # Internal packages
└── cmd/           # Command-line examples
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
    client.WithRetryConfig(client.RetryConfig{
        MaxAttempts: 5,
        InitialInterval: 2 * time.Second,
        MaxInterval: 60 * time.Second,
        BackoffFactor: 1.5,
    }),
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
    client.WithRetryConfig(client.RetryConfig{
        MaxAttempts: 3,
        InitialInterval: 1 * time.Second,
        MaxInterval: 30 * time.Second,
    }),
)
```

## Full Example

Check the [cmd/examples](cmd/examples) directory for complete working examples of all SDK features.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.