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

type fixRetentionTimeKeyType struct{}

var fixRetentionTimeKey = fixRetentionTimeKeyType{}

func WithFixRetentionTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, fixRetentionTimeKey, true)
}

func HasFixRetentionTime(ctx context.Context) bool {
	v, ok := ctx.Value(fixRetentionTimeKey).(bool)
	return ok && v
}

type storageClassKeyType struct{}

var storageClassKey = storageClassKeyType{}

func WithStorageClass(ctx context.Context, storageClass string) context.Context {
	return context.WithValue(ctx, storageClassKey, storageClass)
}

func HasStorageClass(ctx context.Context) bool {
	v, ok := ctx.Value(storageClassKey).(string)
	return ok && v != ""
}
