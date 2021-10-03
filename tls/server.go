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

package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/pem"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type ServerProvider interface {
	NewConfig(ctx context.Context, cfg Config) (*tls.Config, error)
}

type serverProvider struct {
	params ServerProviderParams
}

type ServerProviderParams struct {
	fx.In

	Logger    logging.Logger
	Tracer    trace.Tracer
	PEM       pem.Provider
	Generator Generator `optional:"true"`
}

func NewServerProvider(params ServerProviderParams) ServerProvider {
	return &serverProvider{
		params: params,
	}
}

func (p *serverProvider) NewConfig(ctx context.Context, cfg Config) (*tls.Config, error) {
	ctx, span := p.params.Tracer.Start(ctx, "NewServerTLSConfig")
	defer span.End()

	config := &tls.Config{
		ClientAuth: cfg.ClientAuth,
		MinVersion: tls.VersionTLS12,
	}

	if !strSliceContains(config.NextProtos, "http/1.1") {
		// Enable HTTP/1.1
		config.NextProtos = append(config.NextProtos, "http/1.1")
	}

	if !strSliceContains(config.NextProtos, "h2") {
		// Enable HTTP/2
		config.NextProtos = append([]string{"h2"}, config.NextProtos...)
	}

	if !cfg.GetEnabled() {
		p.params.Logger.Trace("tls disabled")
		return config, nil
	}

	if p.params.Generator != nil {
		ok, certPEMBlock, keyPEMBlock, err := p.params.Generator.GenerateCertificates(ctx)
		if err != nil {
			err = errors.Wrap(err, "tls: failed to generate certificates")
			span.RecordError(err)
			return nil, err
		}

		if ok {
			c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
			if err != nil {
				err = errors.Wrap(err, "tls: failed creating key pair")
				span.RecordError(err)
				return nil, err
			}

			config.Certificates = append(config.Certificates, c)
			return config, nil
		}
	}

	if cfg.Cert.GetEnabled() && cfg.Key.GetEnabled() {
		certPEMBlock, err := p.params.PEM.Load(ctx, pem.CertificateConfig{
			ID:   cfg.Cert.ID,
			File: cfg.Cert.File,
		})
		if err != nil {
			err = errors.Wrap(err, "tls: failed to load certificate")
			span.RecordError(err)
			return nil, err
		}

		p.params.Logger.Trace("certificate loaded")

		keyPEMBlock, err := p.params.PEM.Load(ctx, pem.KeyConfig{
			ID:   cfg.Key.ID,
			File: cfg.Key.File,
		})
		if err != nil {
			err = errors.Wrap(err, "tls: failed to load private key")
			span.RecordError(err)
			return nil, err
		}

		config.ClientCAs = x509.NewCertPool()

		if cfg.RootCA.GetEnabled() {
			rootcaPEMBlock, err := p.params.PEM.Load(ctx, pem.RootCAConfig{
				ID:   cfg.RootCA.ID,
				File: cfg.RootCA.File,
			})
			if err != nil {
				err = errors.Wrap(err, "tls: failed to load RootCA certificates")
				span.RecordError(err)
				return nil, err
			}

			if !config.ClientCAs.AppendCertsFromPEM(rootcaPEMBlock) {
				err = errors.Wrap(err, "tls: failed to append RootCA certificates")
				span.RecordError(err)
				return nil, err
			}
		}

		c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			err = errors.Wrap(err, "tls: failed creating key pair")
			span.RecordError(err)
			return nil, err
		}

		config.Certificates = append(config.Certificates, c)
	}

	return config, nil
}
