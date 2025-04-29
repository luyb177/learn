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
type KafkaConfig struct {
	Addr string `yaml:"addr"`
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

func NewKafkaConfig(vs *ViperSetting) *KafkaConfig {
	var kafkaConfig = &KafkaConfig{}
	vs.ReadSection("kafka", &kafkaConfig)
	return kafkaConfig
}

func NewMysqlConfig(vs *ViperSetting) *MysqlConfig {
	var mysqlConfig = &MysqlConfig{}
	vs.ReadSection("mysql", &mysqlConfig)
	return mysqlConfig
}
