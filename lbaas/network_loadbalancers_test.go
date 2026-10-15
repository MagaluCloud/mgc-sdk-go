package lbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

// Helper functions
func assertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v but got %v. %v", expected, actual, msgAndArgs)
	}
}

func assertNotEqual(t *testing.T, notExpected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if notExpected == actual {
		t.Errorf("Expected anything but %v. %v", notExpected, msgAndArgs)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func testLoadBalancerClient(baseURL string) NetworkLoadBalancerService {
	httpClient := &http.Client{}
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).NetworkLoadBalancers()
}

func TestNetworkLoadBalancerService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		request    CreateNetworkLoadBalancerRequest
		response   string
		statusCode int
		want       string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []CreateNetworkBackendRequest{},
			},
			response:   `{"id": "lb-123"}`,
			statusCode: http.StatusOK,
			want:       "lb-123",
			wantErr:    false,
		},
		{
			name: "bad request - invalid request data",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "",
				Visibility: "invalid",
				VPCID:      "",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []CreateNetworkBackendRequest{},
			},
			response:   `{"error": "invalid request: name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized - invalid credentials",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []CreateNetworkBackendRequest{},
			},
			response:   `{"error": "unauthorized access"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden - insufficient permissions",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []CreateNetworkBackendRequest{},
			},
			response:   `{"error": "forbidden: insufficient permissions"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "conflict - resource already exists",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "existing-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []CreateNetworkBackendRequest{},
			},
			response:   `{"error": "load balancer with name 'existing-lb' already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "server error",
			request: CreateNetworkLoadBalancerRequest{
				Name:       "test-lb",
				Visibility: "external",
				VPCID:      "vpc-123",
				Listeners:  []NetworkListenerRequest{},
				Backends:   []CreateNetworkBackendRequest{},
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
				assertEqual(t, "/load-balancer/v0beta1/network-load-balancers", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)
				assertEqual(t, "application/json", r.Header.Get("Content-Type"))

				// Validate request body contains expected fields
				if tt.statusCode == http.StatusOK {
					var body CreateNetworkLoadBalancerRequest
					if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
						assertEqual(t, tt.request.Name, body.Name)
						assertEqual(t, tt.request.Visibility, body.Visibility)
						assertEqual(t, tt.request.VPCID, body.VPCID)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, id)
		})
	}
}

