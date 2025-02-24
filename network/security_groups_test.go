package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/internal/utils"
)

func TestSecurityGroupService_List(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful list with security groups",
			response: `{
				"security_groups": [
					{"id": "sg1", "name": "test-sg1", "status": "ACTIVE"},
					{"id": "sg2", "name": "test-sg2", "status": "ACTIVE"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "empty list",
			response:   `{"security_groups": []}`,
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v0/security_groups", r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSecurityGroupClient(server.URL)
			groups, err := client.List(context.Background())

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(groups))
		})
	}
}

func TestSecurityGroupService_Get(t *testing.T) {
	basetime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	parsedTime := utils.LocalDateTimeWithoutZone(basetime)

	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		want       *SecurityGroupDetailResponse
		wantErr    bool
	}{
		{
			name: "existing security group",
			id:   "sg1",
			response: `{
				"id": "sg1",
				"name": "test-sg",
				"status": "ACTIVE",
				"external_id": "ext123",
				"rules": [
					{"id": "rule1", "direction": "ingress"}
				],
				"created_at": "2024-01-01T00:00:00",
				"updated": "2024-01-01T00:00:00"
			}`,
			statusCode: http.StatusOK,
			want: &SecurityGroupDetailResponse{
				SecurityGroupResponse: SecurityGroupResponse{
					ID:        helpers.StrPtr("sg1"),
					Name:      helpers.StrPtr("test-sg"),
					Status:    "ACTIVE",
					CreatedAt: &parsedTime,
					Updated:   &parsedTime,
				},
				ExternalID: helpers.StrPtr("ext123"),
				Rules: &[]RuleResponse{
					{ID: helpers.StrPtr("rule1"), Direction: helpers.StrPtr("ingress")},
				},
			},
			wantErr: false,
		},
		{
			name:       "non-existent security group",
			id:         "invalid",
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
				assertEqual(t, fmt.Sprintf("/network/v0/security_groups/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSecurityGroupClient(server.URL)
			group, err := client.Get(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, *tt.want.ID, *group.ID)
			assertEqual(t, *tt.want.Name, *group.Name)
			assertEqual(t, *tt.want.ExternalID, *group.ExternalID)
			assertEqual(t, len(*tt.want.Rules), len(*group.Rules))
		})
	}
}

func TestSecurityGroupService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    SecurityGroupCreateRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful create",
			request: SecurityGroupCreateRequest{
				Name:        "test-sg",
				Description: helpers.StrPtr("test description"),
			},
			response:   `{"id": "sg-new"}`,
			statusCode: http.StatusOK,
			wantID:     "sg-new",
			wantErr:    false,
		},
		{
			name: "missing name",
			request: SecurityGroupCreateRequest{
				Description: helpers.StrPtr("invalid"),
			},
			response:   `{"error": "name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/network/v0/security_groups", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req SecurityGroupCreateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, tt.request.Name, req.Name)
				assertEqual(t, *tt.request.Description, *req.Description)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSecurityGroupClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestSecurityGroupService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "sg1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent security group",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, fmt.Sprintf("/network/v0/security_groups/%s", tt.id), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := testSecurityGroupClient(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func testSecurityGroupClient(baseURL string) SecurityGroupService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).SecurityGroups()
}
