package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccessControlService_Get(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful get access control",
			response: `{
				"name": "default",
				"description": "Default access control",
				"tenant_id": "tenant1",
				"enabled": true,
				"enforce_mfa": false
			}`,
			statusCode: http.StatusOK,
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.AccessControl().Get(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("Get() returned nil result")
			}
		})
	}
}

func TestAccessControlService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    AccessControlCreate
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful create access control",
			request: AccessControlCreate{
				Name:        strPtr("custom-ac"),
				Description: strPtr("Custom access control"),
			},
			response: `{
				"name": "custom-ac",
				"description": "Custom access control",
				"tenant_id": "tenant1",
				"enabled": true,
				"enforce_mfa": false
			}`,
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "create with minimal fields",
			request: AccessControlCreate{
				Name: strPtr("minimal-ac"),
			},
			response: `{
				"name": "minimal-ac",
				"tenant_id": "tenant1",
				"enabled": true,
				"enforce_mfa": false
			}`,
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "server error",
			request: AccessControlCreate{
				Name: strPtr("invalid-ac"),
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.AccessControl().Create(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("Create() returned nil result")
			}
		})
	}
}

func TestAccessControlService_Update(t *testing.T) {
	tests := []struct {
		name       string
		request    AccessControlStatus
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful update status",
			request: AccessControlStatus{
				Status:     boolPtr(true),
				EnforceMFA: boolPtr(false),
			},
			response: `{
				"name": "default",
				"enabled": true,
				"enforce_mfa": false
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "update only status",
			request: AccessControlStatus{
				Status: boolPtr(false),
			},
			response: `{
				"name": "default",
				"enabled": false,
				"enforce_mfa": false
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "update only enforce_mfa",
			request: AccessControlStatus{
				EnforceMFA: boolPtr(true),
			},
			response: `{
				"name": "default",
				"enabled": true,
				"enforce_mfa": true
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "server error",
			request: AccessControlStatus{
				Status: boolPtr(true),
			},
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.AccessControl().Update(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("Update() returned nil result")
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
