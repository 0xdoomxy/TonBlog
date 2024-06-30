package db

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var rs *redis.Client

func init() {
	rs = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	if err := rs.Ping(context.TODO()).Err(); err != nil {
		logrus.Fatalf("redis connect failed, err:%v", err)
	}
}

func GetRedis() *redis.Client {
	return rs
}