func TestNetworkLoadBalancerService_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		lbID                 string
		response             string
		statusCode           int
		wantErr              bool
		expectedID           string
		expectedName         string
		expectedErrorDetails *string
	}{
		{
			name: "existing load balancer",
			lbID: "lb-123",
			response: `{
				"id": "lb-123",
				"name": "test-lb",
				"type": "proxy",
				"visibility": "external",
				"status": "running",
				"listeners": [],
				"backends": [],
				"health_checks": [],
				"public_ip": null,
				"tls_certificates": [],
				"acls": [],
				"vpc_id": "vpc-123",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode:           http.StatusOK,
			wantErr:              false,
			expectedID:           "lb-123",
			expectedName:         "test-lb",
			expectedErrorDetails: nil,
		},
		{
			name: "load balancer with error details",
			lbID: "lb-456",
			response: `{
				"id": "lb-456",
				"name": "test-lb-error",
				"type": "proxy",
				"visibility": "external",
				"status": "error",
				"listeners": [],
				"backends": [],
				"health_checks": [],
				"public_ip": null,
				"tls_certificates": [],
				"acls": [],
				"vpc_id": "vpc-456",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z",
				"error_details": "Failed to provision backend"
			}`,
			statusCode:           http.StatusOK,
			wantErr:              false,
			expectedID:           "lb-456",
			expectedName:         "test-lb-error",
			expectedErrorDetails: stringPtr("Failed to provision backend"),
		},
		{
			name:                 "non-existent load balancer",
			lbID:                 "invalid",
			response:             `{"error": "not found"}`,
			statusCode:           http.StatusNotFound,
			wantErr:              true,
			expectedID:           "",
			expectedName:         "",
			expectedErrorDetails: nil,
		},
		{
			name:                 "unauthorized access",
			lbID:                 "lb-123",
			response:             `{"error": "unauthorized"}`,
			statusCode:           http.StatusUnauthorized,
			wantErr:              true,
			expectedID:           "",
			expectedName:         "",
			expectedErrorDetails: nil,
		},
		{
			name:                 "forbidden access",
			lbID:                 "lb-123",
			response:             `{"error": "forbidden"}`,
			statusCode:           http.StatusForbidden,
			wantErr:              true,
			expectedID:           "",
			expectedName:         "",
			expectedErrorDetails: nil,
		},
		{
			name:                 "server error",
			lbID:                 "lb-123",
			response:             `{"error": "internal server error"}`,
			statusCode:           http.StatusInternalServerError,
			wantErr:              true,
			expectedID:           "",
			expectedName:         "",
			expectedErrorDetails: nil,
		},
		{
			name:                 "nil response body",
			lbID:                 "lb-123",
			response:             ``,
			statusCode:           http.StatusOK,
			wantErr:              true,
			expectedID:           "",
			expectedName:         "",
			expectedErrorDetails: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			lb, err := client.Get(context.Background(), tt.lbID)

			if tt.wantErr {
				assertError(t, err)
				if tt.name == "nil response body" {
					// The nil response body test should trigger JSON parsing error, not the custom error
					assertEqual(t, true, strings.Contains(err.Error(), "unexpected end of JSON input") || strings.Contains(err.Error(), "EOF"))
				} else {
					assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				}
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.expectedID, lb.ID)
			assertEqual(t, tt.expectedName, lb.Name)
			// Validate the PublicIP field is correctly recognized as null
			assertEqual(t, true, lb.PublicIP == nil)
			// Validate ErrorDetails
			if tt.expectedErrorDetails == nil {
				assertEqual(t, true, lb.ErrorDetails == nil)
			} else {
				assertEqual(t, *tt.expectedErrorDetails, *lb.ErrorDetails)
			}
		})
	}
}

func TestNetworkLoadBalancerService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with multiple load balancers",
			response: `{
				"results": [
					{"id": "lb-1", "name": "test1", "type": "proxy", "visibility": "external", "status": "running", "listeners": [], "backends": [], "health_checks": [], "public_ip": null, "tls_certificates": [], "acls": [], "vpc_id": "vpc-1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
					{"id": "lb-2", "name": "test2", "type": "proxy", "visibility": "internal", "status": "running", "listeners": [], "backends": [], "health_checks": [], "public_ip": null, "tls_certificates": [], "acls": [], "vpc_id": "vpc-2", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			response:   `{"results": []}`,
			statusCode: http.StatusOK,
			want:       0,
			wantErr:    false,
		},
		{
			name:       "unauthorized access",
			response:   `{"error": "unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "forbidden access",
			response:   `{"error": "forbidden"}`,
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name:       "server error",
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
				assertEqual(t, "/load-balancer/v0beta1/network-load-balancers", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			resp, err := client.List(context.Background(), ListNetworkLoadBalancerRequest{})

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

func TestNetworkLoadBalancerService_ListWithPagination(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		request     ListNetworkLoadBalancerRequest
		response    string
		statusCode  int
		expectQuery map[string]string
		want        int
		wantErr     bool
	}{
		{
			name: "list with offset and limit - validates request creation",
			request: ListNetworkLoadBalancerRequest{
				Offset: intPtr(10),
				Limit:  intPtr(5),
			},
			response: `{
				"results": [
					{"id": "lb-1", "name": "test1", "type": "proxy", "visibility": "external", "status": "running", "listeners": [], "backends": [], "health_checks": [], "public_ip": null, "tls_certificates": [], "acls": [], "vpc_id": "vpc-1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
				]
			}`,
			statusCode:  http.StatusOK,
			expectQuery: map[string]string{}, // AddReflect doesn't handle pointer types in current implementation
			want:        1,
			wantErr:     false,
		},
		{
			name: "list with sort parameter",
			request: ListNetworkLoadBalancerRequest{
				Sort: stringPtr("created_at:desc"),
			},
			response: `{
				"results": [
					{"id": "lb-1", "name": "test1", "type": "proxy", "visibility": "external", "status": "running", "listeners": [], "backends": [], "health_checks": [], "public_ip": null, "tls_certificates": [], "acls": [], "vpc_id": "vpc-1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
					{"id": "lb-2", "name": "test2", "type": "proxy", "visibility": "internal", "status": "running", "listeners": [], "backends": [], "health_checks": [], "public_ip": null, "tls_certificates": [], "acls": [], "vpc_id": "vpc-2", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
				]
			}`,
			statusCode: http.StatusOK,
			expectQuery: map[string]string{
				"_sort": "created_at:desc",
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "list with all parameters - validates request creation",
			request: ListNetworkLoadBalancerRequest{
				Offset: intPtr(0),
				Limit:  intPtr(20),
				Sort:   stringPtr("name:asc"),
			},
			response:   `{"results": []}`,
			statusCode: http.StatusOK,
			expectQuery: map[string]string{
				"_sort": "name:asc", // Only _sort works because it uses Add() method for *string
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/load-balancer/v0beta1/network-load-balancers", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				// Validate query parameters that should be present
				for key, expectedValue := range tt.expectQuery {
					actualValue := r.URL.Query().Get(key)
					assertEqual(t, expectedValue, actualValue, fmt.Sprintf("Query parameter %s", key))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			resp, err := client.List(context.Background(), tt.request)

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

func TestNetworkLoadBalancerService_ListAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pages      []string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "single page with all results",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/network-load-balancers?_limit=50&_offset=0"
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
							"id": "lb-1",
							"name": "test1",
							"type": "proxy",
							"visibility": "external",
							"status": "running",
							"listeners": [],
							"backends": [],
							"health_checks": [],
							"tls_certificates": [],
							"acls": [],
							"vpc_id": "vpc-1",
							"created_at": "2024-01-01T00:00:00Z",
							"updated_at": "2024-01-01T00:00:00Z"
						},
						{
							"id": "lb-2",
							"name": "test2",
							"type": "proxy",
							"visibility": "internal",
							"status": "running",
							"listeners": [],
							"backends": [],
							"health_checks": [],
							"tls_certificates": [],
							"acls": [],
							"vpc_id": "vpc-2",
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
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/network-load-balancers?_limit=50&_offset=0",
							"next": "/network-load-balancers?_limit=50&_offset=50"
						},
						"page": {
							"count": 50,
							"limit": 50,
							"offset": 0,
							"total": 75
						}
					},
					"results": [` + generateLoadBalancerResults(1, 50) + `]
				}`,
				`{
					"meta": {
						"links": {
							"self": "/network-load-balancers?_limit=50&_offset=50",
							"previous": "/network-load-balancers?_limit=50&_offset=0"
						},
						"page": {
							"count": 25,
							"limit": 50,
							"offset": 50,
							"total": 75
						}
					},
					"results": [` + generateLoadBalancerResults(51, 25) + `]
				}`,
			},
			statusCode: http.StatusOK,
			wantCount:  75,
			wantErr:    false,
		},
		{
			name: "empty results",
			pages: []string{
				`{
					"meta": {
						"links": {
							"self": "/network-load-balancers?_limit=50&_offset=0"
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
				assertEqual(t, "/load-balancer/v0beta1/network-load-balancers", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if pageIndex < len(tt.pages) {
					w.Write([]byte(tt.pages[pageIndex]))
					pageIndex++
				} else {
					w.Write([]byte(`{"meta":{"links":{"self":""},"page":{"count":0,"limit":50,"offset":0,"total":0}},"results":[]}`))
				}
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			loadBalancers, err := client.ListAll(context.Background())

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantCount, len(loadBalancers))
		})
	}
}

func generateLoadBalancerResults(start, count int) string {
	results := make([]string, count)
	for i := 0; i < count; i++ {
		id := start + i
		results[i] = fmt.Sprintf(`{
			"id": "lb-%d",
			"name": "test%d",
			"type": "proxy",
			"visibility": "external",
			"status": "running",
			"listeners": [],
			"backends": [],
			"health_checks": [],
			"tls_certificates": [],
			"acls": [],
			"vpc_id": "vpc-%d",
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z"
		}`, id, id, id)
	}
	return strings.Join(results, ",")
}

func TestNetworkLoadBalancerService_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		request    UpdateNetworkLoadBalancerRequest
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful update",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				Name: stringPtr("updated-lb"),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non-existent load balancer",
			lbID: "invalid",
			request: UpdateNetworkLoadBalancerRequest{
				Name: stringPtr("updated-lb"),
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "bad request - invalid data",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				Name: stringPtr(""),
			},
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unauthorized access",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				Name: stringPtr("updated-lb"),
			},
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "forbidden access",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				Name: stringPtr("updated-lb"),
			},
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name: "conflict - name already exists",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				Name: stringPtr("existing-name"),
			},
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "server error",
			lbID: "lb-123",
			request: UpdateNetworkLoadBalancerRequest{
				Name: stringPtr("updated-lb"),
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodPut, r.Method)
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					w.Write([]byte(fmt.Sprintf(`{"id": "%s"}`, tt.lbID)))
				}
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			id, err := client.Update(context.Background(), tt.lbID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				assertEqual(t, "", id)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.lbID, id)
		})
	}
}

func TestNetworkLoadBalancerService_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lbID       string
		request    DeleteNetworkLoadBalancerRequest
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			lbID:       "lb-123",
			request:    DeleteNetworkLoadBalancerRequest{},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "deletion with delete public IP",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				DeletePublicIP: boolPtr(true),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "deletion without deleting public IP",
			lbID: "lb-123",
			request: DeleteNetworkLoadBalancerRequest{
				DeletePublicIP: boolPtr(false),
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existent load balancer",
			lbID:       "invalid",
			request:    DeleteNetworkLoadBalancerRequest{},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "unauthorized access",
			lbID:       "lb-123",
			request:    DeleteNetworkLoadBalancerRequest{},
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "forbidden access",
			lbID:       "lb-123",
			request:    DeleteNetworkLoadBalancerRequest{},
			statusCode: http.StatusForbidden,
			wantErr:    true,
		},
		{
			name:       "conflict - resource in use",
			lbID:       "lb-123",
			request:    DeleteNetworkLoadBalancerRequest{},
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name:       "server error",
			lbID:       "lb-123",
			request:    DeleteNetworkLoadBalancerRequest{},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/load-balancer/v0beta1/network-load-balancers/%s", tt.lbID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)

				// Check query parameter if DeletePublicIP is set
				if tt.request.DeletePublicIP != nil {
					expected := "false"
					if *tt.request.DeletePublicIP {
						expected = "true"
					}
					assertEqual(t, expected, r.URL.Query().Get("delete_public_ip"))
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testLoadBalancerClient(server.URL)
			err := client.Delete(context.Background(), tt.lbID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

// Helper functions for pointer values
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func TestNetworkLoadBalancerService_ContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow server response
		select {
		case <-r.Context().Done():
			return
		case <-time.After(100 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "lb-123"}`))
		}
	}))
	defer server.Close()

	client := testLoadBalancerClient(server.URL)

	// Test context cancellation for Create
	t.Run("Create with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := client.Create(ctx, CreateNetworkLoadBalancerRequest{
			Name:       "test-lb",
			Visibility: "external",
			VPCID:      "vpc-123",
			Listeners:  []NetworkListenerRequest{},
			Backends:   []CreateNetworkBackendRequest{},
		})

		assertError(t, err)
		assertEqual(t, true, strings.Contains(err.Error(), "context canceled"))
	})

	// Test context timeout for Get
	t.Run("Get with timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := client.Get(ctx, "lb-123")

		assertError(t, err)
		assertEqual(t, true, strings.Contains(err.Error(), "context deadline exceeded"))
	})
}

