package orlop

import (
	"github.com/spf13/pflag"
)

// HasKeyConfig denotes that the object provides Key configuration
type HasKeyConfig interface {
	GetID() string
	GetSecret() string
	GetFile() string
}

// KeyConfig provides key-related configurations
type KeyConfig struct {
	ID     string
	Secret string
	File   string
}

// GetID returns the Vault path to retrieve the key from
func (c KeyConfig) GetID() string {
	return c.ID
}

// GetSecret returns the secret (most likely from an environment variable)
func (c KeyConfig) GetSecret() string {
	return c.Secret
}

// GetFile returns the path to a file containing the key
func (c KeyConfig) GetFile() string {
	return c.File
}

// GetEnabled returns true if the key is enabled
func (c KeyConfig) GetEnabled() bool {
	return len(c.ID) > 0 || len(c.Secret) > 0 || len(c.File) > 0
}

// AddKey adds the key-related parameters
func AddKey(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	flags.String(p("id"), "", "key ID")
	flags.String(p("secret"), "", "key secret")
	AddFile(flags, "key pem file", prefix...)
}
