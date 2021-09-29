package orlop

import (
	"context"
	"go.uber.org/fx"
)

type ProvidesFxOptions interface {
	Options() fx.Option
}

func FxOptions(o ProvidesFxOptions) fx.Option {
	return o.Options()
}

func FxContext(ctx context.Context) fx.Option {
	return fx.Provide(func() context.Context { return ctx })
}
