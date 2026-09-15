package lbaas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

const scopedProjectID = "550e8400-e29b-41d4-a716-446655440000"

func TestNetworkLoadBalancerService_SendsProjectHeader(t *testing.T) {
	tests := []struct {
		name        string
		callOpts    client.CallOptions
		wantPresent bool
	}{
		{"scoped call", client.CallOptions{ProjectID: scopedProjectID}, true},
		{"default project", client.CallOptions{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var present bool

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.Header.Get(client.ProjectIDHeader)
				_, present = r.Header[client.ProjectIDHeader]
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"id":"lb-1","name":"lb","status":"running"}`))
			}))
			defer server.Close()

			core := client.NewMgcClient(
				client.WithAPIKey("test-api-key"),
				client.WithBaseURL(client.MgcUrl(server.URL)),
			)

			_, err := New(core).NetworkLoadBalancers().Get(context.Background(), tt.callOpts, "lb-1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if present != tt.wantPresent {
				t.Fatalf("header present = %v, want %v", present, tt.wantPresent)
			}
			if tt.wantPresent && got != scopedProjectID {
				t.Errorf("header = %q, want %q", got, scopedProjectID)
			}
		})
	}
}
