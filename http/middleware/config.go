package middleware

// Config provides configuration for HTTP logging
type Config struct {
	Enabled bool
	Headers []string `config:",default=X-Forwarded-For"`
}

func (c Config) GetEnabled() bool {
	return c.Enabled
}
