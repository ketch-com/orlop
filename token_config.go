package orlop

import (
	"github.com/spf13/pflag"
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
	KeyMap     KeyConfig
	PrivateKey KeyConfig
	PublicKey  KeyConfig
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

// AddToken adds token parameters to the FlagSet
func AddToken(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	flags.Duration(p("ttl"), 24*time.Hour, "time-to-live for token")
	flags.String(p("issuer"), "", "token issuer value")
	AddKey(flags, append(prefix, "keymap")...)
	AddKey(flags, append(prefix, "privatekey")...)
	AddKey(flags, append(prefix, "publickey")...)
	AddKey(flags, append(prefix, "shared")...)
}
