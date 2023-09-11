package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithURL returns a new context with the given request URL
func WithURL(parent context.Context, requestURL string) context.Context {
	if len(requestURL) == 0 {
		return parent
	}

	return WithValue(parent, URLKey, requestURL)
}

// WithURLOption returns a function to set the given url on a context
func WithURLOption(url string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		if len(URL(ctx)) == 0 {
			ctx = WithURL(ctx, url)
		}
		return ctx
	}
}

// URL returns the request URL or an empty string
func URL(ctx Context) string {
	return Value[string](ctx, URLKey)
}

// RequireURL returns the request URL or an error if not set
func RequireURL(ctx Context) (string, error) {
	return RequireValue[string](ctx, URLKey, errors.Invalid)
}
