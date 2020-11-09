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
	"crypto/x509"
	"github.com/switch-bit/orlop/errors"
	"github.com/switch-bit/orlop/log"
)

const (
	TLSCertificateKey = "certificate"
	TLSPrivateKey = "private_key"
	TLSRootCAKey = "issuing_ca"
)

// NewServerTLSConfig returns a new tls.VaultConfig from the given configuration input
func NewServerTLSConfig(ctx context.Context, cfg HasTLSConfig, vault HasVaultConfig) (*tls.Config, error) {
	var err error

	ctx, span := tracer.Start(ctx, "NewServerTLSConfig")
	defer span.End()

	config := &tls.Config{
		ClientAuth: cfg.GetClientAuth(),
		MinVersion: tls.VersionTLS12,
	}

	if !strSliceContains(config.NextProtos, "http/1.1") {
		// Enable HTTP/1.1
		config.NextProtos = append(config.NextProtos, "http/1.1")
	}

	if !cfg.GetEnabled() {
		log.Trace("tls disabled")
		return config, nil
	}

	if !strSliceContains(config.NextProtos, "h2") {
		// Enable HTTP/2
		config.NextProtos = append([]string{"h2"}, config.NextProtos...)
	}

	t := CloneTLSConfig(cfg)

	err = GenerateCertificates(ctx, vault, cfg.GetGenerate(), &t.Cert.Secret, &t.Key.Secret)
	if err != nil {
		err = errors.Wrap(err, "tls: failed to generate certificates")
		span.RecordError(ctx, err)
		return nil, err
	}

	certPEMBlock, err := LoadKey(ctx, t.GetCert(), vault, TLSCertificateKey)
	if err != nil {
		err = errors.Wrap(err, "tls: failed to load certificate")
		span.RecordError(ctx, err)
		return nil, err
	}

	log.Trace("certificate loaded")

	keyPEMBlock, err := LoadKey(ctx, t.GetKey(), vault, TLSPrivateKey)
	if err != nil {
		err = errors.Wrap(err, "tls: failed to load private key")
		span.RecordError(ctx, err)
		return nil, err
	}

	if t.GetRootCA().GetEnabled() {
		rootcaPEMBlock, err := LoadKey(ctx, t.GetRootCA(), vault, TLSRootCAKey)
		if err == nil {
			config.ClientCAs = x509.NewCertPool()

			if !config.ClientCAs.AppendCertsFromPEM(rootcaPEMBlock) {
				err = errors.Wrap(err, "tls: failed to append root CA certificates")
				span.RecordError(ctx, err)
				return nil, err
			}
		}
	}

	c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		err = errors.Wrap(err, "tls: failed creating key pair")
		span.RecordError(ctx, err)
		return nil, err
	}

	config.Certificates = append(config.Certificates, c)

	return config, nil
}

// NewClientTLSConfig returns a new tls.VaultConfig from the given configuration input
func NewClientTLSConfig(ctx context.Context, cfg HasTLSConfig, vault HasVaultConfig) (*tls.Config, error) {
	ctx, span := tracer.Start(ctx, "NewClientTLSConfig")
	defer span.End()

	if !cfg.GetEnabled() {
		log.Trace("tls disabled")
		return &tls.Config{}, nil
	}

	config := &tls.Config{
		ServerName:         cfg.GetOverride(),
		InsecureSkipVerify: cfg.GetInsecure(),
		MinVersion:         tls.VersionTLS12,
	}

	var err error
	t := CloneTLSConfig(cfg)

	err = GenerateCertificates(ctx, vault, cfg.GetGenerate(), &t.Cert.Secret, &t.Key.Secret)
	if err != nil {
		err = errors.Wrap(err, "tls: failed to generate certificates")
		span.RecordError(ctx, err)
		return nil, err
	}

	if t.GetCert().GetEnabled() && t.GetKey().GetEnabled() {
		certPEMBlock, err := LoadKey(ctx, t.GetCert(), vault, TLSCertificateKey)
		if err != nil {
			err = errors.Wrap(err, "tls: failed to load certificate")
			span.RecordError(ctx, err)
			return nil, err
		}

		keyPEMBlock, err := LoadKey(ctx, t.GetKey(), vault, TLSPrivateKey)
		if err != nil {
			err = errors.Wrap(err, "tls: failed to load private key")
			span.RecordError(ctx, err)
			return nil, err
		}

		config.RootCAs = x509.NewCertPool()

		if !config.RootCAs.AppendCertsFromPEM(certPEMBlock) {
			err = errors.Wrap(err, "tls: failed to append to RootCA certificates")
			span.RecordError(ctx, err)
			return nil, err
		}

		if t.GetRootCA().GetEnabled() {
			rootcaPEMBlock, err := LoadKey(ctx, t.GetRootCA(), vault, TLSRootCAKey)
			if err != nil {
				err = errors.Wrap(err, "tls: failed to load RootCA certificates")
				span.RecordError(ctx, err)
				return nil, err
			}

			if !config.RootCAs.AppendCertsFromPEM(rootcaPEMBlock) {
				err = errors.Wrap(err, "tls: failed to append to RootCA certificates")
				span.RecordError(ctx, err)
				return nil, err
			}
		}

		c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			err = errors.Wrap(err, "tls: failed creating key pair")
			span.RecordError(ctx, err)
			return nil, err
		}

		config.Certificates = append(config.Certificates, c)
	}

	return config, nil
}

func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
