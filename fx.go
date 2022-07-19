// Copyright (c) 2020 Ketch Kloud, Inc.
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

package orlop

import (
	"context"
	"go.ketch.com/lib/orlop/v2/config"
	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/service"
	"go.uber.org/fx"
)

// deprecated: should not need to provide this directly anymore
func FxOptions(c any) fx.Option {
	if cfg, ok := c.(config.Config); ok {
		return cfg.Options()
	}

	return fx.Options()
}

// deprecated: should not need to provide this directly anymore
func FxContext(ctx context.Context) fx.Option {
	return fx.Provide(func() context.Context { return ctx })
}

// Populate is used for testing to populate specific entities for a unit test.
//
// deprecated: Use TestModule instead
func Populate(ctx context.Context, prefix string, _ env.Environment, module fx.Option, targets ...any) error {
	var options []fx.Option
	options = append(options, module)

	if len(targets) > 0 {
		if cfg, ok := targets[0].(config.Config); ok {
			if err := Unmarshal(prefix, cfg); err != nil {
				return err
			}

			options = append(options, cfg.Options())
			targets = targets[1:]
		}

		options = append(options, fx.Populate(targets...))
	}

	app, err := TestModule(prefix, options...)
	if app != nil {
		defer app.Stop(ctx)
	}

	return err
}

// TestModule returns an instantiated fx.App
func TestModule(prefix string, module ...fx.Option) (*fx.App, error) {
	ctx := context.Background()

	env.Test().Load()

	app := fx.New(
		fx.NopLogger,
		fx.Provide(func() context.Context { return ctx }),
		fx.Supply(service.Name(prefix)),
		fx.Supply(logging.TraceLevel),
		Module,
		fx.Options(module...),
	)

	if err := app.Err(); err != nil {
		return nil, err
	}

	if err := app.Start(ctx); err != nil {
		return nil, err
	}

	return app, app.Err()
}
