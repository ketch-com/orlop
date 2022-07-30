// Copyright (c) 2021 Ketch Kloud, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package request

import (
	"context"
	"go.ketch.com/lib/orlop/v2/errors"
	"time"
)

// Key is the key of a request property in context
type Key string

var (
	IDKey         Key = "request_id"
	OperationKey  Key = "operation"
	TimestampKey  Key = "request_ts"
	TenantKey     Key = "tenant"
	URLKey        Key = "request_url"
	OriginatorKey Key = "request_originator"
)

// AllKeys is a slice of all Keys
var AllKeys = []Key{
	IDKey,
	OperationKey,
	TimestampKey,
	TenantKey,
	URLKey,
	OriginatorKey,
}

// HighCardinalityKeys is a map of high-cardinality keys
var HighCardinalityKeys = map[Key]bool{
	IDKey:         true,
	OperationKey:  false,
	TimestampKey:  true,
	TenantKey:     false,
	URLKey:        false,
	OriginatorKey: false,
}

// Setter is a function that adds a string to the context
type Setter func(ctx context.Context, v string) context.Context

// Getter is a function that returns a string from the context
type Getter func(ctx context.Context) string

// Setters is a map from the Key to a Setter for that Key
var Setters = map[Key]Setter{
	IDKey:         WithID,
	OperationKey:  WithOperation,
	TenantKey:     WithTenant,
	URLKey:        WithURL,
	OriginatorKey: WithOriginator,
	TimestampKey: func(ctx context.Context, v string) context.Context {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return WithTimestamp(ctx, t)
		}
		return ctx
	},
}

// Getters is a map from the Key to a Getter for that Key
var Getters = map[Key]Getter{
	IDKey:         ID,
	OperationKey:  Operation,
	TenantKey:     Tenant,
	URLKey:        URL,
	OriginatorKey: Originator,
	TimestampKey: func(ctx context.Context) string {
		if ts := Timestamp(ctx); !ts.IsZero() {
			return ts.Format(time.RFC3339)
		}
		return ""
	},
}

// Value returns the value of a request context key, or a default value if the value is not set
func Value[T any](ctx context.Context, key Key) T {
	if v := ctx.Value(key); v != nil {
		if r, ok := v.(T); ok {
			return r
		}
	}

	var v T
	return v
}

// RequireValue returns the value of a request context key or returns an error if the value is not set
func RequireValue[T any](ctx context.Context, key Key, errDecorator func(error) error) (T, error) {
	if v := ctx.Value(key); v != nil {
		if r, ok := v.(T); ok {
			return r, nil
		}
	}

	var v T
	return v, errDecorator(errors.Errorf("%v not specified", key))
}

// WithValue returns a new context with the given request value
func WithValue[T any](parent context.Context, key Key, v T) context.Context {
	return context.WithValue(parent, key, v)
}

// ID returns the request ID or an empty string
func ID(ctx context.Context) string {
	return Value[string](ctx, IDKey)
}

// RequireID returns the request ID or an error if not set
func RequireID(ctx context.Context) (string, error) {
	return RequireValue[string](ctx, IDKey, errors.Invalid)
}

// WithID returns a new context with the given request ID
func WithID(parent context.Context, requestID string) context.Context {
	return WithValue(parent, IDKey, requestID)
}

// Originator returns the request originator or an empty string
func Originator(ctx context.Context) string {
	return Value[string](ctx, OriginatorKey)
}

// RequireOriginator returns the request originator or an error if not set
func RequireOriginator(ctx context.Context) (string, error) {
	return RequireValue[string](ctx, OriginatorKey, errors.Invalid)
}

// WithOriginator returns a new context with the given request originator
func WithOriginator(parent context.Context, originator string) context.Context {
	return WithValue(parent, OriginatorKey, originator)
}

// URL returns the request URL or an empty string
func URL(ctx context.Context) string {
	return Value[string](ctx, URLKey)
}

// RequireURL returns the request URL or an error if not set
func RequireURL(ctx context.Context) (string, error) {
	return RequireValue[string](ctx, URLKey, errors.Invalid)
}

// WithURL returns a new context with the given request URL
func WithURL(parent context.Context, requestURL string) context.Context {
	return WithValue(parent, URLKey, requestURL)
}

// Timestamp returns the request timestamp or an empty time.Time
func Timestamp(ctx context.Context) time.Time {
	return Value[time.Time](ctx, TimestampKey)
}

// RequireTimestamp returns the request timestamp or an error if not set
func RequireTimestamp(ctx context.Context) (time.Time, error) {
	return RequireValue[time.Time](ctx, TimestampKey, errors.Invalid)
}

// WithTimestamp returns a new context with the given request timestamp
func WithTimestamp(parent context.Context, requestTimestamp time.Time) context.Context {
	return WithValue(parent, TimestampKey, requestTimestamp)
}

// Tenant returns the request Tenant or an empty string
func Tenant(ctx context.Context) string {
	return Value[string](ctx, TenantKey)
}

// RequireTenant returns the request Tenant or an error if not set
func RequireTenant(ctx context.Context) (string, error) {
	return RequireValue[string](ctx, TenantKey, errors.Forbidden)
}

// WithTenant returns a new context with the given request tenant
func WithTenant(parent context.Context, requestTenant string) context.Context {
	return WithValue(parent, TenantKey, requestTenant)
}

// Operation returns the Operation or an empty string
func Operation(ctx context.Context) string {
	return Value[string](ctx, OperationKey)
}

// RequireOperation returns the Operation or an error if not set
func RequireOperation(ctx context.Context) (string, error) {
	return RequireValue[string](ctx, OperationKey, errors.Invalid)
}

// WithOperation returns a new context with the given Operation
func WithOperation(parent context.Context, operation string) context.Context {
	return WithValue(parent, OperationKey, operation)
}

// Values returns a map of the request values from the context
func Values(ctx context.Context, opts ...Option) map[string]string {
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

// SkipHighCardinalityKeysFilter returns true if the key is not high cardinality
func SkipHighCardinalityKeysFilter(k Key) bool {
	return !HighCardinalityKeys[k]
}
