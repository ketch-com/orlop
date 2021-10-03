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

package credentials

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/v2/config"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/parameter"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Tracer  trace.Tracer
	Logger  logging.Logger
	Secrets parameter.ObjectStore
}

func New(params Params) Provider {
	return &providerImpl{
		params: params,
	}
}

type providerImpl struct {
	params Params
}

// GetCredentials retrieves credentials
func (p providerImpl) GetCredentials(ctx context.Context, cfg Config) (creds Credentials, err error) {
	ctx, span := p.params.Tracer.Start(ctx, "GetCredentials")
	defer span.End()

	l := p.params.Logger.WithFields(logrus.Fields{
		"credentials.id": cfg.ID,
	})
	l.Trace("loading credentials")

	if len(cfg.Username) > 0 && len(cfg.Password) > 0 {
		l.Trace("loaded from inline settings")

		return Credentials{
			Username: cfg.Username,
			Password: cfg.Password,
		}, nil
	}

	if len(cfg.ID) == 0 || !config.IsEnabled(p.params.Secrets) {
		err = errors.New("credentials: no credentials specified")
		span.RecordError(err)
		return
	}

	err = p.params.Secrets.ReadObject(ctx, cfg.ID, &creds)
	if err != nil {
		err = errors.Wrap(err, "credentials: not found")
		span.RecordError(err)
		return
	}

	return
}
