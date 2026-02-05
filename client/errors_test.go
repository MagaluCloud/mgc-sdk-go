package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		status     string
		want       string
	}{
		{
			name:       "400 Bad Request",
			statusCode: 400,
			status:     "400 Bad Request",
			want:       "\nHTTP error:\n Status: 400 Bad Request\n Body: ",
		},
		{
			name:       "404 Not Found",
			statusCode: 404,
			status:     "404 Not Found",
			want:       "\nHTTP error:\n Status: 404 Not Found\n Body: ",
		},
		{
			name:       "500 Internal Server Error",
			statusCode: 500,
			status:     "500 Internal Server Error",
			want:       "\nHTTP error:\n Status: 500 Internal Server Error\n Body: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &HTTPError{
				StatusCode: tt.statusCode,
				Status:     tt.status,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("HTTPError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		status     string
		body       string
	}{
		{
			name:       "with body content",
			statusCode: 400,
			status:     "400 Bad Request",
			body:       "invalid request",
		},
		{
			name:       "empty body",
			statusCode: 500,
			status:     "500 Internal Server Error",
			body:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Status:     tt.status,
				Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
			}

			err := NewHTTPError(resp)

			if err.StatusCode != tt.statusCode {
				t.Errorf("NewHTTPError().StatusCode = %v, want %v", err.StatusCode, tt.statusCode)
			}
			if err.Status != tt.status {
				t.Errorf("NewHTTPError().Status = %v, want %v", err.Status, tt.status)
			}
			if string(err.Body) != tt.body {
				t.Errorf("NewHTTPError().Body = %v, want %v", string(err.Body), tt.body)
			}
			if err.Response != resp {
				t.Errorf("NewHTTPError().Response = %v, want %v", err.Response, resp)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		message string
		want    string
	}{
		{
			name:    "regular validation error",
			field:   "email",
			message: "invalid format",
			want:    "validation error: email - invalid format",
		},
		{
			name:    "empty field",
			field:   "",
			message: "field is required",
			want:    "validation error:  - field is required",
		},
		{
			name:    "empty message",
			field:   "password",
			message: "",
			want:    "validation error: password - ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ValidationError{
				Field:   tt.field,
				Message: tt.message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryError_Error(t *testing.T) {
	tests := []struct {
		name      string
		lastError error
		retries   int
		want      string
	}{
		{
			name:      "with HTTP error",
			lastError: &HTTPError{StatusCode: 500, Status: "500 Internal Server Error"},
			retries:   3,
			want:      "max retry attempts reached: \nHTTP error:\n Status: 500 Internal Server Error\n Body: ",
		},
		{
			name:      "with simple error",
			lastError: fmt.Errorf("connection timeout"),
			retries:   5,
			want:      "max retry attempts reached: connection timeout",
		},
		{
			name:      "with nil error",
			lastError: nil,
			retries:   1,
			want:      "max retry attempts reached: <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RetryError{
				LastError: tt.lastError,
				Retries:   tt.retries,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("RetryError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
