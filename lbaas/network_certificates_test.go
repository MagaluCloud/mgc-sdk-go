package lbaas

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func testCertificateClient(baseURL string) NetworkCertificateService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkCertificates()
}

func TestNetworkCertificateService_Create(t *testing.T) {
	t.Parallel()

	// Certificado e chave privada em formato PEM
	certPEM := "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----"
	keyPEM := "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----"

	// Codificar em base64 conforme esperado pela API
	certBase64 := base64.StdEncoding.EncodeToString([]byte(certPEM))
	keyBase64 := base64.StdEncoding.EncodeToString([]byte(keyPEM))

	tests := []struct {
		name       string
		request    CreateNetworkCertificateRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
		errorMsg   string
	}{
		{
			name: "successful creation",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: certBase64,
				PrivateKey:  keyBase64,
			},
			response:   `{"id": "cert-123"}`,
			statusCode: http.StatusOK,
			want:       "cert-123",
			wantErr:    false,
		},
		{
			name: "invalid certificate - not base64 encoded",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: "invalid-base64@#$%",
				PrivateKey:  keyBase64,
			},
			wantErr:  true,
			errorMsg: "certificate is not base64 encoded",
		},
		{
			name: "invalid private key - not base64 encoded",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: certBase64,
				PrivateKey:  "invalid-base64@#$%",
			},
			wantErr:  true,
			errorMsg: "private key is not base64 encoded",
		},
		{
			name: "server error",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: certBase64,
				PrivateKey:  keyBase64,
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "bad request - invalid certificate",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: certBase64,
				PrivateKey:  keyBase64,
			},
			response:   `{"error": "invalid certificate format"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized access",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: certBase64,
				PrivateKey:  keyBase64,
			},
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden access",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: certBase64,
				PrivateKey:  keyBase64,
			},
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "load balancer not found",
			request: CreateNetworkCertificateRequest{
				Name:        "test-cert",
				Certificate: certBase64,
				PrivateKey:  keyBase64,
			},
			response:   `{"error": "load balancer not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "conflict - certificate already exists",
			request: CreateNetworkCertificateRequest{
				Name:        "existing-cert",
				Certificate: certBase64,
				PrivateKey:  keyBase64,
			},
			response:   `{"error": "certificate with this name already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Casos que testam validação de base64 não precisam de servidor HTTP
			if tt.errorMsg != "" {
				client := testCertificateClient("http://dummy-url")
				_, err := client.Create(context.Background(), "lb-123", tt.request)

				assertError(t, err)
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/load-balancer/v0beta1/network-load-balancers/lb-123/tls-certificates", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testCertificateClient(server.URL)
			cert, err := client.Create(context.Background(), "lb-123", tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, cert.ID)
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
			cert, err := client.Get(context.Background(), tt.lbID, tt.certID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
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
			response: `{
				"meta": {
					"current_page": 1,
					"total_count": 2,
					"total_pages": 1,
					"total_results": 2
				},
				"results": [
					{"id": "cert-1", "name": "test1", "certificate": "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
					{"id": "cert-2", "name": "test2", "certificate": "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "empty list",
			lbID: "lb-123",
			response: `{
				"meta": {
					"current_page": 1,
					"total_count": 0,
					"total_pages": 0,
					"total_results": 0
				},
				"results": []
			}`,
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
			resp, err := client.List(context.Background(), tt.lbID, ListNetworkLoadBalancerRequest{})

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(resp.Results))
		})
	}
}

