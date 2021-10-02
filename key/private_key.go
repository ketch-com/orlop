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

package key

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	encoding_pem "encoding/pem"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/logging"
	"go.ketch.com/lib/orlop/pem"
	"go.uber.org/fx"
)

var (
	ErrKeyMustBePEMEncodedPrivateKey = errors.New("key: must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey              = errors.New("key: not a valid RSA private key")
)

type PrivateKeyProvider interface {
	// LoadPrivateKey loads a private key from the given bytes
	LoadPrivateKey(ctx context.Context, cfg Config) (*rsa.PrivateKey, error)
}

type implPrivateKeyProvider struct {
	params PrivateKeyParams
}

type PrivateKeyParams struct {
	fx.In

	Logger logging.Logger
	PEM    pem.Provider
}

func NewPrivateKey(params PrivateKeyParams) PrivateKeyProvider {
	return &implPrivateKeyProvider{
		params: params,
	}
}

// LoadPrivateKey loads a private key from the given bytes
func (p *implPrivateKeyProvider) LoadPrivateKey(ctx context.Context, cfg Config) (*rsa.PrivateKey, error) {
	p.params.Logger.Trace("loading private key")

	key, err := p.params.PEM.Load(ctx, pem.PrivateKeyConfig{
		ID:   cfg.ID,
		File: cfg.File,
	})
	if err != nil {
		return nil, err
	}

	// Parse PEM block
	var block *encoding_pem.Block
	if block, _ = encoding_pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncodedPrivateKey
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}
