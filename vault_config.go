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

// HasVaultConfig denotes that the given object has vault configuration
type HasVaultConfig interface {
	GetEnabled() bool
	GetAddress() string
	GetToken() string
	GetPrefix() string
	GetTLS() HasTLSConfig
}

// VaultConfig provides the configuration options available for Vault
type VaultConfig struct {
	Enabled bool
	Address string
	Token   string
	Prefix  string
	TLS     TLSConfig
}

// GetEnabled returns true if Vault is enabled
func (c VaultConfig) GetEnabled() bool {
	return c.Enabled
}

// GetAddress returns the address of the Vault server
func (c VaultConfig) GetAddress() string {
	return c.Address
}

// GetToken returns the unwrapping token to use for Vault
func (c VaultConfig) GetToken() string {
	return c.Token
}

// GetPrefix returns the prefix to prepend to all paths when requesting Vault
func (c VaultConfig) GetPrefix() string {
	return c.Prefix
}

// GetTLS returns the TLS configuration
func (c VaultConfig) GetTLS() HasTLSConfig {
	return c.TLS
}
