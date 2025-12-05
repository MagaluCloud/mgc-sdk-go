package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMemberService_List(t *testing.T) {
	tests := []struct {
		name       string
		email      *string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list members",
			response: `[
				{"uuid": "uuid1", "email": "user1@example.com", "name": "User 1", "tenant_id": "tenant1"},
				{"uuid": "uuid2", "email": "user2@example.com", "name": "User 2", "tenant_id": "tenant1"}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:  "successful list with email filter",
			email: strPtr("user1@example.com"),
			response: `[
				{"uuid": "uuid1", "email": "user1@example.com", "name": "User 1", "tenant_id": "tenant1"}
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
				if tt.email != nil {
					if r.URL.Query().Get("email") != *tt.email {
						t.Errorf("Expected email query param %s, got %s", *tt.email, r.URL.Query().Get("email"))
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Members().List(context.Background(), tt.email)

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

func TestMemberService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    CreateMember
		response   string
		statusCode int
		wantUUID   string
		wantErr    bool
	}{
		{
			name: "successful create member",
			request: CreateMember{
				Email: "newuser@example.com",
				Roles: []string{"admin"},
			},
			response:   `{"uuid": "new-uuid", "email": "newuser@example.com", "name": "New User", "tenant_id": "tenant1"}`,
			statusCode: http.StatusCreated,
			wantUUID:   "new-uuid",
			wantErr:    false,
		},
		{
			name: "create member with permissions",
			request: CreateMember{
				Email:       "newuser@example.com",
				Permissions: []string{"read:instances"},
			},
			response:   `{"uuid": "new-uuid", "email": "newuser@example.com", "name": "New User", "tenant_id": "tenant1"}`,
			statusCode: http.StatusCreated,
			wantUUID:   "new-uuid",
			wantErr:    false,
		},
		{
			name: "server error",
			request: CreateMember{
				Email: "newuser@example.com",
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
			result, err := client.Members().Create(context.Background(), tt.request)

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

func TestMemberService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		uuid       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			uuid:       "uuid-to-delete",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "empty uuid",
			uuid:       "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			uuid:       "non-existent",
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
			err := client.Members().Delete(context.Background(), tt.uuid)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemberGrantsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		uuid       string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful get grants",
			uuid: "member-uuid",
			response: `{
				"roles": ["admin", "viewer"],
				"permissions": ["read:instances", "write:instances"]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty uuid",
			uuid:       "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "not found",
			uuid:       "non-existent",
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
			result, err := client.Members().Grants().Get(context.Background(), tt.uuid)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetGrants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != nil {
				if len(result.Roles) == 0 && len(result.Permissions) == 0 {
					t.Error("GetGrants() returned empty grants")
				}
			}
		})
	}
}

func TestMemberGrantsService_Add(t *testing.T) {
	tests := []struct {
		name       string
		uuid       string
		request    EditGrant
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful add grants",
			uuid: "member-uuid",
			request: EditGrant{
				Operation: OperationAdd,
				Roles:     []string{"admin"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "successful remove grants",
			uuid: "member-uuid",
			request: EditGrant{
				Operation: OperationRemove,
				Roles:     []string{"viewer"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty uuid",
			uuid:       "",
			request:    EditGrant{Operation: OperationAdd},
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
			err := client.Members().Grants().Add(context.Background(), tt.uuid, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddGrants() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemberGrantsService_BatchUpdate(t *testing.T) {
	tests := []struct {
		name       string
		request    BatchUpdateMembers
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful batch update",
			request: BatchUpdateMembers{
				MemberIDs: []string{"uuid1", "uuid2"},
				Operation: OperationAdd,
				RoleNames: []string{"admin"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "server error",
			request: BatchUpdateMembers{
				MemberIDs: []string{"uuid1"},
				Operation: OperationAdd,
			},
			statusCode: http.StatusInternalServerError,
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
			err := client.Members().Grants().BatchUpdate(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("BatchUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
