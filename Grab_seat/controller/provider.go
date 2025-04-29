package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewGrabController,
	NewMonitorController,
	NewSseController,
)
