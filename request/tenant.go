package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
	"strings"
)

// WithTenant returns a new context with the given request tenant
func WithTenant(parent context.Context, requestTenant string) context.Context {
	if len(requestTenant) == 0 {
		return parent
	}

	return WithValue(parent, TenantKey, requestTenant)
}

// WithTenantOption returns a function to set the given tenant on a context
func WithTenantOption(tenant string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(Tenant(ctx)) == 0 {
			ctx = WithTenant(ctx, tenant)
		}
		return ctx
	}
}

// Tenant returns the request Tenant or an empty string
func Tenant(ctx Context) string {
	tenant := Value[string](ctx, TenantKey)
	parts := strings.Split(tenant, ";")
	return parts[0]
}

// RequireTenant returns the request Tenant or an error if not set
func RequireTenant(ctx Context) (string, error) {
	return RequireValue[string](ctx, TenantKey, errors.Forbidden)
}
