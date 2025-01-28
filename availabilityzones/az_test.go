package availabilityzones

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestService(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		t.Run("successful list with ShowBlocked false", func(t *testing.T) {
			var capturedShowBlocked string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedShowBlocked = r.URL.Query().Get("show_is_blocked")
				response := ListResponse{
					Results: []Region{
						{
							ID: "region1",
							AvailabilityZones: []AvailabilityZone{
								{ID: "az1", BlockType: BlockTypeNone},
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}))
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithTimeout(20*time.Second),
			)
			c := New(cfg, WithGlobalBasePath(client.MgcUrl(ts.URL)))
			service := c.AvailabilityZones()

			_, err := service.List(context.Background(), ListOptions{ShowBlocked: false})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if capturedShowBlocked != "false" {
				t.Errorf("expected show_is_blocked=false, got %s", capturedShowBlocked)
			}
		})

		t.Run("successful list with ShowBlocked true", func(t *testing.T) {
			var capturedShowBlocked string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedShowBlocked = r.URL.Query().Get("show_is_blocked")
				response := ListResponse{
					Results: []Region{
						{
							ID: "region1",
							AvailabilityZones: []AvailabilityZone{
								{ID: "az1", BlockType: BlockTypeNone},
								{ID: "az2", BlockType: BlockTypeTotal},
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}))
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithTimeout(20*time.Second),
			)
			c := New(cfg, WithGlobalBasePath(client.MgcUrl(ts.URL)))
			service := c.AvailabilityZones()

			res, err := service.List(context.Background(), ListOptions{ShowBlocked: true})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if capturedShowBlocked != "true" {
				t.Errorf("expected show_is_blocked=true, got %s", capturedShowBlocked)
			}

			if len(res.Results) != 1 {
				t.Fatalf("expected 1 region, got %d", len(res.Results))
			}
			region := res.Results[0]
			if len(region.AvailabilityZones) != 2 {
				t.Errorf("expected 2 availability zones, got %d", len(region.AvailabilityZones))
			}
		})

		t.Run("successful response parsing", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := ListResponse{
					Results: []Region{
						{
							ID: "region1",
							AvailabilityZones: []AvailabilityZone{
								{ID: "az1", BlockType: BlockTypeNone},
								{ID: "az2", BlockType: BlockTypeReadOnly},
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}))
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithTimeout(20*time.Second),
			)
			c := New(cfg, WithGlobalBasePath(client.MgcUrl(ts.URL)))
			service := c.AvailabilityZones()

			res, err := service.List(context.Background(), ListOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(res.Results) != 1 {
				t.Fatalf("expected 1 region, got %d", len(res.Results))
			}
			region := res.Results[0]
			if region.ID != "region1" {
				t.Errorf("expected region ID 'region1', got %s", region.ID)
			}
			if len(region.AvailabilityZones) != 2 {
				t.Fatalf("expected 2 availability zones, got %d", len(region.AvailabilityZones))
			}
			az1 := region.AvailabilityZones[0]
			if az1.ID != "az1" || az1.BlockType != BlockTypeNone {
				t.Errorf("expected az1 with BlockTypeNone, got %s %s", az1.ID, az1.BlockType)
			}
			az2 := region.AvailabilityZones[1]
			if az2.ID != "az2" || az2.BlockType != BlockTypeReadOnly {
				t.Errorf("expected az2 with BlockTypeReadOnly, got %s %s", az2.ID, az2.BlockType)
			}
		})

		t.Run("error from server", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"error": "server error"}`)
			}))
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithTimeout(20*time.Second),
			)
			c := New(cfg, WithGlobalBasePath(client.MgcUrl(ts.URL)))
			service := c.AvailabilityZones()

			_, err := service.List(context.Background(), ListOptions{})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !contains(err.Error(), "Server Error") {
				t.Errorf("expected 'server error' in error message, got: %v", err)
			}
		})

		t.Run("timeout error", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(50 * time.Millisecond) // Simulate delay
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithTimeout(1*time.Millisecond), // Very short timeout to trigger error
			)
			c := New(cfg, WithGlobalBasePath(client.MgcUrl(ts.URL)))
			service := c.AvailabilityZones()

			_, err := service.List(context.Background(), ListOptions{})
			if err == nil {
				t.Fatal("expected timeout error, got nil")
			}
		})
	})
}

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
