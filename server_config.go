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

// HttpLoggingConfig provides configuration for HTTP logging
type HttpLoggingConfig struct {
	Enabled bool
	Headers []string `config:",default=X-Forwarded-For"`
}

// HasServerConfig denotes an object provides server configuration
type HasServerConfig interface {
	GetBind() string
	GetListen() uint
	GetTLS() HasTLSConfig
	GetLoopback() HasClientConfig
	GetLogging() HttpLoggingConfig
}

// ServerConfig is standard configuration of most server commands
type ServerConfig struct {
	Bind     string `config:"bind,default=0.0.0.0"`
	Listen   uint   `config:"listen,default=5000"`
	TLS      TLSConfig
	Loopback ClientConfig
	Logging  HttpLoggingConfig
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

// GetLoopback returns GRPC gateway loopback client configuration
func (c ServerConfig) GetLoopback() HasClientConfig {
	return c.Loopback
}

// GetLogging returns logging configuration
func (c ServerConfig) GetLogging() HttpLoggingConfig {
	return c.Logging
}
