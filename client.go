package orlop

import (
	"fmt"
	"github.com/switch-bit/orlop/log"
	"google.golang.org/grpc"
)

// Connect creates a new client from configuration
func Connect(cfg HasClientConfig, vault HasVaultConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	l := log.WithField("url", cfg.GetURL())

	if len(cfg.GetURL()) == 0 {
		l.Errorf("client: url required")
		return nil, fmt.Errorf("client: url required")
	}

	if cfg.GetTLS().GetEnabled() {
		l.Trace("tls enabled")
		creds, err := NewClientTLSCredentials(cfg.GetTLS(), vault)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		l.Trace("tls disabled")
		opts = append(opts, grpc.WithInsecure())
	}

	shared := cfg.GetToken().GetShared()
	if len(shared.GetID()) > 0 || len(shared.GetFile()) > 0 || len(shared.GetSecret()) > 0 {
		l.Trace("loading token from configuration")

		s, err := LoadKey(shared, vault, "secret")
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithPerRPCCredentials(SharedContextCredentials{
			token: string(s),
		}))
	} else {
		l.Trace("using context credentials")

		opts = append(opts, grpc.WithPerRPCCredentials(ContextCredentials{}))
	}

	l.Trace("dialling")
	conn, err := grpc.Dial(cfg.GetURL(), opts...)
	if err != nil {
		l.WithError(err).Error("failed dialling")
		return nil, err
	}

	return conn, nil
}
