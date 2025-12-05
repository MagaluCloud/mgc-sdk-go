package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoleService_List(t *testing.T) {
	tests := []struct {
		name       string
		roleName   *string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list roles",
			response: `[
				{"name": "admin", "description": "Admin role", "origin": "system"},
				{"name": "viewer", "description": "Viewer role", "origin": "system"}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:     "successful list with role name filter",
			roleName: strPtr("admin"),
			response: `[
				{"name": "admin", "description": "Admin role", "origin": "system"}
			]`,
			statusCode: http.StatusOK,
			want:       1,
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
				if tt.roleName != nil {
					if r.URL.Query().Get("role_name") != *tt.roleName {
						t.Errorf("Expected role_name query param %s, got %s", *tt.roleName, r.URL.Query().Get("role_name"))
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Roles().List(context.Background(), tt.roleName)

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

func TestRoleService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateRole
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful create role",
			request: CreateRole{
				Name:        "custom-role",
				Description: strPtr("Custom role description"),
				Permissions: []string{"read:instances"},
			},
			response: `[
				{"name": "custom-role", "description": "Custom role description", "origin": "user"}
			]`,
			statusCode: http.StatusCreated,
			want:       1,
			wantErr:    false,
		},
		{
			name: "create role with based role",
			request: CreateRole{
				Name:      "derived-role",
				BasedRole: strPtr("admin"),
			},
			response: `[
				{"name": "derived-role", "origin": "user"}
			]`,
			statusCode: http.StatusCreated,
			want:       1,
			wantErr:    false,
		},
		{
			name: "server error",
			request: CreateRole{
				Name: "invalid-role",
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
			result, err := client.Roles().Create(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("Create() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestRoleService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		roleName   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			roleName:   "role-to-delete",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "empty role name",
			roleName:   "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			roleName:   "non-existent",
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
			err := client.Roles().Delete(context.Background(), tt.roleName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoleService_GetPermissions(t *testing.T) {
	tests := []struct {
		name       string
		roleName   string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:     "successful get permissions",
			roleName: "admin",
			response: `{
				"name": "admin",
				"description": "Admin role",
				"origin": "system",
				"permissions": ["read:instances", "write:instances", "delete:instances"]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty role name",
			roleName:   "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			roleName:   "non-existent",
			statusCode: http.StatusNotFound,
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
			result, err := client.Roles().Permissions(context.Background(), tt.roleName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != nil {
				if result.Name != tt.roleName {
					t.Errorf("GetPermissions() got name = %s, want %s", result.Name, tt.roleName)
				}
			}
		})
	}
}

func TestRoleService_EditPermissions(t *testing.T) {
	tests := []struct {
		name       string
		roleName   string
		request    EditPermissions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name:     "successful add permissions",
			roleName: "admin",
			request: EditPermissions{
				Add: []string{"read:networks"},
			},
			response: `[
				{"name": "admin", "origin": "system"}
			]`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name:     "successful remove permissions",
			roleName: "admin",
			request: EditPermissions{
				Remove: []string{"delete:instances"},
			},
			response: `[
				{"name": "admin", "origin": "system"}
			]`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name:       "empty role name",
			roleName:   "",
			request:    EditPermissions{Add: []string{"read:instances"}},
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
			result, err := client.Roles().EditPermissions(context.Background(), tt.roleName, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("EditPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("EditPermissions() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestRoleService_GetMembers(t *testing.T) {
	tests := []struct {
		name       string
		roleName   string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name:     "successful get members",
			roleName: "admin",
			response: `[
				{"member_uuid": "uuid1", "roles": ["admin"]},
				{"member_uuid": "uuid2", "roles": ["admin", "viewer"]}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty role name",
			roleName:   "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			roleName:   "non-existent",
			statusCode: http.StatusNotFound,
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
			result, err := client.Roles().Members(context.Background(), tt.roleName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("GetMembers() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}
