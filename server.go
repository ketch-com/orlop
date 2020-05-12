// Copyright (c) 2020 SwitchBit, Inc.
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
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/switch-bit/orlop/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	syslog "log"
	"net"
	"net/http"
	"strings"
)

// Serve sets up the server and listens for requests
func Serve(ctx context.Context, serviceName string, options ...ServerOption) error {
	var err error

	// Setup the server options
	serverOptions := &ServerOptions{
		serviceName: serviceName,
		config:      ServerConfig{
			Bind:   "0.0.0.0",
			Listen: 5000,
			TLS:    TLSConfig{},
		},
		tlsProvider: NewSimpleTLSProvider(),
		handlers:    make(map[string]http.Handler),
	}

	// Add default health check
	err = WithHealthCheck(nil).apply(serverOptions)
	if err != nil {
		return err
	}

	// Add default metrics handler
	err = WithPrometheusMetrics().apply(serverOptions)
	if err != nil {
		return err
	}

	// Process all server options (which may override any of the above)
	for _, option := range options {
		err = option.apply(serverOptions)
		if err != nil {
			return err
		}
	}

	addr := fmt.Sprintf("%s:%d", serverOptions.config.GetBind(), serverOptions.config.GetListen())

	l := log.WithField("service", serviceName).WithField("addr", addr)

	// Create the HTTP server
	mux := http.NewServeMux()

	for key, handler := range serverOptions.handlers {
		mux.Handle(key, handler)
	}

	w := log.Writer()
	defer w.Close()

	srv := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: syslog.New(w, "[http]", 0),
	}

	// Setup the Gateway
	if serverOptions.registerServices != nil && len(serverOptions.gatewayHandlers) > 0 {
		gwmux := runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
			runtime.WithForwardResponseOption(redirectFilter),
			runtime.WithOutgoingHeaderMatcher(outgoingHeaderMatcher),
			runtime.WithMarshalerOption("application/octet-stream", &BinaryMarshaler{}),
			runtime.WithMarshalerOption("application/json", &runtime.JSONPb{
				EnumsAsInts:  true,
				EmitDefaults: false,
				OrigName:     false,
			}),
			runtime.WithMarshalerOption("application/javascript", &runtime.JSONPb{
				EnumsAsInts:  true,
				EmitDefaults: false,
				OrigName:     false,
			}),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
				Marshaler: &runtime.JSONPb{
					EnumsAsInts:  true,
					EmitDefaults: false,
					OrigName:     true,
				},
			}),
		)

		// Dial the server
		l.Trace("loading client credentials for loopback")
		t, err := serverOptions.tlsProvider.NewClientTLSConfig(serverOptions.config.GetTLS())
		if err != nil {
			return err
		}

		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(t)),
		}

		l.Trace("dialling grpc")
		conn, err := grpc.Dial(addr, dialOptions...)
		if err != nil {
			return err
		}

		l.Trace("registering gateway handlers")
		for _, gatewayHandler := range serverOptions.gatewayHandlers {
			err = gatewayHandler(ctx, gwmux, conn)
			if err != nil {
				return err
			}
		}

		// Add the JSON gateway
		gatewayHandler, err := NewInstrumentedMetricHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.TLS != nil {
				// Only on TLS per https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
				w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			}

			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
					w.Header().Set("Access-Control-Allow-Headers", headers)
					w.Header().Set("Access-Control-Allow-Methods", methods)
					return
				}
			}

			gwmux.ServeHTTP(w, r)
		})
		if err != nil {
			return err
		}

		mux.Handle(fmt.Sprintf("/%s/", serviceName), gatewayHandler)
	}

	// Setup the GRPC service
	if serverOptions.registerServices != nil {
		var grpcServerOptions []grpc.ServerOption

		// If certificate file and key file have been specified then setup a TLS server
		if serverOptions.config.GetTLS().GetEnabled() {
			l.Trace("tls enabled")

			t, err := serverOptions.tlsProvider.NewServerTLSConfig(serverOptions.config.GetTLS())
			if err != nil {
				return err
			}

			grpcServerOptions = append(grpcServerOptions, grpc.Creds(credentials.NewTLS(t)))
		} else {
			l.Trace("tls disabled")
		}

		// Intercept all request to provide authentication
		if serverOptions.authenticate != nil {
			grpcServerOptions = append(grpcServerOptions, grpc.UnaryInterceptor(serverOptions.authenticate))
		}

		// Setup the gRPC server
		grpcServer := grpc.NewServer(grpcServerOptions...)

		// Finally, add the GRPC handler at the root
		grpcHandler, err := NewInstrumentedMetricHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		})
		if err != nil {
			return err
		}

		// Register all the services
		l.Trace("registering GRPC services")
		err = serverOptions.registerServices(ctx, grpcServer)
		if err != nil {
			return err
		}

		mux.Handle("/", grpcHandler)
	}

	// Start listening
	l.Info("listening")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Serve requests
	if serverOptions.config.GetTLS().GetEnabled() {
		l.Trace("loading server tls certs")
		config, err := serverOptions.tlsProvider.NewServerTLSConfig(serverOptions.config.GetTLS())
		if err != nil {
			ln.Close()

			return err
		}

		ln = tls.NewListener(ln, config)
	}

	defer ln.Close()

	l.Info("serving")
	return srv.Serve(ln)
}

var (
	headers = strings.Join([]string{"Content-Type", "Accept", "Authorization"}, ",")
	methods = strings.Join([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}, ",")
)

func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func outgoingHeaderMatcher(key string) (string, bool) {
	if strings.HasPrefix(strings.ToLower(key), "access-control-") {
		return key, true
	}

	switch strings.ToLower(key) {
	case "cache-control", "expires", "etag", "x-content-type-options", "strict-transport-security", "vary":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func incomingHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-forwarded-for", "x-real-ip":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func redirectFilter(_ context.Context, w http.ResponseWriter, resp proto.Message) error {
	if t, ok := resp.(*Redirect); ok {
		w.Header().Set("Location", t.Location)
		w.WriteHeader(http.StatusFound)
	}

	return nil
}
