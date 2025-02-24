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

func TestRuleService_List(t *testing.T) {
	tests := []struct {
		name            string
		securityGroupID string
		response        string
		statusCode      int
		want            int
		wantErr         bool
	}{
		{
			name:            "successful list with rules",
			securityGroupID: "sg1",
			response: `{
				"rules": [
					{"id": "rule1", "direction": "ingress"},
					{"id": "rule2", "direction": "egress"}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:            "empty rule list",
			securityGroupID: "sg2",
			response:        `{"rules": []}`,
			statusCode:      http.StatusOK,
			want:            0,
			wantErr:         false,
		},
		{
			name:            "invalid security group",
			securityGroupID: "invalid",
			response:        `{"error": "security group not found"}`,
			statusCode:      http.StatusNotFound,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/network/v0/security_groups/%s/rules", tt.securityGroupID)
				assertEqual(t, expectedPath, r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testRulesClient(server.URL)
			rules, err := client.List(context.Background(), tt.securityGroupID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(rules))
		})
	}
}

func TestRuleService_Get(t *testing.T) {
	basetime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	parsedTime := utils.LocalDateTimeWithoutZone(basetime)

	tests := []struct {
		name       string
		ruleID     string
		response   string
		statusCode int
		want       *RuleResponse
		wantErr    bool
	}{
		{
			name:   "existing rule",
			ruleID: "rule1",
			response: `{
				"id": "rule1",
				"direction": "ingress",
				"protocol": "tcp",
				"port_range_min": 80,
				"port_range_max": 80,
				"created_at": "2024-01-01T00:00:00",
				"status": "ACTIVE"
			}`,
			statusCode: http.StatusOK,
			want: &RuleResponse{
				ID:           helpers.StrPtr("rule1"),
				Direction:    helpers.StrPtr("ingress"),
				Protocol:     helpers.StrPtr("tcp"),
				PortRangeMin: helpers.IntPtr(80),
				PortRangeMax: helpers.IntPtr(80),
				CreatedAt:    &parsedTime,
				Status:       "ACTIVE",
			},
			wantErr: false,
		},
		{
			name:       "non-existent rule",
			ruleID:     "invalid",
			response:   `{"error": "rule not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			ruleID:     "rule1",
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
				assertEqual(t, fmt.Sprintf("/network/v0/rules/%s", tt.ruleID), r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testRulesClient(server.URL)
			rule, err := client.Get(context.Background(), tt.ruleID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, *tt.want.ID, *rule.ID)
			assertEqual(t, *tt.want.Direction, *rule.Direction)
			assertEqual(t, *tt.want.Protocol, *rule.Protocol)
		})
	}
}

func TestRuleService_Create(t *testing.T) {
	tests := []struct {
		name            string
		securityGroupID string
		request         RuleCreateRequest
		response        string
		statusCode      int
		wantID          string
		wantErr         bool
	}{
		{
			name:            "successful create",
			securityGroupID: "sg1",
			request: RuleCreateRequest{
				Direction:    helpers.StrPtr("ingress"),
				Protocol:     helpers.StrPtr("tcp"),
				PortRangeMin: helpers.IntPtr(80),
				PortRangeMax: helpers.IntPtr(80),
				EtherType:    "IPv4",
			},
			response:   `{"id": "rule-new"}`,
			statusCode: http.StatusCreated,
			wantID:     "rule-new",
			wantErr:    false,
		},
		{
			name:            "invalid protocol",
			securityGroupID: "sg1",
			request: RuleCreateRequest{
				Direction: helpers.StrPtr("ingress"),
				Protocol:  helpers.StrPtr("invalid"),
			},
			response:   `{"error": "invalid protocol"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:            "invalid security group",
			securityGroupID: "invalid",
			request: RuleCreateRequest{
				Direction: helpers.StrPtr("ingress"),
				Protocol:  helpers.StrPtr("tcp"),
			},
			response:   `{"error": "security group not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/network/v0/security_groups/%s/rules", tt.securityGroupID)
				assertEqual(t, expectedPath, r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req RuleCreateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assertNoError(t, err)
				assertEqual(t, *tt.request.Direction, *req.Direction)
				assertEqual(t, *tt.request.Protocol, *req.Protocol)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testRulesClient(server.URL)
			id, err := client.Create(context.Background(), tt.securityGroupID, tt.request)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestRuleService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		ruleID     string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			ruleID:     "rule1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "non-existent rule",
			ruleID:     "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "rule not found"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			ruleID:     "rule1",
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
				assertEqual(t, fmt.Sprintf("/network/v0/rules/%s", tt.ruleID), r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testRulesClient(server.URL)
			err := client.Delete(context.Background(), tt.ruleID)

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func testRulesClient(baseURL string) RuleService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Rules()
}
