package internal

type Source interface {
	error
	Source() string
}
