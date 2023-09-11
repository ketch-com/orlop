package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithUser returns a new context with the given user ID
func WithUser(parent context.Context, userID string) context.Context {
	if len(userID) == 0 {
		return parent
	}

	return WithValue(parent, UserKey, userID)
}

// WithUserOption returns a function that sets the given User on a context
func WithUserOption(userID string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(User(ctx)) == 0 {
			ctx = WithUser(ctx, userID)
		}
		return ctx
	}
}

// User returns the User ID or an empty string
func User(ctx Context) string {
	return Value[string](ctx, UserKey)
}

// RequireUser returns the request User or an error if not set
func RequireUser(ctx Context) (string, error) {
	return RequireValue[string](ctx, UserKey, errors.Forbidden)
}
