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
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/pem"
	"go.uber.org/fx"
)

var (
	ErrKeyMustBePEMEncodedPublicKey = errors.New("key: must be PEM encoded PKIX public key or certificate")
	ErrNotRSAPublicKey              = errors.New("key: not a valid RSA public key")
)

type PublicKeyProvider interface {
	// LoadPublicKeys loads an array of public keys from the given bytes
	LoadPublicKeys(ctx context.Context, cfg Config) (publicKeys []*rsa.PublicKey, err error)
}

type implPublicKeyProvider struct {
	params PublicKeyParams
}

type PublicKeyParams struct {
	fx.In

	Logger logging.Logger
	PEM    pem.Provider
}

func NewPublicKey(params PublicKeyParams) PublicKeyProvider {
	return &implPublicKeyProvider{
		params: params,
	}
}

// LoadPublicKeys loads an array of public keys from the given bytes
func (p *implPublicKeyProvider) LoadPublicKeys(ctx context.Context, cfg Config) (publicKeys []*rsa.PublicKey, err error) {
	p.params.Logger.Trace("loading public key")

	key, err := p.params.PEM.Load(ctx, pem.PublicKeyConfig{
		ID:   cfg.ID,
		File: cfg.File,
	})
	if err != nil {
		return nil, err
	}

	for len(key) > 0 {
		// Parse PEM block
		var block *encoding_pem.Block
		if block, key = encoding_pem.Decode(key); block == nil {
			return nil, ErrKeyMustBePEMEncodedPublicKey
		}

		// Parse the key
		var parsedKey interface{}
		if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				parsedKey = cert.PublicKey
			} else {
				return nil, err
			}
		}

		var pkey *rsa.PublicKey
		var ok bool
		if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
			return nil, ErrNotRSAPublicKey
		}

		publicKeys = append(publicKeys, pkey)
	}

	return
}
