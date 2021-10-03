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
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/parameter"
	"go.ketch.com/lib/orlop/v2/tls"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type GeneratorParams struct {
	fx.In

	Tracer  trace.Tracer
	Config  GeneratorConfig
	Secrets parameter.Store
}

func NewGenerator(params GeneratorParams) tls.Generator {
	if params.Config.GetEnabled() {
		return &generatorImpl{
			params: params,
		}
	} else {
		return tls.NewNoopGenerator()
	}
}

type generatorImpl struct {
	params GeneratorParams
}

// GenerateCertificates calls Vault to generate a certificate
func (g generatorImpl) GenerateCertificates(ctx context.Context) (ok bool, cert []byte, key []byte, err error) {
	ctx, span := g.params.Tracer.Start(ctx, "GenerateCertificates")
	defer span.End()

	params := map[string]interface{}{
		"common_name": g.params.Config.CommonName,
		"format":      "pem_bundle",
	}
	if len(g.params.Config.AltNames) > 0 {
		params["alt_names"] = g.params.Config.AltNames
	}
	if g.params.Config.TTL.Seconds() > 60 {
		params["ttl"] = g.params.Config.TTL.String()
	}

	// Write the params to the path to generate the certificate
	secret, err := g.params.Secrets.Write(ctx, g.params.Config.Path, params)
	if err != nil {
		err = errors.Wrap(err, "generate: failed to write to Vault")
		span.RecordError(err)
		return false, nil, nil, err
	}

	if secret == nil {
		return false, nil, nil, nil
	}

	// Set the generated certificate and private key
	if d, ok := secret["certificate"]; ok {
		if s, ok := d.(string); ok {
			cert = []byte(s)
		}
	}

	if d, ok := secret["private_key"]; ok {
		if s, ok := d.(string); ok {
			key = []byte(s)
		}
	}

	return true, cert, key, nil
}
