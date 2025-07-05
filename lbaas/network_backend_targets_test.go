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
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkBackends().Targets()
}

func TestNetworkBackendTargetService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		request    CreateNetworkBackendTargetRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateNetworkBackendTargetRequest{
				LoadBalancerID:   "lb-123",
				NetworkBackendID: "backend-123",
				TargetsID:        []string{"target-1", "target-2"},
				TargetsType:      "instance",
			},
			response:   `{"id": "target-123"}`,
			statusCode: http.StatusOK,
			want:       "target-123",
			wantErr:    false,
		},
		{
			name: "server error",
			request: CreateNetworkBackendTargetRequest{
				LoadBalancerID:   "lb-123",
				NetworkBackendID: "backend-123",
				TargetsID:        []string{"target-1", "target-2"},
				TargetsType:      "instance",
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/backends/%s/targets", tt.request.LoadBalancerID, tt.request.NetworkBackendID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testBackendTargetClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

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
			err := client.Delete(context.Background(), DeleteNetworkBackendTargetRequest{
				LoadBalancerID:   tt.lbID,
				NetworkBackendID: tt.backendID,
				TargetID:         tt.targetID,
			})

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

	// Usar um contexto cancelado para forçar erro no newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela imediatamente

	client := testBackendTargetClient("http://dummy-url")

	req := CreateNetworkBackendTargetRequest{
		LoadBalancerID:   "lb-123",
		NetworkBackendID: "backend-123",
		TargetsID:        []string{"target-1", "target-2"},
		TargetsType:      "instance",
	}

	_, err := client.Create(ctx, req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
