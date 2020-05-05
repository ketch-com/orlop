package orlop

import (
	"github.com/spf13/pflag"
	"time"
)

// HasCertGenerationConfig denotes that an object provides certificate generation configuration.
type HasCertGenerationConfig interface {
	GetEnabled() bool
	GetPath() string
	GetCommonName() string
	GetAltNames() string
	GetTTL() time.Duration
}

// CertGenerationConfig provides the certificate generation configuration
type CertGenerationConfig struct {
	Enabled    bool
	Path       string
	CommonName string
	AltNames   string
	TTL        time.Duration
}

// GetEnabled returns true if certificate generation is enabled
func (c CertGenerationConfig) GetEnabled() bool {
	return c.Enabled && len(c.CommonName) > 0
}

// GetPath returns the path to the certificate generator
func (c CertGenerationConfig) GetPath() string {
	return c.Path
}

// GetCommonName returns the common name to use for certificates
func (c CertGenerationConfig) GetCommonName() string {
	return c.CommonName
}

// GetAltNames returns alternate names to use for certificates
func (c CertGenerationConfig) GetAltNames() string {
	return c.AltNames
}

// GetTTL returns the time-to-live for the certificates
func (c CertGenerationConfig) GetTTL() time.Duration {
	return c.TTL
}

// AddCertGenerationConfig adds certificate generation configuration
func AddCertGenerationConfig(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	AddEnabled(flags, "certificate generation enabled", false, prefix...)
	flags.String(p("path"), "/pki/issue/", "path to the issuer")
	flags.String(p("common-name"), "", "common name for the certificate")
	flags.String(p("alt-names"), "", "alt names for the certificate")
	flags.Duration(p("ttl"), 0*time.Hour, "time-to-live for generated certificates")
}
