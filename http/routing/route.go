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
	"go.uber.org/fx"
	"net/http"
)

type Route func(mux chi.Router)

func Method(method string, path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(method, path, handler)
		},
	}
}

func ANY(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Handle(path, handler)
		},
	}
}

func GET(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodGet, path, handler)
		},
	}
}

func HEAD(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodHead, path, handler)
		},
	}
}

func POST(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodPost, path, handler)
		},
	}
}

func PUT(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodPut, path, handler)
		},
	}
}

func PATCH(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodPatch, path, handler)
		},
	}
}

func DELETE(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodDelete, path, handler)
		},
	}
}

func CONNECT(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodConnect, path, handler)
		},
	}
}

func OPTIONS(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodOptions, path, handler)
		},
	}
}

func TRACE(path string, handler http.Handler) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.Method(http.MethodTrace, path, handler)
		},
	}
}
