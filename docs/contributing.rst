Contributing
============

Thank you for your interest in contributing to the MGC Go SDK! This document provides guidelines and information for contributors.

Getting Started
--------------

Prerequisites
~~~~~~~~~~~~

- Go 1.21 or higher
- Git
- A Magalu Cloud account (for testing)

Fork and Clone
~~~~~~~~~~~~~

1. Fork the repository on GitHub
2. Clone your fork locally:

   .. code-block:: bash

      git clone https://github.com/YOUR_USERNAME/mgc-sdk-go.git
      cd mgc-sdk-go

3. Add the upstream repository:

   .. code-block:: bash

      git remote add upstream https://github.com/MagaluCloud/mgc-sdk-go.git

Development Setup
----------------

Install Dependencies
~~~~~~~~~~~~~~~~~~~

.. code-block:: bash

   go mod download
   go mod tidy

Run Tests
~~~~~~~~~

.. code-block:: bash

   # Run all tests
   go test ./...

   # Run tests with coverage
   go test -cover ./...

   # Run tests with verbose output
   go test -v ./...

   # Run tests for a specific package
   go test ./compute

Code Quality
-----------

Code Style
~~~~~~~~~~

The project follows Go's standard formatting and style guidelines:

- Use `gofmt` for code formatting
- Follow the `golint` guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types

Run formatting:

.. code-block:: bash

   go fmt ./...

Linting
~~~~~~~

Install and run linters:

.. code-block:: bash

   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # Run linter
   golangci-lint run

Documentation
~~~~~~~~~~~~

- Add documentation for all exported functions and types
- Follow Go documentation conventions
- Include examples for complex functionality
- Update README.md if adding new features

Testing
-------

Writing Tests
~~~~~~~~~~~~

- Write tests for all new functionality
- Use descriptive test names
- Test both success and error cases
- Mock external dependencies when appropriate

Test Structure
~~~~~~~~~~~~~

Follow this pattern for tests:

.. code-block:: go

   func TestFunctionName(t *testing.T) {
       // Setup
       client := createTestClient()
       
       // Test cases
       tests := []struct {
           name    string
           input   string
           want    string
           wantErr bool
       }{
           {
               name:    "success case",
               input:   "test",
               want:    "expected",
               wantErr: false,
           },
           {
               name:    "error case",
               input:   "",
               want:    "",
               wantErr: true,
           },
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               got, err := FunctionName(client, tt.input)
               if (err != nil) != tt.wantErr {
                   t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                   return
               }
               if got != tt.want {
                   t.Errorf("FunctionName() = %v, want %v", got, tt.want)
               }
           })
       }
   }

Integration Tests
~~~~~~~~~~~~~~~~

For integration tests that require a real API:

.. code-block:: go

   func TestIntegration(t *testing.T) {
       if testing.Short() {
           t.Skip("skipping integration test")
       }
       
       // Your integration test code here
   }

Run integration tests:

.. code-block:: bash

   go test -tags=integration ./...
   # or
   go test -short=false ./...

Adding New Services
------------------

When adding support for a new Magalu Cloud service:

1. Create a new package in the root directory
2. Follow the existing service structure
3. Implement the service interface
4. Add comprehensive tests
5. Update documentation

Service Structure
~~~~~~~~~~~~~~~~

Follow this structure for new services:

.. code-block:: go

   package newservice

   import (
       "context"
       "github.com/MagaluCloud/mgc-sdk-go/client"
   )

   // NewServiceClient represents a client for the new service
   type NewServiceClient struct {
       *client.CoreClient
   }

   // New creates a new NewServiceClient
   func New(core *client.CoreClient) *NewServiceClient {
       return &NewServiceClient{CoreClient: core}
   }

   // ResourceService provides operations for managing resources
   type ResourceService interface {
       List(ctx context.Context, opts ListOptions) ([]Resource, error)
       Get(ctx context.Context, id string) (*Resource, error)
       Create(ctx context.Context, req CreateRequest) (string, error)
       Delete(ctx context.Context, id string) error
   }

   // Resource represents a resource in the service
   type Resource struct {
       ID   string `json:"id"`
       Name string `json:"name"`
       // ... other fields
   }

   // CreateRequest represents the parameters for creating a resource
   type CreateRequest struct {
       Name string `json:"name"`
       // ... other fields
   }

   // ListOptions represents parameters for listing resources
   type ListOptions struct {
       Limit  *int    `json:"limit,omitempty"`
       Offset *int    `json:"offset,omitempty"`
       Sort   *string `json:"sort,omitempty"`
   }

   // Implementation
   type resourceService struct {
       client *NewServiceClient
   }

   func (s *NewServiceClient) Resources() ResourceService {
       return &resourceService{client: s}
   }

   // Implement the interface methods...
