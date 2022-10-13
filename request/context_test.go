package request

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCloneContext(t *testing.T) {
	ctx := context.Background()

	ctx = WithTenant(ctx, "ani")
	ctx = WithID(ctx, "myID")
	ctx, cancelFn := context.WithTimeout(ctx, 5*time.Second)

	newCtx := Clone(ctx)

	assert.Equal(t, Tenant(ctx), Tenant(newCtx))
	assert.Equal(t, ID(ctx), ID(newCtx))
	assert.Equal(t, Operation(ctx), Operation(newCtx))
	assert.Equal(t, URL(ctx), URL(newCtx))
	assert.Equal(t, Originator(ctx), Originator(newCtx))

	d, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.NotNil(t, d)
	cancelFn()
	assert.NotNil(t, ctx.Err())
	_, ok = newCtx.Deadline()
	assert.False(t, ok)
	assert.Nil(t, newCtx.Err())
}
