package client

// MgcUrl represents a MagaluCloud API URL.
// This type is used to ensure type safety when working with API endpoints.
type MgcUrl string

// Predefined MagaluCloud API endpoints for different regions.
const (
	// BrNe1 is the URL for the Brazil Northeast 1 region
	BrNe1 MgcUrl = "https://api.magalu.cloud/br-ne1"
	// BrSe1 is the URL for the Brazil Southeast 1 region
	BrSe1 MgcUrl = "https://api.magalu.cloud/br-se1"
	// BrMgl1 is the URL for the Brazil Magalu region
	BrMgl1 MgcUrl = "https://api.magalu.cloud/br-se-1"
	// Global is the default URL for products that don't have a specific region
	Global MgcUrl = "https://api.magalu.cloud"
)

// String returns the string representation of the MgcUrl.
// This method implements the Stringer interface.
func (m MgcUrl) String() string {
	return string(m)
}
