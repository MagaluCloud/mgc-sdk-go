package lbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func testBackendTargetClient(baseURL string) NetworkBackendTargetService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkBackendTargets()
}

func TestNetworkBackendTargetService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		backendID  string
		request    CreateNetworkBackendTargetRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name:      "successful creation",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkBackendTargetRequest{
				TargetsType: "instance",
				Targets: []NetworkBackendInstanceTargetRequest{
					{
						NicID: stringPtr("nic-1"),
						Port:  80,
					},
				},
			},
			response:   `{"id": "target-123"}`,
			statusCode: http.StatusOK,
			want:       "target-123",
			wantErr:    false,
		},
		{
			name:      "server error",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkBackendTargetRequest{
				TargetsType: "instance",
				Targets: []NetworkBackendInstanceTargetRequest{
					{
						NicID: stringPtr("nic-1"),
						Port:  80,
					},
				},
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s/targets", tt.lbID, tt.backendID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendTargetClient(server.URL)
			id, err := client.Create(context.Background(), tt.lbID, tt.backendID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, id)
		})
	}
}

func TestNetworkBackendTargetService_Replace(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		backendID  string
		request    CreateNetworkBackendTargetRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name:      "successful replace",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkBackendTargetRequest{
				TargetsType: "instance",
				Targets: []NetworkBackendInstanceTargetRequest{
					{
						NicID: stringPtr("nic-1"),
						Port:  80,
					},
					{
						NicID: stringPtr("nic-2"),
						Port:  8080,
					},
				},
			},
			response:   `{"id": "target-456"}`,
			statusCode: http.StatusOK,
			want:       "target-456",
			wantErr:    false,
		},
		{
			name:      "server error",
			lbID:      "lb-123",
			backendID: "backend-123",
			request: CreateNetworkBackendTargetRequest{
				TargetsType: "instance",
				Targets: []NetworkBackendInstanceTargetRequest{
					{
						NicID: stringPtr("nic-1"),
						Port:  80,
					},
				},
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s/targets", tt.lbID, tt.backendID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendTargetClient(server.URL)
			id, err := client.Replace(context.Background(), tt.lbID, tt.backendID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, id)
		})
	}
}

func TestNetworkBackendTargetService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		backendID  string
		targetID   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			lbID:       "lb-123",
			backendID:  "backend-123",
			targetID:   "target-123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent target",
			lbID:       "lb-123",
			backendID:  "backend-123",
			targetID:   "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s/targets/%s", tt.lbID, tt.backendID, tt.targetID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testBackendTargetClient(server.URL)
			err := client.Delete(context.Background(), tt.lbID, tt.backendID, tt.targetID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkBackendTargetService_Create_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendTargetClient("http://dummy-url")

	request := CreateNetworkBackendTargetRequest{
		TargetsType: "instance",
		Targets: []NetworkBackendInstanceTargetRequest{
			{
				NicID: stringPtr("nic-1"),
				Port:  80,
			},
		},
	}

	_, err := client.Create(ctx, "lb-123", "backend-123", request)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkBackendTargetService_Replace_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendTargetClient("http://dummy-url")

	request := CreateNetworkBackendTargetRequest{
		TargetsType: "instance",
		Targets: []NetworkBackendInstanceTargetRequest{
			{
				NicID: stringPtr("nic-1"),
				Port:  80,
			},
		},
	}

	_, err := client.Replace(ctx, "lb-123", "backend-123", request)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkBackendTargetService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testBackendTargetClient("http://dummy-url")

	err := client.Delete(ctx, "lb-123", "backend-123", "target-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
