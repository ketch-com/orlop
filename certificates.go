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

import "encoding/base64"

// CertificateGenerator defines an object that (might) generate certificates
type CertificateGenerator interface {
	GenerateCertificates(cert *string, key *string) error
}

// NewCertificateGenerator returns an appropriate implementation to generate certificates
func NewCertificateGenerator(cfg HasCertGenerationConfig, vault HasVaultConfig) CertificateGenerator {
	// If Vault not enabled or certificate generation not enabled, just return
	if !vault.GetEnabled() || !cfg.GetEnabled() {
		return &NullCertificateGenerator{}
	}

	return &VaultCertificateGenerator{
		cfg:   cfg,
		vault: vault,
	}
}

// NullCertificateGenerator does not generate certificates
type NullCertificateGenerator struct {}

// GenerateCertificates in this case does not generate any certificates
func (v *NullCertificateGenerator) GenerateCertificates(cert *string, key *string) error {
	return nil
}

// VaultCertificateGenerator generates certificates using Vault
type VaultCertificateGenerator struct {
	cfg   HasCertGenerationConfig
	vault HasVaultConfig
}

// GenerateCertificates calls Vault to generate a certificate
func (v *VaultCertificateGenerator) GenerateCertificates(cert *string, key *string) error {
	// Connect to Vault
	client, err := NewVault(v.vault)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"common_name": v.cfg.GetCommonName(),
		"format":      "pem_bundle",
	}
	if len(v.cfg.GetAltNames()) > 0 {
		params["alt_names"] = v.cfg.GetAltNames()
	}
	if v.cfg.GetTTL().Seconds() > 60 {
		params["ttl"] = v.cfg.GetTTL().String()
	}

	// Write the params to the path to generate the certificate
	secret, err := client.Write(v.cfg.GetPath(), params)
	if err != nil {
		return err
	}

	// Set the generated certificate and private key as secrets
	*cert = decodeCertInfo(secret.Data, "certificate")
	*key = decodeCertInfo(secret.Data, "private_key")

	return nil
}

func decodeCertInfo(data map[string]interface{}, key string) string {
	if d, ok := data[key]; ok {
		if s, ok := d.(string); ok {
			return base64.StdEncoding.EncodeToString([]byte(s))
		}
	}

	return ""
}
