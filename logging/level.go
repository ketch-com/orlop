package logging

type Level string
const (
	TraceLevel Level = "trace"
	DebugLevel Level = "debug"
	InfoLevel Level = "info"
	WarnLevel Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)
