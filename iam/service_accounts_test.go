package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceAccountService_List(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list service accounts",
			response: `[
				{
					"uuid": "sa-uuid-1",
					"name": "Service Account 1",
					"description": "First SA",
					"email": "sa1@example.com",
					"tenant": {"uuid": "tenant1", "legal_name": "Tenant 1"}
				},
				{
					"uuid": "sa-uuid-2",
					"name": "Service Account 2",
					"email": "sa2@example.com",
					"tenant": {"uuid": "tenant1", "legal_name": "Tenant 1"}
				}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty response",
			response:   `[]`,
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.ServiceAccounts().List(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("List() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestServiceAccountService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    ServiceAccountCreate
		response   string
		statusCode int
		wantUUID   string
		wantErr    bool
	}{
		{
			name: "successful create service account",
			request: ServiceAccountCreate{
				Name:        "New Service Account",
				Description: "Description",
				Email:       "new-sa@example.com",
			},
			response: `{
				"uuid": "new-sa-uuid",
				"name": "New Service Account",
				"description": "Description",
				"email": "new-sa@example.com",
				"tenant": {"uuid": "tenant1", "legal_name": "Tenant 1"}
			}`,
			statusCode: http.StatusCreated,
			wantUUID:   "new-sa-uuid",
			wantErr:    false,
		},
		{
			name: "server error",
			request: ServiceAccountCreate{
				Name:  "Invalid SA",
				Email: "invalid@example.com",
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
			result, err := client.ServiceAccounts().Create(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.UUID != tt.wantUUID {
				t.Errorf("Create() got UUID = %s, want %s", result.UUID, tt.wantUUID)
			}
		})
	}
}

func TestServiceAccountService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		saUUID     string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			saUUID:     "sa-uuid-to-delete",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "empty uuid",
			saUUID:     "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			saUUID:     "non-existent",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.ServiceAccounts().Delete(context.Background(), tt.saUUID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceAccountService_Edit(t *testing.T) {
	tests := []struct {
		name       string
		saUUID     string
		request    ServiceAccountEdit
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:   "successful edit",
			saUUID: "sa-uuid",
			request: ServiceAccountEdit{
				Name:        strPtr("Updated Name"),
				Description: strPtr("Updated Description"),
			},
			response: `{
				"uuid": "sa-uuid",
				"name": "Updated Name",
				"description": "Updated Description",
				"email": "sa@example.com",
				"tenant": {"uuid": "tenant1", "legal_name": "Tenant 1"}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty uuid",
			saUUID:     "",
			request:    ServiceAccountEdit{Name: strPtr("New Name")},
			statusCode: http.StatusBadRequest,
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
			result, err := client.ServiceAccounts().Edit(context.Background(), tt.saUUID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("Edit() returned nil result")
			}
		})
	}
}

func TestServiceAccountService_GetAPIKeys(t *testing.T) {
	tests := []struct {
		name       string
		saUUID     string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name:   "successful get API keys",
			saUUID: "sa-uuid",
			response: `[
				{
					"uuid": "key-uuid-1",
					"name": "API Key 1",
					"key_pair_id": "key-id-1",
					"scopes": ["read:instances"]
				},
				{
					"uuid": "key-uuid-2",
					"name": "API Key 2",
					"key_pair_id": "key-id-2",
					"scopes": ["write:instances"]
				}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty uuid",
			saUUID:     "",
			statusCode: http.StatusBadRequest,
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
			result, err := client.ServiceAccounts().APIKeys(context.Background(), tt.saUUID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("GetAPIKeys() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestServiceAccountService_CreateAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		saUUID     string
		request    APIKeyServiceAccountCreate
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:   "successful create API key",
			saUUID: "sa-uuid",
			request: APIKeyServiceAccountCreate{
				Name:        "New API Key",
				Description: strPtr("Description"),
				Scopes:      []string{"read:instances"},
			},
			response: `{
				"uuid": "new-key-uuid",
				"name": "New API Key",
				"description": "Description",
				"key_pair_id": "new-key-id",
				"key_pair_secret": "new-key-secret",
				"scopes": ["read:instances"]
			}`,
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:       "empty uuid",
			saUUID:     "",
			request:    APIKeyServiceAccountCreate{Name: "Key"},
			statusCode: http.StatusBadRequest,
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
			result, err := client.ServiceAccounts().CreateAPIKey(context.Background(), tt.saUUID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("CreateAPIKey() returned nil result")
			}
		})
	}
}

func TestServiceAccountService_RevokeAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		saUUID     string
		apikeyUUID string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful revoke",
			saUUID:     "sa-uuid",
			apikeyUUID: "key-uuid",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "empty sa uuid",
			saUUID:     "",
			apikeyUUID: "key-uuid",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "empty apikey uuid",
			saUUID:     "sa-uuid",
			apikeyUUID: "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := testClient(server.URL)
			err := client.ServiceAccounts().RevokeAPIKey(context.Background(), tt.saUUID, tt.apikeyUUID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RevokeAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceAccountService_EditAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		saUUID     string
		apikeyUUID string
		request    APIKeyServiceAccountEditInput
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful edit API key",
			saUUID:     "sa-uuid",
			apikeyUUID: "key-uuid",
			request: APIKeyServiceAccountEditInput{
				Name:        strPtr("Updated Key Name"),
				Description: strPtr("Updated Description"),
				Scopes:      []string{"read:instances", "write:instances"},
			},
			response: `{
				"uuid": "key-uuid",
				"name": "Updated Key Name",
				"description": "Updated Description",
				"scopes": ["read:instances", "write:instances"]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty sa uuid",
			saUUID:     "",
			apikeyUUID: "key-uuid",
			request:    APIKeyServiceAccountEditInput{Name: strPtr("New Name")},
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "empty apikey uuid",
			saUUID:     "sa-uuid",
			apikeyUUID: "",
			request:    APIKeyServiceAccountEditInput{Name: strPtr("New Name")},
			statusCode: http.StatusBadRequest,
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
			result, err := client.ServiceAccounts().EditAPIKey(context.Background(), tt.saUUID, tt.apikeyUUID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("EditAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("EditAPIKey() returned nil result")
			}
		})
	}
}
