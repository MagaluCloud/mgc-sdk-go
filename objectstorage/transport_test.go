package objectstorage

import (
	"net/http"
	"testing"
)

func TestObjectStorageTransport_RoundTrip_ForceDeleteHeader(t *testing.T) {
	t.Parallel()
	var gotHeader string
	mockRT := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		gotHeader = req.Header.Get("X-Force-Container-Delete")
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	})
	tr := &objectStorageTransport{base: mockRT}
	req, _ := http.NewRequest(http.MethodDelete, "http://example.com", nil)
	req = req.WithContext(WithForceDelete(req.Context()))
	_, _ = tr.RoundTrip(req)
	if gotHeader != "true" {
		t.Errorf("expected X-Force-Container-Delete header to be 'true', got '%s'", gotHeader)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFixRetentionTime(t *testing.T) {
	t.Parallel()
	in := []byte(`{"date":"2024-05-27T12:00:00+0000"}`)
	want := []byte(`{"date":"2024-05-27T12:00:00+00:00"}`)
	got := fixRetentionTime(in)
	if string(got) != string(want) {
		t.Errorf("fixRetentionTime() = %s, want %s", got, want)
	}
}
