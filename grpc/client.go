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

package grpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/logging"
	"go.ketch.com/lib/orlop/secret"
	"go.ketch.com/lib/orlop/tls"
	"go.ketch.com/lib/orlop/version"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"strings"
)

type Provider interface {
	Connect(ctx context.Context, cfg Config) (*grpc.ClientConn, error)
}

type Params struct {
	fx.In

	TLS         tls.ClientProvider
	Secret      secret.Provider
	Logger      logging.Logger
	Tracer      trace.Tracer
	DialOptions []grpc.DialOption `group:"grpcDialOptions" optional:"true"`
}

type providerImpl struct {
	params Params
}

func New(params Params) Provider {
	return &providerImpl{
		params: params,
	}
}

// Connect creates a new client from configuration
func (c providerImpl) Connect(ctx context.Context, cfg Config) (*grpc.ClientConn, error) {
	var cancel context.CancelFunc = func() {}

	ctx, span := c.params.Tracer.Start(ctx, cfg.GetName())
	defer span.End()

	opts := c.params.DialOptions

	if len(cfg.GetURL()) == 0 {
		err := errors.Errorf("client: url required for %s", cfg.GetName())
		span.RecordError(err)
		return nil, err
	}

	opts = append(opts, grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	opts = append(opts, grpc.WithChainStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if cfg.TLS.GetEnabled() {
		t, err := c.params.TLS.NewConfig(ctx, cfg.TLS)
		if err != nil {
			span.RecordError(err)
			return nil, errors.Wrap(err, "client: failed to get client TLS config")
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(t)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	shared := cfg.Token
	if len(shared.ID) > 0 || len(shared.Secret) > 0 {
		opts = append(opts, grpc.WithPerRPCCredentials(SharedContextCredentials{
			tokenProvider: func(ctx context.Context) string {
				ctx, span := c.params.Tracer.Start(ctx, "TokenProvider")
				defer span.End()

				s, err := c.params.Secret.Load(ctx, shared)
				if err != nil {
					span.RecordError(err)
					c.params.Logger.WithError(err).Error("client: could not load secret key")
					return ""
				}

				return s
			},
		}))
	} else {
		opts = append(opts, grpc.WithPerRPCCredentials(ContextCredentials{}))
	}

	if cfg.WriteBufferSize > 0 {
		opts = append(opts, grpc.WithWriteBufferSize(cfg.WriteBufferSize))
	}

	if cfg.ReadBufferSize > 0 {
		opts = append(opts, grpc.WithReadBufferSize(cfg.ReadBufferSize))
	}

	if cfg.InitialWindowSize > 0 {
		opts = append(opts, grpc.WithInitialWindowSize(cfg.InitialWindowSize))
	}

	if cfg.InitialConnWindowSize > 0 {
		opts = append(opts, grpc.WithInitialConnWindowSize(cfg.InitialConnWindowSize))
	}

	if cfg.MaxCallRecvMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(cfg.MaxCallRecvMsgSize)))
	}

	if cfg.MaxCallSendMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(cfg.MaxCallSendMsgSize)))
	}

	if cfg.MinConnectTimeout > 0 {
		opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: cfg.MinConnectTimeout,
		}))
	}

	if cfg.Block {
		opts = append(opts, grpc.WithBlock())
	}

	if cfg.ConnTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.ConnTimeout)
		defer cancel()
	}

	ua := fmt.Sprintf("%s/%s", version.Name, version.Version)
	if len(cfg.GetUserAgent()) > 0 {
		ua = cfg.GetUserAgent()
	}
	opts = append(opts, grpc.WithUserAgent(ua))

	c.params.Logger.WithFields(logrus.Fields{
		"name":                  cfg.GetName(),
		"url":                   cfg.GetURL(),
		"connTimeout":           cfg.ConnTimeout,
		"block":                 cfg.Block,
		"initialConnWindowSize": cfg.InitialConnWindowSize,
		"initialWindowSize":     cfg.InitialWindowSize,
		"maxCallRecvMsgSize":    cfg.MaxCallRecvMsgSize,
		"maxCallSendMsgSize":    cfg.MaxCallSendMsgSize,
		"minConnectTimeout":     cfg.MinConnectTimeout,
		"readBufferSize":        cfg.ReadBufferSize,
		"userAgent":             ua,
		"writeBufferSize":       cfg.WriteBufferSize,
	}).Trace("dialling")

	u := cfg.GetURL()
	for _, scheme := range []string{"https://", "http://"} {
		u = strings.TrimPrefix(u, scheme)
	}

	conn, err := grpc.DialContext(ctx, u, opts...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return conn, nil
}
