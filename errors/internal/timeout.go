package internal

type Timeout interface {
	error
	Timeout() bool
}
