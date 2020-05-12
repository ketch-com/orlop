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
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"mime"
	"net/http"
	"strings"
)

// ServerOption provides an interface for utilizing custom server options
type ServerOption interface {
	apply(ctx context.Context, opts *serverOptions) error
}

// serverOptions contain
type serverOptions struct {
	log              *logrus.Entry
	serviceName      string
	addr             string
	config           ServerConfig
	handlers         map[string]http.Handler
	tlsProvider      TLSProvider
	authenticate     func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// serverConfigOption provides the capability to override default server configuration including address, port and TLS
type serverConfigOption struct {
	config HasServerConfig
}

func (o serverConfigOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.config = ServerConfig{
		Bind:   o.config.GetBind(),
		Listen: o.config.GetListen(),
		TLS:    CloneTLSConfig(o.config.GetTLS()),
	}

	opt.addr = fmt.Sprintf("%s:%d", opt.config.GetBind(), opt.config.GetListen())
	opt.log = opt.log.WithField("addr", opt.addr)

	return nil
}

// WithServerConfig returns a new serverConfigOption
func WithServerConfig(config HasServerConfig) ServerOption {
	return &serverConfigOption{
		config: config,
	}
}

// authenticateServerOption is used to specify an authenticator function
type authenticateServerOption struct {
	authenticate func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

func (o authenticateServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.authenticate = o.authenticate
	return nil
}

// WithAuthentication returns a new authenticateServerOption
func WithAuthentication(authenticate func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)) ServerOption {
	return &authenticateServerOption{
		authenticate: authenticate,
	}
}

// tlsServerOption is used to specify TLS config
type tlsServerOption struct {
	tlsConfig TLSConfig
}

func (o tlsServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.config.TLS = o.tlsConfig
	return nil
}

// WithTLS returns a new tlsServerOption
func WithTLS(cfg TLSConfig) ServerOption {
	return &tlsServerOption{
		tlsConfig: cfg,
	}
}

// grpcServerServerOption is used to register a GRPC server
type grpcServerServerOption struct {
	registerServices func(ctx context.Context, grpcServer *grpc.Server) error
}

func (o grpcServerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	var grpcServerOptions []grpc.ServerOption

	// If certificate file and key file have been specified then setup a TLS server
	if opt.config.GetTLS().GetEnabled() {
		opt.log.Trace("tls enabled")

		t, err := opt.tlsProvider.NewServerTLSConfig(opt.config.GetTLS())
		if err != nil {
			return err
		}

		grpcServerOptions = append(grpcServerOptions, grpc.Creds(credentials.NewTLS(t)))
	} else {
		opt.log.Trace("tls disabled")
	}

	// Intercept all request to provide authentication
	if opt.authenticate != nil {
		grpcServerOptions = append(grpcServerOptions, grpc.UnaryInterceptor(opt.authenticate))
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
	opt.log.Trace("registering GRPC services")
	err = o.registerServices(ctx, grpcServer)
	if err != nil {
		return err
	}

	opt.handlers["/"] = grpcHandler
	return nil
}

// WithGRPCServer returns a new grpcServerServerOption
func WithGRPCServer(registerServices func(ctx context.Context, grpcServer *grpc.Server) error) ServerOption {
	return &grpcServerServerOption{
		registerServices: registerServices,
	}
}

// gatewayServerOption is used to specify handlers for a JSON-GRPC gateway
type gatewayServerOption struct {
	gatewayHandlers []func(ctx context.Context, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error
}

func (o gatewayServerOption) apply(ctx context.Context, opt *serverOptions) error {
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
	opt.log.Trace("loading client credentials for loopback")
	t, err := opt.tlsProvider.NewClientTLSConfig(opt.config.GetTLS())
	if err != nil {
		return err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(t)),
	}

	opt.log.Trace("dialling grpc")
	conn, err := grpc.Dial(opt.addr, dialOptions...)
	if err != nil {
		return err
	}

	opt.log.Trace("registering gateway handlers")
	for _, gatewayHandler := range o.gatewayHandlers {
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

	opt.handlers[fmt.Sprintf("/%s/", opt.serviceName)] = gatewayHandler

	return nil
}

// WithGateway returns a new gatewayServerOption
func WithGateway(gatewayHandlers ...func(ctx context.Context, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error) ServerOption {
	return &gatewayServerOption{
		gatewayHandlers: gatewayHandlers,
	}
}

// tlsProviderServerOption is used to specify a TLS provider
type tlsProviderServerOption struct {
	tlsProvider TLSProvider
}

func (o tlsProviderServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.tlsProvider = o.tlsProvider
	return nil
}

// WithTLSProvider returns a new tlsProviderServerOption
func WithTLSProvider(tlsProvider TLSProvider) ServerOption {
	return &tlsProviderServerOption{
		tlsProvider: tlsProvider,
	}
}

// WithHealthCheck specifies a health checker function
func WithHealthCheck(checker HealthChecker) ServerOption {
	return WithHandler("/healthz", &HealthHandler{
		checker: checker,
	})
}

// WithMetrics specifies a metrics handler
func WithMetrics(handler http.Handler) ServerOption {
	return WithHandler("/metrics", handler)
}

// WithPrometheusMetrics specifies to use the Prometheus metrics handler
func WithPrometheusMetrics() ServerOption {
	return WithMetrics(promhttp.Handler())
}

// swaggerHandlerServerOption specifies how to serve swagger
type swaggerHandlerServerOption struct {
	fs http.FileSystem
}

func (o swaggerHandlerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return err
	}

	handler := http.StripPrefix(fmt.Sprintf("/%s/swagger", opt.serviceName), http.FileServer(o.fs))
	opt.handlers[fmt.Sprintf("/%s/swagger/", opt.serviceName)] = handler

	return nil
}

// WithSwagger specifies a swagger handler based off the given file system
func WithSwagger(fs http.FileSystem) ServerOption {
	return &swaggerHandlerServerOption{}
}

// handlerServerOption specifies a custom HTTP handler
type handlerServerOption struct {
	pattern string
	handler http.Handler
}

func (o handlerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.handlers[o.pattern] = o.handler
	return nil
}

// WithHandler returns a handlerServerOption
func WithHandler(pattern string, handler http.Handler) ServerOption {
	return &handlerServerOption{
		pattern: pattern,
		handler: handler,
	}
}

var (
	headers = strings.Join([]string{"Content-Type", "Accept", "Authorization"}, ",")
	methods = strings.Join([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}, ",")
)

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
