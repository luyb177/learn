package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccount,
	NewRedisConfig,
	NewKafkaConfig,
	NewViperSetting,
	NewMysqlConfig,
)
