// Copyright (c) 2020 Ketch Kloud, Inc.
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
	"encoding/pem"
	"fmt"
	"reflect"
)

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

// GetEnabled returns true of TLS is enabled
func (t TLSConfig) GetEnabled() bool {
	return t.Enabled
}

// CloneTLSConfig clones the given TLS configuration
func CloneTLSConfig(cfg TLSConfig) TLSConfig {
	return TLSConfig{
		ClientAuth: cfg.ClientAuth,
		Enabled:    cfg.Enabled,
		Insecure:   cfg.Insecure,
		Override:   cfg.Override,
		Cert:       cfg.Cert,
		Key:        cfg.Key,
		RootCA:     cfg.RootCA,
		Generate:   cfg.Generate,
	}
}

// FromTLSConfig returns a TLSConfig for the given tls.Config
func FromTLSConfig(t *tls.Config) TLSConfig {
	if t == nil || len(t.Certificates) < 1 {
		return TLSConfig{}
	}

	bb := t.Certificates[0].Certificate

	b := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: bb[0],
	})

	return TLSConfig{
		Enabled:  true,
		Insecure: t.InsecureSkipVerify,
		Override: t.ServerName,
		Cert: KeyConfig{
			Secret: b,
		},
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
