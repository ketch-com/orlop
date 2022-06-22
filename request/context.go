package request

import (
	"context"
	"time"
)

type Key string

var (
	IDKey        Key = "request_id"
	URLKey       Key = "request_url"
	TimestampKey Key = "request_timestamp"
	TenantKey    Key = "request_tenant"
)

func Value[T any](ctx context.Context, key Key) T {
	if v := ctx.Value(key); v != nil {
		if r, ok := v.(T); ok {
			return r
		}
	}

	var v T
	return v
}

func WithValue[T any](parent context.Context, key Key, v T) context.Context {
	return context.WithValue(parent, key, v)
}

func ID(ctx context.Context) string {
	return Value[string](ctx, IDKey)
}

func WithID(parent context.Context, requestID string) context.Context {
	return WithValue(parent, IDKey, requestID)
}

func URL(ctx context.Context) string {
	return Value[string](ctx, URLKey)
}

func WithURL(parent context.Context, requestURL string) context.Context {
	return WithValue(parent, URLKey, requestURL)
}

func Timestamp(ctx context.Context) time.Time {
	return Value[time.Time](ctx, TimestampKey)
}

func WithTimestamp(parent context.Context, requestTimestamp time.Time) context.Context {
	return WithValue(parent, TimestampKey, requestTimestamp)
}

func Tenant(ctx context.Context) string {
	return Value[string](ctx, TenantKey)
}

func WithTenant(parent context.Context, requestTenant string) context.Context {
	return WithValue(parent, TenantKey, requestTenant)
}
