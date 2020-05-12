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
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"mime"
	"net/http"
)

type ServerOptions struct {
	serviceName      string
	config           ServerConfig
	handlers         map[string]http.Handler
	tlsProvider      TLSProvider
	registerServices func(ctx context.Context, grpcServer *grpc.Server) error
	authenticate     func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	gatewayHandlers  []func(ctx context.Context, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error
}

type ServerOption interface {
	apply(*ServerOptions) error
}



type ServerConfigOption struct {
	config HasServerConfig
}

func (o ServerConfigOption) apply(opt *ServerOptions) error {
	opt.config = ServerConfig{
		Bind:   o.config.GetBind(),
		Listen: o.config.GetListen(),
		TLS:    CloneTLSConfig(o.config.GetTLS()),
	}
	return nil
}

func WithServerConfig(config HasServerConfig) ServerOption {
	return &ServerConfigOption{
		config: config,
	}
}



type AuthenticateServerOption struct {
	authenticate func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

func (o AuthenticateServerOption) apply(opt *ServerOptions) error {
	opt.authenticate = o.authenticate
	return nil
}

func WithAuthentication(authenticate func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)) ServerOption {
	return &AuthenticateServerOption{
		authenticate: authenticate,
	}
}


type TLSServerOption struct {
	tlsConfig TLSConfig
}

func (o TLSServerOption) apply(opt *ServerOptions) error {
	opt.config.TLS = o.tlsConfig
	return nil
}

func WithTLS(cfg TLSConfig) ServerOption {
	return &TLSServerOption{
		tlsConfig: cfg,
	}
}



type GRPCServerServerOption struct {
	registerServices func(ctx context.Context, grpcServer *grpc.Server) error
}

func (o GRPCServerServerOption) apply(opt *ServerOptions) error {
	opt.registerServices = o.registerServices
	return nil
}

func WithGRPCServer(registerServices func(ctx context.Context, grpcServer *grpc.Server) error) ServerOption {
	return &GRPCServerServerOption{
		registerServices: registerServices,
	}
}


type GatewayServerOption struct {
	gatewayHandlers []func(ctx context.Context, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error
}

func (o GatewayServerOption) apply(opt *ServerOptions) error {
	opt.gatewayHandlers = o.gatewayHandlers
	return nil
}

func WithGateway(gatewayHandlers ...func(ctx context.Context, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error) ServerOption {
	return &GatewayServerOption{
		gatewayHandlers: gatewayHandlers,
	}
}


type TLSProviderServerOption struct {
	tlsProvider TLSProvider
}

func (o TLSProviderServerOption) apply(opt *ServerOptions) error {
	opt.tlsProvider = o.tlsProvider
	return nil
}

func WithTLSProvider(tlsProvider TLSProvider) ServerOption {
	return &TLSProviderServerOption{
		tlsProvider: tlsProvider,
	}
}

func WithHealthCheck(checker HealthChecker) ServerOption {
	return WithHandler("/healthz", &HealthHandler{
		checker: checker,
	})
}

func WithMetrics(handler http.Handler) ServerOption {
	return WithHandler("/metrics", handler)
}

func WithPrometheusMetrics() ServerOption {
	return WithMetrics(promhttp.Handler())
}



type SwaggerHandlerServerOption struct {
	fs http.FileSystem
}

func (o SwaggerHandlerServerOption) apply(opt *ServerOptions) error {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return err
	}

	handler := http.StripPrefix(fmt.Sprintf("/%s/swagger", opt.serviceName), http.FileServer(o.fs))
	opt.handlers[fmt.Sprintf("/%s/swagger/", opt.serviceName)] = handler

	return nil
}

func WithSwagger(fs http.FileSystem) ServerOption {
	return &SwaggerHandlerServerOption{}
}



type HandlerServerOption struct {
	pattern string
	handler http.Handler
}

func (o HandlerServerOption) apply(opt *ServerOptions) error {
	opt.handlers[o.pattern] = o.handler
	return nil
}

func WithHandler(pattern string, handler http.Handler) ServerOption {
	return &HandlerServerOption{
		pattern: pattern,
		handler: handler,
	}
}
