package route

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewApp,
	NewGrabRoot,
	NewMonitorRoot,
	NewSseRoute,
)
