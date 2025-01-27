# MGC Go SDK

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
machineTypes, err := computeClient.MachineTypes.List(context.Background(), compute.MachineTypeListOptions{})
```

### Managing Images

```go
images, err := computeClient.Images().List(context.Background(), compute.ImageListOptions{})
```

## Full Example

Check the [cmd/examples](cmd/examples) directory for complete working examples of all SDK features.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.