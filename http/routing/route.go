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
