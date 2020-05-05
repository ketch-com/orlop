package orlop

import "github.com/spf13/pflag"

// HasClientConfig denotes that an object provides client configuration
type HasClientConfig interface {
	GetTLS() HasTLSConfig
	GetToken() HasTokenConfig
	GetURL() string
}

// ClientConfig is standard configuration of most client commands
type ClientConfig struct {
	URL     string
	Token   TokenConfig
	TLS     TLSConfig
	Headers map[string]string
}

// GetURL returns the URL to contact the server
func (c ClientConfig) GetURL() string {
	return c.URL
}

// GetToken returns the security token configuration information
func (c ClientConfig) GetToken() HasTokenConfig {
	return c.Token
}

// GetTLS returns the TLS configuration
func (c ClientConfig) GetTLS() HasTLSConfig {
	return c.TLS
}

// GetHeaders returns static headers to add to requests
func (c ClientConfig) GetHeaders() map[string]string {
	return c.Headers
}

// AddClient adds the client-related parameters
func AddClient(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	AddTLS(flags, append(prefix, "tls")...)
	AddToken(flags, p("token"))
	flags.String(p("url"), "", "URL to connect to")
	flags.StringToString(p("headers"), nil, "additional headers to send")
}
