package orlop

import (
	"crypto/tls"
	"github.com/spf13/pflag"
)

// HasTLSConfig denotes that an object supports TLS configuration
type HasTLSConfig interface {
	GetCert() KeyConfig
	GetClientAuth() tls.ClientAuthType
	GetEnabled() bool
	GetInsecure() bool
	GetKey() KeyConfig
	GetOverride() string
	GetRootCA() KeyConfig
	GetGenerate() CertGenerationConfig
}

// TLSConfig provides TLS configuration
type TLSConfig struct {
	ClientAuth tls.ClientAuthType
	Enabled    bool
	Insecure   bool
	Override   string
	Cert       KeyConfig
	Key        KeyConfig
	RootCA     KeyConfig
	Generate   CertGenerationConfig
}

// GetCert returns the certificate key configuration
func (t TLSConfig) GetCert() KeyConfig {
	return t.Cert
}

// GetClientAuth returns the Client Authentication Type
func (t TLSConfig) GetClientAuth() tls.ClientAuthType {
	return t.ClientAuth
}

// GetEnabled returns true of TLS is enabled
func (t TLSConfig) GetEnabled() bool {
	return t.Enabled
}

// GetInsecure returns true if insecure mode is enabled (use at your own peril)
func (t TLSConfig) GetInsecure() bool {
	return t.Insecure
}

// GetKey returns the private key configuration
func (t TLSConfig) GetKey() KeyConfig {
	return t.Key
}

// GetOverride returns an override name for the server
func (t TLSConfig) GetOverride() string {
	return t.Override
}

// GetRootCA returns the root CA key configuration
func (t TLSConfig) GetRootCA() KeyConfig {
	return t.RootCA
}

// GetGenerate returns the certificate generation configuration
func (t TLSConfig) GetGenerate() CertGenerationConfig {
	return t.Generate
}

// CloneTLSConfig clones the given TLS configuration
func CloneTLSConfig(cfg HasTLSConfig) TLSConfig {
	return TLSConfig{
		ClientAuth: cfg.GetClientAuth(),
		Enabled:    cfg.GetEnabled(),
		Insecure:   cfg.GetInsecure(),
		Override:   cfg.GetOverride(),
		Cert:       cfg.GetCert(),
		Key:        cfg.GetKey(),
		RootCA:     cfg.GetRootCA(),
		Generate:   cfg.GetGenerate(),
	}
}

// AddTLS adds the TLS-related parameters
func AddTLS(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	AddEnabled(flags, "TLS enabled", true, prefix...)
	AddCertGenerationConfig(flags, append(prefix, "generate")...)
	AddKey(flags, append(prefix, "cert")...)
	AddKey(flags, append(prefix, "key")...)
	AddKey(flags, append(prefix, "rootca")...)
	flags.String(p("override"), "", "server override")
	flags.Bool(p("insecure"), false, "skip verifying insecure certificates")
	flags.Var(newClientAuthValue(), p("clientauth"), "client authentication mode")
}
