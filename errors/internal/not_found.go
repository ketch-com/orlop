package internal

type NotFound interface {
	error
	NotFound() bool
}
