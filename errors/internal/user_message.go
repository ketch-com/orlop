package internal

type UserMessage interface {
	error
	UserMessage() string
}
