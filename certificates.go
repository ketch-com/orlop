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

// GenerateCertificates calls Vault to generate a certificate
func GenerateCertificates(vault HasVaultConfig, cfg HasCertGenerationConfig,
	cert *string, key *string) error {
	// If Vault not enabled or certificate generation not enabled, just return
	if !vault.GetEnabled() || !cfg.GetEnabled() {
		return nil
	}

	// Connect to Vault
	client, err := NewVault(vault)
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
	secret, err := client.Write(cfg.GetPath(), params)
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
