package pem

import (
	"context"
	"go.ketch.com/lib/orlop/secret"
	"go.uber.org/fx"
)

type Provider interface {
	// Load a PEM
	Load(ctx context.Context, cfg Config) ([]byte, error)
}

type Params struct {
	fx.In

	Secrets secret.Provider
}

func New(params Params) Provider {
	return &implProvider{
		params: params,
	}
}

type implProvider struct {
	params Params
}

// Load the PEM bytes based on the config
func (p *implProvider) Load(ctx context.Context, cfg Config) ([]byte, error) {
	s, err := p.params.Secrets.Load(ctx, secret.Config{
		ID:    cfg.GetID(),
		File:  cfg.GetFile(),
		Which: cfg.GetWhich(),
	})
	if err != nil {
		return nil, err
	}

	return []byte(s), nil
}
