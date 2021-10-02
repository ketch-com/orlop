package key

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewPrivateKey,
		NewPublicKey,
	),
)
