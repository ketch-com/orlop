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
	"time"
)

// HasTokenConfig denotes that an object supports providing token configuration
type HasTokenConfig interface {
	GetIssuer() string
	GetKeyMap() HasKeyConfig
	GetPrivateKey() HasKeyConfig
	GetPublicKey() HasKeyConfig
	GetShared() HasKeyConfig
	GetTTL() time.Duration
}

// TokenConfig is the configuration for managing tokens
type TokenConfig struct {
	Issuer     string
	KeyMap     KeyConfig `config:"keymap"`
	PrivateKey KeyConfig `config:"privatekey"`
	PublicKey  KeyConfig `config:"publickey"`
	Shared     KeyConfig
	TTL        time.Duration
}

// GetTTL returns the time-to-live for the token
func (t TokenConfig) GetTTL() time.Duration {
	return t.TTL
}

// GetPrivateKey returns the private key configuration
func (t TokenConfig) GetPrivateKey() HasKeyConfig {
	return t.PrivateKey
}

// GetPublicKey returns the public key configuration
func (t TokenConfig) GetPublicKey() HasKeyConfig {
	return t.PublicKey
}

// GetShared returns a shared secret/key configuration
func (t TokenConfig) GetShared() HasKeyConfig {
	return t.Shared
}

// GetIssuer returns the issuer
func (t TokenConfig) GetIssuer() string {
	return t.Issuer
}

// GetKeyMap returns a key map
func (t TokenConfig) GetKeyMap() HasKeyConfig {
	return t.KeyMap
}
