package client

import (
	"net/http"
	"testing"
)

func TestCallOptionsApply(t *testing.T) {
	tests := []struct {
		name       string
		opts       CallOptions
		wantHeader string
		wantSet    bool
	}{
		{
			name:       "project id sets the header verbatim",
			opts:       CallOptions{ProjectID: "550e8400-e29b-41d4-a716-446655440000"},
			wantHeader: "550e8400-e29b-41d4-a716-446655440000",
			wantSet:    true,
		},
		{
			name:    "zero value omits the header",
			opts:    CallOptions{},
			wantSet: false,
		},
		{
			// An empty header is not the same as an absent one for the API:
			// absent means "default project", present-and-empty is a bad request.
			name:    "empty project id omits the header",
			opts:    CallOptions{ProjectID: ""},
			wantSet: false,
		},
		{
			// The SDK carries the value as given; normalizing is the caller's job.
			name:       "value is not trimmed",
			opts:       CallOptions{ProjectID: " abc "},
			wantHeader: " abc ",
			wantSet:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Fatalf("unexpected error building request: %v", err)
			}

			tt.opts.Apply(req)

			got, ok := req.Header[ProjectIDHeader]
			if ok != tt.wantSet {
				t.Fatalf("header %q present = %v, want %v", ProjectIDHeader, ok, tt.wantSet)
			}
			if tt.wantSet && got[0] != tt.wantHeader {
				t.Errorf("header %q = %q, want %q", ProjectIDHeader, got[0], tt.wantHeader)
			}
		})
	}
}

func TestCallOptionsApplyNilRequest(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Apply panicked on a nil request: %v", r)
		}
	}()

	CallOptions{ProjectID: "abc"}.Apply(nil)
}
