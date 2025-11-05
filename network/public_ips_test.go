package network

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

func TestPublicIPService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with multiple IPs",
			response: `{
				"public_ips": [
					{"id": "ip1", "public_ip": "203.0.113.1"},
					{"id": "ip2", "public_ip": "203.0.113.2"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			response:   `{"public_ips": []}`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v0/public_ips", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testPublicIPClient(server.URL)
			ips, err := client.List(context.Background())

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(ips))
		})
	}
}

func TestPublicIPService_Get(t *testing.T) {
	t.Parallel()
	b, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	parsedTime := utils.LocalDateTimeWithoutZone(b)
	tests := []struct {
		name       string
		ipID       string
		response   string
		statusCode int
		want       *PublicIPResponse
		wantErr    bool
	}{
		{
			name: "existing public IP",
			ipID: "ip1",
			response: `{
				"id": "ip1",
				"public_ip": "203.0.113.5",
				"vpc_id": "vpc1",
				"status": "ACTIVE",
				"created_at": "2024-01-01T00:00:00.000000",
				"updated": "2024-01-02T00:00:00.000000"
			}`,
			statusCode: http.StatusOK,
			want: &PublicIPResponse{
				ID:        helpers.StrPtr("ip1"),
				PublicIP:  helpers.StrPtr("203.0.113.5"),
				VPCID:     helpers.StrPtr("vpc1"),
				Status:    helpers.StrPtr("ACTIVE"),
				Updated:   &parsedTime,
				CreatedAt: &parsedTime,
			},
			wantErr: false,
		},
		{
			name:       "non-existent IP",
			ipID:       "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			ipID:       "ip1",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/public_ips/%s", tt.ipID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testPublicIPClient(server.URL)
			ip, err := client.Get(context.Background(), tt.ipID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, *tt.want.ID, *ip.ID)
			assertEqual(t, *tt.want.PublicIP, *ip.PublicIP)
			assertEqual(t, *tt.want.Status, *ip.Status)
		})
	}
}

func TestPublicIPService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		ipID       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			ipID:       "ip1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent IP",
			ipID:       "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "attached IP",
			ipID:       "ip-attached",
			statusCode: http.StatusConflict,
			response:   `{"error": "ip is attached"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			ipID:       "ip1",
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/public_ips/%s", tt.ipID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testPublicIPClient(server.URL)
			err := client.Delete(context.Background(), tt.ipID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestPublicIPService_AttachDetach(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		method     string
		publicIPID string
		portID     string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful attach",
			method:     "Attach",
			publicIPID: "ip1",
			portID:     "port1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "attach non-existent IP",
			method:     "Attach",
			publicIPID: "invalid",
			portID:     "port1",
			statusCode: http.StatusNotFound,
			response:   `{"error": "public IP not found"}`,
			wantErr:    true,
		},
		{
			name:       "attach to invalid port",
			method:     "Attach",
			publicIPID: "ip1",
			portID:     "invalid",
			statusCode: http.StatusBadRequest,
			response:   `{"error": "invalid port"}`,
			wantErr:    true,
		},
		{
			name:       "successful detach",
			method:     "Detach",
			publicIPID: "ip1",
			portID:     "port1",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "detach non-attached IP",
			method:     "Detach",
			publicIPID: "ip1",
			portID:     "port2",
			statusCode: http.StatusConflict,
			response:   `{"error": "not attached"}`,
			wantErr:    true,
		},
		{
			name:       "server error on attach",
			method:     "Attach",
			publicIPID: "ip1",
			portID:     "port1",
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var expectedPath string
				if tt.method == "Attach" {
					expectedPath = fmt.Sprintf("/network/v0/public_ips/%s/attach/%s", tt.publicIPID, tt.portID)
				} else {
					expectedPath = fmt.Sprintf("/network/v0/public_ips/%s/detach/%s", tt.publicIPID, tt.portID)
				}

				assertEqual(t, expectedPath, r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testPublicIPClient(server.URL)
			var err error

			if tt.method == "Attach" {
				err = client.AttachToPort(context.Background(), tt.publicIPID, tt.portID)
			} else {
				err = client.DetachFromPort(context.Background(), tt.publicIPID, tt.portID)
			}

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func testPublicIPClient(baseURL string) PublicIPService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).PublicIPs()
}
