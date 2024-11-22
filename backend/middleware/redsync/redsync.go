package disync

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var rs *redsync.Redsync

func init() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     viper.GetString("redis.addr"),
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	pool := goredis.NewPool(client)
	rs = redsync.New(pool)
	logrus.Info("init the distribute lock success")
}

func NewMutex(name string) *redsync.Mutex {
	return rs.NewMutex(name)
}
