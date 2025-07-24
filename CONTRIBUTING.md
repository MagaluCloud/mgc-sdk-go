## Contributing to the mgc-sdk-go

First off, thank you for considering contributing to the **mgc-sdk-go**! It's people like you that make the open source community such a great place.

This guide is designed to help you get started, whether you're a seasoned Go developer or just starting out. It covers everything from setting up your environment to creating new service modules and submitting your work.

---

### Learning Go

If you are new to Go, here are some resources to help you get started:

- [A Tour of Go](https://go.dev/tour/welcome/1): An interactive introduction to the basics of Go.
- [Go by Example](https://gobyexample.com/): A hands-on introduction to Go using annotated example programs.
- [Effective Go](https://go.dev/doc/effective_go): Official documentation with tips for writing clear, idiomatic Go code.
- [Go Modules Reference](https://go.dev/ref/mod): A detailed reference for using Go modules.

---

### 1. Getting Started: Your Development Environment

1.  **Install Go**: Ensure you have Go version 1.22 or higher. You can find the latest version on the [official Go download page](https://go.dev/dl/).

2.  **Fork and Clone the Repository**:

    ```bash
    # Clone your fork (replace <your-username> with your GitHub username)
    git clone git@github.com:<your-username>/mgc-sdk-go.git
    cd mgc-sdk-go
    # Add the original repository as an upstream remote
    git remote add upstream https://github.com/MagaluCloud/mgc-sdk-go.git
    ```

3.  **Run the Test Suite**: Make sure everything is working correctly from the start.

    ```bash
    go test ./...
    ```

4.  **Code Formatting and Linting**: We use standard Go tools to keep the code clean and consistent.

    ```bash
    go vet ./...
    go fmt ./...
    ```

---

### 2. Project Structure Overview

```
mgc-sdk-go/
├── client/           # Generic CoreClient for API interaction
├── compute/          # Compute service module
├── blockstorage/     # Block Storage service module
├── helpers/          # Utility functions
├── internal/         # Internal packages, including test data and mocks
└── cmd/examples/     # Usage examples for each service
```

Each service-specific subdirectory (e.g., `compute/`, `blockstorage/`) contains:

- `client.go`: Defines the service-specific client (`type Client struct{ *client.CoreClient }`).
- `service_*.go`: Implements the API methods for the service (e.g., List, Get, Create).
- `types.go`: Contains the request and response data transfer objects (DTOs).
- `service_test.go`: Includes unit tests for the service.

---

### 3. Creating a New Service Module

Let's say you want to add support for a new **Message Queue** service.

1.  **Create a New Branch**:

    ```bash
    git checkout -b feature/mq-service
    ```

2.  **Create the Initial Directory Structure**:

    ```
    mgc-sdk-go/
    └── mq/
        ├── client.go
        ├── queues.go
        └── types.go
    ```

3.  **Implement `client.go`**:

    The `client.go` file is the entry point for the service. It defines a `Client` struct that embeds the `client.CoreClient`, giving it access to the underlying HTTP client and configuration.

    ```go
    package mq

    import "github.com/MagaluCloud/mgc-sdk-go/client"

    // DefaultBasePath is the default base path for the Message Queue API.
    const DefaultBasePath = "/mq"

    // Client is the client for the Message Queue service.
    type Client struct {
        *client.CoreClient
    }

    // New creates a new Message Queue client.
    func New(core *client.CoreClient) *Client {
        return &Client{CoreClient: core.WithBasePath(DefaultBasePath)}
    }
    ```

4.  **Implement the Service (`queues.go`)**:

    The `queues.go` file defines the service's interface and implements its methods. The `Service` interface defines the operations that can be performed on the service's resources (in this case, `queues`).

    ```go
    package mq

    import "context"

    // Service is the interface for the Message Queue service.
    type Service interface {
        List(ctx context.Context, opts ListOptions) ([]Queue, error)
    }

    // queuesService is the implementation of the Service interface.
    type queuesService struct{ c *Client }

    // Queues returns a new queuesService.
    func (c *Client) Queues() Service { return &queuesService{c} }

    // List lists all queues.
    func (s *queuesService) List(ctx context.Context, opts ListOptions) ([]Queue, error) {
        var resp ListResponse
        // The Get method is provided by the CoreClient and handles the underlying HTTP requests and retries.
        err := s.c.Get(ctx, "/queues", opts, &resp)
        return resp.Results, err
    }
    ```

    _Note: The `Get` method is provided by the `CoreClient` and handles the underlying HTTP requests and retries._

5.  **Define the Types (`types.go`)**:

    The `types.go` file defines the data structures used by the service, including request and response objects.

    ```go
    package mq

    // Queue represents a message queue.
    type Queue struct {
        ID   string `json:"queue_id"`
        Name string `json:"name"`
    }

    // ListOptions specifies the optional parameters to the List method.
    type ListOptions struct {
        Limit  *int    `json:"limit,omitempty"`
        Offset *int    `json:"offset,omitempty"`
        Filter *string `json:"filter,omitempty"`
    }

    // ListResponse is the response from the List method.
    type ListResponse struct {
        Results []Queue `json:"results"`
    }
    ```

6.  **Update the `README.md`**: Add the new service to the "Supported Products" section.

7.  **Run Tests and Submit**: Run `go vet`, `go test ./...`, commit your changes, and open a pull request against the `main` branch.

---

### 4. Updating an Existing Module

For example, to add a `Status` field to the `compute.Instance` type:

1.  **Create a New Branch**:

    ```bash
    git checkout -b fix/compute-instance-status
    ```

2.  **Update `compute/types.go`**:

    ```go
    type Instance struct {
        // ... existing fields
        Status *string `json:"status,omitempty"`
    }
    ```

3.  **Update `CreateRequest` (if applicable)**: If the new field can be set on creation, add it to the `CreateRequest` struct.

4.  **Update Tests**: Modify any tests that validate JSON marshaling/unmarshaling to include the new field.

5.  **Run Tests for the Module**:

    ```bash
    go test ./compute
    ```

---

### 5. Unit Testing

#### 5.1 Our Approach

- Tests are written using the **table-driven** methodology.
- We use the `httptest.Server` to mock the API.
- Response JSON payloads are stored in `internal/testdata/*.json`.

#### 5.2 Test Example

```go
package compute_test

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/MagaluCloud/mgc-sdk-go/client"
    "github.com/MagaluCloud/mgc-sdk-go/compute"
)

func TestInstancesList(t *testing.T) {
    // Create a new test server to mock the API.
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set the content type and write the response body.
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"results":[{"instance_id":"i-1","name":"demo"}]}`))
    }))
    // Close the server when the test is finished.
    defer ts.Close()

    // Create a new MgcClient with the test server's URL.
    core := client.NewMgcClient("dummy", client.WithBaseURL(client.MgcUrl(ts.URL)))
    // Create a new compute service client.
    svc := compute.New(core).Instances()

    // Call the List method with a background context and empty options.
    got, err := svc.List(context.Background(), compute.ListOptions{})
    // Check for errors.
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // Check if the response is valid.
    if len(got) != 1 || got[0].ID != "i-1" {
        t.Fatalf("invalid response: %+v", got)
    }
}
```

**Run the full test suite**:

```bash
go test ./...
```

**Generate a coverage report**:

```bash
go test ./... -coverprofile /tmp/cover.out
go tool cover -func=/tmp/cover.out
```

---

### 6. Pull Request Guidelines

- **Small, Focused Commits**: Use imperative mood in your commit messages (e.g., "Add support for X" instead of "Added support for X").
- **Clear PR Description**: Include a changelog entry and a link to the relevant API endpoint documentation.
- **Pass All Checks**: Ensure your code passes `go vet`, `go test`, and any other linters (gofumpt, revive) and SonarCloud checks.
- **Document Breaking Changes**: If your changes are not backward-compatible, document them in the `docs/` directory.

---

### 7. Keeping Your Fork in Sync

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

---

### 8. Using the Makefile

The repository includes a `Makefile` to automate common tasks.

| Target          | Description                                                             |
| --------------- | ----------------------------------------------------------------------- |
| `test`          | Run all tests (`go test ./...`).                                        |
| `test-coverage` | Generate a test coverage report in text and HTML format.                |
| `test-race`     | Run the test suite with the data race detector enabled.                 |
| `test-all`      | Install **gotestsum** and run tests with an aggregated coverage report. |
| `go-fmt`        | Format the code (`gofmt -s -l -w .`).                                   |
| `go-vet`        | Run static analysis checks (`go vet ./...`).                            |
| `readthedocs`   | Build the Sphinx documentation and open it in your browser.             |

#### Daily Workflow

1.  **Before Committing**:

    ```bash
    make go-fmt go-vet test
    ```

2.  **Checking Coverage**:

    ```bash
    make test-coverage
    open coverage.html
    ```

3.  **Hunting for Data Races**:

    ```bash
    make test-race
    ```

---

### Quick Checklist

- [ ] Set up your environment.
- [ ] Created a new branch for your feature or fix.
- [ ] Implemented your changes.
- [ ] Added tests for success and error scenarios.
- [ ] `go vet` and `go test ./...` pass without errors.
- [ ] Opened a pull request that follows the contribution template.

You're now ready to contribute to the **mgc-sdk-go**!
