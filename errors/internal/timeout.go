package internal

import "time"

type Timeout interface {
	error
	Timeout() time.Duration
}
