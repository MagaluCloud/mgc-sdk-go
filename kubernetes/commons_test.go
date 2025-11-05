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

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
