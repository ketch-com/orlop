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
	vault "github.com/hashicorp/vault/api"
	"github.com/switch-bit/orlop/log"
	"net/http"
	"path"
)

// VaultClient is a Vault client
type VaultClient struct {
	cfg    HasVaultConfig
	client *vault.Client
}

// NewVault connects to Vault given the configuration
func NewVault(cfg HasVaultConfig) (*VaultClient, error) {
	var err error

	// First check if Vault is enabled in config, returning if not
	if !cfg.GetEnabled() {
		log.Trace("vault is not enabled")
		return &VaultClient{
			cfg: cfg,
		}, nil
	}

	// Setup the Vault native config
	vc := &vault.Config{
		Address: cfg.GetAddress(),
	}

	// If TLS is enabled, then setup the TLS configuration
	if cfg.GetTLS().GetEnabled() {
		vc.HttpClient = &http.Client{}

		t := http.DefaultTransport.(*http.Transport).Clone()

		t.TLSClientConfig, err = NewClientTLSConfig(cfg.GetTLS(), &VaultConfig{Enabled: false})
		if err != nil {
			return nil, err
		}

		vc.HttpClient.Transport = t
	}

	// Create the vault native client
	client, err := vault.NewClient(vc)
	if err != nil {
		return nil, err
	}

	// Set the token on the client
	client.SetToken(cfg.GetToken())

	return &VaultClient{
		cfg:    cfg,
		client: client,
	}, nil
}

// Reads a secret at the given path
func (c VaultClient) Read(p string) (*vault.Secret, error) {
	// If the client isn't connected (because Vault is not enabled, just return an empty secret).
	if c.client == nil {
		return &vault.Secret{
			RequestID:     "",
			LeaseID:       "",
			LeaseDuration: 0,
			Renewable:     false,
			Data:          make(map[string]interface{}),
			Warnings:      nil,
			Auth:          nil,
			WrapInfo:      nil,
		}, nil
	}

	keyPath := c.prefix(p)

	sec, err := c.client.Logical().Read(keyPath)
	if err != nil {
		log.WithField("vault.path", keyPath).WithError(err).Trace("read")
		return nil, err
	}

	return sec, nil
}

// Writes secret data at the given path
func (c VaultClient) Write(p string, data map[string]interface{}) (*vault.Secret, error) {
	if c.client == nil {
		return &vault.Secret{
			RequestID:     "",
			LeaseID:       "",
			LeaseDuration: 0,
			Renewable:     false,
			Data:          make(map[string]interface{}),
			Warnings:      nil,
			Auth:          nil,
			WrapInfo:      nil,
		}, nil
	}

	keyPath := c.prefix(p)

	sec, err := c.client.Logical().Write(keyPath, data)
	if err != nil {
		log.WithField("vault.path", keyPath).WithError(err).Trace("write")
		return nil, err
	}

	return sec, nil
}

// Delete deletes a secret at the given path
func (c VaultClient) Delete(p string) error {
	if c.client == nil {
		return nil
	}

	keyPath := c.prefix(p)

	_, err := c.client.Logical().Delete(keyPath)
	if err != nil {
		log.WithField("vault.path", keyPath).WithError(err).Trace("delete")
		return err
	}

	return nil
}

func (c VaultClient) prefix(p string) string {
	return path.Join(c.cfg.GetPrefix(), p)
}
