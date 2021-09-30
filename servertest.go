// Copyright (c) 2020 Ketch Kloud, Inc.
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

package orlop

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/log"
	"google.golang.org/grpc"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// TestServer provides functionality for running a test server instance
type TestServer struct {
	*httptest.Server
}

// Connect opens a gRPC client connection to the server
func (s *TestServer) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	return Connect(ctx, s.ClientConfig(), VaultConfig{})
}

// ClientConfig returns a proper ClientConfig for connecting to the server
func (s *TestServer) ClientConfig() ClientConfig {
	return ClientConfig{
		URL: strings.TrimPrefix(s.URL, "https://"),
		TLS: FromTLSConfig(s.TLS),
	}
}

// GrpcTestFunc defines a function called for a GRPC test
type GrpcTestFunc func(ctx context.Context, t *testing.T, conn *grpc.ClientConn)

// RunGrpcTest runs a test function with a client GRPC connection connected to the given server
func RunGrpcTest(ctx context.Context, t *testing.T, s *TestServer, name string, fn GrpcTestFunc) {
	conn, err := s.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer conn.Close()

	t.Run(name, func(t *testing.T) {
		fn(ctx, t, conn)
	})
}

// RunGrpcTestSuite runs a suite of GRPC tests
func RunGrpcTestSuite(ctx context.Context, t *testing.T, serviceName string, options []ServerOption, testCases ...GrpcTestFunc) {
	s, err := NewTestServer(ctx, serviceName, options...)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer s.Close()

	for _, testCase := range testCases {
		RunGrpcTest(ctx, t, s, runtime.FuncForPC(reflect.ValueOf(testCase).Pointer()).Name(), testCase)
	}
}

// NewTestServer sets up the test server and
func NewTestServer(ctx context.Context, serviceName string, options ...ServerOption) (*TestServer, error) {
	var err error

	// Setup the server options
	serverOptions := &serverOptions{
		serviceName: serviceName,
		logger:      log.FromContext(ctx).WithField("service", serviceName),
	}

	options = append([]ServerOption{
		WithServerConfig(ServerConfig{
			Bind:   "0.0.0.0",
			Listen: 5000,
			TLS:    TLSConfig{},
		}),
	}, options...)

	// Process all server options (which may override any of the above)
	for _, option := range options {
		err = option.apply(ctx, serverOptions)
		if err != nil {
			return nil, errors.Wrap(err, "serve: failed to apply options")
		}
	}

	// Create the HTTP server
	mux := chi.NewMux()
	if serverOptions.notFound != nil {
		mux.NotFound(serverOptions.notFound.ServeHTTP)
	}
	if serverOptions.methodNotAllowed != nil {
		mux.MethodNotAllowed(serverOptions.methodNotAllowed.ServeHTTP)
	}

	// Add standard middleware
	mux.Use(DefaultHTTPHeaders(HeaderOptions{
		AllowedOrigins: serverOptions.config.AllowedOrigins,
	}))

	// Add any middlewares registered
	if len(serverOptions.middlewares) > 0 {
		mux.Use(serverOptions.middlewares...)
	}

	for _, option := range options {
		err = option.addHandler(ctx, serverOptions, mux)
		if err != nil {
			return nil, errors.Wrap(err, "serve: failed to add handler")
		}
	}

	// Start listening
	srv := httptest.NewUnstartedServer(mux)
	srv.EnableHTTP2 = true

	// Serve requests
	if serverOptions.config.TLS.GetEnabled() {
		config, err := NewServerTLSConfig(ctx, serverOptions.config.TLS, serverOptions.vault)
		if err != nil {
			return nil, err
		}

		srv.TLS = config
	}

	srv.StartTLS()
	return &TestServer{
		Server: srv,
	}, nil
}
