package objectstorage

import "context"

type forceDeleteKeyType struct{}

var forceDeleteKey = forceDeleteKeyType{}

func WithForceDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, forceDeleteKey, true)
}

func HasForceDelete(ctx context.Context) bool {
	v, ok := ctx.Value(forceDeleteKey).(bool)
	return ok && v
}
