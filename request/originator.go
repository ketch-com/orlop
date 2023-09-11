package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithOriginator returns a new context with the given request originator
func WithOriginator(parent context.Context, originator string) context.Context {
	if len(originator) == 0 {
		return parent
	}

	return WithValue(parent, OriginatorKey, originator)
}

// WithOriginatorOption returns a function that sets the given Originator on a context
func WithOriginatorOption(originator string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(Originator(ctx)) == 0 {
			ctx = WithOriginator(ctx, originator)
		}
		return ctx
	}
}

// Originator returns the request originator or an empty string
func Originator(ctx Context) string {
	return Value[string](ctx, OriginatorKey)
}

// RequireOriginator returns the request originator or an error if not set
func RequireOriginator(ctx Context) (string, error) {
	return RequireValue[string](ctx, OriginatorKey, errors.Invalid)
}
