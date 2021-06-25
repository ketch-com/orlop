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
	"fmt"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/log"
	"go.ketch.com/lib/orlop/version"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"strings"
)

// Connect creates a new client from configuration
func Connect(ctx context.Context, cfg HasClientConfig, vault HasVaultConfig) (*grpc.ClientConn, error) {
	var cancel context.CancelFunc = func() {}

	ctx, span := tracer.Start(ctx, cfg.GetName())
	defer span.End()

	logger := log.FromContext(ctx)

	var opts []grpc.DialOption

	if len(cfg.GetURL()) == 0 {
		err := errors.Errorf("client: url required for %s", cfg.GetName())
		span.RecordError(err)
		return nil, err
	}

	opts = append(opts, grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	opts = append(opts, grpc.WithChainStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if cfg.GetTLS().GetEnabled() {
		t, err := NewClientTLSConfig(ctx, cfg.GetTLS(), vault)
		if err != nil {
			span.RecordError(err)
			return nil, errors.Wrap(err, "client: failed to get client TLS config")
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(t)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	shared := cfg.GetToken().GetShared()
	if len(shared.GetID()) > 0 || len(shared.GetFile()) > 0 || len(shared.GetSecret()) > 0 {
		opts = append(opts, grpc.WithPerRPCCredentials(SharedContextCredentials{
			tokenProvider: func(ctx context.Context) string {
				ctx, span := tracer.Start(ctx, "TokenProvider")
				defer span.End()

				s, err := LoadKey(ctx, shared, vault, "secret")
				if err != nil {
					span.RecordError(err)
					logger.WithError(err).Error("client: could not load secret key")
					return ""
				}

				return string(s)
			},
		}))
	} else {
		opts = append(opts, grpc.WithPerRPCCredentials(ContextCredentials{}))
	}

	if cfg.GetWriteBufferSize() > 0 {
		opts = append(opts, grpc.WithWriteBufferSize(cfg.GetWriteBufferSize()))
	}

	if cfg.GetReadBufferSize() > 0 {
		opts = append(opts, grpc.WithReadBufferSize(cfg.GetReadBufferSize()))
	}

	if cfg.GetInitialWindowSize() > 0 {
		opts = append(opts, grpc.WithInitialWindowSize(cfg.GetInitialWindowSize()))
	}

	if cfg.GetInitialConnWindowSize() > 0 {
		opts = append(opts, grpc.WithInitialConnWindowSize(cfg.GetInitialConnWindowSize()))
	}

	if cfg.GetMaxCallRecvMsgSize() > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(cfg.GetMaxCallRecvMsgSize())))
	}

	if cfg.GetMaxCallSendMsgSize() > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(cfg.GetMaxCallSendMsgSize())))
	}

	if cfg.GetMinConnectTimeout() > 0 {
		opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: cfg.GetMinConnectTimeout(),
		}))
	}

	if cfg.GetBlock() {
		opts = append(opts, grpc.WithBlock())
	}

	if cfg.GetConnTimeout() > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.GetConnTimeout())
		defer cancel()
	}

	ua := fmt.Sprintf("%s/%s", version.Name, version.Version)
	if len(cfg.GetUserAgent()) > 0 {
		ua = cfg.GetUserAgent()
	}
	opts = append(opts, grpc.WithUserAgent(ua))

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"name":                  cfg.GetName(),
		"url":                   cfg.GetURL(),
		"connTimeout":           cfg.GetConnTimeout(),
		"block":                 cfg.GetBlock(),
		"initialConnWindowSize": cfg.GetInitialConnWindowSize(),
		"initialWindowSize":     cfg.GetInitialWindowSize(),
		"maxCallRecvMsgSize":    cfg.GetMaxCallRecvMsgSize(),
		"maxCallSendMsgSize":    cfg.GetMaxCallSendMsgSize(),
		"minConnectTimeout":     cfg.GetMinConnectTimeout(),
		"readBufferSize":        cfg.GetReadBufferSize(),
		"userAgent":             ua,
		"writeBufferSize":       cfg.GetWriteBufferSize(),
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
