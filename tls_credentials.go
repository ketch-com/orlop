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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/switch-bit/orlop/log"
)

// NewServerTLSConfig returns a new tls.VaultConfig from the given configuration input
func NewServerTLSConfig(cfg HasTLSConfig, vault HasVaultConfig) (*tls.Config, error) {
	config := &tls.Config{
		ClientAuth: cfg.GetClientAuth(),
		MinVersion: tls.VersionTLS12,
	}

	if !strSliceContains(config.NextProtos, "http/1.1") {
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

	var err error
	t := CloneTLSConfig(cfg)

	certGenerator := NewCertificateGenerator(cfg.GetGenerate(), vault)

	err = certGenerator.GenerateCertificates(&t.Cert.Secret, &t.Key.Secret)
	if err != nil {
		return nil, err
	}

	certPEMBlock, err := LoadKey(t.GetCert(), vault, "certificate")
	if err == nil {
		log.Trace("certificate loaded")

		keyPEMBlock, err := LoadKey(t.GetKey(), vault, "private_key")
		if err != nil {
			log.WithError(err).Error("error loading private key")
			return nil, err
		}

		if t.GetRootCA().GetEnabled() {
			rootcaPEMBlock, err := LoadKey(t.GetRootCA(), vault, "issuing_ca")
			if err == nil {
				config.ClientCAs = x509.NewCertPool()

				if !config.ClientCAs.AppendCertsFromPEM(rootcaPEMBlock) {
					log.WithError(err).Error("failed to append client CA certificates")
					return nil, fmt.Errorf("tls: failed to append client CA certificates")
				}
			} else {
				log.WithError(err).Error("error loading CA")
			}
		}

		c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			log.WithError(err).Error("error creating key pair")
			return nil, err
		}

		config.Certificates = append(config.Certificates, c)
	} else {
		log.WithError(err).Error("error loading certificate")
	}

	return config, nil
}

// NewClientTLSConfig returns a new tls.VaultConfig from the given configuration input
func NewClientTLSConfig(cfg HasTLSConfig, vault HasVaultConfig) (*tls.Config, error) {
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

	certGenerator := NewCertificateGenerator(cfg.GetGenerate(), vault)

	err = certGenerator.GenerateCertificates(&t.Cert.Secret, &t.Key.Secret)
	if err != nil {
		return nil, err
	}

	certPEMBlock, err := LoadKey(t.GetCert(), vault, "certificate")
	if err == nil {
		keyPEMBlock, err := LoadKey(t.GetKey(), vault, "private_key")
		if err != nil {
			log.WithError(err).Error("error loading private key")
			return nil, err
		}

		config.RootCAs = x509.NewCertPool()

		if !config.RootCAs.AppendCertsFromPEM(certPEMBlock) {
			log.WithError(err).Error("failed to append certificates")
			return nil, fmt.Errorf("tls: failed to append certificates")
		}

		if t.GetRootCA().GetEnabled() {
			rootcaPEMBlock, err := LoadKey(t.GetRootCA(), vault, "issuing_ca")
			if err == nil {
				if !config.RootCAs.AppendCertsFromPEM(rootcaPEMBlock) {
					log.WithError(err).Error("failed to append root certificates")
					return nil, fmt.Errorf("tls: failed to append root certificates")
				}
			} else {
				log.WithError(err).Error("error loading CA")
			}
		}

		c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			log.WithError(err).Error("error creating key pair")
			return nil, err
		}

		config.Certificates = append(config.Certificates, c)
	} else {
		log.WithError(err).Error("error loading certificate")
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
