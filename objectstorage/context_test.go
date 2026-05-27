package objectstorage

import (
	"context"
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
