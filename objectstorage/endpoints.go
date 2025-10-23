package objectstorage

import "fmt"

// Endpoint represents a MagaluObjects endpoint.
type Endpoint string

const (
	// BrSe1 is the Brazil Southeast 1 region endpoint.
	BrSe1 Endpoint = "https://br-se1.magaluobjects.com"

	// BrNe1 is the Brazil Northeast 1 region endpoint.
	BrNe1 Endpoint = "https://br-ne1.magaluobjects.com"
)

// String returns the string representation of the endpoint.
func (e Endpoint) String() string {
	return string(e)
}

// IsValid checks if the endpoint is valid.
func (e Endpoint) IsValid() bool {
	switch e {
	case BrSe1, BrNe1:
		return true
	default:
		return false
	}
}

// ValidateEndpoint validates an endpoint and returns an error if invalid.
func ValidateEndpoint(e Endpoint) error {
	if !e.IsValid() {
		return fmt.Errorf("invalid endpoint: %s", e)
	}
	return nil
}
