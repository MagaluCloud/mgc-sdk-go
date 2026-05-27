package objectstorage

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

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

func TestObjectStorageTransport_BaseTransportError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("transport error")

	transport := &objectStorageTransport{
		base: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, expectedErr
		}),
	}

	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Errorf("failed to create request: %v", err)
	}

	_, err = transport.RoundTrip(req)
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestRoundTrip_DoNotFixRetentionTimeWithoutContext(t *testing.T) {
	original := `{"retention":"2026-05-29T00:12:12+0000"}`

	transport := &objectStorageTransport{
		base: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(original)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	req, err := http.NewRequest(
		http.MethodGet,
		"http://example.com?retention=true",
		nil,
	)
	if err != nil {
		t.Errorf("failed to create request: %v", err)
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error reading body: %v", err)
	}

	if string(body) != original {
		t.Errorf("expected original body %q, got %q", original, string(body))
	}
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
