package health

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type NewLiveHandlerParams struct {
	fx.In

	Checks []Check `group:"livezz" optional:"true"`
}

func NewLiveHandler(params NewHealthHandlerParams) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			h := &Handler{
				checks: params.Checks,
			}

			mux.Handle("/livez", h)
			mux.Handle("/livez/{check}", h)
		},
	}
}
