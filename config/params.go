package config

import "go.uber.org/fx"

type Params struct {
	fx.In

	Defs []Definition `optional:"true"`
}
