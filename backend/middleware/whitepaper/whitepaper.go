package whitepaper

import (
	"blog/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var whitepapre *WhitePaper

type WhitePaper struct {
	delegate *redis.Client
	key      string
}

func init() {
	whitepapre = new(WhitePaper)
	whitepapre.delegate = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	if err := whitepapre.delegate.Ping(context.TODO()).Err(); err != nil {
		logrus.Fatalf("redis connect failed, err:%v", err)
	}
	whitepapre.key = viper.GetString("whitepaper.key")
}
func WhitepaperMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.GetString("address")
		ok, err := whitepapre.delegate.SIsMember(c, whitepapre.key, address).Result()
		if err != nil || !ok {
			c.AbortWithStatusJSON(401, utils.NewFailedResponse("无权限"))
			return
		}
		c.Next()
	}
}

func ExistInWhitePaper(ctx context.Context, address string) (ok bool, err error) {
	ok, err = whitepapre.delegate.SIsMember(ctx, whitepapre.key, address).Result()
	return
}

func AddWhitePaper(ctx context.Context, operator string, address ...string) (err error) {
	var ok bool
	ok, err = whitepapre.delegate.SIsMember(ctx, whitepapre.key, operator).Result()
	if err != nil || !ok {
		if err == nil {
			err = fmt.Errorf("%s cant be allowed to operate white paper")
		}
		return
	}
	_, err = whitepapre.delegate.SAdd(ctx, whitepapre.key, address).Result()
	return
}

func DeleteWhitePaper(ctx context.Context, operator string, address ...string) (err error) {
	var ok bool
	ok, err = whitepapre.delegate.SIsMember(ctx, whitepapre.key, operator).Result()
	if err != nil || !ok {
		if err == nil {
			err = fmt.Errorf("%s cant be allowed to operate white paper")
		}
		return
	}
	_, err = whitepapre.delegate.SRem(ctx, whitepapre.key, address).Result()
	return
}
