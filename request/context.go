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
	"github.com/google/uuid"
	"time"
)

// Key is the key of a request property in context
type Key string

var (
	ConnectionKey  Key = "connectionId"
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

var KeyMap = map[any]string{
	ConnectionKey:  "Connection",
	IDKey:          "Request-Id",
	OperationKey:   "Operation",
	OriginatorKey:  "Request-Originator",
	TenantKey:      "Tenant",
	TimestampKey:   "RequestTS",
	URLKey:         "Request-Url",
	UserKey:        "User-Id",
	IntegrationKey: "Integration",
}

var MetricsKeyMap = map[any]string{
	OperationKey:   "operation",
	OriginatorKey:  "requestOriginator",
	TenantKey:      "tenant",
	ConnectionKey:  "connection",
	IntegrationKey: "integration",
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

// Clone clones the given context into a blank context
func Clone(ctx context.Context) context.Context {
	newCtx := context.Background()

	for _, k := range AllKeys {
		newCtx = context.WithValue(newCtx, k, ctx.Value(k))
	}

	return newCtx
}

// Setup sets up a context with the given options
func Setup(ctx context.Context, opts ...func(ctx context.Context) context.Context) context.Context {
	// Check to see if we have a request ID and if not, populate it and some other of the request fields
	if len(ID(ctx)) == 0 {
		u, err := uuid.NewRandom()
		if err == nil {
			ctx = WithID(ctx, u.String())
		}
	}

	if Timestamp(ctx).IsZero() {
		ctx = WithTimestamp(ctx, time.Now().UTC())
	}

	for _, opt := range opts {
		ctx = opt(ctx)
	}

	return ctx
}
