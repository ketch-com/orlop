package orlop

import (
	"crypto/tls"
)

// TLSProvider defines how TLS configs can be returned based on a setup
type TLSProvider interface {
	// NewClientTLSConfig creates a new client-focused tls.Config for the defined configuration
	NewClientTLSConfig(cfg HasTLSConfig) (*tls.Config, error)

	// NewServerTLSConfig creates a new server-focused tls.Config for the defined configuration
	NewServerTLSConfig(cfg HasTLSConfig) (*tls.Config, error)
}

// NewTLSProvider returns an appropriate TLSProvider based on the Vault configuration
func NewTLSProvider(vault HasVaultConfig) TLSProvider {
	if vault.GetEnabled() {
		return NewVaultTLSProvider(vault)
	}

	return NewSimpleTLSProvider()
}

// VaultTLSProvider is a TLS provider connected to Vault
type VaultTLSProvider struct {
	vault HasVaultConfig
}

// NewVaultTLSProvider returns a VaultTLSProvider
func NewVaultTLSProvider(vault HasVaultConfig) *VaultTLSProvider {
	return &VaultTLSProvider{
		vault: vault,
	}
}

// NewClientTLSConfig creates a new client-focused tls.Config for the defined configuration
func (p *VaultTLSProvider) NewClientTLSConfig(cfg HasTLSConfig) (*tls.Config, error) {
	return NewClientTLSConfig(cfg, p.vault)
}

// NewServerTLSConfig creates a new server-focused tls.Config for the defined configuration
func (p *VaultTLSProvider) NewServerTLSConfig(cfg HasTLSConfig) (*tls.Config, error) {
	return NewServerTLSConfig(cfg, p.vault)
}

// SimpleTLSProvider is a simple TLS provider (not connected to Vault)
type SimpleTLSProvider struct {}

// NewSimpleTLSProvider returns a SimpleTLSProvider
func NewSimpleTLSProvider() *SimpleTLSProvider {
	return &SimpleTLSProvider{}
}

// NewClientTLSConfig creates a new client-focused tls.Config for the defined configuration
func (p *SimpleTLSProvider) NewClientTLSConfig(cfg HasTLSConfig) (*tls.Config, error) {
	return NewClientTLSConfig(cfg, &VaultConfig{})
}

// NewServerTLSConfig creates a new server-focused tls.Config for the defined configuration
func (p *SimpleTLSProvider) NewServerTLSConfig(cfg HasTLSConfig) (*tls.Config, error) {
	return NewServerTLSConfig(cfg, &VaultConfig{})
}
