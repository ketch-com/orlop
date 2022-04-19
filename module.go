package orlop

import (
	"go.ketch.com/lib/orlop/v2/config"
	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/parameter"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	env.Module,
	logging.Module,
	parameter.Module,
)
