package orlop

import (
	"github.com/spf13/pflag"
)

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

// AddVault adds the Vault parameters to the FlagSet
func AddVault(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	AddTLS(flags, append(prefix, "tls")...)
	AddEnabled(flags, "Vault enabled", false, prefix...)
	flags.String(p("address"), "", "Vault address")
	flags.String(p("token"), "", "Vault token")
	flags.String(p("prefix"), "", "Vault prefix")
}
