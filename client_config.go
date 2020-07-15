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

import "time"

// HasClientConfig denotes that an object provides client configuration
type HasClientConfig interface {
	GetTLS() HasTLSConfig
	GetToken() HasTokenConfig
	GetURL() string
	GetHeaders() map[string]string
	GetWriteBufferSize() int
	GetReadBufferSize() int
	GetInitialWindowSize() int32
	GetInitialConnWindowSize() int32
	GetMaxCallRecvMsgSize() int
	GetMinConnectTimeout() time.Duration
	GetBlock() bool
	GetConnTimeout() time.Duration
	GetUserAgent() string
}

// ClientConfig is standard configuration of most client commands
type ClientConfig struct {
	URL                   string
	Token                 TokenConfig
	TLS                   TLSConfig
	Headers               map[string]string
	WriteBufferSize       int
	ReadBufferSize        int
	InitialWindowSize     int32
	InitialConnWindowSize int32
	MaxCallRecvMsgSize    int
	MinConnectTimeout     time.Duration
	Block                 bool
	ConnTimeout           time.Duration
	UserAgent             string
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

func (c ClientConfig) GetWriteBufferSize() int {
	return c.WriteBufferSize
}

func (c ClientConfig) GetReadBufferSize() int {
	return c.ReadBufferSize
}

func (c ClientConfig) GetInitialWindowSize() int32 {
	return c.InitialWindowSize
}

func (c ClientConfig) GetInitialConnWindowSize() int32 {
	return c.InitialConnWindowSize
}

func (c ClientConfig) GetMaxCallRecvMsgSize() int {
	return c.MaxCallRecvMsgSize
}

func (c ClientConfig) GetMinConnectTimeout() time.Duration {
	return c.MinConnectTimeout
}

func (c ClientConfig) GetBlock() bool {
	return c.Block
}

func (c ClientConfig) GetConnTimeout() time.Duration {
	return c.ConnTimeout
}

func (c ClientConfig) GetUserAgent() string {
	return c.UserAgent
}
