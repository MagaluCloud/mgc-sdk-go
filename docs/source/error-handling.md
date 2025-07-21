---
title: Error Handling
---

# Error Handling

## HTTP Errors

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

## Validation Errors

```go
_, err := computeClient.Instances().Create(ctx, compute.CreateRequest{})
if validErr, ok := err.(*client.ValidationError); ok {
    log.Printf("Invalid field %s: %s", validErr.Field, validErr.Message)
}
```

## Error Types and Interfaces

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

## Retries

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
