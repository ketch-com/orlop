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

package client

import (
	"fmt"
	"go.ketch.com/lib/orlop/secret"
	"go.ketch.com/lib/orlop/tls"
	"go.ketch.com/lib/orlop/version"
	"time"
)

// Config is standard configuration of most client commands
type Config struct {
	Name                  string
	URL                   string
	Host                  string
	Port                  int32
	Token                 secret.Config
	TLS                   tls.Config
	Headers               map[string]string
	WriteBufferSize       int
	ReadBufferSize        int
	InitialWindowSize     int32
	InitialConnWindowSize int32
	MaxCallRecvMsgSize    int
	MaxCallSendMsgSize    int
	MinConnectTimeout     time.Duration
	ConnTimeout           time.Duration
	Block                 bool
	UserAgent             string
}

// GetName returns the Name of the client config
func (c Config) GetName() string {
	if len(c.Name) == 0 {
		return "unknown"
	}
	return c.Name
}

// GetURL returns the URL to contact the server
func (c Config) GetURL() string {
	if len(c.URL) == 0 && len(c.Host) > 0 {
		if c.Port != 0 {
			return fmt.Sprintf("https://%s:%d", c.Host, c.Port)
		}

		return fmt.Sprintf("https://%s", c.Host)
	}

	return c.URL
}

// GetUserAgent returns the user agent
func (c Config) GetUserAgent() string {
	if len(c.UserAgent) == 0 {
		return fmt.Sprintf("%s/%s", version.Name, version.Version)
	}

	return c.UserAgent
}
