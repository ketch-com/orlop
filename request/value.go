package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
)

// WithValue returns a new context with the given request value
func WithValue[T any](parent context.Context, key Key, v T) context.Context {
	return context.WithValue(parent, key, v)
}

// Value returns the value of a request context key, or a default value if the value is not set
func Value[T any](ctx Context, key Key) T {
	if v := ctx.Value(key); v != nil {
		if r, ok := v.(T); ok {
			return r
		}
	}

	var v T
	return v
}

// RequireValue returns the value of a request context key or returns an error if the value is not set
func RequireValue[T any](ctx Context, key Key, errDecorator func(error) error) (T, error) {
	if v := ctx.Value(key); v != nil {
		if r, ok := v.(T); ok {
			return r, nil
		}
	}

	var v T
	return v, errDecorator(errors.Errorf("%v not specified", key))
}

// Values returns a map of the request values from the context
func Values(ctx Context, opts ...Option) map[string]string {
	var o options

	for _, opt := range opts {
		opt(&o)
	}

	out := make(map[string]string)

	for k, getter := range Getters {
		if s := getter(ctx); len(s) > 0 {
			skip := false

			for _, filter := range o.filters {
				if !filter(k) {
					skip = true
				}
			}

			if !skip {
				out[string(k)] = s
			}
		}
	}

	return out
}

// Option is a function that sets values on the options structure
type Option func(o *options)

type options struct {
	filters []Filter
}

// Filter is a function that returns true if the given key should be included
type Filter func(k Key) bool

// WithFilter returns a filter that filters request values
func WithFilter(f Filter) Option {
	return func(o *options) {
		o.filters = append(o.filters, f)
	}
}

// SkipHighCardinalityKeysFilter returns true if the key is low cardinality
func SkipHighCardinalityKeysFilter(k Key) bool {
	if v, ok := LowCardinalityKeys[k]; ok {
		return v
	}

	// If we don't know about the key, assume it is "high cardinality"
	return false
}
