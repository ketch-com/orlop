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
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/logging"
	"go.ketch.com/lib/orlop/pem"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type ClientProvider interface {
	NewConfig(ctx context.Context, cfg Config) (*tls.Config, error)
}

type clientProvider struct {
	params ClientProviderParams
}

type ClientProviderParams struct {
	fx.In

	Logger    logging.Logger
	Tracer    trace.Tracer
	PEM       pem.Provider
	Generator Generator `optional:"true"`
}

func NewClientProvider(params ClientProviderParams) ClientProvider {
	return &clientProvider{
		params: params,
	}
}

func (p *clientProvider) NewConfig(ctx context.Context, cfg Config) (*tls.Config, error) {
	ctx, span := p.params.Tracer.Start(ctx, "NewClientTLSConfig")
	defer span.End()

	config := &tls.Config{
		ServerName:         cfg.Override,
		InsecureSkipVerify: cfg.Insecure,
		MinVersion:         tls.VersionTLS12,
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
			ID:   cfg.Cert.GetID(),
			File: cfg.Cert.GetFile(),
		})
		if err != nil {
			err = errors.Wrap(err, "tls: failed to load certificate")
			span.RecordError(err)
			return nil, err
		}

		p.params.Logger.Trace("certificate loaded")

		keyPEMBlock, err := p.params.PEM.Load(ctx, pem.KeyConfig{
			ID:   cfg.Key.GetID(),
			File: cfg.Key.GetFile(),
		})
		if err != nil {
			err = errors.Wrap(err, "tls: failed to load private key")
			span.RecordError(err)
			return nil, err
		}

		config.RootCAs = x509.NewCertPool()

		if !config.RootCAs.AppendCertsFromPEM(certPEMBlock) {
			err = errors.Wrap(err, "tls: failed to append to RootCA certificates")
			span.RecordError(err)
			return nil, err
		}

		if cfg.RootCA.GetEnabled() {
			rootcaPEMBlock, err := p.params.PEM.Load(ctx, pem.RootCAConfig{
				ID:   cfg.RootCA.GetID(),
				File: cfg.RootCA.GetFile(),
			})
			if err != nil {
				err = errors.Wrap(err, "tls: failed to load RootCA certificates")
				span.RecordError(err)
				return nil, err
			}

			if !config.RootCAs.AppendCertsFromPEM(rootcaPEMBlock) {
				err = errors.Wrap(err, "tls: failed to append to RootCA certificates")
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
