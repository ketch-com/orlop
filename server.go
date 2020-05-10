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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/switch-bit/orlop/log"
	"google.golang.org/grpc"
	syslog "log"
	"mime"
	"net"
	"net/http"
	"strings"
)

// HasConfig denotes that an object provides server configuration
type HasConfig interface {
	GetServer() HasServerConfig
	GetVault() HasVaultConfig
}

// RegistersServices is an interface that is implemented by each server
type RegistersServices interface {
	GetConfig() HasConfig
	RegisterServices(ctx context.Context, grpcServer *grpc.Server, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error
}

// HasAuthenticate denotes that an object provides authentication
type HasAuthenticate interface {
	Authenticate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// Serve sets up the server and listens for requests
func Serve(ctx context.Context, serviceName string, server RegistersServices, swagger http.FileSystem) error {
	cfg := server.GetConfig()
	tlsCfg := cfg.GetServer().GetTLS()

	l := log.WithField("service", serviceName)

	var serverOptions []grpc.ServerOption

	// If certificate file and key file have been specified then setup a TLS server
	if tlsCfg.GetEnabled() {
		l.Trace("tls enabled")

		creds, err := NewServerTLSCredentials(cfg.GetServer().GetTLS(), cfg.GetVault())
		if err != nil {
			return err
		}

		serverOptions = append(serverOptions, grpc.Creds(creds))
	} else {
		l.Trace("tls disabled")
	}

	// Intercept all request to provide authentication
	if a, ok := server.(HasAuthenticate); ok {
		serverOptions = append(serverOptions, grpc.UnaryInterceptor(a.Authenticate))
	}

	// Setup the gRPC server
	grpcServer := grpc.NewServer(serverOptions...)

	// Listen on the configured port
	addr := fmt.Sprintf("%s:%d", cfg.GetServer().GetBind(), cfg.GetServer().GetListen())

	// Setup the gateway mux
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
	creds, err := NewClientTLSCredentials(tlsCfg, cfg.GetVault())
	if err != nil {
		return err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	l.Trace("dialling grpc")
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return err
	}

	// Create the HTTP server
	mux := http.NewServeMux()

	// Add the Health check endpoint
	var checker HealthChecker
	if c, ok := server.(HealthChecker); ok {
		checker = c
	}

	mux.Handle("/healthz", &HealthHandler{
		checker: checker,
	})

	// Add the Metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// If swagger is enabled, add the swagger endpoint
	if swagger != nil && cfg.GetServer().GetSwagger().GetEnabled() {
		l.Trace("swagger enabled")

		err = mime.AddExtensionType(".svg", "image/svg+xml")
		if err != nil {
			return err
		}

		swmux := http.StripPrefix(fmt.Sprintf("/%s/swagger", serviceName), http.FileServer(swagger))

		mux.Handle(fmt.Sprintf("/%s/swagger/", serviceName), swmux)
	} else {
		l.Trace("swagger disabled")
	}

	// Add the JSON gateway
	jsonGateway, err := NewInstrumentedMetricHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	mux.Handle(fmt.Sprintf("/%s/", serviceName), jsonGateway)

	// Finally, add the GRPC handler at the root
	grpcHandler, err := NewInstrumentedMetricHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else if h, ok := server.(http.Handler); ok {
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	if err != nil {
		return err
	}

	mux.Handle("/", grpcHandler)

	w := log.Writer()
	defer w.Close()

	srv := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: syslog.New(w, "[http]", 0),
	}

	// Register all the services
	l.Trace("registering services")
	err = server.RegisterServices(ctx, grpcServer, gwmux, conn)
	if err != nil {
		return err
	}

	l.WithField("addr", addr).Info("listening")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	// Server requests
	if tlsCfg.GetEnabled() {
		l.Trace("loading server tls certs")
		config, err := NewServerTLSConfig(cfg.GetServer().GetTLS(), cfg.GetVault())
		if err != nil {
			return err
		}

		l.Info("serving")
		return srv.Serve(tls.NewListener(ln, config))
	}

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
