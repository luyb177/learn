package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ViperSetting struct {
	*viper.Viper
}

func (vs *ViperSetting) ReadSection(k string, v interface{}) error {
	err := vs.UnmarshalKey(k, v) // 将配置文件中某个部分的数据解码到指定变量中
	if err != nil {
		return err
	}
	return nil
}

func NewViperSetting(ConfigPath string) *ViperSetting {
	vp := viper.New()            // 创建实例
	vp.SetConfigFile(ConfigPath) // 配置文件路径
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}
	return &ViperSetting{vp}
}
