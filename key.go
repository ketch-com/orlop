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
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/log"
	"go.opentelemetry.io/otel/attribute"
	"io/ioutil"
)

// LoadKey loads the key material based on the config
func LoadKey(ctx context.Context, cfg HasKeyConfig, vault HasVaultConfig, which string) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "LoadKey")
	defer span.End()

	// If the key is not enabled, return an empty byte array
	if en, ok := cfg.(HasEnabled); ok && !en.GetEnabled() {
		return nil, nil
	}

	fields := logrus.Fields{
		"key":           which,
		"vault.enabled": vault != nil && vault.GetEnabled(),
	}

	method := "none"

	if len(cfg.GetFile()) > 0 {
		method = "file"
		fields["key.file"] = cfg.GetFile()
		span.SetAttributes(attribute.String("key.file", cfg.GetFile()))
	}

	if len(cfg.GetID()) > 0 {
		if vault != nil && vault.GetEnabled() {
			method = "id"
		}
		fields["key.id"] = cfg.GetID()
		span.SetAttributes(attribute.String("key.id", cfg.GetID()))
	}

	if len(cfg.GetSecret()) > 0 {
		method = "secret"
		fields["key.secret"] = "*********"
		span.SetAttributes(attribute.String("key.secret", "*********"))
	}

	fields["method"] = method
	span.SetAttributes(attribute.String("key.method", method))
	l := log.WithFields(fields)

	switch method {
	case "secret":
		l.Trace("key found")
		return cfg.GetSecret(), nil

	case "id":
		client, err := NewVault(ctx, vault)
		if err != nil {
			err = errors.Wrap(err, "key: could not connect to Vault")
			span.RecordError(err)
			return nil, err
		}

		s, err := client.Read(ctx, cfg.GetID())
		if err != nil {
			err = errors.Wrap(err, "key: not found")
			span.RecordError(err)
			return nil, err
		}

		if s == nil || s.Data[which] == nil {
			err = errors.New("key: not found")
			span.RecordError(err)
			return nil, err
		}

		l.Tracef("key found")
		return []byte(s.Data[which].(string)), nil

	case "file":
		key, err := ioutil.ReadFile(cfg.GetFile())
		if err != nil {
			err = errors.Wrap(err, "key: not found")
			span.RecordError(err)
			return nil, err
		}

		l.Trace("key found")
		return key, nil
	}

	err := errors.Errorf("key: no key configured for %s", which)
	span.RecordError(err)
	return nil, err
}

// LoadPrivateKey loads a private key from the given bytes
func LoadPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	log.Trace("loading private key")
	return jwt.ParseRSAPrivateKeyFromPEM(key)
}

// LoadPublicKeys loads an array of public keys from the given bytes
func LoadPublicKeys(key []byte) (publicKeys []*rsa.PublicKey, err error) {
	log.Trace("loading public key")

	for len(key) > 0 {
		// Parse PEM block
		var block *pem.Block
		if block, key = pem.Decode(key); block == nil {
			return nil, jwt.ErrKeyMustBePEMEncoded
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
			return nil, jwt.ErrNotRSAPublicKey
		}

		publicKeys = append(publicKeys, pkey)
	}

	return
}
