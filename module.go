package orlop

import (
	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/parameter"
	"go.ketch.com/lib/orlop/v2/telemetry"
	"go.uber.org/fx"
)

var Module = fx.Options(
	env.Module,
	logging.Module,
	parameter.Module,
	telemetry.Module,
)
