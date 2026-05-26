package containerregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testRegistryID = "0c6bbd87-881a-4bf8-b3f6-4ff3ceacb42c"
	testMemberID   = "ef376048-a0f7-446c-826b-1d5cb3e5a5b0"
	testUserID     = "e3d49354-35d7-4565-b634-65d8b86aa594"
)

func TestMembersService_Add_PostsMemberToRegistry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/container-registry/v0/registries/"+testRegistryID+"/members" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body MemberRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.UserID != testUserID {
			t.Errorf("UserID = %q", body.UserID)
		}
		if body.Role == nil || *body.Role != "developer" {
			t.Errorf("Role = %v", body.Role)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{
			"id": %q,
			"registry_id": %q,
			"user_id": %q,
			"role": "developer",
			"created_at": "2024-05-20T14:30:00Z",
			"updated_at": "2024-05-20T14:30:00Z"
		}`, testMemberID, testRegistryID, testUserID)))
	}))
	defer server.Close()

	role := "developer"
	got, err := testClient(server.URL).Members().Add(context.Background(), testRegistryID, MemberRequest{
		UserID: testUserID,
		Role:   &role,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != testMemberID || got.Role != "developer" {
		t.Errorf("unexpected member: %+v", got)
	}
}

func TestMembersService_List_SendsPaginationParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		q := r.URL.Query()
		if q.Get("_limit") != "25" || q.Get("_offset") != "50" || q.Get("_sort") != "created_at:asc" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"results": [{
				"id": "m-1", "registry_id": "r-1", "user_id": "u-1",
				"role": "guest", "created_at": "t", "updated_at": "t"
			}],
			"meta": {"page": {"count": 1, "limit": 25, "offset": 50, "total": 1}}
		}`))
	}))
	defer server.Close()

	limit, offset := 25, 50
	sort := "created_at:asc"
	got, err := testClient(server.URL).Members().List(context.Background(), testRegistryID, MemberListOptions{
		Limit:  &limit,
		Offset: &offset,
		MemberFilterOptions: MemberFilterOptions{
			Sort: &sort,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Results) != 1 || got.Results[0].ID != "m-1" {
		t.Errorf("unexpected results: %+v", got.Results)
	}
}

func TestMembersService_ListAll_AggregatesPagesUntilExhausted(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		offset := r.URL.Query().Get("_offset")
		switch offset {
		case "0":
			w.WriteHeader(http.StatusOK)
			// Returns exactly `limit` items, signaling more pages.
			results := make([]string, 50)
			for i := range results {
				results[i] = fmt.Sprintf(`{"id":"m-%d","registry_id":"r","user_id":"u","role":"guest","created_at":"t","updated_at":"t"}`, i)
			}
			fmt.Fprintf(w, `{"results":[%s],"meta":{"page":{"count":50,"limit":50,"offset":0,"total":75}}}`,
				joinJSON(results))
		case "50":
			w.WriteHeader(http.StatusOK)
			results := make([]string, 25)
			for i := range results {
				results[i] = fmt.Sprintf(`{"id":"m-%d","registry_id":"r","user_id":"u","role":"guest","created_at":"t","updated_at":"t"}`, 50+i)
			}
			fmt.Fprintf(w, `{"results":[%s],"meta":{"page":{"count":25,"limit":50,"offset":50,"total":75}}}`,
				joinJSON(results))
		default:
			t.Errorf("unexpected offset: %s", offset)
		}
	}))
	defer server.Close()

	all, err := testClient(server.URL).Members().ListAll(context.Background(), testRegistryID, MemberFilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 75 {
		t.Errorf("expected 75 members, got %d", len(all))
	}
	if calls != 2 {
		t.Errorf("expected 2 page fetches, got %d", calls)
	}
}

func TestMembersService_Get_BuildsPathWithRegistryAndMemberID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		want := "/container-registry/v0/registries/" + testRegistryID + "/members/" + testMemberID
		if r.URL.Path != want {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"id": %q, "registry_id": %q, "user_id": %q,
			"role": "guest", "created_at": "t", "updated_at": "t"
		}`, testMemberID, testRegistryID, testUserID)))
	}))
	defer server.Close()

	got, err := testClient(server.URL).Members().Get(context.Background(), testRegistryID, testMemberID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != testMemberID {
		t.Errorf("ID = %q", got.ID)
	}
}

func TestMembersService_Update_PatchesRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var req MemberUpdateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if req.Role != "developer" {
			t.Errorf("Role = %q", req.Role)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"id": %q, "registry_id": %q, "user_id": %q,
			"role": "developer", "created_at": "t", "updated_at": "t"
		}`, testMemberID, testRegistryID, testUserID)))
	}))
	defer server.Close()

	got, err := testClient(server.URL).Members().Update(
		context.Background(), testRegistryID, testMemberID,
		MemberUpdateRequest{Role: "developer"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Role != "developer" {
		t.Errorf("Role = %q", got.Role)
	}
}

func TestMembersService_Delete_RemovesMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := testClient(server.URL).Members().Delete(context.Background(), testRegistryID, testMemberID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMembersService_Add_PropagatesServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"slug":"bad_request","message":"invalid role"}`))
	}))
	defer server.Close()

	_, err := testClient(server.URL).Members().Add(context.Background(), testRegistryID, MemberRequest{UserID: testUserID})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// joinJSON concatenates JSON object strings with commas (helper for paginated fixtures).
func joinJSON(items []string) string {
	out := ""
	for i, it := range items {
		if i > 0 {
			out += ","
		}
		out += it
	}
	return out
}
