// Copyright (c) 2020 Ketch, Inc.
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
	"go.ketch.com/lib/orlop/errors"
)

// GenerateCertificates calls Vault to generate a certificate
//
// deprecated: use GenerateCertificatesContext
func GenerateCertificates(vault HasVaultConfig, cfg HasCertGenerationConfig, cert *[]byte, key *[]byte) error {
	return GenerateCertificatesContext(context.TODO(), vault, cfg, cert, key)
}

// GenerateCertificatesContext calls Vault to generate a certificate
func GenerateCertificatesContext(ctx context.Context, vault HasVaultConfig, cfg HasCertGenerationConfig, cert *[]byte, key *[]byte) error {
	ctx, span := tracer.Start(ctx, "GenerateCertificates")
	defer span.End()

	// If Vault not enabled or certificate generation not enabled, just return
	if !vault.GetEnabled() || !cfg.GetEnabled() {
		return nil
	}

	// Connect to Vault
	client, err := NewVaultContext(ctx, vault)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"common_name": cfg.GetCommonName(),
		"format":      "pem_bundle",
	}
	if len(cfg.GetAltNames()) > 0 {
		params["alt_names"] = cfg.GetAltNames()
	}
	if cfg.GetTTL().Seconds() > 60 {
		params["ttl"] = cfg.GetTTL().String()
	}

	// Write the params to the path to generate the certificate
	secret, err := client.WriteContext(ctx, cfg.GetPath(), params)
	if err != nil {
		err = errors.Wrap(err, "generate: failed to write to Vault")
		span.RecordError(ctx, err)
		return err
	}

	// Set the generated certificate and private key
	if d, ok := secret.Data["certificate"]; ok {
		if s, ok := d.(string); ok {
			*cert = []byte(s)
			return nil
		}
	}

	if d, ok := secret.Data["private_key"]; ok {
		if s, ok := d.(string); ok {
			*key = []byte(s)
			return nil
		}
	}

	return nil
}
