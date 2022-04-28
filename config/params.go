package config

import (
	"context"
	"fmt"

	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/service"
	"go.uber.org/fx"
)

type Definition struct {
	Name   string
	Config any
}

func Option[T any](in ...string) fx.Option {
	var name, annotation string
	if len(in) > 0 {
		name = in[0]
		if len(in) > 1 {
			annotation = fmt.Sprintf(`name:"%s"`, in[1])
		}
	}
	fn := func(ctx context.Context, cfg Provider) (T, error) {
		if c, err := cfg.Get(ctx, name); err != nil {
			return *new(T), err
		} else {
			return *c.(*T), nil
		}
	}

	p := fx.Provide(fn)
	if len(annotation) > 0 {
		p = fx.Provide(
			fx.Annotate(
				fn,
				fx.ResultTags(annotation),
			),
		)
	}

	return fx.Options(
		fx.Supply(
			fx.Annotate(
				Definition{
					Name:   name,
					Config: new(T),
				},
				fx.ResultTags(`group:"configs"`),
			),
		),
		p,
	)
}

type Params struct {
	fx.In

	Environ env.Environ
	Prefix  service.Name
	Defs    []Definition `group:"configs"`
}
