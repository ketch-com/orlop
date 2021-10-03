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
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"mime"
	"net/http"
	"net/http/pprof"
	"strings"
)

// ServerOption provides an interface for utilizing custom server options
type ServerOption interface {
	apply(ctx context.Context, opts *serverOptions) error
	addHandler(ctx context.Context, opt *serverOptions, mux mux) error
}

type mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)

	Method(method, pattern string, handler http.Handler)
	MethodFunc(method, pattern string, handler http.HandlerFunc)
}

type serverOptions struct {
	logger           *logrus.Entry
	serviceName      string
	addr             string
	notFound         http.Handler
	methodNotAllowed http.Handler
	config           ServerConfig
	vault            VaultConfig
	middlewares      []func(http.Handler) http.Handler
	authenticate     func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// serverConfigOption provides the capability to override default server configuration including address, port and TLS
type serverConfigOption struct {
	config ServerConfig
}

func (o serverConfigOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.config = ServerConfig{
		Bind:           o.config.Bind,
		Listen:         o.config.Listen,
		Logging:        o.config.Logging,
		TLS:            CloneTLSConfig(o.config.TLS),
		AllowedOrigins: o.config.AllowedOrigins,
	}

	opt.addr = fmt.Sprintf("%s:%d", opt.config.Bind, opt.config.Listen)
	opt.logger = opt.logger.WithField("addr", opt.addr)

	return nil
}

func (o serverConfigOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	return nil
}

// WithServerConfig returns a new serverConfigOption
func WithServerConfig(config ServerConfig) ServerOption {
	return &serverConfigOption{
		config: config,
	}
}

// loggerServerOption provides capability to provide custom logger
type loggerServerOption struct {
	log *logrus.Entry
}

func (o loggerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.logger = o.log
	return nil
}

func (o loggerServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	return nil
}

// WithLogger returns a new loggerServerOption
func WithLogger(log *logrus.Entry) ServerOption {
	return &loggerServerOption{
		log: log,
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

func (o authenticateServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
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

func (o tlsServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	return nil
}

// WithTLS returns a new tlsServerOption
func WithTLS(cfg TLSConfig) ServerOption {
	return &tlsServerOption{
		tlsConfig: cfg,
	}
}

// grpcServicesServerOption is used to register a GRPC server
type grpcServicesServerOption struct {
	registerServices func(ctx context.Context, grpcServer *grpc.Server)
}

func (o grpcServicesServerOption) apply(ctx context.Context, opt *serverOptions) error {
	return nil
}

func (o grpcServicesServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	var grpcServerOptions []grpc.ServerOption

	// If certificate file and key file have been specified then setup a TLS server
	if opt.config.TLS.GetEnabled() {
		t, err := NewServerTLSConfig(ctx, opt.config.TLS, opt.vault)
		if err != nil {
			return errors.Wrap(err, "server: failed to load server TLS config")
		}

		grpcServerOptions = append(grpcServerOptions, grpc.Creds(credentials.NewTLS(t)))
	}

	grpcServerOptions = append(grpcServerOptions, grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	grpcServerOptions = append(grpcServerOptions, grpc.ChainStreamInterceptor(otelgrpc.StreamServerInterceptor()))

	// Intercept all request to provide authentication
	if opt.authenticate != nil {
		grpcServerOptions = append(grpcServerOptions, grpc.ChainUnaryInterceptor(opt.authenticate))
	}

	// Setup the gRPC server
	grpcServer := grpc.NewServer(grpcServerOptions...)

	// Register all the services
	opt.logger.Trace("registering GRPC services")
	o.registerServices(ctx, grpcServer)

	// Finally, add the GRPC handler at the root
	grpcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else if opt.notFound != nil {
			opt.notFound.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	mux.Handle(fmt.Sprintf("/%s.{service}/*", opt.serviceName), grpcHandler)
	return nil
}

// WithGRPCServices returns a new grpcServicesServerOption
func WithGRPCServices(registerServices func(ctx context.Context, grpcServer *grpc.Server)) ServerOption {
	return &grpcServicesServerOption{
		registerServices: registerServices,
	}
}

// vaultServerOption is used to specify Vault configuration
type vaultServerOption struct {
	vault VaultConfig
}

func (o vaultServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.vault = o.vault
	return nil
}

func (o vaultServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	return nil
}

// WithVault returns a new vaultServerOption
func WithVault(vault VaultConfig) ServerOption {
	return &vaultServerOption{
		vault: vault,
	}
}

// WithHealth specifies a health handler
func WithHealth(checker http.Handler) ServerOption {
	return WithHandler("/healthz", checker)
}

// WithHealthCheck specifies a health checker function
func WithHealthCheck(check string, checker http.Handler) ServerOption {
	return WithGET("/healthz/"+check, checker)
}

// WithMetrics specifies a metrics handler
func WithMetrics(handler http.Handler) ServerOption {
	return WithGET("/metrics", handler)
}

// profileServerOption specifies how to add profiler endpoints
type profileServerOption struct {
}

func (o profileServerOption) apply(ctx context.Context, opt *serverOptions) error {
	return nil
}

func (o profileServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	mux.MethodFunc(http.MethodGet, "/debug/pprof/", pprof.Index)
	mux.MethodFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
	mux.MethodFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
	mux.MethodFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
	mux.MethodFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)

	for _, handler := range []string{"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"} {
		mux.Method(http.MethodGet, fmt.Sprintf("/debug/pprof/%s", handler), pprof.Handler(handler))
	}

	return nil
}

// WithProfiler specifies a profiler handler to provide profiling information to go tool pprof
func WithProfiler() ServerOption {
	return &profileServerOption{}
}

// swaggerHandlerServerOption specifies how to serve swagger
type swaggerHandlerServerOption struct {
	fs http.FileSystem
}

func (o swaggerHandlerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	return nil
}

func (o swaggerHandlerServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return errors.Wrap(err, "failed to add svg extension")
	}

	handler := http.StripPrefix(fmt.Sprintf("/%s/swagger", opt.serviceName), http.FileServer(o.fs))
	mux.Method(http.MethodGet, fmt.Sprintf("/%s/swagger/", opt.serviceName), handler)

	return nil
}

// WithSwagger specifies a swagger handler based off the given file system
func WithSwagger(fs http.FileSystem) ServerOption {
	return &swaggerHandlerServerOption{fs: fs}
}

// handlerServerOption specifies a custom HTTP handler
type handlerServerOption struct {
	method      string
	pattern     string
	handler     http.Handler
	handlerFunc http.HandlerFunc
}

func (o handlerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	return nil
}

func (o handlerServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	if len(o.method) > 0 {
		if o.handler != nil {
			mux.Method(o.method, o.pattern, o.handler)
		} else {
			mux.MethodFunc(o.method, o.pattern, o.handlerFunc)
		}
	} else {
		if o.handler != nil {
			mux.Handle(o.pattern, o.handler)
		} else {
			mux.HandleFunc(o.pattern, o.handlerFunc)
		}
	}

	return nil
}

// WithHttpHandler returns a handlerServerOption
func WithHttpHandler(method string, pattern string, handler http.Handler) ServerOption {
	return &handlerServerOption{
		method:  method,
		pattern: pattern,
		handler: otelhttp.NewHandler(handler, pattern),
	}
}

// WithHandler returns a handlerServerOption
func WithHandler(pattern string, handler http.Handler) ServerOption {
	return WithHttpHandler("", pattern, handler)
}

// WithHandlerFunc returns a handlerServerOption
func WithHandlerFunc(pattern string, handler http.HandlerFunc) ServerOption {
	return WithHttpHandler("", pattern, handler)
}

// WithGET returns a handlerServerOption
func WithGET(pattern string, handler http.Handler) ServerOption {
	return WithHttpHandler(http.MethodGet, pattern, handler)
}

// WithPUT returns a handlerServerOption
func WithPUT(pattern string, handler http.Handler) ServerOption {
	return WithHttpHandler(http.MethodPut, pattern, handler)
}

// WithPOST returns a handlerServerOption
func WithPOST(pattern string, handler http.Handler) ServerOption {
	return WithHttpHandler(http.MethodPost, pattern, handler)
}

// WithDELETE returns a handlerServerOption
func WithDELETE(pattern string, handler http.Handler) ServerOption {
	return WithHttpHandler(http.MethodDelete, pattern, handler)
}

// WithOPTIONS returns a handlerServerOption
func WithOPTIONS(pattern string, handler http.Handler) ServerOption {
	return WithHttpHandler(http.MethodOptions, pattern, handler)
}

// WithPATCH returns a handlerServerOption
func WithPATCH(pattern string, handler http.Handler) ServerOption {
	return WithHttpHandler(http.MethodPatch, pattern, handler)
}

// notFoundHandlerServerOption specifies the handler to invoke when the route is not found
type notFoundHandlerServerOption struct {
	handler http.Handler
}

func (o notFoundHandlerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.notFound = o.handler
	return nil
}

func (o notFoundHandlerServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	return nil
}

// WithNotFoundHandler returns a notFoundHandlerServerOption
func WithNotFoundHandler(handler http.Handler) ServerOption {
	return &notFoundHandlerServerOption{
		handler: otelhttp.NewHandler(handler, "not found"),
	}
}

// methodNotAllowedHandlerServerOption specifies the handler to invoke when the route is not found
type methodNotAllowedHandlerServerOption struct {
	handler http.Handler
}

func (o methodNotAllowedHandlerServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.methodNotAllowed = o.handler
	return nil
}

func (o methodNotAllowedHandlerServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	return nil
}

// WithMethodNotAllowedHandler returns a notFoundHandlerServerOption
func WithMethodNotAllowedHandler(handler http.Handler) ServerOption {
	return &methodNotAllowedHandlerServerOption{
		handler: otelhttp.NewHandler(handler, "method not allowed"),
	}
}

// middlewareServerOption specifies a list of middlewares to add to the router
type middlewareServerOption struct {
	middlewares []func(http.Handler) http.Handler
}

func (o middlewareServerOption) apply(ctx context.Context, opt *serverOptions) error {
	opt.middlewares = append(opt.middlewares, o.middlewares...)
	return nil
}

func (o middlewareServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	return nil
}

// WithMiddleware returns a middlewareServerOption
func WithMiddleware(middlewares ...func(http.Handler) http.Handler) ServerOption {
	return &middlewareServerOption{
		middlewares: middlewares,
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
