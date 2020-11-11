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
	"crypto/tls"
	"fmt"
	"reflect"
)

// HasTLSConfig denotes that an object supports TLS configuration
type HasTLSConfig interface {
	GetCert() KeyConfig
	GetClientAuth() tls.ClientAuthType
	GetEnabled() bool
	GetInsecure() bool
	GetKey() KeyConfig
	GetOverride() string
	GetRootCA() KeyConfig
	GetGenerate() CertGenerationConfig
}

// TLSConfig provides TLS configuration
type TLSConfig struct {
	ClientAuth tls.ClientAuthType `config:"clientauth"`
	Enabled    bool               `config:"enabled,default=true"`
	Insecure   bool
	Override   string
	Cert       KeyConfig
	Key        KeyConfig
	RootCA     KeyConfig `config:"rootca"`
	Generate   CertGenerationConfig
}

// GetCert returns the certificate key configuration
func (t TLSConfig) GetCert() KeyConfig {
	return t.Cert
}

// GetClientAuth returns the Client Authentication Type
func (t TLSConfig) GetClientAuth() tls.ClientAuthType {
	return t.ClientAuth
}

// GetEnabled returns true of TLS is enabled
func (t TLSConfig) GetEnabled() bool {
	return t.Enabled
}

// GetInsecure returns true if insecure mode is enabled (use at your own peril)
func (t TLSConfig) GetInsecure() bool {
	return t.Insecure
}

// GetKey returns the private key configuration
func (t TLSConfig) GetKey() KeyConfig {
	return t.Key
}

// GetOverride returns an override name for the server
func (t TLSConfig) GetOverride() string {
	return t.Override
}

// GetRootCA returns the root CA key configuration
func (t TLSConfig) GetRootCA() KeyConfig {
	return t.RootCA
}

// GetGenerate returns the certificate generation configuration
func (t TLSConfig) GetGenerate() CertGenerationConfig {
	return t.Generate
}

// CloneTLSConfig clones the given TLS configuration
func CloneTLSConfig(cfg HasTLSConfig) TLSConfig {
	return TLSConfig{
		ClientAuth: cfg.GetClientAuth(),
		Enabled:    cfg.GetEnabled(),
		Insecure:   cfg.GetInsecure(),
		Override:   cfg.GetOverride(),
		Cert:       cfg.GetCert(),
		Key:        cfg.GetKey(),
		RootCA:     cfg.GetRootCA(),
		Generate:   cfg.GetGenerate(),
	}
}

func init() {
	RegisterConfigParser("tls.ClientAuthType", func(value reflect.Value, input string) error {
		if len(input) == 0 {
			return nil
		}

		var clientAuthType tls.ClientAuthType

		switch input {
		case "request":
			clientAuthType = tls.RequestClientCert

		case "require":
			clientAuthType = tls.RequireAnyClientCert

		case "verify":
			clientAuthType = tls.VerifyClientCertIfGiven

		case "require-and-verify":
			clientAuthType = tls.RequireAndVerifyClientCert

		case "none":
			clientAuthType = tls.NoClientCert

		default:
			return fmt.Errorf("config: '%s' is not a valid ClientAuth value", input)
		}

		value.Set(reflect.ValueOf(clientAuthType))

		return nil
	})
}
