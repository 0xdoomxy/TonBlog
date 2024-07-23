package db

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 加载配置文件
func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic("load config failed:", err.Error())
	}
}
