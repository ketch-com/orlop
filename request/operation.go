package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithOperation returns a new context with the given Operation
func WithOperation(parent context.Context, operation string) context.Context {
	if len(operation) == 0 {
		return parent
	}

	return WithValue(parent, OperationKey, operation)
}

// WithOperationOption returns a function that sets the given URL on a context
func WithOperationOption(operation string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(Operation(ctx)) == 0 {
			ctx = WithOperation(ctx, operation)
		}
		return ctx
	}
}

// Operation returns the Operation or an empty string
func Operation(ctx Context) string {
	return Value[string](ctx, OperationKey)
}

// RequireOperation returns the Operation or an error if not set
func RequireOperation(ctx Context) (string, error) {
	return RequireValue[string](ctx, OperationKey, errors.Invalid)
}
