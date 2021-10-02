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

package server

import (
	"context"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/logging"
	"go.uber.org/fx"
	"net"
	"net/http"
)

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    logging.Logger
	Listener  net.Listener
	Server    *http.Server
}

func Lifecycle(params Params) {
	s := &serverImpl{
		params: params,
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
}

type serverImpl struct {
	params Params
}

func (s *serverImpl) Start(ctx context.Context) error {
	go s.Run(ctx)
	return nil
}

func (s *serverImpl) Stop(ctx context.Context) error {
	if err := s.params.Server.Shutdown(ctx); err != nil {
		return err
	}

	if err := s.params.Listener.Close(); err != nil {
		return err
	}

	return nil
}

func (s *serverImpl) Run(ctx context.Context) error {
	s.params.Logger.Info("serving")

	if err := s.params.Server.Serve(s.params.Listener); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}
