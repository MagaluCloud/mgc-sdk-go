package client

import (
	"fmt"
	"io"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
	Response   *http.Response
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d %s", e.StatusCode, e.Status)
}

func NewHTTPError(resp *http.Response) *HTTPError {
	body, _ := io.ReadAll(resp.Body)
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
		Response:   resp,
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

type RetryError struct {
	LastError error
	Retries   int
}

func (e *RetryError) Error() string {
	return fmt.Sprintf("max retry attempts reached: %v", e.LastError)
}
