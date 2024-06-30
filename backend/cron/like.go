package cron

import (
	"blog/dao"
	"blog/dao/db"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type likeConsumerCron struct {
	internal *cron.Cron
}

func NewLikeConsumerCron() *likeConsumerCron {
	return &likeConsumerCron{
		internal: cron.New(),
	}
}

func (lcc *likeConsumerCron) Run() {
	lcc.internal.AddJob("*/2 * * * *", cron.FuncJob(func() {
		var cache = db.GetRedis()
		var err error
		var keys []string
		cachekeys := fmt.Sprintf("%s*", viper.GetString("like.cachekey"))
		keys, err = cache.Keys(context.TODO(), cachekeys).Result()
		if err != nil {
			logrus.Errorf("get keys from redis failed: %s", err.Error())
			return
		}
		var likenum uint64
		var articleidStr string
		var articleid uint64
		var found bool
		for _, key := range keys {
			likenum, err = cache.Get(context.TODO(), key).Uint64()
			if err != nil {
				continue
			}
			articleidStr, found = strings.CutPrefix(key, fmt.Sprintf("%s_", viper.GetString("like.cachekey")))
			if !found {
				continue
			}
			articleid, err = strconv.ParseUint(articleidStr, 10, 64)
			if err != nil {
				logrus.Errorf("parse articleid error (articleid:%s,likenum:%d) failed: %s", articleidStr, likenum, err.Error())
				continue
			}
			err = db.GetMysql().Model(&dao.Like{}).Where("article_id = ?", articleid).Update("like_num", likenum).Error
			if err != nil {
				logrus.Errorf("update like num (articleid:%d,likenum:%d) failed: %s", articleid, likenum, err.Error())
				continue
			}
		}
	}))
	lcc.internal.Start()
}
