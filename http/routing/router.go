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
