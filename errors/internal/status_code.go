package internal

type StatusCode interface {
	error
	StatusCode() int
}

type Status interface {
	error
	Status() int
}
