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
	"time"
)

// HasCertGenerationConfig denotes that an object provides certificate generation configuration.
type HasCertGenerationConfig interface {
	GetEnabled() bool
	GetPath() string
	GetCommonName() string
	GetAltNames() string
	GetTTL() time.Duration
}

// CertGenerationConfig provides the certificate generation configuration
type CertGenerationConfig struct {
	Enabled    bool
	Path       string `config:"path,default=/pki/issue/"`
	CommonName string
	AltNames   string
	TTL        time.Duration
}

// GetEnabled returns true if certificate generation is enabled
func (c CertGenerationConfig) GetEnabled() bool {
	return c.Enabled && len(c.CommonName) > 0
}

// GetPath returns the path to the certificate generator
func (c CertGenerationConfig) GetPath() string {
	return c.Path
}

// GetCommonName returns the common name to use for certificates
func (c CertGenerationConfig) GetCommonName() string {
	return c.CommonName
}

// GetAltNames returns alternate names to use for certificates
func (c CertGenerationConfig) GetAltNames() string {
	return c.AltNames
}

// GetTTL returns the time-to-live for the certificates
func (c CertGenerationConfig) GetTTL() time.Duration {
	return c.TTL
}
