package internal

type Parameter interface {
	error
	Parameter() string
}
