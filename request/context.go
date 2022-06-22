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
	"time"
)

// Key is the key of a request property in context
type Key string

var (
	IDKey        Key = "request_id"
	OperationKey Key = "operation"
	TimestampKey Key = "request_ts"
	TenantKey    Key = "tenant"
	URLKey       Key = "request_url"
)

// Value returns the value of a request context key
func Value[T any](ctx context.Context, key Key) T {
	if v := ctx.Value(key); v != nil {
		if r, ok := v.(T); ok {
			return r
		}
	}

	var v T
	return v
}

// WithValue returns a new context with the given request value
func WithValue[T any](parent context.Context, key Key, v T) context.Context {
	return context.WithValue(parent, key, v)
}

// ID returns the request ID or an empty string
func ID(ctx context.Context) string {
	return Value[string](ctx, IDKey)
}

// WithID returns a new context with the given request ID
func WithID(parent context.Context, requestID string) context.Context {
	return WithValue(parent, IDKey, requestID)
}

// URL returns the request URL or an empty string
func URL(ctx context.Context) string {
	return Value[string](ctx, URLKey)
}

// WithURL returns a new context with the given request URL
func WithURL(parent context.Context, requestURL string) context.Context {
	return WithValue(parent, URLKey, requestURL)
}

// Timestamp returns the request timestamp or an empty time.Time
func Timestamp(ctx context.Context) time.Time {
	return Value[time.Time](ctx, TimestampKey)
}

// WithTimestamp returns a new context with the given request timestamp
func WithTimestamp(parent context.Context, requestTimestamp time.Time) context.Context {
	return WithValue(parent, TimestampKey, requestTimestamp)
}

// Tenant returns the request Tenant or an empty string
func Tenant(ctx context.Context) string {
	return Value[string](ctx, TenantKey)
}

// WithTenant returns a new context with the given request tenant
func WithTenant(parent context.Context, requestTenant string) context.Context {
	return WithValue(parent, TenantKey, requestTenant)
}

// Operation returns the Operation or an empty string
func Operation(ctx context.Context) string {
	return Value[string](ctx, OperationKey)
}

// WithOperation returns a new context with the given Operation
func WithOperation(parent context.Context, operation string) context.Context {
	return WithValue(parent, OperationKey, operation)
}

// Values returns a map of the request values from the context
func Values(ctx context.Context) map[string]string {
	out := make(map[string]string)

	if requestID := ID(ctx); len(requestID) > 0 {
		out[string(IDKey)] = requestID
	}

	if requestURL := URL(ctx); len(requestURL) > 0 {
		out[string(URLKey)] = requestURL
	}

	if requestTimestamp := Timestamp(ctx); !requestTimestamp.IsZero() {
		out[string(TimestampKey)] = requestTimestamp.String()
	}

	if requestTenant := Tenant(ctx); len(requestTenant) > 0 {
		out[string(TenantKey)] = requestTenant
	}

	if operation := Tenant(ctx); len(operation) > 0 {
		out[string(OperationKey)] = operation
	}

	return out
}
