package iam

import (
	"net/http"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func testClient(baseURL string) *IAMClient {
	core := client.NewMgcClient(client.WithJWToken("test-token"),
		client.WithHTTPClient(&http.Client{Timeout: 1 * time.Second}),
	)
	return New(core, WithGlobalBasePath(client.MgcUrl(baseURL)))
}

//go:fix inline
func strPtr(s string) *string {
	return new(s)
}
