package orlop

import "github.com/spf13/pflag"

// HasFile denotes an object provides a filename
type HasFile interface {
	GetFile() string
}

// File provides a filename
type File struct {
	File string
}

// GetFile returns the filename
func (f File) GetFile() string {
	return f.File
}

// AddFile adds the file configuration parameters
func AddFile(flags *pflag.FlagSet, usage string, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	flags.String(p("file"), "", usage)
}
