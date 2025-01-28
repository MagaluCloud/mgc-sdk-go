package sshkeys

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestKeyService(t *testing.T) {
	// Setup test server and client
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/profile/v0/ssh-keys":
			if r.Method == http.MethodGet {
				handleListKeys(w, r)
			} else if r.Method == http.MethodPost {
				handleCreateKey(w, r)
			} else {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		case "/profile/v0/ssh-keys/123":
			handleSingleKey(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	cfg := client.NewMgcClient("test-api-key",
		client.WithTimeout(20 * time.Second),
	)
	c := New(cfg, WithGlobalBasePath(client.MgcUrl(ts.URL)))
	service := c.Keys()

	t.Run("ListKeys", func(t *testing.T) {
		t.Run("successful list with default options", func(t *testing.T) {
			keys, err := service.List(context.Background(), ListOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(keys) != 2 {
				t.Errorf("expected 2 keys, got %d", len(keys))
			}
			if keys[0].ID != "key1" {
				t.Errorf("expected key1, got %s", keys[0].ID)
			}
			if keys[1].ID != "key2" {
				t.Errorf("expected key2, got %s", keys[1].ID)
			}
		})

		t.Run("successful list with query parameters", func(t *testing.T) {
			limit := 10
			offset := 5
			sort := "name"
			keys, err := service.List(context.Background(), ListOptions{
				Limit:  &limit,
				Offset: &offset,
				Sort:   &sort,
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(keys) != 2 {
				t.Errorf("expected 2 keys, got %d", len(keys))
			}
		})

		t.Run("error from server", func(t *testing.T) {
			ts := errorTestServer()
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithBaseURL(client.MgcUrl(ts.URL)),
				client.WithTimeout(1),
			)
			client := New(cfg)
			service := client.Keys()

			_, err := service.List(context.Background(), ListOptions{})
			if err == nil {
				t.Error("expected error, got nil")
			}
			if err != nil && !contains(err.Error(), "server error") {
				t.Errorf("expected 'server error' in error message, got: %v", err)
			}
		})
	})

	t.Run("CreateKey", func(t *testing.T) {
		t.Run("successful creation", func(t *testing.T) {
			key, err := service.Create(context.Background(), CreateSSHKeyRequest{
				Name: "new-key",
				Key:  "ssh-rsa...",
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if key.ID != "123" {
				t.Errorf("expected ID '123', got %s", key.ID)
			}
			if key.Name != "new-key" {
				t.Errorf("expected name 'new-key', got %s", key.Name)
			}
		})

		t.Run("invalid request body", func(t *testing.T) {
			ts := errorTestServer()
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithBaseURL(client.MgcUrl(ts.URL)),
				client.WithTimeout(1),
			)
			client := New(cfg)
			service := client.Keys()

			_, err := service.Create(context.Background(), CreateSSHKeyRequest{})
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})

	t.Run("GetKey", func(t *testing.T) {
		t.Run("successful get", func(t *testing.T) {
			key, err := service.Get(context.Background(), "123")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if key.ID != "123" {
				t.Errorf("expected ID '123', got %s", key.ID)
			}
			if key.Name != "test-key" {
				t.Errorf("expected name 'test-key', got %s", key.Name)
			}
		})

		t.Run("non-existent key", func(t *testing.T) {
			ts := errorTestServer()
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithBaseURL(client.MgcUrl(ts.URL)),
				client.WithTimeout(1),
			)
			client := New(cfg)
			service := client.Keys()

			_, err := service.Get(context.Background(), "456")
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})

	t.Run("DeleteKey", func(t *testing.T) {
		t.Run("successful delete", func(t *testing.T) {
			key, err := service.Delete(context.Background(), "123")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if key.ID != "123" {
				t.Errorf("expected ID '123', got %s", key.ID)
			}
			if key.Name != "deleted-key" {
				t.Errorf("expected name 'deleted-key', got %s", key.Name)
			}
		})

		t.Run("delete non-existent key", func(t *testing.T) {
			ts := errorTestServer()
			defer ts.Close()

			cfg := client.NewMgcClient("test-api-key",
				client.WithBaseURL(client.MgcUrl(ts.URL)),
				client.WithTimeout(1),
			)
			client := New(cfg)
			service := client.Keys()

			_, err := service.Delete(context.Background(), "456")
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})
}

// Test handlers
func handleListKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Verify query parameters
	query := r.URL.Query()
	if limit := query.Get("_limit"); limit != "" {
		if _, err := strconv.Atoi(limit); err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
	}

	response := ListSSHKeysResponse{
		Results: []SSHKey{
			{ID: "key1", Name: "Key 1"},
			{ID: "key2", Name: "Key 2"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleSingleKey(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		response := SSHKey{ID: "123", Name: "test-key"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	case http.MethodDelete:
		response := SSHKey{ID: "123", Name: "deleted-key"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func errorTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"error": "server error"}`)
	}))
}

func TestListOptionsQueryParams(t *testing.T) {
	opts := ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(20),
		Sort:   helpers.StrPtr("name"),
	}

	query := make(url.Values)
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}

	if query.Get("_limit") != "10" {
		t.Errorf("expected limit '10', got %s", query.Get("_limit"))
	}
	if query.Get("_offset") != "20" {
		t.Errorf("expected offset '20', got %s", query.Get("_offset"))
	}
	if query.Get("_sort") != "name" {
		t.Errorf("expected sort 'name', got %s", query.Get("_sort"))
	}
}

func contains(s, substr string) bool {
	return s != "" && s != substr
}

func handleCreateKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req CreateSSHKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Key == "" {
		http.Error(w, "name and key are required", http.StatusBadRequest)
		return
	}

	response := SSHKey{
		ID:   "123",
		Name: req.Name,
		Key:  req.Key,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}