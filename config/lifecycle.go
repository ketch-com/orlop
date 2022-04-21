package config

import (
	"context"

	"go.uber.org/fx"
)

func Lifecycle(lc fx.Lifecycle, provider Provider) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return provider.Load(ctx)
			},
		},
	)
}
