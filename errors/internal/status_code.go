package internal

type StatusCode interface {
	error
	StatusCode() int
}
