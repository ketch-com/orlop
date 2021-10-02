// Copyright (c) 2021 Ketch Kloud, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