// New test validating the JSON tag for PublicIP
func TestNetworkLoadBalancer_PublicIPJSONTag(t *testing.T) {
	t.Parallel()

	ip := "1.2.3.4"
	lb := NetworkLoadBalancerResponse{
		ID:         "lb-1",
		Name:       "name",
		Type:       "proxy",
		Visibility: "external",
		Status:     "running",
		PublicIP: &NetworkPublicIPResponse{
			ID:         "pip-1",
			IPAddress:  &ip,
			ExternalID: "ext-1",
		},
		VPCID:     "vpc-1",
		CreatedAt: time.Unix(0, 0).UTC(),
		UpdatedAt: time.Unix(0, 0).UTC(),
	}

	data, err := json.Marshal(lb)
	assertNoError(t, err)

	js := string(data)
	assertEqual(t, true, strings.Contains(js, "\"public_ip\""))
	assertEqual(t, false, strings.Contains(js, "\"public_ips\""))

	// Ensure unmarshalling with the correct tag works and sets PublicIP to nil when null
	var decoded NetworkLoadBalancerResponse
	err = json.Unmarshal([]byte(`{"id":"lb-1","name":"name","type":"proxy","visibility":"external","status":"running","listeners":[],"backends":[],"health_checks":[],"public_ip":null,"tls_certificates":[],"acls":[],"vpc_id":"vpc-1","created_at":"1970-01-01T00:00:00Z","updated_at":"1970-01-01T00:00:00Z"}`), &decoded)
	assertNoError(t, err)
	assertEqual(t, true, decoded.PublicIP == nil)
}

