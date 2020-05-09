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

import "github.com/spf13/pflag"

// HasClientConfig denotes that an object provides client configuration
type HasClientConfig interface {
	GetTLS() HasTLSConfig
	GetToken() HasTokenConfig
	GetURL() string
}

// ClientConfig is standard configuration of most client commands
type ClientConfig struct {
	URL     string
	Token   TokenConfig
	TLS     TLSConfig
	Headers map[string]string
}

// GetURL returns the URL to contact the server
func (c ClientConfig) GetURL() string {
	return c.URL
}

// GetToken returns the security token configuration information
func (c ClientConfig) GetToken() HasTokenConfig {
	return c.Token
}

// GetTLS returns the TLS configuration
func (c ClientConfig) GetTLS() HasTLSConfig {
	return c.TLS
}

// GetHeaders returns static headers to add to requests
func (c ClientConfig) GetHeaders() map[string]string {
	return c.Headers
}

// AddClient adds the client-related parameters
func AddClient(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	AddTLS(flags, append(prefix, "tls")...)
	AddToken(flags, p("token"))
	flags.String(p("url"), "", "URL to connect to")
	flags.StringToString(p("headers"), nil, "additional headers to send")
}
