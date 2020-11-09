package orlop

import (
	"github.com/switch-bit/orlop/version"
	"go.opentelemetry.io/otel/api/global"
)

var tracer = global.Tracer(version.Name)
