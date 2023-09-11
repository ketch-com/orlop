package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithIntegration returns a new context with the given request integration
func WithIntegration(parent context.Context, integration string) context.Context {
	if len(integration) == 0 {
		return parent
	}

	return WithValue(parent, IntegrationKey, integration)
}

// WithIntegrationOption returns a function that sets the given Integration on a context
func WithIntegrationOption(integration string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(Integration(ctx)) == 0 {
			ctx = WithIntegration(ctx, integration)
		}
		return ctx
	}
}

// Integration returns the request ID or an empty string
func Integration(ctx Context) string {
	return Value[string](ctx, IntegrationKey)
}

// RequireIntegration returns the request ID or an error if not set
func RequireIntegration(ctx Context) (string, error) {
	return RequireValue[string](ctx, IntegrationKey, errors.Invalid)
}
