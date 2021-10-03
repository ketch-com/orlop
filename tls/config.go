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
	"crypto/tls"
	"fmt"
	"go.ketch.com/lib/orlop/v2/config"
	"go.ketch.com/lib/orlop/v2/pem"
	"reflect"
)

// Config provides TLS configuration
type Config struct {
	ClientAuth tls.ClientAuthType
	Enabled    bool `config:"enabled,default=true"`
	Insecure   bool
	Override   string
	Cert       pem.CertificateConfig
	Key        pem.KeyConfig
	RootCA     pem.RootCAConfig `config:"rootca"`
}

// GetEnabled returns true of TLS is enabled
func (t Config) GetEnabled() bool {
	return t.Enabled
}

func init() {
	config.RegisterConfigParser("tls.ClientAuthType", func(value reflect.Value, input string) error {
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
