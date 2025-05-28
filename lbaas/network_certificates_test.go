package lbaas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func testCertificateClient(baseURL string) NetworkCertificateService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkCertificates()
}

func TestNetworkCertificateService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		request    CreateNetworkCertificateRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateNetworkCertificateRequest{
				LoadBalancerID: "lb-123",
				Name:           "test-cert",
				Certificate:    "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
				PrivateKey:     "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			response:   `{"id": "cert-123"}`,
			statusCode: http.StatusOK,
			want:       "cert-123",
			wantErr:    false,
		},
		{
			name: "server error",
			request: CreateNetworkCertificateRequest{
				LoadBalancerID: "lb-123",
				Name:           "test-cert",
				Certificate:    "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
				PrivateKey:     "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/tls-certificates", tt.request.LoadBalancerID), r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testCertificateClient(server.URL)
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

func TestNetworkCertificateService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		certID     string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:   "existing certificate",
			lbID:   "lb-123",
			certID: "cert-123",
			response: `{
				"id": "cert-123",
				"name": "test-cert",
				"description": "Test certificate",
				"certificate": "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent certificate",
			lbID:       "lb-123",
			certID:     "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/tls-certificates/%s", tt.lbID, tt.certID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testCertificateClient(server.URL)
			cert, err := client.Get(context.Background(), GetNetworkCertificateRequest{
				LoadBalancerID:   tt.lbID,
				TLSCertificateID: tt.certID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, "cert-123", cert.ID)
			assertEqual(t, "test-cert", cert.Name)
		})
	}
}

func TestNetworkCertificateService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with multiple certificates",
			lbID: "lb-123",
			response: `[
				{"id": "cert-1", "name": "test1", "certificate": "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
				{"id": "cert-2", "name": "test2", "certificate": "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			lbID:       "lb-123",
			response:   `[]`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name:       "server error",
			lbID:       "lb-123",
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
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/tls-certificates", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testCertificateClient(server.URL)
			certs, err := client.List(context.Background(), ListNetworkCertificateRequest{
				LoadBalancerID: tt.lbID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(certs))
		})
	}
}

func TestNetworkCertificateService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		certID     string
		request    UpdateNetworkCertificateRequest
		statusCode int
		wantErr    bool
	}{
		{
			name:   "successful update",
			lbID:   "lb-123",
			certID: "cert-123",
			request: UpdateNetworkCertificateRequest{
				LoadBalancerID:   "lb-123",
				TLSCertificateID: "cert-123",
				Certificate:      "updated-cert",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "non-existent certificate",
			lbID:   "lb-123",
			certID: "invalid",
			request: UpdateNetworkCertificateRequest{
				LoadBalancerID:   "lb-123",
				TLSCertificateID: "invalid",
				Certificate:      "updated-cert",
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/tls-certificates/%s", tt.lbID, tt.certID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testCertificateClient(server.URL)
			err := client.Update(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkCertificateService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		certID     string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			lbID:       "lb-123",
			certID:     "cert-123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent certificate",
			lbID:       "lb-123",
			certID:     "invalid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/tls-certificates/%s", tt.lbID, tt.certID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testCertificateClient(server.URL)
			err := client.Delete(context.Background(), DeleteNetworkCertificateRequest{
				LoadBalancerID:   tt.lbID,
				TLSCertificateID: tt.certID,
			})

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}
