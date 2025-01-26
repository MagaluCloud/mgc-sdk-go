package client

type MgcUrl string

const (
	BrNe1 MgcUrl = "https://api.magalu.cloud/br-ne-1"
	BrSe1 MgcUrl = "https://api.magalu.cloud/br-se-1"
)

func (m MgcUrl) String() string {
	return string(m)
}