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

package vault

import (
	"context"
	vault "github.com/hashicorp/vault/api"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/parameter"
	"go.ketch.com/lib/orlop/v2/tls"
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
