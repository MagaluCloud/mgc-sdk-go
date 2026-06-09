package objectstorage

import (
	"context"
	"sync"
	"testing"
)

func TestWithForceDeleteAndHasForceDelete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	if HasForceDelete(ctx) {
		t.Error("HasForceDelete should be false for default context")
	}
	ctx2 := WithForceDelete(ctx)
	if !HasForceDelete(ctx2) {
		t.Error("HasForceDelete should be true after WithForceDelete")
	}
}

func TestWithFixRetentionTimeAndHasFixRetentionTime(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	if HasFixRetentionTime(ctx) {
		t.Error("HasFixRetentionTime should be false for default context")
	}
	ctx2 := WithFixRetentionTime(ctx)
	if !HasFixRetentionTime(ctx2) {
		t.Error("HasFixRetentionTime should be true after WithFixRetentionTime")
	}
}

func TestWithStorageClass_DefaultContextReturnsFalse(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	if HasStorageClass(ctx) {
		t.Error("HasStorageClass should be false for default context")
	}
}

func TestWithStorageClass_SetValue(t *testing.T) {
	t.Parallel()
	ctx := WithStorageClass(context.Background(), "STANDARD")
	if !HasStorageClass(ctx) {
		t.Error("HasStorageClass should be true after WithStorageClass")
	}
}

func TestWithStorageClass_EmptyValueReturnsFalse(t *testing.T) {
	t.Parallel()
	ctx := WithStorageClass(context.Background(), "")
	if HasStorageClass(ctx) {
		t.Error("HasStorageClass should be false for empty storage class")
	}
}

func TestWithStorageClass_OriginalContextUnchanged(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	_ = WithStorageClass(ctx, "COLD")
	if HasStorageClass(ctx) {
		t.Error("original context should not be modified by WithStorageClass")
	}
}

type mockProgressReporter struct {
	mu           sync.Mutex
	startCalled  bool
	addCalled    bool
	finishCalled bool
	startTotal   int64
	addDelta     int64
}

func (m *mockProgressReporter) Start(total int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.startCalled = true
	m.startTotal = total
}

func (m *mockProgressReporter) Add(delta int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.addCalled = true
	m.addDelta = delta
}

func (m *mockProgressReporter) Finish() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.finishCalled = true
}

func TestWithProgress_DefaultContextReturnsNil(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	if GetProgress(ctx) != nil {
		t.Error("GetProgress should return nil for default context")
	}
}

func TestWithProgress_SetAndGet(t *testing.T) {
	t.Parallel()
	p := &mockProgressReporter{}
	ctx := WithProgress(context.Background(), p)
	got := GetProgress(ctx)
	if got == nil {
		t.Fatal("GetProgress should return the reporter set by WithProgress")
	}
	if got != p {
		t.Error("GetProgress should return the exact reporter passed to WithProgress")
	}
}

func TestWithProgress_ReporterIsCallable(t *testing.T) {
	t.Parallel()
	p := &mockProgressReporter{}
	ctx := WithProgress(context.Background(), p)

	reporter := GetProgress(ctx)
	reporter.Start(100)
	reporter.Add(50)
	reporter.Finish()

	if !p.startCalled || p.startTotal != 100 {
		t.Errorf("Start() not called correctly: startCalled=%v, startTotal=%d", p.startCalled, p.startTotal)
	}
	if !p.addCalled || p.addDelta != 50 {
		t.Errorf("Add() not called correctly: addCalled=%v, addDelta=%d", p.addCalled, p.addDelta)
	}
	if !p.finishCalled {
		t.Error("Finish() was not called")
	}
}

func TestWithProgress_OriginalContextUnchanged(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := &mockProgressReporter{}
	_ = WithProgress(ctx, p)
	if GetProgress(ctx) != nil {
		t.Error("original context should not be modified by WithProgress")
	}
}
