package middleware

import "go.uber.org/fx"

var Modules = fx.Options(
	fx.Provide(
		CORS,
		Logging,
		Metrics,
	),
)
