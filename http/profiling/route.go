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

package profiling

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"net/http"
	"net/http/pprof"
)

func NewHandler() fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			mux.MethodFunc(http.MethodGet, "/debug/pprof/", pprof.Index)
			mux.MethodFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
			mux.MethodFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
			mux.MethodFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
			mux.MethodFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)

			for _, handler := range []string{"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"} {
				mux.Method(http.MethodGet, fmt.Sprintf("/debug/pprof/%s", handler), pprof.Handler(handler))
			}
		},
	}
}