func TestNetworkCertificateService_ListAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		lbID       string
		pages      []string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page with all results",
			lbID: "lb-123",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/tls-certificates?_limit=50&_offset=0"
						},
						"page": {
							"count": 2,
							"limit": 50,
							"offset": 0,
							"total": 2
						}
					},
					"results": [
						{
							"id": "cert-1",
							"name": "test1",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						},
						{
							"id": "cert-2",
							"name": "test2",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						}
					]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name: "multiple pages",
			lbID: "lb-456",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/tls-certificates?_limit=50&_offset=0",
							"next": "/tls-certificates?_limit=50&_offset=50"
						},
						"page": {
							"count": 50,
							"limit": 50,
							"offset": 0,
							"total": 75
						}
					},
					"results": [` + generateCertificateResults(1, 50) + `]
				}`,
				`{
					"meta": {
						"links": {
							"self": "/tls-certificates?_limit=50&_offset=50",
							"previous": "/tls-certificates?_limit=50&_offset=0"
						},
						"page": {
							"count": 25,
							"limit": 50,
							"offset": 50,
							"total": 75
						}
					},
					"results": [` + generateCertificateResults(51, 25) + `]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  75,
			wantErr:    false,
		},
		{
			name: "empty results",
			lbID: "lb-789",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/tls-certificates?_limit=50&_offset=0"
						},
						"page": {
							"count": 0,
							"limit": 50,
							"offset": 0,
							"total": 0
						}
					},
					"results": []
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pageIndex := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s/tls-certificates", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if pageIndex < len(tt.pages) {
					w.Write([]byte(tt.pages[pageIndex]))
					pageIndex++
				} else {
					// Return empty results if we've exhausted pages
					w.Write([]byte(`{"meta":{"links":{"self":""},"page":{"count":0,"limit":50,"offset":0,"total":0}},"results":[]}`))
				}
			}))
			defer server.Close()

			client := testCertificateClient(server.URL)
			certificates, err := client.ListAll(context.Background(), tt.lbID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(certificates))
		})
	}
}

func generateCertificateResults(start, count int) string {
	results := make([]string, count)
	for i := 0; i < count; i++ {
		id := start + i
		results[i] = fmt.Sprintf(`{
			"id": "cert-%d",
			"name": "test%d",
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z"
		}`, id, id)
	}
	return strings.Join(results, ",")
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
				Certificate: "updated-cert",
				PrivateKey:  "updated-key",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "non-existent certificate",
			lbID:   "lb-123",
			certID: "invalid",
			request: UpdateNetworkCertificateRequest{
				Certificate: "updated-cert",
				PrivateKey:  "updated-key",
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
			err := client.Update(context.Background(), tt.lbID, tt.certID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
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
			err := client.Delete(context.Background(), tt.lbID, tt.certID)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestNetworkCertificateService_Create_NewRequestError(t *testing.T) {
	t.Parallel()

	// Criar um cliente com URL base inválida para forçar erro no newRequest
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl("://invalid-url")), // URL malformada
		client.WithHTTPClient(httpClient))

	certificateService := New(core).NetworkCertificates()

	certPEM := "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----"
	keyPEM := "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----"
	certBase64 := base64.StdEncoding.EncodeToString([]byte(certPEM))
	keyBase64 := base64.StdEncoding.EncodeToString([]byte(keyPEM))

	req := CreateNetworkCertificateRequest{
		Name:        "test-cert",
		Certificate: certBase64,
		PrivateKey:  keyBase64,
	}

	_, err := certificateService.Create(context.Background(), "lb-123", req)

	if err == nil {
		t.Error("expected error due to invalid URL, got nil")
	}
}

func TestNetworkCertificateService_Get_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testCertificateClient("http://dummy-url")

	_, err := client.Get(ctx, "lb-123", "cert-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkCertificateService_List_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testCertificateClient("http://dummy-url")

	_, err := client.List(ctx, "lb-123", ListNetworkLoadBalancerRequest{})

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkCertificateService_ListAll_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testCertificateClient("http://dummy-url")

	_, err := client.ListAll(ctx, "lb-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkCertificateService_Update_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testCertificateClient("http://dummy-url")

	req := UpdateNetworkCertificateRequest{
		Certificate: "updated-cert",
		PrivateKey:  "updated-key",
	}

	err := client.Update(ctx, "lb-123", "cert-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkCertificateService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	// Use a canceled context to force an error in newRequest
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := testCertificateClient("http://dummy-url")

	err := client.Delete(ctx, "lb-123", "cert-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
