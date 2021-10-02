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

package grpc

import (
	"context"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/http/server"
	"go.ketch.com/lib/orlop/tls"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type NewServerParams struct {
	fx.In

	Config        server.Config
	TLSProvider   tls.ServerProvider
	ServerOptions []grpc.ServerOption `group:"grpcServerOptions" optional:"true"`
}

func NewServer(ctx context.Context, params NewServerParams) (*grpc.Server, error) {
	grpcServerOptions := params.ServerOptions

	// If certificate file and key file have been specified then set up a TLS server
	if params.Config.TLS.GetEnabled() {
		t, err := params.TLSProvider.NewConfig(ctx, params.Config.TLS)
		if err != nil {
			return nil, errors.Wrap(err, "server: failed to load server TLS config")
		}

		grpcServerOptions = append(grpcServerOptions, grpc.Creds(credentials.NewTLS(t)))
	}

	grpcServerOptions = append(grpcServerOptions, grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	grpcServerOptions = append(grpcServerOptions, grpc.ChainStreamInterceptor(otelgrpc.StreamServerInterceptor()))

	return grpc.NewServer(grpcServerOptions...), nil
}
