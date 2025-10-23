package objectstorage

import (
	"testing"
)

func TestEndpointString(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		expected string
	}{
		{
			name:     "br-se1 endpoint",
			endpoint: BrSe1,
			expected: "https://br-se1.magaluobjects.com",
		},
		{
			name:     "br-ne1 endpoint",
			endpoint: BrNe1,
			expected: "https://br-ne1.magaluobjects.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.endpoint.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestEndpointIsValid(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		expected bool
	}{
		{
			name:     "br-se1 is valid",
			endpoint: BrSe1,
			expected: true,
		},
		{
			name:     "br-ne1 is valid",
			endpoint: BrNe1,
			expected: true,
		},
		{
			name:     "empty endpoint is invalid",
			endpoint: "",
			expected: false,
		},
		{
			name:     "invalid endpoint",
			endpoint: "https://invalid.magaluobjects.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.endpoint.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		endpoint  Endpoint
		wantError bool
	}{
		{
			name:      "br-se1 is valid",
			endpoint:  BrSe1,
			wantError: false,
		},
		{
			name:      "br-ne1 is valid",
			endpoint:  BrNe1,
			wantError: false,
		},
		{
			name:      "empty endpoint is invalid",
			endpoint:  "",
			wantError: true,
		},
		{
			name:      "invalid endpoint",
			endpoint:  "https://invalid.magaluobjects.com",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEndpoint(tt.endpoint)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateEndpoint() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestEndpointConstants(t *testing.T) {
	if BrSe1 != "https://br-se1.magaluobjects.com" {
		t.Errorf("BrSe1 constant has wrong value: %q", BrSe1)
	}

	if BrNe1 != "https://br-ne1.magaluobjects.com" {
		t.Errorf("BrNe1 constant has wrong value: %q", BrNe1)
	}
}
