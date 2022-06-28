package internal

type ErrorCode interface {
	error
	ErrorCode() string
}
