package orlop

import "github.com/spf13/pflag"

// HasEnabled denotes an object provides an Enabled flag
type HasEnabled interface {
	GetEnabled() bool
}

// Enabled provides an Enabled flag
type Enabled struct {
	Enabled bool
}

// GetEnabled returns true if enabled
func (c Enabled) GetEnabled() bool {
	return c.Enabled
}

// AddEnabled adds Enabled flags
func AddEnabled(flags *pflag.FlagSet, usage string, value bool, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	flags.Bool(p("enabled"), value, usage)
}
