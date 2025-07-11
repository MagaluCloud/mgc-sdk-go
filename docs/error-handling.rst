Error Handling
=============

The MGC Go SDK provides comprehensive error handling with detailed error types and information. This section covers how to handle different types of errors that may occur when using the SDK.

Error Types
----------

The SDK defines several error types to help you handle different scenarios:

HTTPError
~~~~~~~~~

`HTTPError` contains details about API request failures:

.. code-block:: go

   type HTTPError struct {
       StatusCode int         // HTTP status code
       Status     string      // Status message
       Body       []byte      // Raw response body
       Response   *http.Response
   }

ValidationError
~~~~~~~~~~~~~~~

`ValidationError` occurs when request parameters are invalid:

.. code-block:: go

   type ValidationError struct {
       Field   string   // Which field failed validation
       Message string   // Why the validation failed
   }

Basic Error Handling
-------------------

Check for errors after each operation:

.. code-block:: go

   instances, err := computeClient.Instances().List(ctx, compute.ListOptions{})
   if err != nil {
       log.Printf("Failed to list instances: %v", err)
       return err
   }

HTTP Error Handling
------------------

Handle specific HTTP status codes:

.. code-block:: go

   err := computeClient.Instances().Delete(ctx, id)
   if httpErr, ok := err.(*client.HTTPError); ok {
       switch httpErr.StatusCode {
       case 404:
           log.Fatal("Instance not found")
       case 403:
           log.Fatal("Permission denied")
       case 429:
           log.Fatal("Rate limit exceeded")
       case 500, 502, 503, 504:
           log.Fatal("Server error, try again later")
       default:
           log.Printf("HTTP %d: %s", httpErr.StatusCode, string(httpErr.Body))
       }
   }

Validation Error Handling
------------------------

Handle validation failures:

.. code-block:: go

   _, err := computeClient.Instances().Create(ctx, compute.CreateRequest{})
   if validErr, ok := err.(*client.ValidationError); ok {
       log.Printf("Invalid field %s: %s", validErr.Field, validErr.Message)
   }

Comprehensive Error Handling
---------------------------

Use type switches for comprehensive error handling:

.. code-block:: go

   err := computeClient.Instances().Delete(ctx, id)
   switch e := err.(type) {
   case *client.HTTPError:
       // Handle HTTP errors (404, 403, etc)
       fmt.Printf("HTTP %d: %s\n", e.StatusCode, e.Body)
   case *client.ValidationError:
       // Handle validation failures
       fmt.Printf("Invalid %s: %s\n", e.Field, e.Message)
   case nil:
       // No error
       fmt.Println("Operation completed successfully")
   default:
       // Handle other errors (context timeout, network issues)
       fmt.Printf("Unexpected error: %v\n", err)
   }

Common HTTP Status Codes
-----------------------

The following HTTP status codes are commonly returned by the Magalu Cloud API:

200 - OK
~~~~~~~~
Request completed successfully.

201 - Created
~~~~~~~~~~~~
Resource created successfully.

400 - Bad Request
~~~~~~~~~~~~~~~~
Invalid request parameters or malformed request.

401 - Unauthorized
~~~~~~~~~~~~~~~~~
Authentication required or invalid API token.

403 - Forbidden
~~~~~~~~~~~~~~
Insufficient permissions to perform the operation.

404 - Not Found
~~~~~~~~~~~~~~
The requested resource does not exist.

409 - Conflict
~~~~~~~~~~~~~
Resource conflict (e.g., trying to create a resource that already exists).

422 - Unprocessable Entity
~~~~~~~~~~~~~~~~~~~~~~~~~
Request validation failed.

429 - Too Many Requests
~~~~~~~~~~~~~~~~~~~~~~
Rate limit exceeded. Wait before making more requests.

500 - Internal Server Error
~~~~~~~~~~~~~~~~~~~~~~~~~~
Server error. Try again later.

502 - Bad Gateway
~~~~~~~~~~~~~~~~
Gateway error. Try again later.

503 - Service Unavailable
~~~~~~~~~~~~~~~~~~~~~~~~
Service temporarily unavailable. Try again later.

504 - Gateway Timeout
~~~~~~~~~~~~~~~~~~~~
Gateway timeout. Try again later.

Retry Logic
----------

The client automatically retries on network errors and 5xx responses:

.. code-block:: go

   client := client.NewMgcClient(
       apiToken,
       client.WithRetryConfig(
           3, // maxAttempts
           1 * time.Second, // initialInterval
           30 * time.Second, // maxInterval
           1.5, // backoffFactor
       ),
   )

Custom Retry Logic
~~~~~~~~~~~~~~~~~

For more control, implement custom retry logic:

.. code-block:: go

   func retryOperation(operation func() error, maxRetries int) error {
       var lastErr error
       
       for i := 0; i < maxRetries; i++ {
           err := operation()
           if err == nil {
               return nil
           }
           
           // Check if error is retryable
           if httpErr, ok := err.(*client.HTTPError); ok {
               if httpErr.StatusCode >= 500 {
                   lastErr = err
                   time.Sleep(time.Duration(i+1) * time.Second)
                   continue
               }
           }
           
           // Don't retry on client errors
           return err
       }
       
       return fmt.Errorf("operation failed after %d retries: %v", maxRetries, lastErr)
   }

