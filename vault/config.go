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

package vault

import (
	"fmt"
	"go.ketch.com/lib/orlop/tls"
	"time"
)

// Config provides the configuration options available for Vault
type Config struct {
	Enabled bool
	Address string
	Host    string
	Port    int32
	Token   string
	Prefix  string
	TLS     tls.Config
}

// GetEnabled returns true if Vault is enabled
func (c Config) GetEnabled() bool {
	return c.Enabled
}

// GetURL returns the URL to contact the server
func (c Config) GetURL() string {
	if len(c.Address) == 0 && len(c.Host) > 0 {
		if c.Port != 0 {
			return fmt.Sprintf("https://%s:%d", c.Host, c.Port)
		}

		return fmt.Sprintf("https://%s", c.Host)
	}

	return c.Address
}

// GeneratorConfig provides the certificate generation configuration
type GeneratorConfig struct {
	Enabled    bool
	Path       string `config:"path,default=/pki/issue/"`
	CommonName string
	AltNames   string
	TTL        time.Duration
}

// GetEnabled returns true if certificate generation is enabled
func (c GeneratorConfig) GetEnabled() bool {
	return c.Enabled && len(c.CommonName) > 0
}
