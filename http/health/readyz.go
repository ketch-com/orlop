package health

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type NewReadyHandlerParams struct {
	fx.In

	Checks []Check `group:"readyz" optional:"true"`
}

func NewReadyHandler(params NewHealthHandlerParams) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			h := &Handler{
				checks: params.Checks,
			}

			mux.Handle("/readyz", h)
			mux.Handle("/readyz/{check}", h)
		},
	}
}
