Configuration
=============

The MGC Go SDK provides extensive configuration options to customize client behavior, including timeouts, retry logic, logging, and region selection.

Basic Configuration
------------------

The simplest way to create a client is with just an API token:

.. code-block:: go

   c := client.NewMgcClient(apiToken)

Advanced Configuration
---------------------

You can customize the client behavior using configuration options:

.. code-block:: go

   import (
       "time"
       "log/slog"
       "github.com/MagaluCloud/mgc-sdk-go/client"
   )

   c := client.NewMgcClient(
       apiToken,
       client.WithTimeout(5*time.Minute),
       client.WithUserAgent("my-app/1.0"),
       client.WithLogger(slog.Default().With("service", "mgc")),
       client.WithRetryConfig(
           3, // maxAttempts
           1*time.Second, // initialInterval
           30*time.Second, // maxInterval
           1.5, // backoffFactor
       ),
       client.WithBaseURL(client.BrSe1),
   )

Available Options
----------------

Timeout Configuration
~~~~~~~~~~~~~~~~~~~~

Set a custom timeout for all requests:

.. code-block:: go

   c := client.NewMgcClient(
       apiToken,
       client.WithTimeout(30*time.Second),
   )

User Agent
~~~~~~~~~~

Set a custom User-Agent header:

.. code-block:: go

   c := client.NewMgcClient(
       apiToken,
       client.WithUserAgent("my-application/1.0.0"),
   )

Logging
~~~~~~~

Configure structured logging:

.. code-block:: go

   import "log/slog"

   logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
   c := client.NewMgcClient(
       apiToken,
       client.WithLogger(logger.With("service", "mgc")),
   )

Retry Configuration
~~~~~~~~~~~~~~~~~~

Customize retry behavior:

.. code-block:: go

   c := client.NewMgcClient(
       apiToken,
       client.WithRetryConfig(
           5, // maxAttempts
           2*time.Second, // initialInterval
           60*time.Second, // maxInterval
           2.0, // backoffFactor
       ),
   )

Custom HTTP Client
~~~~~~~~~~~~~~~~~

Use a custom HTTP client:

.. code-block:: go

   import "net/http"

   httpClient := &http.Client{
       Timeout: 30 * time.Second,
       Transport: &http.Transport{
           MaxIdleConns: 100,
           IdleConnTimeout: 90 * time.Second,
       },
   }

   c := client.NewMgcClient(
       apiToken,
       client.WithHTTPClient(httpClient),
   )

Custom Headers
~~~~~~~~~~~~~

Add custom headers to all requests:

.. code-block:: go

   c := client.NewMgcClient(
       apiToken,
       client.WithCustomHeader("X-Custom-Header", "custom-value"),
   )

Regions
-------

Magalu Cloud operates in multiple regions. You can specify which region to use:

Brazil South East 1 (BR-SE1) - Default
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   c := client.NewMgcClient(
       apiToken,
       client.WithBaseURL(client.BrSe1),
   )

Brazil North East 1 (BR-NE1)
~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   c := client.NewMgcClient(
       apiToken,
       client.WithBaseURL(client.BrNe1),
   )

Global Services
--------------

Some services operate globally and ignore region settings:

- SSH Keys Management
- Availability Zones

These services automatically use the global endpoint (api.magalu.cloud):

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/sshkeys"

   // Even if core client has a region set
   c := client.NewMgcClient(apiToken, client.WithBaseURL(client.BrSe1))

   // Global services will ignore the region and use global endpoint
   sshClient := sshkeys.New(c) // Uses api.magalu.cloud

To use a custom endpoint for global services:

.. code-block:: go

   sshClient := sshkeys.New(c, sshkeys.WithGlobalBasePath("custom-endpoint"))
```

Request Tracking
---------------

You can track requests across systems using request IDs:

.. code-block:: go

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

Request ID requirements:

- Must be a valid UUIDv4 string (e.g., "123e4567-e89b-12d3-a456-426614174000")
- Included in the request as `X-Request-ID` header
- Logged in the client's logger
- Returned in the response headers for tracking

Configuration Examples
---------------------

Production Configuration
~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   func createProductionClient(apiToken string) *client.CoreClient {
       return client.NewMgcClient(
           apiToken,
           client.WithTimeout(60*time.Second),
           client.WithUserAgent("production-app/1.0.0"),
           client.WithLogger(slog.Default().With("env", "production")),
           client.WithRetryConfig(5, 2*time.Second, 60*time.Second, 2.0),
           client.WithBaseURL(client.BrSe1),
       )
   }

Development Configuration
~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   func createDevelopmentClient(apiToken string) *client.CoreClient {
       return client.NewMgcClient(
           apiToken,
           client.WithTimeout(30*time.Second),
           client.WithUserAgent("dev-app/1.0.0"),
           client.WithLogger(slog.Default().With("env", "development")),
           client.WithRetryConfig(3, 1*time.Second, 10*time.Second, 1.5),
           client.WithBaseURL(client.BrSe1),
       )
   }

Testing Configuration
~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   func createTestClient(apiToken string) *client.CoreClient {
       return client.NewMgcClient(
           apiToken,
           client.WithTimeout(10*time.Second),
           client.WithUserAgent("test-app/1.0.0"),
           client.WithRetryConfig(1, 100*time.Millisecond, 1*time.Second, 1.0),
       )
   }

Environment-Based Configuration
-----------------------------

You can create different configurations based on environment variables:

.. code-block:: go

   func createClient(apiToken string) *client.CoreClient {
       env := os.Getenv("APP_ENV")
       
       switch env {
       case "production":
           return client.NewMgcClient(
               apiToken,
               client.WithTimeout(60*time.Second),
               client.WithRetryConfig(5, 2*time.Second, 60*time.Second, 2.0),
               client.WithBaseURL(client.BrSe1),
           )
       case "development":
           return client.NewMgcClient(
               apiToken,
               client.WithTimeout(30*time.Second),
               client.WithRetryConfig(3, 1*time.Second, 10*time.Second, 1.5),
               client.WithBaseURL(client.BrSe1),
           )
       default:
           return client.NewMgcClient(
               apiToken,
               client.WithTimeout(10*time.Second),
               client.WithRetryConfig(1, 100*time.Millisecond, 1*time.Second, 1.0),
           )
       }
   }

Next Steps
----------

After configuring your client:

1. Start using the services (see :doc:`services/index`)
2. Learn about error handling (see :doc:`error-handling`)
3. Check out examples (see :doc:`examples`) 