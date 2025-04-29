package config

type AccountConfig struct {
	Username string
	Password string
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type QQConfig struct {
	Email string `yaml:"email"`
	Key   string `yaml:"key"`
}
type MysqlConfig struct {
	Addr string `yaml:"addr"`
}

func NewAccount(vs *ViperSetting) *AccountConfig {
	var account = &AccountConfig{}
	vs.ReadSection("account", &account)
	return account
}

func NewRedisConfig(vs *ViperSetting) *RedisConfig {
	var redisConfig = &RedisConfig{}
	vs.ReadSection("redis", &redisConfig)
	return redisConfig
}

func NewQQConfig(vs *ViperSetting) *QQConfig {
	var qqConfig = &QQConfig{}
	vs.ReadSection("qq", &qqConfig)
	return qqConfig
}

func NewMysqlConfig(vs *ViperSetting) *MysqlConfig {
	var mysqlConfig = &MysqlConfig{}
	vs.ReadSection("mysql", &mysqlConfig)
	return mysqlConfig
}
