package mgc_http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

const testProjectID = "550e8400-e29b-41d4-a716-446655440000"

type scopeProbe struct {
	mu      sync.Mutex
	seen    []string
	present []bool
	calls   atomic.Int32
	status  int
}

func (p *scopeProbe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	values, ok := r.Header[client.ProjectIDHeader]
	p.present = append(p.present, ok)
	if ok {
		p.seen = append(p.seen, values[0])
	} else {
		p.seen = append(p.seen, "")
	}
	p.mu.Unlock()

	status := p.status
	if status == 0 {
		status = http.StatusOK
	}
	if status != http.StatusOK && p.calls.Add(1) == 1 {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (p *scopeProbe) last() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	n := len(p.seen)
	if n == 0 {
		return "", false
	}
	return p.seen[n-1], p.present[n-1]
}

func newProbeClient(t *testing.T, probe *scopeProbe) (*client.Config, func()) {
	t.Helper()

	ts := httptest.NewServer(probe)
	core := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(ts.URL)),
		client.WithTimeout(5*time.Second),
		client.WithRetryConfig(3, 10*time.Millisecond, 50*time.Millisecond, 1.5),
	)

	return core.GetConfig(), ts.Close
}

func requestFuncFor(cfg *client.Config) NewRequestFunc {
	return func(ctx context.Context, method, path string, body any) (*http.Request, error) {
		return NewRequest[any](cfg, ctx, method, path, nil)
	}
}

func TestWithCallOptions_ProjectHeader(t *testing.T) {
	tests := []struct {
		name        string
		opts        client.CallOptions
		wantPresent bool
		wantValue   string
	}{
		{
			name:        "scoped call sends the project verbatim",
			opts:        client.CallOptions{ProjectID: testProjectID},
			wantPresent: true,
			wantValue:   testProjectID,
		},
		{
			name:        "zero scope sends no header",
			opts:        client.CallOptions{},
			wantPresent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			probe := &scopeProbe{}
			cfg, closeFn := newProbeClient(t, probe)
			defer closeFn()

			_, err := ExecuteSimpleRequestWithRespBody[map[string]string](
				context.Background(), WithCallOptions(requestFuncFor(cfg), tt.opts), cfg,
				http.MethodPost, "/test", nil, nil,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, present := probe.last()
			if present != tt.wantPresent {
				t.Fatalf("header present = %v, want %v", present, tt.wantPresent)
			}
			if tt.wantPresent && got != tt.wantValue {
				t.Errorf("header = %q, want %q", got, tt.wantValue)
			}
		})
	}
}

func TestWithCallOptions_GetKeepsProjectHeader(t *testing.T) {
	probe := &scopeProbe{}
	cfg, closeFn := newProbeClient(t, probe)
	defer closeFn()

	_, err := ExecuteSimpleRequestWithRespBody[map[string]string](
		context.Background(), WithCallOptions(requestFuncFor(cfg), client.CallOptions{ProjectID: testProjectID}), cfg,
		http.MethodGet, "/test", nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, present := probe.last()
	if !present || got != testProjectID {
		t.Errorf("GET header = %q (present %v), want %q", got, present, testProjectID)
	}
}

// S4: Do clones the request per attempt; the scope must reach the retry too.
func TestWithCallOptions_ProjectHeaderSurvivesRetry(t *testing.T) {
	probe := &scopeProbe{status: http.StatusInternalServerError}
	cfg, closeFn := newProbeClient(t, probe)
	defer closeFn()

	_, err := ExecuteSimpleRequestWithRespBody[map[string]string](
		context.Background(), WithCallOptions(requestFuncFor(cfg), client.CallOptions{ProjectID: testProjectID}), cfg,
		http.MethodPost, "/test", nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	probe.mu.Lock()
	defer probe.mu.Unlock()

	if len(probe.seen) < 2 {
		t.Fatalf("expected a retry, got %d request(s)", len(probe.seen))
	}
	for i, got := range probe.seen {
		if got != testProjectID {
			t.Errorf("attempt %d header = %q, want %q", i+1, got, testProjectID)
		}
	}
}

func TestWithCallOptions_ConcurrentScopesDoNotLeak(t *testing.T) {
	const otherProjectID = "11111111-2222-4333-8444-555555555555"

	seen := make(chan string, 2)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen <- r.Header.Get(client.ProjectIDHeader)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer ts.Close()

	core := client.NewMgcClient(
		client.WithAPIKey("test-api-key"),
		client.WithBaseURL(client.MgcUrl(ts.URL)),
		client.WithTimeout(5*time.Second),
	)
	cfg := core.GetConfig()

	var wg sync.WaitGroup
	for _, projectID := range []string{testProjectID, otherProjectID} {
		wg.Add(1)
		go func(projectID string) {
			defer wg.Done()
			_, _ = ExecuteSimpleRequestWithRespBody[map[string]string](
				context.Background(), WithCallOptions(requestFuncFor(cfg), client.CallOptions{ProjectID: projectID}), cfg,
				http.MethodGet, "/test", nil, nil,
			)
		}(projectID)
	}
	wg.Wait()
	close(seen)

	got := map[string]bool{}
	for v := range seen {
		got[v] = true
	}
	if !got[testProjectID] || !got[otherProjectID] {
		t.Errorf("each call must carry its own project, server saw %v", got)
	}
}

func TestWithCallOptions_NoRespBodySendsProjectHeader(t *testing.T) {
	probe := &scopeProbe{}
	cfg, closeFn := newProbeClient(t, probe)
	defer closeFn()

	err := ExecuteSimpleRequest(
		context.Background(), WithCallOptions(requestFuncFor(cfg), client.CallOptions{ProjectID: testProjectID}), cfg,
		http.MethodDelete, "/test", nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, present := probe.last()
	if !present || got != testProjectID {
		t.Errorf("header = %q (present %v), want %q", got, present, testProjectID)
	}
}

func TestWithCallOptions_PropagatesRequestError(t *testing.T) {
	wantErr := errors.New("build failed")
	reqf := WithCallOptions(func(context.Context, string, string, any) (*http.Request, error) {
		return nil, wantErr
	}, client.CallOptions{ProjectID: testProjectID})

	req, err := reqf(context.Background(), http.MethodGet, "/test", nil)
	if req != nil || !errors.Is(err, wantErr) {
		t.Fatalf("got (%v, %v), want (nil, %v)", req, err, wantErr)
	}
}

func TestExecuteSimpleRequest_SendsNoProjectHeader(t *testing.T) {
	probe := &scopeProbe{}
	cfg, closeFn := newProbeClient(t, probe)
	defer closeFn()

	_, err := ExecuteSimpleRequestWithRespBody[map[string]string](
		context.Background(), requestFuncFor(cfg), cfg,
		http.MethodGet, "/test", nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := probe.last(); present {
		t.Error("unscoped helper must not send the project header")
	}
}
