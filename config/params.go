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

func Option[T any](name ...string) fx.Option {
	var n, annotation string
	if len(name) > 0 {
		n = name[0]
		if len(name) > 1 {
			annotation = fmt.Sprintf(`name:"%s"`, name[1])
		}
	}
	fn := func(ctx context.Context, cfg Provider) (T, error) {
		if c, err := cfg.Get(ctx, n); err != nil {
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
					Name:   n,
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
