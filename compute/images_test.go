package compute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       ImageListOptions
		response   *string
		statusCode int
		want       int
		wantErr    bool
		checkQuery func(*testing.T, *http.Request)
	}{
		{
			name: "basic list",
			opts: ImageListOptions{},
			response: strPtr(`{
				"images": [
					{"id": "img1", "name": "ubuntu-20.04", "status": "active"},
					{"id": "img2", "name": "centos-8", "status": "active"}
				]
			}`),
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: ImageListOptions{
				Limit:  intPtr(1),
				Offset: intPtr(1),
			},
			response: strPtr(`{
				"images": [
					{"id": "img2", "name": "centos-8", "status": "active"}
				]
			}`),
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_limit") != "1" {
					t.Errorf("expected limit=1, got %s", r.URL.Query().Get("_limit"))
				}
				if r.URL.Query().Get("_offset") != "1" {
					t.Errorf("expected offset=1, got %s", r.URL.Query().Get("_offset"))
				}
			},
		},
		{
			name: "with sorting",
			opts: ImageListOptions{
				Sort: strPtr("platform:asc"),
			},
			response: strPtr(`{
				"images": [
					{"id": "img1", "name": "ubuntu-20.04", "status": "active"},
					{"id": "img2", "name": "centos-8", "status": "active"}
				]
			}`),
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_sort") != "platform:asc" {
					t.Errorf("expected sort=platform:asc, got %s", r.URL.Query().Get("_sort"))
				}
			},
		},
		{
			name: "with labels",
			opts: ImageListOptions{
				Labels: []string{"prod", "latest"},
			},
			response: strPtr(`{
				"images": [
					{"id": "img1", "name": "ubuntu-20.04", "status": "active", "labels": ["prod", "latest"]}
				]
			}`),
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_labels") != "prod,latest" {
					t.Errorf("expected labels=prod,latest, got %s", r.URL.Query().Get("_labels"))
				}
			},
		},
		{
			name: "with availability zone",
			opts: ImageListOptions{
				AvailabilityZone: strPtr("zone1"),
			},
			response: strPtr(`{
				"images": [
					{"id": "img1", "name": "ubuntu-20.04", "status": "active", "availability_zones": ["zone1"]}
				]
			}`),
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("availability-zone") != "zone1" {
					t.Errorf("expected availability-zone=zone1, got %s", r.URL.Query().Get("availability-zone"))
				}
			},
		},
		{
			name:       "server error",
			opts:       ImageListOptions{},
			response:   strPtr(`{"error": "internal server error"}`),
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "empty response",
			opts:       ImageListOptions{},
			response:   strPtr(""),
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "response is nil",
			opts:       ImageListOptions{},
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "malformed json",
			opts:       ImageListOptions{},
			response:   strPtr(`{"images": [{"id": "broken"}`),
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid pagination values",
			opts: ImageListOptions{
				Limit:  intPtr(-1),
				Offset: intPtr(-1),
			},
			response:   strPtr(`{"error": "invalid pagination parameters"}`),
			statusCode: http.StatusBadRequest,
			wantErr:    true,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("_limit") != "-1" {
					t.Errorf("expected limit=-1, got %s", r.URL.Query().Get("_limit"))
				}
				if r.URL.Query().Get("_offset") != "-1" {
					t.Errorf("expected offset=-1, got %s", r.URL.Query().Get("_offset"))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkQuery != nil {
					tt.checkQuery(t, r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(*tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.Images().List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("List() got %v images, want %v", len(got), tt.want)
			}
		})
	}
}

func TestImageService_Concurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"images": []}`))
	}))
	defer server.Close()

	client := testClient(server.URL)
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.Images().List(ctx, ImageListOptions{})
			if err != nil {
				t.Errorf("concurrent List() error = %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
