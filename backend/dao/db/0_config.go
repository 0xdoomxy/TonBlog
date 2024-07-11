package db

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 加载配置文件
func init() {
	go func() {
		logrus.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic("load config failed:", err.Error())
	}
}
