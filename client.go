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
	"fmt"
	"github.com/switch-bit/orlop/log"
	"google.golang.org/grpc"
)

// Connect creates a new client from configuration
func Connect(cfg HasClientConfig, vault HasVaultConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	l := log.WithField("url", cfg.GetURL())

	if len(cfg.GetURL()) == 0 {
		l.Errorf("client: url required")
		return nil, fmt.Errorf("client: url required")
	}

	if cfg.GetTLS().GetEnabled() {
		l.Trace("tls enabled")
		creds, err := NewClientTLSCredentials(cfg.GetTLS(), vault)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		l.Trace("tls disabled")
		opts = append(opts, grpc.WithInsecure())
	}

	shared := cfg.GetToken().GetShared()
	if len(shared.GetID()) > 0 || len(shared.GetFile()) > 0 || len(shared.GetSecret()) > 0 {
		l.Trace("loading token from configuration")

		s, err := LoadKey(shared, vault, "secret")
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithPerRPCCredentials(SharedContextCredentials{
			token: string(s),
		}))
	} else {
		l.Trace("using context credentials")

		opts = append(opts, grpc.WithPerRPCCredentials(ContextCredentials{}))
	}

	l.Trace("dialling")
	conn, err := grpc.Dial(cfg.GetURL(), opts...)
	if err != nil {
		l.WithError(err).Error("failed dialling")
		return nil, err
	}

	return conn, nil
}
