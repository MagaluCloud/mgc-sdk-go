package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScopeService_GetGroupsAndProductsAndScopes(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "successful get groups and products and scopes",
			response: `[
				{
					"uuid": "group-uuid-1",
					"name": "Compute Group",
					"api_products": [
						{
							"uuid": "product-uuid-1",
							"name": "compute",
							"scopes": [
								{"uuid": "scope-uuid-1", "name": "read:instances", "title": "Read Instances"},
								{"uuid": "scope-uuid-2", "name": "write:instances", "title": "Write Instances"}
							]
						}
					]
				},
				{
					"uuid": "group-uuid-2",
					"name": "Network Group",
					"api_products": [
						{
							"uuid": "product-uuid-2",
							"name": "network",
							"scopes": [
								{"uuid": "scope-uuid-3", "name": "read:networks", "title": "Read Networks"}
							]
						}
					]
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
			result, err := client.Scopes().GroupsAndProductsAndScopes(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroupsAndProductsAndScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("GetGroupsAndProductsAndScopes() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}
