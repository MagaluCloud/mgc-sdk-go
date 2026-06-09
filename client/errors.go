package client

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPError represents an error that occurred during an HTTP request.
// This error type includes the HTTP status code, status message, and response body.
type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
	Response   *http.Response
}

// Error returns a string representation of the HTTP error.
// This method implements the error interface.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("\nHTTP error:\n Status: %s\n Body: %s", e.Status, e.Body)
}

// NewHTTPError creates a new HTTPError from an HTTP response.
// This function reads the response body and creates an error with all available information.
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
	Field   string
	Message string
}

// Error returns a string representation of the validation error.
// This method implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// RetryError represents an error that occurred after exhausting all retry attempts.
// This error type includes the last error encountered and the number of retries attempted.
type RetryError struct {
	LastError error
	Retries   int
}

// Error returns a string representation of the retry error.
// This method implements the error interface.
func (e *RetryError) Error() string {
	return fmt.Sprintf("max retry attempts reached: %v", e.LastError)
}
