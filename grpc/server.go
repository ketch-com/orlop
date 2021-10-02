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
