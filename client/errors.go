// Package client provides the core client functionality for the MagaluCloud SDK.
// This package contains the main client implementation, configuration options, and error handling.
package client

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPError represents an error that occurred during an HTTP request.
// This error type includes the HTTP status code, status message, and response body.
type HTTPError struct {
	// StatusCode is the HTTP status code from the response
	StatusCode int
	// Status is the HTTP status message
	Status string
	// Body contains the response body as bytes
	Body []byte
	// Response is the original HTTP response object
	Response *http.Response
}

// Error returns a string representation of the HTTP error.
// This method implements the error interface.
//
// Returns:
//   - string: A formatted error message including status code and status
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d %s", e.StatusCode, e.Status)
}

// NewHTTPError creates a new HTTPError from an HTTP response.
// This function reads the response body and creates an error with all available information.
//
// Parameters:
//   - resp: The HTTP response that caused the error
//
// Returns:
//   - *HTTPError: A new HTTP error instance
func NewHTTPError(resp *http.Response) *HTTPError {
	body, _ := io.ReadAll(resp.Body)
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
		Response:   resp,
	}
}

// ValidationError represents an error that occurred during input validation.
// This error type includes the field that failed validation and a descriptive message.
type ValidationError struct {
	// Field is the name of the field that failed validation
	Field string
	// Message is a descriptive error message
	Message string
}

// Error returns a string representation of the validation error.
// This method implements the error interface.
//
// Returns:
//   - string: A formatted error message including field name and validation message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// RetryError represents an error that occurred after exhausting all retry attempts.
// This error type includes the last error encountered and the number of retries attempted.
type RetryError struct {
	// LastError is the error from the final retry attempt
	LastError error
	// Retries is the number of retry attempts that were made
	Retries int
}

// Error returns a string representation of the retry error.
// This method implements the error interface.
//
// Returns:
//   - string: A formatted error message indicating max retry attempts were reached
func (e *RetryError) Error() string {
	return fmt.Sprintf("max retry attempts reached: %v", e.LastError)
}
