package internal

type Temporary interface {
	error
	Temporary() bool
}