Error Context
-------------

Add context to errors for better debugging:

.. code-block:: go

   func createInstance(ctx context.Context, client *compute.Client, req compute.CreateRequest) error {
       id, err := client.Instances().Create(ctx, req)
       if err != nil {
           return fmt.Errorf("failed to create instance %s: %w", req.Name, err)
       }
       
       log.Printf("Created instance with ID: %s", id)
       return nil
   }

Error Logging
------------

Log errors with appropriate levels:

.. code-block:: go

   import "log/slog"

   err := computeClient.Instances().Delete(ctx, id)
   if err != nil {
       if httpErr, ok := err.(*client.HTTPError); ok {
           switch httpErr.StatusCode {
           case 404:
               slog.Info("Instance not found, may have been deleted", "id", id)
           case 403:
               slog.Error("Permission denied", "id", id, "error", err)
           case 429:
               slog.Warn("Rate limit exceeded", "error", err)
           default:
               slog.Error("HTTP error", "status", httpErr.StatusCode, "error", err)
           }
       } else {
           slog.Error("Unexpected error", "error", err)
       }
   }

Error Recovery
-------------

Implement error recovery strategies:

.. code-block:: go

   func safeDeleteInstance(ctx context.Context, client *compute.Client, id string) error {
       err := client.Instances().Delete(ctx, id)
       if err != nil {
           if httpErr, ok := err.(*client.HTTPError); ok {
               if httpErr.StatusCode == 404 {
                   // Instance already deleted, consider it a success
                   return nil
               }
           }
           return err
       }
       return nil
   }

Best Practices
-------------

1. **Always check for errors**
   - Don't ignore error return values
   - Handle errors appropriately for your use case

2. **Use specific error types**
   - Check for `HTTPError` and `ValidationError` types
   - Handle different status codes appropriately

3. **Provide context**
   - Wrap errors with additional context
   - Include relevant information in error messages

4. **Log appropriately**
   - Use appropriate log levels (info, warn, error)
   - Include relevant context in log messages

5. **Implement retry logic**
   - Use the built-in retry configuration
   - Implement custom retry logic for specific scenarios

6. **Handle rate limits**
   - Check for 429 status codes
   - Implement exponential backoff

7. **Validate input**
   - Validate parameters before making requests
   - Handle validation errors gracefully

Example: Complete Error Handling
-------------------------------

Here's a complete example showing comprehensive error handling:

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "log"
       "time"

       "github.com/MagaluCloud/mgc-sdk-go/client"
       "github.com/MagaluCloud/mgc-sdk-go/compute"
       "github.com/MagaluCloud/mgc-sdk-go/helpers"
   )

   func main() {
       // Initialize client with retry configuration
       c := client.NewMgcClient(
           os.Getenv("MGC_API_TOKEN"),
           client.WithRetryConfig(3, 1*time.Second, 30*time.Second, 1.5),
       )

       computeClient := compute.New(c)
       ctx := context.Background()

       // Create an instance with error handling
       err := createInstanceWithRetry(ctx, computeClient, "test-instance")
       if err != nil {
           log.Fatalf("Failed to create instance: %v", err)
       }
   }

   func createInstanceWithRetry(ctx context.Context, client *compute.Client, name string) error {
       req := compute.CreateRequest{
           Name: name,
           MachineType: compute.IDOrName{
               Name: helpers.StrPtr("BV1-1-40"),
           },
           Image: compute.IDOrName{
               Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
           },
       }

       // Try to create the instance
       id, err := client.Instances().Create(ctx, req)
       if err != nil {
           return handleCreateError(err, name)
       }

       log.Printf("Successfully created instance %s with ID: %s", name, id)
       return nil
   }

   func handleCreateError(err error, name string) error {
       switch e := err.(type) {
       case *client.HTTPError:
           switch e.StatusCode {
           case 400:
               return fmt.Errorf("invalid request for instance %s: %s", name, string(e.Body))
           case 403:
               return fmt.Errorf("insufficient permissions to create instance %s", name)
           case 409:
               return fmt.Errorf("instance %s already exists", name)
           case 422:
               return fmt.Errorf("validation failed for instance %s: %s", name, string(e.Body))
           case 429:
               return fmt.Errorf("rate limit exceeded while creating instance %s", name)
           case 500, 502, 503, 504:
               return fmt.Errorf("server error while creating instance %s, try again later", name)
           default:
               return fmt.Errorf("HTTP %d error creating instance %s: %s", e.StatusCode, name, string(e.Body))
           }
       case *client.ValidationError:
           return fmt.Errorf("validation error for instance %s: field %s - %s", name, e.Field, e.Message)
       default:
           return fmt.Errorf("unexpected error creating instance %s: %v", name, err)
       }
   } 