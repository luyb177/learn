package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewGrabServiceImpl,
	NewMonitorServiceImpl,
	NewSseServiceImpl,
)
