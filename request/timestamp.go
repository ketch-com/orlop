package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
	"time"
)

// WithTimestamp returns a new context with the given request timestamp
func WithTimestamp(parent context.Context, requestTimestamp time.Time) context.Context {
	if requestTimestamp.IsZero() {
		requestTimestamp = time.Now().UTC()
	}

	return WithValue(parent, TimestampKey, requestTimestamp)
}

// WithTimestampOption returns a function to set the given timestamp
func WithTimestampOption(requestTimestamp time.Time) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if Timestamp(ctx).IsZero() {
			ctx = WithTimestamp(ctx, requestTimestamp)
		}
		return ctx
	}
}

// Timestamp returns the request timestamp or an empty time.Time
func Timestamp(ctx Context) time.Time {
	return Value[time.Time](ctx, TimestampKey)
}

// RequireTimestamp returns the request timestamp or an error if not set
func RequireTimestamp(ctx Context) (time.Time, error) {
	return RequireValue[time.Time](ctx, TimestampKey, errors.Invalid)
}
