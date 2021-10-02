package profiling

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.ketch.com/lib/orlop/http/routing"
	"net/http"
	"net/http/pprof"
)

func NewHandler() routing.Route {
	return func(mux chi.Router) {
		mux.MethodFunc(http.MethodGet, "/debug/pprof/", pprof.Index)
		mux.MethodFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
		mux.MethodFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
		mux.MethodFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
		mux.MethodFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)

		for _, handler := range []string{"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"} {
			mux.Method(http.MethodGet, fmt.Sprintf("/debug/pprof/%s", handler), pprof.Handler(handler))
		}
	}
}