```

Adding New Features
------------------

When adding new features to existing services:

1. Follow the existing patterns
2. Add proper error handling
3. Include comprehensive tests
4. Update documentation
5. Consider backward compatibility

Breaking Changes
~~~~~~~~~~~~~~~

If your changes include breaking changes:

1. Create a new major version
2. Document the changes clearly
3. Provide migration guides
4. Update examples

Pull Request Process
-------------------

1. Create a feature branch:

   .. code-block:: bash

      git checkout -b feature/your-feature-name

2. Make your changes and commit them:

   .. code-block:: bash

      git add .
      git commit -m "feat: add new feature description"

3. Push your branch:

   .. code-block:: bash

      git push origin feature/your-feature-name

4. Create a Pull Request on GitHub

Commit Message Format
~~~~~~~~~~~~~~~~~~~~

Follow the conventional commit format:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `style:` for formatting changes
- `refactor:` for code refactoring
- `test:` for adding tests
- `chore:` for maintenance tasks

Examples:

.. code-block:: text

   feat: add support for new API endpoint
   fix: handle 404 errors properly
   docs: update installation instructions
   test: add integration tests for compute service

Pull Request Guidelines
~~~~~~~~~~~~~~~~~~~~~~

- Provide a clear description of the changes
- Include tests for new functionality
- Update documentation if needed
- Ensure all tests pass
- Follow the code style guidelines

Code Review
-----------

All contributions require code review:

1. Address review comments promptly
2. Make requested changes
3. Ensure CI checks pass
4. Keep the PR up to date with main branch

CI/CD Pipeline
--------------

The project uses GitHub Actions for continuous integration:

- **Code Quality**: Runs linters and code analysis
- **Unit Tests**: Runs all unit tests
- **Integration Tests**: Runs integration tests (if applicable)
- **Documentation**: Builds and validates documentation

Local Development
----------------

Running Examples
~~~~~~~~~~~~~~~

Test your changes with the provided examples:

.. code-block:: bash

   cd cmd/examples/compute
   go run main.go

Building Documentation
~~~~~~~~~~~~~~~~~~~~~

Build the documentation locally:

.. code-block:: bash

   # Install Sphinx
   pip install sphinx sphinx-rtd-theme

   # Build docs
   cd docs
   make html

   # View docs
   open _build/html/index.html

Debugging
---------

Enable Debug Logging
~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   import "log/slog"

   logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
       Level: slog.LevelDebug,
   }))

   c := client.NewMgcClient(
       apiToken,
       client.WithLogger(logger),
   )

Request Tracing
~~~~~~~~~~~~~~

Use request IDs for tracing:

.. code-block:: go

   import (
       "context"
       "github.com/google/uuid"
       "github.com/MagaluCloud/mgc-sdk-go/client"
   )

   requestID := uuid.New().String()
   ctx := context.WithValue(context.Background(), client.RequestIDKey, requestID)

   // The request ID will be included in logs and headers
   instances, err := computeClient.Instances().List(ctx, compute.ListOptions{})

Getting Help
-----------

If you need help with development:

1. Check the existing documentation
2. Look at existing code examples
3. Open an issue on GitHub
4. Join the community discussions

Community Guidelines
-------------------

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Follow the project's code of conduct

License
-------

By contributing to this project, you agree that your contributions will be licensed under the Apache 2.0 License.

Thank you for contributing to the MGC Go SDK! 