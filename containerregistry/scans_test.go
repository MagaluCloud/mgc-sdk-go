package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	testRepositoryID = "0c6bbd87-881a-4bf8-b3f6-4ff3ceacb42c"
	testDigestOrTag  = "latest"
	testScanID       = "00000000-0000-0000-0000-000000000000"
)

func TestScansService_Schedule_PostsToImageScansEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		want := fmt.Sprintf("/container-registry/v1/registries/%s/repositories/%s/images/%s/scans",
			testRegistryID, testRepositoryID, testDigestOrTag)
		if r.URL.Path != want {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{
			"id": %q, "digest": "sha256:000000000",
			"registry_id": %q, "repository_id": %q,
			"input_tag": "latest", "status": "pending",
			"created_at": "2024-05-15T19:56:47Z", "updated_at": "2024-05-15T19:56:47Z"
		}`, testScanID, testRegistryID, testRepositoryID)))
	}))
	defer server.Close()

	got, err := testClient(server.URL).Scans().Schedule(context.Background(), ScheduleScanInput{
		RegistryID:   testRegistryID,
		RepositoryID: testRepositoryID,
		DigestOrTag:  testDigestOrTag,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != testScanID || got.Status != "pending" {
		t.Errorf("unexpected scan: %+v", got)
	}
}

func TestScansService_ListByImage_SendsStatusFilterAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("_limit") != "10" || q.Get("_offset") != "0" || q.Get("status") != "completed" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"results": [{
				"id": "s-1", "digest": "sha256:abc", "status": "completed",
				"created_at": "t", "updated_at": "t"
			}],
			"meta": {"page": {"count": 1, "limit": 10, "offset": 0, "total": 1}}
		}`))
	}))
	defer server.Close()

	limit, offset := 10, 0
	status := "completed"
	got, err := testClient(server.URL).Scans().ListByImage(
		context.Background(),
		ListByImageInput{RegistryID: testRegistryID, RepositoryID: testRepositoryID, DigestOrTag: testDigestOrTag},
		ListImageScansOptions{
			Limit:  &limit,
			Offset: &offset,
			ImageScansFilterOptions: ImageScansFilterOptions{
				Status: &status,
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(got.Results))
	}
}

func TestScansService_Get_ReturnsScanWithChildren(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/container-registry/v1/scans/"+testScanID {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"id": %q, "digest": "sha256:000000000", "status": "completed",
			"severity_summary": {"total": 10, "low": 3, "medium": 3, "high": 3, "critical": 1, "fixable": 6},
			"created_at": "t", "updated_at": "t",
			"child_scans": [
				{"id":"child-1","digest":"sha256:111","status":"completed","created_at":"t","updated_at":"t"}
			]
		}`, testScanID)))
	}))
	defer server.Close()

	got, err := testClient(server.URL).Scans().Get(context.Background(), testScanID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.SeveritySummary == nil || got.SeveritySummary.Total != 10 {
		t.Errorf("unexpected severity_summary: %+v", got.SeveritySummary)
	}
	if len(got.ChildScans) != 1 || got.ChildScans[0].ID != "child-1" {
		t.Errorf("unexpected child_scans: %+v", got.ChildScans)
	}
}

func TestScansService_ListVulnerabilities_RepeatsSeverityQueryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		gotSeverities := q["severity"]
		if !reflect.DeepEqual(gotSeverities, []string{"high", "critical"}) {
			t.Errorf("expected ?severity=high&severity=critical, got %v (raw=%s)", gotSeverities, r.URL.RawQuery)
		}
		if q.Get("package_name") != "openssl" || q.Get("cve_id") != "CVE-2024-1" || q.Get("fixable") != "true" {
			t.Errorf("missing scalar filters: %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"results": [{
				"cve_id": "CVE-2024-1", "severity": "high", "package_name": "openssl",
				"current_version": "1.0", "fixable": true, "is_allowlisted": false,
				"description": "x", "links": [], "cvss": {"preferred": null, "max": null, "by_source": {}}
			}],
			"meta": {"page": {"count": 1, "limit": 50, "offset": 0, "total": 1}}
		}`))
	}))
	defer server.Close()

	pkg, cve := "openssl", "CVE-2024-1"
	fixable := true
	got, err := testClient(server.URL).Scans().ListVulnerabilities(
		context.Background(), testScanID,
		ListVulnerabilitiesOptions{
			VulnerabilitiesFilterOptions: VulnerabilitiesFilterOptions{
				Severity:    []string{"high", "critical"},
				PackageName: &pkg,
				CveID:       &cve,
				Fixable:     &fixable,
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Results) != 1 || got.Results[0].CveID != "CVE-2024-1" {
		t.Errorf("unexpected results: %+v", got.Results)
	}
}

func TestScansService_Stop_PostsToStopEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/container-registry/v1/scans/"+testScanID+"/stop" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	if err := testClient(server.URL).Scans().Stop(context.Background(), testScanID); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestScansService_Schedule_PropagatesServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"slug":"not_found","message":"image not found"}`))
	}))
	defer server.Close()

	_, err := testClient(server.URL).Scans().Schedule(context.Background(), ScheduleScanInput{
		RegistryID: testRegistryID, RepositoryID: testRepositoryID, DigestOrTag: "missing",
	})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
