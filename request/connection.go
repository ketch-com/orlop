package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithConnection returns a new context with the given request Connection
func WithConnection(parent context.Context, connectionID string) context.Context {
	if len(connectionID) == 0 {
		return parent
	}

	return WithValue(parent, ConnectionKey, connectionID)
}

// WithConnectionOption returns a function that sets the given Connection on a context
func WithConnectionOption(connection string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(Connection(ctx)) == 0 {
			ctx = WithConnection(ctx, connection)
		}
		return ctx
	}
}

// Connection returns the request ID or an empty string
func Connection(ctx Context) string {
	return Value[string](ctx, ConnectionKey)
}

// RequireConnection returns the request ID or an error if not set
func RequireConnection(ctx Context) (string, error) {
	return RequireValue[string](ctx, ConnectionKey, errors.Invalid)
}
