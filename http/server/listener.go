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

package server

import (
	"context"
	crypto_tls "crypto/tls"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/tls"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"net"
)

type NewListenerParams struct {
	fx.In

	Config    Config
	Logger    logging.Logger
	Tracer    trace.Tracer
	ServerTLS tls.ServerProvider
}

func NewListener(ctx context.Context, params NewListenerParams) (net.Listener, error) {
	ctx, span := params.Tracer.Start(ctx, "NewListener")
	defer span.End()

	// Start listening
	params.Logger.Info("listening")

	listener, err := net.Listen("tcp", params.Config.Addr())
	if err != nil {
		err = errors.Wrapf(err, "failed to listen on %s", params.Config.Addr())
		span.RecordError(err)
		return nil, err
	}

	// If TLS is not enabled, return the listener as it is
	if !params.Config.TLS.GetEnabled() {
		return listener, nil
	}

	// Since TLS is enabled, get the tls.Config and return a crypto/tls listener
	config, err := params.ServerTLS.NewConfig(ctx, params.Config.TLS)
	if err != nil {
		_ = listener.Close()

		err = errors.Wrap(err, "failed to get server TLS config")
		span.RecordError(err)
		return nil, err
	}

	return crypto_tls.NewListener(listener, config), nil
}
