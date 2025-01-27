package client

type MgcUrl string

const (
	// BrNe1 is the URL for the Brazil Northeast 1 region
	BrNe1 MgcUrl = "https://api.magalu.cloud/br-ne1"
	// BrSe1 is the URL for the Brazil Southeast 1 region
	BrSe1 MgcUrl = "https://api.magalu.cloud/br-se1"
	// BrMgl1 is the URL for the Brazil Magalu region
	BrMgl1 MgcUrl = "https://api.magalu.cloud/br-se-1"
)

// String returns the string representation of the MgcUrl
func (m MgcUrl) String() string {
	return string(m)
}