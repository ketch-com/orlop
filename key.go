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
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/switch-bit/orlop/log"
	"io/ioutil"
)

// LoadPrivateKeyFromVault loads public key material from vault
func LoadPrivateKeyFromVault(cfg HasKeyConfig, vault HasVaultConfig) (*rsa.PrivateKey, error) {
	key, err := LoadKey(cfg, vault, "private")
	if err != nil {
		return nil, err
	}

	return LoadPrivateKey(key)
}

// LoadPublicKeysFromVault loads public key material from vault
func LoadPublicKeysFromVault(cfg HasKeyConfig, vault HasVaultConfig) ([]*rsa.PublicKey, error) {
	key, err := LoadKey(cfg, vault, "public")
	if err != nil {
		return nil, err
	}

	return LoadPublicKeys(key)
}

// LoadKey loads the key material based on the config
func LoadKey(cfg HasKeyConfig, vault HasVaultConfig, which string) ([]byte, error) {
	secret := cfg.GetSecret()
	if len(secret) > 0 {
		secret = "***********"
	}
	l := log.WithFields(logrus.Fields{
		"key":           which,
		"key.id":        cfg.GetID(),
		"key.file":      cfg.GetFile(),
		"key.secret":    secret,
		"vault.enabled": vault != nil && vault.GetEnabled(),
	})
	l.Debug("loading key")

	if len(cfg.GetSecret()) > 0 {
		l.WithField("method", "secret").Trace("key provided as secret")
		return base64.StdEncoding.DecodeString(cfg.GetSecret())
	}

	if len(cfg.GetID()) > 0 && vault != nil && vault.GetEnabled() {
		l = l.WithField("method", "id")
		l.Tracef("using vault to lookup ID %s", cfg.GetID())

		client, err := NewVault(vault)
		if err != nil {
			l.WithError(err).Error("could not connect to Vault")
			return nil, err
		}

		s, err := client.Read(cfg.GetID())
		if err != nil {
			l.WithError(err).Error("key not found")
			return nil, err
		}

		if s != nil && s.Data[which] != nil {
			l.Tracef("key found")
			return []byte(s.Data[which].(string)), nil
		}

		l.Tracef("key not found")
		return nil, fmt.Errorf("key: could not load key from %s", cfg.GetID())
	}

	if len(cfg.GetFile()) > 0 {
		l = l.WithField("method", "file").WithField("file", cfg.GetFile())
		l.Trace("loading key from file")

		key, err := ioutil.ReadFile(cfg.GetFile())
		if err != nil {
			l.Trace("key not found")
			return nil, err
		}

		l.Trace("key found")

		return key, nil
	}

	return nil, fmt.Errorf("key: no key configured for %s", which)
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
