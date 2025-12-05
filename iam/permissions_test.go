package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPermissionService_ProductsAndPermissions(t *testing.T) {
	tests := []struct {
		name        string
		productName *string
		response    string
		statusCode  int
		want        int
		wantErr     bool
	}{
		{
			name: "successful get products and permissions",
			response: `[
				{
					"name": "compute",
					"permissions": [
						{"name": "read:instances", "description": "Read instances"},
						{"name": "write:instances", "description": "Write instances"}
					]
				},
				{
					"name": "network",
					"permissions": [
						{"name": "read:networks", "description": "Read networks"}
					]
				}
			]`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name:        "successful get with product name filter",
			productName: strPtr("compute"),
			response: `[
				{
					"name": "compute",
					"permissions": [
						{"name": "read:instances", "description": "Read instances"}
					]
				}
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
				if tt.productName != nil {
					if r.URL.Query().Get("product_name") != *tt.productName {
						t.Errorf("Expected product_name query param %s, got %s", *tt.productName, r.URL.Query().Get("product_name"))
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			result, err := client.Permissions().ProductsAndPermissions(context.Background(), tt.productName)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProductsAndPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.want {
				t.Errorf("ProductsAndPermissions() got = %d, want %d", len(result), tt.want)
			}
		})
	}
}
