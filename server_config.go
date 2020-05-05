package orlop

import (
	"github.com/spf13/pflag"
)

// HasServerConfig denotes an object provides server configuration
type HasServerConfig interface {
	GetBind() string
	GetListen() uint
	GetTLS() HasTLSConfig
	GetSwagger() HasEnabled
}

// ServerConfig is standard configuration of most server commands
type ServerConfig struct {
	Bind    string
	Listen  uint
	TLS     TLSConfig
	Swagger Enabled
}

// GetBind returns the address to bind to
func (c ServerConfig) GetBind() string {
	return c.Bind
}

// GetListen returns the port to listen on
func (c ServerConfig) GetListen() uint {
	return c.Listen
}

// GetTLS returns TLS configuration
func (c ServerConfig) GetTLS() HasTLSConfig {
	return c.TLS
}

// GetSwagger returns Enabled for Swagger
func (c ServerConfig) GetSwagger() HasEnabled {
	return c.Swagger
}

// NewServerConfig returns a new, unmarshaled standard ServerConfig
func NewServerConfig() (*ServerConfig, error) {
	c := new(ServerConfig)

	err := Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// AddServer adds the server-related parameters
func AddServer(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	AddTLS(flags, append(prefix, "tls")...)
	AddEnabled(flags, "enable Swagger", false, append(prefix, "swagger")...)
	flags.String(p("bind"), "0.0.0.0", "address to bind to")
	flags.Uint(p("listen"), 5000, "port to bind to")
}
