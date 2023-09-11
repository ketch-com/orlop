package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithID returns a new context with the given request ID
func WithID(parent context.Context, requestID string) context.Context {
	if len(requestID) == 0 {
		return parent
	}

	return WithValue(parent, IDKey, requestID)
}

// WithIDOption returns a function that sets the given ID on a context
func WithIDOption(id string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(ID(ctx)) == 0 {
			ctx = WithID(ctx, id)
		}
		return ctx
	}
}

// ID returns the request ID or an empty string
func ID(ctx Context) string {
	return Value[string](ctx, IDKey)
}

// RequireID returns the request ID or an error if not set
func RequireID(ctx Context) (string, error) {
	return RequireValue[string](ctx, IDKey, errors.Invalid)
}
