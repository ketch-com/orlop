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
	"crypto/tls"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/log"
	"go.uber.org/fx"
	syslog "log"
	"net"
	"net/http"
)

type ServeLifecycleParams struct {
	fx.In

	Lifecycle     fx.Lifecycle
	ServiceName   string `name:"serviceName"`
	ServerOptions []ServerOption
}

func ServeLifecycle(params ServeLifecycleParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := Serve(ctx, params.ServiceName, params.ServerOptions...)
				if err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
	})
}

// Serve sets up the server and listens for requests
func Serve(ctx context.Context, serviceName string, options ...ServerOption) error {
	var err error

	ctx, span := tracer.Start(ctx, "Serve")
	defer span.End()

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
		WithHealth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})),
		WithMetrics(NewMetricsHandler()),
	}, options...)

	if serverOptions.config.Profiling.GetEnabled() {
		options = append(options, WithProfiler())
	}

	// Process all server options (which may override any of the above)
	for _, option := range options {
		err = option.apply(ctx, serverOptions)
		if err != nil {
			err = errors.Wrap(err, "serve: failed to apply options")
			span.RecordError(err)
			return err
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
	mux.Use(Metrics)
	mux.Use(DefaultHTTPHeaders)
	if serverOptions.config.Logging.Enabled {
		mux.Use(Logging(serverOptions.config.Logging))
	}

	// Add any middlewares registered
	if len(serverOptions.middlewares) > 0 {
		mux.Use(serverOptions.middlewares...)
	}

	for _, option := range options {
		err = option.addHandler(ctx, serverOptions, mux)
		if err != nil {
			err = errors.Wrap(err, "serve: failed to add handler")
			span.RecordError(err)
			return err
		}
	}

	// Start listening
	serverOptions.logger.Info("listening")
	ln, err := net.Listen("tcp", serverOptions.addr)
	if err != nil {
		err = errors.Wrap(err, "serve: failed to listen")
		span.RecordError(err)
		return err
	}

	// Serve requests
	if serverOptions.config.GetTLS().GetEnabled() {
		config, err := NewServerTLSConfig(ctx, serverOptions.config.GetTLS(), serverOptions.vault)
		if err != nil {
			ln.Close()

			err = errors.Wrap(err, "serve: failed to get server TLS config")
			span.RecordError(err)
			return err
		}

		ln = tls.NewListener(ln, config)
	}

	defer ln.Close()

	serverOptions.logger.Info("serving")

	w := serverOptions.logger.WriterLevel(logrus.WarnLevel)
	defer w.Close()

	srv := &http.Server{
		Addr:     serverOptions.addr,
		Handler:  mux,
		ErrorLog: syslog.New(w, "[http]", 0),
	}

	if err = srv.Serve(ln); err != nil {
		err = errors.Wrap(err, "serve: failed to serve")
		span.RecordError(err)
		return err
	}

	return nil
}
