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

package routing

import (
	"github.com/go-chi/chi/v5"
	"go.ketch.com/lib/orlop/http/middleware"
	"go.uber.org/fx"
	"net/http"
)

type Params struct {
	fx.In

	Middlewares      []middleware.Middleware `group:"middleware,flatten"`
	Routes           []Route                 `group:"routes,flatten"`
	MethodNotAllowed http.HandlerFunc        `name:"methodNotAllowed" optional:"true"`
	NotFound         http.HandlerFunc        `name:"notFound" optional:"true"`
}

func BuildRouter(params Params) chi.Router {
	var mux chi.Router = chi.NewMux()

	for _, route := range params.Routes {
		route(mux)
	}

	for _, mw := range params.Middlewares {
		mux.Use(mw)
	}

	if params.MethodNotAllowed != nil {
		mux.MethodNotAllowed(params.MethodNotAllowed)
	}

	if params.NotFound != nil {
		mux.NotFound(params.NotFound)
	}

	return mux
}

func MuxToRouter(mux *chi.Mux) chi.Router {
	return mux
}

func RouterToHandler(mux chi.Router) http.Handler {
	return mux
}

func HandlerFuncToHandler(h http.HandlerFunc) http.Handler {
	return h
}

func MethodNotAllowed(h http.HandlerFunc) fx.Option {
	return fx.Supply(
		fx.Annotated{
			Name:   "methodNotAllowed",
			Target: h,
		},
	)
}

func NotFound(h http.HandlerFunc) fx.Option {
	return fx.Supply(
		fx.Annotated{
			Name:   "notFound",
			Target: h,
		},
	)
}
