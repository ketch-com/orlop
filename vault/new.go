package vault

import (
	"context"
	vault "github.com/hashicorp/vault/api"
	"go.ketch.com/lib/orlop/logging"
	"go.ketch.com/lib/orlop/parameter"
	"go.ketch.com/lib/orlop/tls"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"net/http"
)

type Params struct {
	fx.In

	Config   Config
	Logger   logging.Logger
	Tracer   trace.Tracer
	Provider tls.ClientProvider
}

// New connects to Vault given the configuration
func New(ctx context.Context, params Params) (parameter.ObjectStore, error) {
	var err error

	// First check if Vault is enabled in config, returning if not
	if !params.Config.GetEnabled() {
		params.Logger.Trace("vault is not enabled")
		return parameter.NewNoopStore(), nil
	}

	// Setup the Vault native config
	vc := &vault.Config{
		Address: params.Config.Address,
	}

	// If TLS is enabled, then setup the TLS configuration
	if params.Config.TLS.GetEnabled() {
		vc.HttpClient = &http.Client{}

		t := http.DefaultTransport.(*http.Transport).Clone()

		t.TLSClientConfig, err = params.Provider.NewConfig(ctx, params.Config.TLS)
		if err != nil {
			return nil, err
		}

		vc.HttpClient.Transport = t
	}

	// Create the vault native client
	c, err := vault.NewClient(vc)
	if err != nil {
		return nil, err
	}

	// Set the token on the client
	c.SetToken(params.Config.Token)

	return &client{
		params: params,
		client: c,
	}, nil
}