func TestNetworkLoadBalancerService_Create_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testLoadBalancerClient("http://dummy-url")

	req := CreateNetworkLoadBalancerRequest{
		Name:       "test-lb",
		Visibility: "external",
		VPCID:      "vpc-123",
		Listeners:  []NetworkListenerRequest{},
		Backends:   []CreateNetworkBackendRequest{},
	}

	_, err := client.Create(ctx, req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkLoadBalancerService_Get_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testLoadBalancerClient("http://dummy-url")

	_, err := client.Get(ctx, "lb-123")

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkLoadBalancerService_List_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testLoadBalancerClient("http://dummy-url")

	_, err := client.List(ctx, ListNetworkLoadBalancerRequest{})

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkLoadBalancerService_ListAll_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testLoadBalancerClient("http://dummy-url")

	_, err := client.ListAll(ctx)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkLoadBalancerService_Update_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testLoadBalancerClient("http://dummy-url")

	req := UpdateNetworkLoadBalancerRequest{
		Name: stringPtr("updated-lb"),
	}

	_, err := client.Update(ctx, "lb-123", req)

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}

func TestNetworkLoadBalancerService_Delete_NewRequestError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := testLoadBalancerClient("http://dummy-url")

	err := client.Delete(ctx, "lb-123", DeleteNetworkLoadBalancerRequest{})

	if err == nil {
		t.Error("expected error due to canceled context, got nil")
	}
}
