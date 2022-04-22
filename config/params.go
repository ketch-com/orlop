package config

import (
	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/service"
	"go.uber.org/fx"
)

type Definition struct {
	Name   string
	Config any
}

func ConfigOption(name string, config any) fx.Option {
	return fx.Supply(fx.Annotate(Definition{
		Name:   name,
		Config: config,
	}, fx.ResultTags(`group:"configs"`)))
}

type Params struct {
	fx.In

	Environ env.Environ
	Prefix  service.Name
	Defs    []Definition `group:"configs"`
}
