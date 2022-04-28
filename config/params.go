package config

import (
	"context"

	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/service"
	"go.uber.org/fx"
)

type Definition struct {
	Name   string
	Config any
}

func Option[T any](name ...string) fx.Option {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}

	return fx.Options(
		fx.Supply(
			fx.Annotate(
				Definition{
					Name:   n,
					Config: new(T),
				},
				fx.ResultTags(`group:"configs"`),
			),
		),

		fx.Provide(func(ctx context.Context, cfg Provider) (T, error) {
			if c, err := cfg.Get(ctx, n); err != nil {
				return *new(T), err
			} else {
				return *c.(*T), nil
			}
		}),
	)
}

type Params struct {
	fx.In

	Environ env.Environ
	Prefix  service.Name
	Defs    []Definition `group:"configs"`
}
