package request

import (
	"context"
)

type HeaderSetter interface {
	Set(key, value string)
}

// CopyHeaders copies the context to the output headers
func CopyHeaders[T HeaderSetter](ctx context.Context, h T) T {
	for k, v := range Values(ctx) {
		h.Set(k, v)
	}

	return h
}
