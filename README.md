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
├── client/         # Base client implementation
├── helpers/        # Utility functions
├── sdk/           # Service-specific implementations
│   ├── virtualmachine/
│   └── ... (other services)
└── examples/      # Usage examples
```

## Usage Examples

### Initializing the Client

```go
import (
    "github.com/MagaluCloud/mgc-sdk-go/client"
    "github.com/MagaluCloud/mgc-sdk-go/sdk/virtualmachine"
)

apiToken := os.Getenv("MGC_API_TOKEN")
c := client.New(apiToken)
vmClient := virtualmachine.New(c)
```

### Listing Virtual Machines

```go
instances, err := vmClient.Instances().List(context.Background(), virtualmachine.ListOptions{
    Limit:  helpers.IntPtr(10),
    Offset: helpers.IntPtr(0),
    Expand: []string{"machine-type", "image"},
})
```

### Creating a Virtual Machine

```go
createReq := virtualmachine.CreateRequest{
    Name: "my-test-vm",
    MachineType: virtualmachine.IDOrName{
        Name: helpers.StrPtr("BV1-1-40"),
    },
    Image: virtualmachine.IDOrName{
        Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
    },
    Network: &virtualmachine.CreateParametersNetwork{
        AssociatePublicIp: helpers.BoolPtr(false),
    },
    SshKeyName: helpers.StrPtr("my-ssh-key"),
}

id, err := vmClient.Instances().Create(context.Background(), createReq)
```

### Managing Virtual Machines

```go
// Get instance details
instance, err := vmClient.Instances().Get(ctx, id, []string{"network"})

// Rename instance
err = vmClient.Instances().Rename(ctx, instanceID, "new-name")

// Change machine type
err = vmClient.Instances().Retype(ctx, instanceID, virtualmachine.RetypeRequest{
    MachineType: virtualmachine.IDOrName{
        Name: helpers.StrPtr("BV2-2-20"),
    },
})

// Delete instance
err = vmClient.Instances().Delete(ctx, instanceID, true)
```

## Full Example

Check the [examples directory](examples/) for complete working examples of all SDK features.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.