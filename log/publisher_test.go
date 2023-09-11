package log_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.ketch.com/lib/orlop/v2/log"
	"testing"
)

type SomeEvent struct {
	ID   string
	Name string
}

func TestLog(t *testing.T) {
	ctx := context.Background()

	pub := log.NewPublisher()

	err := pub.PublishEvent(ctx, "SomeEvent", &SomeEvent{
		ID:   "123",
		Name: "one two three",
	})
	require.NoError(t, err)
}
