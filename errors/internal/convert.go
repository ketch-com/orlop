package internal

// Converter converts the given error to a standard. If a conversion happened, then ok is true.
type Converter func(err error) (error, bool)

var Converters []Converter

func RegisterConverter(converter func(err error) (error, bool)) {
	Converters = append(Converters, converter)
}
