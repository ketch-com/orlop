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
	"github.com/spf13/pflag"
)

// HasServerConfig denotes an object provides server configuration
type HasServerConfig interface {
	GetBind() string
	GetListen() uint
	GetTLS() HasTLSConfig
	GetSwagger() HasEnabled
}

// ServerConfig is standard configuration of most server commands
type ServerConfig struct {
	Bind    string
	Listen  uint
	TLS     TLSConfig
	Swagger Enabled
}

// GetBind returns the address to bind to
func (c ServerConfig) GetBind() string {
	return c.Bind
}

// GetListen returns the port to listen on
func (c ServerConfig) GetListen() uint {
	return c.Listen
}

// GetTLS returns TLS configuration
func (c ServerConfig) GetTLS() HasTLSConfig {
	return c.TLS
}

// GetSwagger returns Enabled for Swagger
func (c ServerConfig) GetSwagger() HasEnabled {
	return c.Swagger
}

// NewServerConfig returns a new, unmarshaled standard ServerConfig
func NewServerConfig() (*ServerConfig, error) {
	c := new(ServerConfig)

	err := Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// AddServer adds the server-related parameters
func AddServer(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	AddTLS(flags, append(prefix, "tls")...)
	AddEnabled(flags, "enable Swagger", false, append(prefix, "swagger")...)
	flags.String(p("bind"), "0.0.0.0", "address to bind to")
	flags.Uint(p("listen"), 5000, "port to bind to")
}
