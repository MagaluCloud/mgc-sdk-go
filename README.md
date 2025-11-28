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
- GitHub: [https://github.com/MagaluCloud/mgc-sdk-go](https://github.com/MagaluCloud/mgc-sdk-go)

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
  - Schedulers
- Object Storage
  - Buckets
  - Objects
- SSH Keys
- Availability Zones
- Audit
  - Events
  - Events Types
- Database as a Service (DBaaS)
  - Instances
  - Instance Types
  - Snapshots
  - Replicas
  - Engines
  - Clusters
  - Parameters
- Container Registry
  - Repositories
  - Registries
  - Images
  - Credentials
  - Proxy Caches
- Kubernetes
  - Clusters
  - Flavors
  - Info
  - Nodepool
  - Version
- Load Balancer as a Service (LBaaS)
  - Load Balancers
  - Listeners
  - Backends
  - Backend Targets
  - Health Checks
  - Certificates
  - Network ACLs
- Network
  - VPCs
  - Subnets
  - Ports
  - Security Groups
  - Rules
  - Public IPs
  - Subnet Pools
  - NAT Gateways

## Authentication

The SDK uses an API token for authentication.

```go
client := client.New("your-api-token")
```

Find more information about how to generate an API token in the [Magalu Cloud documentation](https://docs.magalu.cloud/docs/devops-tools/api-keys/overview).

## Regions

The Magalu Cloud API is region-based, and each service is available in specific regions. You can set the region on the client to interact with a specific region.

### Brazil South East 1 (BR-SE1) - Default

```go
core := client.NewMgcClient(apiToken, client.WithBaseURL(client.BrSe1))
```

### Brazil North East 1 (BR-NE1)

```go
core := client.NewMgcClient(apiToken, client.WithBaseURL(client.BrNe1))
```

## Global Services

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

## Project Structure

```
mgc-sdk-go/
├── client/         # Base client implementation and configuration
├── compute/        # Compute service API (instances, images, machine types)
├── objectstorage/  # Object Storage service API (buckets, objects)
├── helpers/        # Utility functions
├── internal/       # Internal packages
└── cmd/            # Examples
```

## Usage Examples

### Object Storage Usage

The Object Storage service provides an interface for managing buckets and objects in MagaluCloud Object Storage.

#### Endpoint Configuration

By default, the client uses the **BR-SE1 (Brazil Southeast 1)** endpoint. You can override this using the `WithEndpoint()` option to connect to a different region:

- `BrSe1`: Brazil Southeast 1 (default)
- `BrNe1`: Brazil Northeast 1

#### Creating a Client

##### Default Endpoint (BR-SE1)

```go
import (
    "context"
    "time"

    "github.com/MagaluCloud/mgc-sdk-go/client"
    "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
)

apiToken := os.Getenv("MGC_API_TOKEN")
c := client.NewMgcClient(client.WithAPIKey(apiToken))

accessKey := os.Getenv("MGC_OBJECT_STORAGE_ACCESS_KEY")
secretKey := os.Getenv("MGC_OBJECT_STORAGE_SECRET_KEY")

// Uses BR-SE1 by default
osClient, err := objectstorage.New(c, accessKey, secretKey)
```

##### Custom Region Configuration (BR-NE1)

```go
// To use a different region, pass the WithEndpoint option
osClient, err := objectstorage.New(c, accessKey, secretKey,
    objectstorage.WithEndpoint(objectstorage.BrNe1))
```

#### Bucket Operations

##### Listing Buckets

```go
buckets, err := osClient.Buckets().List(context.Background(), objectstorage.BucketListOptions{})
for _, bucket := range buckets {
    fmt.Printf("Bucket: %s (Created: %s)\n", bucket.Name, bucket.CreationDate)
}
```

##### Creating a Bucket

```go
err := osClient.Buckets().Create(context.Background(), "my-bucket")
```

##### Checking if a Bucket Exists

```go
exists, err := osClient.Buckets().Exists(context.Background(), "my-bucket")
if exists {
    fmt.Println("Bucket exists")
}
```

##### Deleting a Bucket

```go
err := osClient.Buckets().Delete(context.Background(), "my-bucket")
```

##### Bucket Policies

Get the policy of a bucket:

```go
policy, err := osClient.Buckets().GetPolicy(context.Background(), "my-bucket")
if policy != nil {
    fmt.Printf("Policy: %+v\n", policy)
}
```

Set a bucket policy:

```go
policy := &objectstorage.Policy{
    Version: "2012-10-17",
    Statement: []objectstorage.Statement{
        {
            Effect:   "Allow",
            Action:   "s3:GetObject",
            Resource: "arn:aws:s3:::my-bucket/*",
        },
    },
}
err := osClient.Buckets().SetPolicy(context.Background(), "my-bucket", policy)
```

Delete a bucket policy:

```go
err := osClient.Buckets().DeletePolicy(context.Background(), "my-bucket")
```

##### Bucket Locking

Lock a bucket (enables Object Lock):

```go
err := osClient.Buckets().LockBucket(context.Background(), "my-bucket")
```

Unlock a bucket (disables Object Lock):

```go
err := osClient.Buckets().UnlockBucket(context.Background(), "my-bucket")
```

Check bucket lock status:

```go
locked, err := osClient.Buckets().GetBucketLockStatus(context.Background(), "my-bucket")
if locked {
    fmt.Println("Bucket is locked")
}
```

##### CORS Configuration

Set CORS configuration:

```go
corsConfig := &objectstorage.CORSConfiguration{
    CORSRules: []objectstorage.CORSRule{
        {
            AllowedOrigins: []string{"https://example.com"},
            AllowedMethods: []string{"GET", "PUT"},
            AllowedHeaders: []string{"*"},
            MaxAgeSeconds:  3600,
        },
    },
}
err := osClient.Buckets().SetCORS(context.Background(), "my-bucket", corsConfig)
```

Get CORS configuration:

```go
corsConfig, err := osClient.Buckets().GetCORS(context.Background(), "my-bucket")
if corsConfig != nil {
    fmt.Printf("CORS Rules: %+v\n", corsConfig.CORSRules)
}
```

Delete CORS configuration:

```go
err := osClient.Buckets().DeleteCORS(context.Background(), "my-bucket")
```

##### Versioning

Enable versioning:

```go
err := osClient.Buckets().EnableVersioning(context.Background(), "my-bucket")
```

Suspend versioning:

```go
err := osClient.Buckets().SuspendVersioning(context.Background(), "my-bucket")
```

Get versioning status:

```go
status, err := osClient.Buckets().GetVersioningStatus(context.Background(), "my-bucket")
fmt.Printf("Versioning Status: %s\n", status.Status)
```

#### Object Operations

##### Uploading an Object

```go
data := []byte("Hello, World!")
err := osClient.Objects().Upload(context.Background(), "my-bucket", "hello.txt", data, "text/plain")
```

##### Downloading an Object

```go
data, err := osClient.Objects().Download(context.Background(), "my-bucket", "hello.txt", nil)
fmt.Printf("Downloaded: %s\n", string(data))
```

Download a specific version:

```go
opts := &objectstorage.DownloadOptions{VersionID: "version-id"}
data, err := osClient.Objects().Download(context.Background(), "my-bucket", "hello.txt", opts)
```

##### Streaming Downloads

```go
reader, err := osClient.Objects().DownloadStream(context.Background(), "my-bucket", "large-file.mp4", nil)
if err == nil {
    defer reader.Close()
    // Process the stream
    io.Copy(os.Stdout, reader)
}
```

##### Listing Objects

List objects with pagination:

```go
opts := objectstorage.ObjectListOptions{
    Limit:  intPtr(10),
    Offset: intPtr(0),
    Prefix: "folder/",
}
objects, err := osClient.Objects().List(context.Background(), "my-bucket", opts)
for _, obj := range objects {
    fmt.Printf("Object: %s (Size: %d)\n", obj.Key, obj.Size)
}
```

List all objects (without pagination):

```go
opts := objectstorage.ObjectFilterOptions{Prefix: "documents/"}
objects, err := osClient.Objects().ListAll(context.Background(), "my-bucket", opts)
fmt.Printf("Total objects: %d\n", len(objects))
```

##### Deleting an Object

```go
err := osClient.Objects().Delete(context.Background(), "my-bucket", "hello.txt", nil)
```

Delete a specific version:

```go
opts := &objectstorage.DeleteOptions{VersionID: "version-id"}
err := osClient.Objects().Delete(context.Background(), "my-bucket", "hello.txt", opts)
```

##### Getting Object Metadata

```go
metadata, err := osClient.Objects().Metadata(context.Background(), "my-bucket", "hello.txt")
if metadata != nil {
    fmt.Printf("Size: %d, Content-Type: %s\n", metadata.Size, metadata.ContentType)
}
```

##### Object Locking

Lock an object with retention:

```go
retainUntil := time.Now().Add(24 * time.Hour)
err := osClient.Objects().LockObject(context.Background(), "my-bucket", "important.txt", retainUntil)
```

Unlock an object:

```go
err := osClient.Objects().UnlockObject(context.Background(), "my-bucket", "important.txt")
```

Check object lock status:

```go
locked, err := osClient.Objects().GetObjectLockStatus(context.Background(), "my-bucket", "important.txt")
if locked {
    fmt.Println("Object is locked")
}
```

##### Versioning

List versions of an object:

```go
versions, err := osClient.Objects().ListVersions(context.Background(), "my-bucket", "hello.txt", nil)
for _, version := range versions {
    fmt.Printf("Version: %s (Size: %d)\n", version.VersionID, version.Size)
}
```

### Initializing the Client

```go
import (
    "github.com/MagaluCloud/mgc-sdk-go/client"
    "github.com/MagaluCloud/mgc-sdk-go/compute"
)

apiToken := os.Getenv("MGC_API_TOKEN")
c := client.NewMgcClient(client.WithAPIKey(apiToken))
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

### Pagination Helpers (ListAll)

Many services expose a convenience `ListAll` method that transparently walks through all paginated results and returns a single in‑memory slice. Use these helpers when you need the full dataset and do not require manual pagination control.

How it works:

- Starts at `_offset=0` with a service-specific default `_limit`
- Repeats requests, incrementing offset by the page size
- Stops when a page returns fewer items than the limit (or zero)
- Accumulates results into a single slice which is returned

Typical signatures:

- `Service().List(ctx, ListOptions{ Limit, Offset, ... })` // single page
- `Service().ListAll(ctx, FilterOptions{ ... })` // all pages (no pagination fields)

Filter option structs intentionally remove `Limit` / `Offset` and keep only filterable fields (status, engine ID, type, source ID, expand, etc.).

Example (Compute): list all images that match a filter

```go
computeClient := compute.New(c)
images, err := computeClient.Images().ListAll(
    context.Background(),
    compute.ImageFilterOptions{
        // Example filter fields (adjust as available)
        Name: helpers.StrPtr("ubuntu"),
    },
)
if err != nil {
    log.Fatalf("list all images: %v", err)
}
fmt.Printf("Fetched %d images\n", len(images))
```

When to use which:

- Use `List` if you need streaming/partial processing, custom limits, pagination UI, or to avoid loading large datasets entirely into memory.
- Use `ListAll` for simplicity when result counts are manageable or for setup/administrative scripts.

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

### Advanced HTTP Client Usage

For more advanced use cases, you can directly use the `mgc_http` package. This is useful when you need to interact with API endpoints that are not yet fully supported by the SDK.

#### `ExecuteSimpleRequestWithRespBody`

Use this function when you expect a response body that needs to be unmarshaled into a struct.

```go
import (
    "context"
    "github.com/MagaluCloud/mgc-sdk-go/internal/http"
    "github.com/MagaluCloud/mgc-sdk-go/client"
)

type MyResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// Create a new request function
newRequest := func(ctx context.Context, method, path string, body any) (*http.Request, error) {
    return mgc_http.NewRequest(myClient.Config, ctx, method, path, body)
}

// Make the request
var myResponse MyResponse
_, err := mgc_http.ExecuteSimpleRequestWithRespBody(
    context.Background(),
    newRequest,
    myClient.Config,
    "GET",
    "/my-endpoint",
    nil, // No request body
    nil, // No query parameters
    &myResponse,
)
```

#### `ExecuteSimpleRequest`

Use this function for requests where you don't expect a response body, such as `DELETE` requests or other operations that return a `204 No Content` status.

```go
err := mgc_http.ExecuteSimpleRequest(
    context.Background(),
    newRequest,
    myClient.Config,
    "DELETE",
    "/my-endpoint/123",
    nil, // No request body
    nil, // No query parameters
)
```

#### Handling Query Parameters

To add query parameters to your request, you can use the `helpers.NewQueryParams` function.

```go
import "github.com/MagaluCloud/mgc-sdk-go/helpers"

// ...

req, _ := http.NewRequest("GET", "", nil)
queryParams := helpers.NewQueryParams(req)
queryParams.Add("name", helpers.StrPtr("my-resource"))
queryParams.Add("limit", helpers.IntPtr(10))

// The query parameters will be encoded and added to the request URL
_, err := mgc_http.ExecuteSimpleRequestWithRespBody(
    // ...
    queryParams,
    // ...
)
```

#### Handling Optional Fields in JSON

To prevent optional fields from being included in the JSON payload when they are not set, use the `omitempty` tag in your struct definitions. This is standard Go practice and is fully supported by the SDK.

```go
type MyRequest struct {
    RequiredField string  `json:"required_field"`
    OptionalField *string `json:"optional_field,omitempty"`
}

// Create a request with only the required field
reqBody := MyRequest{
    RequiredField: "some-value",
}

// The resulting JSON will be: {"required_field":"some-value"}
// The optional_field will be omitted.
```

## Full Example

Check the [cmd/examples](cmd/examples) directory for complete working examples of all SDK features.

## Contributing

We welcome contributions from the community! Please read our [Contributing Guide](CONTRIBUTING.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
