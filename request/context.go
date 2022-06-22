package request

import (
	"context"
	"time"
)

var requestIDKey = struct{}{}
var requestURLKey = struct{}{}
var requestTimestampKey = struct{}{}
var requestTenantKey = struct{}{}

func ID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if r, ok := v.(string); ok {
			return r
		}
	}

	return ""
}

func WithID(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, requestIDKey, requestID)
}

func URL(ctx context.Context) string {
	if v := ctx.Value(requestURLKey); v != nil {
		if r, ok := v.(string); ok {
			return r
		}
	}

	return ""
}

func WithURL(parent context.Context, requestURL string) context.Context {
	return context.WithValue(parent, requestURLKey, requestURL)
}

func Timestamp(ctx context.Context) time.Time {
	if v := ctx.Value(requestTimestampKey); v != nil {
		if r, ok := v.(time.Time); ok {
			return r
		}
	}

	return time.Time{}
}

func WithTimestamp(parent context.Context, requestTimestamp time.Time) context.Context {
	return context.WithValue(parent, requestTimestampKey, requestTimestamp)
}

func Tenant(ctx context.Context) string {
	if v := ctx.Value(requestTenantKey); v != nil {
		if r, ok := v.(string); ok {
			return r
		}
	}

	return ""
}

func WithTenant(parent context.Context, requestTenant string) context.Context {
	return context.WithValue(parent, requestTenantKey, requestTenant)
}
