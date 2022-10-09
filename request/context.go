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
	"go.ketch.com/lib/orlop/v2/metadata"
	"time"
)

// Key is the key of a request property in context
type Key metadata.Key

// AllKeys is a slice of all Keys
var AllKeys = metadata.AllKeys

// HighCardinalityKeys is a map of high-cardinality keys
var HighCardinalityKeys = metadata.HighCardinalityKeys

// Setter is a function that adds a string to the context
type Setter metadata.Setter

// Getter is a function that returns a string from the context
type Getter metadata.Getter

// Setters is a map from the Key to a Setter for that Key
var Setters = metadata.Setters

// Getters is a map from the Key to a Getter for that Key
var Getters = metadata.Getters

// Value returns the value of a request context key, or a default value if the value is not set
func Value[T any](ctx context.Context, key Key) T {
	return metadata.Value[T](ctx, metadata.Key(key))
}

// RequireValue returns the value of a request context key or returns an error if the value is not set
func RequireValue[T any](ctx context.Context, key Key, errDecorator func(error) error) (T, error) {
	return metadata.RequireValue[T](ctx, metadata.Key(key), errDecorator)
}

// WithValue returns a new context with the given request value
func WithValue[T any](parent context.Context, key Key, v T) context.Context {
	return metadata.WithValue(parent, metadata.Key(key), v)
}

// ID returns the request ID or an empty string
func ID(ctx context.Context) string {
	return metadata.ID(ctx)
}

// RequireID returns the request ID or an error if not set
func RequireID(ctx context.Context) (string, error) {
	return metadata.RequireID(ctx)
}

// WithID returns a new context with the given request ID
func WithID(parent context.Context, requestID string) context.Context {
	return metadata.WithID(parent, requestID)
}

// Originator returns the request originator or an empty string
func Originator(ctx context.Context) string {
	return metadata.Originator(ctx)
}

// RequireOriginator returns the request originator or an error if not set
func RequireOriginator(ctx context.Context) (string, error) {
	return metadata.RequireOriginator(ctx)
}

// WithOriginator returns a new context with the given request originator
func WithOriginator(parent context.Context, originator string) context.Context {
	return metadata.WithOriginator(parent, originator)
}

// URL returns the request URL or an empty string
func URL(ctx context.Context) string {
	return metadata.URL(ctx)
}

// RequireURL returns the request URL or an error if not set
func RequireURL(ctx context.Context) (string, error) {
	return metadata.RequireURL(ctx)
}

// WithURL returns a new context with the given request URL
func WithURL(parent context.Context, requestURL string) context.Context {
	return metadata.WithURL(parent, requestURL)
}

// Timestamp returns the request timestamp or an empty time.Time
func Timestamp(ctx context.Context) time.Time {
	return metadata.Timestamp(ctx)
}

// RequireTimestamp returns the request timestamp or an error if not set
func RequireTimestamp(ctx context.Context) (time.Time, error) {
	return metadata.RequireTimestamp(ctx)
}

// WithTimestamp returns a new context with the given request timestamp
func WithTimestamp(parent context.Context, requestTimestamp time.Time) context.Context {
	return metadata.WithTimestamp(parent, requestTimestamp)
}

// Tenant returns the request Tenant or an empty string
func Tenant(ctx context.Context) string {
	return metadata.Tenant(ctx)
}

// RequireTenant returns the request Tenant or an error if not set
func RequireTenant(ctx context.Context) (string, error) {
	return metadata.RequireTenant(ctx)
}

// WithTenant returns a new context with the given request tenant
func WithTenant(parent context.Context, requestTenant string) context.Context {
	return metadata.WithTenant(parent, requestTenant)
}

// Operation returns the Operation or an empty string
func Operation(ctx context.Context) string {
	return metadata.Operation(ctx)
}

// RequireOperation returns the Operation or an error if not set
func RequireOperation(ctx context.Context) (string, error) {
	return metadata.RequireOperation(ctx)
}

// WithOperation returns a new context with the given Operation
func WithOperation(parent context.Context, operation string) context.Context {
	return metadata.WithOperation(parent, operation)
}

// Values returns a map of the request values from the context
func Values(ctx context.Context, opts ...Option) map[string]string {
	var mdOpts []metadata.Option
	for _, opt := range opts {
		mdOpts = append(mdOpts, metadata.Option(opt))
	}
	return metadata.Values(ctx, mdOpts...)
}

// Option is a function that sets values on the options structure
type Option metadata.Option

// Filter is a function that returns true if the given key should be included
type Filter metadata.Filter

// WithFilter returns a filter that filters request values
func WithFilter(f Filter) Option {
	return Option(metadata.WithFilter(metadata.Filter(f)))
}

// SkipHighCardinalityKeysFilter returns true if the key is not high cardinality
func SkipHighCardinalityKeysFilter(k Key) bool {
	return metadata.SkipHighCardinalityKeysFilter(metadata.Key(k))
}
