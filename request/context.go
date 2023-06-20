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
	"strings"
	"time"

	"go.ketch.com/lib/orlop/v2/errors"
)

// Key is the key of a request property in context
type Key string

var (
	ConnectionKey  Key = "connection"
	IDKey          Key = "requestId"
	OperationKey   Key = "operation"
	OriginatorKey  Key = "requestOriginator"
	TenantKey      Key = "tenant"
	TimestampKey   Key = "requestTS"
	URLKey         Key = "requestUrl"
	UserKey        Key = "userId"
	IntegrationKey Key = "integration"
)

// AllKeys is a slice of all Keys
var AllKeys = []Key{
	ConnectionKey,
	IDKey,
	OperationKey,
	OriginatorKey,
	TenantKey,
	TimestampKey,
	URLKey,
	UserKey,
	IntegrationKey,
}

// LowCardinalityKeys is a map of high-cardinality keys
var LowCardinalityKeys = map[Key]bool{
	OperationKey:  true,
	OriginatorKey: true,
	TenantKey:     true,
}

// Setter is a function that adds a string to the context
type Setter func(ctx context.Context, v string) context.Context

// Getter is a function that returns a string from the context
type Getter func(ctx Context) string

// Setters is a map from the Key to a Setter for that Key
var Setters = map[Key]Setter{
	ConnectionKey:  WithConnection,
	IDKey:          WithID,
	OperationKey:   WithOperation,
	OriginatorKey:  WithOriginator,
	TenantKey:      WithTenant,
	URLKey:         WithURL,
	UserKey:        WithUser,
	IntegrationKey: WithIntegration,
	TimestampKey: func(ctx context.Context, v string) context.Context {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return WithTimestamp(ctx, t)
		}
		return ctx
	},
}

// Getters is a map from the Key to a Getter for that Key
var Getters = map[Key]Getter{
	ConnectionKey:  Connection,
	IDKey:          ID,
	OperationKey:   Operation,
	OriginatorKey:  Originator,
	TenantKey:      Tenant,
	URLKey:         URL,
	UserKey:        User,
	IntegrationKey: Integration,
	TimestampKey: func(ctx Context) string {
		if ts := Timestamp(ctx); !ts.IsZero() {
			return ts.Format(time.RFC3339)
		}
		return ""
	},
}

type Context interface {
	Value(key any) any
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

// WithValue returns a new context with the given request value
func WithValue[T any](parent context.Context, key Key, v T) context.Context {
	return context.WithValue(parent, key, v)
}

// WithConnection returns a new context with the given request Connection
func WithConnection(parent context.Context, connectionID string) context.Context {
	return WithValue(parent, ConnectionKey, connectionID)
}

// Connection returns the request ID or an empty string
func Connection(ctx Context) string {
	return Value[string](ctx, ConnectionKey)
}

// RequireConnection returns the request ID or an error if not set
func RequireConnection(ctx Context) (string, error) {
	return RequireValue[string](ctx, ConnectionKey, errors.Invalid)
}

// WithIntegration returns a new context with the given request integration
func WithIntegration(parent context.Context, integration string) context.Context {
	return WithValue(parent, IntegrationKey, integration)
}

// Integration returns the request ID or an empty string
func Integration(ctx Context) string {
	return Value[string](ctx, IntegrationKey)
}

// RequireIntegration returns the request ID or an error if not set
func RequireIntegration(ctx Context) (string, error) {
	return RequireValue[string](ctx, IntegrationKey, errors.Invalid)
}

// WithID returns a new context with the given request ID
func WithID(parent context.Context, requestID string) context.Context {
	return WithValue(parent, IDKey, requestID)
}

// ID returns the request ID or an empty string
func ID(ctx Context) string {
	return Value[string](ctx, IDKey)
}

// RequireID returns the request ID or an error if not set
func RequireID(ctx Context) (string, error) {
	return RequireValue[string](ctx, IDKey, errors.Invalid)
}

// WithUser returns a new context with the given user ID
func WithUser(parent context.Context, userID string) context.Context {
	return WithValue(parent, UserKey, userID)
}

// User returns the User ID or an empty string
func User(ctx Context) string {
	return Value[string](ctx, UserKey)
}

// RequireUser returns the request User or an error if not set
func RequireUser(ctx Context) (string, error) {
	return RequireValue[string](ctx, UserKey, errors.Forbidden)
}

// Originator returns the request originator or an empty string
func Originator(ctx Context) string {
	return Value[string](ctx, OriginatorKey)
}

// RequireOriginator returns the request originator or an error if not set
func RequireOriginator(ctx Context) (string, error) {
	return RequireValue[string](ctx, OriginatorKey, errors.Invalid)
}

// WithOriginator returns a new context with the given request originator
func WithOriginator(parent context.Context, originator string) context.Context {
	return WithValue(parent, OriginatorKey, originator)
}

// URL returns the request URL or an empty string
func URL(ctx Context) string {
	return Value[string](ctx, URLKey)
}

// RequireURL returns the request URL or an error if not set
func RequireURL(ctx Context) (string, error) {
	return RequireValue[string](ctx, URLKey, errors.Invalid)
}

// WithURL returns a new context with the given request URL
func WithURL(parent context.Context, requestURL string) context.Context {
	return WithValue(parent, URLKey, requestURL)
}

// Timestamp returns the request timestamp or an empty time.Time
func Timestamp(ctx Context) time.Time {
	return Value[time.Time](ctx, TimestampKey)
}

// RequireTimestamp returns the request timestamp or an error if not set
func RequireTimestamp(ctx Context) (time.Time, error) {
	return RequireValue[time.Time](ctx, TimestampKey, errors.Invalid)
}

// WithTimestamp returns a new context with the given request timestamp
func WithTimestamp(parent context.Context, requestTimestamp time.Time) context.Context {
	return WithValue(parent, TimestampKey, requestTimestamp)
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

// WithTenant returns a new context with the given request tenant
func WithTenant(parent context.Context, requestTenant string) context.Context {
	return WithValue(parent, TenantKey, requestTenant)
}

// Operation returns the Operation or an empty string
func Operation(ctx Context) string {
	return Value[string](ctx, OperationKey)
}

// RequireOperation returns the Operation or an error if not set
func RequireOperation(ctx Context) (string, error) {
	return RequireValue[string](ctx, OperationKey, errors.Invalid)
}

// WithOperation returns a new context with the given Operation
func WithOperation(parent context.Context, operation string) context.Context {
	return WithValue(parent, OperationKey, operation)
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

func Clone(ctx context.Context) context.Context {
	newCtx := context.Background()

	for _, k := range AllKeys {
		v := ctx.Value(k)
		newCtx = context.WithValue(newCtx, k, v)
	}

	return newCtx
}
