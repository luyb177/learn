package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccount,
	NewRedisConfig,
	NewViperSetting,
	NewQQConfig,
	NewMysqlConfig,
)
