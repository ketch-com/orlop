package health

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type NewHealthHandlerParams struct {
	fx.In

	Checks []Check `group:"healthz" optional:"true"`
}

func NewHealthHandler(params NewHealthHandlerParams) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			h := &Handler{
				checks: params.Checks,
			}

			mux.Handle("/healthz", h)
			mux.Handle("/healthz/{check}", h)
		},
	}
}
