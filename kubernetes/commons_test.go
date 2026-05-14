package kubernetes

import (
	"net/http"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func testClient(baseURL string) *KubernetesClient {
	core := client.NewMgcClient(client.WithJWToken("test-token"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(&http.Client{Timeout: 1 * time.Second}),
	)
	return New(core)
}

//go:fix inline
func intPtr(i int) *int {
	return new(i)
}

//go:fix inline
func strPtr(s string) *string {
	return new(s)
}
