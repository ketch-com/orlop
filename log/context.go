package log

type Context interface {
	Value(key any) any
}
