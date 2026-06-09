package containerregistry

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersService_Create_ProvisionsAuthenticatedUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/container-registry/v0/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if len(body) > 0 && string(body) != "null" {
			t.Errorf("expected empty body, got %q", string(body))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": "e3d49354-35d7-4565-b634-65d8b86aa594",
			"username": "e3d49354-35d7-4565-b634-65d8b86aa594",
			"created_at": "2024-05-15T19:56:47Z"
		}`))
	}))
	defer server.Close()

	got, err := testClient(server.URL).Users().Create(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "e3d49354-35d7-4565-b634-65d8b86aa594" {
		t.Errorf("ID = %q", got.ID)
	}
	if got.Username != "e3d49354-35d7-4565-b634-65d8b86aa594" {
		t.Errorf("Username = %q", got.Username)
	}
	if got.CreatedAt != "2024-05-15T19:56:47Z" {
		t.Errorf("CreatedAt = %q", got.CreatedAt)
	}
}

func TestUsersService_Create_PropagatesServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"slug":"unauthorized","message":"invalid token"}`))
	}))
	defer server.Close()

	got, err := testClient(server.URL).Users().Create(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got != nil {
		t.Errorf("expected nil result on error, got %+v", got)
	}
}

func TestUsersService_Get_ReturnsAuthenticatedUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/container-registry/v0/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "e3d49354-35d7-4565-b634-65d8b86aa594",
			"username": "e3d49354-35d7-4565-b634-65d8b86aa594",
			"created_at": "2024-05-15T19:56:47Z"
		}`))
	}))
	defer server.Close()

	got, err := testClient(server.URL).Users().Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "e3d49354-35d7-4565-b634-65d8b86aa594" {
		t.Errorf("ID = %q", got.ID)
	}
}

func TestUsersService_Delete_SendsUserIDInPath(t *testing.T) {
	const userID = "e3d49354-35d7-4565-b634-65d8b86aa594"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/container-registry/v0/users/"+userID {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := testClient(server.URL).Users().Delete(context.Background(), userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUsersService_Delete_PropagatesServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"slug":"not_found","message":"user not found"}`))
	}))
	defer server.Close()

	err := testClient(server.URL).Users().Delete(context.Background(), "missing-id")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
